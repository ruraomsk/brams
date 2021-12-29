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
	"github.com/ruraomsk/brams/setup"
	"github.com/ruraomsk/brams/tester"
)

func init() {
	if _, err := toml.DecodeFile("brams.toml", &setup.Set); err != nil {
		fmt.Printf("Can't load config file %s\n", err.Error())
	}
}
func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	err := logger.Init(setup.Set.LogPath + "log")
	if err != nil {
		fmt.Printf("Error init log subsystem %s\n", err.Error())
		return
	}
	stop := make(chan interface{})
	dbstop := make(chan interface{})

	drive.SetupBrams(setup.Set)
	drive.StartBrams(dbstop)

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
