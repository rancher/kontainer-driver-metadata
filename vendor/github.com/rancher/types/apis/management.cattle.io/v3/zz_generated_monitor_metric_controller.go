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
	MonitorMetricGroupVersionKind = schema.GroupVersionKind{
		Version: Version,
		Group:   GroupName,
		Kind:    "MonitorMetric",
	}
	MonitorMetricResource = metav1.APIResource{
		Name:         "monitormetrics",
		SingularName: "monitormetric",
		Namespaced:   true,

		Kind: MonitorMetricGroupVersionKind.Kind,
	}

	MonitorMetricGroupVersionResource = schema.GroupVersionResource{
		Group:    GroupName,
		Version:  Version,
		Resource: "monitormetrics",
	}
)

func init() {
	resource.Put(MonitorMetricGroupVersionResource)
}

func NewMonitorMetric(namespace, name string, obj MonitorMetric) *MonitorMetric {
	obj.APIVersion, obj.Kind = MonitorMetricGroupVersionKind.ToAPIVersionAndKind()
	obj.Name = name
	obj.Namespace = namespace
	return &obj
}

type MonitorMetricList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MonitorMetric `json:"items"`
}

type MonitorMetricHandlerFunc func(key string, obj *MonitorMetric) (runtime.Object, error)

type MonitorMetricChangeHandlerFunc func(obj *MonitorMetric) (runtime.Object, error)

type MonitorMetricLister interface {
	List(namespace string, selector labels.Selector) (ret []*MonitorMetric, err error)
	Get(namespace, name string) (*MonitorMetric, error)
}

type MonitorMetricController interface {
	Generic() controller.GenericController
	Informer() cache.SharedIndexInformer
	Lister() MonitorMetricLister
	AddHandler(ctx context.Context, name string, handler MonitorMetricHandlerFunc)
	AddFeatureHandler(ctx context.Context, enabled func() bool, name string, sync MonitorMetricHandlerFunc)
	AddClusterScopedHandler(ctx context.Context, name, clusterName string, handler MonitorMetricHandlerFunc)
	AddClusterScopedFeatureHandler(ctx context.Context, enabled func() bool, name, clusterName string, handler MonitorMetricHandlerFunc)
	Enqueue(namespace, name string)
	Sync(ctx context.Context) error
	Start(ctx context.Context, threadiness int) error
}

type MonitorMetricInterface interface {
	ObjectClient() *objectclient.ObjectClient
	Create(*MonitorMetric) (*MonitorMetric, error)
	GetNamespaced(namespace, name string, opts metav1.GetOptions) (*MonitorMetric, error)
	Get(name string, opts metav1.GetOptions) (*MonitorMetric, error)
	Update(*MonitorMetric) (*MonitorMetric, error)
	Delete(name string, options *metav1.DeleteOptions) error
	DeleteNamespaced(namespace, name string, options *metav1.DeleteOptions) error
	List(opts metav1.ListOptions) (*MonitorMetricList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	DeleteCollection(deleteOpts *metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Controller() MonitorMetricController
	AddHandler(ctx context.Context, name string, sync MonitorMetricHandlerFunc)
	AddFeatureHandler(ctx context.Context, enabled func() bool, name string, sync MonitorMetricHandlerFunc)
	AddLifecycle(ctx context.Context, name string, lifecycle MonitorMetricLifecycle)
	AddFeatureLifecycle(ctx context.Context, enabled func() bool, name string, lifecycle MonitorMetricLifecycle)
	AddClusterScopedHandler(ctx context.Context, name, clusterName string, sync MonitorMetricHandlerFunc)
	AddClusterScopedFeatureHandler(ctx context.Context, enabled func() bool, name, clusterName string, sync MonitorMetricHandlerFunc)
	AddClusterScopedLifecycle(ctx context.Context, name, clusterName string, lifecycle MonitorMetricLifecycle)
	AddClusterScopedFeatureLifecycle(ctx context.Context, enabled func() bool, name, clusterName string, lifecycle MonitorMetricLifecycle)
}

type monitorMetricLister struct {
	controller *monitorMetricController
}

func (l *monitorMetricLister) List(namespace string, selector labels.Selector) (ret []*MonitorMetric, err error) {
	err = cache.ListAllByNamespace(l.controller.Informer().GetIndexer(), namespace, selector, func(obj interface{}) {
		ret = append(ret, obj.(*MonitorMetric))
	})
	return
}

func (l *monitorMetricLister) Get(namespace, name string) (*MonitorMetric, error) {
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
			Group:    MonitorMetricGroupVersionKind.Group,
			Resource: "monitorMetric",
		}, key)
	}
	return obj.(*MonitorMetric), nil
}

