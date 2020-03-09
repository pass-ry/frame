package server

import (
	"context"
	"net/http"
	"time"

	"gitlab.ifchange.com/data/cordwood/log"
	router "gitlab.ifchange.com/data/cordwood/rpc/rpc-router"
)

type Server struct {
	*http.Server
	ctx      context.Context
	isTLS    bool
	certFile string
	keyFile  string
}

func NewServer(ctx context.Context, addr string, r *router.Router) *Server {
	server := &Server{
		Server: new(http.Server),
		ctx:    ctx,
	}
	server.Server.Addr = addr
	server.Server.Handler = r
	return server
}

func NewServerTLS(ctx context.Context, addr string, r *router.Router, certFile, keyFile string) *Server {
	server := &Server{
		Server: new(http.Server),
		ctx:    ctx,
	}
	server.Server.Addr = addr
	server.Server.Handler = r

	server.isTLS = true
	server.certFile = certFile
	server.keyFile = keyFile
	return server
}

func (server *Server) GraceRun() {
	go func() {
		<-server.ctx.Done()
		ctx, cancel := context.WithTimeout(context.Background(),
			time.Duration(10)*time.Second)
		defer cancel()
		err := server.Shutdown(ctx)
		log.Infof("HTTPD %s is gracefully stopped %v",
			server.Addr, err)
	}()

	var err error
	if server.isTLS {
		err = server.ListenAndServeTLS(server.certFile, server.keyFile)
	} else {
		err = server.ListenAndServe()
	}

	log.Infof("HTTPD %s is stopped %v",
		server.Addr, err)
}
