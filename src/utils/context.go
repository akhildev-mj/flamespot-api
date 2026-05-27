package utils

import (
	"context"
	"flamespot-api/src/config"
)

func NewDBContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), config.DefaultTimeout)
}
