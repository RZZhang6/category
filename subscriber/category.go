package subscriber

import (
	"context"
	log "github.com/micro/go-micro/v2/logger"

	category "category/proto/category"
)

type Category struct{}

func (e *Category) Handle(ctx context.Context, msg *category.Message) error {
	log.Info("Handler Received message: ", msg.Say)
	return nil
}

func Handler(ctx context.Context, msg *category.Message) error {
	log.Info("Function Received message: ", msg.Say)
	return nil
}
