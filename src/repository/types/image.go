/**
 * @Time : 2019-09-06 11:46
 * @Author : solacowa@gmail.com
 * @File : image
 * @Software: GoLand
 */

package types

import "time"

type Image struct {
	ID                 int64      `gorm:"column:id;primary_key" json:"id"`
	ImageName          string     `gorm:"column:image_name" json:"image_name"`
	Extension          string     `gorm:"column:extension" json:"extension"`
	ImagePath          string     `gorm:"column:image_path" json:"image_path"`
	RealPath           string     `gorm:"column:real_path" json:"real_path"`
	ImageStatus        int        `gorm:"column:image_status" json:"image_status"`
	ImageSize          string     `gorm:"column:image_size" json:"image_size"`
	Md5                string     `gorm:"column:md5" json:"md5"`
	ClientOriginalMame string     `gorm:"column:client_original_mame" json:"client_original_mame"`
	PostID             int64      `gorm:"column:post_id" json:"post_id"`
	CreatedAt          time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt          time.Time  `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt          *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}

func (p *Image) TableName() string {
	return "images"
}
