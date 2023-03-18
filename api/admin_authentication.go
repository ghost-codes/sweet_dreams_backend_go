package api

import (
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
