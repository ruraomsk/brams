package main

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/ruraomsk/TLServer/logger"
	"github.com/ruraomsk/brams/drive"
	"github.com/ruraomsk/brams/netcom"
	"github.com/ruraomsk/brams/tester"
)

type Setup struct {
	DbPath  string `toml:"dbpath"`
	LogPath string `toml:"logpath"`
	Step    int    `toml:"step"`
	Timeout int    `toml:"timeout"`
	Port    int    `toml:"port"`
}

var set Setup

func init() {
	if _, err := toml.DecodeFile("brams.toml", &set); err != nil {
		fmt.Printf("Can't load config file %s\n", err.Error())
	}
}
func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	err := logger.Init(set.LogPath + "log")
	if err != nil {
		fmt.Printf("Error init log subsystem %s\n", err.Error())
		return
	}
	drive.SetPath(set.DbPath)
	for _, db := range drive.GetListFilesDbs() {
		drive.AddDb(db)
	}
	stop := make(chan interface{})
	dbstop := make(chan interface{})
	go drive.WorkerDB(set.Step, dbstop)
	go netcom.ServerCommections(set.Port, time.Duration(set.Step*int(time.Second)), stop)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	fmt.Println("server start")
	logger.Info.Println("server start")
	tester.CreateDb()
	go tester.Update("mema", "1", true)
	go tester.Update("memb", "1", false)
	go tester.Update("memc", "2", false)
	go tester.Update("memd", "3", true)
	go tester.Update("mema", "4", false)
	go tester.Update("memb", "2", true)
	go tester.Update("memc", "3", false)
	go tester.Update("memd", "4", true)
	select {
	case <-stop:
		{
			fmt.Println("Wait make stop...")
			dbstop <- 1
			break
		}
	case <-c:
		{
			fmt.Println("Wait make abort...")
			dbstop <- 1
			break
		}
	}
	time.Sleep(3 * time.Second)
	fmt.Println("server stop")
	logger.Info.Println("server stop")
}
