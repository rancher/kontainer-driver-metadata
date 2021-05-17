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
	SettingGroupVersionKind = schema.GroupVersionKind{
		Version: Version,
		Group:   GroupName,
		Kind:    "Setting",
	}
	SettingResource = metav1.APIResource{
		Name:         "settings",
		SingularName: "setting",
		Namespaced:   false,
		Kind:         SettingGroupVersionKind.Kind,
	}

	SettingGroupVersionResource = schema.GroupVersionResource{
		Group:    GroupName,
		Version:  Version,
		Resource: "settings",
	}
)

func init() {
	resource.Put(SettingGroupVersionResource)
}

func NewSetting(namespace, name string, obj Setting) *Setting {
	obj.APIVersion, obj.Kind = SettingGroupVersionKind.ToAPIVersionAndKind()
	obj.Name = name
	obj.Namespace = namespace
	return &obj
}

type SettingList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Setting `json:"items"`
}

type SettingHandlerFunc func(key string, obj *Setting) (runtime.Object, error)

type SettingChangeHandlerFunc func(obj *Setting) (runtime.Object, error)

type SettingLister interface {
	List(namespace string, selector labels.Selector) (ret []*Setting, err error)
	Get(namespace, name string) (*Setting, error)
}

type SettingController interface {
	Generic() controller.GenericController
	Informer() cache.SharedIndexInformer
	Lister() SettingLister
	AddHandler(ctx context.Context, name string, handler SettingHandlerFunc)
	AddFeatureHandler(ctx context.Context, enabled func() bool, name string, sync SettingHandlerFunc)
	AddClusterScopedHandler(ctx context.Context, name, clusterName string, handler SettingHandlerFunc)
	AddClusterScopedFeatureHandler(ctx context.Context, enabled func() bool, name, clusterName string, handler SettingHandlerFunc)
	Enqueue(namespace, name string)
	Sync(ctx context.Context) error
	Start(ctx context.Context, threadiness int) error
}

type SettingInterface interface {
	ObjectClient() *objectclient.ObjectClient
	Create(*Setting) (*Setting, error)
	GetNamespaced(namespace, name string, opts metav1.GetOptions) (*Setting, error)
	Get(name string, opts metav1.GetOptions) (*Setting, error)
	Update(*Setting) (*Setting, error)
	Delete(name string, options *metav1.DeleteOptions) error
	DeleteNamespaced(namespace, name string, options *metav1.DeleteOptions) error
	List(opts metav1.ListOptions) (*SettingList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	DeleteCollection(deleteOpts *metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Controller() SettingController
	AddHandler(ctx context.Context, name string, sync SettingHandlerFunc)
	AddFeatureHandler(ctx context.Context, enabled func() bool, name string, sync SettingHandlerFunc)
	AddLifecycle(ctx context.Context, name string, lifecycle SettingLifecycle)
	AddFeatureLifecycle(ctx context.Context, enabled func() bool, name string, lifecycle SettingLifecycle)
	AddClusterScopedHandler(ctx context.Context, name, clusterName string, sync SettingHandlerFunc)
	AddClusterScopedFeatureHandler(ctx context.Context, enabled func() bool, name, clusterName string, sync SettingHandlerFunc)
	AddClusterScopedLifecycle(ctx context.Context, name, clusterName string, lifecycle SettingLifecycle)
	AddClusterScopedFeatureLifecycle(ctx context.Context, enabled func() bool, name, clusterName string, lifecycle SettingLifecycle)
}

type settingLister struct {
	controller *settingController
}

func (l *settingLister) List(namespace string, selector labels.Selector) (ret []*Setting, err error) {
	err = cache.ListAllByNamespace(l.controller.Informer().GetIndexer(), namespace, selector, func(obj interface{}) {
		ret = append(ret, obj.(*Setting))
	})
	return
}

func (l *settingLister) Get(namespace, name string) (*Setting, error) {
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
			Group:    SettingGroupVersionKind.Group,
			Resource: "setting",
		}, key)
	}
	return obj.(*Setting), nil
}

type settingController struct {
	controller.GenericController
}

