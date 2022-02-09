package trace

import (
	"context"

	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

func LogFromContext(ctx context.Context) logr.Logger {
	span := SpanFromContext(ctx)
	logger := log.FromContext(ctx)
	return span.Log(logger)
}
