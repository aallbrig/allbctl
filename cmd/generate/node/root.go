/*
Copyright © 2020 Andrew Allbright

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
package node

import (
	"github.com/aallbrig/allbctl/pkg"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "node",
	Short: "code generators for the node runtime",
	Run: func(cmd *cobra.Command, args []string) {
		pkg.HelpTextIfEmpty(cmd, args)
	},
}
