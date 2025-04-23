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

import "k8s.io/apimachinery/pkg/runtime"

func (in *RecycleItem) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

func (in *RecycleItem) DeepCopy() *RecycleItem {
	if in == nil {
		return nil
	}

	out := new(RecycleItem)
	in.DeepCopyInto(out)
	return out
}

func (in *RecycleItem) DeepCopyInto(out *RecycleItem) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ObjectMeta = in.ObjectMeta
	out.Object = in.Object
}

func (in *RecycleItemList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

func (in *RecycleItemList) DeepCopy() *RecycleItemList {
	if in == nil {
		return nil
	}
	out := new(RecycleItemList)
	in.DeepCopyInto(out)
	return out
}

func (in *RecycleItemList) DeepCopyInto(out *RecycleItemList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)

	if in.Items != nil {
		out.Items = make([]RecycleItem, len(in.Items))
		for i := range in.Items {
			in.Items[i].DeepCopyInto(&out.Items[i])
		}
	}
}
