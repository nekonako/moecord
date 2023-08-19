package main

import (
	"context"
	"flag"
	"fmt"

	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/nekonako/moecord/config"
	"github.com/nekonako/moecord/infra"
	"github.com/nekonako/moecord/oauth"
	"github.com/nekonako/moecord/websocket"
	"github.com/rs/zerolog/log"
)

func main() {

	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish")
	flag.Parse()

	config, err := config.New()
	if err != nil {
		panic(err)
	}

	infra := infra.New(config)
	httpServer := newHttpServer(config, infra)
	ws, err := websocket.New(config)
	if err != nil {
		panic(err)
	}

	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			log.Error().Err(err).Msg("http server error")
		}
	}()

	go ws.ListenConnection()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c

	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	httpServer.Shutdown(ctx)
	log.Info().Msg("shutting down")
	os.Exit(0)

}

func newHttpServer(c *config.Config, infra *infra.Infra) *http.Server {

	r := mux.NewRouter()

	oauth := oauth.New(c, infra)
	oauth.InitRouter(r)

	origins := handlers.AllowedOrigins([]string{"*"})
	headers := handlers.AllowedHeaders([]string{"Access-Control-Allow-Origin", "Content-Type"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS"})

	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", c.Api.Host, c.Api.Port),
		WriteTimeout: time.Duration(c.Api.WriteTimeout) * time.Second,
		ReadTimeout:  time.Duration(c.Api.ReadTimeout) * time.Second,
		IdleTimeout:  time.Duration(c.Api.IdleTimeout) * time.Second,
		Handler:      handlers.CORS(origins, headers, methods)(r),
	}

	return srv

}
