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

package main

import (
	"fmt"
	"os"
	"strings"
	"unicode"

	"stl-go/grow-with-stl-go/pkg/log"

	"github.com/spf13/cobra"
)

var (
	shift   int
	input   string
	decrypt = false
	rootCmd = &cobra.Command{
		Use:     "grow-with-stl-go",
		Short:   "grow-with-stl-go is a sample go application developed by stl-go for demonstration purposes, this is its REST example client",
		Run:     runCeaserCipher,
		Version: version(),
	}
)

func init() {
	// Add a 'version' command, in addition to the '--version' option that is auto created
	rootCmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Show version",
		Long:  "Version for grow with stl-go Ceaser Cipher binary",
		Run: func(cmd *cobra.Command, _ []string) {
			out := cmd.OutOrStdout()
			fmt.Fprintln(out, "grow with stl-go Ceaser Cipher version", version())
		},
	})

	// Add the shift
	rootCmd.Flags().IntVar(
		&shift,
		"shift",
		4,
		"The shift used in the ceaser cipher",
	)

	// Add the string to encipher / decipher
	rootCmd.Flags().StringVarP(
		&input,
		"input",
		"i",
		"some text here",
		"The string to encipher/decipher",
	)

	rootCmd.Flags().BoolVarP(
		&decrypt,
		"decipher",
		"d",
		false,
		"decipher the input",
	)
}

// Version returns application version
func version() string {
	// pull the file version if it's available
	if fileVersion, err := os.ReadFile("version"); err == nil {
		return string(fileVersion)
	}
	return "dev-version"
}

func runCeaserCipher(_ *cobra.Command, _ []string) {
	log.Infof("Enciphering: '%s'", input)
	cipherText := encipher(input, shift)
	log.Infof("Ciphertext: '%s'", cipherText)
	log.Infof("Deciphered text: '%s'", decipher(cipherText, shift))
}

func encipher(text string, shift int) string {
	var ciphertext strings.Builder
	for _, char := range text {
		if unicode.IsLower(char) && char >= 'a' && char <= 'z' {
			ciphertext.WriteRune((char-'a'+rune(shift))%26 + 'a')
			continue
		}
		if unicode.IsUpper(char) && char >= 'A' && char <= 'A' {
			ciphertext.WriteRune((char-'A'+rune(shift))%26 + 'A')
			continue
		}
		ciphertext.WriteRune(char)
	}
	return ciphertext.String()
}

// Function to perform Caesar Cipher decryption
func decipher(ciphertext string, shift int) string {
	return encipher(ciphertext, 26-shift) // Decryption is just the reverse shift
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Error(err)
		os.Exit(1)
	}
}
