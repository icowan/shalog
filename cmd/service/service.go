package service

import (
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	kitlogrus "github.com/go-kit/kit/log/logrus"
	"github.com/go-kit/kit/metrics/prometheus"
	"github.com/icowan/blog/src/cmd"
	"github.com/icowan/blog/src/config"
	"github.com/icowan/blog/src/logging"
	"github.com/icowan/blog/src/mysql"
	"github.com/icowan/blog/src/pkg/about"
	"github.com/icowan/blog/src/pkg/account"
	"github.com/icowan/blog/src/pkg/api"
	"github.com/icowan/blog/src/pkg/board"
	"github.com/icowan/blog/src/pkg/home"
	"github.com/icowan/blog/src/pkg/post"
	"github.com/icowan/blog/src/pkg/reward"
	"github.com/icowan/blog/src/repository"
	"github.com/jinzhu/gorm"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

const (
	DefaultHttpPort   = ":8080"
	DefaultConfigPath = "./app.cfg"
	DefaultStaticPath = "./static/"
)

var (
	httpAddr   = envString("HTTP_ADDR", DefaultHttpPort)
	configPath = envString("CONFIG_PATH", DefaultConfigPath)
	staticPath = envString("STATIC_PATH", DefaultStaticPath)

	rootCmd = &cobra.Command{
		Use:               "server",
		Short:             "开普勒平台服务端",
		SilenceErrors:     true,
		DisableAutoGenTag: true,
		Long: `# 开普勒平台服务端
您可以通过改命令来启动您的服务
可用的配置类型：
[start]
有关开普勒平台的相关概述，请参阅 https://github.com/nsini/kplcloud
`,
	}

	startCmd = &cobra.Command{
		Use:   "start",
		Short: "启动服务",
		Example: `## 启动命令
blog start -p :8080 -c ./app.cfg
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			start()
			return nil
		},
	}

	logger     log.Logger
	store      repository.Repository
	db         *gorm.DB
	adminEmail string
)

func init() {

	rootCmd.PersistentFlags().StringVarP(&httpAddr, "http.port", "p", DefaultHttpPort, "服务启动的端口: :8080")
	rootCmd.PersistentFlags().StringVarP(&configPath, "config.path", "c", DefaultConfigPath, "配置文件路径: ./app.yaml")
	startCmd.PersistentFlags().StringVarP(&staticPath, "static.path", "s", DefaultStaticPath, "静态文件目录: ./static/")

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

	if cf.GetString("server", "logs_path") != "" {
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
	var boardSvc board.Service
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

	// board
	boardSvc = board.NewService(logger)

	// admin account
	accountSvc := account.NewService(logger, store, cf)
	accountSvc = account.NewLoggingServer(logger, accountSvc)

	httpLogger := log.With(logger, "component", "http")

	mux := http.NewServeMux()

	mux.Handle("/account/", account.MakeHTTPHandler(accountSvc, httpLogger, store))
	mux.Handle("/search", post.MakeHandler(ps, httpLogger, store))
	mux.Handle("/post", post.MakeHandler(ps, httpLogger, store))
	mux.Handle("/post/", post.MakeHandler(ps, httpLogger, store))
	mux.Handle("/about", about.MakeHandler(aboutMe, httpLogger))
	mux.Handle("/api/", api.MakeHandler(apiSvc, httpLogger, store, cf))
	mux.Handle("/board", board.MakeHandler(boardSvc, httpLogger))
	mux.Handle("/reward", reward.MakeHandler(httpLogger))
	mux.Handle("/", home.MakeHandler(homeSvc, httpLogger))

	http.Handle("/metrics", promhttp.Handler())
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("./views/tonight/images/"))))
	http.Handle("/fonts/", http.StripPrefix("/fonts/", http.FileServer(http.Dir("./views/tonight/fonts/"))))
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./views/tonight/css/"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("./views/tonight/js/"))))
	http.Handle("/", accessControl(mux, logger))

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

func accessControl(h http.Handler, logger log.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

		if r.Method == "OPTIONS" {
			return
		}

		_ = level.Info(logger).Log("remote-addr", r.RemoteAddr, "uri", r.RequestURI, "method", r.Method, "length", r.ContentLength)

		h.ServeHTTP(w, r)
	})
}
