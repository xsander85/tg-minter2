package main

import (
	"flag"
	"fmt"
	"runtime"
	"time"

	"github.com/xsander85/tg-minter2/pkg/config"
	"github.com/xsander85/tg-minter2/pkg/minter"
	"github.com/xsander85/tg-minter2/pkg/telegram"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "./config/config.json", "path to config file")
}

func main() {
	runtime.GOMAXPROCS(2)
	fmt.Println("Start bot")

	var config config.Config

	if err := config.LoadConfig(configPath); err != nil {
		panic(err)
	}
	fmt.Println(config)
	minterChanel := make(chan string, 100)

	minter := minter.New(&config)
	telegram := telegram.New(&config.Telegram)

	defer minter.Exit()
	go func() { minter.Run(minterChanel) }()
	go func() { telegram.Run(minter) }()
	go func() {
		println(".")
		for {
			telegram.Send(0, <-minterChanel)
		}

	}()

	for {
		time.Sleep(1000)
	}
}
