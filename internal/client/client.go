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

package client

import (
	"context"

	"github.com/ketches/kube-recycle-bin/internal/api"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	scheme = runtime.NewScheme()
)

func init() {
	api.AddToScheme(scheme)
}

var (
	cli              client.Client
	recycleItemCli   RecycleItemInterface
	recyclePolicyCli RecyclePolicyInterface
)

type RecycleItemInterface interface {
	Create(ctx context.Context, obj *api.RecycleItem, opts client.CreateOptions) error
	Get(ctx context.Context, name string, opts client.GetOptions) (*api.RecycleItem, error)
	List(ctx context.Context, opts client.ListOptions) (*api.RecycleItemList, error)
	Update(ctx context.Context, obj *api.RecycleItem, opts client.UpdateOptions) error
	Delete(ctx context.Context, name string, opts client.DeleteOptions) error
}

type RecyclePolicyInterface interface {
	Create(ctx context.Context, obj *api.RecyclePolicy, opts client.CreateOptions) error
	Get(ctx context.Context, name string, opts client.GetOptions) (*api.RecyclePolicy, error)
	List(ctx context.Context, opts client.ListOptions) (*api.RecyclePolicyList, error)
	Update(ctx context.Context, obj *api.RecyclePolicy, opts client.UpdateOptions) error
	Delete(ctx context.Context, name string, opts client.DeleteOptions) error
}
