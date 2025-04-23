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
	"testing"

	"github.com/ketches/kube-recycle-bin/internal/api"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func TestCreateRecycleItem(t *testing.T) {
	// Create a new RecycleItem object
	obj := &api.RecycleItem{
		TypeMeta: metav1.TypeMeta{
			APIVersion: api.GroupVersion.String(),
			Kind:       api.RecycleItemKind,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-recycle-item",
		},
		Object: api.Object{
			Group:     "apps",
			Version:   "v1",
			Kind:      "Deployment",
			Namespace: "default",
			Name:      "nginx-deployment",
			Raw: []byte(`apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
spec:
  selector:
    matchLabels:
      app: nginx-deployment
  template:
    metadata:
      labels:
        app: nginx-deployment
    spec:
      containers:
      - name: nginx-deployment
        image: nginx
        resources:
          requests:
            memory: "64Mi"
            cpu: "100m"
          limits:
            memory: "128Mi"
            cpu: "250m"
        ports:
        - containerPort: 80`),
		},
	}

	// Create the RecycleItem object
	if err := RecycleItem().Create(context.Background(), obj, client.CreateOptions{}); err != nil {
		t.Fatalf("✗ failed to create RecycleItem object: %v", err)
	}

	// Retrieve the RecycleItem object
	if got, err := RecycleItem().Get(context.Background(), obj.Name, client.GetOptions{}); err != nil {
		t.Fatalf("✗ failed to get RecycleItem object: %v", err)
	} else if got.Name != obj.Name {
		t.Errorf("got name %s, want %s", got.Name, obj.Name)
	}

	// Clean up the RecycleItem object
	if err := RecycleItem().Delete(context.Background(), obj.Name, client.DeleteOptions{}); err != nil {
		t.Fatalf("✗ failed to delete RecycleItem object: %v", err)
	}
}

func TestCreateRecyclePolicy(t *testing.T) {
	// Create a new RecyclePolicy object
	obj := &api.RecyclePolicy{
		TypeMeta: metav1.TypeMeta{
			APIVersion: api.GroupVersion.String(),
			Kind:       api.RecyclePolicyKind,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-recycle-policy",
		},
		Target: api.Target{
			Group:      "apps",
			Resource:   "deployments",
			Namespaces: []string{"default"},
		},
	}

	// Create the RecyclePolicy object
	if err := RecyclePolicy().Create(context.Background(), obj, client.CreateOptions{}); err != nil {
		t.Fatalf("✗ failed to create RecyclePolicy object: %v", err)
	}

	// Retrieve the RecyclePolicy object
	if got, err := RecyclePolicy().Get(context.Background(), obj.Name, client.GetOptions{}); err != nil {
		t.Fatalf("✗ failed to get RecyclePolicy object: %v", err)
	} else if got.Name != obj.Name {
		t.Errorf("✗ got name %s, want %s", got.Name, obj.Name)
	}

	// Clean up the RecyclePolicy object
	if err := RecyclePolicy().Delete(context.Background(), obj.Name, client.DeleteOptions{}); err != nil {
		t.Fatalf("✗ failed to delete RecyclePolicy object: %v", err)
	}
}
