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

package log_test

import (
	"bytes"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"stl-go/grow-with-stl-go/pkg/log"
)

var logFormatRegex = regexp.MustCompile(`^\[stl-go\] .*`)

func TestLoggingTrace(t *testing.T) {
	testRequire := require.New(t)

	t.Run("TraceViewable", func(t *testing.T) {
		output := new(bytes.Buffer)
		log.Init("TRACE", output)

		log.Debug("TraceViewable args ", 5)
		actual := output.String()

		expected := "TraceViewable args 5"
		testRequire.Regexp(logFormatRegex, actual)
		actualArray := strings.Split(actual, "]")
		actual = strings.TrimSpace(actualArray[len(actualArray)-1])
		require.Equal(t, expected, actual)
	})

	t.Run("TracefViewable", func(t *testing.T) {
		output := new(bytes.Buffer)
		log.Init("TRACE", output)

		log.Debugf("%s %d", "TracefViewable args", 5)
		actual := output.String()

		expected := "TracefViewable args 5"
		testRequire.Regexp(logFormatRegex, actual)
		actualArray := strings.Split(actual, "]")
		actual = strings.TrimSpace(actualArray[len(actualArray)-1])
		require.Equal(t, expected, actual)
	})

	t.Run("TraceNotViewable", func(t *testing.T) {
		output := new(bytes.Buffer)
		log.Init("FATAL", output)

		log.Debug("TraceNotViewable args ", 5)
		require.Equal(t, "", output.String())
	})

	t.Run("TracefNotViewable", func(t *testing.T) {
		output := new(bytes.Buffer)
		log.Init("FATAL", output)

		log.Debugf("%s %d", "TracefNotViewable args", 5)
		require.Equal(t, "", output.String())
	})
}

func TestLoggingDebug(t *testing.T) {
	testRequire := require.New(t)

	t.Run("DebugViewable", func(t *testing.T) {
		output := new(bytes.Buffer)
		log.Init("DEBUG", output)

		log.Debug("DebugViewable args ", 5)
		actual := output.String()

		expected := "DebugViewable args 5"
		testRequire.Regexp(logFormatRegex, actual)
		actualArray := strings.Split(actual, "]")
		actual = strings.TrimSpace(actualArray[len(actualArray)-1])
		require.Equal(t, expected, actual)
	})

	t.Run("DebugfViewable", func(t *testing.T) {
		output := new(bytes.Buffer)
		log.Init("DEBUG", output)

		log.Debugf("%s %d", "DebugfViewable args", 5)
		actual := output.String()

		expected := "DebugfViewable args 5"
		testRequire.Regexp(logFormatRegex, actual)
		actualArray := strings.Split(actual, "]")
		actual = strings.TrimSpace(actualArray[len(actualArray)-1])
		require.Equal(t, expected, actual)
	})

	t.Run("DebugNotViewable", func(t *testing.T) {
		output := new(bytes.Buffer)
		log.Init("FATAL", output)

		log.Debug("DebugNotViewable args ", 5)
		require.Equal(t, "", output.String())
	})

	t.Run("DebugfNotViewable", func(t *testing.T) {
		output := new(bytes.Buffer)
		log.Init("FATAL", output)

		log.Debugf("%s %d", "DebugfNotViewable args", 5)
		require.Equal(t, "", output.String())
	})
}

