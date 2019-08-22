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
	ClusterAlertGroupGroupVersionKind = schema.GroupVersionKind{
		Version: Version,
		Group:   GroupName,
		Kind:    "ClusterAlertGroup",
	}
	ClusterAlertGroupResource = metav1.APIResource{
		Name:         "clusteralertgroups",
		SingularName: "clusteralertgroup",
		Namespaced:   true,

		Kind: ClusterAlertGroupGroupVersionKind.Kind,
	}

	ClusterAlertGroupGroupVersionResource = schema.GroupVersionResource{
		Group:    GroupName,
		Version:  Version,
		Resource: "clusteralertgroups",
	}
)

func init() {
	resource.Put(ClusterAlertGroupGroupVersionResource)
}

func NewClusterAlertGroup(namespace, name string, obj ClusterAlertGroup) *ClusterAlertGroup {
	obj.APIVersion, obj.Kind = ClusterAlertGroupGroupVersionKind.ToAPIVersionAndKind()
	obj.Name = name
	obj.Namespace = namespace
	return &obj
}

type ClusterAlertGroupList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ClusterAlertGroup `json:"items"`
}

type ClusterAlertGroupHandlerFunc func(key string, obj *ClusterAlertGroup) (runtime.Object, error)

type ClusterAlertGroupChangeHandlerFunc func(obj *ClusterAlertGroup) (runtime.Object, error)

type ClusterAlertGroupLister interface {
	List(namespace string, selector labels.Selector) (ret []*ClusterAlertGroup, err error)
	Get(namespace, name string) (*ClusterAlertGroup, error)
}

type ClusterAlertGroupController interface {
	Generic() controller.GenericController
	Informer() cache.SharedIndexInformer
	Lister() ClusterAlertGroupLister
	AddHandler(ctx context.Context, name string, handler ClusterAlertGroupHandlerFunc)
	AddFeatureHandler(ctx context.Context, enabled func() bool, name string, sync ClusterAlertGroupHandlerFunc)
	AddClusterScopedHandler(ctx context.Context, name, clusterName string, handler ClusterAlertGroupHandlerFunc)
	AddClusterScopedFeatureHandler(ctx context.Context, enabled func() bool, name, clusterName string, handler ClusterAlertGroupHandlerFunc)
	Enqueue(namespace, name string)
	Sync(ctx context.Context) error
	Start(ctx context.Context, threadiness int) error
}

type ClusterAlertGroupInterface interface {
	ObjectClient() *objectclient.ObjectClient
	Create(*ClusterAlertGroup) (*ClusterAlertGroup, error)
	GetNamespaced(namespace, name string, opts metav1.GetOptions) (*ClusterAlertGroup, error)
	Get(name string, opts metav1.GetOptions) (*ClusterAlertGroup, error)
	Update(*ClusterAlertGroup) (*ClusterAlertGroup, error)
	Delete(name string, options *metav1.DeleteOptions) error
	DeleteNamespaced(namespace, name string, options *metav1.DeleteOptions) error
	List(opts metav1.ListOptions) (*ClusterAlertGroupList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	DeleteCollection(deleteOpts *metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Controller() ClusterAlertGroupController
	AddHandler(ctx context.Context, name string, sync ClusterAlertGroupHandlerFunc)
	AddFeatureHandler(ctx context.Context, enabled func() bool, name string, sync ClusterAlertGroupHandlerFunc)
	AddLifecycle(ctx context.Context, name string, lifecycle ClusterAlertGroupLifecycle)
	AddFeatureLifecycle(ctx context.Context, enabled func() bool, name string, lifecycle ClusterAlertGroupLifecycle)
	AddClusterScopedHandler(ctx context.Context, name, clusterName string, sync ClusterAlertGroupHandlerFunc)
	AddClusterScopedFeatureHandler(ctx context.Context, enabled func() bool, name, clusterName string, sync ClusterAlertGroupHandlerFunc)
	AddClusterScopedLifecycle(ctx context.Context, name, clusterName string, lifecycle ClusterAlertGroupLifecycle)
	AddClusterScopedFeatureLifecycle(ctx context.Context, enabled func() bool, name, clusterName string, lifecycle ClusterAlertGroupLifecycle)
}

type clusterAlertGroupLister struct {
	controller *clusterAlertGroupController
}

func (l *clusterAlertGroupLister) List(namespace string, selector labels.Selector) (ret []*ClusterAlertGroup, err error) {
	err = cache.ListAllByNamespace(l.controller.Informer().GetIndexer(), namespace, selector, func(obj interface{}) {
		ret = append(ret, obj.(*ClusterAlertGroup))
	})
	return
}

func (l *clusterAlertGroupLister) Get(namespace, name string) (*ClusterAlertGroup, error) {
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
			Group:    ClusterAlertGroupGroupVersionKind.Group,
			Resource: "clusterAlertGroup",
		}, key)
	}
	return obj.(*ClusterAlertGroup), nil
}

