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
	GroupMemberGroupVersionKind = schema.GroupVersionKind{
		Version: Version,
		Group:   GroupName,
		Kind:    "GroupMember",
	}
	GroupMemberResource = metav1.APIResource{
		Name:         "groupmembers",
		SingularName: "groupmember",
		Namespaced:   false,
		Kind:         GroupMemberGroupVersionKind.Kind,
	}

	GroupMemberGroupVersionResource = schema.GroupVersionResource{
		Group:    GroupName,
		Version:  Version,
		Resource: "groupmembers",
	}
)

func init() {
	resource.Put(GroupMemberGroupVersionResource)
}

func NewGroupMember(namespace, name string, obj GroupMember) *GroupMember {
	obj.APIVersion, obj.Kind = GroupMemberGroupVersionKind.ToAPIVersionAndKind()
	obj.Name = name
	obj.Namespace = namespace
	return &obj
}

type GroupMemberList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []GroupMember `json:"items"`
}

type GroupMemberHandlerFunc func(key string, obj *GroupMember) (runtime.Object, error)

type GroupMemberChangeHandlerFunc func(obj *GroupMember) (runtime.Object, error)

type GroupMemberLister interface {
	List(namespace string, selector labels.Selector) (ret []*GroupMember, err error)
	Get(namespace, name string) (*GroupMember, error)
}

type GroupMemberController interface {
	Generic() controller.GenericController
	Informer() cache.SharedIndexInformer
	Lister() GroupMemberLister
	AddHandler(ctx context.Context, name string, handler GroupMemberHandlerFunc)
	AddFeatureHandler(ctx context.Context, enabled func() bool, name string, sync GroupMemberHandlerFunc)
	AddClusterScopedHandler(ctx context.Context, name, clusterName string, handler GroupMemberHandlerFunc)
	AddClusterScopedFeatureHandler(ctx context.Context, enabled func() bool, name, clusterName string, handler GroupMemberHandlerFunc)
	Enqueue(namespace, name string)
	Sync(ctx context.Context) error
	Start(ctx context.Context, threadiness int) error
}

type GroupMemberInterface interface {
	ObjectClient() *objectclient.ObjectClient
	Create(*GroupMember) (*GroupMember, error)
	GetNamespaced(namespace, name string, opts metav1.GetOptions) (*GroupMember, error)
	Get(name string, opts metav1.GetOptions) (*GroupMember, error)
	Update(*GroupMember) (*GroupMember, error)
	Delete(name string, options *metav1.DeleteOptions) error
	DeleteNamespaced(namespace, name string, options *metav1.DeleteOptions) error
	List(opts metav1.ListOptions) (*GroupMemberList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	DeleteCollection(deleteOpts *metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Controller() GroupMemberController
	AddHandler(ctx context.Context, name string, sync GroupMemberHandlerFunc)
	AddFeatureHandler(ctx context.Context, enabled func() bool, name string, sync GroupMemberHandlerFunc)
	AddLifecycle(ctx context.Context, name string, lifecycle GroupMemberLifecycle)
	AddFeatureLifecycle(ctx context.Context, enabled func() bool, name string, lifecycle GroupMemberLifecycle)
	AddClusterScopedHandler(ctx context.Context, name, clusterName string, sync GroupMemberHandlerFunc)
	AddClusterScopedFeatureHandler(ctx context.Context, enabled func() bool, name, clusterName string, sync GroupMemberHandlerFunc)
	AddClusterScopedLifecycle(ctx context.Context, name, clusterName string, lifecycle GroupMemberLifecycle)
	AddClusterScopedFeatureLifecycle(ctx context.Context, enabled func() bool, name, clusterName string, lifecycle GroupMemberLifecycle)
}

type groupMemberLister struct {
	controller *groupMemberController
}

func (l *groupMemberLister) List(namespace string, selector labels.Selector) (ret []*GroupMember, err error) {
	err = cache.ListAllByNamespace(l.controller.Informer().GetIndexer(), namespace, selector, func(obj interface{}) {
		ret = append(ret, obj.(*GroupMember))
	})
	return
}

func (l *groupMemberLister) Get(namespace, name string) (*GroupMember, error) {
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
			Group:    GroupMemberGroupVersionKind.Group,
			Resource: "groupMember",
		}, key)
	}
	return obj.(*GroupMember), nil
}

type groupMemberController struct {
	controller.GenericController
}

func (c *groupMemberController) Generic() controller.GenericController {
	return c.GenericController
}

