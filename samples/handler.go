package main

import (
	"github.com/gin-gonic/gin"
	"github.com/leyle/go-api-starter/ginhelper"
)

func GetHandler(ctx *AppOption) {
	ctx.Logger.Info().Msg("called GetHandler")
	ginhelper.ReturnOKJson(ctx.C, "")
}

func PostHandler(ctx *AppOption) {
	ctx.Logger.Info().Msg("called PostHandler")
	ginhelper.ReturnOKJson(ctx.C, "")
}

type CreateSSTForm struct {
	UserId string `json:"userId"`
}

func CreateSSTHandler(ctx *AppOption) {
	var form CreateSSTForm
	err := ctx.C.BindJSON(&form)
	ginhelper.StopExec(err)

	if form.UserId == "" {
		ginhelper.Return400Json(ctx.C, 4000, "userId cannot be empty string")
		return
	}

	result := ctx.SST.GenerateToken(ctx.C.Request.Context(), form.UserId)
	ginhelper.StopExec(result.Err)

	resp := gin.H{
		"token":  result.Token,
		"userId": form.UserId,
	}

	ginhelper.ReturnOKJson(ctx.C, resp)
	return
}
