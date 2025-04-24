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

	krbclient "github.com/ketches/kube-recycle-bin/internal/client"
	"github.com/ketches/kube-recycle-bin/internal/completion"
	"github.com/ketches/kube-recycle-bin/pkg/kube"
	"github.com/ketches/kube-recycle-bin/pkg/tlog"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type RestoreFlags struct {
	ObjectResource  string
	ObjectNamespace string
}

var restoreFlags RestoreFlags

// restoreCmd represents the restore command
var restoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "Restore recycled resource objects from RecycleItem",
	Args:  cobra.MinimumNArgs(1),
	Example: `
# Restore RecycleItem with names foo and bar
krb-cli restore foo bar

# Restore RecycleItem deployments foo and filter by object resource deployments
krb-cli restore --object-resource deployments foo

# Restore RecycleItem deployments foo-deploy, service foo-svc and filter by object namespace dev
krb-cli restore --object-namespace dev foo-deploy foo-svc
`,

	Run: func(cmd *cobra.Command, args []string) {
		runRestore(args)
	},
	ValidArgsFunction: completion.RecycleItem,
}

func init() {
	rootCmd.AddCommand(restoreCmd)

	restoreCmd.Flags().StringVarP(&restoreFlags.ObjectResource, "object-resource", "", "", "Restore recycled resource objects filtered by the specified object resource")
	restoreCmd.Flags().StringVarP(&restoreFlags.ObjectNamespace, "object-namespace", "", "", "Restore recycled resource objects filtered by the specified object namespace")

	restoreCmd.RegisterFlagCompletionFunc("object-resource", completion.RecycleItemGroupResource)
	restoreCmd.RegisterFlagCompletionFunc("object-namespace", completion.RecycleItemNamespace)
}

func runRestore(args []string) {
	if len(args) == 0 {
		tlog.Panicf("✗ please specify recycle items to restore.")
	}

	for _, recycleItemName := range args {
		recycleItem, err := krbclient.RecycleItem().Get(context.Background(), recycleItemName, client.GetOptions{})
		if err != nil {
			tlog.Printf("✗ failed to get RecycleItem [%s]: %v, ignored.", recycleItemName, err)
			continue
		}

		unstructuredObj, err := recycleItem.Object.Unstructured()
		if err != nil {
			tlog.Printf("✗ failed to get unstructured object from RecycleItem [%s]: %v, ignored.", recycleItemName, err)
			continue
		}

		if _, err := kube.DynamicClient().Resource(recycleItem.Object.GroupVersionResource()).Namespace(recycleItem.Object.Namespace).Create(context.Background(), unstructuredObj, metav1.CreateOptions{}); err != nil {
			tlog.Printf("✗ failed to restore recycled resource object [%s]: %v", recycleItem.Object.Key().String(), err)
		} else {
			tlog.Printf("✓ restored recycled resource object [%s: %s] done.", recycleItem.Object.Resource, recycleItem.Object.Key().String())
			// delete the recycle item after successful restore
			if err := krbclient.RecycleItem().Delete(context.Background(), recycleItemName, client.DeleteOptions{}); err != nil {
				tlog.Printf("✗ failed to automatically delete RecycleItem [%s] after restore: %v", recycleItemName, err)
			} else {
				tlog.Printf("✓ automatically deleted RecycleItem [%s] after restore.", recycleItemName)
			}
		}
	}
}