type clusterAlertGroupController struct {
	controller.GenericController
}

func (c *clusterAlertGroupController) Generic() controller.GenericController {
	return c.GenericController
}

func (c *clusterAlertGroupController) Lister() ClusterAlertGroupLister {
	return &clusterAlertGroupLister{
		controller: c,
	}
}

func (c *clusterAlertGroupController) AddHandler(ctx context.Context, name string, handler ClusterAlertGroupHandlerFunc) {
	c.GenericController.AddHandler(ctx, name, func(key string, obj interface{}) (interface{}, error) {
		if obj == nil {
			return handler(key, nil)
		} else if v, ok := obj.(*ClusterAlertGroup); ok {
			return handler(key, v)
		} else {
			return nil, nil
		}
	})
}

func (c *clusterAlertGroupController) AddFeatureHandler(ctx context.Context, enabled func() bool, name string, handler ClusterAlertGroupHandlerFunc) {
	c.GenericController.AddHandler(ctx, name, func(key string, obj interface{}) (interface{}, error) {
		if !enabled() {
			return nil, nil
		} else if obj == nil {
			return handler(key, nil)
		} else if v, ok := obj.(*ClusterAlertGroup); ok {
			return handler(key, v)
		} else {
			return nil, nil
		}
	})
}

func (c *clusterAlertGroupController) AddClusterScopedHandler(ctx context.Context, name, cluster string, handler ClusterAlertGroupHandlerFunc) {
	c.GenericController.AddHandler(ctx, name, func(key string, obj interface{}) (interface{}, error) {
		if obj == nil {
			return handler(key, nil)
		} else if v, ok := obj.(*ClusterAlertGroup); ok && controller.ObjectInCluster(cluster, obj) {
			return handler(key, v)
		} else {
			return nil, nil
		}
	})
}

func (c *clusterAlertGroupController) AddClusterScopedFeatureHandler(ctx context.Context, enabled func() bool, name, cluster string, handler ClusterAlertGroupHandlerFunc) {
	c.GenericController.AddHandler(ctx, name, func(key string, obj interface{}) (interface{}, error) {
		if !enabled() {
			return nil, nil
		} else if obj == nil {
			return handler(key, nil)
		} else if v, ok := obj.(*ClusterAlertGroup); ok && controller.ObjectInCluster(cluster, obj) {
			return handler(key, v)
		} else {
			return nil, nil
		}
	})
}

type clusterAlertGroupFactory struct {
}

func (c clusterAlertGroupFactory) Object() runtime.Object {
	return &ClusterAlertGroup{}
}

func (c clusterAlertGroupFactory) List() runtime.Object {
	return &ClusterAlertGroupList{}
}

func (s *clusterAlertGroupClient) Controller() ClusterAlertGroupController {
	s.client.Lock()
	defer s.client.Unlock()

	c, ok := s.client.clusterAlertGroupControllers[s.ns]
	if ok {
		return c
	}

	genericController := controller.NewGenericController(ClusterAlertGroupGroupVersionKind.Kind+"Controller",
		s.objectClient)

	c = &clusterAlertGroupController{
		GenericController: genericController,
	}

	s.client.clusterAlertGroupControllers[s.ns] = c
	s.client.starters = append(s.client.starters, c)

	return c
}

