package router

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/haandol/devops-academy-eda-demo/pkg/constant"
	"github.com/haandol/devops-academy-eda-demo/pkg/port/primaryport/routerport"
	"github.com/haandol/devops-academy-eda-demo/pkg/service"
	"github.com/haandol/devops-academy-eda-demo/pkg/util/cerrors"
	"github.com/haandol/devops-academy-eda-demo/pkg/util/o11y"
)

type HotelRouter struct {
	BaseRouter
	hotelService *service.HotelService
}

func NewHotelRouter(
	hotelService *service.HotelService,
) *HotelRouter {
	return &HotelRouter{
		hotelService: hotelService,
	}
}

func (r *HotelRouter) Route(rg routerport.RouterGroup) {
	g := rg.Group("/hotels")
	g.Handle("PUT", "/error/", r.WrappedHandler(r.InjectErrorHandler))
	g.Handle("GET", "/error/", r.WrappedHandler(r.GetErrorFlagHandler))
}

// @Summary toggle error injection
// @Schemes
// @Description toggle error injection on creating hotel
// @Tags hotels
// @Accept json
// @Produce json
// @Success 200 {object} bool
// @Router /hotels/error [put]
func (r *HotelRouter) InjectErrorHandler(c *gin.Context) *cerrors.CodedError {
	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Second*10)
	defer cancel()

	span := o11y.SpanFromContext(ctx)

	var req struct {
		Flag bool `json:"flag"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		return cerrors.New(constant.ErrBadRequest, err)
	}
	span.SetAttributes(
		o11y.AttrString("Flag", fmt.Sprintf("%v", req.Flag)),
	)

	if err := r.hotelService.InjectError(ctx, req.Flag); err != nil {
		span.RecordError(err)
		span.SetStatus(o11y.GetStatus(err))
		return cerrors.New(constant.ErrFailToCreateTrip, err)
	}

	return r.Success(c, req.Flag)
}

// @Summary check error flag
// @Schemes
// @Description return error flag
// @Tags trips
// @Accept json
// @Produce json
// @Success 200 {object} bool
// @Router /hotels/error [get]
func (r *HotelRouter) GetErrorFlagHandler(c *gin.Context) *cerrors.CodedError {
	errorFlag := r.hotelService.GetErrorFlag(c.Request.Context())
	return r.Success(c, errorFlag)
}
