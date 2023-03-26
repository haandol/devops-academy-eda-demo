package router

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/haandol/devops-academy-eda-demo/pkg/constant"
	"github.com/haandol/devops-academy-eda-demo/pkg/port/primaryport/routerport"
	"github.com/haandol/devops-academy-eda-demo/pkg/service"
	"github.com/haandol/devops-academy-eda-demo/pkg/util/cerrors"
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
	g.Handle("POST", "/", r.WrappedHandler(r.ToggleErrorInjection))
	g.Handle("GET", "/", r.WrappedHandler(r.GetErrorFlag))
}

// @Summary toggle error injection
// @Schemes
// @Description toggle error injection on creating hotel
// @Tags hotels
// @Accept json
// @Produce json
// @Success 200 {object} bool
// @Router /hotels/error [post]
func (r *HotelRouter) ToggleErrorInjection(c *gin.Context) *cerrors.CodedError {
	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Second*10)
	defer cancel()

	errorFlag, err := r.hotelService.ToggleErrorInjection(ctx)
	if err != nil {
		return cerrors.New(constant.ErrFailToCreateTrip, err)
	}

	return r.Success(c, errorFlag)
}

// @Summary check error flag
// @Schemes
// @Description return error flag
// @Tags trips
// @Accept json
// @Produce json
// @Success 200 {object} bool
// @Router /hotels/error [get]
func (r *HotelRouter) GetErrorFlag(c *gin.Context) *cerrors.CodedError {
	errorFlag := r.hotelService.GetErrorFlag(c.Request.Context())

	return r.Success(c, errorFlag)
}
