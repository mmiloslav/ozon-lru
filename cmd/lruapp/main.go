// основной пакет
package main

import (
	"api"
	"config"
	"lru"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	log "github.com/sirupsen/logrus"
)

func main() {
	sigint := make(chan os.Signal, 1)
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt, syscall.SIGINT)
	signal.Notify(sigterm, syscall.SIGTERM)

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()

		for {
			select {
			case <-sigint:
				log.Println("Received SIGINT")
				os.Exit(0)
			case <-sigterm:
				log.Println("Received SIGTERM")
				os.Exit(0)
			}
		}
	}()

	go func() {
		processRequests()
	}()

	wg.Wait()
}

func processRequests() {
	conf, err := config.InitConf()
	if err != nil {
		log.Errorf("failed to init config with error [%s]", err.Error())
		os.Exit(1)
	}

	config.SetLogLevel(conf.LogLevel)

	lruCache := lru.NewLRUCache(conf.CacheSize)
	cacheHandler := api.NewCacheHandler(lruCache, conf.DefaultCacheTTL)

	router := api.NewRouter(cacheHandler)
	err = http.ListenAndServe(conf.ServerHostPort, router)
	if err != nil {
		log.Fatalf("failed to listen and serve with error [%s]", err.Error())
		os.Exit(1)
	}
}
