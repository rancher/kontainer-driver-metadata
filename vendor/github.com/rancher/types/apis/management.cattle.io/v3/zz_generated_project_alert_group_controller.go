package v3

import (
	"context"
	"time"

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
	ProjectAlertGroupGroupVersionKind = schema.GroupVersionKind{
		Version: Version,
		Group:   GroupName,
		Kind:    "ProjectAlertGroup",
	}
	ProjectAlertGroupResource = metav1.APIResource{
		Name:         "projectalertgroups",
		SingularName: "projectalertgroup",
		Namespaced:   true,

		Kind: ProjectAlertGroupGroupVersionKind.Kind,
	}

	ProjectAlertGroupGroupVersionResource = schema.GroupVersionResource{
		Group:    GroupName,
		Version:  Version,
		Resource: "projectalertgroups",
	}
)

func init() {
	resource.Put(ProjectAlertGroupGroupVersionResource)
}

func NewProjectAlertGroup(namespace, name string, obj ProjectAlertGroup) *ProjectAlertGroup {
	obj.APIVersion, obj.Kind = ProjectAlertGroupGroupVersionKind.ToAPIVersionAndKind()
	obj.Name = name
	obj.Namespace = namespace
	return &obj
}

type ProjectAlertGroupList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ProjectAlertGroup `json:"items"`
}

type ProjectAlertGroupHandlerFunc func(key string, obj *ProjectAlertGroup) (runtime.Object, error)

type ProjectAlertGroupChangeHandlerFunc func(obj *ProjectAlertGroup) (runtime.Object, error)

type ProjectAlertGroupLister interface {
	List(namespace string, selector labels.Selector) (ret []*ProjectAlertGroup, err error)
	Get(namespace, name string) (*ProjectAlertGroup, error)
}

type ProjectAlertGroupController interface {
	Generic() controller.GenericController
	Informer() cache.SharedIndexInformer
	Lister() ProjectAlertGroupLister
	AddHandler(ctx context.Context, name string, handler ProjectAlertGroupHandlerFunc)
	AddFeatureHandler(ctx context.Context, enabled func() bool, name string, sync ProjectAlertGroupHandlerFunc)
	AddClusterScopedHandler(ctx context.Context, name, clusterName string, handler ProjectAlertGroupHandlerFunc)
	AddClusterScopedFeatureHandler(ctx context.Context, enabled func() bool, name, clusterName string, handler ProjectAlertGroupHandlerFunc)
	Enqueue(namespace, name string)
	EnqueueAfter(namespace, name string, after time.Duration)
	Sync(ctx context.Context) error
	Start(ctx context.Context, threadiness int) error
}

type ProjectAlertGroupInterface interface {
	ObjectClient() *objectclient.ObjectClient
	Create(*ProjectAlertGroup) (*ProjectAlertGroup, error)
	GetNamespaced(namespace, name string, opts metav1.GetOptions) (*ProjectAlertGroup, error)
	Get(name string, opts metav1.GetOptions) (*ProjectAlertGroup, error)
	Update(*ProjectAlertGroup) (*ProjectAlertGroup, error)
	Delete(name string, options *metav1.DeleteOptions) error
	DeleteNamespaced(namespace, name string, options *metav1.DeleteOptions) error
	List(opts metav1.ListOptions) (*ProjectAlertGroupList, error)
	ListNamespaced(namespace string, opts metav1.ListOptions) (*ProjectAlertGroupList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	DeleteCollection(deleteOpts *metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Controller() ProjectAlertGroupController
	AddHandler(ctx context.Context, name string, sync ProjectAlertGroupHandlerFunc)
	AddFeatureHandler(ctx context.Context, enabled func() bool, name string, sync ProjectAlertGroupHandlerFunc)
	AddLifecycle(ctx context.Context, name string, lifecycle ProjectAlertGroupLifecycle)
	AddFeatureLifecycle(ctx context.Context, enabled func() bool, name string, lifecycle ProjectAlertGroupLifecycle)
	AddClusterScopedHandler(ctx context.Context, name, clusterName string, sync ProjectAlertGroupHandlerFunc)
	AddClusterScopedFeatureHandler(ctx context.Context, enabled func() bool, name, clusterName string, sync ProjectAlertGroupHandlerFunc)
	AddClusterScopedLifecycle(ctx context.Context, name, clusterName string, lifecycle ProjectAlertGroupLifecycle)
	AddClusterScopedFeatureLifecycle(ctx context.Context, enabled func() bool, name, clusterName string, lifecycle ProjectAlertGroupLifecycle)
}

type projectAlertGroupLister struct {
	controller *projectAlertGroupController
}

func (l *projectAlertGroupLister) List(namespace string, selector labels.Selector) (ret []*ProjectAlertGroup, err error) {
	err = cache.ListAllByNamespace(l.controller.Informer().GetIndexer(), namespace, selector, func(obj interface{}) {
		ret = append(ret, obj.(*ProjectAlertGroup))
	})
	return
}

func (l *projectAlertGroupLister) Get(namespace, name string) (*ProjectAlertGroup, error) {
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
			Group:    ProjectAlertGroupGroupVersionKind.Group,
			Resource: "projectAlertGroup",
		}, key)
	}
	return obj.(*ProjectAlertGroup), nil
}

