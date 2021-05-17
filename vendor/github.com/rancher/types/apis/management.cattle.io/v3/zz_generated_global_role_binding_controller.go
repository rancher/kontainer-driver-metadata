package v3

import (
	"context"

	"github.com/rancher/norman/controller"
	"github.com/rancher/norman/objectclient"
	"github.com/rancher/norman/resource"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
)

var (
	GlobalRoleBindingGroupVersionKind = schema.GroupVersionKind{
		Version: Version,
		Group:   GroupName,
		Kind:    "GlobalRoleBinding",
	}
	GlobalRoleBindingResource = metav1.APIResource{
		Name:         "globalrolebindings",
		SingularName: "globalrolebinding",
		Namespaced:   false,
		Kind:         GlobalRoleBindingGroupVersionKind.Kind,
	}

	GlobalRoleBindingGroupVersionResource = schema.GroupVersionResource{
		Group:    GroupName,
		Version:  Version,
		Resource: "globalrolebindings",
	}
)

func init() {
	resource.Put(GlobalRoleBindingGroupVersionResource)
}

func NewGlobalRoleBinding(namespace, name string, obj GlobalRoleBinding) *GlobalRoleBinding {
	obj.APIVersion, obj.Kind = GlobalRoleBindingGroupVersionKind.ToAPIVersionAndKind()
	obj.Name = name
	obj.Namespace = namespace
	return &obj
}

type GlobalRoleBindingList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []GlobalRoleBinding `json:"items"`
}

type GlobalRoleBindingHandlerFunc func(key string, obj *GlobalRoleBinding) (runtime.Object, error)

type GlobalRoleBindingChangeHandlerFunc func(obj *GlobalRoleBinding) (runtime.Object, error)

type GlobalRoleBindingLister interface {
	List(namespace string, selector labels.Selector) (ret []*GlobalRoleBinding, err error)
	Get(namespace, name string) (*GlobalRoleBinding, error)
}

type GlobalRoleBindingController interface {
	Generic() controller.GenericController
	Informer() cache.SharedIndexInformer
	Lister() GlobalRoleBindingLister
	AddHandler(ctx context.Context, name string, handler GlobalRoleBindingHandlerFunc)
	AddFeatureHandler(ctx context.Context, enabled func() bool, name string, sync GlobalRoleBindingHandlerFunc)
	AddClusterScopedHandler(ctx context.Context, name, clusterName string, handler GlobalRoleBindingHandlerFunc)
	AddClusterScopedFeatureHandler(ctx context.Context, enabled func() bool, name, clusterName string, handler GlobalRoleBindingHandlerFunc)
	Enqueue(namespace, name string)
	Sync(ctx context.Context) error
	Start(ctx context.Context, threadiness int) error
}

