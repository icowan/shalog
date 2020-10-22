package post

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/icowan/config"
	"github.com/icowan/shalog/src/middleware"
	"github.com/icowan/shalog/src/repository"
	"github.com/icowan/shalog/src/repository/types"
	"github.com/mozillazg/go-pinyin"
	"github.com/pkg/errors"
	"strconv"
	"strings"
	"time"
)

var (
	ErrInvalidArgument     = errors.New("invalid argument")
	ErrPostCreate          = errors.New("发布失败 ")
	ErrPostFind            = errors.New("查询失败")
	ErrPostUpdate          = errors.New("更新失败")
	ErrPostParams          = errors.New("参数错误")
	ErrPostParamTitle      = errors.New("标题不能为空")
	ErrPostParamCategories = errors.New("请选择至少一个分类")
)

type Service interface {
	// 详情页信息
	Get(ctx context.Context, id int64) (rs map[string]interface{}, err error)

	// 列表页
	List(ctx context.Context, order, by, category string, pageSize, offset int) (rs []map[string]interface{}, count int64, other map[string]interface{}, err error)

	// 受欢迎的
	Popular(ctx context.Context) (rs []map[string]interface{}, err error)

	// 更新点赞
	Awesome(ctx context.Context, id int64) (err error)

	// 搜索文章
	Search(ctx context.Context, keyword, tag string, categoryId int64, offset, pageSize int) (posts []*types.Post, total int64, err error)

	// 创建新文章
	NewPost(ctx context.Context, title, description, content string,
		postStatus repository.PostStatus, categories, tags []string, markdown bool, imageId int64) (id int64, err error)

	// 编辑内容 ps: 参数意思就不写了,变量名称就是意思...
	Put(ctx context.Context, id int64, title, description, content string,
		postStatus repository.PostStatus, categories, tags []string, markdown bool, imageId int64) (err error)

	// 删除文章
	Delete(ctx context.Context, id int64) (err error)

	// 恢复文章
	Restore(ctx context.Context, id int64) (err error)

	// 后端列表
	AdminList(ctx context.Context, order, by, category, tag string, pageSize, offset int, keyword string) (posts []*types.Post, total int64, err error)

	// 后台获取详情
	Detail(ctx context.Context, id int64) (rs *types.Post, err error)

	// 标星的会显示在首页Banner上
	Star(ctx context.Context, id int64) (err error)
}

type service struct {
	repository repository.Repository
	logger     log.Logger
	config     *config.Config
}

func (c *service) Star(ctx context.Context, id int64) (err error) {
	p, err := c.repository.Post().FindOnce(id)
	if err != nil {
		return errors.Wrap(err, ErrPostFind.Error())
	}

	if p.Star == 1 {
		p.Star = -1
	} else {
		p.Star = 1
	}
	if err = c.repository.Post().Update(p); err != nil {
		err = errors.Wrap(err, ErrPostUpdate.Error())
	}
	return
}

func (c *service) Detail(ctx context.Context, id int64) (rs *types.Post, err error) {
	rs, err = c.repository.Post().Find(id)
	if err != nil {
		return nil, errors.Wrap(err, ErrPostFind.Error())
	}

	for k, v := range rs.Images {
		v.ImagePath = imageUrl(v.ImagePath, c.config.GetString(config.SectionServer, repository.SettingGlobalDomainImage.String()))
		rs.Images[k] = v
	}

	return
}

func (c *service) AdminList(ctx context.Context, order, by, category, tag string, pageSize, offset int, keyword string) (posts []*types.Post, total int64, err error) {
	userId, _ := ctx.Value(middleware.ContextUserId).(int64)
	posts, total, err = c.repository.Post().FindAll(userId, order, by, offset, pageSize, keyword)
	for k, v := range posts {
		var imgs []types.Image
		for _, img := range v.Images {
			img.ImagePath = imageUrl(img.ImagePath, c.config.GetString(config.SectionServer, repository.SettingGlobalDomainImage.String()))
			imgs = append(imgs, img)
		}
		posts[k].Images = imgs
	}
	return
}

