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
	"fmt"
	"time"

	"github.com/go-logr/logr"
	ctrl "sigs.k8s.io/controller-runtime"
)

const resultCallDepth = 2

type Result struct{}

func (r *Result) Finish(ctx context.Context) (ctrl.Result, error) {
	log := logr.FromContextOrDiscard(ctx).WithCallDepth(resultCallDepth)
	log.Info("Reconciliation finished successfully")

	return ctrl.Result{}, nil
}

func (r *Result) Requeue(ctx context.Context) (ctrl.Result, error) {
	log := logr.FromContextOrDiscard(ctx).WithCallDepth(resultCallDepth)
	log.Info("Reconciliation requeued", "requeue", "now")

	return ctrl.Result{Requeue: true}, nil
}

func (r *Result) RequeueAfter(ctx context.Context, duration time.Duration) (ctrl.Result, error) {
	log := logr.FromContextOrDiscard(ctx).WithCallDepth(resultCallDepth)
	log.Info("Reconciliation requeued", "requeue", fmt.Sprintf("in %s", duration))

	return ctrl.Result{RequeueAfter: duration}, nil
}

func (r *Result) RequeueOnErr(ctx context.Context, err error) (ctrl.Result, error) {
	log := logr.FromContextOrDiscard(ctx).WithCallDepth(resultCallDepth)
	log.Info("Reconciliation requeued due to error", "requeue", "now")

	return ctrl.Result{}, err
}

func (r *Result) RequeueOnErrAfter(ctx context.Context, err error, duration time.Duration) (ctrl.Result, error) {
	log := logr.FromContextOrDiscard(ctx).WithCallDepth(resultCallDepth)
	log.Info("Reconciliation requeued due to error", "requeue", fmt.Sprintf("in %s", duration))

	return ctrl.Result{RequeueAfter: duration}, err
}
