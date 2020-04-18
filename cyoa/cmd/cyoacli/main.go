package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ardanlabs/conf"
	"github.com/pkg/errors"
)

var appName = "cyoa"

func main() {
	if err := run(); err != nil {
		log.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {

	var cfg struct {
		Source string `type:"string" conf:"default:json"`
		Json   struct {
			File string `type:"string"`
			Dir  string `type:"string" conf:"default:../../json"`
		}
		Log string `type:"string" conf:"default:/tmp/cyoa_%s.log"`
	}

	if err := conf.Parse(os.Args[1:], appName, &cfg); err != nil {
		usage, err := conf.Usage(appName, &cfg)
		if err != nil {
			return errors.Wrap(err, "generating config usage")
		}
		fmt.Println(usage)
		return nil
	}
	fmt.Printf("Config: %+v\n", cfg)
	return nil
}
