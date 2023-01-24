package main

import (
	"context"
	"fmt"
	"github.com/OpenFilWallet/OpenFilWallet/build"
	"github.com/OpenFilWallet/msgboat/conf"
	"github.com/OpenFilWallet/msgboat/server"
	logging "github.com/ipfs/go-log/v2"
	"github.com/urfave/cli/v2"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var log = logging.Logger("msgboat")

func main() {
	_ = logging.SetLogLevel("*", "INFO")

	app := &cli.App{
		Name:                 "msgboat",
		Usage:                "msg boat for OpenFilWallet",
		Version:              build.Version(),
		EnableBashCompletion: true,
		Commands: []*cli.Command{
			runCmd,
		},
	}

	if err := app.Run(os.Args); err != nil {
		os.Stderr.WriteString("Error: " + err.Error() + "\n")
	}
}

var runCmd = &cli.Command{
	Name:  "run",
	Usage: "Start process",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "port",
			Usage: "api port",
			Value: "6689",
		},
	},
	Action: func(cctx *cli.Context) error {
		if err := conf.LocalConfig(); err != nil {
			return err
		}

		var closeCh = make(chan struct{})

		walletServer, err := server.NewBoat(conf.GetNodes())
		if err != nil {
			return fmt.Errorf("new Wallet fail: %s", err.Error())
		}

		router := walletServer.NewRouter()

		endpoint := "0.0.0.0:" + cctx.String("port")

		s := &http.Server{
			Addr:         endpoint,
			Handler:      router,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		}

		log.Infow("start msgboat server", "endpoint", endpoint)
		go func() {
			if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("s.ListenAndServe err: %v", err)
			}
		}()

		quit := make(chan os.Signal)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		close(closeCh)

		log.Info("shutting down msgboat server...")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := s.Shutdown(ctx); err != nil {
			log.Fatal("server forced to shutdown:", err)
		}

		log.Info("msgboat server exit")

		return nil
	},
}
