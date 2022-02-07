# reconciler

A library to avoid overstuffed Reconcile functions of Kubernetes operators as suggested in
this [blog post](https://cloud.redhat.com/blog/7-best-practices-for-writing-kubernetes-operators-an-sre-perspective)

## Features

* Avoid overstuffed functions through handlers in a chain of responsibility
* Wrap the reconciliation results to improve the readability of your code

## Getting Started

Install by running:

```shell
go install github.com/angelokurtis/reconciler@latest
```

Split the Reconcile responsibilities into handlers with one single purpose and chain them in your controller:

```go
import "github.com/angelokurtis/reconciler"

// code omitted

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