func (c *settingController) Generic() controller.GenericController {
	return c.GenericController
}

func (c *settingController) Lister() SettingLister {
	return &settingLister{
		controller: c,
	}
}

func (c *settingController) AddHandler(ctx context.Context, name string, handler SettingHandlerFunc) {
	c.GenericController.AddHandler(ctx, name, func(key string, obj interface{}) (interface{}, error) {
		if obj == nil {
			return handler(key, nil)
		} else if v, ok := obj.(*Setting); ok {
			return handler(key, v)
		} else {
			return nil, nil
		}
	})
}

func (c *settingController) AddFeatureHandler(ctx context.Context, enabled func() bool, name string, handler SettingHandlerFunc) {
	c.GenericController.AddHandler(ctx, name, func(key string, obj interface{}) (interface{}, error) {
		if !enabled() {
			return nil, nil
		} else if obj == nil {
			return handler(key, nil)
		} else if v, ok := obj.(*Setting); ok {
			return handler(key, v)
		} else {
			return nil, nil
		}
	})
}

func (c *settingController) AddClusterScopedHandler(ctx context.Context, name, cluster string, handler SettingHandlerFunc) {
	c.GenericController.AddHandler(ctx, name, func(key string, obj interface{}) (interface{}, error) {
		if obj == nil {
			return handler(key, nil)
		} else if v, ok := obj.(*Setting); ok && controller.ObjectInCluster(cluster, obj) {
			return handler(key, v)
		} else {
			return nil, nil
		}
	})
}

func (c *settingController) AddClusterScopedFeatureHandler(ctx context.Context, enabled func() bool, name, cluster string, handler SettingHandlerFunc) {
	c.GenericController.AddHandler(ctx, name, func(key string, obj interface{}) (interface{}, error) {
		if !enabled() {
			return nil, nil
		} else if obj == nil {
			return handler(key, nil)
		} else if v, ok := obj.(*Setting); ok && controller.ObjectInCluster(cluster, obj) {
			return handler(key, v)
		} else {
			return nil, nil
		}
	})
}

type settingFactory struct {
}

func (c settingFactory) Object() runtime.Object {
	return &Setting{}
}

func (c settingFactory) List() runtime.Object {
	return &SettingList{}
}

func (s *settingClient) Controller() SettingController {
	s.client.Lock()
	defer s.client.Unlock()

	c, ok := s.client.settingControllers[s.ns]
	if ok {
		return c
	}

	genericController := controller.NewGenericController(SettingGroupVersionKind.Kind+"Controller",
		s.objectClient)

	c = &settingController{
		GenericController: genericController,
	}

	s.client.settingControllers[s.ns] = c
	s.client.starters = append(s.client.starters, c)

	return c
}

type settingClient struct {
	client       *Client
	ns           string
	objectClient *objectclient.ObjectClient
	controller   SettingController
}

func (s *settingClient) ObjectClient() *objectclient.ObjectClient {
	return s.objectClient
}

func (s *settingClient) Create(o *Setting) (*Setting, error) {
	obj, err := s.objectClient.Create(o)
	return obj.(*Setting), err
}

func (s *settingClient) Get(name string, opts metav1.GetOptions) (*Setting, error) {
	obj, err := s.objectClient.Get(name, opts)
	return obj.(*Setting), err
}

func (s *settingClient) GetNamespaced(namespace, name string, opts metav1.GetOptions) (*Setting, error) {
	obj, err := s.objectClient.GetNamespaced(namespace, name, opts)
	return obj.(*Setting), err
}

func (s *settingClient) Update(o *Setting) (*Setting, error) {
	obj, err := s.objectClient.Update(o.Name, o)
	return obj.(*Setting), err
}

func (s *settingClient) Delete(name string, options *metav1.DeleteOptions) error {
	return s.objectClient.Delete(name, options)
}

func (s *settingClient) DeleteNamespaced(namespace, name string, options *metav1.DeleteOptions) error {
	return s.objectClient.DeleteNamespaced(namespace, name, options)
}

func (s *settingClient) List(opts metav1.ListOptions) (*SettingList, error) {
	obj, err := s.objectClient.List(opts)
	return obj.(*SettingList), err
}

