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
	"os"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/ketches/kube-recycle-bin/internal/api"
	krbclient "github.com/ketches/kube-recycle-bin/internal/client"
	"github.com/ketches/kube-recycle-bin/internal/completion"
	"github.com/ketches/kube-recycle-bin/pkg/tlog"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/util/duration"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type GetRecyclePoliciesFlags struct {
	Group         string
	Version       string
	Kind          string
	Namespace     string
	AllNamespaces bool
}

var getRecycylPoliciesFlags GetRecyclePoliciesFlags

// getRecyclePoliciesCmd represents the get recycle policies command
var getRecyclePoliciesCmd = &cobra.Command{
	Use:     "recyclepolicies",
	Aliases: []string{"rp", "recyclepolicy"},
	Short:   "Get recycle policies",
	Long:    `Get recycle policies. This command retrieves the specified RecyclePolicy resources.`,
	Example: `
# Get RecyclePolicy with names foo and bar
krb-cli get recyclepolicies foo bar

# Get all RecyclePolicy
krb-cli get recyclepolicies
`,
	Run: func(cmd *cobra.Command, args []string) {
		runGetRecyclePolicies(args)
	},
	ValidArgsFunction: completion.RecyclePolicy,
}

func init() {
	getCmd.AddCommand(getRecyclePoliciesCmd)

	getRecyclePoliciesCmd.Flags().StringVarP(&getRecycylPoliciesFlags.Group, "group", "g", "", "List resources of the specified group")
	getRecyclePoliciesCmd.Flags().StringVarP(&getRecycylPoliciesFlags.Kind, "kind", "k", "", "List resources of the specified kind")
	getRecyclePoliciesCmd.Flags().StringVarP(&getRecycylPoliciesFlags.Version, "version", "v", "", "List resources with the specified version")
	getRecyclePoliciesCmd.Flags().StringVarP(&getRecycylPoliciesFlags.Namespace, "namespace", "n", "default", "List resources in the specified namespace")
	getRecyclePoliciesCmd.Flags().BoolVarP(&getRecycylPoliciesFlags.AllNamespaces, "all-namespaces", "A", false, "List resources from all namespaces")

	getRecyclePoliciesCmd.MarkFlagsMutuallyExclusive("namespace", "all-namespaces")
}

func runGetRecyclePolicies(args []string) {
	var result api.RecyclePolicyList
	if len(args) > 0 {
		for _, objName := range args {
			obj, err := krbclient.RecyclePolicy().Get(context.Background(), objName, client.GetOptions{})
			if err != nil {
				tlog.Errorf("✗ failed to get RecyclePolicy [%s]: %v, skipping.", objName, err)
				continue
			}
			result.Items = append(result.Items, *obj)
		}
	} else {
		list, err := krbclient.RecyclePolicy().List(context.Background(), client.ListOptions{})
		if err != nil {
			tlog.Panicf("✗ failed to list RecyclePolicy: %v", err)
			return
		}
		result = *list
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"name", "group", "resource", "namespaces", "age"})

	for _, obj := range result.Items {
		t.AppendRow(table.Row{obj.Name, obj.Target.Group, obj.Target.Resource, obj.Target.Namespaces, duration.HumanDuration(time.Since(obj.CreationTimestamp.Time))}, table.RowConfig{
			AutoMerge: true,
		})
	}
	t.SetStyle(KrbTableStyle)
	t.Render()
}
