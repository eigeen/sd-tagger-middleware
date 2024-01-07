package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func ErrBadRequest(ctx *gin.Context, reason string) {
	ctx.JSON(http.StatusBadRequest, ErrResp{
		Error:  "HTTPException",
		Detail: reason,
		Errors: "",
	})
	ctx.Abort()
	return
}

func ErrImageEncoding(ctx *gin.Context, reason string) {
	ctx.JSON(http.StatusInternalServerError, ErrResp{
		Error:  "HTTPException",
		Detail: reason,
		Errors: "",
	})
	ctx.Abort()
	return
}

func ErrInternalServer(ctx *gin.Context, reason string) {
	ctx.JSON(http.StatusInternalServerError, ErrResp{
		Error:  "HTTPException",
		Detail: reason,
		Errors: "",
	})
	ctx.Abort()
	return
}
