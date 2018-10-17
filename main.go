package main

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
	"./logc"
	"./conf"
	"./db"
	"./util"
	"./web/routes"
	"./http"
	"github.com/valyala/fasthttp"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			logc.Error(err)
		}
	}()
	defer func() {
		if err := recover(); err != nil {
			logc.Error(err)
			logc.Error(string(util.PanicTrace(5)))
		}
	}()
	config := conf.NewConf()
	logc.InitLogger(config.AppPath, config.LogLevel)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	var stopLock sync.Mutex
	go func() {
		//阻塞程序运行，直到收到终止的信号
		<-signalChan
		stopLock.Lock()
		defer stopLock.Unlock()
		// 此处写入程序退出要执行的代码

		if logc.IsInfo() {
			logc.Info("已安全退出系统")
		}
		db.Close()
		os.Exit(0)
	}()
	Startup(routes.Bind())
}

var router *http.Router

func Startup(initRouter *http.Router) {
	router = initRouter
	var config = conf.NewConf()
	logc.Info("程序运行路径:" + config.AppPath)
	logc.Info("Web服务器地址:" + config.ServerURL)
	fasthttp.ListenAndServe(config.ServerURL, http.HttpHandler(config.AppPath, router))
}