type monitorMetricController struct {
	controller.GenericController
}

func (c *monitorMetricController) Generic() controller.GenericController {
	return c.GenericController
}

func (c *monitorMetricController) Lister() MonitorMetricLister {
	return &monitorMetricLister{
		controller: c,
	}
}

func (c *monitorMetricController) AddHandler(ctx context.Context, name string, handler MonitorMetricHandlerFunc) {
	c.GenericController.AddHandler(ctx, name, func(key string, obj interface{}) (interface{}, error) {
		if obj == nil {
			return handler(key, nil)
		} else if v, ok := obj.(*MonitorMetric); ok {
			return handler(key, v)
		} else {
			return nil, nil
		}
	})
}

func (c *monitorMetricController) AddFeatureHandler(ctx context.Context, enabled func() bool, name string, handler MonitorMetricHandlerFunc) {
	c.GenericController.AddHandler(ctx, name, func(key string, obj interface{}) (interface{}, error) {
		if !enabled() {
			return nil, nil
		} else if obj == nil {
			return handler(key, nil)
		} else if v, ok := obj.(*MonitorMetric); ok {
			return handler(key, v)
		} else {
			return nil, nil
		}
	})
}

func (c *monitorMetricController) AddClusterScopedHandler(ctx context.Context, name, cluster string, handler MonitorMetricHandlerFunc) {
	c.GenericController.AddHandler(ctx, name, func(key string, obj interface{}) (interface{}, error) {
		if obj == nil {
			return handler(key, nil)
		} else if v, ok := obj.(*MonitorMetric); ok && controller.ObjectInCluster(cluster, obj) {
			return handler(key, v)
		} else {
			return nil, nil
		}
	})
}

func (c *monitorMetricController) AddClusterScopedFeatureHandler(ctx context.Context, enabled func() bool, name, cluster string, handler MonitorMetricHandlerFunc) {
	c.GenericController.AddHandler(ctx, name, func(key string, obj interface{}) (interface{}, error) {
		if !enabled() {
			return nil, nil
		} else if obj == nil {
			return handler(key, nil)
		} else if v, ok := obj.(*MonitorMetric); ok && controller.ObjectInCluster(cluster, obj) {
			return handler(key, v)
		} else {
			return nil, nil
		}
	})
}

type monitorMetricFactory struct {
}

func (c monitorMetricFactory) Object() runtime.Object {
	return &MonitorMetric{}
}

func (c monitorMetricFactory) List() runtime.Object {
	return &MonitorMetricList{}
}

func (s *monitorMetricClient) Controller() MonitorMetricController {
	s.client.Lock()
	defer s.client.Unlock()

	c, ok := s.client.monitorMetricControllers[s.ns]
	if ok {
		return c
	}

	genericController := controller.NewGenericController(MonitorMetricGroupVersionKind.Kind+"Controller",
		s.objectClient)

	c = &monitorMetricController{
		GenericController: genericController,
	}

	s.client.monitorMetricControllers[s.ns] = c
	s.client.starters = append(s.client.starters, c)

	return c
}

type monitorMetricClient struct {
	client       *Client
	ns           string
	objectClient *objectclient.ObjectClient
	controller   MonitorMetricController
}

func (s *monitorMetricClient) ObjectClient() *objectclient.ObjectClient {
	return s.objectClient
}

func (s *monitorMetricClient) Create(o *MonitorMetric) (*MonitorMetric, error) {
	obj, err := s.objectClient.Create(o)
	return obj.(*MonitorMetric), err
}

func (s *monitorMetricClient) Get(name string, opts metav1.GetOptions) (*MonitorMetric, error) {
	obj, err := s.objectClient.Get(name, opts)
	return obj.(*MonitorMetric), err
}

func (s *monitorMetricClient) GetNamespaced(namespace, name string, opts metav1.GetOptions) (*MonitorMetric, error) {
	obj, err := s.objectClient.GetNamespaced(namespace, name, opts)
	return obj.(*MonitorMetric), err
}

func (s *monitorMetricClient) Update(o *MonitorMetric) (*MonitorMetric, error) {
	obj, err := s.objectClient.Update(o.Name, o)
	return obj.(*MonitorMetric), err
}

func (s *monitorMetricClient) Delete(name string, options *metav1.DeleteOptions) error {
	return s.objectClient.Delete(name, options)
}

func (s *monitorMetricClient) DeleteNamespaced(namespace, name string, options *metav1.DeleteOptions) error {
	return s.objectClient.DeleteNamespaced(namespace, name, options)
}

