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
	"bytes"
	"context"
	"encoding/json"
	"os"
	"strings"
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
	"sigs.k8s.io/yaml"
)

type GetRecyclePoliciesFlags struct {
	TargetResource  string
	TargetNamespace string
	OutputFormat    string
}

var getRecyclePoliciesFlags GetRecyclePoliciesFlags

// getRecyclePoliciesCmd represents the get recycle policies command
var getRecyclePoliciesCmd = &cobra.Command{
	Use:     "recyclepolicies",
	Aliases: []string{"rp", "recyclepolicy"},
	Short:   "Get recycle policies",
	Long:    `Get recycle policies. This command retrieves the specified RecyclePolicy resources.`,
	Example: `
# Get all RecyclePolicy
krb-cli get recyclepolicies

# Get RecyclePolicy with names foo and bar
krb-cli get recyclepolicies foo bar

# Get RecyclePolicy for deployments resource
krb-cli get recyclepolicies --target-resource deployments

# Get RecyclePolicy for default namespace
krb-cli get recyclepolicies --target-namespace default
`,
	Run: func(cmd *cobra.Command, args []string) {
		runGetRecyclePolicies(args)
	},
	ValidArgsFunction: completion.RecyclePolicy,
}

func init() {
	getCmd.AddCommand(getRecyclePoliciesCmd)

	getRecyclePoliciesCmd.Flags().StringVarP(&getRecyclePoliciesFlags.TargetResource, "target-resource", "", "", "List recycle policies filtered by the specified target resource")
	getRecyclePoliciesCmd.Flags().StringVarP(&getRecyclePoliciesFlags.TargetNamespace, "target-namespace", "", "", "List recycle policies filtered by the specified target namespace")
	getRecyclePoliciesCmd.Flags().StringVarP(&getRecyclePoliciesFlags.OutputFormat, "output", "o", "", "Output format. One of: json|yaml")

	getRecyclePoliciesCmd.RegisterFlagCompletionFunc("target-resource", completion.RecyclePolicyGroupResource)
	getRecyclePoliciesCmd.RegisterFlagCompletionFunc("target-namespace", completion.RecyclePolicyNamespace)
	getRecyclePoliciesCmd.RegisterFlagCompletionFunc("output", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"json", "yaml"}, cobra.ShellCompDirectiveDefault
	})
}

func runGetRecyclePolicies(args []string) {
	var result api.RecyclePolicyList
	if len(args) > 0 {
		for _, name := range args {
			obj, err := krbclient.RecyclePolicy().Get(context.Background(), name, client.GetOptions{})
			if err != nil {
				tlog.Errorf("✗ failed to get RecyclePolicy [%s]: %v, ignored.", name, err)
				continue
			}
			result.Items = append(result.Items, *obj)
		}
	} else {
		labelSet := labels.Set{}
		if getRecyclePoliciesFlags.TargetResource != "" {
			if gvr, err := kube.GetPreferredGroupVersionResourceFor(getRecyclePoliciesFlags.TargetResource); err != nil {
				tlog.Panicf("✗ failed to get preferred group version resource: %v", err)
			} else {
				labelSet["krb.ketches.cn/target-gr"] = gvr.GroupResource().String()
			}
		}
		list, err := krbclient.RecyclePolicy().List(context.Background(), client.ListOptions{
			LabelSelector: labels.SelectorFromSet(labelSet),
			Namespace:     getRecyclePoliciesFlags.TargetNamespace,
		})
		if err != nil {
			tlog.Panicf("✗ failed to list RecyclePolicy: %v", err)
			return
		}
		result = *list
	}

	if len(result.Items) == 0 {
		tlog.Println("No recycle items found.")
		return
	}

	switch getRecyclePoliciesFlags.OutputFormat {
	case "yaml":
		output, err := yaml.Marshal(result)
		if err != nil {
			tlog.Panicf("failed to marshal recycle policies to yaml: %v", err)
		}
		tlog.Print(string(output))
	case "json":
		y, err := yaml.Marshal(result)
		if err != nil {
			tlog.Panicf("failed to marshal recycle policies to json: %v", err)
		}
		j, err := yaml.YAMLToJSON(y)
		if err != nil {
			tlog.Panicf("failed to convert recycle policies to json: %v", err)
		}
		var output bytes.Buffer
		if err := json.Indent(&output, j, "", "  "); err != nil {
			tlog.Panicf("failed to indent recycle policies json: %v", err)
		}
		tlog.Println(output.String())
	default:
		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"Name", "Target GR", "Target Namespaces", "Age"})

		for _, obj := range result.Items {
			t.AppendRow(table.Row{obj.Name, obj.Target.GroupResource().String(), strings.Join(obj.Target.Namespaces, ","), duration.HumanDuration(time.Since(obj.CreationTimestamp.Time))}, table.RowConfig{
				AutoMerge: true,
			})
		}
		t.SetStyle(KrbTableStyle)
		t.Render()
	}

}