type clusterAlertGroupClient struct {
	client       *Client
	ns           string
	objectClient *objectclient.ObjectClient
	controller   ClusterAlertGroupController
}

func (s *clusterAlertGroupClient) ObjectClient() *objectclient.ObjectClient {
	return s.objectClient
}

func (s *clusterAlertGroupClient) Create(o *ClusterAlertGroup) (*ClusterAlertGroup, error) {
	obj, err := s.objectClient.Create(o)
	return obj.(*ClusterAlertGroup), err
}

func (s *clusterAlertGroupClient) Get(name string, opts metav1.GetOptions) (*ClusterAlertGroup, error) {
	obj, err := s.objectClient.Get(name, opts)
	return obj.(*ClusterAlertGroup), err
}

func (s *clusterAlertGroupClient) GetNamespaced(namespace, name string, opts metav1.GetOptions) (*ClusterAlertGroup, error) {
	obj, err := s.objectClient.GetNamespaced(namespace, name, opts)
	return obj.(*ClusterAlertGroup), err
}

func (s *clusterAlertGroupClient) Update(o *ClusterAlertGroup) (*ClusterAlertGroup, error) {
	obj, err := s.objectClient.Update(o.Name, o)
	return obj.(*ClusterAlertGroup), err
}

func (s *clusterAlertGroupClient) Delete(name string, options *metav1.DeleteOptions) error {
	return s.objectClient.Delete(name, options)
}

func (s *clusterAlertGroupClient) DeleteNamespaced(namespace, name string, options *metav1.DeleteOptions) error {
	return s.objectClient.DeleteNamespaced(namespace, name, options)
}

func (s *clusterAlertGroupClient) List(opts metav1.ListOptions) (*ClusterAlertGroupList, error) {
	obj, err := s.objectClient.List(opts)
	return obj.(*ClusterAlertGroupList), err
}

func (s *clusterAlertGroupClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	return s.objectClient.Watch(opts)
}

// Patch applies the patch and returns the patched deployment.
func (s *clusterAlertGroupClient) Patch(o *ClusterAlertGroup, patchType types.PatchType, data []byte, subresources ...string) (*ClusterAlertGroup, error) {
	obj, err := s.objectClient.Patch(o.Name, o, patchType, data, subresources...)
	return obj.(*ClusterAlertGroup), err
}

func (s *clusterAlertGroupClient) DeleteCollection(deleteOpts *metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	return s.objectClient.DeleteCollection(deleteOpts, listOpts)
}

func (s *clusterAlertGroupClient) AddHandler(ctx context.Context, name string, sync ClusterAlertGroupHandlerFunc) {
	s.Controller().AddHandler(ctx, name, sync)
}

func (s *clusterAlertGroupClient) AddFeatureHandler(ctx context.Context, enabled func() bool, name string, sync ClusterAlertGroupHandlerFunc) {
	s.Controller().AddFeatureHandler(ctx, enabled, name, sync)
}

func (s *clusterAlertGroupClient) AddLifecycle(ctx context.Context, name string, lifecycle ClusterAlertGroupLifecycle) {
	sync := NewClusterAlertGroupLifecycleAdapter(name, false, s, lifecycle)
	s.Controller().AddHandler(ctx, name, sync)
}

func (s *clusterAlertGroupClient) AddFeatureLifecycle(ctx context.Context, enabled func() bool, name string, lifecycle ClusterAlertGroupLifecycle) {
	sync := NewClusterAlertGroupLifecycleAdapter(name, false, s, lifecycle)
	s.Controller().AddFeatureHandler(ctx, enabled, name, sync)
}

func (s *clusterAlertGroupClient) AddClusterScopedHandler(ctx context.Context, name, clusterName string, sync ClusterAlertGroupHandlerFunc) {
	s.Controller().AddClusterScopedHandler(ctx, name, clusterName, sync)
}