func (s *monitorMetricClient) List(opts metav1.ListOptions) (*MonitorMetricList, error) {
	obj, err := s.objectClient.List(opts)
	return obj.(*MonitorMetricList), err
}

func (s *monitorMetricClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	return s.objectClient.Watch(opts)
}

// Patch applies the patch and returns the patched deployment.
func (s *monitorMetricClient) Patch(o *MonitorMetric, patchType types.PatchType, data []byte, subresources ...string) (*MonitorMetric, error) {
	obj, err := s.objectClient.Patch(o.Name, o, patchType, data, subresources...)
	return obj.(*MonitorMetric), err
}

func (s *monitorMetricClient) DeleteCollection(deleteOpts *metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	return s.objectClient.DeleteCollection(deleteOpts, listOpts)
}

func (s *monitorMetricClient) AddHandler(ctx context.Context, name string, sync MonitorMetricHandlerFunc) {
	s.Controller().AddHandler(ctx, name, sync)
}

func (s *monitorMetricClient) AddFeatureHandler(ctx context.Context, enabled func() bool, name string, sync MonitorMetricHandlerFunc) {
	s.Controller().AddFeatureHandler(ctx, enabled, name, sync)
}

func (s *monitorMetricClient) AddLifecycle(ctx context.Context, name string, lifecycle MonitorMetricLifecycle) {
	sync := NewMonitorMetricLifecycleAdapter(name, false, s, lifecycle)
	s.Controller().AddHandler(ctx, name, sync)
}

func (s *monitorMetricClient) AddFeatureLifecycle(ctx context.Context, enabled func() bool, name string, lifecycle MonitorMetricLifecycle) {
	sync := NewMonitorMetricLifecycleAdapter(name, false, s, lifecycle)
	s.Controller().AddFeatureHandler(ctx, enabled, name, sync)
}

func (s *monitorMetricClient) AddClusterScopedHandler(ctx context.Context, name, clusterName string, sync MonitorMetricHandlerFunc) {
	s.Controller().AddClusterScopedHandler(ctx, name, clusterName, sync)
}

func (s *monitorMetricClient) AddClusterScopedFeatureHandler(ctx context.Context, enabled func() bool, name, clusterName string, sync MonitorMetricHandlerFunc) {
	s.Controller().AddClusterScopedFeatureHandler(ctx, enabled, name, clusterName, sync)
}

func (s *monitorMetricClient) AddClusterScopedLifecycle(ctx context.Context, name, clusterName string, lifecycle MonitorMetricLifecycle) {
	sync := NewMonitorMetricLifecycleAdapter(name+"_"+clusterName, true, s, lifecycle)
	s.Controller().AddClusterScopedHandler(ctx, name, clusterName, sync)
}

func (s *monitorMetricClient) AddClusterScopedFeatureLifecycle(ctx context.Context, enabled func() bool, name, clusterName string, lifecycle MonitorMetricLifecycle) {
	sync := NewMonitorMetricLifecycleAdapter(name+"_"+clusterName, true, s, lifecycle)
	s.Controller().AddClusterScopedFeatureHandler(ctx, enabled, name, clusterName, sync)
}

type MonitorMetricIndexer func(obj *MonitorMetric) ([]string, error)

type MonitorMetricClientCache interface {
	Get(namespace, name string) (*MonitorMetric, error)
	List(namespace string, selector labels.Selector) ([]*MonitorMetric, error)

	Index(name string, indexer MonitorMetricIndexer)
	GetIndexed(name, key string) ([]*MonitorMetric, error)
}

type MonitorMetricClient interface {
	Create(*MonitorMetric) (*MonitorMetric, error)
	Get(namespace, name string, opts metav1.GetOptions) (*MonitorMetric, error)
	Update(*MonitorMetric) (*MonitorMetric, error)
	Delete(namespace, name string, options *metav1.DeleteOptions) error
	List(namespace string, opts metav1.ListOptions) (*MonitorMetricList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)

	Cache() MonitorMetricClientCache

	OnCreate(ctx context.Context, name string, sync MonitorMetricChangeHandlerFunc)
	OnChange(ctx context.Context, name string, sync MonitorMetricChangeHandlerFunc)
	OnRemove(ctx context.Context, name string, sync MonitorMetricChangeHandlerFunc)
	Enqueue(namespace, name string)

	Generic() controller.GenericController
	ObjectClient() *objectclient.ObjectClient
	Interface() MonitorMetricInterface
}

type monitorMetricClientCache struct {
	client *monitorMetricClient2
}

type monitorMetricClient2 struct {
	iface      MonitorMetricInterface
	controller MonitorMetricController
}

func (n *monitorMetricClient2) Interface() MonitorMetricInterface {
	return n.iface
}

