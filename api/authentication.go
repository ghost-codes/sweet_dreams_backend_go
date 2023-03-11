package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	db "github.com/gost-codes/sweet_dreams/db/sqlc"
	"github.com/gost-codes/sweet_dreams/util"
	"github.com/lib/pq"
)

type createUserWithEmailPasswordReq struct {
	Username  string `json:"username" binding:"required,min=8,alphanum"`
	FirstName string `json:"first_name" binding:"required,alphanum"`
	LastName  string `json:"last_name" binding:"required,alphanum"`
	Email     string `json:"email" binding:"required,email"`
	Contact   string `json:"contact"`
	Password  string `json:"password" binding:"required,min=8"`
}

type UserResponse struct {
	ID                int64      `json:"id"`
	Username          string     `json:"username"`
	FirstName         string     `json:"first_name"`
	LastName          string     `json:"last_name"`
	Email             string     `json:"email"`
	AvatarUrl         *string    `json:"avatar_url"`
	Contact           *string    `json:"contact"`
	PasswordChangedAt time.Time  `json:"password_changed_at"`
	VerifiedAt        *time.Time `json:"verified_at"`
	CreatedAt         time.Time  `json:"created_at"`
}

type loginUserResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	User         UserResponse `json:"user"`
}

func newUserResponse(user db.User) UserResponse {
	return UserResponse{
		Username:          user.Username,
		Email:             user.Email,
		FirstName:         user.FirstName,
		LastName:          user.LastName,
		VerifiedAt:        user.VerifiedAt,
		AvatarUrl:         user.AvatarUrl,
		ID:                user.ID,
		Contact:           user.Contact,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}
}

func (server *Server) createUserWithEmailPassword(ctx *gin.Context) {
	req := createUserWithEmailPasswordReq{}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hasedPassword, err := util.HashedPassword(req.Password)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	args := db.CreateuserParams{
		Username:       req.Username,
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		Email:          req.Email,
		Contact:        &req.Contact,
		HashedPassword: hasedPassword,
		SecurityKey:    uuid.NewString(),
	}

	user, err := server.store.Createuser(ctx, args)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "uniue_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}

		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := newUserResponse(user)

	accessToken, _, err := server.tokenMaker.CreateToken(req.Username, user.SecurityKey, server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return

	}
	refreshToken, _, err := server.tokenMaker.CreateToken(req.Username, user.SecurityKey, server.config.RefreshTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return

	}

	ctx.JSON(http.StatusOK, loginUserResponse{User: res, RefreshToken: refreshToken, AccessToken: accessToken})
}

type loginUserReq struct {
	UsernameEmail string `json:"username_email" binding:"required"`
	Password      string `json:"password" binding:"required"`
}

func (server *Server) loginWithEmailPassword(ctx *gin.Context) {
	req := loginUserReq{}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUser(ctx, req.UsernameEmail)
	if err != nil {
		if err == sql.ErrNoRows {
			err = fmt.Errorf("user does not exist:%v", err)
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if err := util.ComparePassword(req.Password, user.HashedPassword); err != nil {
		err := fmt.Errorf("invalid credentials")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	accessToken, _, err := server.tokenMaker.CreateToken(user.Username, user.SecurityKey, server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	refresh_token, _, err := server.tokenMaker.CreateToken(user.Username, user.SecurityKey, server.config.RefreshTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	userRes := newUserResponse(user)

	ctx.JSON(http.StatusOK, loginUserResponse{RefreshToken: refresh_token, AccessToken: accessToken, User: userRes})

}
