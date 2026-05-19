package utils

import (
	"context"

	"flamespot-api/config"
)

func NewDBContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), config.DefaultTimeout)
}
