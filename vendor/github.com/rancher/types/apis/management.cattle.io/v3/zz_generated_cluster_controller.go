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
	ClusterGroupVersionKind = schema.GroupVersionKind{
		Version: Version,
		Group:   GroupName,
		Kind:    "Cluster",
	}
	ClusterResource = metav1.APIResource{
		Name:         "clusters",
		SingularName: "cluster",
		Namespaced:   false,
		Kind:         ClusterGroupVersionKind.Kind,
	}

	ClusterGroupVersionResource = schema.GroupVersionResource{
		Group:    GroupName,
		Version:  Version,
		Resource: "clusters",
	}
)

func init() {
	resource.Put(ClusterGroupVersionResource)
}

func NewCluster(namespace, name string, obj Cluster) *Cluster {
	obj.APIVersion, obj.Kind = ClusterGroupVersionKind.ToAPIVersionAndKind()
	obj.Name = name
	obj.Namespace = namespace
	return &obj
}

type ClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Cluster `json:"items"`
}

type ClusterHandlerFunc func(key string, obj *Cluster) (runtime.Object, error)

type ClusterChangeHandlerFunc func(obj *Cluster) (runtime.Object, error)

type ClusterLister interface {
	List(namespace string, selector labels.Selector) (ret []*Cluster, err error)
	Get(namespace, name string) (*Cluster, error)
}

type ClusterController interface {
	Generic() controller.GenericController
	Informer() cache.SharedIndexInformer
	Lister() ClusterLister
	AddHandler(ctx context.Context, name string, handler ClusterHandlerFunc)
	AddFeatureHandler(ctx context.Context, enabled func() bool, name string, sync ClusterHandlerFunc)
	AddClusterScopedHandler(ctx context.Context, name, clusterName string, handler ClusterHandlerFunc)
	AddClusterScopedFeatureHandler(ctx context.Context, enabled func() bool, name, clusterName string, handler ClusterHandlerFunc)
	Enqueue(namespace, name string)
	Sync(ctx context.Context) error
	Start(ctx context.Context, threadiness int) error
}

type ClusterInterface interface {
	ObjectClient() *objectclient.ObjectClient
	Create(*Cluster) (*Cluster, error)
	GetNamespaced(namespace, name string, opts metav1.GetOptions) (*Cluster, error)
	Get(name string, opts metav1.GetOptions) (*Cluster, error)
	Update(*Cluster) (*Cluster, error)
	Delete(name string, options *metav1.DeleteOptions) error
	DeleteNamespaced(namespace, name string, options *metav1.DeleteOptions) error
	List(opts metav1.ListOptions) (*ClusterList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	DeleteCollection(deleteOpts *metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Controller() ClusterController
	AddHandler(ctx context.Context, name string, sync ClusterHandlerFunc)
	AddFeatureHandler(ctx context.Context, enabled func() bool, name string, sync ClusterHandlerFunc)
	AddLifecycle(ctx context.Context, name string, lifecycle ClusterLifecycle)
	AddFeatureLifecycle(ctx context.Context, enabled func() bool, name string, lifecycle ClusterLifecycle)
	AddClusterScopedHandler(ctx context.Context, name, clusterName string, sync ClusterHandlerFunc)
	AddClusterScopedFeatureHandler(ctx context.Context, enabled func() bool, name, clusterName string, sync ClusterHandlerFunc)
	AddClusterScopedLifecycle(ctx context.Context, name, clusterName string, lifecycle ClusterLifecycle)
	AddClusterScopedFeatureLifecycle(ctx context.Context, enabled func() bool, name, clusterName string, lifecycle ClusterLifecycle)
}

type clusterLister struct {
	controller *clusterController
}

func (l *clusterLister) List(namespace string, selector labels.Selector) (ret []*Cluster, err error) {
	err = cache.ListAllByNamespace(l.controller.Informer().GetIndexer(), namespace, selector, func(obj interface{}) {
		ret = append(ret, obj.(*Cluster))
	})
	return
}

func (l *clusterLister) Get(namespace, name string) (*Cluster, error) {
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
			Group:    ClusterGroupVersionKind.Group,
			Resource: "cluster",
		}, key)
	}
	return obj.(*Cluster), nil
}

