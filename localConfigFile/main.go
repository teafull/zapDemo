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
	"io/ioutil"
	"log"
	"os"

	"github.com/fsnotify/fsnotify"
	"gopkg.in/yaml.v2"
)

func main() {
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

	// read config file
	fileByte, err := ioutil.ReadFile("localConfigFile/properties.yaml")
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	lp := LoggerProperties{}
	err = yaml.Unmarshal(fileByte, &lp)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Println(lp)

	jsonStr, err := json.Marshal(lp)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Println(string(jsonStr))
}

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

type LoggerProperties struct {
	RootLogger logProperties `yaml:"rootLogger" json:"rootLogger,omitempty"`
}
