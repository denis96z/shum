package api

import (
	"fmt"
	"net/http"

	"github.com/maildealru/shum/pkg/shum/conf"

	"github.com/gin-gonic/gin"
)

func WriteOK(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, okBody{})
}

func WriteOKWithBody(ctx *gin.Context, body interface{}) {
	ctx.JSON(http.StatusOK, body)
}

func WriteFailWithBody(ctx *gin.Context, body interface{}) {
	ctx.JSON(http.StatusUnprocessableEntity, body)
}

func WriteErrWithBody(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusInternalServerError, makeErrBody(err))
}

type okBody struct{}

func (v okBody) MarshalJSON() ([]byte, error) {
	return []byte(`{"status":"ok"}`), nil
}

type errBody struct {
	err error
}

func makeErrBody(err error) errBody {
	return errBody{err: err}
}

func (v errBody) MarshalJSON() ([]byte, error) {
	b := []byte(fmt.Sprintf(
		`{"status":"error","message":%q}`, v.err.Error(),
	))
	return b, nil
}

func WithAuth(handler func(ctx *gin.Context), authConf conf.AuthConfig) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		id, secret, ok := ctx.Request.BasicAuth()
		if !ok || !authConf.IsAuthOK(id, secret) {
			ctx.JSON(http.StatusUnauthorized, makeUnauthorizedBody())
			return
		}
		handler(ctx)
	}
}

type unauthorizedBody struct{}

func makeUnauthorizedBody() unauthorizedBody {
	return unauthorizedBody{}
}

func (v unauthorizedBody) MarshalJSON() ([]byte, error) {
	return []byte(`{"status":"fail","message":"unauthorized"}`), nil
}
