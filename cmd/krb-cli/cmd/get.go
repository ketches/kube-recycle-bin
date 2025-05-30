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
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

var KrbTableStyle = table.Style{
	Name:    "KrbTableStyle",
	Box:     table.StyleBoxDefault,
	Color:   table.ColorOptionsDefault,
	Format:  table.FormatOptionsDefault,
	HTML:    table.DefaultHTMLOptions,
	Options: table.OptionsNoBordersAndSeparators,
	Size:    table.SizeOptionsDefault,
	Title:   table.TitleOptionsDefault,
}

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get krb resources",
	// Run:   func(cmd *cobra.Command, args []string) {},
}

func init() {
	rootCmd.AddCommand(getCmd)
}
