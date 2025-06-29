package command

import (
	"errors"
	"io"
	"os"
	"testing"

	"github.com/gardenbed/charm/ui"
	"github.com/stretchr/testify/assert"

	"github.com/gardenbed/emerge/internal/ebnf/parser/spec"
	"github.com/gardenbed/emerge/internal/generate/golang"
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
			r, w, err := os.Pipe()
			assert.NoError(t, err)

			orig := os.Stdout
			os.Stdout = w
			defer func() {
				os.Stdout = orig
			}()

			err = tc.c.PrintHelp()

			if tc.expectedError != "" {
				assert.EqualError(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
				assert.NoError(t, w.Close())

				out, err := io.ReadAll(r)
				assert.NotNil(t, out)
				assert.NoError(t, err)
			}
		})
	}
}

func TestCommand_Run(t *testing.T) {
	tests := []struct {
		name                 string
		c                    *Command
		args                 []string
		expectedErrorStrings []string
	}{
		{
			name: "Error_NoFile",
			c: &Command{
				UI:    ui.NewNop(),
				funcs: funcs{},
			},
			args: []string{},
			expectedErrorStrings: []string{
				`no input file specified, please provide a file path`,
			},
		},
		{
			name: "Error_FileNotExist",
			c: &Command{
				UI:    ui.NewNop(),
				funcs: funcs{},
			},
			args: []string{
				"missing.grammar",
			},
			expectedErrorStrings: []string{
				`open missing.grammar: no such file or directory`,
			},
		},
		{
			name: "Error_ParseFails",
			c: &Command{
				UI: ui.NewNop(),
				funcs: funcs{
					Parse: func(string, io.Reader) (*spec.Spec, error) {
						return nil, errors.New("error on parsing the input")
					},
				},
			},
			args: []string{
				"../ebnf/fixture/test.success.grammar",
			},
			expectedErrorStrings: []string{
				`error on parsing the input`,
			},
		},
		{
			name: "Error_GenerateFails",
			c: &Command{
				UI: ui.NewNop(),
				funcs: funcs{
					Parse: func(string, io.Reader) (*spec.Spec, error) {
						return &spec.Spec{}, nil
					},
					Generate: func(ui.UI, *golang.Params) error {
						return errors.New("error on generating the parser")
					},
				},
			},
			args: []string{
				"../ebnf/fixture/test.success.grammar",
			},
			expectedErrorStrings: []string{
				`error on generating the parser`,
			},
		},
		{
			name: "Success",
			c: &Command{
				UI: ui.NewNop(),
				funcs: funcs{
					Parse: func(string, io.Reader) (*spec.Spec, error) {
						return &spec.Spec{}, nil
					},
					Generate: func(ui.UI, *golang.Params) error {
						return nil
					},
				},

				Out:   "/path/to/destination",
				Name:  "override",
				Debug: false,
			},
			args: []string{
				"../ebnf/fixture/test.success.grammar",
			},
			expectedErrorStrings: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.c.Run(tc.args)

			if len(tc.expectedErrorStrings) == 0 {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				s := err.Error()
				for _, expectedErrorString := range tc.expectedErrorStrings {
					assert.Contains(t, s, expectedErrorString)
				}
			}
		})
	}
}
