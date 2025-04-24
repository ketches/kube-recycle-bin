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
	"bytes"
	"encoding/json"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/rand"
	"sigs.k8s.io/yaml"
)

type RecycleItem struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`

	Object RecycledObject `json:"object"`
}

type RecycledObject struct {
	Group     string `json:"group,omitempty"`
	Version   string `json:"version"`
	Kind      string `json:"kind"`
	Resource  string `json:"resource"`
	Namespace string `json:"namespace,omitempty"`
	Name      string `json:"name"`
	Raw       []byte `json:"raw"`
}

type RecycleItemList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []RecycleItem `json:"items"`
}

func NewRecycleItem(recycledObj *RecycledObject) *RecycleItem {
	gr := metav1.GroupResource{
		Group:    recycledObj.Group,
		Resource: recycledObj.Resource,
	}
	labels := map[string]string{
		"krb.ketches.cn/object-name": recycledObj.Name,
		"krb.ketches.cn/object-gr":   gr.String(),
		"krb.ketches.cn/recycled-at": fmt.Sprintf("%d", metav1.Now().Unix()),
	}
	if recycledObj.Namespace != "" {
		labels["krb.ketches.cn/object-namespace"] = recycledObj.Namespace
	}

	return &RecycleItem{
		TypeMeta: metav1.TypeMeta{
			APIVersion: GroupVersion.String(),
			Kind:       RecycleItemKind,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   fmt.Sprintf("%s-%s", recycledObj.Name, rand.String(8)),
			Labels: labels,
		},
		Object: *recycledObj,
	}
}

func (obj *RecycledObject) Key() types.NamespacedName {
	return types.NamespacedName{
		Name:      obj.Name,
		Namespace: obj.Namespace,
	}
}

func (obj *RecycledObject) GroupVersionKind() schema.GroupVersionKind {
	return schema.GroupVersionKind{
		Group:   obj.Group,
		Version: obj.Version,
		Kind:    obj.Kind,
	}
}

func (obj *RecycledObject) GroupVersionResource() schema.GroupVersionResource {
	return schema.GroupVersionResource{
		Group:    obj.Group,
		Version:  obj.Version,
		Resource: obj.Resource,
	}

}

func (obj *RecycledObject) GroupResource() schema.GroupResource {
	return schema.GroupResource{
		Group:    obj.Group,
		Resource: obj.Resource,
	}
}

func (obj *RecycledObject) ObjectGroupKind() schema.GroupKind {
	return schema.GroupKind{
		Group: obj.Group,
		Kind:  obj.Kind,
	}
}

func (obj *RecycledObject) Unstructured() (*unstructured.Unstructured, error) {
	unstructuredObj := &unstructured.Unstructured{}
	if err := json.Unmarshal(obj.Raw, unstructuredObj); err != nil {
		return nil, err
	}

	// Remove the resourceVersion field from the metadata, so it
	// doesn't cause conflicts when creating a new object.
	unstructured.RemoveNestedField(unstructuredObj.Object, "metadata", "resourceVersion")

	return unstructuredObj, nil
}

func (obj *RecycledObject) JSON() string {
	return string(obj.Raw)
}

func (obj *RecycledObject) IndentedJSON() (string, error) {
	var out bytes.Buffer
	if err := json.Indent(&out, obj.Raw, "", "  "); err != nil {
		return "", err
	}

	return out.String(), nil
}

func (obj *RecycledObject) YAML() (string, error) {
	b, err := yaml.JSONToYAML(obj.Raw)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
