package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/gost-codes/sweet_dreams/db/sqlc"
	"github.com/gost-codes/sweet_dreams/token"
)

const (
	MaternityNursing = "maternityNursing"
	GiftPackage      = "giftPackage"
)

type createBookingReq struct {
	Type          string    `json:"type" binding:"required,bookingType"`
	PreferedNurse *int64    `json:"prefered_nurse"`
	StartDate     time.Time `json:"start_date" binding:"required"`
	EndDate       time.Time `json:"end_date" binding:"required"`
	Long          float64   `json:"long" binding:"required"`
	Lat           float64   `json:"lat" binding:"required"`
}

func (server *Server) createBooking(ctx *gin.Context) {
	authUser := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	if authUser == nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("Something went wrong")))
		return
	}

	var req createBookingReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	args := db.CreateBookingRequestParams{
		UserID:        authUser.UserID,
		Type:          req.Type,
		PreferedNurse: req.PreferedNurse,
		StartDate:     req.StartDate,
		EndDate:       req.EndDate,
		Point:         req.Lat,
		Point_2:       req.Long,
	}

	bookingReq, err := server.store.CreateBookingRequest(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, bookingReq)
}

type bookingRequesParams struct {
	Count int32 `form:"count" binding:"default=20,min=5"`
	Page  int32 `form:"page" binding:"default=1"`
}

func (server *Server) listUserBookingReqs(ctx *gin.Context) {

	var queryParams bookingRequesParams

	userPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	if err := ctx.ShouldBindQuery(&queryParams); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	args := db.ListUserBookingReqsParams{
		UserID: userPayload.UserID,
		Limit:  queryParams.Count,
		Offset: (queryParams.Page - 1) * queryParams.Count,
	}
	reqs, err := server.store.ListUserBookingReqs(ctx, args)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusOK, []int{})
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, reqs)
}

type bookingReqID struct {
	ID int64 `uri:"id" binding:"required"`
}

func (server *Server) bookingReq(ctx *gin.Context) {

	var params bookingReqID

	userPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	if err := ctx.ShouldBindUri(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	args := db.GetBookingByIDParams{
		UserID: userPayload.UserID,
		ID:     params.ID,
	}
	reqs, err := server.store.GetBookingByID(ctx, args)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, genericResponse("Booking request not found"))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, reqs)
}

func (server *Server) deleteBookingReq(ctx *gin.Context) {

	var params bookingReqID

	userPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	if err := ctx.ShouldBindUri(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	args := db.DeleteBookingByIDParams{
		UserID: userPayload.UserID,
		ID:     params.ID,
	}
	err := server.store.DeleteBookingByID(ctx, args)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, genericResponse("Booking request not found"))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, genericResponse("Booking successfully deleted"))
}

func (server *Server) adminBookings(ctx *gin.Context) {
	bookingReqs, err := server.store.GetAllBookingsByAdmin(ctx)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusOK, []int{})
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, bookingReqs)
}

func (server *Server) adminBookingsByID(ctx *gin.Context) {
	var req bookingReqID

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	bookingReq, err := server.store.GetBookingsByAdminByID(ctx, req.ID)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(fmt.Errorf("Booking not found")))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, bookingReq)
}

type approveBooking struct {
	RequestID     int64   `json:"request_id" binding:"required"`
	AssignedNurse int64   `json:"assigned_nurse" binding:"required"`
	UserID        int64   `json:"user_id" binding:"required"`
	ApprovedBy    int64   `json:"approved_by" binding:"required"`
	Status        string  `json:"status" binding:"required"`
	Notes         *string `json:"notes"`
}

func (server *Server) adminApproveBooking(ctx *gin.Context) {
	var req approveBooking

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	_, err := server.store.GetBookingByID(ctx, db.GetBookingByIDParams{
		UserID: req.UserID,
		ID:     req.RequestID,
	})

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	adminPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	args := db.CreateApprovalParams{
		RequestID:     req.RequestID,
		UserID:        req.UserID,
		Notes:         req.Notes,
		AssignedNurse: req.AssignedNurse,
		ApprovedBy:    adminPayload.UserID,
		Status:        req.Status,
	}
	approval, err := server.store.CreateApproval(ctx, args)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, approval)
}

func (server *Server) deleteApprovals(ctx *gin.Context) {
	var param bookingReqID

	if err := ctx.ShouldBindUri(&param); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := server.store.DeleteApprovals(ctx, param.ID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, genericResponse("Approval successfully deleted"))
}

func (server *Server) userApprovals(ctx *gin.Context) {
	var pagination paginationData

	if err := ctx.ShouldBindQuery(&pagination); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	userPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	args := db.GetUserApprovalsParams{
		UserID: userPayload.UserID,
		Limit:  pagination.Count,
		Offset: (pagination.Page - 1) * pagination.Count,
	}

	approvals, err := server.store.GetUserApprovals(ctx, args)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusOK, []int{})
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, approvals)

}
func (server *Server) adminGetApprovals(ctx *gin.Context) {
	var pagination paginationData

	if err := ctx.ShouldBindQuery(&pagination); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	args := db.AdminGetUserApprovalsParams{

		Limit:  pagination.Count,
		Offset: (pagination.Page - 1) * pagination.Count,
	}

	approvals, err := server.store.AdminGetUserApprovals(ctx, args)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusOK, []int{})
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, approvals)

}
