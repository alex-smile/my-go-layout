package api

import (
	"errors"

	"github.com/gin-gonic/gin"

	"mygo/template/pkg/errorx"
	"mygo/template/pkg/store"
	"mygo/template/pkg/util"
)

func HandleError(c *gin.Context, err error) {
	if errors.Is(err, store.ErrRecordNotFound) {
		util.NotFoundJSONResponse(c, err.Error())
		return
	}

	var myError *errorx.MygoError
	if errors.As(err, &myError) {
		util.BaseErrorJSONResponse(c, string(myError.Code), err.Error(), myError.StatusCode)
		return
	}

	util.SystemErrorJSONResponse(c, err)
}
