package main

import (
	"fmt"
	"net/http"
	"nh-downloader/internal/config"
	"nh-downloader/internal/rpc/nh_rpc"
	"nh-downloader/routes"
	"nh-downloader/utils/logs"
	"os"
)

func main() {
	config.Init()
	err := createFolders()
	if err != nil {
		fmt.Println("[main] createfolders failed")
		panic(err)
	}
	logs.Init()
	nh_rpc.Init()
	routes.Init()

	server := &http.Server{
		Addr:    fmt.Sprintf(":%v", config.Port()),
		Handler: routes.Router,
	}

	logs.Info("Server listening on port:", config.Port())
	logs.Fatal(server.ListenAndServe())
}

func createFolders() error {
	_, err := os.Stat(config.DownloadPath())
	if os.IsNotExist(err) {
		err = os.Mkdir(config.DownloadPath(), os.ModePerm)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	_, err = os.Stat(config.CachePath())
	if os.IsNotExist(err) {
		err = os.Mkdir(config.CachePath(), os.ModePerm)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	_, err = os.Stat(config.LogsPath())
	if os.IsNotExist(err) {
		err = os.Mkdir(config.LogsPath(), os.ModePerm)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	return nil
}
