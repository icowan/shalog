package service

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	kitlogrus "github.com/go-kit/kit/log/logrus"
	"github.com/go-kit/kit/metrics/prometheus"
	"github.com/icowan/shalog/src/cmd"
	"github.com/icowan/shalog/src/config"
	"github.com/icowan/shalog/src/logging"
	"github.com/icowan/shalog/src/mysql"
	"github.com/icowan/shalog/src/pkg/about"
	"github.com/icowan/shalog/src/pkg/account"
	"github.com/icowan/shalog/src/pkg/api"
	"github.com/icowan/shalog/src/pkg/category"
	"github.com/icowan/shalog/src/pkg/home"
	"github.com/icowan/shalog/src/pkg/image"
	"github.com/icowan/shalog/src/pkg/link"
	"github.com/icowan/shalog/src/pkg/post"
	"github.com/icowan/shalog/src/pkg/setting"
	"github.com/icowan/shalog/src/pkg/tag"
	"github.com/icowan/shalog/src/repository"
	"github.com/jinzhu/gorm"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	DefaultHttpPort   = ":8080"
	DefaultConfigPath = "./app.cfg"
	//DefaultStaticPath = "./static/"
	AdminViewPath   = "./views/admin/"
	DefaultUsername = "root"
	DefaultPassword = "admin"
	DefaultSQL      = "./database/db.sql"
	DefaultImage    = "local"
)

