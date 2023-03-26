package router

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/haandol/devops-academy-eda-demo/pkg/constant"
	"github.com/haandol/devops-academy-eda-demo/pkg/port/primaryport/routerport"
	"github.com/haandol/devops-academy-eda-demo/pkg/service"
	"github.com/haandol/devops-academy-eda-demo/pkg/util"
	"github.com/haandol/devops-academy-eda-demo/pkg/util/cerrors"
)

type CreateTripRequest struct {
	TripID string `json:"tripId" binding:"required" validate:"required"`
}

type TripRouter struct {
	BaseRouter
	tripService *service.TripService
}

func NewTripRouter(
	tripService *service.TripService,
) *TripRouter {
	return &TripRouter{
		tripService: tripService,
	}
}

func (r *TripRouter) Route(rg routerport.RouterGroup) {
	g := rg.Group("/trips")
	g.Handle("POST", "/", r.WrappedHandler(r.CreateHandler))
	g.Handle("GET", "/", r.WrappedHandler(r.ListHandler))
	g.Handle("GET", "/hotels/error", r.WrappedHandler(r.GetInjectionStatusHandler))
	g.Handle("POST", "/hotels/error", r.WrappedHandler(r.InjectErrorHandler))
}

// @Summary create new trip
// @Schemes
// @Description create new trip
// @Tags trips
// @Accept json
// @Produce json
// @Param trip body dto.Trip true "trip id is required"
// @Success 200 {object} dto.Trip
// @Router /trips [post]
func (r *TripRouter) CreateHandler(c *gin.Context) *cerrors.CodedError {
	req := &CreateTripRequest{}
	if err := c.ShouldBindJSON(req); err != nil {
		return cerrors.New(constant.ErrBadRequest, err)
	}

	if err := util.ValidateStruct(req); err != nil {
		return cerrors.New(constant.ErrInvalidRequest, err)
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Second*10)
	defer cancel()

	trip, err := r.tripService.Create(ctx, req.TripID)
	if err != nil {
		return cerrors.New(constant.ErrFailToCreateTrip, err)
	}

	return r.Success(c, trip)
}

// @Summary list all trips
// @Schemes
// @Description list all trips
// @Tags trips
// @Accept json
// @Produce json
// @Success 200 {object} []dto.Trip
// @Router /trips [get]
func (r *TripRouter) ListHandler(c *gin.Context) *cerrors.CodedError {
	trips, err := r.tripService.List(c.Request.Context())
	if err != nil {
		return cerrors.New(constant.ErrFailToListTrip, err)
	}

	return r.Success(c, trips)
}

// @Summary get error injection status from hotel service
// @Schemes
// @Description get error injection status from hotel service
// @Tags trips
// @Accept json
// @Produce json
// @Success 200 {object} bool
// @Router /trips/hotels/error [get]
func (r *TripRouter) GetInjectionStatusHandler(c *gin.Context) *cerrors.CodedError {
	injectionStatus, err := r.tripService.GetInjectionStatus(c.Request.Context())
	if err != nil {
		return cerrors.New(constant.ErrInjectionError, err)
	}

	return r.Success(c, injectionStatus)
}

// @Summary inject error to hotel service
// @Schemes
// @Description inject error to hotel service
// @Tags trips
// @Accept json
// @Produce json
// @Success 200 {object} bool
// @Router /trips/hotels/error [post]
func (r *TripRouter) InjectErrorHandler(c *gin.Context) *cerrors.CodedError {
	injectionStatus, err := r.tripService.InjectErrorHandler(c.Request.Context())
	if err != nil {
		return cerrors.New(constant.ErrInjectionError, err)
	}

	return r.Success(c, injectionStatus)
}
