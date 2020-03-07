/**
 * @Time : 2019-09-06 11:30
 * @Author : solacowa@gmail.com
 * @File : service_gen
 * @Software: GoLand
 */

package service

import (
	"errors"
	"github.com/go-kit/kit/log/level"
	"github.com/go-sql-driver/mysql"
	"github.com/icowan/blog/src/repository/types"
	"regexp"
)

var (
	ErrBadFormat = errors.New("invalid format")
	//ErrUnresolvableHost = errors.New("unresolvable host")

	emailRegexp = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)

func validateFormat(email string) error {
	if !emailRegexp.MatchString(email) {
		return ErrBadFormat
	}
	return nil
}

func importToDb() error {
	if _, err := store.User().Find("1"); err != nil {
		switch err.(type) {
		case *mysql.MySQLError:
			e := err.(*mysql.MySQLError)
			if e.Number == 1146 {
				goto CREATE
			}
		}
		return nil
	}

	return nil

CREATE:
	_ = level.Debug(logger).Log("create table", db.CreateTable(types.User{}).Error)
	_ = level.Debug(logger).Log("create table", db.CreateTable(types.Post{}).Error)
	_ = level.Debug(logger).Log("create table", db.CreateTable(types.Image{}).Error)
	_ = level.Debug(logger).Log("create table", db.CreateTable(types.Comment{}).Error)
	_ = level.Debug(logger).Log("create table", db.CreateTable(types.Tag{}).Error)
	_ = level.Debug(logger).Log("create table", db.CreateTable(types.Category{}).Error)
	_ = level.Debug(logger).Log("create table", db.CreateTable(types.Setting{}).Error)
	_ = level.Debug(logger).Log("create table", db.CreateTable(types.Link{}).Error)

	//u := strings.Split(adminEmail, "@")

	//member := &types.User{
	//	Email:      adminEmail,
	//	Username:   u[0],
	//	Password:   null.StringFrom(encode.EncodePassword(adminPassword, appKey)),
	//	Roles:      roles,
	//	Namespaces: nss,
	//}
	//return store.User().(member)

	return nil

}

func logLevel(logLevel string) (opt level.Option) {
	switch logLevel {
	case "warn":
		opt = level.AllowWarn()
	case "error":
		opt = level.AllowError()
	case "debug":
		opt = level.AllowDebug()
	case "info":
		opt = level.AllowInfo()
	case "all":
		opt = level.AllowAll()
	default:
		opt = level.AllowNone()
	}

	return
}
