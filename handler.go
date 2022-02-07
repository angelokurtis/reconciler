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

package reconciler

import (
	"context"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Handler interface {
	Reconcile(ctx context.Context, obj client.Object) (ctrl.Result, error)
	Next(ctx context.Context, obj client.Object) (ctrl.Result, error)
	setNext(next Handler)
}

func Chain(handlers ...Handler) Handler {
	handlers = append(handlers, &finisher{})
	handlers = append([]Handler{&tracer{}}, handlers...)
	var last Handler
	for i := len(handlers) - 1; i >= 0; i-- {
		current := handlers[i]
		current.setNext(last)
		last = current
	}
	return last
}