func (c *groupMemberController) Lister() GroupMemberLister {
	return &groupMemberLister{
		controller: c,
	}
}

func (c *groupMemberController) AddHandler(ctx context.Context, name string, handler GroupMemberHandlerFunc) {
	c.GenericController.AddHandler(ctx, name, func(key string, obj interface{}) (interface{}, error) {
		if obj == nil {
			return handler(key, nil)
		} else if v, ok := obj.(*GroupMember); ok {
			return handler(key, v)
		} else {
			return nil, nil
		}
	})
}

func (c *groupMemberController) AddFeatureHandler(ctx context.Context, enabled func() bool, name string, handler GroupMemberHandlerFunc) {
	c.GenericController.AddHandler(ctx, name, func(key string, obj interface{}) (interface{}, error) {
		if !enabled() {
			return nil, nil
		} else if obj == nil {
			return handler(key, nil)
		} else if v, ok := obj.(*GroupMember); ok {
			return handler(key, v)
		} else {
			return nil, nil
		}
	})
}

func (c *groupMemberController) AddClusterScopedHandler(ctx context.Context, name, cluster string, handler GroupMemberHandlerFunc) {
	c.GenericController.AddHandler(ctx, name, func(key string, obj interface{}) (interface{}, error) {
		if obj == nil {
			return handler(key, nil)
		} else if v, ok := obj.(*GroupMember); ok && controller.ObjectInCluster(cluster, obj) {
			return handler(key, v)
		} else {
			return nil, nil
		}
	})
}

func (c *groupMemberController) AddClusterScopedFeatureHandler(ctx context.Context, enabled func() bool, name, cluster string, handler GroupMemberHandlerFunc) {
	c.GenericController.AddHandler(ctx, name, func(key string, obj interface{}) (interface{}, error) {
		if !enabled() {
			return nil, nil
		} else if obj == nil {
			return handler(key, nil)
		} else if v, ok := obj.(*GroupMember); ok && controller.ObjectInCluster(cluster, obj) {
			return handler(key, v)
		} else {
			return nil, nil
		}
	})
}

type groupMemberFactory struct {
}

func (c groupMemberFactory) Object() runtime.Object {
	return &GroupMember{}
}

func (c groupMemberFactory) List() runtime.Object {
	return &GroupMemberList{}
}

func (s *groupMemberClient) Controller() GroupMemberController {
	s.client.Lock()
	defer s.client.Unlock()

	c, ok := s.client.groupMemberControllers[s.ns]
	if ok {
		return c
	}

	genericController := controller.NewGenericController(GroupMemberGroupVersionKind.Kind+"Controller",
		s.objectClient)

	c = &groupMemberController{
		GenericController: genericController,
	}

	s.client.groupMemberControllers[s.ns] = c
	s.client.starters = append(s.client.starters, c)

	return c
}

type groupMemberClient struct {
	client       *Client
	ns           string
	objectClient *objectclient.ObjectClient
	controller   GroupMemberController
}

func (s *groupMemberClient) ObjectClient() *objectclient.ObjectClient {
	return s.objectClient
}

func (s *groupMemberClient) Create(o *GroupMember) (*GroupMember, error) {
	obj, err := s.objectClient.Create(o)
	return obj.(*GroupMember), err
}

func (s *groupMemberClient) Get(name string, opts metav1.GetOptions) (*GroupMember, error) {
	obj, err := s.objectClient.Get(name, opts)
	return obj.(*GroupMember), err
}

func (s *groupMemberClient) GetNamespaced(namespace, name string, opts metav1.GetOptions) (*GroupMember, error) {
	obj, err := s.objectClient.GetNamespaced(namespace, name, opts)
	return obj.(*GroupMember), err
}

func (s *groupMemberClient) Update(o *GroupMember) (*GroupMember, error) {
	obj, err := s.objectClient.Update(o.Name, o)
	return obj.(*GroupMember), err
}

func (s *groupMemberClient) Delete(name string, options *metav1.DeleteOptions) error {
	return s.objectClient.Delete(name, options)
}

func (s *groupMemberClient) DeleteNamespaced(namespace, name string, options *metav1.DeleteOptions) error {
	return s.objectClient.DeleteNamespaced(namespace, name, options)
}

func (s *groupMemberClient) List(opts metav1.ListOptions) (*GroupMemberList, error) {
	obj, err := s.objectClient.List(opts)
	return obj.(*GroupMemberList), err
}

