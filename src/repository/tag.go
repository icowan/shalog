/**
 * @Time : 2019-09-10 11:05
 * @Author : solacowa@gmail.com
 * @File : meta
 * @Software: GoLand
 */

package repository

import (
	"github.com/icowan/shalog/src/repository/types"
	"github.com/jinzhu/gorm"
)

type MetaType string

const (
	MetaCategory MetaType = "category"
	MetaTag      MetaType = "tag"
)

type TagRepository interface {
	FirstOrCreate(name string) (meta *types.Tag, err error)
	All(limit int) (metas []*types.Tag, err error)
	FindPostByName(name string) (meta types.Tag, err error)
	FindPostIdsByName(name string) (meta types.Tag, err error)
	Find(id int64) (meta types.Tag, err error)
	FindByIds(ids []int64) (tags []types.Tag, err error)
	FindByNames(names []string) (tags []types.Tag, err error)
	CleanByPostId(id int64) (err error)
	List(tagName string, limit, offset int) (metas []*types.Tag, count int64, err error)
	Delete(id int64) (err error)
	Update(id int64, name string) (err error)
	TagCountById(id int64) int
	UpdateCount(tag *types.Tag) (err error)
}

type tag struct {
	db *gorm.DB
}

func (c *tag) UpdateCount(tag *types.Tag) (err error) {
	return c.db.Model(&types.Tag{Id: tag.Id}).Where("id = ?", tag.Id).Update(tag).Error
}

func (c *tag) TagCountById(id int64) int {
	return c.db.Model(&types.Tag{Id: id}).Where("id = ?", id).Association("Posts").Count()
}

func (c *tag) Update(id int64, name string) (err error) {
	// todo: 如果名称已经存在需要返回错误信息
	return c.db.Model(&types.Tag{}).Where("id = ?", id).Update(&types.Tag{
		Id:   id,
		Name: name,
	}).Error
}

func (c *tag) Delete(id int64) (err error) {
	if err = c.db.Model(&types.Tag{Id: id}).
		Where("id = ?", id).
		Association("Posts").
		Clear().Error; err == nil {
		err = c.db.Model(&types.Tag{}).Where("id = ?", id).Delete(&types.Tag{Id: id}).Error
	}
	return
}

func (c *tag) List(tagName string, limit, offset int) (metas []*types.Tag, count int64, err error) {
	query := c.db.Model(&types.Tag{})
	if tagName != "" {
		query = query.Where("name link ?", "%"+tagName+"%")
	}
	err = query.Count(&count).Offset(offset).Limit(limit).Order("id desc").Find(&metas).Error
	return
}

func (c *tag) FindByNames(names []string) (tags []types.Tag, err error) {
	err = c.db.Model(&types.Tag{}).Where("name in (?)", names).Find(&tags).Error
	return
}

func (c *tag) CleanByPostId(id int64) (err error) {
	panic("implement me")
}

func (c *tag) FindByIds(ids []int64) (tags []types.Tag, err error) {
	err = c.db.Model(&types.Tag{}).Where("id in (?)", ids).Find(&tags).Error
	return
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

func (c *tag) All(limit int) (metas []*types.Tag, err error) {
	err = c.db.Model(&types.Tag{}).Order("count desc").Limit(limit).Find(&metas).Error
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
