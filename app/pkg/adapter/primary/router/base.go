package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/haandol/devops-academy-eda-demo/pkg/util/cerrors"
)

type BaseRouter struct{}

func (r BaseRouter) WrappedHandler(f func(c *gin.Context) *cerrors.CodedError) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := f(c); err != nil {
			c.AbortWithStatusJSON(
				http.StatusInternalServerError,
				gin.H{"status": false, "code": err.Code, "message": err.Error()},
			)
		}
	}
}

func (r BaseRouter) Success(c *gin.Context, data any) *cerrors.CodedError {
	c.JSON(http.StatusOK, gin.H{"status": true, "data": data})
	return nil
}
