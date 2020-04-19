package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/ardanlabs/conf"
	"github.com/pkg/errors"
	"github.com/roloum/gophercises/cyoa/internal/cyoa"
	cyoadb "github.com/roloum/gophercises/cyoa/internal/db"
)

var appName = "cyoaweb"

func main() {
	if err := run(); err != nil {
		log.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {

	//Define config struct to parse
	var cfg struct {
		Datastore string `type:"string" conf:"default:json"`
		Json      struct {
			File    string `type:"string"`
			Dir     string `type:"string" conf:"default:../../json"`
			Chapter string `type:"string" conf:"default:intro"`
		}
	}

	log := log.New(os.Stdout, appName,
		log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)

	log.Println("Loading configuration")

	//Parse configuration
	if err := conf.Parse(os.Args[1:], appName, &cfg); err != nil {
		usage, err := conf.Usage(appName, &cfg)
		if err != nil {
			return errors.Wrap(err, "generating config usage")
		}
		fmt.Println(usage)
		return nil
	}

	log.Printf("Configuration loaded.\n%+v\n", cfg)

	//Default datastore is json
	dataStore := cyoa.NewDataStore(&cyoadb.JSON{Log: log})

	story, err := dataStore.LoadStory(cfg.Json.File)
	if err != nil {
		return errors.Wrap(err, "Loading story")
	}

	chapterTpl, err := template.ParseFiles("template/chapter.html")
	if err != nil {
		return errors.Wrap(err, "Parsing HTML template")
	}

	chapterHandler := cyoa.NewChapterHTTPHandler(story, cfg.Json.Chapter, log,
		cyoa.WithNewTemplate(chapterTpl))
	sm := http.NewServeMux()
	sm.Handle("/", chapterHandler)

	server := &http.Server{
		Addr:    ":8000",
		Handler: sm,
	}

	//Run server
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	//Block program until error is thrown or kill signal is received
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan
	log.Println("Terminating.", sig)

	//Shutdown gracefully
	tc, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	cancel()
	server.Shutdown(tc)

	return nil
}
