package main

import (
	"github.com/leyle/go-api-starter/ginhelper"
	"github.com/leyle/server-side-token/sstapp"
)

func Auth(ctx *AppOption) {
	token := ctx.C.Request.Header.Get(sstapp.ServerSideTokenHeaderName)
	if token == "" {
		ctx.Logger.Warn().Msgf("no [%s] in request headers", sstapp.ServerSideTokenHeaderName)
		ginhelper.Return401Json(ctx.C, "No token")
		return
	}

	result := ctx.SST.VerifyToken(ctx.C.Request.Context(), token)
	if !result.OK {
		ctx.Logger.Warn().Msgf("invalid token[%s], %s", token, result.Msg)
		ginhelper.Return401Json(ctx.C, result.Msg)
		return
	}

	// token is valid
	ctx.Logger.Debug().Str("token", token).Str("userId", result.Msg).Msg("token is valid")

	// check role and permission, if needed

	// save current user info into context

	ctx.C.Next()
}
