/**
* @File: main
* @Author: zf
* @Copyright 2019 zsf5110@163.com.  All rights reserved.
* @Version: 1.0.0
* @Date: 2021/12/12 下午8:42
* @License: MIT License
 */

package main

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func main() {
	hookRoot := lumberjack.Logger{
		Filename:   "./logs/root.log", // 日志文件路径
		MaxSize:    10,                // 每个日志文件保存的最大尺寸 单位：M
		MaxBackups: 1,                 // 日志文件最多保存多少个备份
		MaxAge:     28,                // 文件最多保存多少天
		Compress:   true,              // 是否压缩
	}

	hookPlugin := lumberjack.Logger{
		Filename:   "./logs/plugin.log", // 日志文件路径
		MaxSize:    10,                  // 每个日志文件保存的最大尺寸 单位：M
		MaxBackups: 1,                   // 日志文件最多保存多少个备份
		MaxAge:     28,                  // 文件最多保存多少天
		Compress:   true,                // 是否压缩
	}

	// 公用编码器
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "linenum",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder, // 小写编码器
		EncodeTime:     zapcore.ISO8601TimeEncoder,    // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}

	// 设置日志级别
	atomicLevelRoot := zap.NewAtomicLevel()
	atomicLevelRoot.SetLevel(zap.InfoLevel)
	coreRoot := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),                   // 编码器配置
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(&hookRoot)), // 打印到控制台和文件
		atomicLevelRoot, // 日志级别
	)

	atomicLevelPlugin := zap.NewAtomicLevel()
	atomicLevelPlugin.SetLevel(zap.InfoLevel)
	corePlugin := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(&hookPlugin)),
		atomicLevelPlugin,
	)

	coreTeePlugin := zapcore.NewTee( // print logs of plugin to root that is plugin parent, and it self。
		corePlugin,
		coreRoot, // Print log to the root node by default
	)
	loggerRoot := zap.New(coreRoot, zap.AddCaller(), zap.Fields(zap.String("serviceName", "Root")))
	loggerPlugin := zap.New(coreTeePlugin, zap.AddCaller(), zap.Fields(zap.String("serviceName", "Plugin"))) // print logs to two file

	loggerRoot.Info("root test")
	loggerPlugin.Info("plugin test")
}
