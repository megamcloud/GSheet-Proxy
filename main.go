package main

import (
	"context"
	"fmt"
	"git.anphabe.net/event/anphabe-event-hub/config"
	"git.anphabe.net/event/anphabe-event-hub/infrastructure/injection"
	"go.uber.org/zap"
	"net/http"
	"os"
	"sync"
	"time"
)

var (
	cfg    *config.ConfigurationInfo
	logger *zap.Logger
	err    error
)

func InitHttp(cfg *config.ConfigurationInfo) *http.Server {

	srv := &http.Server{
		Addr:    cfg.Server,
		Handler: injection.InitGin(),
	}

	go func() {
		if err1 := srv.ListenAndServeTLS("./tls/eventHub.pem","./tls/eventHub-key.pem"); err1 != http.ErrServerClosed {
			logger.Error(err1.Error())
		}
	}()

	outboundAddr := injection.GetOutboundAddress(cfg)
	fmt.Println("Listening at: " + outboundAddr)

	injection.Openbrowser(outboundAddr+"/hello")

	return srv
}


func main() {
	var wg sync.WaitGroup

	cfg = injection.InitConfig()
	logger = injection.InitLogger(cfg)

	defer func() {
		if err != nil {
			logger.Error(err.Error())
		}

		wg.Wait()
		logger.Info("Bye bye.")
		logger.Sync()

		if err != nil {
			os.Exit(1)
		} else {
			os.Exit(0)
		}
	}()

	// Run Importer
	importRunner := injection.InitSourceKeeper()

	wg.Add(1)
	importRunner.Start(&wg)

	// Setup Gin Server
	server := InitHttp(cfg)

	// Signal Handling
	quit := injection.SignalsHandle()
	<-quit

	importRunner.Stop()

	ctx, cancelFunc := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancelFunc()

	if err1 := server.Shutdown(ctx); err1 != nil {
		logger.Error(err1.Error())
	}
}