type GlobalRoleBindingInterface interface {
	ObjectClient() *objectclient.ObjectClient
	Create(*GlobalRoleBinding) (*GlobalRoleBinding, error)
	GetNamespaced(namespace, name string, opts metav1.GetOptions) (*GlobalRoleBinding, error)
	Get(name string, opts metav1.GetOptions) (*GlobalRoleBinding, error)
	Update(*GlobalRoleBinding) (*GlobalRoleBinding, error)
	Delete(name string, options *metav1.DeleteOptions) error
	DeleteNamespaced(namespace, name string, options *metav1.DeleteOptions) error
	List(opts metav1.ListOptions) (*GlobalRoleBindingList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	DeleteCollection(deleteOpts *metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Controller() GlobalRoleBindingController
	AddHandler(ctx context.Context, name string, sync GlobalRoleBindingHandlerFunc)
	AddFeatureHandler(ctx context.Context, enabled func() bool, name string, sync GlobalRoleBindingHandlerFunc)
	AddLifecycle(ctx context.Context, name string, lifecycle GlobalRoleBindingLifecycle)
	AddFeatureLifecycle(ctx context.Context, enabled func() bool, name string, lifecycle GlobalRoleBindingLifecycle)
	AddClusterScopedHandler(ctx context.Context, name, clusterName string, sync GlobalRoleBindingHandlerFunc)
	AddClusterScopedFeatureHandler(ctx context.Context, enabled func() bool, name, clusterName string, sync GlobalRoleBindingHandlerFunc)
	AddClusterScopedLifecycle(ctx context.Context, name, clusterName string, lifecycle GlobalRoleBindingLifecycle)
	AddClusterScopedFeatureLifecycle(ctx context.Context, enabled func() bool, name, clusterName string, lifecycle GlobalRoleBindingLifecycle)
}

type globalRoleBindingLister struct {
	controller *globalRoleBindingController
}

func (l *globalRoleBindingLister) List(namespace string, selector labels.Selector) (ret []*GlobalRoleBinding, err error) {
	err = cache.ListAllByNamespace(l.controller.Informer().GetIndexer(), namespace, selector, func(obj interface{}) {
		ret = append(ret, obj.(*GlobalRoleBinding))
	})
	return
}

func (l *globalRoleBindingLister) Get(namespace, name string) (*GlobalRoleBinding, error) {
	var key string
	if namespace != "" {
		key = namespace + "/" + name
	} else {
		key = name
	}
	obj, exists, err := l.controller.Informer().GetIndexer().GetByKey(key)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(schema.GroupResource{
			Group:    GlobalRoleBindingGroupVersionKind.Group,
			Resource: "globalRoleBinding",
		}, key)
	}
	return obj.(*GlobalRoleBinding), nil
}

type globalRoleBindingController struct {
	controller.GenericController
}

func (c *globalRoleBindingController) Generic() controller.GenericController {
	return c.GenericController
}

func (c *globalRoleBindingController) Lister() GlobalRoleBindingLister {
	return &globalRoleBindingLister{
		controller: c,
	}
}

func (c *globalRoleBindingController) AddHandler(ctx context.Context, name string, handler GlobalRoleBindingHandlerFunc) {
	c.GenericController.AddHandler(ctx, name, func(key string, obj interface{}) (interface{}, error) {
		if obj == nil {
			return handler(key, nil)
		} else if v, ok := obj.(*GlobalRoleBinding); ok {
			return handler(key, v)
		} else {
			return nil, nil
		}
	})
}

func (c *globalRoleBindingController) AddFeatureHandler(ctx context.Context, enabled func() bool, name string, handler GlobalRoleBindingHandlerFunc) {
	c.GenericController.AddHandler(ctx, name, func(key string, obj interface{}) (interface{}, error) {
		if !enabled() {
			return nil, nil
		} else if obj == nil {
			return handler(key, nil)
		} else if v, ok := obj.(*GlobalRoleBinding); ok {
			return handler(key, v)
		} else {
			return nil, nil
		}
	})
}

func (c *globalRoleBindingController) AddClusterScopedHandler(ctx context.Context, name, cluster string, handler GlobalRoleBindingHandlerFunc) {
	c.GenericController.AddHandler(ctx, name, func(key string, obj interface{}) (interface{}, error) {
		if obj == nil {
			return handler(key, nil)
		} else if v, ok := obj.(*GlobalRoleBinding); ok && controller.ObjectInCluster(cluster, obj) {
			return handler(key, v)
		} else {
			return nil, nil
		}
	})
}

func (c *globalRoleBindingController) AddClusterScopedFeatureHandler(ctx context.Context, enabled func() bool, name, cluster string, handler GlobalRoleBindingHandlerFunc) {
	c.GenericController.AddHandler(ctx, name, func(key string, obj interface{}) (interface{}, error) {
		if !enabled() {
			return nil, nil
		} else if obj == nil {
			return handler(key, nil)
		} else if v, ok := obj.(*GlobalRoleBinding); ok && controller.ObjectInCluster(cluster, obj) {
			return handler(key, v)
		} else {
			return nil, nil
		}
	})
}

type globalRoleBindingFactory struct {
}

func (c globalRoleBindingFactory) Object() runtime.Object {
	return &GlobalRoleBinding{}
}

func (c globalRoleBindingFactory) List() runtime.Object {
	return &GlobalRoleBindingList{}
}

func (s *globalRoleBindingClient) Controller() GlobalRoleBindingController {
	s.client.Lock()
	defer s.client.Unlock()

	c, ok := s.client.globalRoleBindingControllers[s.ns]
	if ok {
		return c
	}

	genericController := controller.NewGenericController(GlobalRoleBindingGroupVersionKind.Kind+"Controller",
		s.objectClient)

	c = &globalRoleBindingController{
		GenericController: genericController,
	}

	s.client.globalRoleBindingControllers[s.ns] = c
	s.client.starters = append(s.client.starters, c)

	return c
}

type globalRoleBindingClient struct {
	client       *Client
	ns           string
	objectClient *objectclient.ObjectClient
	controller   GlobalRoleBindingController
}

func (s *globalRoleBindingClient) ObjectClient() *objectclient.ObjectClient {
	return s.objectClient
}

func (s *globalRoleBindingClient) Create(o *GlobalRoleBinding) (*GlobalRoleBinding, error) {
	obj, err := s.objectClient.Create(o)
	return obj.(*GlobalRoleBinding), err
}

func (s *globalRoleBindingClient) Get(name string, opts metav1.GetOptions) (*GlobalRoleBinding, error) {
	obj, err := s.objectClient.Get(name, opts)
	return obj.(*GlobalRoleBinding), err
}

func (s *globalRoleBindingClient) GetNamespaced(namespace, name string, opts metav1.GetOptions) (*GlobalRoleBinding, error) {
	obj, err := s.objectClient.GetNamespaced(namespace, name, opts)
	return obj.(*GlobalRoleBinding), err
}

func (s *globalRoleBindingClient) Update(o *GlobalRoleBinding) (*GlobalRoleBinding, error) {
	obj, err := s.objectClient.Update(o.Name, o)
	return obj.(*GlobalRoleBinding), err
}

func (s *globalRoleBindingClient) Delete(name string, options *metav1.DeleteOptions) error {
	return s.objectClient.Delete(name, options)
}

func (s *globalRoleBindingClient) DeleteNamespaced(namespace, name string, options *metav1.DeleteOptions) error {
	return s.objectClient.DeleteNamespaced(namespace, name, options)
}

func (s *globalRoleBindingClient) List(opts metav1.ListOptions) (*GlobalRoleBindingList, error) {
	obj, err := s.objectClient.List(opts)
	return obj.(*GlobalRoleBindingList), err
}

func (s *globalRoleBindingClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	return s.objectClient.Watch(opts)
}

// Patch applies the patch and returns the patched deployment.
func (s *globalRoleBindingClient) Patch(o *GlobalRoleBinding, patchType types.PatchType, data []byte, subresources ...string) (*GlobalRoleBinding, error) {
	obj, err := s.objectClient.Patch(o.Name, o, patchType, data, subresources...)
	return obj.(*GlobalRoleBinding), err
}

func (s *globalRoleBindingClient) DeleteCollection(deleteOpts *metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	return s.objectClient.DeleteCollection(deleteOpts, listOpts)
}

func (s *globalRoleBindingClient) AddHandler(ctx context.Context, name string, sync GlobalRoleBindingHandlerFunc) {
	s.Controller().AddHandler(ctx, name, sync)
}

func (s *globalRoleBindingClient) AddFeatureHandler(ctx context.Context, enabled func() bool, name string, sync GlobalRoleBindingHandlerFunc) {
	s.Controller().AddFeatureHandler(ctx, enabled, name, sync)
}

func (s *globalRoleBindingClient) AddLifecycle(ctx context.Context, name string, lifecycle GlobalRoleBindingLifecycle) {
	sync := NewGlobalRoleBindingLifecycleAdapter(name, false, s, lifecycle)
	s.Controller().AddHandler(ctx, name, sync)
}

func (s *globalRoleBindingClient) AddFeatureLifecycle(ctx context.Context, enabled func() bool, name string, lifecycle GlobalRoleBindingLifecycle) {
	sync := NewGlobalRoleBindingLifecycleAdapter(name, false, s, lifecycle)
	s.Controller().AddFeatureHandler(ctx, enabled, name, sync)
}

func (s *globalRoleBindingClient) AddClusterScopedHandler(ctx context.Context, name, clusterName string, sync GlobalRoleBindingHandlerFunc) {
	s.Controller().AddClusterScopedHandler(ctx, name, clusterName, sync)
}

func (s *globalRoleBindingClient) AddClusterScopedFeatureHandler(ctx context.Context, enabled func() bool, name, clusterName string, sync GlobalRoleBindingHandlerFunc) {
	s.Controller().AddClusterScopedFeatureHandler(ctx, enabled, name, clusterName, sync)
}

func (s *globalRoleBindingClient) AddClusterScopedLifecycle(ctx context.Context, name, clusterName string, lifecycle GlobalRoleBindingLifecycle) {
	sync := NewGlobalRoleBindingLifecycleAdapter(name+"_"+clusterName, true, s, lifecycle)
	s.Controller().AddClusterScopedHandler(ctx, name, clusterName, sync)
}

func (s *globalRoleBindingClient) AddClusterScopedFeatureLifecycle(ctx context.Context, enabled func() bool, name, clusterName string, lifecycle GlobalRoleBindingLifecycle) {
	sync := NewGlobalRoleBindingLifecycleAdapter(name+"_"+clusterName, true, s, lifecycle)
	s.Controller().AddClusterScopedFeatureHandler(ctx, enabled, name, clusterName, sync)
}

type GlobalRoleBindingIndexer func(obj *GlobalRoleBinding) ([]string, error)

type GlobalRoleBindingClientCache interface {
	Get(namespace, name string) (*GlobalRoleBinding, error)
	List(namespace string, selector labels.Selector) ([]*GlobalRoleBinding, error)

	Index(name string, indexer GlobalRoleBindingIndexer)
	GetIndexed(name, key string) ([]*GlobalRoleBinding, error)
}

type GlobalRoleBindingClient interface {
	Create(*GlobalRoleBinding) (*GlobalRoleBinding, error)
	Get(namespace, name string, opts metav1.GetOptions) (*GlobalRoleBinding, error)
	Update(*GlobalRoleBinding) (*GlobalRoleBinding, error)
	Delete(namespace, name string, options *metav1.DeleteOptions) error
	List(namespace string, opts metav1.ListOptions) (*GlobalRoleBindingList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)

	Cache() GlobalRoleBindingClientCache

	OnCreate(ctx context.Context, name string, sync GlobalRoleBindingChangeHandlerFunc)
	OnChange(ctx context.Context, name string, sync GlobalRoleBindingChangeHandlerFunc)
	OnRemove(ctx context.Context, name string, sync GlobalRoleBindingChangeHandlerFunc)
	Enqueue(namespace, name string)

	Generic() controller.GenericController
	ObjectClient() *objectclient.ObjectClient
	Interface() GlobalRoleBindingInterface
}

type globalRoleBindingClientCache struct {
	client *globalRoleBindingClient2
}

type globalRoleBindingClient2 struct {
	iface      GlobalRoleBindingInterface
	controller GlobalRoleBindingController
}

func (n *globalRoleBindingClient2) Interface() GlobalRoleBindingInterface {
	return n.iface
}

func (n *globalRoleBindingClient2) Generic() controller.GenericController {
	return n.iface.Controller().Generic()
}

func (n *globalRoleBindingClient2) ObjectClient() *objectclient.ObjectClient {
	return n.Interface().ObjectClient()
}

func (n *globalRoleBindingClient2) Enqueue(namespace, name string) {
	n.iface.Controller().Enqueue(namespace, name)
}

func (n *globalRoleBindingClient2) Create(obj *GlobalRoleBinding) (*GlobalRoleBinding, error) {
	return n.iface.Create(obj)
}

func (n *globalRoleBindingClient2) Get(namespace, name string, opts metav1.GetOptions) (*GlobalRoleBinding, error) {
	return n.iface.GetNamespaced(namespace, name, opts)
}

func (n *globalRoleBindingClient2) Update(obj *GlobalRoleBinding) (*GlobalRoleBinding, error) {
	return n.iface.Update(obj)
}

func (n *globalRoleBindingClient2) Delete(namespace, name string, options *metav1.DeleteOptions) error {
	return n.iface.DeleteNamespaced(namespace, name, options)
}

func (n *globalRoleBindingClient2) List(namespace string, opts metav1.ListOptions) (*GlobalRoleBindingList, error) {
	return n.iface.List(opts)
}

func (n *globalRoleBindingClient2) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	return n.iface.Watch(opts)
}

