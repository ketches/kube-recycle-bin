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
	"github.com/ketches/kube-recycle-bin/pkg/kube"
	"github.com/ketches/kube-recycle-bin/pkg/tlog"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type recycleItemClient struct {
	client.Client
}

func RecycleItem() RecycleItemInterface {
	if recycleItemCli == nil {
		if cli == nil {
			var err error
			cli, err = client.New(kube.RestConfig(), client.Options{Scheme: scheme})
			if err != nil {
				tlog.Fatalf("✗ failed to create client: %v", err)
			}
		}

		recycleItemCli = &recycleItemClient{
			Client: cli,
		}
	}
	return recycleItemCli
}

func (c *recycleItemClient) Create(ctx context.Context, obj *api.RecycleItem, opts client.CreateOptions) error {
	return c.Client.Create(ctx, obj, &opts)
}

func (c *recycleItemClient) Get(ctx context.Context, name string, opts client.GetOptions) (*api.RecycleItem, error) {
	var obj api.RecycleItem
	if err := c.Client.Get(ctx, types.NamespacedName{
		Name: name,
	}, &obj, &opts); err != nil {
		return nil, err
	}
	return &obj, nil
}

func (c *recycleItemClient) List(ctx context.Context, opts client.ListOptions) (*api.RecycleItemList, error) {
	var objList api.RecycleItemList
	if err := c.Client.List(ctx, &objList, &opts); err != nil {
		return nil, err
	}
	return &objList, nil
}

func (c *recycleItemClient) Update(ctx context.Context, obj *api.RecycleItem, opts client.UpdateOptions) error {
	if err := c.Client.Update(ctx, obj, &opts); err != nil {
		return err
	}
	return nil
}

func (c *recycleItemClient) Delete(ctx context.Context, name string, opts client.DeleteOptions) error {
	if err := c.Client.Delete(ctx, &api.RecycleItem{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}, &opts); err != nil {
		return err
	}
	return nil
}