func (s *groupMemberClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	return s.objectClient.Watch(opts)
}

// Patch applies the patch and returns the patched deployment.
func (s *groupMemberClient) Patch(o *GroupMember, patchType types.PatchType, data []byte, subresources ...string) (*GroupMember, error) {
	obj, err := s.objectClient.Patch(o.Name, o, patchType, data, subresources...)
	return obj.(*GroupMember), err
}

func (s *groupMemberClient) DeleteCollection(deleteOpts *metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	return s.objectClient.DeleteCollection(deleteOpts, listOpts)
}

func (s *groupMemberClient) AddHandler(ctx context.Context, name string, sync GroupMemberHandlerFunc) {
	s.Controller().AddHandler(ctx, name, sync)
}

func (s *groupMemberClient) AddFeatureHandler(ctx context.Context, enabled func() bool, name string, sync GroupMemberHandlerFunc) {
	s.Controller().AddFeatureHandler(ctx, enabled, name, sync)
}

func (s *groupMemberClient) AddLifecycle(ctx context.Context, name string, lifecycle GroupMemberLifecycle) {
	sync := NewGroupMemberLifecycleAdapter(name, false, s, lifecycle)
	s.Controller().AddHandler(ctx, name, sync)
}

func (s *groupMemberClient) AddFeatureLifecycle(ctx context.Context, enabled func() bool, name string, lifecycle GroupMemberLifecycle) {
	sync := NewGroupMemberLifecycleAdapter(name, false, s, lifecycle)
	s.Controller().AddFeatureHandler(ctx, enabled, name, sync)
}

func (s *groupMemberClient) AddClusterScopedHandler(ctx context.Context, name, clusterName string, sync GroupMemberHandlerFunc) {
	s.Controller().AddClusterScopedHandler(ctx, name, clusterName, sync)
}

func (s *groupMemberClient) AddClusterScopedFeatureHandler(ctx context.Context, enabled func() bool, name, clusterName string, sync GroupMemberHandlerFunc) {
	s.Controller().AddClusterScopedFeatureHandler(ctx, enabled, name, clusterName, sync)
}

func (s *groupMemberClient) AddClusterScopedLifecycle(ctx context.Context, name, clusterName string, lifecycle GroupMemberLifecycle) {
	sync := NewGroupMemberLifecycleAdapter(name+"_"+clusterName, true, s, lifecycle)
	s.Controller().AddClusterScopedHandler(ctx, name, clusterName, sync)
}

func (s *groupMemberClient) AddClusterScopedFeatureLifecycle(ctx context.Context, enabled func() bool, name, clusterName string, lifecycle GroupMemberLifecycle) {
	sync := NewGroupMemberLifecycleAdapter(name+"_"+clusterName, true, s, lifecycle)
	s.Controller().AddClusterScopedFeatureHandler(ctx, enabled, name, clusterName, sync)
}

type GroupMemberIndexer func(obj *GroupMember) ([]string, error)

type GroupMemberClientCache interface {
	Get(namespace, name string) (*GroupMember, error)
	List(namespace string, selector labels.Selector) ([]*GroupMember, error)

	Index(name string, indexer GroupMemberIndexer)
	GetIndexed(name, key string) ([]*GroupMember, error)
}

type GroupMemberClient interface {
	Create(*GroupMember) (*GroupMember, error)
	Get(namespace, name string, opts metav1.GetOptions) (*GroupMember, error)
	Update(*GroupMember) (*GroupMember, error)
	Delete(namespace, name string, options *metav1.DeleteOptions) error
	List(namespace string, opts metav1.ListOptions) (*GroupMemberList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)

	Cache() GroupMemberClientCache

	OnCreate(ctx context.Context, name string, sync GroupMemberChangeHandlerFunc)
	OnChange(ctx context.Context, name string, sync GroupMemberChangeHandlerFunc)
	OnRemove(ctx context.Context, name string, sync GroupMemberChangeHandlerFunc)
	Enqueue(namespace, name string)

	Generic() controller.GenericController
	ObjectClient() *objectclient.ObjectClient
	Interface() GroupMemberInterface
}

type groupMemberClientCache struct {
	client *groupMemberClient2
}

type groupMemberClient2 struct {
	iface      GroupMemberInterface
	controller GroupMemberController
}

func (n *groupMemberClient2) Interface() GroupMemberInterface {
	return n.iface
}

