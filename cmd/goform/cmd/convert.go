// Copyright © 2018 wener <wenermail@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/wenerme/goform"
	_ "github.com/wenerme/goform/jsonfile"
	_ "github.com/wenerme/goform/metfile"
	"github.com/wenerme/letsgo/fs"
	"io"
	"os"
	"strings"
)

// convertCmd represents the convert command
var convertCfg = _convertCfg{}

type _convertCfg struct {
	to   string
	from string

	inputOptions  []string
	outputOptions []string

	inputOptionMap  map[string]string
	outputOptionMap map[string]string
}

func (cfg *_convertCfg) PreCheck() error {
	if !wfs.Exists(cfg.from) {
		return errors.Errorf("file not exists %v", convertCfg.from)
	}
	cfg.inputOptionMap = make(map[string]string)
	cfg.outputOptionMap = make(map[string]string)

	parseOption(cfg.inputOptions, cfg.inputOptionMap)
	parseOption(cfg.outputOptions, cfg.outputOptionMap)
	return nil
}

func parseOption(o []string, m map[string]string) {
	for _, v := range o {
		key := v
		value := ""
		if strings.ContainsRune(v, '=') {
			split := strings.Split(v, "=")
			key = strings.TrimSpace(split[0])
			value = strings.TrimSpace(split[1])
		}
		m[key] = value
	}
}

var convertCmd = &cobra.Command{
	Use:     "convert",
	Aliases: []string{"c"},
	Short:   "Convert between different form",
	Long:    ``,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if err = convertCfg.PreCheck(); err != nil {
			return
		}

		fromExt, fromBase := wfs.Ext(convertCfg.from)
		toExt := convertCfg.to
		toBase := fromBase
		toFn := toBase + toExt
		var form goform.Form

		if strings.ContainsRune(convertCfg.to, '.') {
			toExt, toBase = wfs.Ext(convertCfg.to)
			toFn = convertCfg.to
		}
		inFile, err := os.Open(convertCfg.from)
		if err != nil {
			return err
		}
		// input
		input := goform.FormInput{
			Extension:  fromExt,
			Options:    convertCfg.inputOptionMap,
			Reader:     inFile,
			TargetForm: goform.NewFormWithExtension(toExt),
		}
		// output
		output := goform.FormOutput{
			Extension: toExt,
			Options:   convertCfg.outputOptionMap,
		}
		if toBase == "-" {
			output.Writer = os.Stdout
		} else {
			output.Writer, err = os.OpenFile(toFn, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, os.ModePerm)
			if err != nil {
				return err
			}
			defer func() {
				if c, ok := output.Writer.(io.Closer); ok {
					c.Close()
				}
			}()
		}

		// do read
		if form, err = goform.ReadForm(input); err != nil {
			if err == goform.EUnsupportedFileFormat {
				return errors.Errorf("unsupported input format: %v", fromExt)
			}
			return
		}

		// do write
		if err = goform.WriteForm(form, output); err != nil {
			if err == goform.EUnsupportedFileFormat {
				return errors.Errorf("unsupported output format: %v", toExt)
			}
			return
		}

		return err
	},
}

func init() {
	rootCmd.AddCommand(convertCmd)

	convertCmd.Flags().StringVarP(&convertCfg.from, "input", "i", "", "Input file")
	convertCmd.Flags().StringVarP(&convertCfg.to, "output", "o", "-", "Output file")

	convertCmd.Flags().StringSliceVarP(&convertCfg.inputOptions, "input-option", "I", nil, "Input options")
	convertCmd.Flags().StringSliceVarP(&convertCfg.outputOptions, "output-option", "O", nil, "Output options")
}
