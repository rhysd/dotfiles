package dotfiles

import (
	"os"
	"path"
	"strings"
	"testing"
)

func TestGetMappingsConfigDirNotExist(t *testing.T) {
	m, err := GetMappings("unknown_directory")
	if err != nil {
		t.Fatal(err)
	}
	if len(m) == 0 {
		t.Errorf("Mappings should not be empty. Default value is not set.")
	}
	if m[".vimrc"] == "" {
		t.Errorf("Any platform default value must have '.vimrc' mapping. %v", m)
	}
}

func TestGetMappingsConfigFileNotExist(t *testing.T) {
	if err := os.MkdirAll("_test_config", os.ModeDir|os.ModePerm); err != nil {
		panic(err)
	}
	defer os.Remove("_test_config")

	m, err := GetMappings("_test_config")
	if err != nil {
		t.Fatal(err)
	}
	if len(m) == 0 {
		t.Errorf("Mappings should not be empty. Default value is not set.")
	}
	if m[".vimrc"] == "" {
		t.Errorf("Any platform default value must have '.vimrc' mapping. %v", m)
	}
}

func TestGetMappingsUnknownPlatform(t *testing.T) {
	m, err := GetMappingsForPlatform("unknown", "unknown_directory")
	if err != nil {
		t.Fatal(err)
	}
	if len(m) != 0 {
		t.Fatalf("Unknown mappings for unknown platform %v", m)
	}
}

func createTestJson(fname, contents string) {
	if err := os.MkdirAll("_test_config", os.ModeDir|os.ModePerm); err != nil {
		panic(err)
	}

	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	f, err := os.OpenFile(path.Join(cwd, "_test_config", fname), os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		os.RemoveAll("_test_config")
		panic(err)
	}
	defer f.Close()

	_, err = f.WriteString(contents)
	if err != nil {
		os.RemoveAll("_test_config")
		panic(err)
	}
}

func TestGetMappingsMappingsJson(t *testing.T) {
	createTestJson("mappings.json", `
	{
		"some_file": "/path/to/some_file",
		".vimrc": "/override/path/vimrc",
		".conf": "~/path/in/home"
	}
	`)
	defer os.RemoveAll("_test_config")

	m, err := GetMappingsForPlatform("unknown", "_test_config")
	if err != nil {
		t.Fatal(err)
	}
	if !m["some_file"].Compare("/path/to/some_file") {
		t.Errorf("Mapping value set in mappings.json is wrong: '%s'", m["some_file"])
	}
	if !m[".vimrc"].Compare("/override/path/vimrc") {
		t.Errorf("Mapping should be overridden but actually '%s'", m[".vimrc"])
	}
	if !path.IsAbs(string(m[".conf"])) {
		t.Errorf("'~' must be converted to absolute path: %s", m[".conf"])
	}

	m, err = GetMappingsForPlatform("darwin", "_test_config")
	if err != nil {
		t.Fatal(err)
	}
	if !m["some_file"].Compare("/path/to/some_file") {
		t.Errorf("Mapping value set in mappings.json is wrong: '%s' in Darwin", m["some_file"])
	}
	if !m[".vimrc"].Compare("/override/path/vimrc") {
		t.Errorf("Mapping should be overridden but actually '%s' for Darwin platform", m[".vimrc"])
	}
}

func TestGetMappingsPlatformSpecificMappingsJson(t *testing.T) {
	createTestJson("mappings_darwin.json", `
	{
		"some_file": "/path/to/some_file",
		".vimrc": "/override/path/vimrc"
	}
	`)
	defer os.RemoveAll("_test_config")

	m, err := GetMappingsForPlatform("darwin", "_test_config")
	if err != nil {
		t.Fatal(err)
	}
	if !m["some_file"].Compare("/path/to/some_file") {
		t.Errorf("Mapping value set in mappings_darwin.json is wrong: '%s' in Darwin", m["some_file"])
	}
	if !m[".vimrc"].Compare("/override/path/vimrc") {
		t.Errorf("Mapping should be overridden by mappings_darwin.json but actually '%s'", m[".vimrc"])
	}

	m, err = GetMappingsForPlatform("windows", "_test_config")
	if err != nil {
		t.Fatal(err)
	}
	if !m["some_file"].IsEmpty() {
		t.Errorf("Different configuration must not be loaded but actually some_file was linked to '%s'", m["some_file"])
	}

	// Note: Consider '~' prefix in JSON path value
	if !strings.HasSuffix(string(m[".vimrc"]), DefaultMappings["windows"][".vimrc"][1:]) {
		t.Errorf("Mapping should not be overridden by mappings_darwin.json on different platform (Windows) but actually '%s'", m[".vimrc"])
	}
}

func TestGetMappingsInvalidJson(t *testing.T) {
	createTestJson("mappings.json", `
	{
		"some_file":
	`)
	defer os.RemoveAll("_test_config")

	_, err := GetMappings("_test_config")
	if err == nil {
		t.Fatalf("Invalid Json configuration must raise a parse error")
	}
}

func TestGetMappingsInvalidPathValue(t *testing.T) {
	createTestJson("mappings.json", `
	{
		"some_file": "relative-path"
	}`)
	defer os.RemoveAll("_test_config")

	_, err := GetMappings("_test_config")
	if err == nil {
		t.Fatalf("Relative path must be checked")
	}
}
