package jsonfile

import (
	"encoding/json"
	"github.com/wenerme/goform"
)

func init() {
	goform.RegisterForm(goform.FormRegistration{
		Extensions: []string{".json"},
		Read: func(input goform.FormInput) (form goform.Form, err error) {
			form = input.TargetForm
			if form == nil {
				form = make(map[string]interface{})
			}
			err = json.NewDecoder(input.Reader).Decode(&form)
			return
		},
		Write: func(form goform.Form, output goform.FormOutput) (err error) {
			_, indent := output.Options["indent"]

			encoder := json.NewEncoder(output.Writer)
			if indent {
				encoder.SetIndent("", "  ")
			}
			err = encoder.Encode(form)
			return
		},
	})
}
