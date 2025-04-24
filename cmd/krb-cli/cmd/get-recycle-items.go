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
	"github.com/ketches/kube-recycle-bin/pkg/kube"
	"github.com/ketches/kube-recycle-bin/pkg/tlog"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/duration"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type GetRecycleItemFlags struct {
	ObjectResource  string
	ObjectNamespace string
}

var getRecycleItemFlags GetRecycleItemFlags

// getRecycleItemCmd represents the get recycle item command
var getRecycleItemCmd = &cobra.Command{
	Use:     "recycleitems",
	Aliases: []string{"ri", "recycleitem"},
	Short:   "Get recycle items",
	Long:    `Get recycle items. This command retrieves the specified RecycleItem resources.`,
	Example: `
# Get all RecycleItems
krb-cli get ri

# Get RecycleItems with names foo and bar
krb-cli get ri foo bar

# Get RecycleItems recycled from deployments resource
krb-cli get ri --object-resource deployments

# Get RecycleItems recycled from dev namespace
krb-cli get ri --object-namespace dev
`,
	Run: func(cmd *cobra.Command, args []string) {
		runGetRecycleItems(args)
	},
	ValidArgsFunction: completion.RecycleItem,
}

func init() {
	getCmd.AddCommand(getRecycleItemCmd)

	getRecycleItemCmd.Flags().StringVarP(&getRecycleItemFlags.ObjectResource, "object-resource", "", "", "List recycled resource objects filtered by the specified object resource")
	getRecycleItemCmd.Flags().StringVarP(&getRecycleItemFlags.ObjectNamespace, "object-namespace", "", "", "List recycled resource objects filtered by the specified object namespace")

	getRecycleItemCmd.RegisterFlagCompletionFunc("object-resource", completion.RecycleItemGroupResource)
	getRecycleItemCmd.RegisterFlagCompletionFunc("object-namespace", completion.RecycleItemNamespace)
}

func runGetRecycleItems(args []string) {
	var result api.RecycleItemList

	if len(args) > 0 {
		for _, name := range args {
			obj, err := krbclient.RecycleItem().Get(context.Background(), name, client.GetOptions{})
			if err != nil {
				tlog.Errorf("✗ failed to get RecycleItem [%s]: %v, skipping.", name, err)
				continue
			}
			result.Items = append(result.Items, *obj)
		}
	} else {
		labelSet := labels.Set{}
		if getRecycleItemFlags.ObjectNamespace != "" {
			labelSet["krb.ketches.cn/object-namespace"] = getRecycleItemFlags.ObjectNamespace
		}
		if getRecycleItemFlags.ObjectResource != "" {
			tlog.Println(getRecycleItemFlags.ObjectResource)
			if gvr, err := kube.GetPreferredGroupVersionResourceFor(getRecycleItemFlags.ObjectResource); err != nil {
				tlog.Panicf("✗ failed to get preferred group version resource: %v", err)
			} else {
				labelSet["krb.ketches.cn/object-gr"] = gvr.GroupResource().String()
			}
		}

		list, err := krbclient.RecycleItem().List(context.Background(), client.ListOptions{
			LabelSelector: labels.SelectorFromSet(labelSet),
		})
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
		t.AppendRow(table.Row{obj.Name, obj.Object.Group, obj.Object.Version, obj.Object.Resource, obj.Object.Namespace, obj.Object.Name, duration.HumanDuration(time.Since(obj.CreationTimestamp.Time))}, table.RowConfig{
			AutoMerge: true,
		})
	}
	t.SetStyle(KrbTableStyle)
	t.Render()
}
