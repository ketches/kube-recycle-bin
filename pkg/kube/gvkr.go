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
	"fmt"
	"slices"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/restmapper"
)

// GetAllGroupResources returns all group resources in the cluster.
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

// GetPreferredGroupVersionResourceFor returns the preferred GroupVersionResource from the given resource name.
// resource name can be plural, singular, short names or grouped resource name like deployments.apps
func GetPreferredGroupVersionResourceFor(resource string) (*schema.GroupVersionResource, error) {
	gvr, gr := schema.ParseResourceArg(resource)
	if gvr == nil {
		gvr = &schema.GroupVersionResource{
			Resource: gr.Resource,
			Group:    gr.Group,
		}
	}

	discoveryClient := DiscoveryClient()

	apiResourceLists, err := discoveryClient.ServerPreferredResources()
	if err != nil {
		return nil, err
	}

	for _, resourceList := range apiResourceLists {
		gv, err := schema.ParseGroupVersion(resourceList.GroupVersion)
		if err != nil {
			continue
		}

		if gr.Group != "" && gv.Group != gr.Group && gvr.Group != "" && gv.Group != gvr.Group {
			continue
		}

		for _, res := range resourceList.APIResources {
			if res.Name == gr.Resource || res.SingularName == gr.Resource || slices.Contains(res.ShortNames, gr.Resource) {
				return &schema.GroupVersionResource{
					Group:    gv.Group,
					Version:  gv.Version,
					Resource: res.Name,
				}, nil
			}
		}
	}
	return nil, fmt.Errorf("can not find preferred GroupVersionResource for resource %s", resource)
}

// IsResourceNamespaced checks if the given GroupVersionResource is namespaced.
func IsResourceNamespaced(gvr schema.GroupVersionResource) (bool, error) {
	discoveryClient := DiscoveryClient()

	apiResourceList, err := discoveryClient.ServerResourcesForGroupVersion(gvr.GroupVersion().String())
	if err != nil {
		return false, err
	}

	for _, apiResource := range apiResourceList.APIResources {
		if apiResource.Name == gvr.Resource {
			return apiResource.Namespaced, nil
		}
	}
	return false, fmt.Errorf("can not assert if resource %s is namespaced", gvr.GroupResource().String())
}
