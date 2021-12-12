/**
* @File: main
* @Author: zf
* @Copyright 2019 zsf5110@163.com.  All rights reserved.
* @Version: 1.0.0
* @Date: 2021/12/12 下午5:22
* @License: MIT License
 */

package main

import (
	"go.uber.org/zap"
	"time"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	logger.Info("failed to fetch URL",
		// Structured context as strongly typed Field values.
		zap.String("url", "https://github.com/uber-go/zap"),
		zap.Int("attempt", 3),
		zap.Duration("backoff", time.Second),
	)

	logger.Info("logger debug", zap.String("Name", "zhangfeng"))
}
