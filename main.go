/**
 * @Time : 2019-09-06 10:52
 * @Author : solacowa@gmail.com
 * @File : main
 * @Software: GoLand
 */

package main

import (
	"github.com/icowan/shalog/cmd/service"
	"runtime"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	service.Run()
}

//type Post struct {
//	Id int64
//}
//
//func test()  {
//	var posts []Post
//
//	var ps []*Post
//	posts = append(posts, Post{Id:1}, Post{Id:2}, Post{Id:3}, Post{Id:4})
//
//	for _, v := range posts {
//		t := v
//		ps = append(ps, &t)
//	}
//
//	for _, v := range ps {
//		fmt.Println(v.Id)
//	}
//}
