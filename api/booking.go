package api

import (
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