func (n *globalRoleBindingClientCache) Get(namespace, name string) (*GlobalRoleBinding, error) {
	return n.client.controller.Lister().Get(namespace, name)
}

func (n *globalRoleBindingClientCache) List(namespace string, selector labels.Selector) ([]*GlobalRoleBinding, error) {
	return n.client.controller.Lister().List(namespace, selector)
}

func (n *globalRoleBindingClient2) Cache() GlobalRoleBindingClientCache {
	n.loadController()
	return &globalRoleBindingClientCache{
		client: n,
	}
}

func (n *globalRoleBindingClient2) OnCreate(ctx context.Context, name string, sync GlobalRoleBindingChangeHandlerFunc) {
	n.loadController()
	n.iface.AddLifecycle(ctx, name+"-create", &globalRoleBindingLifecycleDelegate{create: sync})
}

func (n *globalRoleBindingClient2) OnChange(ctx context.Context, name string, sync GlobalRoleBindingChangeHandlerFunc) {
	n.loadController()
	n.iface.AddLifecycle(ctx, name+"-change", &globalRoleBindingLifecycleDelegate{update: sync})
}

func (n *globalRoleBindingClient2) OnRemove(ctx context.Context, name string, sync GlobalRoleBindingChangeHandlerFunc) {
	n.loadController()
	n.iface.AddLifecycle(ctx, name, &globalRoleBindingLifecycleDelegate{remove: sync})
}

