package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/haandol/devops-academy-eda-demo/pkg/util"
	"go.uber.org/ratelimit"
)

func LeakBucket(rps int) gin.HandlerFunc {
	logger := util.GetLogger().With(
		"module", "Middleware",
		"func", "LeakBucket",
	)

	var limiter ratelimit.Limiter
	if rps == 0 {
		limiter = ratelimit.NewUnlimited()
	} else {
		limiter = ratelimit.New(rps)
	}

	prev := time.Now()
	return func(c *gin.Context) {
		now := limiter.Take()
		logger.Debugf("%v", now.Sub(prev))
		prev = now
		c.Next()
	}
}
