package decorator

import (
	"context"
	"github.com/sirupsen/logrus"
)

func ApplyRequestDecorators[H any, R any](handler RequestHandler[H, R], logger *logrus.Entry) RequestHandler[H, R] {
	return requestLoggingDecorator[H, R]{
		base:   handler,
		logger: logger,
	}
}

type RequestHandler[Req any, Res any] interface {
	Handle(ctx context.Context, r Req) (Res, error)
}
