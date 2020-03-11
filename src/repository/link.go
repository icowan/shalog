package repository

import (
	"github.com/icowan/shalog/src/repository/types"
	"github.com/jinzhu/gorm"
	"time"
)

type LinkRepository interface {
	Add(name, link, icon string, state LinkState) (err error)
	Delete(id int64) (err error)
	List() (links []*types.Link, err error)
	FindByState(state int) (links []types.Link, err error)
	Find(id int64) (link types.Link, err error)
	Update(link *types.Link) (err error)
	FindAll() (links []*types.Link, err error)
}

type LinkState int

const (
	LinkStateApply LinkState = iota
	LinkStatePass
)

func (l LinkState) Int() int {
	return int(l)
}

type link struct {
	db *gorm.DB
}

func (l *link) FindAll() (links []*types.Link, err error) {
	err = l.db.Find(&links).Error
	return
}

func (l *link) Update(link *types.Link) (err error) {
	return l.db.Model(link).Where("id = ?", link.Id).Update(link).Error
}

func (l *link) Find(id int64) (link types.Link, err error) {
	err = l.db.Where("id = ?", id).First(&link).Error
	return
}

func (l *link) FindByState(state int) (links []types.Link, err error) {
	err = l.db.Where("state = ?", state).Find(&links).Error
	return
}

func (l *link) Add(name, link, icon string, state LinkState) (err error) {
	return l.db.Save(&types.Link{
		Name:      name,
		Link:      link,
		Icon:      icon,
		State:     state.Int(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}).Error
}

func (l *link) Delete(id int64) (err error) {
	return l.db.Where("id = ?", id).Delete(&types.Link{
		Id: id,
	}).Error
}

func (l *link) List() (links []*types.Link, err error) {
	err = l.db.Where("state = 1").Find(&links).Error
	return
}

func NewLinkRepository(db *gorm.DB) LinkRepository {
	return &link{db: db}
}
