package golang

import (
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gardenbed/charm/ui"

	"github.com/gardenbed/emerge/internal/ebnf/parser/spec"
)

func TestIsIDValid(t *testing.T) {
	tests := []struct {
		name         string
		expectedBool bool
	}{
		{
			name:         "4ever",
			expectedBool: false,
		},
		{
			name:         "complex",
			expectedBool: false,
		},
		{
			name:         "pascal",
			expectedBool: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expectedBool, isIDValid(tc.name))
		})
	}
}

func TestGenerate(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "emerge-test-")
	assert.NoError(t, err)

	defer func() {
		assert.NoError(t, os.RemoveAll(tempDir))
	}()

	tests := []struct {
		name               string
		params             *Params
		expectedFiles      []string
		expectedErrorRegex string
	}{
		{
			name: "PathNotExist",
			params: &Params{
				Debug: false,
				Path:  filepath.Join(tempDir, "missing"),
			},
			expectedErrorRegex: `output path does not exist: ".+/emerge-test-.+/missing"`,
		},
		{
			name: "Success",
			params: &Params{
				Debug: false,
				Path:  tempDir,
				Spec: &spec.Spec{
					Name:        "foo",
					Definitions: definitions,
					Grammar:     grammars[0],
					Precedences: precedences[0],
				},
			},
			expectedFiles: []string{
				"foo.go",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			u := ui.NewNop()
			err := Generate(u, tc.params)

			if tc.expectedErrorRegex != "" {
				assert.Error(t, err)

				re := regexp.MustCompile(tc.expectedErrorRegex)
				assert.True(t, re.MatchString(err.Error()), "%q DOES NOT INCLUDE %q", err, tc.expectedErrorRegex)
			} else {
				assert.NoError(t, err)

				for _, expectedFile := range tc.expectedFiles {
					filename := filepath.Join(tc.params.Path, tc.params.Spec.Name, expectedFile)
					assert.FileExists(t, filename)
				}
			}
		})
	}
}

func TestGenerator_prepare(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "emerge-test-")
	assert.NoError(t, err)

	defer func() {
		assert.NoError(t, os.RemoveAll(tempDir))
	}()

	tempFile, err := os.CreateTemp("", "emerge-test-file-")
	assert.NoError(t, err)

	defer func() {
		assert.NoError(t, os.RemoveAll(tempFile.Name()))
	}()

	tests := []struct {
		name               string
		g                  *generator
		expectedErrorRegex string
	}{
		{
			name: "PathNotExist",
			g: &generator{
				UI: ui.NewNop(),
				Params: &Params{
					Path: filepath.Join(tempDir, "missing"),
				},
			},
			expectedErrorRegex: `output path does not exist: ".+/emerge-test-.+/missing"`,
		},
		{
			name: "PathNotDirectory",
			g: &generator{
				UI: ui.NewNop(),
				Params: &Params{
					Path: tempFile.Name(),
				},
			},
			expectedErrorRegex: `output path is not a directory: ".+/emerge-test-file-.+"`,
		},
		{
			name: "InvalidPackage",
			g: &generator{
				UI: ui.NewNop(),
				Params: &Params{
					Path: tempDir,
					Spec: &spec.Spec{
						Name: "\x00",
					},
				},
			},
			expectedErrorRegex: `invalid package name: \x00`,
		},
		{
			name: "Success",
			g: &generator{
				UI: ui.NewNop(),
				Params: &Params{
					Path: tempDir,
					Spec: &spec.Spec{
						Name: "foo",
					},
				},
			},
			expectedErrorRegex: ``,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.g.prepare()

			if tc.expectedErrorRegex == "" {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)

				re := regexp.MustCompile(tc.expectedErrorRegex)
				assert.True(t, re.MatchString(err.Error()), "%q DOES NOT INCLUDE %q", err, tc.expectedErrorRegex)
			}
		})
	}
}