type projectAlertGroupController struct {
	controller.GenericController
}

func (c *projectAlertGroupController) Generic() controller.GenericController {
	return c.GenericController
}

func (c *projectAlertGroupController) Lister() ProjectAlertGroupLister {
	return &projectAlertGroupLister{
		controller: c,
	}
}

func (c *projectAlertGroupController) AddHandler(ctx context.Context, name string, handler ProjectAlertGroupHandlerFunc) {
	c.GenericController.AddHandler(ctx, name, func(key string, obj interface{}) (interface{}, error) {
		if obj == nil {
			return handler(key, nil)
		} else if v, ok := obj.(*ProjectAlertGroup); ok {
			return handler(key, v)
		} else {
			return nil, nil
		}
	})
}

func (c *projectAlertGroupController) AddFeatureHandler(ctx context.Context, enabled func() bool, name string, handler ProjectAlertGroupHandlerFunc) {
	c.GenericController.AddHandler(ctx, name, func(key string, obj interface{}) (interface{}, error) {
		if !enabled() {
			return nil, nil
		} else if obj == nil {
			return handler(key, nil)
		} else if v, ok := obj.(*ProjectAlertGroup); ok {
			return handler(key, v)
		} else {
			return nil, nil
		}
	})
}

func (c *projectAlertGroupController) AddClusterScopedHandler(ctx context.Context, name, cluster string, handler ProjectAlertGroupHandlerFunc) {
	c.GenericController.AddHandler(ctx, name, func(key string, obj interface{}) (interface{}, error) {
		if obj == nil {
			return handler(key, nil)
		} else if v, ok := obj.(*ProjectAlertGroup); ok && controller.ObjectInCluster(cluster, obj) {
			return handler(key, v)
		} else {
			return nil, nil
		}
	})
}

func (c *projectAlertGroupController) AddClusterScopedFeatureHandler(ctx context.Context, enabled func() bool, name, cluster string, handler ProjectAlertGroupHandlerFunc) {
	c.GenericController.AddHandler(ctx, name, func(key string, obj interface{}) (interface{}, error) {
		if !enabled() {
			return nil, nil
		} else if obj == nil {
			return handler(key, nil)
		} else if v, ok := obj.(*ProjectAlertGroup); ok && controller.ObjectInCluster(cluster, obj) {
			return handler(key, v)
		} else {
			return nil, nil
		}
	})
}

type projectAlertGroupFactory struct {
}

func (c projectAlertGroupFactory) Object() runtime.Object {
	return &ProjectAlertGroup{}
}

func (c projectAlertGroupFactory) List() runtime.Object {
	return &ProjectAlertGroupList{}
}

func (s *projectAlertGroupClient) Controller() ProjectAlertGroupController {
	s.client.Lock()
	defer s.client.Unlock()

	c, ok := s.client.projectAlertGroupControllers[s.ns]
	if ok {
		return c
	}

	genericController := controller.NewGenericController(ProjectAlertGroupGroupVersionKind.Kind+"Controller",
		s.objectClient)

	c = &projectAlertGroupController{
		GenericController: genericController,
	}

	s.client.projectAlertGroupControllers[s.ns] = c
	s.client.starters = append(s.client.starters, c)

	return c
}

type projectAlertGroupClient struct {
	client       *Client
	ns           string
	objectClient *objectclient.ObjectClient
	controller   ProjectAlertGroupController
}

func (s *projectAlertGroupClient) ObjectClient() *objectclient.ObjectClient {
	return s.objectClient
}

func (s *projectAlertGroupClient) Create(o *ProjectAlertGroup) (*ProjectAlertGroup, error) {
	obj, err := s.objectClient.Create(o)
	return obj.(*ProjectAlertGroup), err
}

func (s *projectAlertGroupClient) Get(name string, opts metav1.GetOptions) (*ProjectAlertGroup, error) {
	obj, err := s.objectClient.Get(name, opts)
	return obj.(*ProjectAlertGroup), err
}