func (n *monitorMetricClient2) Generic() controller.GenericController {
	return n.iface.Controller().Generic()
}

func (n *monitorMetricClient2) ObjectClient() *objectclient.ObjectClient {
	return n.Interface().ObjectClient()
}

func (n *monitorMetricClient2) Enqueue(namespace, name string) {
	n.iface.Controller().Enqueue(namespace, name)
}

func (n *monitorMetricClient2) Create(obj *MonitorMetric) (*MonitorMetric, error) {
	return n.iface.Create(obj)
}

func (n *monitorMetricClient2) Get(namespace, name string, opts metav1.GetOptions) (*MonitorMetric, error) {
	return n.iface.GetNamespaced(namespace, name, opts)
}

func (n *monitorMetricClient2) Update(obj *MonitorMetric) (*MonitorMetric, error) {
	return n.iface.Update(obj)
}

func (n *monitorMetricClient2) Delete(namespace, name string, options *metav1.DeleteOptions) error {
	return n.iface.DeleteNamespaced(namespace, name, options)
}

func (n *monitorMetricClient2) List(namespace string, opts metav1.ListOptions) (*MonitorMetricList, error) {
	return n.iface.List(opts)
}

func (n *monitorMetricClient2) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	return n.iface.Watch(opts)
}

func (n *monitorMetricClientCache) Get(namespace, name string) (*MonitorMetric, error) {
	return n.client.controller.Lister().Get(namespace, name)
}

func (n *monitorMetricClientCache) List(namespace string, selector labels.Selector) ([]*MonitorMetric, error) {
	return n.client.controller.Lister().List(namespace, selector)
}

func (n *monitorMetricClient2) Cache() MonitorMetricClientCache {
	n.loadController()
	return &monitorMetricClientCache{
		client: n,
	}
}

func (n *monitorMetricClient2) OnCreate(ctx context.Context, name string, sync MonitorMetricChangeHandlerFunc) {
	n.loadController()
	n.iface.AddLifecycle(ctx, name+"-create", &monitorMetricLifecycleDelegate{create: sync})
}

func (n *monitorMetricClient2) OnChange(ctx context.Context, name string, sync MonitorMetricChangeHandlerFunc) {
	n.loadController()
	n.iface.AddLifecycle(ctx, name+"-change", &monitorMetricLifecycleDelegate{update: sync})
}

func (n *monitorMetricClient2) OnRemove(ctx context.Context, name string, sync MonitorMetricChangeHandlerFunc) {
	n.loadController()
	n.iface.AddLifecycle(ctx, name, &monitorMetricLifecycleDelegate{remove: sync})
}

func (n *monitorMetricClientCache) Index(name string, indexer MonitorMetricIndexer) {
	err := n.client.controller.Informer().GetIndexer().AddIndexers(map[string]cache.IndexFunc{
		name: func(obj interface{}) ([]string, error) {
			if v, ok := obj.(*MonitorMetric); ok {
				return indexer(v)
			}
			return nil, nil
		},
	})

	if err != nil {
		panic(err)
	}
}

func (n *monitorMetricClientCache) GetIndexed(name, key string) ([]*MonitorMetric, error) {
	var result []*MonitorMetric
	objs, err := n.client.controller.Informer().GetIndexer().ByIndex(name, key)
	if err != nil {
		return nil, err
	}
	for _, obj := range objs {
		if v, ok := obj.(*MonitorMetric); ok {
			result = append(result, v)
		}
	}

	return result, nil
}

func (n *monitorMetricClient2) loadController() {
	if n.controller == nil {
		n.controller = n.iface.Controller()
	}
}

type monitorMetricLifecycleDelegate struct {
	create MonitorMetricChangeHandlerFunc
	update MonitorMetricChangeHandlerFunc
	remove MonitorMetricChangeHandlerFunc
}

func (n *monitorMetricLifecycleDelegate) HasCreate() bool {
	return n.create != nil
}

func (n *monitorMetricLifecycleDelegate) Create(obj *MonitorMetric) (runtime.Object, error) {
	if n.create == nil {
		return obj, nil
	}
	return n.create(obj)
}

func (n *monitorMetricLifecycleDelegate) HasFinalize() bool {
	return n.remove != nil
}

func (n *monitorMetricLifecycleDelegate) Remove(obj *MonitorMetric) (runtime.Object, error) {
	if n.remove == nil {
		return obj, nil
	}
	return n.remove(obj)
}

func (n *monitorMetricLifecycleDelegate) Updated(obj *MonitorMetric) (runtime.Object, error) {
	if n.update == nil {
		return obj, nil
	}
	return n.update(obj)
}