func (c *service) Restore(ctx context.Context, id int64) (err error) {
	post, err := c.repository.Post().FindOnce(id)
	if err != nil {
		_ = level.Error(c.logger).Log("repository.Post", "FindOnce", "err", err.Error())
		return errors.Wrap(err, ErrPostFind.Error())
	}

	post.DeletedAt = nil

	err = c.repository.Post().Update(post)
	if err != nil {
		_ = level.Error(c.logger).Log("repository.Post", "Update", "err", err.Error())
		err = errors.Wrap(err, ErrPostUpdate.Error())
	}
	return
}

func (c *service) Delete(ctx context.Context, id int64) (err error) {
	post, err := c.repository.Post().FindOnce(id)
	if err != nil {
		_ = level.Error(c.logger).Log("repository.Post", "FindOnce", "err", err.Error())
		return errors.Wrap(err, ErrPostFind.Error())
	}

	t := time.Now()
	post.DeletedAt = &t

	err = c.repository.Post().Update(post)
	if err != nil {
		_ = level.Error(c.logger).Log("repository.Post", "Update", "err", err.Error())
		err = errors.Wrap(err, ErrPostUpdate.Error())
	}

	return nil
}

func (c *service) Put(ctx context.Context, id int64, title, description, content string,
	postStatus repository.PostStatus, categories, tags []string, markdown bool, imageId int64) (err error) {

	// todo: 是否需要验证是否为文章本人编辑呢？
	// userId, _ := ctx.Value(middleware.ContextUserId).(int64)

	post, err := c.repository.Post().FindOnce(id)
	if err != nil {
		_ = level.Error(c.logger).Log("repository.Post", "FindOnce", "err", err.Error())
		return errors.Wrap(err, ErrPostFind.Error())
	}

	categoryList, err := c.repository.Category().FindByNames(categories)
	if err != nil {
		_ = level.Error(c.logger).Log("repository.Category", "FindByIds", "err", err.Error())
		return
	}
	var tagList []types.Tag
	for _, v := range tags {
		if t, err := c.repository.Tag().FirstOrCreate(v); err == nil {
			tagList = append(tagList, *t)
		} else {
			_ = level.Error(c.logger).Log("repository.Tag", "FirstOrCreate", "err", err.Error())
		}
	}

	if post.PushTime == nil && postStatus == repository.PostStatusPublish {
		t := time.Now()
		post.PushTime = &t
	} else {
		post.PushTime = nil
	}

	// 清除分类关系表数据
	//_ = c.repository.Category().CleanByPostId(post.ID)

	// 清除tag关系表数据
	//_ = c.repository.Tag().CleanByPostId(post.ID)

	// 清除images的关系数据

	post.Title = title
	post.Description = description
	post.Content = content
	post.PostStatus = postStatus.String()
	post.Categories = categoryList
	post.Tags = tagList
	post.IsMarkdown = markdown

	var imageExists bool
	for _, v := range post.Images {
		if v.ID == imageId {
			imageExists = true
			break
		}
	}

	if !imageExists {
		if img, e := c.repository.Image().FindById(imageId); e == nil {
			var imgs []types.Image
			imgs = append(imgs, img)
			post.Images = imgs
		}
	}

	err = c.repository.Post().Update(post)
	if err != nil {
		_ = level.Error(c.logger).Log("repository.Post", "Update", "err", err.Error())
		err = errors.Wrap(err, ErrPostUpdate.Error())
	}

	return
}

