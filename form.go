package goform

import (
	"github.com/Sirupsen/logrus"
	"io"
	"sort"
)

//type FormDefinition struct {
//	PrimaryExtension string
//	Extensions       []string
//}
//
//type FormReader interface {
//	ReadForm() (Form, error)
//}
//type FormWriter interface {
//	WriteForm() error
//}

type Form interface {
}

type FormInput struct {
	Filename  string
	Extension string

	Reader     io.Reader
	Options    map[string]string
	TargetForm Form
}

type FormOutput struct {
	Filename  string
	Extension string

	Writer  io.Writer
	Options map[string]string
}

type FormRegistration struct {
	Extensions []string
	Priority   int

	CanRead  func(input FormInput) bool
	CanWrite func(form Form, output FormOutput) bool

	Read  func(input FormInput) (form Form, err error)
	Write func(form Form, output FormOutput) (err error)

	NewForm func() Form
}

var registrations []FormRegistration

func RegisterForm(r FormRegistration) {
	if r.Write != nil && r.CanWrite == nil {
		r.CanWrite = func(form Form, output FormOutput) bool {
			for _, v := range r.Extensions {
				if v == output.Extension {
					return true
				}
			}
			return false
		}
	}
	if r.Read != nil && r.CanRead == nil {
		r.CanRead = func(input FormInput) bool {
			for _, v := range r.Extensions {
				if v == input.Extension {
					return true
				}
			}
			return false
		}
	}
	registrations = append(registrations, r)

	sort.Slice(registrations, func(i, j int) bool {
		return registrations[i].Priority > registrations[j].Priority
	})

	logrus.Debug("register form with extension ", r.Extensions)
}

func RegisteredFormExtension() []string {
	all := make([]string, 0)
	for _, r := range registrations {
		all = append(all, r.Extensions...)
	}
	return all
}

func ReadForm(file FormInput) (form Form, err error) {
	for _, r := range registrations {
		if r.CanRead != nil && r.CanRead(file) {
			return r.Read(file)
		}
	}
	return nil, EUnsupportedFileFormat
}

func WriteForm(form Form, file FormOutput) (err error) {
	for _, r := range registrations {
		if r.CanWrite != nil && r.CanWrite(form, file) {
			return r.Write(form, file)
		}
	}
	return EUnsupportedFileFormat
}

func NewFormWithExtension(ext string) (form Form) {
	for _, r := range registrations {
		for _, v := range r.Extensions {
			if v == ext {
				if r.NewForm != nil {
					form = r.NewForm()
				}
				if form != nil {
					return
				}
			}
		}
	}
	return
}
