package reconciler

import (
	"context"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Interface[T client.Object] interface {
	Reconcile(ctx context.Context, resource T) (ctrl.Result, error)
}
