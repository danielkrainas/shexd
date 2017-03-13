package server

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/danielkrainas/gobag/context"
	"github.com/rs/cors"
	"github.com/urfave/negroni"

	"github.com/danielkrainas/shex/api/server/handlers"
	"github.com/danielkrainas/shex/registry/actions"
	"github.com/danielkrainas/shex/registry/configuration"
)

func New(ctx context.Context, config configuration.HTTPConfig, actionPack actions.Pack) (*Server, error) {
	api, err := handlers.NewApi(actionPack)
	if err != nil {
		return nil, fmt.Errorf("error creating server api: %v", err)
	}

	log := acontext.GetLogger(ctx)
	n := negroni.New()

	n.Use(cors.New(cors.Options{
		AllowedOrigins:   config.CORS.Origins,
		AllowedMethods:   config.CORS.Methods,
		AllowCredentials: true,
		AllowedHeaders:   config.CORS.Headers,
		Debug:            config.Debug,
	}))

	n.UseFunc(handlers.Logging)
	n.Use(handlers.Context(ctx))
	n.Use(&negroni.Recovery{
		Logger:     negroni.ALogger(log),
		PrintStack: true,
		StackAll:   true,
	})

	n.Use(handlers.Alive("/"))
	n.UseFunc(handlers.TrackErrors)
	n.UseHandler(api)

	s := &Server{
		Context: ctx,
		api:     api,
		config:  config,
		server: &http.Server{
			Addr:    config.Addr,
			Handler: n,
		},
	}

	return s, nil
}

type Server struct {
	context.Context
	config configuration.HTTPConfig
	server *http.Server
	api    *handlers.Api
}

func (server *Server) ListenAndServe() error {
	config := server.config
	ln, err := net.Listen("tcp", config.Addr)
	if err != nil {
		return err
	}

	acontext.GetLogger(server).Infof("listening on %v", ln.Addr())
	return server.server.Serve(ln)
}