func (s *projectAlertGroupClient) GetNamespaced(namespace, name string, opts metav1.GetOptions) (*ProjectAlertGroup, error) {
	obj, err := s.objectClient.GetNamespaced(namespace, name, opts)
	return obj.(*ProjectAlertGroup), err
}

func (s *projectAlertGroupClient) Update(o *ProjectAlertGroup) (*ProjectAlertGroup, error) {
	obj, err := s.objectClient.Update(o.Name, o)
	return obj.(*ProjectAlertGroup), err
}

func (s *projectAlertGroupClient) Delete(name string, options *metav1.DeleteOptions) error {
	return s.objectClient.Delete(name, options)
}

func (s *projectAlertGroupClient) DeleteNamespaced(namespace, name string, options *metav1.DeleteOptions) error {
	return s.objectClient.DeleteNamespaced(namespace, name, options)
}

func (s *projectAlertGroupClient) List(opts metav1.ListOptions) (*ProjectAlertGroupList, error) {
	obj, err := s.objectClient.List(opts)
	return obj.(*ProjectAlertGroupList), err
}

func (s *projectAlertGroupClient) ListNamespaced(namespace string, opts metav1.ListOptions) (*ProjectAlertGroupList, error) {
	obj, err := s.objectClient.ListNamespaced(namespace, opts)
	return obj.(*ProjectAlertGroupList), err
}

func (s *projectAlertGroupClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	return s.objectClient.Watch(opts)
}

// Patch applies the patch and returns the patched deployment.
func (s *projectAlertGroupClient) Patch(o *ProjectAlertGroup, patchType types.PatchType, data []byte, subresources ...string) (*ProjectAlertGroup, error) {
	obj, err := s.objectClient.Patch(o.Name, o, patchType, data, subresources...)
	return obj.(*ProjectAlertGroup), err
}

func (s *projectAlertGroupClient) DeleteCollection(deleteOpts *metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	return s.objectClient.DeleteCollection(deleteOpts, listOpts)
}

func (s *projectAlertGroupClient) AddHandler(ctx context.Context, name string, sync ProjectAlertGroupHandlerFunc) {
	s.Controller().AddHandler(ctx, name, sync)
}

func (s *projectAlertGroupClient) AddFeatureHandler(ctx context.Context, enabled func() bool, name string, sync ProjectAlertGroupHandlerFunc) {
	s.Controller().AddFeatureHandler(ctx, enabled, name, sync)
}

func (s *projectAlertGroupClient) AddLifecycle(ctx context.Context, name string, lifecycle ProjectAlertGroupLifecycle) {
	sync := NewProjectAlertGroupLifecycleAdapter(name, false, s, lifecycle)
	s.Controller().AddHandler(ctx, name, sync)
}

func (s *projectAlertGroupClient) AddFeatureLifecycle(ctx context.Context, enabled func() bool, name string, lifecycle ProjectAlertGroupLifecycle) {
	sync := NewProjectAlertGroupLifecycleAdapter(name, false, s, lifecycle)
	s.Controller().AddFeatureHandler(ctx, enabled, name, sync)
}

func (s *projectAlertGroupClient) AddClusterScopedHandler(ctx context.Context, name, clusterName string, sync ProjectAlertGroupHandlerFunc) {
	s.Controller().AddClusterScopedHandler(ctx, name, clusterName, sync)
}

func (s *projectAlertGroupClient) AddClusterScopedFeatureHandler(ctx context.Context, enabled func() bool, name, clusterName string, sync ProjectAlertGroupHandlerFunc) {
	s.Controller().AddClusterScopedFeatureHandler(ctx, enabled, name, clusterName, sync)
}

func (s *projectAlertGroupClient) AddClusterScopedLifecycle(ctx context.Context, name, clusterName string, lifecycle ProjectAlertGroupLifecycle) {
	sync := NewProjectAlertGroupLifecycleAdapter(name+"_"+clusterName, true, s, lifecycle)
	s.Controller().AddClusterScopedHandler(ctx, name, clusterName, sync)
}

func (s *projectAlertGroupClient) AddClusterScopedFeatureLifecycle(ctx context.Context, enabled func() bool, name, clusterName string, lifecycle ProjectAlertGroupLifecycle) {
	sync := NewProjectAlertGroupLifecycleAdapter(name+"_"+clusterName, true, s, lifecycle)
	s.Controller().AddClusterScopedFeatureHandler(ctx, enabled, name, clusterName, sync)
}