func (n *groupMemberClient2) Generic() controller.GenericController {
	return n.iface.Controller().Generic()
}

func (n *groupMemberClient2) ObjectClient() *objectclient.ObjectClient {
	return n.Interface().ObjectClient()
}

func (n *groupMemberClient2) Enqueue(namespace, name string) {
	n.iface.Controller().Enqueue(namespace, name)
}

func (n *groupMemberClient2) Create(obj *GroupMember) (*GroupMember, error) {
	return n.iface.Create(obj)
}

func (n *groupMemberClient2) Get(namespace, name string, opts metav1.GetOptions) (*GroupMember, error) {
	return n.iface.GetNamespaced(namespace, name, opts)
}

func (n *groupMemberClient2) Update(obj *GroupMember) (*GroupMember, error) {
	return n.iface.Update(obj)
}

func (n *groupMemberClient2) Delete(namespace, name string, options *metav1.DeleteOptions) error {
	return n.iface.DeleteNamespaced(namespace, name, options)
}

func (n *groupMemberClient2) List(namespace string, opts metav1.ListOptions) (*GroupMemberList, error) {
	return n.iface.List(opts)
}

func (n *groupMemberClient2) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	return n.iface.Watch(opts)
}

func (n *groupMemberClientCache) Get(namespace, name string) (*GroupMember, error) {
	return n.client.controller.Lister().Get(namespace, name)
}

func (n *groupMemberClientCache) List(namespace string, selector labels.Selector) ([]*GroupMember, error) {
	return n.client.controller.Lister().List(namespace, selector)
}

func (n *groupMemberClient2) Cache() GroupMemberClientCache {
	n.loadController()
	return &groupMemberClientCache{
		client: n,
	}
}

func (n *groupMemberClient2) OnCreate(ctx context.Context, name string, sync GroupMemberChangeHandlerFunc) {
	n.loadController()
	n.iface.AddLifecycle(ctx, name+"-create", &groupMemberLifecycleDelegate{create: sync})
}

func (n *groupMemberClient2) OnChange(ctx context.Context, name string, sync GroupMemberChangeHandlerFunc) {
	n.loadController()
	n.iface.AddLifecycle(ctx, name+"-change", &groupMemberLifecycleDelegate{update: sync})
}

func (n *groupMemberClient2) OnRemove(ctx context.Context, name string, sync GroupMemberChangeHandlerFunc) {
	n.loadController()
	n.iface.AddLifecycle(ctx, name, &groupMemberLifecycleDelegate{remove: sync})
}

func (n *groupMemberClientCache) Index(name string, indexer GroupMemberIndexer) {
	err := n.client.controller.Informer().GetIndexer().AddIndexers(map[string]cache.IndexFunc{
		name: func(obj interface{}) ([]string, error) {
			if v, ok := obj.(*GroupMember); ok {
				return indexer(v)
			}
			return nil, nil
		},
	})

	if err != nil {
		panic(err)
	}
}

func (n *groupMemberClientCache) GetIndexed(name, key string) ([]*GroupMember, error) {
	var result []*GroupMember
	objs, err := n.client.controller.Informer().GetIndexer().ByIndex(name, key)
	if err != nil {
		return nil, err
	}
	for _, obj := range objs {
		if v, ok := obj.(*GroupMember); ok {
			result = append(result, v)
		}
	}

	return result, nil
}

func (n *groupMemberClient2) loadController() {
	if n.controller == nil {
		n.controller = n.iface.Controller()
	}
}

type groupMemberLifecycleDelegate struct {
	create GroupMemberChangeHandlerFunc
	update GroupMemberChangeHandlerFunc
	remove GroupMemberChangeHandlerFunc
}

func (n *groupMemberLifecycleDelegate) HasCreate() bool {
	return n.create != nil
}

func (n *groupMemberLifecycleDelegate) Create(obj *GroupMember) (runtime.Object, error) {
	if n.create == nil {
		return obj, nil
	}
	return n.create(obj)
}

func (n *groupMemberLifecycleDelegate) HasFinalize() bool {
	return n.remove != nil
}

func (n *groupMemberLifecycleDelegate) Remove(obj *GroupMember) (runtime.Object, error) {
	if n.remove == nil {
		return obj, nil
	}
	return n.remove(obj)
}

func (n *groupMemberLifecycleDelegate) Updated(obj *GroupMember) (runtime.Object, error) {
	if n.update == nil {
		return obj, nil
	}
	return n.update(obj)
}