func (c *service) NewPost(ctx context.Context, title, description, content string,
	postStatus repository.PostStatus, categoryNames, tagNames []string, markdown bool, imageId int64) (id int64, err error) {

	userId, _ := ctx.Value(middleware.ContextUserId).(int64)

	categories, err := c.repository.Category().FindByNames(categoryNames)
	if err != nil {
		_ = level.Error(c.logger).Log("repository.Category", "FindByIds", "err", err.Error())
		return
	}
	var tags []types.Tag
	for _, v := range tagNames {
		if t, err := c.repository.Tag().FirstOrCreate(v); err == nil {
			tags = append(tags, *t)
		} else {
			_ = level.Error(c.logger).Log("repository.Tag", "FirstOrCreate", "err", err.Error())
		}
	}

	var pushTime *time.Time
	if postStatus == repository.PostStatusPublish {
		t := time.Now()
		pushTime = &t
	}

	var slug string
	// todo: 没有Gcc 好像不太好使，windows需要自行安装, mac,linux下好像自带
	slug = strings.Join(pinyin.LazyConvert(title, nil), "-")

	// todo: 如果数据不全均为草稿

	post := types.Post{
		Content:     content,
		Description: description,
		Slug:        slug,     // todo: 在transport转换成拼音
		IsMarkdown:  markdown, // todo: 考虑传参进来
		ReadNum:     1,
		Reviews:     1,
		Awesome:     1,
		Title:       title,
		UserID:      userId,
		PostStatus:  postStatus.String(),
		PushTime:    pushTime,
		Tags:        tags,
		Categories:  categories,
	}

	if err = c.repository.Post().Create(&post); err != nil {
		err = errors.Wrap(err, ErrPostCreate.Error())
	}

	return post.ID, nil
}

func (c *service) Search(ctx context.Context, keyword, tag string, categoryId int64, offset, pageSize int) (posts []*types.Post, total int64, err error) {
	if keyword != "" {
		return c.repository.Post().Search(keyword, categoryId, offset, pageSize)
	}

	if tag != "" {
		tagInfo, err := c.repository.Tag().FindPostIdsByName(tag)
		if err != nil {
			_ = level.Warn(c.logger).Log("repository.Tag", "FindPostByName", "err", err.Error())
			return nil, 0, nil
		}

		return c.repository.Post().FindByIds(tagInfo.PostIds, categoryId, offset, pageSize)
	}

	return
}

func (c *service) Awesome(ctx context.Context, id int64) (err error) {
	post, err := c.repository.Post().FindOnce(id)
	if err != nil {
		return
	}
	post.Awesome += 1
	return c.repository.Post().Update(post)
}

func (c *service) Get(ctx context.Context, id int64) (rs map[string]interface{}, err error) {
	detail, err := c.repository.Post().Find(id)
	if err != nil {
		return
	}

	if detail == nil {
		return nil, repository.PostNotFound
	}

	if detail.PushTime == nil || detail.DeletedAt != nil || repository.PostStatus(detail.PostStatus) != repository.PostStatusPublish {
		return nil, repository.PostNotFound
	}

	var headerImage string

	if image, err := c.repository.Image().FindByPostIdLast(id); err == nil && image != nil {
		headerImage = c.config.GetString(config.SectionServer, repository.SettingGlobalDomainImage.String()) + "/" + image.ImagePath
	}

	var category types.Category
	for _, v := range detail.Categories {
		category = v
		break
	}

	// prev
	prev, _ := c.repository.Post().Prev(detail.PushTime, []int64{category.Id})
	// next
	next, _ := c.repository.Post().Next(detail.PushTime, []int64{category.Id})

	populars, _ := c.Popular(ctx)
	return map[string]interface{}{
		"content":      detail.Content,
		"title":        detail.Title,
		"publish_at":   detail.PushTime.Format("2006/01/02 15:04:05"),
		"updated_at":   detail.UpdatedAt,
		"author":       detail.User.Username,
		"comment":      detail.Reviews,
		"banner_image": headerImage,
		"read_num":     strconv.Itoa(int(detail.ReadNum)),
		"description":  strings.Replace(detail.Description, "\n", "", -1),
		"tags":         detail.Tags,
		"populars":     populars,
		"prev":         prev,
		"next":         next,
		"awesome":      detail.Awesome,
		"id":           detail.ID,
	}, nil
}

