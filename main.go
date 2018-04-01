package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/moogar0880/nap-bot/bot"
	"github.com/moogar0880/nap-bot/config"
)

type arrayFlags []string

func (i *arrayFlags) String() string {
	return fmt.Sprintf("%q", *i)
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

var configFiles arrayFlags

func main() {

	versionFlag := flag.Bool("version", false, "Version")
	flag.Var(&configFiles, "file", "A config file to load. May be specified multiple times")
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

	// TODO: this is a bit of a dirty hack, we should move this into a
	// structured config file
	for _, f := range configFiles {
		if _, err := os.Stat(f); err == nil {
			data, err := ioutil.ReadFile(f)
			if err != nil {
				panic(err.Error())
			}
			config.AuthToken = string(data)
			fmt.Println(len(config.AuthToken))
		}
	}

	bot := bot.New(*config)
	if err := bot.Run(context.Background()); err != nil {
		panic(err.Error())
	}

	<-bot.Done
}
