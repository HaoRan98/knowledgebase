package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"NULL/knowledgebase/models"
	"NULL/knowledgebase/pkg/logging"
	"NULL/knowledgebase/pkg/setting"
	"NULL/knowledgebase/pkg/util"
	"NULL/knowledgebase/routers"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
)

func init() {
	setting.Setup()
	logging.Setup()
	util.Setup()
	models.Setup()
}

func main() {
	gin.SetMode(setting.ServerSetting.RunMode)
	if len(os.Args) == 2 {
		models.InitDb()
		log.Println("*******init browse info over*****")
		log.Println("*******please rerun the program*****")
		return
	}

	routersInit := routers.InitRouter()
	readTimeout := setting.ServerSetting.ReadTimeout
	writeTimeout := setting.ServerSetting.WriteTimeout
	endPoint := fmt.Sprintf(":%d", setting.ServerSetting.HttpPort)
	maxHeaderBytes := 1 << 20

	server := &http.Server{
		Addr:           endPoint,
		Handler:        routersInit,
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
		MaxHeaderBytes: maxHeaderBytes,
	}
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	log.Println(dir)
	log.Printf("[info] start http server listening %s", endPoint)

	//创建监听退出chan
	c := make(chan os.Signal)
	//监听指定信号 ctrl+c kill
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		log.Println("监听进程启动...")
		for s := range c {
			switch s {
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				log.Println("退出", s)
				ExitFunc()
			default:
				log.Println("other signal", s)
			}
		}
	}()

	go func() {
		models.Info, err = models.GetInfoCount()
		if err != nil {
			log.Fatalf("初始化读取搜索统计错误:%v", err)
		}
		log.Printf("当前已累计搜索%d次", models.Info.Browse)
	}()

	err = server.ListenAndServe()
	if err != nil {
		log.Printf("init listen server fail:%v", err)
	}
}

func ExitFunc() {
	log.Println("进程断开,开始存储InfoCount...")
	if err := models.SaveInfoCount(models.Info.ID, models.Info.Browse); err != nil {
		log.Println("存储InfoCount err:", err)
	}
	log.Println("存储InfoCount完成...")
	os.Exit(0)
}
