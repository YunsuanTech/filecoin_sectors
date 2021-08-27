package main

import (
	"flag"
	"github.com/e421083458/filecoin_sectors/router"
	"github.com/e421083458/filecoin_sectors/task"
	"github.com/e421083458/golang_common/lib"
	"github.com/robfig/cron"
	"os"
	"os/signal"
	"syscall"
)

var (
	config = flag.String("config", "", "input config file like ./conf/dev/")
)

func main() {
	flag.Parse()
	c := cron.New()
	c.AddFunc("0 */1 * * * ?", func() {
		task.LoopSector()
	})
	c.Start()
	//fmt.Println(123)
	//fmt.Println(*config)
	lib.InitModule(*config, []string{"base", "mysql"})
	//lib.InitModule("./conf/dev/",[]string{"base","mysql",})
	defer lib.Destroy()
	router.HttpServerRun()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	router.HttpServerStop()
}
