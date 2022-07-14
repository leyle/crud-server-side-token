package sstapp

import (
	"github.com/gin-gonic/gin"
	"github.com/leyle/go-api-starter/ginhelper"
	"github.com/rs/zerolog"
)

// auth middleware for gin framework

func GinAuthMiddleware(sst *SSTokenOption) gin.HandlerFunc {
	return func(c *gin.Context) {
		// get logger from context
		logger := zerolog.Ctx(c.Request.Context())
		logger.Debug().Msg("try to check sst token...")

		token := c.Request.Header.Get(ServerSideTokenHeaderName)
		if token == "" {
			logger.Warn().Msgf("no [%s] in request headers", ServerSideTokenHeaderName)
			ginhelper.Return401Json(c, "No token")
			return
		}

		result := sst.VerifyToken(c.Request.Context(), token)
		if !result.OK {
			logger.Warn().Msgf("invalid token[%s], %s", token, result.Msg)
			ginhelper.Return401Json(c, result.Msg)
			return
		}
		logger.Debug().Str("token", token).Str("userId", result.Msg).Msg("sst token is valid")

		c.Next()
	}
}
