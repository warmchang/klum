/*
Copyright 2020 Rancher Labs, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by main. DO NOT EDIT.

package v1alpha1

import (
	"context"
	"time"

	v1alpha1 "github.com/ibuildthecloud/klum/pkg/apis/klum.cattle.io/v1alpha1"
	clientset "github.com/ibuildthecloud/klum/pkg/generated/clientset/versioned/typed/klum.cattle.io/v1alpha1"
	informers "github.com/ibuildthecloud/klum/pkg/generated/informers/externalversions/klum.cattle.io/v1alpha1"
	listers "github.com/ibuildthecloud/klum/pkg/generated/listers/klum.cattle.io/v1alpha1"
	"github.com/rancher/wrangler/pkg/generic"
	"k8s.io/apimachinery/pkg/api/equality"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
)

type KubeconfigHandler func(string, *v1alpha1.Kubeconfig) (*v1alpha1.Kubeconfig, error)

type KubeconfigController interface {
	generic.ControllerMeta
	KubeconfigClient

	OnChange(ctx context.Context, name string, sync KubeconfigHandler)
	OnRemove(ctx context.Context, name string, sync KubeconfigHandler)
	Enqueue(name string)
	EnqueueAfter(name string, duration time.Duration)

	Cache() KubeconfigCache
}

type KubeconfigClient interface {
	Create(*v1alpha1.Kubeconfig) (*v1alpha1.Kubeconfig, error)
	Update(*v1alpha1.Kubeconfig) (*v1alpha1.Kubeconfig, error)

	Delete(name string, options *metav1.DeleteOptions) error
	Get(name string, options metav1.GetOptions) (*v1alpha1.Kubeconfig, error)
	List(opts metav1.ListOptions) (*v1alpha1.KubeconfigList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.Kubeconfig, err error)
}

type KubeconfigCache interface {
	Get(name string) (*v1alpha1.Kubeconfig, error)
	List(selector labels.Selector) ([]*v1alpha1.Kubeconfig, error)

	AddIndexer(indexName string, indexer KubeconfigIndexer)
	GetByIndex(indexName, key string) ([]*v1alpha1.Kubeconfig, error)
}

type KubeconfigIndexer func(obj *v1alpha1.Kubeconfig) ([]string, error)

type kubeconfigController struct {
	controllerManager *generic.ControllerManager
	clientGetter      clientset.KubeconfigsGetter
	informer          informers.KubeconfigInformer
	gvk               schema.GroupVersionKind
}

func NewKubeconfigController(gvk schema.GroupVersionKind, controllerManager *generic.ControllerManager, clientGetter clientset.KubeconfigsGetter, informer informers.KubeconfigInformer) KubeconfigController {
	return &kubeconfigController{
		controllerManager: controllerManager,
		clientGetter:      clientGetter,
		informer:          informer,
		gvk:               gvk,
	}
}

func FromKubeconfigHandlerToHandler(sync KubeconfigHandler) generic.Handler {
	return func(key string, obj runtime.Object) (ret runtime.Object, err error) {
		var v *v1alpha1.Kubeconfig
		if obj == nil {
			v, err = sync(key, nil)
		} else {
			v, err = sync(key, obj.(*v1alpha1.Kubeconfig))
		}
		if v == nil {
			return nil, err
		}
		return v, err
	}
}

func (c *kubeconfigController) Updater() generic.Updater {
	return func(obj runtime.Object) (runtime.Object, error) {
		newObj, err := c.Update(obj.(*v1alpha1.Kubeconfig))
		if newObj == nil {
			return nil, err
		}
		return newObj, err
	}
}

func UpdateKubeconfigDeepCopyOnChange(client KubeconfigClient, obj *v1alpha1.Kubeconfig, handler func(obj *v1alpha1.Kubeconfig) (*v1alpha1.Kubeconfig, error)) (*v1alpha1.Kubeconfig, error) {
	if obj == nil {
		return obj, nil
	}

	copyObj := obj.DeepCopy()
	newObj, err := handler(copyObj)
	if newObj != nil {
		copyObj = newObj
	}
	if obj.ResourceVersion == copyObj.ResourceVersion && !equality.Semantic.DeepEqual(obj, copyObj) {
		return client.Update(copyObj)
	}

	return copyObj, err
}

func (c *kubeconfigController) AddGenericHandler(ctx context.Context, name string, handler generic.Handler) {
	c.controllerManager.AddHandler(ctx, c.gvk, c.informer.Informer(), name, handler)
}

func (c *kubeconfigController) AddGenericRemoveHandler(ctx context.Context, name string, handler generic.Handler) {
	removeHandler := generic.NewRemoveHandler(name, c.Updater(), handler)
	c.controllerManager.AddHandler(ctx, c.gvk, c.informer.Informer(), name, removeHandler)
}

func (c *kubeconfigController) OnChange(ctx context.Context, name string, sync KubeconfigHandler) {
	c.AddGenericHandler(ctx, name, FromKubeconfigHandlerToHandler(sync))
}

func (c *kubeconfigController) OnRemove(ctx context.Context, name string, sync KubeconfigHandler) {
	removeHandler := generic.NewRemoveHandler(name, c.Updater(), FromKubeconfigHandlerToHandler(sync))
	c.AddGenericHandler(ctx, name, removeHandler)
}

func (c *kubeconfigController) Enqueue(name string) {
	c.controllerManager.Enqueue(c.gvk, c.informer.Informer(), "", name)
}

func (c *kubeconfigController) EnqueueAfter(name string, duration time.Duration) {
	c.controllerManager.EnqueueAfter(c.gvk, c.informer.Informer(), "", name, duration)
}

func (c *kubeconfigController) Informer() cache.SharedIndexInformer {
	return c.informer.Informer()
}

func (c *kubeconfigController) GroupVersionKind() schema.GroupVersionKind {
	return c.gvk
}

func (c *kubeconfigController) Cache() KubeconfigCache {
	return &kubeconfigCache{
		lister:  c.informer.Lister(),
		indexer: c.informer.Informer().GetIndexer(),
	}
}

func (c *kubeconfigController) Create(obj *v1alpha1.Kubeconfig) (*v1alpha1.Kubeconfig, error) {
	return c.clientGetter.Kubeconfigs().Create(obj)
}

func (c *kubeconfigController) Update(obj *v1alpha1.Kubeconfig) (*v1alpha1.Kubeconfig, error) {
	return c.clientGetter.Kubeconfigs().Update(obj)
}

func (c *kubeconfigController) Delete(name string, options *metav1.DeleteOptions) error {
	return c.clientGetter.Kubeconfigs().Delete(name, options)
}

func (c *kubeconfigController) Get(name string, options metav1.GetOptions) (*v1alpha1.Kubeconfig, error) {
	return c.clientGetter.Kubeconfigs().Get(name, options)
}

func (c *kubeconfigController) List(opts metav1.ListOptions) (*v1alpha1.KubeconfigList, error) {
	return c.clientGetter.Kubeconfigs().List(opts)
}

func (c *kubeconfigController) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	return c.clientGetter.Kubeconfigs().Watch(opts)
}

func (c *kubeconfigController) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.Kubeconfig, err error) {
	return c.clientGetter.Kubeconfigs().Patch(name, pt, data, subresources...)
}

type kubeconfigCache struct {
	lister  listers.KubeconfigLister
	indexer cache.Indexer
}

func (c *kubeconfigCache) Get(name string) (*v1alpha1.Kubeconfig, error) {
	return c.lister.Get(name)
}

func (c *kubeconfigCache) List(selector labels.Selector) ([]*v1alpha1.Kubeconfig, error) {
	return c.lister.List(selector)
}

func (c *kubeconfigCache) AddIndexer(indexName string, indexer KubeconfigIndexer) {
	utilruntime.Must(c.indexer.AddIndexers(map[string]cache.IndexFunc{
		indexName: func(obj interface{}) (strings []string, e error) {
			return indexer(obj.(*v1alpha1.Kubeconfig))
		},
	}))
}

func (c *kubeconfigCache) GetByIndex(indexName, key string) (result []*v1alpha1.Kubeconfig, err error) {
	objs, err := c.indexer.ByIndex(indexName, key)
	if err != nil {
		return nil, err
	}
	for _, obj := range objs {
		result = append(result, obj.(*v1alpha1.Kubeconfig))
	}
	return result, nil
}