func TestGenerator_generateCore(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "emerge-test-")
	assert.NoError(t, err)

	defer func() {
		assert.NoError(t, os.RemoveAll(tempDir))
	}()

	tests := []struct {
		name                 string
		g                    *generator
		expectedErrorRegexes []string
	}{
		{
			name: "PackageDirNotExist",
			g: &generator{
				UI: ui.NewNop(),
				Params: &Params{
					Path: tempDir,
					Spec: &spec.Spec{
						Name: "foo",
					},
				},
			},
			expectedErrorRegexes: []string{
				`open .+/foo/foo.go: no such file or directory`,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.g.generateCore()

			if len(tc.expectedErrorRegexes) == 0 {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)

				for _, expectedErrorRegex := range tc.expectedErrorRegexes {
					re := regexp.MustCompile(expectedErrorRegex)
					assert.True(t, re.MatchString(err.Error()), "%q DOES NOT INCLUDE %q", err, expectedErrorRegex)
				}
			}
		})
	}
}

func TestGenerator_generateLexer(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "emerge-test-")
	assert.NoError(t, err)

	defer func() {
		assert.NoError(t, os.RemoveAll(tempDir))
	}()

	tests := []struct {
		name                 string
		g                    *generator
		expectedErrorRegexes []string
	}{
		{
			name: "InvalidRegexDefinitions",
			g: &generator{
				UI: ui.NewNop(),
				Params: &Params{
					Path: tempDir,
					Spec: &spec.Spec{
						Name: "expr",
						Definitions: []*spec.TerminalDef{
							{Terminal: "ID", Kind: spec.RegexDef, Value: "[A-Z"},
							{Terminal: "NUM", Kind: spec.RegexDef, Value: "[0-9"},
						},
					},
				},
			},
			expectedErrorRegexes: []string{
				`2 errors occurred:`,
				`"ID": invalid regular expression: \[A-Z`,
				`"NUM": invalid regular expression: \[0-9`,
			},
		},
		{
			name: "PackageDirNotExist",
			g: &generator{
				UI: ui.NewNop(),
				Params: &Params{
					Path: tempDir,
					Spec: &spec.Spec{
						Name:        "expr",
						Definitions: definitions,
					},
				},
			},
			expectedErrorRegexes: []string{
				`open .+/expr/expr.go: no such file or directory`,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.g.generateLexer()

			if len(tc.expectedErrorRegexes) == 0 {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)

				for _, expectedErrorRegex := range tc.expectedErrorRegexes {
					re := regexp.MustCompile(expectedErrorRegex)
					assert.True(t, re.MatchString(err.Error()), "%q DOES NOT INCLUDE %q", err, expectedErrorRegex)
				}
			}
		})
	}
}

func TestGenerator_renderTemplate(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "emerge-test-")
	assert.NoError(t, err)

	defer func() {
		assert.NoError(t, os.RemoveAll(tempDir))
	}()

	tests := []struct {
		name               string
		g                  *generator
		filename           string
		data               any
		expectedErrorRegex string
	}{
		{
			name: "InvalidFilename",
			g: &generator{
				UI: ui.NewNop(),
			},
			filename:           "missing.go.tmpl",
			data:               nil,
			expectedErrorRegex: `open templates/missing.go.tmpl: file does not exist`,
		},
		{
			name: "PackageDirNotExist",
			g: &generator{
				UI: ui.NewNop(),
				Params: &Params{
					Path: tempDir,
					Spec: &spec.Spec{
						Name: "foo",
					},
				},
			},
			filename:           "core.go.tmpl",
			data:               nil,
			expectedErrorRegex: `open .+/foo/foo.go: no such file or directory`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.g.renderTemplate(tc.filename, tc.data)

			if tc.expectedErrorRegex == "" {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)

				re := regexp.MustCompile(tc.expectedErrorRegex)
				assert.True(t, re.MatchString(err.Error()), "%q DOES NOT INCLUDE %q", err, tc.expectedErrorRegex)
			}
		})
	}
}