var (
	httpAddr   = envString("HTTP_ADDR", DefaultHttpPort)
	configPath = envString("CONFIG_PATH", DefaultConfigPath)
	username   = envString("USERNAME", DefaultUsername)
	password   = envString("PASSWORD", DefaultPassword)
	sqlPath    = envString("SQL_PATH", DefaultSQL)
	imagePath  = envString("IMAGE_PATH", DefaultImage)

	appKey = ""

	rootCmd = &cobra.Command{
		Use:               "server",
		Short:             "",
		SilenceErrors:     true,
		DisableAutoGenTag: true,
		Long: `# 博客系统
可用的配置类型：
[start]
有关本系统的相关概述，请参阅 https://github.com/icowan/shalog
`,
	}

	startCmd = &cobra.Command{
		Use:   "start",
		Short: "启动服务",
		Example: `## 启动命令
shalog start -p :8080 -c ./app.cfg
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			start()
			return nil
		},
	}

	logger log.Logger
	store  repository.Repository
	db     *gorm.DB
)

func init() {

	rootCmd.PersistentFlags().StringVarP(&httpAddr, "http.port", "p", DefaultHttpPort, "服务启动的端口: :8080")
	rootCmd.PersistentFlags().StringVarP(&configPath, "config.path", "c", DefaultConfigPath, "配置文件路径: ./app.yaml")
	//rootCmd.PersistentFlags().StringVarP(&username, "username", "u", DefaultUsername, "初始化用户名")
	//rootCmd.PersistentFlags().StringVarP(&password, "password", "P", DefaultPassword, "初始化密码")
	rootCmd.PersistentFlags().StringVarP(&sqlPath, "sql.path", "s", DefaultSQL, "初始化数据库SQL文件")
	startCmd.PersistentFlags().StringVarP(&imagePath, "static.path", "i", DefaultImage, "是否使用本地资源对图片进行处理: local | remote")

	cmd.AddFlags(rootCmd)
	rootCmd.AddCommand(startCmd)
}

func Run() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}

func start() {

	cf, err := config.NewConfig(configPath)
	if err != nil {
		panic(err)
	}

	appKey = cf.GetString(config.SectionServer, "app_key")

	if cf.GetString(config.SectionServer, "logs_path") != "" {
		logrusLogger, err := logging.LogrusLogger(cf.GetString("server", "logs_path"))
		if err != nil {
			panic(err)
		}
		logger = kitlogrus.NewLogrusLogger(logrusLogger)
	} else {
		logger = log.NewLogfmtLogger(log.StdlibWriter{})
	}
	logger = log.With(logger, "caller", log.DefaultCaller)
	logger = level.NewFilter(logger, logLevel(cf.GetString("server", "log_level")))

	db, err = mysql.NewDb(logger, cf)
	if err != nil {
		_ = level.Error(logger).Log("db", "connect", "err", err)
		os.Exit(1)
	}
	defer func() {
		if err = db.Close(); err != nil {
			panic(err)
		}
	}()

	store = repository.NewRepository(db)

	if err = importToDb(); err != nil {
		_ = level.Error(logger).Log("import", "db", "err", err.Error())
		os.Exit(1)
	}

	fieldKeys := []string{"method"}

	var ps post.Service
	var aboutMe about.Service
	var homeSvc home.Service
	var apiSvc api.Service

	// post
	ps = post.NewService(logger, cf, store)
	ps = post.NewLoggingService(logger, ps)

	// api
	apiSvc = api.NewService(logger, cf, store)
	apiSvc = api.NewLoggingService(logger, apiSvc)
	apiSvc = api.NewInstrumentingService(
		prometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "api",
			Subsystem: "post_service",
			Name:      "request_count",
			Help:      "Number of requests received.",
		}, fieldKeys),
		prometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "api",
			Subsystem: "post_service",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, fieldKeys),
		apiSvc,
	)

	// home
	homeSvc = home.NewService(logger, cf, store)
	homeSvc = home.NewLoggingService(logger, homeSvc)

	// about
	aboutMe = about.NewService(logger)
	aboutMe = about.NewLoggingService(logger, aboutMe)

	// link
	linkSvc := link.NewService(logger, store)
	linkSvc = link.NewLoggingServer(logger, linkSvc)

	// admin account
	accountSvc := account.NewService(logger, store, cf)
	accountSvc = account.NewLoggingServer(logger, accountSvc)

	// setting
	settingSvc := setting.NewService(logger, store, cf)
	settingSvc = setting.NewLoggingServer(logger, settingSvc)

	// category
	categorySvc := category.NewService(logger, store)
	categorySvc = category.NewLoggingServer(logger, categorySvc)

	// tag
	tagSvc := tag.NewService(logger, store)
	tagSvc = tag.NewLoggingServer(logger, tagSvc)

	// image
	imageSvc := image.NewService(logger, store, cf)
	imageSvc = image.NewLoggingServer(logger, imageSvc)

	settings, err := settingSvc.List(context.Background())
	if err != nil {
		_ = level.Error(logger).Log("SettingSvc", "List", "err", err.Error())
		return
	}

	sets := map[string]string{}

	for _, v := range settings {
		sets[v.Key] = v.Value
		cf.SetValue(config.SectionServer, v.Key, v.Value)
	}

	sets["service-start-time"] = time.Now().Format("20060102150405")

	httpLogger := log.With(logger, "component", "http")

	mux := http.NewServeMux()

	mux.Handle("/search", post.MakeHandler(ps, httpLogger, store, sets))
	mux.Handle("/post", post.MakeHandler(ps, httpLogger, store, sets))
	mux.Handle("/post/", post.MakeHandler(ps, httpLogger, store, sets))
	mux.Handle("/about", about.MakeHandler(aboutMe, httpLogger, sets))
	mux.Handle("/link/", link.MakeHTTPHandler(linkSvc, httpLogger))

	mux.Handle("/account/", account.MakeHTTPHandler(accountSvc, httpLogger))
	mux.Handle("/setting", setting.MakeHTTPHandler(settingSvc, httpLogger))
	mux.Handle("/setting/", setting.MakeHTTPHandler(settingSvc, httpLogger))
	mux.Handle("/category/", category.MakeHTTPHandler(categorySvc, httpLogger))
	mux.Handle("/tag/", tag.MakeHTTPHandler(tagSvc, httpLogger))
	mux.Handle("/image/", image.MakeHTTPHandler(imageSvc, httpLogger, sets))

	mux.Handle("/api/", api.MakeHandler(apiSvc, httpLogger, store, cf))
	//mux.Handle("/board", board.MakeHandler(boardSvc, httpLogger))
	mux.Handle("/", home.MakeHandler(homeSvc, httpLogger, sets))

	viewsPath := sets[repository.SettingViewTemplate.String()]

	http.Handle("/metrics", promhttp.Handler())
	http.Handle("/image/", http.StripPrefix("/image/", http.FileServer(http.Dir(viewsPath+"/image/"))))
	http.Handle("/fonts/", http.StripPrefix("/fonts/", http.FileServer(http.Dir(viewsPath+"/fonts/"))))
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir(viewsPath+"/css/"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir(viewsPath+"/js/"))))
	http.Handle("/admin/", http.StripPrefix("/admin/", http.FileServer(http.Dir(AdminViewPath))))

	imageDomain := sets[repository.SettingGlobalDomainImage.String()]
	u, _ := url.Parse(imageDomain)

	if imagePath == DefaultImage {
		http.Handle(u.Path+"/", image.MakeHTTPHandler(imageSvc, logger, sets))
	} else {
		storagePath := sets[repository.SettingSiteMediaUploadPath.String()]
		http.Handle(u.Path+"/", http.StripPrefix(u.Path+"/", http.FileServer(http.Dir(storagePath))))
	}

	handlers := make(map[string]string, 3)
	if cf.GetBool("cors", "allow") {
		handlers["Access-Control-Allow-Origin"] = cf.GetString("cors", "origin")
		handlers["Access-Control-Allow-Methods"] = cf.GetString("cors", "methods")
		handlers["Access-Control-Allow-Headers"] = cf.GetString("cors", "headers")
	}
	http.Handle("/", accessControl(mux, logger, handlers))

	errs := make(chan error, 2)
	go func() {
		_ = level.Debug(logger).Log("transport", "http", "address", httpAddr, "msg", "listening")
		errs <- http.ListenAndServe(httpAddr, nil)
	}()
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	_ = level.Error(logger).Log("terminated", <-errs)
}

func envString(env, fallback string) string {
	e := os.Getenv(env)
	if e == "" {
		return fallback
	}
	return e
}

func accessControl(h http.Handler, logger log.Logger, headers map[string]string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for key, val := range headers {
			w.Header().Set(key, val)
		}
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Connection", "keep-alive")

		if r.Method == "OPTIONS" {
			return
		}

		_ = level.Info(logger).Log("remote-addr", r.RemoteAddr, "uri", r.RequestURI, "method", r.Method, "length", r.ContentLength)

		h.ServeHTTP(w, r)
	})
}
