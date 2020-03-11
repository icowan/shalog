package api

import (
	"context"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/icowan/shalog/src/config"
	"github.com/icowan/shalog/src/repository"
	"github.com/icowan/shalog/src/repository/types"
	"github.com/icowan/shalog/src/util/file"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"mime"
	"os"
	"strconv"
	"strings"
	"time"
)

var PostNotFound = errors.New("post not found!")

type Service interface {
	Authentication(ctx context.Context, req postRequest) (rs getUsersBlogsResponse, err error)

	// 发布内容
	Post(ctx context.Context, req postRequest) (rs newPostResponse, err error)

	// 获取文章
	GetPost(ctx context.Context, id int64) (rs *getPostResponse, err error)

	// 编辑文章
	EditPost(ctx context.Context, id int64, req postRequest) (rs newPostResponse, err error)

	// 获取分类列表
	GetCategories(ctx context.Context) (rs *getCategoriesResponse, err error) // todo 需要调整 不应该让service返回xml

	// 上传流媒体类型的文件
	MediaObject(ctx context.Context, req postRequest) (rs *getPostResponse, err error)

	// 上传图片资源
	UploadImage(ctx context.Context, image uploadImageRequest) (res imageResponse, err error)
}

type service struct {
	repository repository.Repository
	logger     log.Logger
	config     *config.Config
}

type PostFields string

const (
	PostStatus      PostFields = "post_status"
	PostType        PostFields = "post_type"
	PostCategories  PostFields = "categories"
	PostTitle       PostFields = "title"
	PostDateCreated PostFields = "dateCreated"
	PostWpSlug      PostFields = "wp_slug"
	PostDescription PostFields = "description"
	PostKeywords    PostFields = "mt_keywords"
	MediaOverwrite  PostFields = "overwrite"
	MediaBits       PostFields = "bits"
	MediaName       PostFields = "name"
	MediaType       PostFields = "type"
)

func (c *service) UploadImage(ctx context.Context, image uploadImageRequest) (res imageResponse, err error) {

	if err = ioutil.WriteFile("/tmp/"+image.Image.Filename, image.Image.Body, 0666); err != nil {
		_ = level.Error(c.logger).Log("ioutil", "WriteFile", "err", err.Error())
		return
	}

	// 先存，再验证，失败再删除
	f, err := os.Open("/tmp/" + image.Image.Filename)
	if err != nil {
		_ = level.Error(c.logger).Log("os", "Open", "err", err.Error())
		return
	}

	defer func() {
		if err = f.Close(); err != nil {
			_ = level.Error(c.logger).Log("f", "Close", "err", err.Error())
		}
	}()

	md5h := md5.New()
	if _, err = io.Copy(md5h, f); err != nil {
		_ = level.Error(c.logger).Log("io", "Copy", "err", err.Error())
		return
	}

	fileSha := fmt.Sprintf("%x", md5h.Sum([]byte("")))

	if dbImg, err := c.repository.Image().FindImageByMd5(fileSha); err == nil && dbImg != nil {
		_ = level.Error(c.logger).Log("c.image", "ExistsImageByMd5", "err", "file is exists.")
		size, _ := strconv.ParseInt(dbImg.ImageSize, 10, 64)
		return imageResponse{
			Filename:  dbImg.ClientOriginalMame,
			Storename: dbImg.ImageName,
			Size:      size,
			Path:      dbImg.ImagePath,
			Hash:      dbImg.Md5,
			Timestamp: dbImg.CreatedAt.Unix(),
			Url:       c.config.GetString(config.SectionServer, repository.SettingGlobalDomainImage.String()) + "/" + dbImg.ImagePath + c.config.GetString(config.SectionServer, repository.SettingSiteContentImageSuffix.String()),
		}, nil
	}

	if err != nil {
		_ = level.Error(c.logger).Log("c.image", "ExistsImageByMd5", "err", "file is exists.")
		return
	}

	simPath := time.Now().Format("2006/01/") + fileSha[len(fileSha)-5:len(fileSha)-3] + "/" + fileSha[24:26] + "/" + fileSha[16:17] + fileSha[12:13] + "/"
	filePath := c.config.GetString(config.SectionServer, repository.SettingSiteMediaUploadPath.String()) + "/" + simPath
	if !file.PathExist(filePath) {
		if err = os.MkdirAll(filePath, os.ModePerm); err != nil {
			_ = level.Error(c.logger).Log("os", "MkdirAll", "err", err.Error())
			return
		}
	}

	var extName = ".jpg" // image.Image.Header[""]
	if exts, err := mime.ExtensionsByType(image.Image.Header.Get("Content-Type")); err == nil {
		extName = exts[0]
	}

	fileName := time.Now().Format("20060102") + "-" + fileSha + extName
	fileFullPath := filePath + fileName

	//if err = os.Rename("/tmp/"+image.Image.Filename, fileFullPath); err != nil {
	//	_ = level.Error(c.logger).Log("os", "Rename", "err", err.Error())
	//	return
	//}
	if err = moveFile("/tmp/"+image.Image.Filename, fileFullPath); err != nil {
		_ = level.Error(c.logger).Log("os", "moveFile", "err", err.Error())
		return
	}

	saveImage := &types.Image{
		ImageName:          fileName,
		Extension:          extName,
		ImagePath:          simPath + fileName,
		RealPath:           fileFullPath,
		ImageStatus:        0,
		ImageSize:          strconv.Itoa(int(image.Image.Size)),
		Md5:                fileSha,
		ClientOriginalMame: image.Image.Filename,
	}
	// 存入数据库
	if err = c.repository.Image().AddImage(saveImage); err != nil {
		_ = level.Error(c.logger).Log("c.image", "AddImage", "err", err.Error())
		return
	}

	size, _ := strconv.ParseInt(saveImage.ImageSize, 10, 64)
	return imageResponse{
		Filename:  saveImage.ClientOriginalMame,
		Storename: saveImage.ImageName,
		Size:      size,
		Path:      saveImage.ImagePath,
		Hash:      saveImage.Md5,
		Timestamp: saveImage.CreatedAt.Unix(),
		Url:       c.config.GetString(config.SectionServer, repository.SettingGlobalDomainImage.String()) + "/" + saveImage.ImagePath + c.config.GetString(config.SectionServer, repository.SettingSiteContentImageSuffix.String()),
	}, nil

}

