package main

import (
	log "github.com/cihub/seelog"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"net/http"
	_ "net/http/pprof"
	"push-go/push"
	"runtime"
	"strings"
)

func main() {
	go func() {
		log.Info(http.ListenAndServe("localhost:6060", nil))
	}()

	viper.SetConfigFile("./config/App.json")
	if err := viper.ReadInConfig(); err != nil {
		panic("加载配置文件出错")
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		if strings.Contains(e.Name, "App.json") {
			log.Info("Reload App.json Success")
		}
	})

	logger, err := log.LoggerFromConfigAsFile("./config/seelog.xml")
	if err != nil {
		log.Critical("err parsing config log file", err)
		return
	}
	log.ReplaceLogger(logger)
	defer log.Flush()

	log.Info("cpu: ", runtime.NumCPU())
	runtime.GOMAXPROCS(runtime.NumCPU() - 1)

	//Util.InsertDataToDB()

	push.Run()
}
