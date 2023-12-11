package middleware

import (
	"net/http"

	"rip/internal/pkg/api/consts"

	"github.com/gin-gonic/gin"
)

func WithModeratorCheck(gCtx *gin.Context) {
	if !gCtx.GetBool(consts.ModeratorCtxParam) {
		gCtx.AbortWithStatus(http.StatusForbidden)
		return
	}
}
