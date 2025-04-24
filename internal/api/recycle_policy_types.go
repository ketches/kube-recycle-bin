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

package api

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/rand"
)

type RecyclePolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`

	Target RecycleTarget `json:"target"`
}

type RecycleTarget struct {
	Group      string   `json:"group,omitempty"`
	Resource   string   `json:"resource,omitempty"`
	Namespaces []string `json:"namespaces,omitempty"`
}

type RecyclePolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []RecyclePolicy `json:"items"`
}

func NewRecyclePolicy(gvr schema.GroupVersionResource, targetNamespaces []string) *RecyclePolicy {
	if len(targetNamespaces) == 0 {
		targetNamespaces = []string{metav1.NamespaceAll}
	}

	labels := map[string]string{
		"krb.ketches.cn/target-gr": gvr.GroupResource().String(),
	}
	for _, ns := range targetNamespaces {
		if ns != metav1.NamespaceAll {
			labels["krb.ketches.cn/target-namespace-"+ns] = "true"
		}
	}

	return &RecyclePolicy{
		TypeMeta: metav1.TypeMeta{
			APIVersion: GroupVersion.String(),
			Kind:       RecyclePolicyKind,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   "recycle-" + gvr.Resource + "-" + rand.String(8),
			Labels: labels,
		},
		Target: RecycleTarget{
			Group:      gvr.Group,
			Resource:   gvr.Resource,
			Namespaces: targetNamespaces,
		},
	}
}

func (rt *RecycleTarget) GroupResource() schema.GroupResource {
	return schema.GroupResource{
		Group:    rt.Group,
		Resource: rt.Resource,
	}
}
