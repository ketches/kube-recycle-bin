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

type GetRecycleItemFlags struct {
	Group         string
	Version       string
	Kind          string
	Namespace     string
	AllNamespaces bool
}

var getRecycylItemFlags GetRecycleItemFlags

// getRecycleItemCmd represents the get recycle item command
var getRecycleItemCmd = &cobra.Command{
	Use:     "recycleitems",
	Aliases: []string{"ri", "recycleitem"},
	Short:   "Get recycle items",
	Long:    `Get recycle items. This command retrieves the specified RecycleItem resources.`,
	Example: `
# Get RecycleItems with names foo and bar
krb-cli get recycleitems foo bar

# Get all RecyclePolicies
krb-cli get recyclepolicies
`,
	Run: func(cmd *cobra.Command, args []string) {
		runGetRecycleItems(args)
	},
	ValidArgsFunction: completion.RecycleItem,
}

func init() {
	getCmd.AddCommand(getRecycleItemCmd)

	getRecycleItemCmd.Flags().StringVarP(&getRecycylItemFlags.Group, "group", "g", "", "List resources of the specified group")
	getRecycleItemCmd.Flags().StringVarP(&getRecycylItemFlags.Kind, "kind", "k", "", "List resources of the specified kind")
	getRecycleItemCmd.Flags().StringVarP(&getRecycylItemFlags.Version, "version", "v", "", "List resources with the specified version")
	getRecycleItemCmd.Flags().StringVarP(&getRecycylItemFlags.Namespace, "namespace", "n", "default", "List resources in the specified namespace")
	getRecycleItemCmd.Flags().BoolVarP(&getRecycylItemFlags.AllNamespaces, "all-namespaces", "A", false, "List resources from all namespaces")

	getRecycleItemCmd.MarkFlagsMutuallyExclusive("namespace", "all-namespaces")
}

func runGetRecycleItems(args []string) {
	var result api.RecycleItemList
	if len(args) > 0 {
		for _, objName := range args {
			obj, err := krbclient.RecycleItem().Get(context.Background(), objName, client.GetOptions{})
			if err != nil {
				tlog.Errorf("✗ failed to get RecycleItem [%s]: %v, skipping.", objName, err)
				continue
			}
			result.Items = append(result.Items, *obj)
		}
	} else {
		list, err := krbclient.RecycleItem().List(context.Background(), client.ListOptions{})
		if err != nil {
			tlog.Panicf("✗ failed to list RecycleItem: %v", err)
			return
		}
		result = *list
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"name", "group", "version", "kind", "namespace", "name", "age"})

	for _, obj := range result.Items {
		t.AppendRow(table.Row{obj.Name, obj.Object.Group, obj.Object.Version, obj.Object.Kind, obj.Object.Namespace, obj.Object.Name, duration.HumanDuration(time.Since(obj.CreationTimestamp.Time))}, table.RowConfig{
			AutoMerge: true,
		})
	}
	t.SetStyle(KrbTableStyle)
	t.Render()
}