func (n *globalRoleBindingClientCache) Index(name string, indexer GlobalRoleBindingIndexer) {
	err := n.client.controller.Informer().GetIndexer().AddIndexers(map[string]cache.IndexFunc{
		name: func(obj interface{}) ([]string, error) {
			if v, ok := obj.(*GlobalRoleBinding); ok {
				return indexer(v)
			}
			return nil, nil
		},
	})

	if err != nil {
		panic(err)
	}
}

func (n *globalRoleBindingClientCache) GetIndexed(name, key string) ([]*GlobalRoleBinding, error) {
	var result []*GlobalRoleBinding
	objs, err := n.client.controller.Informer().GetIndexer().ByIndex(name, key)
	if err != nil {
		return nil, err
	}
	for _, obj := range objs {
		if v, ok := obj.(*GlobalRoleBinding); ok {
			result = append(result, v)
		}
	}

	return result, nil
}

func (n *globalRoleBindingClient2) loadController() {
	if n.controller == nil {
		n.controller = n.iface.Controller()
	}
}

type globalRoleBindingLifecycleDelegate struct {
	create GlobalRoleBindingChangeHandlerFunc
	update GlobalRoleBindingChangeHandlerFunc
	remove GlobalRoleBindingChangeHandlerFunc
}

func (n *globalRoleBindingLifecycleDelegate) HasCreate() bool {
	return n.create != nil
}

func (n *globalRoleBindingLifecycleDelegate) Create(obj *GlobalRoleBinding) (runtime.Object, error) {
	if n.create == nil {
		return obj, nil
	}
	return n.create(obj)
}

func (n *globalRoleBindingLifecycleDelegate) HasFinalize() bool {
	return n.remove != nil
}

func (n *globalRoleBindingLifecycleDelegate) Remove(obj *GlobalRoleBinding) (runtime.Object, error) {
	if n.remove == nil {
		return obj, nil
	}
	return n.remove(obj)
}

func (n *globalRoleBindingLifecycleDelegate) Updated(obj *GlobalRoleBinding) (runtime.Object, error) {
	if n.update == nil {
		return obj, nil
	}
	return n.update(obj)
}