func (s *settingClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	return s.objectClient.Watch(opts)
}

// Patch applies the patch and returns the patched deployment.
func (s *settingClient) Patch(o *Setting, patchType types.PatchType, data []byte, subresources ...string) (*Setting, error) {
	obj, err := s.objectClient.Patch(o.Name, o, patchType, data, subresources...)
	return obj.(*Setting), err
}

func (s *settingClient) DeleteCollection(deleteOpts *metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	return s.objectClient.DeleteCollection(deleteOpts, listOpts)
}

func (s *settingClient) AddHandler(ctx context.Context, name string, sync SettingHandlerFunc) {
	s.Controller().AddHandler(ctx, name, sync)
}

func (s *settingClient) AddFeatureHandler(ctx context.Context, enabled func() bool, name string, sync SettingHandlerFunc) {
	s.Controller().AddFeatureHandler(ctx, enabled, name, sync)
}

func (s *settingClient) AddLifecycle(ctx context.Context, name string, lifecycle SettingLifecycle) {
	sync := NewSettingLifecycleAdapter(name, false, s, lifecycle)
	s.Controller().AddHandler(ctx, name, sync)
}

func (s *settingClient) AddFeatureLifecycle(ctx context.Context, enabled func() bool, name string, lifecycle SettingLifecycle) {
	sync := NewSettingLifecycleAdapter(name, false, s, lifecycle)
	s.Controller().AddFeatureHandler(ctx, enabled, name, sync)
}

func (s *settingClient) AddClusterScopedHandler(ctx context.Context, name, clusterName string, sync SettingHandlerFunc) {
	s.Controller().AddClusterScopedHandler(ctx, name, clusterName, sync)
}

func (s *settingClient) AddClusterScopedFeatureHandler(ctx context.Context, enabled func() bool, name, clusterName string, sync SettingHandlerFunc) {
	s.Controller().AddClusterScopedFeatureHandler(ctx, enabled, name, clusterName, sync)
}

func (s *settingClient) AddClusterScopedLifecycle(ctx context.Context, name, clusterName string, lifecycle SettingLifecycle) {
	sync := NewSettingLifecycleAdapter(name+"_"+clusterName, true, s, lifecycle)
	s.Controller().AddClusterScopedHandler(ctx, name, clusterName, sync)
}

func (s *settingClient) AddClusterScopedFeatureLifecycle(ctx context.Context, enabled func() bool, name, clusterName string, lifecycle SettingLifecycle) {
	sync := NewSettingLifecycleAdapter(name+"_"+clusterName, true, s, lifecycle)
	s.Controller().AddClusterScopedFeatureHandler(ctx, enabled, name, clusterName, sync)
}

type SettingIndexer func(obj *Setting) ([]string, error)

type SettingClientCache interface {
	Get(namespace, name string) (*Setting, error)
	List(namespace string, selector labels.Selector) ([]*Setting, error)

	Index(name string, indexer SettingIndexer)
	GetIndexed(name, key string) ([]*Setting, error)
}

type SettingClient interface {
	Create(*Setting) (*Setting, error)
	Get(namespace, name string, opts metav1.GetOptions) (*Setting, error)
	Update(*Setting) (*Setting, error)
	Delete(namespace, name string, options *metav1.DeleteOptions) error
	List(namespace string, opts metav1.ListOptions) (*SettingList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)

	Cache() SettingClientCache

	OnCreate(ctx context.Context, name string, sync SettingChangeHandlerFunc)
	OnChange(ctx context.Context, name string, sync SettingChangeHandlerFunc)
	OnRemove(ctx context.Context, name string, sync SettingChangeHandlerFunc)
	Enqueue(namespace, name string)

	Generic() controller.GenericController
	ObjectClient() *objectclient.ObjectClient
	Interface() SettingInterface
}

type settingClientCache struct {
	client *settingClient2
}

type settingClient2 struct {
	iface      SettingInterface
	controller SettingController
}

func (n *settingClient2) Interface() SettingInterface {
	return n.iface
}

