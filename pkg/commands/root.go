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
	"os"
	"os/signal"
	"path/filepath"
	"reflect"
	"runtime"
	"runtime/pprof"
	"syscall"
	"time"

	"github.com/spf13/cobra"

	"stl-go/grow-with-stl-go/pkg/admin"
	"stl-go/grow-with-stl-go/pkg/audit"
	"stl-go/grow-with-stl-go/pkg/configs"
	"stl-go/grow-with-stl-go/pkg/log"
	"stl-go/grow-with-stl-go/pkg/seeds"
	"stl-go/grow-with-stl-go/pkg/weather"
	"stl-go/grow-with-stl-go/pkg/webservice"
)

// rootCmd represents the base command when called without any subcommands
var (
	rootCmd = &cobra.Command{
		Use:     "grow-with-stl-go",
		Short:   "grow-with-stl-go is a sample go application developed by stl-go for demonstration purposes",
		Run:     launch,
		Version: Version(),
	}

	cpuProfile bool
	memProfile bool
)

func init() {
	// Add a 'version' command, in addition to the '--version' option that is auto created
	rootCmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Show version",
		Long:  "Version for grow with stl-go binary",
		Run: func(cmd *cobra.Command, _ []string) {
			out := cmd.OutOrStdout()

			fmt.Fprintln(out, "grow with stl-go version", Version())
		},
	})

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

	rootCmd.Flags().BoolVarP(
		&cpuProfile,
		"cpuProfile",
		"C",
		false,
		"This will enable pprof CPU profile",
	)

	rootCmd.Flags().BoolVarP(
		&memProfile,
		"memProfile",
		"m",
		false,
		"This will enable pprof memory profile",
	)

	// Add the logging level flag
	rootCmd.Flags().IntVar(
		&log.LogLevel,
		"loglevel",
		4,
		"This will set the log level, anything at or below that level will be viewed, all others suppressed\n"+
			"  6 -- Trace\n"+
			"  5 -- Debug\n"+
			"  4 -- Info\n"+
			"  3 -- Warn\n"+
			"  2 -- Error\n"+
			"  1 -- Fatal\n",
	)
}

func setupProfile() error {
	now := time.Now().UnixMilli()
	if cpuProfile && configs.GrowSTLGo != nil && configs.GrowSTLGo.DataDir != nil {
		cpuProfileFile := filepath.Join(*configs.GrowSTLGo.DataDir, fmt.Sprintf("cpuProfile.%d", now))
		f, err := os.Create(cpuProfileFile)
		if err != nil {
			return fmt.Errorf("could not create CPU profile: %s.  Error: %s", cpuProfileFile, err)
		}
		defer f.Close()
		if err := pprof.StartCPUProfile(f); err != nil {
			return fmt.Errorf("could not start CPU profile: %s", err)
		}
		defer pprof.StopCPUProfile()
		log.Infof("cpu profiling sent to file: %s", cpuProfileFile)
	}

	if memProfile && configs.GrowSTLGo != nil && configs.GrowSTLGo.DataDir != nil {
		memProfileFile := filepath.Join(*configs.GrowSTLGo.DataDir, fmt.Sprintf("memProfile.%d", now))
		f, err := os.Create(memProfileFile)
		if err != nil {
			return fmt.Errorf("could not create memory profile: %s.  Error: %s", memProfileFile, err)
		}
		runtime.GC() // get up-to-date statistics
		// Lookup("allocs") creates a profile similar to go test -memprofile.
		// Alternatively, use Lookup("heap") for a profile
		// that has inuse_space as the default index.
		if err := pprof.Lookup("allocs").WriteTo(f, 0); err != nil {
			return fmt.Errorf("could not write memory profile: %s", err)
		}
		defer f.Close()
		log.Infof("memory profiling sent to file: %s", memProfileFile)
	}
	return nil
}

func launch(_ *cobra.Command, _ []string) {
	// read the config file
	if err := configs.SetGrowSTLGoConfig(); err != nil {
		log.Fatalf("error in config %s: %s", *configs.ConfigFile, err)
	}

	if err := setupProfile(); err != nil {
		log.Fatalf("error setting up profiling: %s", err)
	}

	// kick off the init functions for the various packages
	for _, function := range []func() error{audit.Init, seeds.Init, admin.Init, weather.Init} {
		if err := function(); err != nil {
			log.Fatalf("error calling function %s cannot continue to start", runtime.FuncForPC(reflect.ValueOf(function).Pointer()).Name())
		}
	}

	// start webservice and listen for the the ctl + c to exit
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Info("Closing all active sessions")
		webservice.Shutdown()
		configs.ShutdownSQLite()
		log.Info("Exiting the webservice")
		os.Exit(0)
	}()

	webservice.WebServer()
}

// Version returns the version number for the cobra command
func Version() string {
	return configs.Version
}

// Execute is called from the main program and kicks this whole shindig off
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Error(err)
		os.Exit(1)
	}
}
