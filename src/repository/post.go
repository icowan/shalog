package repository

import (
	"errors"
	"github.com/icowan/shalog/src/repository/types"
	"github.com/jinzhu/gorm"
	"time"
)

var (
	PostNotFound = errors.New("post not found!")
)

type PostRepository interface {
	Find(id int64) (res *types.Post, err error)
	FindBy(categoryIds []int64, order, by string, pageSize, offset int) ([]types.Post, int64, error)
	Popular() (posts []*types.Post, err error)
	SetReadNum(p *types.Post) error
	Create(p *types.Post) error
	Update(p *types.Post) error
	Stars() (res []*types.Post, err error)
	Index() (res []*types.Post, err error)
	Prev(publishTime *time.Time, action []int64) (res *types.Post, err error)
	Next(publishTime *time.Time, action []int64) (res *types.Post, err error)
	Count() (total int64, err error)
	FindOnce(id int64) (res *types.Post, err error)
	Search(keyword string, categoryId int64, offset, pageSize int) (res []*types.Post, count int64, err error)
	FindByIds(ids []int64, categoryId int64, offset, pageSize int) (res []*types.Post, count int64, err error)
	FindByCategoryId(categoryId int64, limit int) (posts []types.Post, err error)
	FindAll(userId int64, order, by string, offset, pageSize int, keyword string) (posts []*types.Post, count int64, err error)
}

type PostStatus string

const (
	PostStatusPublish PostStatus = "publish"
	PostStatusDraft   PostStatus = "draft"
)

func (c PostStatus) String() string {
	return string(c)
}

type post struct {
	db *gorm.DB
}

func (c *post) FindAll(userId int64, order, by string, offset, pageSize int, keyword string) (posts []*types.Post, count int64, err error) {
	db := c.db.Model(&types.Post{}).Preload("Categories").Preload("Images", func(db *gorm.DB) *gorm.DB {
		return db.Select("image_name,image_path,post_id")
	})
	if userId > 0 {
		db = db.Where("user_id = ?", userId)
	}
	if keyword != "" {
		db = db.Where("title like ? OR content like ?", `%`+keyword+`%`, `%`+keyword+`%`)
	}
	err = db.Order(gorm.Expr(by + " " + order)).
		Count(&count).
		Offset(offset).
		Limit(pageSize).
		Find(&posts).Error
	return
}

func (c *post) FindByIds(ids []int64, categoryId int64, offset, pageSize int) (res []*types.Post, count int64, err error) {
	err = c.db.Model(&types.Post{}).
		Where("id in (?)", ids).
		//Where("action in (?)", categoryId).
		Where("push_time IS NOT NULL").
		Where("post_status = ?", PostStatusPublish).
		Order("push_time desc").
		Count(&count).
		Offset(offset).Limit(pageSize).Find(&res).Error
	return
}

func (c *post) Search(keyword string, categoryId int64, offset, pageSize int) (res []*types.Post, count int64, err error) {
	err = c.db.Model(&types.Post{}).
		Where("title like ? OR content like ?", `%`+keyword+`%`, `%`+keyword+`%`).
		//Where("action in (?)", categoryId).
		Where("push_time IS NOT NULL").
		Where("post_status = ?", PostStatusPublish).
		Order("push_time desc").
		Count(&count).
		Offset(offset).Limit(pageSize).Find(&res).Error
	return
}

func (c *post) Count() (total int64, err error) {
	err = c.db.Model(&types.Post{}).Where("post_status = ?", PostStatusPublish).Count(&total).Error
	return
}

func (c *post) Prev(publishTime *time.Time, action []int64) (res *types.Post, err error) {
	var p types.Post
	err = c.db.Where("push_time < ?", publishTime).
		//Where("action in (?)", action). // todo: 注释掉后可能会是个bug ，下一帖跟下一贴的bug
		Where("post_status = ?", PostStatusPublish).
		Order("push_time desc").Limit(1).First(&p).Error
	return &p, err
}

