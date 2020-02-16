/**
 * @Time : 2019-09-06 11:45
 * @Author : solacowa@gmail.com
 * @File : post
 * @Software: GoLand
 */

package types

import (
	"time"
)

type Post struct {
	ID          int64      `gorm:"column:id;primary_key" json:"id"`
	Action      int        `gorm:"column:action" json:"action"`
	Content     string     `gorm:"column:content;type:text" json:"content"`
	Description string     `gorm:"column:description" json:"description"`
	Slug        string     `gorm:"column:slug" json:"slug"`
	IsMarkdown  bool       `gorm:"column:is_markdown" json:"is_markdown"`
	ReadNum     int64      `gorm:"column:read_num" json:"read_num"`
	Reviews     int64      `gorm:"column:reviews" json:"reviews"`
	Star        int        `gorm:"column:star" json:"star"`
	Status      int        `gorm:"column:status" json:"status"`
	Title       string     `gorm:"column:title" json:"title"`
	UserID      int64      `gorm:"column:user_id" json:"user_id"`
	PostStatus  string     `gorm:"column:post_status" json:"post_status"`
	CreatedAt   time.Time  `gorm:"column:created_at;type:datetime" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"column:updated_at;type:datetime" json:"updated_at"`
	PushTime    *time.Time `gorm:"column:push_time" json:"push_time"`
	DeletedAt   *time.Time `gorm:"column:deleted_at;type:datetime" json:"deleted_at"`
	Awesome     int        `orm:"column:awesome" json:"awesome"`
	User        User       `gorm:"ForeignKey:id;AssociationForeignKey:UserId"`
	Tags        []Tag      `gorm:"many2many:post_tags;foreignkey:id;association_foreignkey:id;association_jointable_foreignkey:tag_id;jointable_foreignkey:post_id;" json:"tags"`
	Categories  []Category `gorm:"many2many:post_categories;foreignkey:id;association_foreignkey:id;association_jointable_foreignkey:category_id;jointable_foreignkey:post_id;" json:"categories"`
	Images      []Image    `gorm:"foreignkey:post_id" json:"images"`
}

func (p *Post) TableName() string {
	return "posts"
}