type clusterController struct {
	controller.GenericController
}

func (c *clusterController) Generic() controller.GenericController {
	return c.GenericController
}

func (c *clusterController) Lister() ClusterLister {
	return &clusterLister{
		controller: c,
	}
}

func (c *clusterController) AddHandler(ctx context.Context, name string, handler ClusterHandlerFunc) {
	c.GenericController.AddHandler(ctx, name, func(key string, obj interface{}) (interface{}, error) {
		if obj == nil {
			return handler(key, nil)
		} else if v, ok := obj.(*Cluster); ok {
			return handler(key, v)
		} else {
			return nil, nil
		}
	})
}

func (c *clusterController) AddFeatureHandler(ctx context.Context, enabled func() bool, name string, handler ClusterHandlerFunc) {
	c.GenericController.AddHandler(ctx, name, func(key string, obj interface{}) (interface{}, error) {
		if !enabled() {
			return nil, nil
		} else if obj == nil {
			return handler(key, nil)
		} else if v, ok := obj.(*Cluster); ok {
			return handler(key, v)
		} else {
			return nil, nil
		}
	})
}

func (c *clusterController) AddClusterScopedHandler(ctx context.Context, name, cluster string, handler ClusterHandlerFunc) {
	c.GenericController.AddHandler(ctx, name, func(key string, obj interface{}) (interface{}, error) {
		if obj == nil {
			return handler(key, nil)
		} else if v, ok := obj.(*Cluster); ok && controller.ObjectInCluster(cluster, obj) {
			return handler(key, v)
		} else {
			return nil, nil
		}
	})
}

func (c *clusterController) AddClusterScopedFeatureHandler(ctx context.Context, enabled func() bool, name, cluster string, handler ClusterHandlerFunc) {
	c.GenericController.AddHandler(ctx, name, func(key string, obj interface{}) (interface{}, error) {
		if !enabled() {
			return nil, nil
		} else if obj == nil {
			return handler(key, nil)
		} else if v, ok := obj.(*Cluster); ok && controller.ObjectInCluster(cluster, obj) {
			return handler(key, v)
		} else {
			return nil, nil
		}
	})
}

type clusterFactory struct {
}

func (c clusterFactory) Object() runtime.Object {
	return &Cluster{}
}

func (c clusterFactory) List() runtime.Object {
	return &ClusterList{}
}

func (s *clusterClient) Controller() ClusterController {
	s.client.Lock()
	defer s.client.Unlock()

	c, ok := s.client.clusterControllers[s.ns]
	if ok {
		return c
	}

	genericController := controller.NewGenericController(ClusterGroupVersionKind.Kind+"Controller",
		s.objectClient)

	c = &clusterController{
		GenericController: genericController,
	}

	s.client.clusterControllers[s.ns] = c
	s.client.starters = append(s.client.starters, c)

	return c
}

type clusterClient struct {
	client       *Client
	ns           string
	objectClient *objectclient.ObjectClient
	controller   ClusterController
}

func (s *clusterClient) ObjectClient() *objectclient.ObjectClient {
	return s.objectClient
}

func (s *clusterClient) Create(o *Cluster) (*Cluster, error) {
	obj, err := s.objectClient.Create(o)
	return obj.(*Cluster), err
}

func (s *clusterClient) Get(name string, opts metav1.GetOptions) (*Cluster, error) {
	obj, err := s.objectClient.Get(name, opts)
	return obj.(*Cluster), err
}

func (s *clusterClient) GetNamespaced(namespace, name string, opts metav1.GetOptions) (*Cluster, error) {
	obj, err := s.objectClient.GetNamespaced(namespace, name, opts)
	return obj.(*Cluster), err
}

func (s *clusterClient) Update(o *Cluster) (*Cluster, error) {
	obj, err := s.objectClient.Update(o.Name, o)
	return obj.(*Cluster), err
}

func (s *clusterClient) Delete(name string, options *metav1.DeleteOptions) error {
	return s.objectClient.Delete(name, options)
}

func (s *clusterClient) DeleteNamespaced(namespace, name string, options *metav1.DeleteOptions) error {
	return s.objectClient.DeleteNamespaced(namespace, name, options)
}

