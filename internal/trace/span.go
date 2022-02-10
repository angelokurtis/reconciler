package trace

import (
	"context"
	"errors"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
)

type Span struct{ trace.Span }

func SpanFromContext(ctx context.Context) Span {
	return Span{Span: trace.SpanFromContext(ctx)}
}

func (s *Span) Error(err error) error {
	if err != nil {
		s.RecordError(err)
		s.SetStatus(codes.Error, err.Error())
		var serr *kerrors.StatusError
		if errors.As(err, &serr) {
			status := serr.Status()
			s.SetAttributes(attribute.Int64("code", int64(status.Code)))
			s.SetAttributes(attribute.String("reason", string(status.Reason)))
		}
	}
	return err
}