func (c *service) List(ctx context.Context, order, by, category string, pageSize, offset int) (rs []map[string]interface{},
	count int64, other map[string]interface{}, err error) {
	// 取列表 判断搜索、分类、Tag条件
	// 取最多阅读

	categoriesCh := make(chan []*types.Category)
	tagsCh := make(chan []*types.Tag)

	go func() {
		cateRes, _ := c.repository.Category().FindAll()
		categoriesCh <- cateRes
	}()

	go func() {
		tagsRes, _ := c.repository.Tag().All(20)
		tagsCh <- tagsRes
	}()

	var posts []types.Post
	if category != "" {
		if category, total, err := c.repository.Category().FindByName(category, pageSize, offset); err == nil {
			for _, v := range category.Posts {
				posts = append(posts, v)
			}
			count = total
		}
	} else {
		var categoryIds []int64
		if cates, err := c.repository.Category().FindAll(); err == nil {
			for _, v := range cates {
				categoryIds = append(categoryIds, v.Id)
			}
		}
		posts, count, err = c.repository.Post().FindBy(categoryIds, order, by, pageSize, offset)
		if err != nil {
			_ = level.Warn(c.logger).Log("repository.Post", "FindBy", "err", err.Error())
			return
		}
	}

	var postIds []int64
	for _, post := range posts {
		postIds = append(postIds, post.ID)
	}

	images, err := c.repository.Image().FindByPostIds(postIds)
	if err == nil && images == nil {
		_ = level.Warn(c.logger).Log("c.image.FindByPostIds", "postIds", "err", err)
	}

	imageMap := make(map[int64]string, len(images))
	for _, image := range images {
		imageMap[image.PostID] = imageUrl(image.ImagePath, c.config.GetString(config.SectionServer, repository.SettingGlobalDomainImage.String()))
	}

	_ = c.logger.Log("count", count)

	for _, val := range posts {
		imageUrl, ok := imageMap[val.ID]
		if !ok {
			_ = c.logger.Log("postId", val.ID, "image", ok)
		}
		rs = append(rs, map[string]interface{}{
			"id":         strconv.FormatUint(uint64(val.ID), 10),
			"title":      val.Title,
			"desc":       val.Description,
			"publish_at": val.PushTime.Format("2006/01/02 15:04:05"),
			"image_url":  imageUrl,
			"comment":    val.Reviews,
			"author":     val.User.Username,
			"tags":       val.Tags,
		})
	}

	categories := <-categoriesCh
	tags := <-tagsCh
	close(categoriesCh)
	close(tagsCh)
	other = map[string]interface{}{
		"tags":       tags,
		"category":   category,
		"categories": categories,
	}

	return
}

func (c *service) Popular(ctx context.Context) (rs []map[string]interface{}, err error) {

	posts, err := c.repository.Post().Popular()
	if err != nil {
		return
	}

	var postIds []int64

	for _, post := range posts {
		postIds = append(postIds, post.ID)
	}

	images, err := c.repository.Image().FindByPostIds(postIds)
	if err == nil && images == nil {
		_ = c.logger.Log("c.image.FindByPostIds", "postIds", "err", err)
	}

	imageMap := make(map[int64]string, len(images))
	for _, image := range images {
		imageMap[image.PostID] = imageUrl(image.ImagePath, c.config.GetString(config.SectionServer, repository.SettingGlobalDomainImage.String()))
	}

	for _, post := range posts {
		imageUrl, ok := imageMap[post.ID]
		if !ok {
			_ = c.logger.Log("postId", post.ID, "image", ok)
		}

		desc := []rune(post.Description)
		if len(desc) > 40 {
			desc = desc[:40]
		}
		rs = append(rs, map[string]interface{}{
			"title":     post.Title,
			"desc":      string(desc),
			"id":        post.ID,
			"image_url": imageUrl,
		})
	}

	return
}

func imageUrl(path, imageDomain string) string {
	return imageDomain + "/" + path
}

func NewService(logger log.Logger, cf *config.Config, repo repository.Repository) Service {
	return &service{
		repository: repo,
		logger:     logger,
		config:     cf,
	}
}
