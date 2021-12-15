/**
* @File: main
* @Author: zf
* @Copyright 2019 zsf5110@163.com.  All rights reserved.
* @Version: 1.0.0
* @Date: 2021/12/13 下午10:38
* @License: MIT License
 */

package main

import (
	"encoding/json"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
)

const (
	logfile = "localConfigFile/properties.yaml"
)

func main() {

	// 沒有日志配置文件時，使用默認的日志配置
	debugLoggerHook := lumberjack.Logger{
		Filename:   "./logs/localConfigFile.debug.log", // 日志文件路径
		MaxSize:    100,                                // 每个日志文件保存的最大尺寸 单位：M
		MaxAge:     7,                                  // 文件最多保存多少天
		MaxBackups: 10,                                 // 日志文件最多保存多少个备份
		Compress:   true,                               // 是否压缩
	}
	errorLoggerHook := lumberjack.Logger{ // 僅輸出錯誤級別以上的日志
		Filename:   "./logs/localConfigFile.error.log",
		MaxSize:    50,
		MaxAge:     7,
		MaxBackups: 10,
		Compress:   true,
	}

	// 公用编码器
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "lv",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stack",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}

	// 设置日志级别
	debugAtomicLevel := zap.NewAtomicLevel()
	debugAtomicLevel.SetLevel(zap.DebugLevel)
	debugCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(&debugLoggerHook), zapcore.AddSync(os.Stdout)),
		debugAtomicLevel,
	)
	errorCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(&errorLoggerHook)),
		zap.WarnLevel,
	)
	coreTee := zapcore.NewTee( // print logs of plugin to root that is plugin parent, and it self。
		debugCore,
		errorCore, // Print log to the root node by default
	)
	drvDebug := zap.New(coreTee, zap.AddCaller(), zap.Fields(zap.String("serviceName", "Root")))

	drvDebug.Info("info 1...................")
	drvDebug.Debug("debug 1.................")
	drvDebug.Warn("Warn 1.................")
	drvDebug.Error("Error 1.................")

	fmt.Println("")
	debugAtomicLevel.SetLevel(zap.InfoLevel)

	drvDebug.Info("info 2...................")
	drvDebug.Debug("debug 2.................")
	drvDebug.Warn("Warn   3.................")
	drvDebug.Error("Error   3.................")

	return

	// first load log config file
	lp, err := ReadLogConfigFile(logfile)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	jsonTmp, _ := json.Marshal(lp)
	fmt.Println(string(jsonTmp))

	// monitor config file
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Println("event:", event)

				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("modified file:", event.Name)

					// Reload file
					lp, err := ReadLogConfigFile(logfile)
					if err != nil {
						fmt.Println(err)
					} else {
						jsonTmp, _ := json.Marshal(lp)
						fmt.Println(string(jsonTmp))
					}
				} else {
					// Unwanted operation
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			case exist := <-done:
				log.Println(exist)
				return
			}
		}
	}()
	err = watcher.Add("localConfigFile/properties.yaml")
	if err != nil {
		log.Fatal(err)
	}
	<-done
}

// ReadLogConfigFile read log config from local config file.
func ReadLogConfigFile(configFile string) (LoggerProperties, error) {
	fileByte, err := ioutil.ReadFile(configFile)
	if err != nil {
		return LoggerProperties{}, err
	}

	lp := LoggerProperties{}
	err = yaml.Unmarshal(fileByte, &lp)
	if err != nil {
		return LoggerProperties{}, err
	}

	return lp, nil
}

// logProperties define logger property
type logProperties struct {
	File          string `yaml:"file,omitempty" json:"file,omitempty"`
	MaxSize       int    `yaml:"maxSize,omitempty" json:"maxSize,omitempty"`
	MaxAge        int    `yaml:"maxAge,omitempty" json:"maxAge,omitempty"`
	MaxBackup     int    `yaml:"maxBackup,omitempty" json:"maxBackup,omitempty"`
	Interval      int    `yaml:"interval,omitempty" json:"interval,omitempty"`
	CallerSkip    int    `yaml:"callerSkip,omitempty" json:"callerSkip,omitempty"`
	EnableConsole bool   `yaml:"enableConsole,omitempty" json:"enableConsole,omitempty"`
	Async         bool   `yaml:"async,omitempty" json:"async,omitempty"`
	Level         string `yaml:"level,omitempty" json:"level,omitempty"`
	AddCaller     bool   `yaml:"addCaller,omitempty" json:"addCaller,omitempty"`
	PatternLayout string `yaml:"patternLayout,omitempty" json:"patternLayout,omitempty"`

	// appender properties
	Append  bool   `yaml:"append,omitempty" json:"append,omitempty"`
	LogName string `yaml:"logName,omitempty" json:"logName,omitempty"`

	Appenders []logProperties `yaml:"appenders,omitempty" json:"appenders,omitempty"`
}

// LoggerProperties define log
type LoggerProperties struct {
	RootLogger logProperties `yaml:"rootLogger" json:"rootLogger,omitempty"`
}
