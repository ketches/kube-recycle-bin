/*
Copyright © 2025 The Ketches Authors.

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

package cmd

import (
	"context"

	"github.com/ketches/kube-recycle-bin/internal/api"
	krbclient "github.com/ketches/kube-recycle-bin/internal/client"
	"github.com/ketches/kube-recycle-bin/internal/completion"
	"github.com/ketches/kube-recycle-bin/pkg/kube"
	"github.com/ketches/kube-recycle-bin/pkg/tlog"
	"github.com/spf13/cobra"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type RecycleFlags struct {
	TargetNamespaces []string
}

var recycleFlags RecycleFlags

// recycleCmd represents the create policy command
var recycleCmd = &cobra.Command{
	Use:   "recycle",
	Short: "Recycle specified resources",
	Long:  `Recycle specified resources. This command creates a RecyclePolicy for the specified resource type.`,
	Example: `# Recycle Deployment in dev and prod namespace
krb-cli recycle deployments -n dev,prod

# Recycle service in all namespaces
krb-cli recycle services
`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		runRecycle(args)
	},
	ValidArgsFunction: completion.KubeGroupResources,
}

func init() {
	rootCmd.AddCommand(recycleCmd)

	recycleCmd.Flags().StringSliceVarP(&recycleFlags.TargetNamespaces, "target-namespaces", "n", []string{}, "Create a RecyclePolicy with specific target namespaces")
}

func runRecycle(args []string) {
	if len(args) == 0 {
		tlog.Panicf("✗ please specify a resource to recycle.")
	}

	for _, resource := range args {
		gvr, err := kube.GetPreferredGroupVersionResourceFor(resource)
		if err != nil {
			tlog.Errorf("✗ failed to get gvr from resource name: %v, ignored.", err)
			continue
		}
		if gvr == nil {
			tlog.Errorf("✗ no resources found for %s, ignored.", resource)
			continue
		}

		recycleItem := api.NewRecyclePolicy(*gvr, recycleFlags.TargetNamespaces)
		if err := krbclient.RecyclePolicy().Create(context.Background(), recycleItem, client.CreateOptions{}); err != nil {
			tlog.Panicf("✗ failed to create recycle policy: %v, ignored.", err)
			continue
		}
		tlog.Printf("✓ create recycle policy [%s] done.", recycleItem.Name)
	}
}
