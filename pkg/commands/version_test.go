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
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

var testRootCmd = &cobra.Command{
	Use:     "grow-with-stl-go",
	Short:   "grow-with-stl-go is a sample go application developed by stl-go",
	Run:     launch,
	Version: Version(),
}

func TestVersionExecuteCommand(t *testing.T) {
	// Add a 'versioin' command
	rootCmd.AddCommand(versionCmd())

	cmd := testRootCmd
	b := bytes.NewBufferString("")
	cmd.SetOut(b)
	cmd.SetArgs([]string{"-v"})

	err := cmd.Execute()
	if err != nil {
		t.Fatal(err)
	}

	out, err := io.ReadAll(b)
	if err != nil {
		t.Fatal(err)
	}

	if strings.TrimSpace(string(out)) != "grow-with-stl-go version (dev-version)" {
		t.Fatalf("expected \"%s\" gpt \"%s\"", "grow-with-stl-go version (dev-version)", string(out))
	}
}
