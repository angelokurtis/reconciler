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
	"time"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Funcs struct{ next Handler }

func (f *Funcs) setNext(next Handler) { f.next = next }

func (f *Funcs) Next(ctx context.Context, obj client.Object) (ctrl.Result, error) {
	return f.next.Reconcile(ctx, obj)
}

func (f *Funcs) Finish(ctx context.Context) (ctrl.Result, error) {
	return ctrl.Result{}, nil
}

func (f *Funcs) Requeue(ctx context.Context) (ctrl.Result, error) {
	return ctrl.Result{Requeue: true}, nil
}

func (f *Funcs) RequeueAfter(ctx context.Context, duration time.Duration) (ctrl.Result, error) {
	return ctrl.Result{RequeueAfter: duration}, nil
}

func (f *Funcs) RequeueOnErr(ctx context.Context, err error) (ctrl.Result, error) {
	return ctrl.Result{}, err
}

func (f *Funcs) RequeueOnErrAfter(ctx context.Context, err error, duration time.Duration) (ctrl.Result, error) {
	return ctrl.Result{RequeueAfter: duration}, err
}
