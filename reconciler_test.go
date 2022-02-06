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

package reconcilers_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logr "sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/angelokurtis/reconcilers"
)

func TestOperations(t *testing.T) {
	t.Run("should throw error when an operation fails", func(t *testing.T) {
		r := reconcilers.Chain(
			&withoutError{},
			&withError{},
			&withRequeueIn2min{},
		)
		_, err := r.Reconcile(context.TODO(), nil)
		if err == nil {
			t.Errorf("expect error, got %+v", err)
		}
	})
	t.Run("should requeue in 2m when an operation have this result", func(t *testing.T) {
		r := reconcilers.Chain(
			&withRequeueIn2min{},
			&withoutError{},
		)
		result, err := r.Reconcile(context.TODO(), nil)
		if err != nil || result.RequeueAfter != 2*time.Minute {
			t.Errorf("expect requeue in 2min, got result=%+v err=%+v", result, err)
		}
	})
	t.Run("should requeue now when an operation have this result", func(t *testing.T) {
		r := reconcilers.Chain(
			&withRequeue{},
			&withoutError{},
		)
		result, err := r.Reconcile(context.TODO(), nil)
		if err != nil || !result.Requeue {
			t.Errorf("expect requeue true, got result=%+v err=%+v", result, err)
		}
	})
	t.Run("should finish all operations when none requeue or fail", func(t *testing.T) {
		r := reconcilers.Chain(
			&withoutError{},
			&withoutError{},
		)
		result, err := r.Reconcile(context.TODO(), nil)
		if err != nil || result.Requeue || result.RequeueAfter > 0 {
			t.Errorf("expect empty result, got result=%+v err=%+v", result, err)
		}
	})
}

type withoutError struct{ reconcilers.Funcs }

func (w *withoutError) Reconcile(ctx context.Context, obj client.Object) (ctrl.Result, error) {
	logr.FromContext(ctx).V(0).Info(fmt.Sprintf("%T", w))
	return w.Next(ctx, obj)
}

type withError struct{ reconcilers.Funcs }

func (w *withError) Reconcile(ctx context.Context, obj client.Object) (ctrl.Result, error) {
	return w.RequeueOnErr(ctx, errors.New("reconcile with error"))
}

type withRequeue struct{ reconcilers.Funcs }

func (w *withRequeue) Reconcile(ctx context.Context, obj client.Object) (ctrl.Result, error) {
	return w.Requeue(ctx)
}

type withRequeueIn2min struct{ reconcilers.Funcs }

func (w *withRequeueIn2min) Reconcile(ctx context.Context, obj client.Object) (ctrl.Result, error) {
	return w.RequeueAfter(ctx, 2*time.Minute)
}
