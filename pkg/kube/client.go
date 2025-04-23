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
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	controllerruntime "sigs.k8s.io/controller-runtime"
)

var (
	restConfig      *rest.Config
	client          kubernetes.Interface
	dynamicClient   dynamic.Interface
	discoveryClient discovery.DiscoveryInterface
)

func RestConfig() *rest.Config {
	if restConfig == nil {
		restConfig = controllerruntime.GetConfigOrDie()
	}
	return restConfig
}

func Client() kubernetes.Interface {
	if client == nil {
		client = kubernetes.NewForConfigOrDie(RestConfig())
	}
	return client
}

func DynamicClient() dynamic.Interface {
	if dynamicClient == nil {
		dynamicClient = dynamic.NewForConfigOrDie(RestConfig())
	}
	return dynamicClient
}

func DiscoveryClient() discovery.DiscoveryInterface {
	if discoveryClient == nil {
		discoveryClient = Client().Discovery()
	}
	return discoveryClient
}
