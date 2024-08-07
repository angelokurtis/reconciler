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
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type Handler[T client.Object] interface {
	reconcile.ObjectReconciler[T]
	Next(ctx context.Context, resource T) (ctrl.Result, error)
	setNext(next Handler[T])
}

func Chain[T client.Object](handlers ...Handler[T]) Handler[T] {
	handlers = append(handlers, &finisher[T]{})
	handlers = append([]Handler[T]{&initializer[T]{}}, handlers...)

	var last Handler[T]

	for i := len(handlers) - 1; i >= 0; i-- {
		current := handlers[i]
		current.setNext(last)
		last = current
	}

	return last
}
