/**
* @File: main
* @Author: zf
* @Copyright 2019 zsf5110@163.com.  All rights reserved.
* @Version: 1.0.0
* @Date: 2021/12/12 下午5:32
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
		TimeKey:        "tv",                          // json时时间键
		LevelKey:       "l",                           // json时日志等级键
		NameKey:        "N",                           // json时日志记录器键
		CallerKey:      "l",                           // json时日志文件信息键
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
	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "all.log", // 输出文件
		MaxSize:    1,         // 日志文件最大大小（单位：MB）
		MaxBackups: 3,         // 保留的旧日志文件最大数量
		MaxAge:     28,        // 保存日期
		Compress:   true,
	})

	dynamicLevel := zap.NewAtomicLevel()

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(NewEncoderConfig()), // 日志形式 json形式
		//友好形式
		//zapcore.NewConsoleEncoder(NewEncoderConfig())

		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), w), // 日志输出流可以添加多个

		dynamicLevel, // 日志等级
	)
	logger := zap.New(core, zap.AddCaller())
	defer logger.Sync()

	logger.Info("this is info.", zap.Int("id", 1))
	logger.Info("this is info.", zap.Int("id", 1))
	logger.Info("this is info.", zap.Int("id", 1))

	dynamicLevel.SetLevel(zap.DebugLevel)

	for true {
		logger.Debug("this is info.", zap.Int("id", 1))
		logger.Debug("this is info.", zap.Int("id", 1))
		time.Sleep(time.Millisecond)
	}
}
