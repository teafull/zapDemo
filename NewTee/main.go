/**
* @File: main
* @Author: zf
* @Copyright 2019 zsf5110@163.com.  All rights reserved.
* @Version: 1.0.0
* @Date: 2021/12/12 下午5:50
* @License: MIT License
 */

package main

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"time"
)

// 格式化时间
func TimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

func NewEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		// Keys can be anything except the empty string.
		TimeKey:        "T",                           // json时时间键
		LevelKey:       "L",                           // json时日志等级键
		NameKey:        "N",                           // json时日志记录器名
		CallerKey:      "C",                           // json时日志文件信息键
		MessageKey:     "M",                           // json时日志消息键
		StacktraceKey:  "S",                           // json时堆栈键
		LineEnding:     zapcore.DefaultLineEnding,     // 友好日志换行符
		EncodeLevel:    zapcore.CapitalLevelEncoder,   // 友好日志等级名大小写（info INFO）
		EncodeTime:     TimeEncoder,                   // 友好日志时日期格式化
		EncodeDuration: zapcore.StringDurationEncoder, // 时间序列化
		EncodeCaller:   zapcore.ShortCallerEncoder,    // 日志文件信息（包/文件.go:行号）
	}
}

func main() {
	// 动态日志等级
	dynamicLevel := zap.NewAtomicLevel()

	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "all1.log", // 输出文件
		MaxSize:    500,        // 日志文件最大大小（单位：MB）
		MaxBackups: 3,          // 保留的旧日志文件最大数量
		MaxAge:     28,         // 保存日期
	})

	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.InfoLevel
	})

	core := zapcore.NewTee(
		// 有好的格式、输出控制台、动态等级
		zapcore.NewCore(zapcore.NewConsoleEncoder(NewEncoderConfig()), os.Stdout, dynamicLevel),
		// json格式、输出文件、处定义等级规则
		zapcore.NewCore(zapcore.NewJSONEncoder(zap.NewDevelopmentEncoderConfig()), w, highPriority),
	)
	logger := zap.New(core, zap.AddCaller())
	defer logger.Sync()

	// 将当前日志等级设置为Debug
	// 注意日志等级低于设置的等级，日志文件也不分记录
	dynamicLevel.SetLevel(zap.DebugLevel)

	logger.Info("this is info.", zap.Int("ID", 1))
}