func (c *service) Authentication(ctx context.Context, req postRequest) (rs getUsersBlogsResponse, err error) {

	return
}

func (c *service) EditPost(ctx context.Context, id int64, req postRequest) (rs newPostResponse, err error) {

	post, err := c.repository.Post().Find(id)
	if err != nil {
		return
	}

	if post == nil {
		return rs, PostNotFound
	}

	var postStatus, postType, postTitle, slug, description string
	var categories []string
	var keywords []string
	var postDateCreated time.Time

	for _, member := range req.Params.Param[3].Value.Struct.Member {
		_ = c.logger.Log("member", member.Name)
		switch PostFields(member.Name) {
		case PostStatus:
			postStatus = member.Value.String
		case PostType:
			postType = member.Value.String
		case PostCategories:
			for _, val := range member.Value.Array.Data {
				categories = append(categories, val.Value.String)
			}
		case PostTitle:
			postTitle = member.Value.String
		case PostDateCreated:
			load, _ := time.LoadLocation("Asia/Shanghai")
			if postDateCreated, err = time.ParseInLocation("20060102T15:04:05Z", member.Value.DateTimeIso8601, load); err == nil {
				_ = c.logger.Log("time", "Parse", "err", err)
				postDateCreated = postDateCreated.Add(8 * 3600 * time.Second)
			} else {
				postDateCreated = time.Now()
			}
		case PostWpSlug:
			slug = member.Value.String
		case PostDescription:
			description = member.Value.String
		case PostKeywords:
			keywords = strings.Split(member.Value.String, ",")
		}
	}

	_ = c.logger.Log("req.Params.Param[4].Value.Boolean", req.Params.Param[4].Value.Boolean)

	publishStatus, _ := strconv.Atoi(req.Params.Param[4].Value.Boolean) // todo 1: 已发布，0: 草稿

	_ = c.logger.Log("postStatus", postStatus, "postType", postType, "categories", categories, "postDateCreated", postDateCreated.Format("2006-01-02 15:04:05"), "postTitle", postTitle, "slug", slug, "keywords", keywords)

	desc := []rune(description)
	if len(desc) > 100 {
		desc = desc[:100]
	}

	var tags []types.Tag

	for _, v := range keywords {
		tag, err := c.repository.Tag().FirstOrCreate(v)
		if err != nil {
			_ = level.Error(c.logger).Log("Tag", "FirstOrCreate", "err", err.Error())
			continue
		}
		tags = append(tags, *tag)
	}

	var cates []types.Category
	for _, v := range categories {
		category, err := c.repository.Category().FirstOrCreate(v)
		if err != nil {
			_ = level.Error(c.logger).Log("Category", "FirstOrCreate", "err", err.Error())
			continue
		}
		cates = append(cates, *category)
	}

	post.Title = postTitle
	post.Content = description
	post.Description = string(desc)
	post.IsMarkdown = true
	post.Status = publishStatus
	post.Categories = cates
	post.Tags = tags
	post.PostStatus = postStatus
	post.Slug = slug

	if post.PushTime == nil && repository.PostStatus(postStatus) == repository.PostStatusPublish {
		now := time.Now()
		post.PushTime = &now
	}

	if err = c.repository.Post().Update(post); err != nil {
		return
	}

	rs.Params.Param.Value.String = strconv.Itoa(int(post.ID))

	return
}

