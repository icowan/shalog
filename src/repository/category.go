/**
 * @Time : 2019-09-10 11:15
 * @Author : solacowa@gmail.com
 * @File : category
 * @Software: GoLand
 */

package repository

import (
	"github.com/jinzhu/gorm"
	"github.com/nsini/blog/src/repository/types"
)

type CategoryRepository interface {
	FirstOrCreate(name string) (cate *types.Category, err error)
	FindAll() (res []*types.Category, err error)
}

type category struct {
	db *gorm.DB
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