func TestLoggingInfo(t *testing.T) {
	testRequire := require.New(t)

	t.Run("InfoViewable", func(t *testing.T) {
		output := new(bytes.Buffer)
		log.Init("INFO", output)

		log.Info("InfoViewable args ", 5)
		actual := output.String()

		expected := "InfoViewable args 5"
		testRequire.Regexp(logFormatRegex, actual)
		actualArray := strings.Split(actual, "]")
		actual = strings.TrimSpace(actualArray[len(actualArray)-1])
		require.Equal(t, expected, actual)
	})

	t.Run("InfofViewable", func(t *testing.T) {
		output := new(bytes.Buffer)
		log.Init("INFO", output)

		log.Infof("%s %d", "InfofViewable args", 5)
		actual := output.String()

		expected := "InfofViewable args 5"
		testRequire.Regexp(logFormatRegex, actual)
		actualArray := strings.Split(actual, "]")
		actual = strings.TrimSpace(actualArray[len(actualArray)-1])
		require.Equal(t, expected, actual)
	})

	t.Run("InfoNotViewable", func(t *testing.T) {
		output := new(bytes.Buffer)
		log.Init("FATAL", output)

		log.Info("InfoNotViewable args ", 5)
		require.Equal(t, "", output.String())
	})

	t.Run("InfofNotViewable", func(t *testing.T) {
		output := new(bytes.Buffer)
		log.Init("FATAL", output)

		log.Infof("%s %d", "InfofNotViewable args", 5)
		require.Equal(t, "", output.String())
	})
}

func TestLoggingWarn(t *testing.T) {
	testRequire := require.New(t)

	t.Run("WarnViewable", func(t *testing.T) {
		output := new(bytes.Buffer)
		log.Init("WARN", output)

		log.Warn("WarnViewable args ", 5)
		actual := output.String()

		expected := "WarnViewable args 5"
		testRequire.Regexp(logFormatRegex, actual)
		actualArray := strings.Split(actual, "]")
		actual = strings.TrimSpace(actualArray[len(actualArray)-1])
		require.Equal(t, expected, actual)
	})

	t.Run("WarnfViewable", func(t *testing.T) {
		output := new(bytes.Buffer)
		log.Init("WARN", output)

		log.Warnf("%s %d", "WarnfViewable args", 5)
		actual := output.String()

		expected := "WarnfViewable args 5"
		testRequire.Regexp(logFormatRegex, actual)
		actualArray := strings.Split(actual, "]")
		actual = strings.TrimSpace(actualArray[len(actualArray)-1])
		require.Equal(t, expected, actual)
	})

	t.Run("WarnNotViewable", func(t *testing.T) {
		output := new(bytes.Buffer)
		log.Init("FATAL", output)

		log.Warn("WarnNotViewable args ", 5)
		require.Equal(t, "", output.String())
	})

	t.Run("WarnfNotViewable", func(t *testing.T) {
		output := new(bytes.Buffer)
		log.Init("FATAL", output)

		log.Warnf("%s %d", "WarnfNotViewable args", 5)
		require.Equal(t, "", output.String())
	})
}

func TestLoggingError(t *testing.T) {
	testRequire := require.New(t)

	t.Run("ErrorViewable", func(t *testing.T) {
		output := new(bytes.Buffer)
		log.Init("ERROR", output)

		log.Error("ErrorViewable args ", 5)
		actual := output.String()

		expected := "ErrorViewable args 5"
		testRequire.Regexp(logFormatRegex, actual)
		actualArray := strings.Split(actual, "]")
		actual = strings.TrimSpace(actualArray[len(actualArray)-1])
		require.Equal(t, expected, actual)
	})

	t.Run("ErrorfViewable", func(t *testing.T) {
		output := new(bytes.Buffer)
		log.Init("ERROR", output)

		log.Errorf("%s %d", "ErrorfViewable args", 5)
		actual := output.String()

		expected := "ErrorfViewable args 5"
		testRequire.Regexp(logFormatRegex, actual)
		actualArray := strings.Split(actual, "]")
		actual = strings.TrimSpace(actualArray[len(actualArray)-1])
		require.Equal(t, expected, actual)
	})

	t.Run("ErrorNotViewable", func(t *testing.T) {
		output := new(bytes.Buffer)
		log.Init("FATAL", output)

		log.Warn("ErrorNotViewable args ", 5)
		require.Equal(t, "", output.String())
	})

	t.Run("ErrorfNotViewable", func(t *testing.T) {
		output := new(bytes.Buffer)
		log.Init("FATAL", output)

		log.Warnf("%s %d", "ErrorfNotViewable args", 5)
		require.Equal(t, "", output.String())
	})
}
