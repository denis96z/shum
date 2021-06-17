package api

import (
	"bytes"
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

func WriteErrWithAndOutputBody(ctx *gin.Context, err error, stdout, stderr bytes.Buffer) {
	ctx.JSON(http.StatusInternalServerError, makeErrWithOutputBody(err, stdout, stderr))
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

type errWithOutputBody struct {
	err    error
	stdout bytes.Buffer
	stderr bytes.Buffer
}

func makeErrWithOutputBody(err error, stdout, stderr bytes.Buffer) errWithOutputBody {
	return errWithOutputBody{
		err: err, stdout: stdout, stderr: stderr,
	}
}

func (v errWithOutputBody) MarshalJSON() ([]byte, error) {
	b := []byte(fmt.Sprintf(
		`{"status":"error","message":%q,"stdout":%q,"stderr":%q}`,
		v.err.Error(), v.stdout.String(), v.stderr.String(),
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
