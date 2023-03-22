package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	db "github.com/gost-codes/sweet_dreams/db/sqlc"
	"github.com/gost-codes/sweet_dreams/token"
	"github.com/gost-codes/sweet_dreams/util"
	"github.com/gost-codes/sweet_dreams/worker"
	"github.com/hibiken/asynq"
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

// @Summary      signup new user
// @Description  create new user using email and password
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        body    body   	createUserWithEmailPasswordReq  true " "
// @Success      200  	{object}   	loginUserResponse
// @response     default  {object}  	ErrorResponse
// @Router       /signup [post]
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

	args := db.CreateUserTxParams{

		CreateuserParams: db.CreateuserParams{
			Username:       req.Username,
			FirstName:      req.FirstName,
			LastName:       req.LastName,
			Email:          req.Email,
			Contact:        &req.Contact,
			HashedPassword: &hasedPassword,
			SecurityKey:    uuid.NewString(),
		},
		AterTx: func(user db.User) error {
			//TODO: send verification email to client using redis
			opts := []asynq.Option{
				asynq.MaxRetry(10),
				asynq.ProcessIn(10 * time.Second),
				asynq.Queue(worker.CriticalQueue),
			}
			err = server.taskDistributor.DistributeTaskSendVerifyEmail(ctx, &worker.PayloadSendVerifyEmail{UserID: user.ID}, opts...)
			if err != nil {

				return fmt.Errorf("failed to distribute send verified email task: %w", err)
			}
			return nil

		},
	}

	//-------------------> TODO: convert to TX
	createUserTxResult, err := server.store.CreateUserTx(ctx, args)
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

	//------------------>
	res := newUserResponse(createUserTxResult.User)

	accessToken, _, err := server.tokenMaker.CreateToken(createUserTxResult.User.ID, createUserTxResult.User.SecurityKey, server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return

	}
	refreshToken, _, err := server.tokenMaker.CreateToken(createUserTxResult.User.ID, createUserTxResult.User.SecurityKey, server.config.RefreshTokenDuration)
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