func (c *service) MediaObject(ctx context.Context, req postRequest) (rs *getPostResponse, err error) {
	var overwrite bool
	var bits, mediaName, mediaType string
	for _, val := range req.Params.Param[3].Value.Struct.Member {
		switch PostFields(val.Name) {
		case MediaOverwrite:
			overwrite, _ = strconv.ParseBool(val.Value.Boolean)
		case MediaBits:
			bits = val.Value.Base64
		case MediaName:
			mediaName = strings.TrimSpace(strings.ToLower(val.Value.String))
		case MediaType:
			mediaType = val.Value.String
		}
	}

	bits = strings.TrimSpace(strings.Trim(bits, "\n"))
	bits = strings.Replace(bits, " ", "", -1)
	dist, err := base64.StdEncoding.DecodeString(bits)

	if err != nil {
		_ = c.logger.Log("base64", "DecodeString", "err", err.Error())
		return
	}

	if err = ioutil.WriteFile("/tmp/"+mediaName, dist, 0666); err != nil {
		_ = c.logger.Log("ioutil", "WriteFile", "err", err.Error())
		return
	}

	// 先存，再验证，失败再删除
	f, err := os.Open("/tmp/" + mediaName)
	if err != nil {
		_ = c.logger.Log("os", "Open", "err", err.Error())
		return
	}
	defer func() {
		if err = f.Close(); err != nil {
			_ = c.logger.Log("f", "Close", "err", err.Error())
		}
	}()

	md5h := md5.New()
	if _, err = io.Copy(md5h, f); err != nil {
		_ = c.logger.Log("io", "Copy", "err", err.Error())
		return
	}

	var fileSize int64
	if fileInfo, err := os.Stat("/tmp/" + mediaName); err == nil {
		fileSize = fileInfo.Size()
	}

	//defer func() {
	//	if err = os.Remove("/tmp/"+mediaName); err != nil {
	//		_ = c.logger.Log("os", "Remove", "err", err.Error())
	//	}
	//}()

	fileSha := fmt.Sprintf("%x", md5h.Sum([]byte("")))

	// todo 进行数据md5值验证 需要不需要返回地址呢？
	if c.repository.Image().ExistsImageByMd5(fileSha) {
		_ = c.logger.Log("c.image", "ExistsImageByMd5", "err", "file is exists.")
		return
	}

	simPath := time.Now().Format("2006/01/") + fileSha[len(fileSha)-5:len(fileSha)-3] + "/" + fileSha[24:26] + "/" + fileSha[16:17] + fileSha[12:13] + "/"
	filePath := c.config.GetString(config.SectionServer, repository.SettingSiteMediaUploadPath.String()) + "/" + simPath
	if !file.PathExist(filePath) {
		if err = os.MkdirAll(filePath, os.ModePerm); err != nil {
			_ = c.logger.Log("os", "MkdirAll", "err", err.Error())
			return
		}
	}

	var extName = ".jpg"
	if exts, err := mime.ExtensionsByType(mediaType); err == nil {
		extName = exts[0]
	}

	fileName := time.Now().Format("20060102") + "-" + fileSha + extName
	fileFullPath := filePath + fileName

	if err = os.Rename("/tmp/"+mediaName, fileFullPath); err != nil {
		_ = c.logger.Log("os", "Rename", "err", err.Error())
		return
	}

	// 存入数据库
	if err = c.repository.Image().AddImage(&types.Image{
		ImageName:          fileName,
		Extension:          extName,
		ImagePath:          simPath + fileName,
		RealPath:           fileFullPath,
		ImageStatus:        0,
		ImageSize:          strconv.Itoa(int(fileSize)),
		Md5:                fileSha,
		ClientOriginalMame: mediaName,
	}); err != nil {
		_ = c.logger.Log("c.image", "AddImage", "err", err.Error())
		return
	}

	// todo 返回图片的xml response

	_ = c.logger.Log("overwrite", overwrite, "bits", "", "mediaName", mediaName, "mediaType", mediaType, "fileSha", fileSha, "fileName", fileName)

	var members []member

	members = append(members, member{
		Name: "id",
		Value: memberValue{
			String: "0",
		},
	}, member{
		Name: "file",
		Value: memberValue{
			String: mediaName,
		},
	}, member{
		Name: "url",
		Value: memberValue{
			String: c.config.GetString(config.SectionServer, repository.SettingGlobalDomainImage.String()) + "/" + simPath + fileName + c.config.GetString(config.SectionServer, repository.SettingSiteContentImageSuffix.String()),
		},
	}, member{
		Name: "type",
		Value: memberValue{
			String: "",
		},
	})

	resp := getPostResponse{}
	resp.Params.Param.Value.Struct.Member = members

	return &resp, nil
}

