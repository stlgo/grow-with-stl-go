/*
 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

     https://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	// version will be overridden by ldflags supplied in Makefile
	version = "(dev-version)"
)

func versionCmd() *cobra.Command {
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Show version",
		Long:  "Version for grow with stl-go binary",
		Run: func(cmd *cobra.Command, _ []string) {
			out := cmd.OutOrStdout()

			fmt.Fprintln(out, "grow with stl-go version", Version())
		},
	}
	return versionCmd
}

// Version returns application version
func Version() string {
	return version
}
