package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"rip/internal/pkg/redis"
	"rip/internal/pkg/repo"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type loginReq struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type loginResp struct {
	ExpiresIn   time.Duration `json:"expires_in"`
	AccessToken string        `json:"access_token"`
	TokenType   string        `json:"token_type"`
}

// @Summary		Login
// @Tags		Auth
// @Description	Login account
// @Accept		json
// @Produce		json
// @Param		userInput	body		loginReq		true	"username and password"
// @Success		200			{object}	loginResp				"User created"
// @Failure		400			{object}	error				"Incorrect input"
// @Failure		500			{object}	error				"Server error"
// @Router		/api/auth/login [post]
func login(r repo.Repository, secret string, jwtExpiresIn time.Duration) func(c *gin.Context) {
	return func(gCtx *gin.Context) {
		req := &loginReq{}

		err := json.NewDecoder(gCtx.Request.Body).Decode(req)
		if err != nil {
			gCtx.AbortWithError(http.StatusBadRequest, err)
			return
		}

		hash := generateHashString(req.Password)
		userUUID, isModerator, err := r.CheckUser(req.Login, hash)

		if err != nil {
			gCtx.AbortWithStatus(http.StatusForbidden) // отдаем 403 ответ в знак того что доступ запрещен
			return
		}

		// значит проверка пройдена
		// генерируем ему jwt
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, &repo.JWTClaims{
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(jwtExpiresIn).Unix(),
				IssuedAt:  time.Now().Unix(),
				Issuer:    "rust-admin",
			},
			UserUUID:    userUUID,
			IsModerator: isModerator,
		})
		if token == nil {
			gCtx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("token is nil"))
			return
		}

		strToken, err := token.SignedString([]byte(secret))
		if err != nil {
			gCtx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("cant create str token"))
			return
		}

		gCtx.JSON(http.StatusOK, loginResp{
			ExpiresIn:   jwtExpiresIn,
			AccessToken: strToken,
			TokenType:   "Bearer",
		})
	}
}

type registerReq struct {
	Login       string `json:"login"`
	Password    string `json:"password"`
	IsModerator bool   `json:"is_moderator"`
}

type registerResp struct {
	UserUUID uuid.UUID `json:"user_uuid"`
}

// @Summary		Register
// @Tags		Auth
// @Description	Create account
// @Accept		json
// @Produce		json
// @Param		user	body		registerReq	true	"User info"
// @Success		200		{object}	registerResp		"User created"
// @Failure		400		{object}	error			"Incorrect input"
// @Failure		500		{object}	error			"Server error"
// @Router		/api/auth/register [post]
func register(r repo.Repository) func(c *gin.Context) {
	return func(gCtx *gin.Context) {
		req := &registerReq{}

		err := json.NewDecoder(gCtx.Request.Body).Decode(req)
		if err != nil {
			gCtx.AbortWithError(http.StatusBadRequest, err)
			return
		}

		userUUID, err := r.CreateUser(req.Login, generateHashString(req.Password), req.IsModerator)
		if err != nil {
			gCtx.AbortWithError(http.StatusBadRequest, err)
			return
		}

		gCtx.JSON(http.StatusOK, registerResp{
			UserUUID: userUUID,
		})
	}
}

// @Summary		Logout
// @Tags		Auth
// @Description	Logout
// @Accept		json
// @Produce		json
// @Success		200		{object}	string				"Logout success"
// @Failure		400		{object}	error			"Incorrect input"
// @Failure		500		{object}	error			"Server error"
// @Router		/api/auth/logout [get]
func logout(r repo.Repository, jwtExpiresIn time.Duration, redisCl *redis.RedisClient) func(c *gin.Context) {
	return func(gCtx *gin.Context) {
		jwtStr := getJWTStr(gCtx)

		// сохраняем в блеклист редиса
		err := redisCl.WriteJWTToBlacklist(gCtx.Request.Context(), jwtStr, jwtExpiresIn)
		if err != nil {
			gCtx.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		gCtx.Status(http.StatusOK)
	}
}