func (n *settingClient2) Generic() controller.GenericController {
	return n.iface.Controller().Generic()
}

func (n *settingClient2) ObjectClient() *objectclient.ObjectClient {
	return n.Interface().ObjectClient()
}

func (n *settingClient2) Enqueue(namespace, name string) {
	n.iface.Controller().Enqueue(namespace, name)
}

func (n *settingClient2) Create(obj *Setting) (*Setting, error) {
	return n.iface.Create(obj)
}

func (n *settingClient2) Get(namespace, name string, opts metav1.GetOptions) (*Setting, error) {
	return n.iface.GetNamespaced(namespace, name, opts)
}

func (n *settingClient2) Update(obj *Setting) (*Setting, error) {
	return n.iface.Update(obj)
}

func (n *settingClient2) Delete(namespace, name string, options *metav1.DeleteOptions) error {
	return n.iface.DeleteNamespaced(namespace, name, options)
}

func (n *settingClient2) List(namespace string, opts metav1.ListOptions) (*SettingList, error) {
	return n.iface.List(opts)
}

func (n *settingClient2) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	return n.iface.Watch(opts)
}

func (n *settingClientCache) Get(namespace, name string) (*Setting, error) {
	return n.client.controller.Lister().Get(namespace, name)
}

func (n *settingClientCache) List(namespace string, selector labels.Selector) ([]*Setting, error) {
	return n.client.controller.Lister().List(namespace, selector)
}

func (n *settingClient2) Cache() SettingClientCache {
	n.loadController()
	return &settingClientCache{
		client: n,
	}
}

func (n *settingClient2) OnCreate(ctx context.Context, name string, sync SettingChangeHandlerFunc) {
	n.loadController()
	n.iface.AddLifecycle(ctx, name+"-create", &settingLifecycleDelegate{create: sync})
}

func (n *settingClient2) OnChange(ctx context.Context, name string, sync SettingChangeHandlerFunc) {
	n.loadController()
	n.iface.AddLifecycle(ctx, name+"-change", &settingLifecycleDelegate{update: sync})
}

func (n *settingClient2) OnRemove(ctx context.Context, name string, sync SettingChangeHandlerFunc) {
	n.loadController()
	n.iface.AddLifecycle(ctx, name, &settingLifecycleDelegate{remove: sync})
}

func (n *settingClientCache) Index(name string, indexer SettingIndexer) {
	err := n.client.controller.Informer().GetIndexer().AddIndexers(map[string]cache.IndexFunc{
		name: func(obj interface{}) ([]string, error) {
			if v, ok := obj.(*Setting); ok {
				return indexer(v)
			}
			return nil, nil
		},
	})

	if err != nil {
		panic(err)
	}
}

func (n *settingClientCache) GetIndexed(name, key string) ([]*Setting, error) {
	var result []*Setting
	objs, err := n.client.controller.Informer().GetIndexer().ByIndex(name, key)
	if err != nil {
		return nil, err
	}
	for _, obj := range objs {
		if v, ok := obj.(*Setting); ok {
			result = append(result, v)
		}
	}

	return result, nil
}

func (n *settingClient2) loadController() {
	if n.controller == nil {
		n.controller = n.iface.Controller()
	}
}

type settingLifecycleDelegate struct {
	create SettingChangeHandlerFunc
	update SettingChangeHandlerFunc
	remove SettingChangeHandlerFunc
}

func (n *settingLifecycleDelegate) HasCreate() bool {
	return n.create != nil
}

func (n *settingLifecycleDelegate) Create(obj *Setting) (runtime.Object, error) {
	if n.create == nil {
		return obj, nil
	}
	return n.create(obj)
}

func (n *settingLifecycleDelegate) HasFinalize() bool {
	return n.remove != nil
}

func (n *settingLifecycleDelegate) Remove(obj *Setting) (runtime.Object, error) {
	if n.remove == nil {
		return obj, nil
	}
	return n.remove(obj)
}

func (n *settingLifecycleDelegate) Updated(obj *Setting) (runtime.Object, error) {
	if n.update == nil {
		return obj, nil
	}
	return n.update(obj)
}
