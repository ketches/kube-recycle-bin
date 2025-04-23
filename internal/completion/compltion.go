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

package completion

import (
	"context"

	krbclient "github.com/ketches/kube-recycle-bin/internal/client"
	"github.com/ketches/kube-recycle-bin/pkg/kube"
	"github.com/ketches/kube-recycle-bin/pkg/tlog"
	"github.com/spf13/cobra"
	client "sigs.k8s.io/controller-runtime/pkg/client"
)

// None is a shell completion function that does nothing.
func None(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return nil, cobra.ShellCompDirectiveNoFileComp
}

// KubeGroupResources is a shell completion function that lists all group resources.
func KubeGroupResources(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) > 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	resources, err := kube.GetAllGroupResources()
	if err != nil {
		tlog.Printf("✗ list resources: %v", err)
	}

	return resources, cobra.ShellCompDirectiveNoFileComp
}

// RecycleItem is a shell completion function that lists all recycle items.
func RecycleItem(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	list, err := krbclient.RecycleItem().List(context.Background(), client.ListOptions{})
	if err != nil {
		tlog.Printf("✗ failed to list recycle items: %v", err)
	}
	var completions []string
	for _, obj := range list.Items {
		completions = append(completions, obj.Name)
	}

	return completions, cobra.ShellCompDirectiveNoFileComp
}

// RecyclePolicy is a shell completion function that lists all recycle policies.
func RecyclePolicy(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	list, err := krbclient.RecyclePolicy().List(context.Background(), client.ListOptions{})
	if err != nil {
		tlog.Printf("✗ failed to list recycle policies: %v", err)
	}
	var completions []string
	for _, obj := range list.Items {
		completions = append(completions, obj.Name)
	}

	return completions, cobra.ShellCompDirectiveNoFileComp
}
