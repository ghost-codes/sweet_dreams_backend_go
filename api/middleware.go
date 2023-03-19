package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	db "github.com/gost-codes/sweet_dreams/db/sqlc"
	"github.com/gost-codes/sweet_dreams/token"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "payload"
)

func authMiddleware(tokenMaker token.Maker, store *db.Store) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorization := ctx.GetHeader(authorizationHeaderKey)

		if len(authorization) == 0 {
			err := errors.New("Authorization header not provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		fields := strings.Fields(authorization)
		if len(fields) < 2 {
			err := errors.New("invalid authorization header format")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		if authorizationTypeBearer != strings.ToLower(fields[0]) {
			err := errors.New("invalid authorization header type")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		accessToken := fields[1]

		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		user, err := store.GetUser(ctx, payload.Username)
		if err != nil {
			if err == sql.ErrNoRows {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, map[string]string{"error": fmt.Sprintf("user with username %s not found", payload.Username)})
				return
			}
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, map[string]string{"error": "error occured"})
			return
		}

		if payload.SecurityKey != user.SecurityKey {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, map[string]string{"error": "stale credentials. please login again"})
			return
		}

		ctx.Set(authorizationPayloadKey, payload)
		ctx.Next()
	}
}

func adminAuthMiddleware(tokenMaker token.Maker, store *db.Store, is_super_check bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorization := ctx.GetHeader(authorizationHeaderKey)

		if len(authorization) == 0 {
			err := errors.New("Authorization header not provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		fields := strings.Fields(authorization)
		if len(fields) < 2 {
			err := errors.New("invalid authorization header format")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		if authorizationTypeBearer != strings.ToLower(fields[0]) {
			err := errors.New("invalid authorization header type")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		accessToken := fields[1]

		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		if is_super_check {
			admin, err := store.GetAdmin(ctx, payload.Username)
			if err != nil {
				if err == sql.ErrNoRows {
					ctx.AbortWithStatusJSON(http.StatusUnauthorized, map[string]string{"error": fmt.Sprintf("user with username %s not found", payload.Username)})
					return
				}
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, map[string]string{"error": "error occured"})
				return
			}

			if !admin.IsSuper {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("only super admins are allowed this fuction")))
				return
			}
		}

		ctx.Set(authorizationPayloadKey, payload)
		ctx.Next()
	}
}
