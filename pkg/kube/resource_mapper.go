/*
Copyright 2025 The Ketches Authors.

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

package kube

import (
	"slices"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery/cached/memory"
	"k8s.io/client-go/restmapper"
)

func GetAllGroupResources() ([]string, error) {
	discoveryClient := DiscoveryClient()

	apiResourceLists, err := discoveryClient.ServerPreferredResources()
	if err != nil {
		return nil, err
	}

	var result []string
	for _, resourceList := range apiResourceLists {
		for _, res := range resourceList.APIResources {
			result = append(result, schema.GroupResource{
				Group:    res.Group,
				Resource: res.Name,
			}.String())
		}
	}
	return result, nil
}

// GetResourceNameFromGroupVersionKind returns the resource name from the given GroupVersionKind.
func GetResourceNameFromGroupVersionKind(gvk schema.GroupVersionKind) (string, error) {
	discoveryClient := DiscoveryClient()

	groupResources, err := restmapper.GetAPIGroupResources(discoveryClient)
	if err != nil {
		return "", err
	}
	mapper := restmapper.NewDiscoveryRESTMapper(groupResources)

	mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		return "", err
	}

	return mapping.Resource.Resource, nil
}

// GetGroupVersionResourceFromGroupVersionKind returns the GroupVersionResource from the given GroupVersionKind.
func GetGroupVersionResourceFromGroupVersionKind(gvk schema.GroupVersionKind) (schema.GroupVersionResource, error) {
	resource, err := GetResourceNameFromGroupVersionKind(gvk)
	if err != nil {
		return schema.GroupVersionResource{}, err
	}

	return schema.GroupVersionResource{
		Group:    gvk.Group,
		Version:  gvk.Version,
		Resource: resource,
	}, nil
}

// GetGroupVersionKindFromResourceName returns the GroupVersionKind from the given resource name.
// resource name can be plural, singular or short names.
func GetGroupVersionKindFromResourceName(resourceName string) ([]schema.GroupVersionKind, error) {
	discoveryClient := DiscoveryClient()

	apiResourceLists, err := discoveryClient.ServerPreferredResources()
	if err != nil {
		return nil, err
	}

	var result []schema.GroupVersionKind
	for _, resourceList := range apiResourceLists {
		gv, err := schema.ParseGroupVersion(resourceList.GroupVersion)
		if err != nil {
			continue
		}

		for _, res := range resourceList.APIResources {
			if res.Name == resourceName || res.SingularName == resourceName || slices.Contains(res.ShortNames, resourceName) {
				result = append(result, schema.GroupVersionKind{
					Group:   gv.Group,
					Version: gv.Version,
					Kind:    res.Kind,
				})
			}
		}
	}
	return result, nil
}

// GetGroupVersionResourceFromResourceName returns the GroupVersionResource from the given resource name.
// resource name can be plural, singular or short names.
func GetGroupVersionResourceFromResourceName(resourceName string) ([]schema.GroupVersionResource, error) {
	discoveryClient := DiscoveryClient()

	apiResourceLists, err := discoveryClient.ServerPreferredResources()
	if err != nil {
		return nil, err
	}

	var result []schema.GroupVersionResource
	for _, resourceList := range apiResourceLists {
		gv, err := schema.ParseGroupVersion(resourceList.GroupVersion)
		if err != nil {
			continue
		}

		gr := schema.ParseGroupResource(resourceName)
		for _, res := range resourceList.APIResources {

			if gr.Group != "" && res.Group != "" && res.Group != gr.Group {
				continue
			}
			if res.Name == gr.Resource || res.SingularName == gr.Resource || slices.Contains(res.ShortNames, gr.Resource) {
				result = append(result, schema.GroupVersionResource{
					Group:    gv.Group,
					Version:  gv.Version,
					Resource: res.Name,
				})
			}
		}
	}
	return result, nil
}

// GetPreferredGroupVersionResource returns the preferred GroupVersionResource for the given fuzz resource name.
// resource name can be plural, singular, short names or grouped resource name like deployments.apps, deployments.v1.apps.
func GetPreferredGroupVersionResource(fuzzResourceName string) (*schema.GroupVersionResource, error) {
	discoveryClient := DiscoveryClient()
	restMapper := restmapper.NewDeferredDiscoveryRESTMapper(memory.NewMemCacheClient(discoveryClient))

	gvr, err := restMapper.ResourceFor(schema.GroupVersionResource{
		Resource: fuzzResourceName,
	})
	if err != nil {
		return nil, err
	}

	return &gvr, nil
}
