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
	"github.com/nekonako/moecord/internal/auth"
	"github.com/nekonako/moecord/internal/channel"
	"github.com/nekonako/moecord/internal/message"
	"github.com/nekonako/moecord/internal/server"
	"github.com/nekonako/moecord/internal/websocket"
	"github.com/rs/zerolog/log"
	otelMiddleware "go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
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

	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			log.Error().Err(err).Msg("http server error")
		}
	}()

	log.Info().Msg("server is running on " + config.Api.Host + ":" + fmt.Sprint(config.Api.Port))

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

	r.Use(otelMiddleware.Middleware(c.Apm.ServiceName))

	oauth := auth.New(c, infra)
	oauth.InitRouter(r)

	server := server.New(c, infra)
	server.InitRouter(r)

	channel := channel.New(c, infra)
	channel.InitRouter(r)

	message := message.New(c, infra)
	message.InitRouter(r)

	ws := websocket.New(c, infra)
	ws.InitRouter(r)

	origins := handlers.AllowedOrigins([]string{"*"})
	headers := handlers.AllowedHeaders([]string{"Access-Control-Allow-Origin", "Content-Type"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS", "PUT", "PATCH", "DELETE"})

	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", c.Api.Host, c.Api.Port),
		WriteTimeout: time.Duration(c.Api.WriteTimeout) * time.Second,
		ReadTimeout:  time.Duration(c.Api.ReadTimeout) * time.Second,
		IdleTimeout:  time.Duration(c.Api.IdleTimeout) * time.Second,
		Handler:      handlers.CORS(origins, headers, methods)(r),
	}

	return srv

}
