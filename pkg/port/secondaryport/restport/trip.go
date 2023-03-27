package restport

import (
	"context"
)

type TripRestAdapter interface {
	GetInjectionStatus(ctx context.Context) (bool, error)
	InjectError(ctx context.Context, flag bool) error
}
