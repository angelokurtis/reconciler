package reconciler

import (
	"context"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Interface interface {
	Reconcile(ctx context.Context, obj client.Object) (ctrl.Result, error)
}
