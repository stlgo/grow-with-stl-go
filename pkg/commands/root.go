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
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"stl-go/grow-with-stl-go/pkg/configs"
	"stl-go/grow-with-stl-go/pkg/log"
	"stl-go/grow-with-stl-go/pkg/webservice"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "grow-with-stl-go",
	Short:   "grow-with-stl-go is a sample go application developed by stl-go for demonstration purposes",
	Run:     launch,
	Version: Version(),
}

func init() {
	// Add a 'version' command, in addition to the '--version' option that is auto created
	rootCmd.AddCommand(versionCmd())

	// Add the config file Flag
	configFile := "../../etc/grow-with-stl-go.json"
	configs.ConfigFile = &configFile
	rootCmd.Flags().StringVarP(
		configs.ConfigFile,
		"conf",
		"c",
		"etc/grow-with-stl-go.json",
		"This will set the location of the conf file needed to start the UI",
	)

	// Add the logging level flag
	rootCmd.Flags().IntVar(
		&log.LogLevel,
		"loglevel",
		6,
		"This will set the log level, anything at or below that level will be viewed, all others suppressed\n"+
			"  6 -- Trace\n"+
			"  5 -- Debug\n"+
			"  4 -- Info\n"+
			"  3 -- Warn\n"+
			"  2 -- Error\n"+
			"  1 -- Fatal\n",
	)
}

func launch(_ *cobra.Command, _ []string) {
	// read the config file
	if err := configs.SetGrowSTLGoConfig(); err != nil {
		log.Fatalf("error in config %s: %s", *configs.ConfigFile, err)
	}

	// start webservice and listen for the the ctl + c to exit
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Info("Closing all active sessions")
		webservice.Shutdown()
		log.Info("Exiting the webservice")
		os.Exit(0)
	}()

	webservice.WebServer()
}

// Execute is called from the main program and kicks this whole shindig off
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Error(err)
		os.Exit(1)
	}
}
