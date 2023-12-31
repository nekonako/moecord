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
	"github.com/nekonako/moecord/internal/profile"
	"github.com/nekonako/moecord/internal/server"
	"github.com/nekonako/moecord/internal/sfu"
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

	log.Info().Msg("server is running on " + config.Http.Host + ":" + fmt.Sprint(config.Http.Port))

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

	sfu := sfu.New(c)

	ws := websocket.New(c, infra, sfu)
	ws.InitRouter(r)

	oauth := auth.New(c, infra)
	oauth.InitRouter(r)

	server := server.New(c, infra, ws)
	server.InitRouter(r)

	channel := channel.New(c, infra)
	channel.InitRouter(r)

	message := message.New(c, infra, ws)
	message.InitRouter(r)

	profile := profile.New(c, infra)
	profile.InitRouter(r)

	origins := handlers.AllowedOrigins([]string{"*"})
	headers := handlers.AllowedHeaders([]string{"Access-Control-Allow-Origin", "Content-Type"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS", "PUT", "PATCH", "DELETE"})

	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", c.Http.Host, c.Http.Port),
		WriteTimeout: time.Duration(c.Http.WriteTimeout) * time.Second,
		ReadTimeout:  time.Duration(c.Http.ReadTimeout) * time.Second,
		IdleTimeout:  time.Duration(c.Http.IdleTimeout) * time.Second,
		Handler:      handlers.CORS(origins, headers, methods)(r),
	}

	return srv

}
