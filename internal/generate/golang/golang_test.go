package golang

import (
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gardenbed/charm/ui"
	"github.com/moorara/algo/automata"
	"github.com/moorara/algo/parser/lr"
	"github.com/moorara/algo/parser/lr/lookahead"

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
					Debug: false,
					Path:  filepath.Join(tempDir, "missing"),
				},
			},
			expectedErrorRegex: `output path does not exist: ".+/emerge-test-.+/missing"`,
		},
		{
			name: "PathNotDirectory",
			g: &generator{
				UI: ui.NewNop(),
				Params: &Params{
					Debug: false,
					Path:  tempFile.Name(),
				},
			},
			expectedErrorRegex: `output path is not a directory: ".+/emerge-test-file-.+"`,
		},
		{
			name: "InvalidPackage",
			g: &generator{
				UI: ui.NewNop(),
				Params: &Params{
					Debug: false,
					Path:  tempDir,
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
					Debug: false,
					Path:  tempDir,
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
					Debug: false,
					Path:  tempDir,
					Spec: &spec.Spec{
						Name: "foo",
					},
				},
			},
			expectedErrorRegexes: []string{
				`open .+/foo/foo.go: no such file or directory`,
			},
		},
		{
			name: "Success",
			g: &generator{
				UI: ui.NewNop(),
				Params: &Params{
					Debug: false,
					Path:  tempDir,
					Spec: &spec.Spec{
						Name: "",
					},
				},
			},
			expectedErrorRegexes: nil,
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
					Debug: false,
					Path:  tempDir,
					Spec: &spec.Spec{
						Name: "foo",
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
					Debug: false,
					Path:  tempDir,
					Spec: &spec.Spec{
						Name:        "foo",
						Definitions: definitions,
					},
				},
			},
			expectedErrorRegexes: []string{
				`open .+/foo/foo.go: no such file or directory`,
			},
		},
		{
			name: "Success",
			g: &generator{
				UI: ui.NewNop(),
				Params: &Params{
					Debug: false,
					Path:  tempDir,
					Spec: &spec.Spec{
						Name:        "",
						Definitions: definitions,
					},
				},
			},
			expectedErrorRegexes: nil,
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

func TestGenerator_generateParser(t *testing.T) {
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
			name: "ParsingTableFails",
			g: &generator{
				UI: ui.NewNop(),
				Params: &Params{
					Debug: false,
					Spec: &spec.Spec{
						Name:        "foo",
						Grammar:     grammars[0],
						Precedences: lr.PrecedenceLevels{},
					},
				},
			},
			expectedErrorRegexes: []string{
				`error on building LALR\(1\) parsing table:`,
				`Error:      Ambiguous Grammar`,
				`Cause:      Multiple conflicts in the parsing table:`,
				`Resolution: Specify associativity and precedence for these Terminals/Productions:`,
			},
		},
		{
			name: "PackageDirNotExist",
			g: &generator{
				UI: ui.NewNop(),
				Params: &Params{
					Debug: false,
					Path:  tempDir,
					Spec: &spec.Spec{
						Name:        "foo",
						Grammar:     grammars[0],
						Precedences: precedences[0],
					},
				},
			},
			expectedErrorRegexes: []string{
				`open .+/foo/foo.go: no such file or directory`,
			},
		},
		{
			name: "Success",
			g: &generator{
				UI: ui.NewNop(),
				Params: &Params{
					Debug: false,
					Path:  tempDir,
					Spec: &spec.Spec{
						Name:        "",
						Grammar:     grammars[0],
						Precedences: precedences[0],
					},
				},
			},
			expectedErrorRegexes: nil,
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
					Debug: false,
					Path:  tempDir,
					Spec: &spec.Spec{
						Name: "foo",
					},
				},
			},
			filename:           "core.go.tmpl",
			data:               nil,
			expectedErrorRegex: `open .+/foo/foo.go: no such file or directory`,
		},
		{
			name: "Success",
			g: &generator{
				UI: ui.NewNop(),
				Params: &Params{
					Debug: false,
					Path:  tempDir,
					Spec: &spec.Spec{
						Name: "",
					},
				},
			},
			filename:           "core.go.tmpl",
			data:               nil,
			expectedErrorRegex: ``,
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

func TestGenerator_generateLexerGraph(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "emerge-test-")
	assert.NoError(t, err)

	defer func() {
		assert.NoError(t, os.RemoveAll(tempDir))
	}()

	tests := []struct {
		name               string
		g                  *generator
		dfa                *automata.DFA
		assocs             []spec.FinalTerminalAssociation
		expectedErrorRegex string
	}{
		{
			name: "DebugFalse",
			g: &generator{
				UI: ui.NewNop(),
				Params: &Params{
					Debug: false,
				},
			},
			dfa:                nil,
			assocs:             nil,
			expectedErrorRegex: ``,
		},
		{
			name: "PackageDirNotExist",
			g: &generator{
				UI: ui.NewNop(),
				Params: &Params{
					Debug: true,
					Path:  tempDir,
					Spec: &spec.Spec{
						Name: "foo",
					},
				},
			},
			dfa:                nil,
			assocs:             nil,
			expectedErrorRegex: `open .+/foo/lexer.dot: no such file or directory`,
		},
		{
			name: "Success",
			g: &generator{
				UI: ui.NewNop(),
				Params: &Params{
					Debug: true,
					Path:  tempDir,
					Spec: &spec.Spec{
						Name: "",
					},
				},
			},
			dfa: automata.NewDFABuilder().
				SetStart(0).
				SetFinal([]automata.State{1}).
				AddTransition(0, '1', '9', 1).
				AddTransition(1, '0', '9', 1).
				Build(),
			assocs: []spec.FinalTerminalAssociation{
				{
					Final:    automata.NewStates(1),
					Terminal: "NUM",
					Kind:     spec.RegexDef,
					Value:    "[1-9][0-9]+",
				},
			},
			expectedErrorRegex: ``,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.g.generateLexerGraph(tc.dfa, tc.assocs)

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

func TestGenerator_generateParsingTable(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "emerge-test-")
	assert.NoError(t, err)

	defer func() {
		assert.NoError(t, os.RemoveAll(tempDir))
	}()

	T, err := lookahead.BuildParsingTable(grammars[0], precedences[0])
	assert.NoError(t, err)

	tests := []struct {
		name               string
		g                  *generator
		T                  *lr.ParsingTable
		expectedErrorRegex string
	}{
		{
			name: "DebugFalse",
			g: &generator{
				UI: ui.NewNop(),
				Params: &Params{
					Debug: false,
				},
			},
			T:                  nil,
			expectedErrorRegex: ``,
		},
		{
			name: "PackageDirNotExist",
			g: &generator{
				UI: ui.NewNop(),
				Params: &Params{
					Debug: true,
					Path:  tempDir,
					Spec: &spec.Spec{
						Name: "foo",
					},
				},
			},
			T:                  nil,
			expectedErrorRegex: `open .+/foo/parser.txt: no such file or directory`,
		},
		{
			name: "Success",
			g: &generator{
				UI: ui.NewNop(),
				Params: &Params{
					Debug: true,
					Path:  tempDir,
					Spec: &spec.Spec{
						Name: "",
					},
				},
			},
			T:                  T,
			expectedErrorRegex: ``,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.g.generateParsingTable(tc.T)

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
