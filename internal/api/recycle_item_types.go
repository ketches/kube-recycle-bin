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
	"encoding/json"
	"fmt"

	"github.com/ketches/kube-recycle-bin/pkg/kube"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/rand"
)

type RecycleItem struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`

	Object Object `json:"object"`
}

type Object struct {
	Group     string `json:"group,omitempty"`
	Version   string `json:"version"`
	Kind      string `json:"kind"`
	Namespace string `json:"namespace,omitempty"`
	Name      string `json:"name"`
	Raw       []byte `json:"raw"`
}

type RecycleItemList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []RecycleItem `json:"items"`
}

func NewRecycleItem(gvk metav1.GroupVersionKind, namespace, name string, raw []byte) *RecycleItem {
	labels := map[string]string{
		"krb.ketches.cn/object-name":    name,
		"krb.ketches.cn/object-kind":    gvk.Kind,
		"krb.ketches.cn/object-version": gvk.Version,
		"krb.ketches.cn/recycled-at":    fmt.Sprintf("%d", metav1.Now().Unix()),
	}
	if namespace != "" {
		labels["krb.ketches.cn/object-namespace"] = namespace
	}
	if gvk.Group != "" {
		labels["krb.ketches.cn/object-group"] = gvk.Group
	}

	return &RecycleItem{
		TypeMeta: metav1.TypeMeta{
			APIVersion: GroupVersion.String(),
			Kind:       RecycleItemKind,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   fmt.Sprintf("%s-%s", name, rand.String(8)),
			Labels: labels,
		},
		Object: Object{
			Group:     gvk.Group,
			Version:   gvk.Version,
			Kind:      gvk.Kind,
			Namespace: namespace,
			Name:      name,
			Raw:       raw,
		},
	}
}

func (r *RecycleItem) ObjectName() string {
	return r.Object.Name
}

func (r *RecycleItem) ObjectNamespace() string {
	return r.Object.Namespace
}

func (r *RecycleItem) ObjectNamespacedName() types.NamespacedName {
	return types.NamespacedName{
		Name:      r.Object.Name,
		Namespace: r.Object.Namespace,
	}
}

func (r *RecycleItem) ObjectGroupVersionKind() schema.GroupVersionKind {
	return schema.GroupVersionKind{
		Group:   r.Object.Group,
		Version: r.Object.Version,
		Kind:    r.Object.Kind,
	}
}

func (r *RecycleItem) ObjectGroupVersionResource() schema.GroupVersionResource {
	gvr, err := kube.GetGroupVersionResourceFromGroupVersionKind(r.ObjectGroupVersionKind())
	if err != nil {
		return schema.GroupVersionResource{}
	}
	return gvr
}

func (r *RecycleItem) ObjectRaw() []byte {
	return r.Object.Raw
}

func (r *RecycleItem) ObjectUnstructured() (*unstructured.Unstructured, error) {
	unstructuredObj := &unstructured.Unstructured{}
	if err := json.Unmarshal(r.Object.Raw, unstructuredObj); err != nil {
		return nil, err
	}

	sanitizeMetadata(unstructuredObj)

	return unstructuredObj, nil
}

func sanitizeMetadata(u *unstructured.Unstructured) {
	// Remove the resourceVersion field from the metadata, so it
	// doesn't cause conflicts when creating a new object.
	unstructured.RemoveNestedField(u.Object, "metadata", "resourceVersion")
}
