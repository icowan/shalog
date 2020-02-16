package repository

import (
	"database/sql"
	"github.com/jinzhu/gorm"
	"github.com/nsini/blog/src/repository/types"
	"time"
)

var (
	_ = time.Second
	_ = sql.LevelDefault
)

type comment struct {
	db *gorm.DB
}

type CommentRepository interface {
	FindByPostId(postId int64) (res *types.Comment, err error)
}

func NewCommentRepository(db *gorm.DB) CommentRepository {
	return &comment{db: db}
}

func (c *comment) FindByPostId(postId int64) (res *types.Comment, err error) {
	var i types.Comment
	if err = c.db.Last(&i, "post_id=?", postId).Error; err != nil {
		return
	}
	return
}