func (c *service) Post(ctx context.Context, req postRequest) (rs newPostResponse, err error) {
	userId := ctx.Value(UserIdContext).(int64)

	_ = c.logger.Log("methodName", req.MethodName, "username", req.Params.Param[1].Value.String, "password", req.Params.Param[2].Value.String)

	var postStatus, postType, postTitle, slug, description string
	var categories []string
	var keywords []string
	var postDateCreated time.Time

	for _, member := range req.Params.Param[3].Value.Struct.Member {
		_ = c.logger.Log("member", member.Name)
		switch PostFields(member.Name) {
		case PostStatus:
			postStatus = member.Value.String
		case PostType:
			postType = member.Value.String
		case PostCategories:
			for _, val := range member.Value.Array.Data {
				categories = append(categories, val.Value.String)
			}
		case PostTitle:
			postTitle = member.Value.String
		case PostDateCreated:
			load, _ := time.LoadLocation("Asia/Shanghai")
			if postDateCreated, err = time.ParseInLocation("20060102T15:04:05Z", member.Value.DateTimeIso8601, load); err == nil {
				_ = c.logger.Log("time", "Parse", "err", err)
				postDateCreated = postDateCreated.Add(8 * 3600 * time.Second)
			} else {
				postDateCreated = time.Now()
			}
		case PostWpSlug:
			slug = member.Value.String
		case PostDescription:
			description = member.Value.String
		case PostKeywords:
			keywords = strings.Split(member.Value.String, ",")
		}
	}

	var tags []types.Tag

	for _, v := range keywords {
		tag, err := c.repository.Tag().FirstOrCreate(v)
		if err != nil {
			_ = level.Error(c.logger).Log("Tag", "FirstOrCreate", "err", err.Error())
			continue
		}
		tags = append(tags, *tag)
	}

	var cates []types.Category
	for _, v := range categories {
		category, err := c.repository.Category().FirstOrCreate(v)
		if err != nil {
			_ = level.Error(c.logger).Log("Category", "FirstOrCreate", "err", err.Error())
			continue
		}
		cates = append(cates, *category)
	}

	_ = c.logger.Log("req.Params.Param[4].Value.Boolean", req.Params.Param[4].Value.Boolean)

	publishStatus, _ := strconv.Atoi(req.Params.Param[4].Value.Boolean) // todo 1: 已发布，0: 草稿

	_ = c.logger.Log("postStatus", postStatus, "postType", postType, "categories", categories, "postDateCreated", postDateCreated.Format("2006-01-02 15:04:05"), "postTitle", postTitle, "slug", slug, "keywords", keywords)

	desc := []rune(description)
	if len(desc) > 100 {
		desc = desc[:100]
	}

	now := time.Now()

	p := types.Post{
		Title:       postTitle,
		Content:     description,
		Description: string(desc),
		IsMarkdown:  true, // todo 想办法怎么验证一下
		UserID:      userId,
		Status:      publishStatus,
		ReadNum:     1,
		Awesome:     1,
		PostStatus:  postStatus,
		Tags:        tags,
		Categories:  cates,
		Slug:        slug,
	}

	if publishStatus == 1 {
		p.PushTime = &now
	}

	if err = c.repository.Post().Create(&p); err != nil {
		return
	}

	rs.Params.Param.Value.String = strconv.Itoa(int(p.ID))

	return
}

