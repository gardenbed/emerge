package command

import (
	"errors"
	"io"
	"testing"

	"github.com/gardenbed/charm/ui"
	"github.com/stretchr/testify/assert"

	"github.com/gardenbed/emerge/internal/ebnf/parser/spec"
	"github.com/gardenbed/emerge/internal/generate"
)

func TestNew(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		u := ui.NewNop()
		cmd, err := New(u)

		assert.NotNil(t, cmd)
		assert.NoError(t, err)
	})
}

func TestCommand_PrintHelp(t *testing.T) {
	tests := []struct {
		name          string
		c             *Command
		expectedError string
	}{
		{
			name: "OK",
			c: &Command{
				UI:      ui.NewNop(),
				Verbose: true,
			},
			expectedError: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.c.PrintHelp()

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func TestCommand_Run(t *testing.T) {
	tests := []struct {
		name          string
		funcs         funcs
		args          []string
		expectedError string
	}{
		{
			name:          "Error_NoFile",
			funcs:         funcs{},
			args:          []string{},
			expectedError: "no input file specified, please provide a file path",
		},
		{
			name:  "Error_FileNotExist",
			funcs: funcs{},
			args: []string{
				"missing.grammar",
			},
			expectedError: "open missing.grammar: no such file or directory",
		},
		{
			name: "Error_ParseFails",
			funcs: funcs{
				Parse: func(string, io.Reader) (*spec.Spec, error) {
					return nil, errors.New("error on parsing the input")
				},
			},
			args: []string{
				"../ebnf/fixture/test.success.grammar",
			},
			expectedError: "error on parsing the input",
		},
		{
			name: "Error_GenerateFails",
			funcs: funcs{
				Parse: func(string, io.Reader) (*spec.Spec, error) {
					return &spec.Spec{}, nil
				},
				Generate: func(ui.UI, *generate.Params) error {
					return errors.New("error on generating the parser")
				},
			},
			args: []string{
				"../ebnf/fixture/test.success.grammar",
			},
			expectedError: "error on generating the parser",
		},
		{
			name: "Success",
			funcs: funcs{
				Parse: func(string, io.Reader) (*spec.Spec, error) {
					return &spec.Spec{}, nil
				},
				Generate: func(ui.UI, *generate.Params) error {
					return nil
				},
			},
			args: []string{
				"../ebnf/fixture/test.success.grammar",
			},
			expectedError: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c := &Command{
				UI:    ui.NewNop(),
				funcs: tc.funcs,
			}

			err := c.Run(tc.args)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}
