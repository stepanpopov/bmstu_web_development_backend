package middleware

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"rip/internal/pkg/api/consts"
	myRedis "rip/internal/pkg/redis"
	"rip/internal/pkg/repo"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt"
)

func WithAuthCheck(secret string, redisCl *myRedis.RedisClient) func(gCtx *gin.Context) {
	return func(gCtx *gin.Context) {
		jwtStr := gCtx.GetHeader("Authorization")
		if !strings.HasPrefix(jwtStr, consts.JwtPrefix) {
			gCtx.AbortWithStatus(http.StatusForbidden)
			return
		}
		// отрезаем префикс
		jwtStr = jwtStr[len(consts.JwtPrefix):]

		err := redisCl.CheckJWTInBlacklist(gCtx, jwtStr)
		if err == nil { // значит что токен в блеклисте
			gCtx.AbortWithStatus(http.StatusForbidden)
			return
		}
		if !errors.Is(err, redis.Nil) { // значит что это не ошибка отсуствия - внутренняя ошибка
			gCtx.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		parsedToken, err := jwt.ParseWithClaims(jwtStr, &repo.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})
		if err != nil {
			gCtx.AbortWithStatus(http.StatusForbidden)
			log.Println(err)
			return
		}

		myClaims := parsedToken.Claims.(*repo.JWTClaims)
		gCtx.Set(consts.ModeratorCtxParam, myClaims.IsModerator)
		gCtx.Set(consts.UserUUIDCtxParam, myClaims.UserUUID)
	}
}
