/**
 * @Time : 2019-09-10 11:15
 * @Author : solacowa@gmail.com
 * @File : category
 * @Software: GoLand
 */

package repository

import (
	"github.com/icowan/blog/src/repository/types"
	"github.com/jinzhu/gorm"
)

type CategoryRepository interface {
	FirstOrCreate(name string) (cate *types.Category, err error)
	FindAll() (res []*types.Category, err error)
	FindByName(name string, pageSize, offset int) (cate types.Category, count int64, err error)
	Find(id int64) (cate types.Category, err error)
	FindByIds(ids []int64) (categories []types.Category, err error)
	CleanByPostId(id int64) (err error)
}

type category struct {
	db *gorm.DB
}

func (c *category) CleanByPostId(id int64) (err error) {
	panic("implement me")
}

func (c *category) FindByIds(ids []int64) (categories []types.Category, err error) {
	err = c.db.Model(&types.Category{}).Where("id in (?)", ids).Find(&categories).Error
	return
}

func (c *category) Find(id int64) (cate types.Category, err error) {
	err = c.db.Where("id = ?", id).First(&cate).Error
	return
}

func (c *category) FindByName(name string, pageSize, offset int) (cate types.Category, count int64, err error) {
	var res []types.Post
	err = c.db.Model(&types.Category{}).Preload("Posts", func(db *gorm.DB) *gorm.DB {
		return db.Model(&types.Post{}).
			Where("id in (SELECT post_id FROM post_categories WHERE category_id = ?)", cate.Id).
			Where("push_time IS NOT NULL").
			Where("post_status = ?", PostStatusPublish).
			Order("push_time desc").
			Count(&count).
			Offset(offset).Limit(pageSize).Find(&res)
	}).Where("name = ?", name).First(&cate).Error
	return
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &category{db: db}
}

func (c *category) FirstOrCreate(name string) (cate *types.Category, err error) {
	cc := types.Category{
		Name:        name,
		Description: name,
	}
	err = c.db.FirstOrCreate(&cc, types.Category{
		Name: name,
	}).Error

	return &cc, err
}

func (c *category) FindAll() (res []*types.Category, err error) {
	err = c.db.Find(&res).Error
	return
}