func (c *post) Next(publishTime *time.Time, action []int64) (res *types.Post, err error) {
	var p types.Post
	err = c.db.Where("push_time > ?", publishTime).
		//Where("action in (?)", action). // todo: 注释掉后可能会是个bug ，下一帖跟下一贴的bug
		Where("post_status = ?", PostStatusPublish).
		Order("push_time asc").Limit(1).First(&p).Error
	return &p, err
}

func NewPostRepository(db *gorm.DB) PostRepository {
	return &post{db: db}
}

func (c *post) Index() (res []*types.Post, err error) {
	err = c.db.Where("post_status = ?", PostStatusPublish).
		Preload("Images").
		Preload("User").
		Preload("Tags").
		Order(gorm.Expr("push_time DESC")).
		Limit(10).Find(&res).Error

	return
}

func (c *post) Stars() (res []*types.Post, err error) {
	err = c.db.Where("star = 1").
		Where("post_status = ?", PostStatusPublish).
		Preload("Images").
		Order(gorm.Expr("push_time DESC")).
		Limit(7).Find(&res).Error
	return
}

func (c *post) Update(p *types.Post) error {
	return c.db.Model(p).Where("id = ?", p.ID).Update(p).Error
}

func (c *post) Find(id int64) (res *types.Post, err error) {
	var p types.Post

	if err = c.db.Model(&p).
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id,username")
		}).
		Preload("Categories").
		Preload("Tags").
		Preload("Images", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, image_name, extension, image_path, post_id")
		}).
		Find(&p, "id = ?", id).Error; err != nil {
		return nil, PostNotFound
	}
	return &p, nil
}

func (c *post) FindOnce(id int64) (res *types.Post, err error) {
	var p types.Post

	if err = c.db.Model(&p).
		First(&p, "id = ?", id).Error; err != nil {
		return nil, PostNotFound
	}
	return &p, nil
}

func (c *post) FindByCategoryId(categoryId int64, limit int) (posts []types.Post, err error) {
	err = c.db.Model(&types.Post{}).
		Where("id in (SELECT post_id FROM post_categories WHERE category_id = ?)", categoryId).
		Where("push_time IS NOT NULL").
		Where("post_status = ?", PostStatusPublish).
		Preload("Images").
		Order("push_time desc").
		Limit(limit).
		Find(&posts).Error
	return
}

func (c *post) FindBy(categoryIds []int64, order, by string, pageSize, offset int) (posts []types.Post, count int64, err error) {
	var categories []types.Category
	err = c.db.Model(&types.Category{}).Preload("Posts", func(db *gorm.DB) *gorm.DB {
		var ids []int64
		for _, v := range categories {
			ids = append(ids, v.Id)
		}
		return db.Model(&types.Post{}).
			Preload("User", func(db *gorm.DB) *gorm.DB {
				return db.Select("id,username")
			}).
			Preload("Tags").
			Where("id in (SELECT post_id FROM post_categories WHERE category_id in (?))", ids).
			Where("push_time IS NOT NULL").
			Where("post_status = ?", PostStatusPublish).
			Order(gorm.Expr(by + " " + order)).
			Count(&count).
			Offset(offset).Limit(pageSize).Find(&posts)
	}).Where("id in (?)", categoryIds).Find(&categories).Error

	return
}

func (c *post) Popular() (posts []*types.Post, err error) {
	if err = c.db.Order("read_num DESC").Limit(9).Find(&posts).Error; err != nil {
		return
	}
	return
}

func (c *post) SetReadNum(p *types.Post) error {
	p.ReadNum += 1
	return c.db.Exec("UPDATE `posts` SET `read_num` = ?  WHERE `posts`.`deleted_at` IS NULL AND `posts`.`id` = ?", p.ReadNum, p.ID).Error
}

func (c *post) Create(p *types.Post) error {
	if err := c.db.Save(p).Error; err != nil {
		return err
	}
	return nil
}
