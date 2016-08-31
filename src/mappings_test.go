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

func getwd() string {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return cwd
}

func createTestJson(fname, contents string) {
	if err := os.MkdirAll("_test_config", os.ModeDir|os.ModePerm); err != nil {
		panic(err)
	}

	cwd := getwd()

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

func TestGetMappingsEmptyKey(t *testing.T) {
	createTestJson("mappings.json", `
	{
		"": "/path/to/somewhere"
	}
	`)
	defer os.RemoveAll("_test_config")

	_, err := GetMappings("_test_config")
	if err == nil {
		t.Fatalf("Empty key must raise an error")
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

func mapping(k string, v string) Mappings {
	cwd := getwd()
	m := make(Mappings, 1)
	m[k] = AbsolutePath(path.Join(cwd, v))
	return m
}

func openFile(n string) *os.File {
	cwd := getwd()
	f, err := os.OpenFile(path.Join(cwd, n), os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		panic(err)
	}
	_, err = f.WriteString("this file is for test")
	if err != nil {
		panic(err)
	}
	return f
}

func isSymlinkTo(n, d string) bool {
	cwd := getwd()
	source := path.Join(cwd, n)
	s, err := os.Lstat(source)
	if err != nil {
		panic(err)
	}
	if s.Mode()&os.ModeSymlink != os.ModeSymlink {
		return false
	}
	dist, err := os.Readlink(source)
	if err != nil {
		panic(err)
	}
	return dist == path.Join(cwd, d)
}

func TestLinkNormalFile(t *testing.T) {
	m := mapping("._test_source.conf", "_test.conf")
	f := openFile("._test_source.conf")
	defer f.Close()

	err := m.CreateAllLinks(false)
	if err != nil {
		t.Fatal(err)
	}

	if !isSymlinkTo("_test.conf", "._test_source.conf") {
		t.Fatalf("Symbolic link not found")
	}
	defer os.Remove("_test.conf")
	defer os.Remove("._test_source.conf")
}

func TestLinkToNonExistingDir(t *testing.T) {
	m := mapping("._source.conf", "_dist_dir/_dist.conf")
	f := openFile("._source.conf")
	defer f.Close()

	err := m.CreateAllLinks(false)
	if err != nil {
		t.Fatal(err)
	}

	if !isSymlinkTo("_dist_dir/_dist.conf", "._source.conf") {
		t.Fatalf("Symbolic link not found. Directory was not generated to put symlink into?")
	}
	defer os.RemoveAll("_dist_dir")
	defer os.Remove("._source.conf")
}