// @Summary      log existing user in
// @Description  log existing users in with email and password
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        body    body   	loginUserReq  true " "
// @Success      200  	{object}   	loginUserResponse
// @response     default  {object}  	ErrorResponse
// @Router       /login [post]
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

	if user.HashedPassword == nil {
		err := fmt.Errorf("password has not been set for this account with email, check email for an OTP to set password, or sign in with your preferred social")
		//TODO: setup email service to send otp
		//TODO: setup set or reset password
		ctx.JSON(http.StatusNotAcceptable, errorResponse(err))
		return
	}

	if err := util.ComparePassword(req.Password, *user.HashedPassword); err != nil {
		err := fmt.Errorf("invalid credentials")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	accessToken, _, err := server.tokenMaker.CreateToken(user.ID, user.SecurityKey, server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	refresh_token, _, err := server.tokenMaker.CreateToken(user.ID, user.SecurityKey, server.config.RefreshTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	userRes := newUserResponse(user)

	ctx.JSON(http.StatusOK, loginUserResponse{RefreshToken: refresh_token, AccessToken: accessToken, User: userRes})
}

type signInUserSocialReq struct {
	TokenId string `json:"token_id" binding:"required"`
}

func (server *Server) signInUserSocial(ctx *gin.Context) {
	req := signInUserSocialReq{}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	//Authenticate social tokenId from firebase
	client, err := server.firebase.Auth(ctx)
	if err != nil {
		err = fmt.Errorf("error getting Auth client: %v\n", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	token, err := client.VerifyIDToken(ctx, req.TokenId)
	if err != nil {
		err = fmt.Errorf("user unauthenticated: %v\n", err)
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	userRecord, err := client.GetUser(ctx, token.UID)

	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	//check if user already in database
	user, err := server.store.GetUser(ctx, userRecord.Email)
	if err != nil {
		if err != sql.ErrNoRows {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		//if not create user
		names := strings.Fields(userRecord.DisplayName)
		firstName := names[0]
		lastName := ""
		if len(names) > 1 {
			lastName = strings.Join(names[1:len(names)-1], " ")
		}

		if userRecord.ProviderID != "twitter.com" && userRecord.ProviderID != "google.com" && userRecord.ProviderID != "apple.com" {
			err := fmt.Errorf("social:%v not supported", userRecord.ProviderID)
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}

		arg := db.CreateuserParams{
			Username:      strings.Join(strings.Fields(userRecord.DisplayName), ""),
			FirstName:     firstName,
			LastName:      lastName,
			Email:         userRecord.Email,
			Contact:       &userRecord.PhoneNumber,
			AvatarUrl:     &userRecord.PhotoURL,
			SecurityKey:   uuid.NewString(),
			TwitterSocial: userRecord.ProviderID == "twitter.com",
			GoogleSocial:  userRecord.ProviderID == "google.com",
			AppleSocial:   userRecord.ProviderID == "apple.com",
		}
		user, err = server.store.Createuser(ctx, arg)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}
	//if yes check if appropriate social has been checked and update if needed
	if user.AppleSocial == false && user.GoogleSocial == false && user.TwitterSocial == false {

		switch userRecord.ProviderID {
		case "twitter.com":
			user.TwitterSocial = true
		case "google.com":
			user.GoogleSocial = true
		case "apple.com":
			user.AppleSocial = true
		}
		arg := db.UpdateUserParams{
			ID:             user.ID,
			FirstName:      user.FirstName,
			LastName:       user.LastName,
			Username:       user.Username,
			AvatarUrl:      user.AvatarUrl,
			Email:          user.Email,
			Contact:        user.Contact,
			HashedPassword: user.HashedPassword,
			SecurityKey:    user.SecurityKey,
			CreatedAt:      user.CreatedAt,
			VerifiedAt:     &user.PasswordChangedAt,
			TwitterSocial:  user.TwitterSocial,
			AppleSocial:    user.AppleSocial,
			GoogleSocial:   user.GoogleSocial,
		}

		user, err = server.store.UpdateUser(ctx, arg)
	}

	//then login to proceed
	accessToken, _, err := server.tokenMaker.CreateToken(user.ID, user.SecurityKey, server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	refresh_token, _, err := server.tokenMaker.CreateToken(user.ID, user.SecurityKey, server.config.RefreshTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	userRes := newUserResponse(user)

	ctx.JSON(http.StatusOK, loginUserResponse{RefreshToken: refresh_token, AccessToken: accessToken, User: userRes})
}

type verifyEmailReq struct {
	ID   int64  `form:"id" binding:"required"`
	Code string `form:"code" binding:"required,min=64"`
}

// @Summary      verify email
// @Description  verify email using short code
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        body    body   	loginUserReq  true " "
// @Success      200  	{object}    UserResponse
// @response     default  {object}  	ErrorResponse
// @Router       /verify_email [post]
func (server *Server) verifyEmail(ctx *gin.Context) {
	var req verifyEmailReq

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	verifyEmail, err := server.store.GetVerifyEmail(ctx, req.ID)
	if err != nil || verifyEmail.SecretKey != req.Code {
		err := fmt.Errorf("invalid verification link:%w", err)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if time.Now().After(verifyEmail.ExpiredAt) {
		err := fmt.Errorf("verification link has expired")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// replace user's email when user data email is different from verify_email email
	user, err := server.store.GetUser(ctx, *verifyEmail.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	now := time.Now()
	arg := db.UpdateUserParams{
		ID:                user.ID,
		Username:          user.Username,
		Email:             verifyEmail.Email,
		FirstName:         user.FirstName,
		LastName:          user.LastName,
		HashedPassword:    user.HashedPassword,
		AvatarUrl:         user.AvatarUrl,
		Contact:           user.Contact,
		SecurityKey:       uuid.NewString(),
		VerifiedAt:        &now,
		TwitterSocial:     user.TwitterSocial,
		AppleSocial:       user.AppleSocial,
		GoogleSocial:      user.GoogleSocial,
		CreatedAt:         user.CreatedAt,
		PasswordChangedAt: user.PasswordChangedAt,
	}

	updatedUser, err := server.store.UpdateUser(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, newUserResponse(updatedUser))
}

// @Summary      send verification email
// @Description  short code is sent to the user's email for verification
// @Tags         Authentication
// @Security 	bearerAuth
// @Accept       json
// @Produce      json
// @Success      200  	string    "verification email successfully sent"
// @response     default  {object}  	ErrorResponse
// @Router       /send_verification_email [get]
func (server *Server) sendVerificationEmail(ctx *gin.Context) {
	payload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	err := server.taskDistributor.DistributeTaskSendVerifyEmail(ctx, &worker.PayloadSendVerifyEmail{
		UserID: payload.UserID,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, map[string]string{"success": "verification email successfully sent"})
}
