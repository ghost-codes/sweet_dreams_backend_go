package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/gost-codes/sweet_dreams/db/sqlc"
	"github.com/gost-codes/sweet_dreams/util"
	"github.com/gost-codes/sweet_dreams/worker"
	"github.com/hibiken/asynq"
	"github.com/lib/pq"
)

type createAdminReq struct {
	Username string `json:"username" binding:"required,min=5"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

func (server *Server) createAdmin(ctx *gin.Context) {
	req := createAdminReq{}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hasedPassword, err := util.HashedPassword(req.Password)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	args := db.CreateAdminTxParams{
		CreateAdminParams: db.CreateAdminParams{
			Username:       req.Username,
			FullName:       req.FullName,
			Email:          req.Email,
			HashedPassword: hasedPassword,
		},
		AfterCreate: func(admin db.Admin) error {
			//TODO: send verification email to client using redis
			opts := []asynq.Option{
				asynq.MaxRetry(10),
				asynq.ProcessIn(10 * time.Second),
				asynq.Queue(worker.CriticalQueue),
			}
			err = server.taskDistributor.DistributeTaskSendAdminEmail(ctx, &worker.PayloadSendAdminEmail{Username: req.Username, Password: req.Password}, opts...)
			if err != nil {

				return fmt.Errorf("failed to distribute send verified email task: %w", err)
			}
			return nil

		},
	}

	//-------------------> TODO: convert to TX
	_, err = server.store.CreateAdminTX(ctx, args)
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

	ctx.JSON(http.StatusOK, map[string]string{"message": "admin successfully created"})
}

type adminLoginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type admin struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	FullName  string    `json:"full_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func Newadmin(Admin db.Admin) admin {
	return admin{
		ID:        Admin.ID,
		Username:  Admin.Username,
		FullName:  Admin.FullName,
		Email:     Admin.Email,
		CreatedAt: Admin.CreatedAt,
	}
}

type adminLoginResponse struct {
	Admin       admin  `json:"admin"`
	AccessToken string `json:"access_token"`
}

func (server *Server) adminLogin(ctx *gin.Context) {
	var req adminLoginReq

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	admin, err := server.store.GetAdmin(ctx, req.Username)

	if err != nil {
		if err == sql.ErrNoRows {
			err := fmt.Errorf("admin with credential %s does not exist", req.Username)
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if err := util.ComparePassword(req.Password, admin.HashedPassword); err != nil {
		err := fmt.Errorf("invalid credentials")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	accessToken, _, err := server.tokenMaker.CreateToken(admin.ID, "", time.Duration(1*time.Hour))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := adminLoginResponse{
		Admin:       Newadmin(admin),
		AccessToken: accessToken,
	}

	ctx.JSON(http.StatusOK, res)

}
