package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/maildealru/shum/pkg/shum/conf"

	"github.com/gin-gonic/gin"
)

type API struct {
	config *conf.Config
	server *http.Server
}

func NewAPI(conf *conf.Config) *API {
	srv := &API{
		config: conf,
	}

	r := gin.Default()

	r.GET("/status", srv.HandleStatus)
	r.POST("/cmd/:name", WithAuth(srv.HandleCmd, conf.Auth))

	srv.server = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", conf.Server.Addr, conf.Server.Port),
		Handler: r,
	}

	return srv
}

func (srv *API) Start() error {
	if srv.config.Server.TLS.Enabled {
		return srv.server.ListenAndServeTLS(
			srv.config.Server.TLS.CertPath, srv.config.Server.TLS.KeyPath,
		)
	}
	return srv.server.ListenAndServe()
}

func (srv *API) Stop() error {
	ctx, cancel := context.WithTimeout(
		context.Background(), srv.config.Server.Shutdown.Timeout,
	)
	defer cancel()

	return srv.server.Shutdown(ctx)
}
