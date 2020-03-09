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
	"github.com/icowan/shalom/src/encode"
	"github.com/icowan/shalom/src/repository/types"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"
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
	db = db.Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci")
	_ = level.Debug(logger).Log("create table", db.CreateTable(types.User{}).Error)
	_ = level.Debug(logger).Log("create table", db.CreateTable(types.Post{}).Error)
	_ = level.Debug(logger).Log("create table", db.CreateTable(types.Image{}).Error)
	_ = level.Debug(logger).Log("create table", db.CreateTable(types.Comment{}).Error)
	_ = level.Debug(logger).Log("create table", db.CreateTable(types.Tag{}).Error)
	_ = level.Debug(logger).Log("create table", db.CreateTable(types.Category{}).Error)
	_ = level.Debug(logger).Log("create table", db.CreateTable(types.Setting{}).Error)
	_ = level.Debug(logger).Log("create table", db.CreateTable(types.Link{}).Error)

	path, err := filepath.Abs(sqlPath)
	if err != nil {
		_ = level.Error(logger).Log("filepath", "Abs", "err", err.Error())
		return err
	}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		_ = level.Error(logger).Log("ioutil", "ReadFile", "err", err.Error())
		return err
	}

	ds := strings.Split(string(data), "');")
	for _, v := range ds {
		if !strings.Contains(v, "');") {
			v += "');"
		}
		if v == "');" {
			continue
		}
		db.Exec(v)
	}

	return store.User().Create(username, encode.EncodePassword(password, appKey), "")
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