func (s *clusterAlertGroupClient) AddClusterScopedFeatureHandler(ctx context.Context, enabled func() bool, name, clusterName string, sync ClusterAlertGroupHandlerFunc) {
	s.Controller().AddClusterScopedFeatureHandler(ctx, enabled, name, clusterName, sync)
}

func (s *clusterAlertGroupClient) AddClusterScopedLifecycle(ctx context.Context, name, clusterName string, lifecycle ClusterAlertGroupLifecycle) {
	sync := NewClusterAlertGroupLifecycleAdapter(name+"_"+clusterName, true, s, lifecycle)
	s.Controller().AddClusterScopedHandler(ctx, name, clusterName, sync)
}

func (s *clusterAlertGroupClient) AddClusterScopedFeatureLifecycle(ctx context.Context, enabled func() bool, name, clusterName string, lifecycle ClusterAlertGroupLifecycle) {
	sync := NewClusterAlertGroupLifecycleAdapter(name+"_"+clusterName, true, s, lifecycle)
	s.Controller().AddClusterScopedFeatureHandler(ctx, enabled, name, clusterName, sync)
}

type ClusterAlertGroupIndexer func(obj *ClusterAlertGroup) ([]string, error)

type ClusterAlertGroupClientCache interface {
	Get(namespace, name string) (*ClusterAlertGroup, error)
	List(namespace string, selector labels.Selector) ([]*ClusterAlertGroup, error)

	Index(name string, indexer ClusterAlertGroupIndexer)
	GetIndexed(name, key string) ([]*ClusterAlertGroup, error)
}

type ClusterAlertGroupClient interface {
	Create(*ClusterAlertGroup) (*ClusterAlertGroup, error)
	Get(namespace, name string, opts metav1.GetOptions) (*ClusterAlertGroup, error)
	Update(*ClusterAlertGroup) (*ClusterAlertGroup, error)
	Delete(namespace, name string, options *metav1.DeleteOptions) error
	List(namespace string, opts metav1.ListOptions) (*ClusterAlertGroupList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)

	Cache() ClusterAlertGroupClientCache

	OnCreate(ctx context.Context, name string, sync ClusterAlertGroupChangeHandlerFunc)
	OnChange(ctx context.Context, name string, sync ClusterAlertGroupChangeHandlerFunc)
	OnRemove(ctx context.Context, name string, sync ClusterAlertGroupChangeHandlerFunc)
	Enqueue(namespace, name string)

	Generic() controller.GenericController
	ObjectClient() *objectclient.ObjectClient
	Interface() ClusterAlertGroupInterface
}

type clusterAlertGroupClientCache struct {
	client *clusterAlertGroupClient2
}

type clusterAlertGroupClient2 struct {
	iface      ClusterAlertGroupInterface
	controller ClusterAlertGroupController
}

func (n *clusterAlertGroupClient2) Interface() ClusterAlertGroupInterface {
	return n.iface
}

func (n *clusterAlertGroupClient2) Generic() controller.GenericController {
	return n.iface.Controller().Generic()
}

func (n *clusterAlertGroupClient2) ObjectClient() *objectclient.ObjectClient {
	return n.Interface().ObjectClient()
}

func (n *clusterAlertGroupClient2) Enqueue(namespace, name string) {
	n.iface.Controller().Enqueue(namespace, name)
}

func (n *clusterAlertGroupClient2) Create(obj *ClusterAlertGroup) (*ClusterAlertGroup, error) {
	return n.iface.Create(obj)
}

func (n *clusterAlertGroupClient2) Get(namespace, name string, opts metav1.GetOptions) (*ClusterAlertGroup, error) {
	return n.iface.GetNamespaced(namespace, name, opts)
}

func (n *clusterAlertGroupClient2) Update(obj *ClusterAlertGroup) (*ClusterAlertGroup, error) {
	return n.iface.Update(obj)
}