func (c *service) GetPost(ctx context.Context, id int64) (rs *getPostResponse, err error) {

	_ = c.logger.Log("postId", id)

	post, err := c.repository.Post().Find(id)
	if err != nil {
		return nil, PostNotFound
	}
	var categoryName string
	for _, v := range post.Categories {
		categoryName = v.Name
		break
	}

	var tags []string
	for _, tag := range post.Tags {
		tags = append(tags, tag.Name)
	}

	var members []member

	members = append(members, member{
		Name: "userid",
		Value: memberValue{
			String: strconv.Itoa(int(post.UserID)),
		},
	}, member{
		Name: "postid",
		Value: memberValue{
			String: strconv.Itoa(int(post.ID)),
		},
	}, member{
		Name: "description",
		Value: memberValue{
			String: post.Description,
		},
	}, member{
		Name: "title",
		Value: memberValue{
			String: post.Title,
		},
	}, member{
		Name: "link",
		Value: memberValue{
			String: "/post/" + strconv.Itoa(int(post.ID)),
		},
	}, member{
		Name: "mt_keywords",
		Value: memberValue{
			String: strings.Join(tags, ","),
		},
	}, member{
		Name: "wp_slug",
		Value: memberValue{
			String: post.Slug,
		},
	}, member{
		Name: "wp_author",
		Value: memberValue{
			String: "",
		},
	}, member{
		Name: "wp_author_id",
		Value: memberValue{
			String: "",
		},
	}, member{
		Name: "date_created_gmt",
		Value: memberValue{
			String: post.CreatedAt.String(),
		},
	}, member{
		Name: "post_status",
		Value: memberValue{
			String: "publish",
		},
	},
		member{
			Name: "categories",
			Value: memberValue{
				Array: array{
					Data: data{
						Value: categoryName,
					},
				},
			},
		},
		member{
			Name: "sticky",
			Value: memberValue{
				String: "0",
			},
		})

	resp := new(getPostResponse)

	resp.Params.Param.Value.Struct.Member = members
	return resp, nil
}

func (c *service) GetCategories(ctx context.Context) (rs *getCategoriesResponse, err error) {
	categorys, err := c.repository.Category().FindAll()
	if err != nil {
		return
	}

	resp := new(getCategoriesResponse)
	var dvs []dataValue
	for _, v := range categorys {
		dvs = append(dvs, dataValue{
			Struct: valStruct{
				Member: []member{
					{Name: "categoryId", Value: memberValue{String: strconv.Itoa(int(v.Id))}},
					{Name: "parentId", Value: memberValue{String: strconv.Itoa(int(v.ParentId))}},
					{Name: "categoryName", Value: memberValue{String: v.Name}},
					{Name: "description", Value: memberValue{String: v.Description}},
					{Name: "title", Value: memberValue{String: v.Name}},
				},
			},
		})
	}

	resp.Params.Param.Value.Array.Data.Value = dvs

	return resp, nil
}

func NewService(logger log.Logger, cf *config.Config, repository repository.Repository) Service {
	return &service{
		repository: repository,
		logger:     logger,
		config:     cf,
	}
}

func moveFile(sourcePath, destPath string) error {
	inputFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("Couldn't open source file: %s", err)
	}

	defer func() {
		_ = inputFile.Close()
	}()

	outputFile, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("Couldn't open dest file: %s", err)
	}
	defer func() {
		_ = outputFile.Close()
	}()
	_, err = io.Copy(outputFile, inputFile)
	if err != nil {
		return fmt.Errorf("Writing to output file failed: %s", err)
	}
	err = os.Remove(sourcePath)
	if err != nil {
		return fmt.Errorf("Failed removing original file: %s", err)
	}
	return nil
}
