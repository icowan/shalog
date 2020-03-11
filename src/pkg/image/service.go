package image

import (
	"context"
	"crypto/md5"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/icowan/shalog/src/config"
	"github.com/icowan/shalog/src/repository"
	"github.com/icowan/shalog/src/repository/types"
	file2 "github.com/icowan/shalog/src/util/file"
	"github.com/nfnt/resize"
	"github.com/pkg/errors"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	ErrFileNotFound = errors.New("图片不存在.")
	ErrFileParams   = errors.New("处理参数错误")
)

type Service interface {
	List(ctx context.Context, pageSize, offset int) (images []types.Image, count int64, err error)
	UploadMedia(ctx context.Context, f *multipart.FileHeader) (resImg *imageResponse, err error)

	// 实时渲染图片
	Get(ctx context.Context, req *url.URL) ([]byte, error)
}

type service struct {
	logger     log.Logger
	repository repository.Repository
	config     *config.Config
}

// 没有参数直接展示
// 判断有没有切好的，有就直接输出
// 如果没有切好的，切完存本地 然后输出
func (s *service) Get(ctx context.Context, u *url.URL) (res []byte, err error) {
	setting := ctx.Value("settings").(map[string]string)

	mediaUploadPath := setting[repository.SettingSiteMediaUploadPath.String()]
	domainImage := setting[repository.SettingGlobalDomainImage.String()]

	domain, _ := url.Parse(domainImage)

	filePath := strings.ReplaceAll(u.Path, domain.Path, mediaUploadPath)

	query := strings.Split(u.RawQuery, "/")

	q := u.Query()

	if (len(query) % 2) != 0 {
		q = map[string][]string{}
	}
	if !file2.PathExist(filePath) {
		// 如果文件不存在，返回404
		err = ErrFileNotFound
		return
	}

	if len(q) == 0 {
		return ioutil.ReadFile(filePath)
	}

	params := map[string]string{}
	for i := 0; i < len(query)-1; i++ {
		v := strings.Split(query[i+1], "|")[0]
		params[query[i]] = v
	}

	f, err := os.Open(filePath)
	if err != nil {
		err = errors.Wrap(err, ErrFileNotFound.Error())
		return
	}

	defer func() {
		_ = f.Close()
	}()

	var width, height, quality, clip int

	if w, ok := params["w"]; ok {
		width, _ = strconv.Atoi(w)
	}
	if h, ok := params["h"]; ok {
		height, _ = strconv.Atoi(h)
	}
	if q, ok := params["q"]; ok {
		quality, _ = strconv.Atoi(q)
	}
	// clip: 1 裁切, clip: 2 压缩
	if iv, ok := params["imageView2"]; ok {
		clip, _ = strconv.Atoi(iv)
	}

	filePaths := strings.Split(filePath, "/")
	fileName := filePaths[len(filePaths)-1]
	names := strings.Split(fileName, ".")
	dst := names[0] + "-" + strconv.Itoa(clip) + "-" + strconv.Itoa(width) + "x" + strconv.Itoa(height) + "-" + strconv.Itoa(quality) + "." + names[1]
	dist := strings.ReplaceAll(filePath, fileName, dst)
	if file2.PathExist(dist) {
		return ioutil.ReadFile(dist)
	}

	fOut, err := os.Create(dist)
	if err != nil {
		_ = level.Error(s.logger).Log("os", "Create", "err", err.Error())
		return ioutil.ReadFile(filePath)
	}
	defer func() {
		_ = fOut.Close()
	}()

	if clip == 2 {
		if err = s.scale(f, fOut, width, height, quality); err != nil {
			_ = level.Error(s.logger).Log("s", "scale", "err", err.Error())
			return
		}
	} else if clip == 1 {
		if err = s.clip(f, fOut, 0, 0, width, height, quality); err != nil {
			_ = level.Error(s.logger).Log("s", "scale", "err", err.Error())
			return
		}
	}

	return ioutil.ReadFile(dist)
}

func (s *service) scale(in io.Reader, out io.Writer, width, height, quality int) error {
	origin, fm, err := image.Decode(in)
	if err != nil {
		return err
	}
	if width == 0 || height == 0 {
		//width = origin.Bounds().Max.X
		height = origin.Bounds().Max.Y
	}
	if quality == 0 {
		quality = 100
	}

	canvas := resize.Thumbnail(uint(width), uint(height), origin, resize.Lanczos3)

	//return jpeg.Encode(out, canvas, &jpeg.Options{quality})

	switch fm {
	case "jpeg":
		return jpeg.Encode(out, canvas, &jpeg.Options{quality})
	case "png":
		return png.Encode(out, canvas)
	case "gif":
		return gif.Encode(out, canvas, &gif.Options{})
	default:
		return errors.New("ERROR FORMAT")
	}
}

