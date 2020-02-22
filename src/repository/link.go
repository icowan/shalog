package repository

import (
	"github.com/icowan/blog/src/repository/types"
	"github.com/jinzhu/gorm"
	"time"
)

type LinkRepository interface {
	Add(name, link, icon string) (err error)
	Delete(id int64) (err error)
	List() (links []types.Link, err error)
	FindByState(state int) (links []types.Link, err error)

	// 要啥更新，直接删除再加
	Update(id int64, name, link, icon string) (err error)
}

type link struct {
	db *gorm.DB
}

func (l *link) FindByState(state int) (links []types.Link, err error) {
	err = l.db.Where("state = ?", state).Find(&links).Error
	return
}

func (l *link) Add(name, link, icon string) (err error) {
	return l.db.Save(&types.Link{
		Name:      name,
		Link:      link,
		Icon:      icon,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}).Error
}

func (l *link) Delete(id int64) (err error) {
	return l.db.Where("id = ?", id).Delete(&types.Link{
		Id: id,
	}).Error
}

func (l *link) Update(id int64, name, link, icon string) (err error) {
	panic("implement me")
}

func (l *link) List() (links []types.Link, err error) {
	err = l.db.Find(&links).Error
	return
}

func NewLinkRepository(db *gorm.DB) LinkRepository {
	return &link{db: db}
}
