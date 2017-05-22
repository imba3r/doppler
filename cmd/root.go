// Copyright © 2017 imba3r
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"fmt"
	"os"

	"github.com/imba3r/doppler/doppler"
	"github.com/spf13/cobra"
)

var (
	skipName bool
	skipHash bool

	outputJSON bool
	absolute   bool
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "doppler dir [dirN]",
	Short: "Doppler locates duplicate files by name and/or hash.",
	Long: `Doppler locates duplicate files by name and/or hash.

Directories are searched recursively. If multiple directories are specified
the results of the search spans all of them. Symlinks are not followed!`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Need at least one directory.")
			return
		}
		if skipHash && skipName {
			fmt.Println("Nothing to do if both name and hash check are skipped.")
			return
		}
		s := doppler.NewScanner(skipName, skipHash, absolute)
		s.ScanDirs(args)
		if s.HasErrors() {
			for f, err := range s.ErrMap {
				os.Stderr.WriteString(f + ": " + err)
			}
		}
		if !s.FoundDuplicates() {
			fmt.Println("No duplicates found!")
			return
		}
		if outputJSON {
			s.PrintJSON()
		} else {
			s.Print()
		}
	},
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	RootCmd.Flags().BoolVarP(&outputJSON, "json", "j", false, "Print result in JSON format (more details).")
	RootCmd.Flags().BoolVarP(&absolute, "absolute", "a", false, "Print absolute file paths.")
	RootCmd.Flags().BoolVarP(&skipName, "skip-name", "", false, "Skip check for duplicate file names.")
	RootCmd.Flags().BoolVarP(&skipHash, "skip-hash", "", false, "Skip check for duplicate file content.")
}
