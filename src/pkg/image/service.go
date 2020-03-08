package image

import (
	"context"
	"crypto/md5"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/icowan/shalom/src/config"
	"github.com/icowan/shalom/src/repository"
	"github.com/icowan/shalom/src/repository/types"
	file2 "github.com/icowan/shalom/src/util/file"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"os"
	"strconv"
	"time"
)

type Service interface {
	List(ctx context.Context, pageSize, offset int) (images []types.Image, count int64, err error)
	UploadMedia(ctx context.Context, f *multipart.FileHeader) (resImg *imageResponse, err error)

	// 实时渲染图片
	Get(ctx context.Context, path string)
}

type service struct {
	logger     log.Logger
	repository repository.Repository
	config     *config.Config
}

func (s *service) Get(ctx context.Context, path string) {
	panic("implement me")
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
