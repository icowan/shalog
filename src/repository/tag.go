/**
 * @Time : 2019-09-10 11:05
 * @Author : solacowa@gmail.com
 * @File : meta
 * @Software: GoLand
 */

package repository

import (
	"github.com/icowan/blog/src/repository/types"
	"github.com/jinzhu/gorm"
)

type MetaType string

const (
	MetaCategory MetaType = "category"
	MetaTag      MetaType = "tag"
)

type TagRepository interface {
	FirstOrCreate(name string) (meta *types.Tag, err error)
	List(limit int) (metas []*types.Tag, err error)
	FindPostByName(name string) (meta types.Tag, err error)
	FindPostIdsByName(name string) (meta types.Tag, err error)
	Find(id int64) (meta types.Tag, err error)
}

type tag struct {
	db *gorm.DB
}

func (c *tag) FindPostIdsByName(name string) (meta types.Tag, err error) {
	// name需要加索引
	var tag types.Tag
	err = c.db.Model(&types.Tag{}).First(&tag, "name = ?", name).Error
	if err != nil {
		return
	}

	var ids []struct {
		PostId int64
	}

	err = c.db.Table("post_tags").Select("post_id").Where("tag_id = ?", tag.Id).Find(&ids).Error

	for _, v := range ids {
		meta.PostIds = append(meta.PostIds, v.PostId)
	}

	return
}
func (c *tag) FindPostByName(name string) (meta types.Tag, err error) {
	// name需要加索引
	err = c.db.Model(&types.Tag{}).Preload("Posts").
		First(&meta, "name = ?", name).Error
	return
}

func (c *tag) Find(id int64) (meta types.Tag, err error) {
	err = c.db.First(&meta, "id = ?", id).Error
	return
}

func (c *tag) List(limit int) (metas []*types.Tag, err error) {
	err = c.db.Model(&types.Tag{}).Order("id desc").Limit(limit).Find(&metas).Error
	return
}

func NewTagRepository(db *gorm.DB) TagRepository {
	return &tag{db: db}
}

func (c *tag) FirstOrCreate(name string) (tag *types.Tag, err error) {
	t := types.Tag{
		Name:        name,
		Description: name,
	}
	err = c.db.FirstOrCreate(&t, types.Tag{
		Name: name,
	}).Error

	return &t, err
}