func (s *clusterClient) List(opts metav1.ListOptions) (*ClusterList, error) {
	obj, err := s.objectClient.List(opts)
	return obj.(*ClusterList), err
}

func (s *clusterClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	return s.objectClient.Watch(opts)
}

// Patch applies the patch and returns the patched deployment.
func (s *clusterClient) Patch(o *Cluster, patchType types.PatchType, data []byte, subresources ...string) (*Cluster, error) {
	obj, err := s.objectClient.Patch(o.Name, o, patchType, data, subresources...)
	return obj.(*Cluster), err
}

func (s *clusterClient) DeleteCollection(deleteOpts *metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	return s.objectClient.DeleteCollection(deleteOpts, listOpts)
}

func (s *clusterClient) AddHandler(ctx context.Context, name string, sync ClusterHandlerFunc) {
	s.Controller().AddHandler(ctx, name, sync)
}

func (s *clusterClient) AddFeatureHandler(ctx context.Context, enabled func() bool, name string, sync ClusterHandlerFunc) {
	s.Controller().AddFeatureHandler(ctx, enabled, name, sync)
}

func (s *clusterClient) AddLifecycle(ctx context.Context, name string, lifecycle ClusterLifecycle) {
	sync := NewClusterLifecycleAdapter(name, false, s, lifecycle)
	s.Controller().AddHandler(ctx, name, sync)
}

func (s *clusterClient) AddFeatureLifecycle(ctx context.Context, enabled func() bool, name string, lifecycle ClusterLifecycle) {
	sync := NewClusterLifecycleAdapter(name, false, s, lifecycle)
	s.Controller().AddFeatureHandler(ctx, enabled, name, sync)
}

func (s *clusterClient) AddClusterScopedHandler(ctx context.Context, name, clusterName string, sync ClusterHandlerFunc) {
	s.Controller().AddClusterScopedHandler(ctx, name, clusterName, sync)
}

func (s *clusterClient) AddClusterScopedFeatureHandler(ctx context.Context, enabled func() bool, name, clusterName string, sync ClusterHandlerFunc) {
	s.Controller().AddClusterScopedFeatureHandler(ctx, enabled, name, clusterName, sync)
}

func (s *clusterClient) AddClusterScopedLifecycle(ctx context.Context, name, clusterName string, lifecycle ClusterLifecycle) {
	sync := NewClusterLifecycleAdapter(name+"_"+clusterName, true, s, lifecycle)
	s.Controller().AddClusterScopedHandler(ctx, name, clusterName, sync)
}

func (s *clusterClient) AddClusterScopedFeatureLifecycle(ctx context.Context, enabled func() bool, name, clusterName string, lifecycle ClusterLifecycle) {
	sync := NewClusterLifecycleAdapter(name+"_"+clusterName, true, s, lifecycle)
	s.Controller().AddClusterScopedFeatureHandler(ctx, enabled, name, clusterName, sync)
}

type ClusterIndexer func(obj *Cluster) ([]string, error)

type ClusterClientCache interface {
	Get(namespace, name string) (*Cluster, error)
	List(namespace string, selector labels.Selector) ([]*Cluster, error)

	Index(name string, indexer ClusterIndexer)
	GetIndexed(name, key string) ([]*Cluster, error)
}

type ClusterClient interface {
	Create(*Cluster) (*Cluster, error)
	Get(namespace, name string, opts metav1.GetOptions) (*Cluster, error)
	Update(*Cluster) (*Cluster, error)
	Delete(namespace, name string, options *metav1.DeleteOptions) error
	List(namespace string, opts metav1.ListOptions) (*ClusterList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)

	Cache() ClusterClientCache

	OnCreate(ctx context.Context, name string, sync ClusterChangeHandlerFunc)
	OnChange(ctx context.Context, name string, sync ClusterChangeHandlerFunc)
	OnRemove(ctx context.Context, name string, sync ClusterChangeHandlerFunc)
	Enqueue(namespace, name string)

	Generic() controller.GenericController
	ObjectClient() *objectclient.ObjectClient
	Interface() ClusterInterface
}

type clusterClientCache struct {
	client *clusterClient2
}

type clusterClient2 struct {
	iface      ClusterInterface
	controller ClusterController
}

func (n *clusterClient2) Interface() ClusterInterface {
	return n.iface
}

