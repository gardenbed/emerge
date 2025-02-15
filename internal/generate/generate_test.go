package generate

import (
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/gardenbed/charm/ui"
	"github.com/moorara/algo/parser/lr"
	"github.com/stretchr/testify/assert"

	"github.com/gardenbed/emerge/internal/ebnf/parser/spec"
)

func TestGenerate(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "emerge-test-")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

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
				Debug:       false,
				Path:        tempDir,
				Package:     "pascal",
				Definitions: []*spec.TerminalDef{},
				Grammar:     nil,
				Precedences: lr.PrecedenceLevels{},
			},
			expectedFiles: []string{
				"errors.go",
				"types.go",
				"stack.go",
				"input.go",
				"lexer.go",
				"parser.go",
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
				assert.True(t, re.MatchString(err.Error()), "%q DOES NOT MATCH %q", tc.expectedErrorRegex, err)
			} else {
				assert.NoError(t, err)

				for _, expectedFile := range tc.expectedFiles {
					filename := filepath.Join(tc.params.Path, tc.params.Package, expectedFile)
					assert.FileExists(t, filename)
				}
			}
		})
	}
}

func TestGenerator_prepare(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "emerge-test-")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	tempFile, err := os.CreateTemp("", "emerge-test-file-")
	assert.NoError(t, err)
	defer os.RemoveAll(tempFile.Name())

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
					Path:    tempDir,
					Package: "\x00",
				},
			},
			expectedErrorRegex: `invalid package name: \x00`,
		},
		{
			name: "PathReadOnly",
			g: &generator{
				UI: ui.NewNop(),
				Params: &Params{
					Path:    "/opt",
					Package: "pascal",
				},
			},
			expectedErrorRegex: `error on creating package directory: mkdir /opt/pascal: permission denied`,
		},
		{
			name: "Success",
			g: &generator{
				UI: ui.NewNop(),
				Params: &Params{
					Path:    tempDir,
					Package: "pascal",
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
				assert.True(t, re.MatchString(err.Error()), "%q DOES NOT MATCH %q", tc.expectedErrorRegex, err)
			}
		})
	}
}

func TestGenerator_generateCore(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "emerge-test-")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

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
					Path:    tempDir,
					Package: "pascal",
				},
			},
			expectedErrorRegexes: []string{
				`open .+/pascal/errors.go: no such file or directory`,
				`open .+/pascal/types.go: no such file or directory`,
				`open .+/pascal/stack.go: no such file or directory`,
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
					assert.True(t, re.MatchString(err.Error()), "%q DOES NOT MATCH %q", expectedErrorRegex, err)
				}
			}
		})
	}
}

func TestGenerator_generateLexer(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "emerge-test-")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

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
					Path:    tempDir,
					Package: "pascal",
				},
			},
			expectedErrorRegexes: []string{
				`open .+/pascal/input.go: no such file or directory`,
				`open .+/pascal/lexer.go: no such file or directory`,
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
					assert.True(t, re.MatchString(err.Error()), "%q DOES NOT MATCH %q", expectedErrorRegex, err)
				}
			}
		})
	}
}

func TestGenerator_generateParser(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "emerge-test-")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

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
					Path:    tempDir,
					Package: "pascal",
				},
			},
			expectedErrorRegexes: []string{
				`open .+/pascal/parser.go: no such file or directory`,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.g.generateParser()

			if len(tc.expectedErrorRegexes) == 0 {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)

				for _, expectedErrorRegex := range tc.expectedErrorRegexes {
					re := regexp.MustCompile(expectedErrorRegex)
					assert.True(t, re.MatchString(err.Error()), "%q DOES NOT MATCH %q", expectedErrorRegex, err)
				}
			}
		})
	}
}

func TestGenerator_renderTemplate(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "emerge-test-")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

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
			filename:           "missing.go",
			data:               nil,
			expectedErrorRegex: `open templates/missing.go.tmpl: file does not exist`,
		},
		{
			name: "PackageDirNotExist",
			g: &generator{
				UI: ui.NewNop(),
				Params: &Params{
					Path:    tempDir,
					Package: "pascal",
				},
			},
			filename:           "types.go",
			data:               nil,
			expectedErrorRegex: `open .+/pascal/types.go: no such file or directory`,
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
				assert.True(t, re.MatchString(err.Error()), "%q DOES NOT MATCH %q", tc.expectedErrorRegex, err)
			}
		})
	}
}
