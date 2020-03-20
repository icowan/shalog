/**
 * @Time: 2020/3/20 14:04
 * @Author: solacowa@gmail.com
 * @File: service_job
 * @Software: GoLand
 */

package service

import (
	"context"
	"github.com/go-kit/kit/log/level"
	"github.com/icowan/shalog/src/pkg/tag"
	"github.com/robfig/cron"
)

func countTag(c *cron.Cron, service tag.Service, spec string) {
	//spec := "0 */6 * * *" // 每6时执行

	if err := c.AddFunc(spec, func() {
		if err := service.UpdateTagCount(context.Background()); err != nil {
			_ = level.Error(logger).Log("service", "UpdateTagCount", "err", err.Error())
		}
	}); err != nil {
		_ = level.Error(logger).Log("c", "AddFunc", "err", err.Error())
	}
}
