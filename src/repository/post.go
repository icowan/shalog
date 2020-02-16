package repository

import (
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/nsini/blog/src/repository/types"
	"time"
)

var (
	PostNotFound = errors.New("post not found!")
)

type PostRepository interface {
	Find(id int64) (res *types.Post, err error)
	FindBy(action []int, order, by string, pageSize, offset int) ([]*types.Post, int64, error)
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
}

type PostStatus string

const (
	PostStatusPublish PostStatus = "publish"
	PostStatusDraft   PostStatus = "draft"
)

type post struct {
	db *gorm.DB
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
		Where("action in (?)", action).
		Where("post_status = ?", PostStatusPublish).
		Order("push_time desc").Limit(1).First(&p).Error
	return &p, err
}

func (c *post) Next(publishTime *time.Time, action []int64) (res *types.Post, err error) {
	var p types.Post
	err = c.db.Where("push_time > ?", publishTime).
		Where("action in (?)", action).
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
		Preload("User").
		Preload("Categories").
		Preload("Tags").
		Preload("Images").
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

func (c *post) FindBy(action []int, order, by string, pageSize, offset int) ([]*types.Post, int64, error) {
	posts := make([]*types.Post, 0)
	var count int64
	if err := c.db.Model(&posts).Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("id,username")
	}).Preload("Tags").
		Where("action in (?)", action).
		Where("push_time IS NOT NULL").
		Where("post_status = ?", PostStatusPublish).
		Order(gorm.Expr(by + " " + order)).
		Count(&count).
		Offset(offset).Limit(pageSize).Find(&posts).Error; err != nil {
		return nil, 0, err
	}
	return posts, count, nil
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
