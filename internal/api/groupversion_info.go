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
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	Group   = "krb.ketches.cn"
	Version = "v1"
)

const (
	RecycleItemKind       = "RecycleItem"
	RecycleItemListKind   = "RecycleItemList"
	RecyclePolicyKind     = "RecyclePolicy"
	RecyclePolicyListKind = "RecyclePolicyList"
)

var (
	GroupVersion = schema.GroupVersion{
		Group:   Group,
		Version: Version,
	}
	SchemeBuilder = runtime.NewSchemeBuilder(addKnownTypes)
	AddToScheme   = SchemeBuilder.AddToScheme
)

// addKnownTypes registers the known types for the API group.
func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(GroupVersion,
		&RecycleItem{},
		&RecycleItemList{},
		&RecyclePolicy{},
		&RecyclePolicyList{},
	)
	metav1.AddToGroupVersion(scheme, GroupVersion)
	return nil
}