func (n *clusterClient2) Generic() controller.GenericController {
	return n.iface.Controller().Generic()
}

func (n *clusterClient2) ObjectClient() *objectclient.ObjectClient {
	return n.Interface().ObjectClient()
}

func (n *clusterClient2) Enqueue(namespace, name string) {
	n.iface.Controller().Enqueue(namespace, name)
}

func (n *clusterClient2) Create(obj *Cluster) (*Cluster, error) {
	return n.iface.Create(obj)
}

func (n *clusterClient2) Get(namespace, name string, opts metav1.GetOptions) (*Cluster, error) {
	return n.iface.GetNamespaced(namespace, name, opts)
}

func (n *clusterClient2) Update(obj *Cluster) (*Cluster, error) {
	return n.iface.Update(obj)
}

func (n *clusterClient2) Delete(namespace, name string, options *metav1.DeleteOptions) error {
	return n.iface.DeleteNamespaced(namespace, name, options)
}

func (n *clusterClient2) List(namespace string, opts metav1.ListOptions) (*ClusterList, error) {
	return n.iface.List(opts)
}

func (n *clusterClient2) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	return n.iface.Watch(opts)
}

func (n *clusterClientCache) Get(namespace, name string) (*Cluster, error) {
	return n.client.controller.Lister().Get(namespace, name)
}

func (n *clusterClientCache) List(namespace string, selector labels.Selector) ([]*Cluster, error) {
	return n.client.controller.Lister().List(namespace, selector)
}

func (n *clusterClient2) Cache() ClusterClientCache {
	n.loadController()
	return &clusterClientCache{
		client: n,
	}
}

func (n *clusterClient2) OnCreate(ctx context.Context, name string, sync ClusterChangeHandlerFunc) {
	n.loadController()
	n.iface.AddLifecycle(ctx, name+"-create", &clusterLifecycleDelegate{create: sync})
}

func (n *clusterClient2) OnChange(ctx context.Context, name string, sync ClusterChangeHandlerFunc) {
	n.loadController()
	n.iface.AddLifecycle(ctx, name+"-change", &clusterLifecycleDelegate{update: sync})
}

func (n *clusterClient2) OnRemove(ctx context.Context, name string, sync ClusterChangeHandlerFunc) {
	n.loadController()
	n.iface.AddLifecycle(ctx, name, &clusterLifecycleDelegate{remove: sync})
}

func (n *clusterClientCache) Index(name string, indexer ClusterIndexer) {
	err := n.client.controller.Informer().GetIndexer().AddIndexers(map[string]cache.IndexFunc{
		name: func(obj interface{}) ([]string, error) {
			if v, ok := obj.(*Cluster); ok {
				return indexer(v)
			}
			return nil, nil
		},
	})

	if err != nil {
		panic(err)
	}
}

func (n *clusterClientCache) GetIndexed(name, key string) ([]*Cluster, error) {
	var result []*Cluster
	objs, err := n.client.controller.Informer().GetIndexer().ByIndex(name, key)
	if err != nil {
		return nil, err
	}
	for _, obj := range objs {
		if v, ok := obj.(*Cluster); ok {
			result = append(result, v)
		}
	}

	return result, nil
}

func (n *clusterClient2) loadController() {
	if n.controller == nil {
		n.controller = n.iface.Controller()
	}
}

type clusterLifecycleDelegate struct {
	create ClusterChangeHandlerFunc
	update ClusterChangeHandlerFunc
	remove ClusterChangeHandlerFunc
}

func (n *clusterLifecycleDelegate) HasCreate() bool {
	return n.create != nil
}

func (n *clusterLifecycleDelegate) Create(obj *Cluster) (runtime.Object, error) {
	if n.create == nil {
		return obj, nil
	}
	return n.create(obj)
}

func (n *clusterLifecycleDelegate) HasFinalize() bool {
	return n.remove != nil
}

func (n *clusterLifecycleDelegate) Remove(obj *Cluster) (runtime.Object, error) {
	if n.remove == nil {
		return obj, nil
	}
	return n.remove(obj)
}

func (n *clusterLifecycleDelegate) Updated(obj *Cluster) (runtime.Object, error) {
	if n.update == nil {
		return obj, nil
	}
	return n.update(obj)
}
