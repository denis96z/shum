package api

import (
	"github.com/gin-gonic/gin"
)

func (srv *API) HandleStatus(ctx *gin.Context) {
	WriteOK(ctx)
}
