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
	"github.com/joaoribeirodasilva/teos/common/conf"
	"github.com/joaoribeirodasilva/teos/common/database"
	"github.com/joaoribeirodasilva/teos/common/service_log"
)

type Server struct {
	conf    *conf.Conf
	Service *gin.Engine
}

func New(db *database.Db, conf *conf.Conf) *Server {

	s := new(Server)
	s.conf = conf
	s.Service = gin.Default()

	return s
}

func (s *Server) Listen() error {

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", s.conf.Service.BindAddress, s.conf.Service.BindPort),
		Handler: s.Service.Handler(),
	}

	service_log.Info("COMMON::SERVER::Listen", "starting server at %s", srv.Addr)
	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			service_log.Error(0, http.StatusBadRequest, "COMMON::SERVER::Listen", "", "listenning. ERR: %s", err.Error())
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
	service_log.Info("COMMON::SERVER::Listen", "shuting down server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		service_log.Error(0, http.StatusBadRequest, "COMMON::SERVER::Listen", "", "shutingdon server. ERR: %s", err.Error())
	}
	// catching ctx.Done(). timeout of 5 seconds.
	select {
	case <-ctx.Done():
		service_log.Info("COMMON::SERVER::Listen", "wait for 5 seconds...")
	}
	service_log.Info("COMMON::SERVER::Listen", "server terminated")

	return nil
}