func (n *clusterAlertGroupClient2) Delete(namespace, name string, options *metav1.DeleteOptions) error {
	return n.iface.DeleteNamespaced(namespace, name, options)
}

func (n *clusterAlertGroupClient2) List(namespace string, opts metav1.ListOptions) (*ClusterAlertGroupList, error) {
	return n.iface.List(opts)
}

func (n *clusterAlertGroupClient2) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	return n.iface.Watch(opts)
}

func (n *clusterAlertGroupClientCache) Get(namespace, name string) (*ClusterAlertGroup, error) {
	return n.client.controller.Lister().Get(namespace, name)
}

func (n *clusterAlertGroupClientCache) List(namespace string, selector labels.Selector) ([]*ClusterAlertGroup, error) {
	return n.client.controller.Lister().List(namespace, selector)
}

func (n *clusterAlertGroupClient2) Cache() ClusterAlertGroupClientCache {
	n.loadController()
	return &clusterAlertGroupClientCache{
		client: n,
	}
}

func (n *clusterAlertGroupClient2) OnCreate(ctx context.Context, name string, sync ClusterAlertGroupChangeHandlerFunc) {
	n.loadController()
	n.iface.AddLifecycle(ctx, name+"-create", &clusterAlertGroupLifecycleDelegate{create: sync})
}

func (n *clusterAlertGroupClient2) OnChange(ctx context.Context, name string, sync ClusterAlertGroupChangeHandlerFunc) {
	n.loadController()
	n.iface.AddLifecycle(ctx, name+"-change", &clusterAlertGroupLifecycleDelegate{update: sync})
}

func (n *clusterAlertGroupClient2) OnRemove(ctx context.Context, name string, sync ClusterAlertGroupChangeHandlerFunc) {
	n.loadController()
	n.iface.AddLifecycle(ctx, name, &clusterAlertGroupLifecycleDelegate{remove: sync})
}

func (n *clusterAlertGroupClientCache) Index(name string, indexer ClusterAlertGroupIndexer) {
	err := n.client.controller.Informer().GetIndexer().AddIndexers(map[string]cache.IndexFunc{
		name: func(obj interface{}) ([]string, error) {
			if v, ok := obj.(*ClusterAlertGroup); ok {
				return indexer(v)
			}
			return nil, nil
		},
	})

	if err != nil {
		panic(err)
	}
}

func (n *clusterAlertGroupClientCache) GetIndexed(name, key string) ([]*ClusterAlertGroup, error) {
	var result []*ClusterAlertGroup
	objs, err := n.client.controller.Informer().GetIndexer().ByIndex(name, key)
	if err != nil {
		return nil, err
	}
	for _, obj := range objs {
		if v, ok := obj.(*ClusterAlertGroup); ok {
			result = append(result, v)
		}
	}

	return result, nil
}

func (n *clusterAlertGroupClient2) loadController() {
	if n.controller == nil {
		n.controller = n.iface.Controller()
	}
}

type clusterAlertGroupLifecycleDelegate struct {
	create ClusterAlertGroupChangeHandlerFunc
	update ClusterAlertGroupChangeHandlerFunc
	remove ClusterAlertGroupChangeHandlerFunc
}

func (n *clusterAlertGroupLifecycleDelegate) HasCreate() bool {
	return n.create != nil
}

func (n *clusterAlertGroupLifecycleDelegate) Create(obj *ClusterAlertGroup) (runtime.Object, error) {
	if n.create == nil {
		return obj, nil
	}
	return n.create(obj)
}

func (n *clusterAlertGroupLifecycleDelegate) HasFinalize() bool {
	return n.remove != nil
}

func (n *clusterAlertGroupLifecycleDelegate) Remove(obj *ClusterAlertGroup) (runtime.Object, error) {
	if n.remove == nil {
		return obj, nil
	}
	return n.remove(obj)
}

func (n *clusterAlertGroupLifecycleDelegate) Updated(obj *ClusterAlertGroup) (runtime.Object, error) {
	if n.update == nil {
		return obj, nil
	}
	return n.update(obj)
}
