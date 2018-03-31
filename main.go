package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/moogar0880/nap-bot/bot"
	"github.com/moogar0880/nap-bot/config"
)

func main() {

	versionFlag := flag.Bool("version", false, "Version")
	flag.Parse()

	if *versionFlag {
		fmt.Println("Git Commit:", GitCommit)
		fmt.Println("Version:", Version)
		if VersionPrerelease != "" {
			fmt.Println("Version PreRelease:", VersionPrerelease)
		}
		return
	}

	config := config.LoadConfig()
	bot := bot.New(*config)
	if err := bot.Run(context.Background()); err != nil {
		panic(err.Error())
	}

	<-bot.Done
}
