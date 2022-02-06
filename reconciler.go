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

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type (
	Reconciler interface {
		Reconcile(ctx context.Context, obj client.Object) (ctrl.Result, error)
		SetNext(next Reconciler)
		Next(ctx context.Context, obj client.Object) (ctrl.Result, error)
	}
)

func Chain(reconcilers ...Reconciler) Reconciler {
	reconcilers = append(reconcilers, &finisher{})
	reconcilers = append([]Reconciler{&tracer{}}, reconcilers...)
	var last Reconciler
	for i := len(reconcilers) - 1; i >= 0; i-- {
		current := reconcilers[i]
		current.SetNext(last)
		last = current
	}
	return last
}
