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

package cmd

import (
	"context"
	"strings"

	krbclient "github.com/ketches/kube-recycle-bin/internal/client"
	"github.com/ketches/kube-recycle-bin/internal/completion"
	"github.com/ketches/kube-recycle-bin/pkg/tlog"
	"github.com/spf13/cobra"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ViewFlags struct {
	ObjectResource  string
	ObjectNamespace string
	OutputFormat    string
}

var viewFlags ViewFlags

var viewCmd = &cobra.Command{
	Use:     "view",
	Aliases: []string{"show", "display"},
	Short:   "View recyceled resource objects from RecycleItem",
	Example: `
# View recycled resource objects from RecycleItem with names foo and bar
krb-cli view foo bar

# View recycled resource objects from RecycleItem deployments foo and filter by object resource deployments
krb-cli view --object-resource deployments foo

# View recycled resource objects from RecycleItem deployments foo-deploy, service foo-svc and filter by object namespace dev
krb-cli view --object-namespace dev foo-deploy foo-svc

# View recycled resource objects from RecycleItem with names foo and bar in JSON format
krb-cli view foo bar --output json
`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		runView(args)
	},
	ValidArgsFunction: completion.RecycleItem,
}

func init() {
	rootCmd.AddCommand(viewCmd)

	viewCmd.Flags().StringVarP(&viewFlags.ObjectResource, "object-resource", "", "", "View recycled resource objects filtered by the specified object resource")
	viewCmd.Flags().StringVarP(&viewFlags.ObjectNamespace, "object-namespace", "", "", "View recycled resource objects filtered by the specified object namespace")
	viewCmd.Flags().StringVarP(&viewFlags.OutputFormat, "output", "o", "yaml", "Output format. One of: json|yaml, default is yaml")

	viewCmd.RegisterFlagCompletionFunc("object-resource", completion.RecycleItemGroupResource)
	viewCmd.RegisterFlagCompletionFunc("object-namespace", completion.RecycleItemNamespace)
	viewCmd.RegisterFlagCompletionFunc("output", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"json", "yaml"}, cobra.ShellCompDirectiveNoFileComp
	})
}

func runView(args []string) {
	if len(args) == 0 {
		tlog.Panicf("✗ please specify recycle items to view.")
	}

	var objContents []string
	for _, recycleItemName := range args {
		recycleItem, err := krbclient.RecycleItem().Get(context.Background(), recycleItemName, client.GetOptions{})
		if err != nil {
			tlog.Printf("✗ failed to get RecycleItem [%s]: %v, ignored.", recycleItemName, err)
			continue
		}

		switch viewFlags.OutputFormat {
		case "json":
			objContent, err := recycleItem.Object.IndentedJSON()
			if err != nil {
				tlog.Printf("✗ failed to view recycled resource object [%s: %s] from RecycleItem [%s] in JSON format: %s, error: %v", recycleItem.Object.Kind, recycleItem.Object.Key().String(), recycleItem.Name, recycleItem.Name, err)
				continue
			}
			objContents = append(objContents, objContent)
			// tlog.Printf("✓ recycled resource object [%s: %s] from RecycleItem [%s]: \n%s \n", recycleItem.Object.Kind, recycleItem.Object.Key().String(), recycleItem.Name, objContent)
		default:
			objContent, err := recycleItem.Object.YAML()
			if err != nil {
				tlog.Printf("✗ failed to view recycled resource object [%s: %s] from RecycleItem [%s] in YAML format: %s, error: %v", recycleItem.Object.Kind, recycleItem.Object.Key().String(), recycleItem.Name, recycleItem.Name, err)
				continue
			}
			objContents = append(objContents, objContent)
			// tlog.Printf("✓ recycled resource object [%s: %s] from RecycleItem [%s]: \n%s \n", recycleItem.Object.Kind, recycleItem.Object.Key().String(), recycleItem.Name, objContent)
		}
	}
	tlog.Print(strings.Join(objContents, "\n---\n"))
}
