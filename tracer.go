/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package reconcilers

import (
	"context"
	"fmt"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type tracer struct{ Funcs }

func (t *tracer) Reconcile(ctx context.Context, obj client.Object) (ctrl.Result, error) {
	l := log.FromContext(ctx)
	l.Info("Reconciler has been triggered")
	result, err := t.next.Reconcile(ctx, obj)
	switch {
	case err != nil:
		l.Error(err, "Reconciler error")
		return result, err
	case result.RequeueAfter > 0:
		l.Info("Successfully reconciled!", "requeue", fmt.Sprintf("in %s", result.RequeueAfter))
		return result, nil
	case result.Requeue:
		l.Info("Successfully reconciled!", "requeue", "now")
		return result, nil
	}
	l.Info("Successfully reconciled!")
	return result, nil
}
