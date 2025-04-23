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

func (in *RecyclePolicy) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

func (in *RecyclePolicy) DeepCopy() *RecyclePolicy {
	if in == nil {
		return nil
	}

	out := new(RecyclePolicy)
	in.DeepCopyInto(out)
	return out
}

func (in *RecyclePolicy) DeepCopyInto(out *RecyclePolicy) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ObjectMeta = in.ObjectMeta
	out.Target = in.Target
}

func (in *RecyclePolicyList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

func (in *RecyclePolicyList) DeepCopy() *RecyclePolicyList {
	if in == nil {
		return nil
	}
	out := new(RecyclePolicyList)
	in.DeepCopyInto(out)
	return out
}

func (in *RecyclePolicyList) DeepCopyInto(out *RecyclePolicyList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)

	if in.Items != nil {
		out.Items = make([]RecyclePolicy, len(in.Items))
		for i := range in.Items {
			in.Items[i].DeepCopyInto(&out.Items[i])
		}
	}
}
