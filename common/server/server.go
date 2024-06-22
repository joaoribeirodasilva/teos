package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joaoribeirodasilva/teos/common/logger"
)

type Server struct {
	addr    string
	port    int
	Service *gin.Engine
}

func New(addr string, port int) *Server {

	s := new(Server)
	s.addr = addr
	s.port = port
	s.Service = gin.Default()

	return s
}

func (s *Server) Listen() *logger.HttpError {

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", s.addr, s.port),
		Handler: s.Service.Handler(),
	}
	logger.Info("starting server at %s", srv.Addr)

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error(logger.LogStatusInternalServerError, nil, "listening", err, nil)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("shuting down server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error(logger.LogStatusInternalServerError, nil, "shuting down server", err, nil)
	}
	// catching ctx.Done(). timeout of 5 seconds.
	select {
	case <-ctx.Done():
		logger.Info("wait for 5 seconds...")
	}
	logger.Info("server terminated")

	return nil
}
