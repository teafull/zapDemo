/**
* @File: main
* @Author: zf
* @Copyright 2019 zsf5110@163.com.  All rights reserved.
* @Version: 1.0.0
* @Date: 2021/12/12 下午5:29
* @License: MIT License
 */

package main

import (
	"go.uber.org/zap"
	"time"
)

func main() {
	dynamicLevel := zap.NewAtomicLevel()

	// 默认Info级别
	logcfg := zap.NewProductionConfig()
	logcfg.Level = dynamicLevel
	logger, err := logcfg.Build()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	// 获取当前日志等级
	level := dynamicLevel.String()
	logger.Info("this info log", zap.String("logLevel", level))

	// 设置当前日志等级为Debug
	dynamicLevel.SetLevel(zap.DebugLevel)
	logger.Debug("this is debug log.", zap.Int64("st", time.Now().Unix()))

	// http设置日志等级
	// 只支持get、put
	// 格式：{"level":"info"}
	// dynamicLevel.ServeHTTP
}