func (s *service) clip(in io.Reader, out io.Writer, x0, y0, x1, y1, quality int) error {
	origin, fm, err := image.Decode(in)
	if err != nil {
		return err
	}

	// 先压缩成 1280 再进行裁切
	if origin.Bounds().Max.Y < origin.Bounds().Max.X {
		origin = resize.Resize(0, uint(y1), origin, resize.Lanczos3)
		x0 = (origin.Bounds().Max.X / 2) - (x1 / 2)
		x1 += x0
	} else {
		origin = resize.Resize(uint(x1), 0, origin, resize.Lanczos3)
		y0 = (origin.Bounds().Max.Y / 2) - (y1 / 2)
		y1 += y0
	}

	switch fm {
	case "jpeg":
		img := origin.(*image.YCbCr)
		subImg := img.SubImage(image.Rect(x0, y0, x1, y1)).(*image.YCbCr)
		return jpeg.Encode(out, subImg, &jpeg.Options{quality})
	case "png":
		switch origin.(type) {
		case *image.NRGBA:
			img := origin.(*image.NRGBA)
			subImg := img.SubImage(image.Rect(x0, y0, x1, y1)).(*image.NRGBA)
			return png.Encode(out, subImg)
		case *image.RGBA:
			img := origin.(*image.RGBA)
			subImg := img.SubImage(image.Rect(x0, y0, x1, y1)).(*image.RGBA)
			return png.Encode(out, subImg)
		}
	case "gif":
		img := origin.(*image.Paletted)
		subImg := img.SubImage(image.Rect(x0, y0, x1, y1)).(*image.Paletted)
		return gif.Encode(out, subImg, &gif.Options{})
	default:
		return errors.New("ERROR FORMAT")
	}
	return nil
}

func (s *service) UploadMedia(ctx context.Context, f *multipart.FileHeader) (resImg *imageResponse, err error) {
	file, err := f.Open()
	if err != nil {
		return
	}
	b, err := ioutil.ReadAll(file)
	if err != nil {
		return
	}

	settings := ctx.Value("settings").(map[string]string)
	uploadPath := settings[repository.SettingSiteMediaUploadPath.String()]
	domainUrl := settings[repository.SettingGlobalDomainImage.String()]

	var extName = ".jpg" // image.Image.Header[""]
	if exts, err := mime.ExtensionsByType(f.Header.Get("Content-Type")); err == nil {
		extName = exts[0]
	}

	md5h := md5.New()
	md5h.Write(b)
	fileSha := fmt.Sprintf("%x", md5h.Sum([]byte("")))

	if dbImg, err := s.repository.Image().FindImageByMd5(fileSha); err == nil && dbImg != nil {
		_ = level.Error(s.logger).Log("c.image", "ExistsImageByMd5", "err", "file is exists.")
		return &imageResponse{
			Id:        dbImg.ID,
			Filename:  dbImg.ClientOriginalMame,
			Storename: dbImg.ImageName,
			Size:      dbImg.ImageSize,
			Path:      dbImg.ImagePath,
			Hash:      dbImg.Md5,
			Timestamp: dbImg.CreatedAt.Unix(),
			Url:       domainUrl + "/" + dbImg.ImagePath + s.config.GetString(config.SectionServer, repository.SettingSiteContentImageSuffix.String()),
		}, nil
	}

	fileName := time.Now().Format("20060102") + "-" + fileSha + extName
	simPath := time.Now().Format("2006/01/") + fileSha[len(fileSha)-5:len(fileSha)-3] + "/" + fileSha[24:26] + "/" + fileSha[16:17] + fileSha[12:13] + "/"
	fileFullPath := uploadPath + simPath + fileName

	if !file2.PathExist(uploadPath + simPath) {
		if err = os.MkdirAll(uploadPath+simPath, os.ModePerm); err != nil {
			_ = level.Error(s.logger).Log("os", "MkdirAll", "err", err.Error())
			return
		}
	}

	img := types.Image{
		ImageName:          fileName,
		Extension:          extName,
		ImagePath:          simPath + fileName,
		RealPath:           fileFullPath,
		ImageStatus:        0,
		ImageSize:          strconv.Itoa(int(f.Size)),
		Md5:                "",
		ClientOriginalMame: f.Filename,
		PostID:             0,
	}

	if err = ioutil.WriteFile(fileFullPath, b, 0666); err != nil {
		_ = level.Error(s.logger).Log("ioutil", "WriteFile", "err", err.Error())
		return
	}

	if err = s.repository.Image().AddImage(&img); err != nil {
		_ = level.Error(s.logger).Log("c.image", "AddImage", "err", err.Error())
		return
	}

	return &imageResponse{
		Id:        img.ID,
		Filename:  img.ClientOriginalMame,
		Storename: img.ImageName,
		Size:      img.ImageSize,
		Path:      img.ImagePath,
		Hash:      img.Md5,
		Timestamp: img.CreatedAt.Unix(),
		Url:       domainUrl + "/" + img.ImagePath + s.config.GetString(config.SectionServer, repository.SettingSiteContentImageSuffix.String()),
	}, nil
}

func (s *service) List(ctx context.Context, pageSize, offset int) (images []types.Image, count int64, err error) {
	images, count, err = s.repository.Image().FindAll(pageSize, offset)
	if err != nil {
		_ = level.Error(s.logger).Log("repository.Image", "FindAll", "err", err.Error())
		return
	}

	for k, v := range images {
		images[k].ImagePath = imageUrl(v.ImagePath, s.config.GetString(config.SectionServer, repository.SettingGlobalDomainImage.String()))
	}

	return
}

func NewService(logger log.Logger, repository repository.Repository, config *config.Config) Service {
	return &service{
		logger:     logger,
		repository: repository,
		config:     config,
	}
}

func imageUrl(path, imageDomain string) string {
	return imageDomain + "/" + path
}
