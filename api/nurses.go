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
	Page  int32 `form:"page" binding:"default=1,min=1"`
	Count int32 `form:"count" binding:"default=20,min=5,max=20"`
}

// @Tags         Admin endpoints
// @Security 	bearerAuth
// @Accept       json
// @Produce      json
// @Param      	body 		body 		createNurseReq 	true	 " "
// @Success     200  		string    	"nurse has been created successfully"
// @response    default  	{object}  	ErrorResponse
// @Router      /nurses/create 	[post]
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

// @Tags         Booking requests
// @Security 	bearerAuth
// @Accept       json
// @Produce      json
// @Param      count 		query	int false " "
// @Param      	page 		query	int false " "
// @Success     200  		{array}    db.Nurse
// @response    default  	{object}  	ErrorResponse
// @Router      /nurses 	[get]
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

// @Summary      send verification email
// @Description  short code is sent to the user's email for verification
// @Tags         Booking requests
// @Security 	bearerAuth
// @Accept       json
// @Produce      json
// @Param      	id 		path	int false " "
// @Success     200  		{object} 	db.Nurse
// @response    default  	{object}  	ErrorResponse
// @Router      /nurses/{id}	[get]
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

// @Tags         Admin endpoints
// @Security 	bearerAuth
// @Accept       json
// @Produce      json
// @Param      	id 		path	int false " "
// @Success     200  		string    "nurse deleted successfully"
// @response    default  	{object}  	ErrorResponse
// @Router      /nurses/{id}	[delete]
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
