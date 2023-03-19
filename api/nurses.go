package api

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/gost-codes/sweet_dreams/db/sqlc"
)

type createNurseReq struct {
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Contact  string `json:"contact" binding:"required"`
}

type paginationData struct {
	Page  int32 `form:"page" binding:"required,min=1"`
	Count int32 `form:"count" binding:"required,min=5,max=20"`
}

func (server *Server) createNurse(ctx *gin.Context) {
	var req createNurseReq

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	args := db.CreateNurseParams{
		Email:          req.Email,
		FullName:       req.FullName,
		Contact:        req.Contact,
		ProfilePicture: nil, //TODO: add s3 bucket functionality to server
	}

	_, err := server.store.CreateNurse(ctx, args)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, genericResponse("nurse has be created successfully"))
}

func (server *Server) fetchNurses(ctx *gin.Context) {
	var pagination paginationData

	if err := ctx.ShouldBindQuery(&pagination); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	nurses, err := server.store.ListNurses(ctx, db.ListNursesParams{Limit: pagination.Count,
		Offset: (pagination.Page - 1) * pagination.Count})

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusOK, []int{})
		} else {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		}
		return
	}

	ctx.JSON(http.StatusOK, nurses)
}

func (server *Server) fetchNurse(ctx *gin.Context) {
	type idParam struct {
		ID int64 `json:"id" binding:"required"`
	}
	var id idParam

	if err := ctx.ShouldBindUri(&id); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	nurse, err := server.store.GetNurse(ctx, id.ID)

	if err != nil {
		if err == sql.ErrNoRows {
			err = fmt.Errorf("nurse with id %v not found", id.ID)
			ctx.JSON(http.StatusNotFound, errorResponse(err))
		} else {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		}
		return
	}

	ctx.JSON(http.StatusOK, nurse)
}

func (server *Server) deleteNurse(ctx *gin.Context) {
	type idParam struct {
		ID int64 `json:"id" binding:"required"`
	}
	var id idParam

	if err := ctx.ShouldBindUri(&id); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := server.store.DeleteNurse(ctx, id.ID)

	if err != nil {
		if err == sql.ErrNoRows {
			err = fmt.Errorf("nurse with id %v not found", id.ID)
			ctx.JSON(http.StatusNotFound, errorResponse(err))
		} else {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		}
		return
	}

	ctx.JSON(http.StatusOK, genericResponse("nurse deleted successfully"))
}
