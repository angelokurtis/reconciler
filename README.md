# reconciler

[![Go Reference](https://pkg.go.dev/badge/github.com/angelokurtis/reconciler.svg)](https://pkg.go.dev/github.com/angelokurtis/reconciler)
[![Build Status](https://github.com/angelokurtis/reconciler/actions/workflows/test.yml/badge.svg?branch=main)](https://github.com/angelokurtis/reconciler/actions)
[![Coverage](https://codecov.io/gh/angelokurtis/reconciler/branch/main/graph/badge.svg)](https://codecov.io/gh/angelokurtis/reconciler)

A library to avoid overstuffed Reconcile functions of Kubernetes operators as suggested in
this [blog post](https://cloud.redhat.com/blog/7-best-practices-for-writing-kubernetes-operators-an-sre-perspective)

## Features

* Avoid overstuffed functions through handlers in a chain of responsibility
* Wrap the reconciliation results to improve the readability of your code

## Getting Started

Install by running:

```shell
go get github.com/angelokurtis/reconciler
```

Split the Reconcile responsibilities into handlers with one single purpose and chain them in your controller:

```go
// package omitted

import "github.com/angelokurtis/reconciler"
// other imports

func (r *MemcachedReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	// Fetch the Memcached instance
	m := &cachev1alpha1.Memcached{}
	err := r.Get(ctx, req.NamespacedName, m)
	if err != nil {
		if errors.IsNotFound(err) {
			return r.Finish(ctx) // Ignoring since object must be deleted
		}
		return r.RequeueOnErr(ctx, err) // Failed to get Memcached
	}

	// Chains reconcile handlers
	return reconciler.Chain(
		&memcached.DeploymentCreation{Client: r.Client, Scheme: r.Scheme},
		&memcached.Status{Client: r.Client},
	).Reconcile(ctx, m)
}
```

Your handler can compose [reconciler.Funcs](https://github.com/angelokurtis/reconciler/blob/main/funcs.go#L27) that
already implement most of the [interface](https://github.com/angelokurtis/reconciler/blob/main/handler.go#L26) functions
so you just need to put your logic
on [Reconcile(context.Context, client.Object)](https://github.com/angelokurtis/reconciler/blob/main/handler.go#L27)
implementation:

```go
// package omitted

import "github.com/angelokurtis/reconciler"
// other imports

type DeploymentCreation struct {
	reconciler.Funcs
	// other fields
}

func (d *DeploymentCreation) Reconcile(ctx context.Context, obj client.Object) (ctrl.Result, error) {
	// TODO(user): your logic here
	return d.Next(ctx, obj)
}
```
