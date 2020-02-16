/**
 * @Time : 2019-09-12 09:42
 * @Author : solacowa@gmail.com
 * @File : logging
 * @Software: GoLand
 */

package logging

import (
	"github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	"path/filepath"
	"time"
)

func LogrusLogger(filePath string) (*logrus.Logger, error) {
	//path, fileName := filepath.Split(filePath)
	linkFile, err := filepath.Abs(filePath)
	if err != nil {
		return nil, err
	}

	logrusLogger := logrus.New()
	writer, err := rotatelogs.New(
		linkFile+"-%Y-%m-%d",
		rotatelogs.WithLinkName(linkFile),         // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(time.Hour*24*365),   // 文件最大保存时间
		rotatelogs.WithRotationTime(time.Hour*24), // 日志切割时间间隔
	)
	if err != nil {
		logrusLogger.Error("Init log failed, err:", err)
		return nil, err
	}

	logrusLogger.SetOutput(writer)
	logrusLogger.SetFormatter(&logrus.TextFormatter{
		DisableColors:   false,
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
	})

	return logrusLogger, nil
}
