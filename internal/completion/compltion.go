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
	"slices"

	krbclient "github.com/ketches/kube-recycle-bin/internal/client"
	"github.com/ketches/kube-recycle-bin/pkg/kube"
	"github.com/ketches/kube-recycle-bin/pkg/tlog"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/labels"
	client "sigs.k8s.io/controller-runtime/pkg/client"
)

// None is a shell completion function that does nothing.
func None(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return nil, cobra.ShellCompDirectiveNoFileComp
}

func RecycleItemGroupResource(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	list, err := krbclient.RecycleItem().List(context.Background(), client.ListOptions{})
	if err != nil {
		tlog.Printf("✗ failed to list recycle items: %v", err)
		return nil, cobra.ShellCompDirectiveError
	}

	var result []string
	for _, item := range list.Items {
		result = append(result, item.Object.GroupResource().String())
	}

	return result, cobra.ShellCompDirectiveNoFileComp
}

func RecycleItemNamespace(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	list, err := krbclient.RecycleItem().List(context.Background(), client.ListOptions{})
	if err != nil {
		tlog.Printf("✗ failed to list recycle items: %v", err)
		return nil, cobra.ShellCompDirectiveError
	}

	var result []string
	for _, item := range list.Items {
		result = append(result, item.Object.Namespace)
	}

	return result, cobra.ShellCompDirectiveNoFileComp
}

// KubeGroupResources is a shell completion function that lists all group resources.
func KubeGroupResources(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) > 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	resources, err := kube.GetAllGroupResources()
	if err != nil {
		tlog.Printf("✗ failed to get all group resources: %v", err)
		return nil, cobra.ShellCompDirectiveError
	}

	var result []string
	for _, resource := range resources {
		if slices.Contains(args, resource) {
			continue
		}
		result = append(result, resource)
	}

	return result, cobra.ShellCompDirectiveNoFileComp
}

// RecycleItem is a shell completion function that lists all recycle items.
func RecycleItem(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	labelSet := labels.Set{}
	objectNamespace, _ := cmd.Flags().GetString("object-namespace")
	if objectNamespace != "" {
		labelSet["krb.ketches.cn/object-namespace"] = objectNamespace
	}
	objectResource, _ := cmd.Flags().GetString("object-resource")
	if objectResource != "" {
		if gvr, err := kube.GetPreferredGroupVersionResourceFor(objectResource); err != nil {
			tlog.Printf("✗ failed to get preferred group version resource: %v", err)
		} else {
			labelSet["krb.ketches.cn/object-gr"] = gvr.GroupResource().String()
		}
	}

	list, err := krbclient.RecycleItem().List(context.Background(), client.ListOptions{
		LabelSelector: labels.SelectorFromSet(labelSet),
	})
	if err != nil {
		tlog.Printf("✗ failed to list recycle items: %v", err)
		return nil, cobra.ShellCompDirectiveError
	}

	var result []string
	for _, obj := range list.Items {
		if slices.Contains(args, obj.Name) {
			continue
		}
		result = append(result, obj.Name)
	}

	return result, cobra.ShellCompDirectiveNoFileComp
}

func RecyclePolicyGroupResource(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	list, err := krbclient.RecyclePolicy().List(context.Background(), client.ListOptions{})
	if err != nil {
		tlog.Printf("✗ failed to list recycle items: %v", err)
		return nil, cobra.ShellCompDirectiveError
	}

	var result []string
	for _, item := range list.Items {
		result = append(result, item.Target.GroupResource().String())
	}

	return result, cobra.ShellCompDirectiveNoFileComp
}

func RecyclePolicyNamespace(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	ri, err := krbclient.RecyclePolicy().List(context.Background(), client.ListOptions{})
	if err != nil {
		tlog.Printf("✗ failed to list recycle items: %v", err)
		return nil, cobra.ShellCompDirectiveError
	}

	var result []string
	for _, item := range ri.Items {
		for _, ns := range item.Target.Namespaces {
			if ns == "" {
				continue
			}
			result = append(result, ns)
		}
	}

	return result, cobra.ShellCompDirectiveNoFileComp
}

// RecyclePolicy is a shell completion function that lists all recycle policies.
func RecyclePolicy(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	labelSet := labels.Set{}
	targetNamespace, _ := cmd.Flags().GetString("target-namespace")
	if targetNamespace != "" {
		labelSet["krb.ketches.cn/target-namespace"] = targetNamespace
	}
	targetResource, _ := cmd.Flags().GetString("target-resource")
	if targetResource != "" {
		if gvr, err := kube.GetPreferredGroupVersionResourceFor(targetResource); err != nil {
			tlog.Printf("✗ failed to get preferred group version resource: %v", err)
		} else {
			labelSet["krb.ketches.cn/target-gr"] = gvr.GroupResource().String()
		}
	}

	list, err := krbclient.RecyclePolicy().List(context.Background(), client.ListOptions{
		LabelSelector: labels.SelectorFromSet(labelSet),
	})
	if err != nil {
		tlog.Printf("✗ failed to list recycle policies: %v", err)
		return nil, cobra.ShellCompDirectiveError
	}

	var result []string
	for _, obj := range list.Items {
		if slices.Contains(args, obj.Name) {
			continue
		}
		result = append(result, obj.Name)
	}

	return result, cobra.ShellCompDirectiveNoFileComp
}
