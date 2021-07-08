/*
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
// Code generated by lister-gen. DO NOT EDIT.

package v1alpha2

import (
	v1alpha2 "github.com/openservicemesh/osm/pkg/apis/config/v1alpha2"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// MeshConfigLister helps list MeshConfigs.
// All objects returned here must be treated as read-only.
type MeshConfigLister interface {
	// List lists all MeshConfigs in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha2.MeshConfig, err error)
	// MeshConfigs returns an object that can list and get MeshConfigs.
	MeshConfigs(namespace string) MeshConfigNamespaceLister
	MeshConfigListerExpansion
}

// meshConfigLister implements the MeshConfigLister interface.
type meshConfigLister struct {
	indexer cache.Indexer
}

// NewMeshConfigLister returns a new MeshConfigLister.
func NewMeshConfigLister(indexer cache.Indexer) MeshConfigLister {
	return &meshConfigLister{indexer: indexer}
}

// List lists all MeshConfigs in the indexer.
func (s *meshConfigLister) List(selector labels.Selector) (ret []*v1alpha2.MeshConfig, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha2.MeshConfig))
	})
	return ret, err
}

// MeshConfigs returns an object that can list and get MeshConfigs.
func (s *meshConfigLister) MeshConfigs(namespace string) MeshConfigNamespaceLister {
	return meshConfigNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// MeshConfigNamespaceLister helps list and get MeshConfigs.
// All objects returned here must be treated as read-only.
type MeshConfigNamespaceLister interface {
	// List lists all MeshConfigs in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha2.MeshConfig, err error)
	// Get retrieves the MeshConfig from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1alpha2.MeshConfig, error)
	MeshConfigNamespaceListerExpansion
}

// meshConfigNamespaceLister implements the MeshConfigNamespaceLister
// interface.
type meshConfigNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all MeshConfigs in the indexer for a given namespace.
func (s meshConfigNamespaceLister) List(selector labels.Selector) (ret []*v1alpha2.MeshConfig, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha2.MeshConfig))
	})
	return ret, err
}

// Get retrieves the MeshConfig from the indexer for a given namespace and name.
func (s meshConfigNamespaceLister) Get(name string) (*v1alpha2.MeshConfig, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha2.Resource("meshconfig"), name)
	}
	return obj.(*v1alpha2.MeshConfig), nil
}
