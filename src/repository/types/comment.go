/**
 * @Time : 2019-09-06 11:49
 * @Author : solacowa@gmail.com
 * @File : comment
 * @Software: GoLand
 */

package types

import "time"

type Comment struct {
	ID        int       `gorm:"column:id;primary_key" json:"id"`
	PostID    int64     `gorm:"column:post_id" json:"post_id"`
	LogID     int64     `gorm:"column:log_id" json:"log_id"`
	APIUserID int       `gorm:"column:api_user_id" json:"api_user_id"`
	APIAction string    `gorm:"column:api_action" json:"api_action"`
	APIPostID int64     `gorm:"column:api_post_id" json:"api_post_id"`
	ThreadID  int64     `gorm:"column:thread_id" json:"thread_id"`
	ThreadKey string    `gorm:"column:thread_key" json:"thread_key"`
	CommentIP string    `gorm:"column:comment_ip" json:"comment_ip"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	Message   string    `gorm:"column:message" json:"message"`
	Status    string    `gorm:"column:status" json:"status"`
	ParentID  int64     `gorm:"column:parent_id" json:"parent_id"`
	Type      int       `gorm:"column:type" json:"type"`
	Agent     string    `gorm:"column:agent" json:"agent"`
}

func (*Comment) TableName() string {
	return "comment"
}
