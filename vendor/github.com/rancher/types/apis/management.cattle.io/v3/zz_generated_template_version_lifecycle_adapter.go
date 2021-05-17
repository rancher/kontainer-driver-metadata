package v3

import (
	"github.com/rancher/norman/lifecycle"
	"github.com/rancher/norman/resource"
	"k8s.io/apimachinery/pkg/runtime"
)

type TemplateVersionLifecycle interface {
	Create(obj *TemplateVersion) (runtime.Object, error)
	Remove(obj *TemplateVersion) (runtime.Object, error)
	Updated(obj *TemplateVersion) (runtime.Object, error)
}

type templateVersionLifecycleAdapter struct {
	lifecycle TemplateVersionLifecycle
}

func (w *templateVersionLifecycleAdapter) HasCreate() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasCreate()
}

func (w *templateVersionLifecycleAdapter) HasFinalize() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasFinalize()
}

func (w *templateVersionLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Create(obj.(*TemplateVersion))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *templateVersionLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Remove(obj.(*TemplateVersion))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *templateVersionLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Updated(obj.(*TemplateVersion))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewTemplateVersionLifecycleAdapter(name string, clusterScoped bool, client TemplateVersionInterface, l TemplateVersionLifecycle) TemplateVersionHandlerFunc {
	if clusterScoped {
		resource.PutClusterScoped(TemplateVersionGroupVersionResource)
	}
	adapter := &templateVersionLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *TemplateVersion) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}
