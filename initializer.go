package reconciler

import (
	"context"

	kerrors "k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

func Init(obj client.Object) Initializer {
	return &initializer{obj: obj}
}

type Initializer interface {
	WithManager(mgr manager.Manager) Builder
	WithReader(reader client.Reader) Builder
}

type Builder interface {
	Chain(handlers ...Handler) Chain
}

type Chain interface {
	Reconcile(ctx context.Context, key client.ObjectKey) (ctrl.Result, error)
}

type initializer struct {
	obj client.Object
}

func (b *initializer) WithManager(mgr manager.Manager) Builder {
	return &builder{obj: b.obj, reader: mgr.GetClient()}
}

func (b *initializer) WithReader(reader client.Reader) Builder {
	return &builder{obj: b.obj, reader: reader}
}

type builder struct {
	obj    client.Object
	reader client.Reader
}

func (c *builder) Chain(handlers ...Handler) Chain {
	handlers = append(handlers, &finisher{})
	handlers = append([]Handler{&tracer{}}, handlers...)
	var last Handler
	for i := len(handlers) - 1; i >= 0; i-- {
		current := handlers[i]
		current.setNext(last)
		last = current
	}
	return &chain{obj: c.obj, reader: c.reader, handlers: last}
}

type chain struct {
	Results
	obj      client.Object
	reader   client.Reader
	handlers Handler
}

func (c *chain) Reconcile(ctx context.Context, key client.ObjectKey) (ctrl.Result, error) {
	err := c.reader.Get(ctx, key, c.obj)
	if kerrors.IsNotFound(err) {
		return c.Finish(ctx) // Ignoring since object must be deleted
	}
	if err != nil {
		return c.RequeueOnErr(ctx, err)
	}
	return c.handlers.Reconcile(ctx, c.obj)
}
