/**
 * @Time : 2019-09-11 17:43
 * @Author : solacowa@gmail.com
 * @File : client
 * @Software: GoLand
 */

package main

import (
	"flag"
	"fmt"
	"github.com/icowan/blog/src/repository"
	"github.com/icowan/blog/src/repository/types"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
	"strconv"
	"strings"
	"time"
)

var host, port, user, password, dbname string

func main() {

	flag.StringVar(&host, "host", "127.0.0.1", "help message for host")
	flag.StringVar(&port, "port", "3306", "help message for port")
	flag.StringVar(&user, "user", "blog", "help message for user")
	flag.StringVar(&password, "password", "admin", "help message for password")
	flag.StringVar(&dbname, "dbname", "blog", "help message for dbname")
	flag.Parse()

	syncPostCategories()
	//changeImage()
	return

	dburl1 := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true&loc=Local&timeout=20m&collation=utf8mb4_unicode_ci",
		"blog", "blog*16347",
		"129.204.31.254", "30596", "superCong")

	dburl2 := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true&loc=Local&timeout=20m&collation=utf8mb4_unicode_ci",
		"blog", "blog*16347",
		"129.204.31.254", "30596", "blog")

	db1, err := gorm.Open("mysql", dburl1)
	if err != nil {
		log.Fatal(err)
	}

	db2, err := gorm.Open("mysql", dburl2)
	if err != nil {
		log.Fatal(err)
	}

	type oldImg struct {
		Id                 int64     `json:"id"`
		ImageName          string    `json:"image_name"`
		Extension          string    `json:"extension"`
		ImagePath          string    `json:"image_path"`
		RealPath           string    `json:"real_path"`
		ImageTime          time.Time `json:"image_time"`
		ImageStatus        int       `json:"image_status"`
		Md5                string    `json:"md5"`
		ClientOriginalName string    `json:"client_original_name"`
		PostId             int64     `json:"post_id"`
	}

	var images []oldImg

	if err := db1.Raw("select * from images").Scan(&images).Error; err != nil {
		log.Fatal(err)
	}

	for _, v := range images {
		log.Println(db2.Save(&types.Image{
			ImageName:          v.ImageName,
			Extension:          v.Extension,
			ImagePath:          strings.ReplaceAll(v.ImagePath, "/mnt/storage/uploads/images/", ""),
			RealPath:           v.RealPath,
			ImageSize:          strconv.Itoa(v.ImageStatus),
			Md5:                v.Md5,
			ClientOriginalMame: v.ClientOriginalName,
			PostID:             v.PostId,
			CreatedAt:          v.ImageTime,
		}).Error)
	}

	type oldPost struct {
		Id          int64     `json:"id"`
		Title       string    `json:"title"`
		Description string    `json:"description"`
		Content     string    `json:"content"`
		IsMarkdown  int       `json:"is_markdown"`
		Status      int       `json:"status"`
		ReadNum     int64     `json:"read_num"`
		Reviews     int64     `json:"reviews"`
		Star        int64     `json:"star"`
		PushTime    time.Time `json:"push_time"`
		CreatedAt   time.Time `json:"created_at"`
		Action      int       `json:"action"`
	}

	var posts []oldPost

	if err = db1.Raw("select * from posts").Scan(&posts).Error; err != nil {
		log.Fatal(err)
	}

	defer func() {
		_ = db2.Close()
		_ = db1.Close()
	}()

	for _, v := range posts {
		var isMarkdown bool
		if v.IsMarkdown == 1 {
			isMarkdown = true
		}
		log.Println(db2.Save(&types.Post{
			ID:          v.Id,
			Action:      v.Action,
			Content:     v.Content,
			Description: v.Description,
			IsMarkdown:  isMarkdown,
			ReadNum:     v.ReadNum,
			Reviews:     v.Reviews,
			Star:        int(v.Star),
			Status:      v.Status,
			Title:       v.Title,
			UserID:      1,
			PostStatus:  string(repository.PostStatusPublish),
			CreatedAt:   v.CreatedAt,
			UpdatedAt:   v.PushTime,
			PushTime:    &v.PushTime,
			DeletedAt:   nil,
		}).Error)
	}
}

func changeImage() {
	dburl2 := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true&loc=Local&timeout=20m&collation=utf8mb4_unicode_ci",
		"blog", "blog*16347",
		"129.204.31.254", "30596", "blog")

	db2, err := gorm.Open("mysql", dburl2)
	if err != nil {
		log.Fatal(err)
	}

	var images []types.Image

	if err := db2.Raw("select * from images where image_path like 'uploads/images%'").Scan(&images).Error; err != nil {
		log.Fatal(err)
	}

	for _, v := range images {
		v.ImagePath = strings.ReplaceAll(v.RealPath, "/mnt/storage/uploads/images/", "")
		err = db2.Model(&types.Image{}).Where("id = ?", v.ID).Update(v).Error
		log.Println(err)
	}
}

func syncPostCategories() {
	dburl := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true&loc=Local&timeout=20m&collation=utf8mb4_unicode_ci",
		user, password,
		host, port, dbname)

	db, err := gorm.Open("mysql", dburl)
	if err != nil {
		log.Fatal(err)
	}

	var posts []types.Post

	if err = db.Find(&posts).Error; err != nil {
		log.Fatal(err)
	}

	repo := repository.NewRepository(db)

	for _, v := range posts {
		cate, _ := repo.Category().Find(int64(v.Action))

		cates := []types.Category{
			cate,
		}

		v.Categories = cates
		if err = repo.Post().Update(&v); err != nil {
			fmt.Println(err.Error())
		}
	}
}
