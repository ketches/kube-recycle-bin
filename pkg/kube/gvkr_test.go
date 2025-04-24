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
	"testing"

	"k8s.io/apimachinery/pkg/runtime/schema"
)

func TestGetResourceNameFromGroupVersionKind(t *testing.T) {
	testdata := []struct {
		name    string
		gvk     schema.GroupVersionKind
		desired string
	}{
		{
			name: "gvk-deployment",
			gvk: schema.GroupVersionKind{
				Group:   "apps",
				Version: "v1",
				Kind:    "Deployment",
			},
			desired: "deployments",
		},
		{
			name: "gvk-pod",
			gvk: schema.GroupVersionKind{
				Group:   "",
				Version: "v1",
				Kind:    "Pod",
			},
			desired: "pods",
		},
	}

	for _, tt := range testdata {
		t.Run(tt.name, func(t *testing.T) {
			resourceName, err := GetResourceNameFromGroupVersionKind(tt.gvk)
			if err != nil {
				t.Fatalf("✗ failed to get resource name: %v", err)
			}

			if resourceName != tt.desired {
				t.Errorf("✗ expected %s, got %s", tt.desired, resourceName)
			}
		})
	}
}

func TestGetGroupVersionKindFromResourceName(t *testing.T) {
	testdata := []struct {
		name     string
		resource string
		desired  []schema.GroupVersionKind
	}{
		{
			name:     "resource-pods",
			resource: "pods",
			desired: []schema.GroupVersionKind{
				{
					Group:   "",
					Version: "v1",
					Kind:    "Pod",
				},
			},
		},
		{
			name:     "resource-deployments",
			resource: "deployments",
			desired: []schema.GroupVersionKind{
				{
					Group:   "apps",
					Version: "v1",
					Kind:    "Deployment",
				},
			},
		},
	}

	for _, tt := range testdata {
		t.Run(tt.name, func(t *testing.T) {
			gvks, err := GetGroupVersionKindFromResourceName(tt.resource)
			if err != nil {
				t.Fatalf("✗ failed to get group version kind: %v", err)
			}

			if len(gvks) != len(tt.desired) {
				t.Errorf("✗ expected %d gvks, got %d", len(tt.desired), len(gvks))
				return
			}

			for i, gvk := range gvks {
				if gvk.Group != tt.desired[i].Group || gvk.Version != tt.desired[i].Version || gvk.Kind != tt.desired[i].Kind {
					t.Errorf("✗ expected %v, got %v", tt.desired[i], gvk)
				}
			}
		})
	}
}

func TestGetGroupVersionResourceFromResourceName(t *testing.T) {
	testdata := []struct {
		name     string
		resource string
		desired  schema.GroupVersionResource
	}{
		{
			name:     "resource-pods",
			resource: "pods",
			desired: schema.GroupVersionResource{
				Group:    "",
				Version:  "v1",
				Resource: "pods",
			},
		},
		{
			name:     "resource-po-singular",
			resource: "po",
			desired: schema.GroupVersionResource{
				Group:    "",
				Version:  "v1",
				Resource: "pods",
			},
		},
		{
			name:     "resource-deployments",
			resource: "deployments",
			desired: schema.GroupVersionResource{
				Group:    "apps",
				Version:  "v1",
				Resource: "deployments",
			},
		},
		{
			name:     "resource-deployments-with-group",
			resource: "deployments.apps",
			desired: schema.GroupVersionResource{
				Group:    "apps",
				Version:  "v1",
				Resource: "deployments",
			},
		},
		{
			name:     "resource-deployments-with-group-version",
			resource: "deployments.v1.apps",
			desired: schema.GroupVersionResource{
				Group:    "apps",
				Version:  "v1",
				Resource: "deployments",
			},
		},
		{
			name:     "resource-roles",
			resource: "roles",
			desired: schema.GroupVersionResource{
				Group:    "rbac.authorization.k8s.io",
				Version:  "v1",
				Resource: "roles",
			},
		},
	}

	for _, tt := range testdata {
		t.Run(tt.name, func(t *testing.T) {
			gvr, err := GetPreferredGroupVersionResourceFor(tt.resource)
			if err != nil {
				t.Fatalf("✗ get group version resource: %v", err)
			}

			if gvr.Group != tt.desired.Group || gvr.Version != tt.desired.Version || gvr.Resource != tt.desired.Resource {
				t.Errorf("expected %v, got %v", tt.desired, gvr)
			}
		})
	}
}
