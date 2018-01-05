package dotfiles

import (
	"os"
	"strings"
	"testing"

	"github.com/rhysd/abspath"
)

func TestGetMappingsConfigDirNotExist(t *testing.T) {
	p, err := abspath.ExpandFrom("unknown_directory")
	if err != nil {
		panic(err)
	}
	m, err := GetMappings(p)
	if err != nil {
		t.Fatal(err)
	}
	if len(m) == 0 {
		t.Errorf("Mappings should not be empty. Default value is not set.")
	}
	if m[".vimrc"].String() == "" {
		t.Errorf("Any platform default value must have '.vimrc' mapping. %v", m)
	}
}

func TestGetMappingsConfigFileNotExist(t *testing.T) {
	if err := os.MkdirAll("_test_config", os.ModeDir|os.ModePerm); err != nil {
		panic(err)
	}
	defer os.Remove("_test_config")

	p, err := abspath.ExpandFrom("_test_config")
	if err != nil {
		panic(err)
	}
	m, err := GetMappings(p)
	if err != nil {
		t.Fatal(err)
	}
	if len(m) == 0 {
		t.Errorf("Mappings should not be empty. Default value is not set.")
	}
	if m[".vimrc"].String() == "" {
		t.Errorf("Any platform default value must have '.vimrc' mapping. %v", m)
	}
}

func TestGetMappingsUnknownPlatform(t *testing.T) {
	p, err := abspath.ExpandFrom("unknown_directory")
	if err != nil {
		panic(err)
	}

	m, err := GetMappingsForPlatform("unknown", p)
	if err != nil {
		t.Fatal(err)
	}
	if len(m) != 0 {
		t.Fatalf("Unknown mappings for unknown platform %v", m)
	}
}

func getcwd() abspath.AbsPath {
	cwd, err := abspath.Getwd()
	if err != nil {
		panic(err)
	}
	return cwd
}

func createTestJson(fname, contents string) {
	if err := os.MkdirAll("_test_config", os.ModeDir|os.ModePerm); err != nil {
		panic(err)
	}

	f, err := os.OpenFile(getcwd().Join("_test_config", fname).String(), os.O_CREATE|os.O_RDWR, 0644)
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

	p, err := abspath.ExpandFrom("_test_config")
	if err != nil {
		panic(err)
	}

	m, err := GetMappingsForPlatform("unknown", p)
	if err != nil {
		t.Fatal(err)
	}
	if m["some_file"].String() != "/path/to/some_file" {
		t.Errorf("Mapping value set in mappings.json is wrong: '%s'", m["some_file"])
	}
	if m[".vimrc"].String() != "/override/path/vimrc" {
		t.Errorf("Mapping should be overridden but actually '%s'", m[".vimrc"])
	}

	m, err = GetMappingsForPlatform("darwin", p)
	if err != nil {
		t.Fatal(err)
	}
	if m["some_file"].String() != "/path/to/some_file" {
		t.Errorf("Mapping value set in mappings.json is wrong: '%s' in Darwin", m["some_file"])
	}
	if m[".vimrc"].String() != "/override/path/vimrc" {
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

	p, err := abspath.ExpandFrom("_test_config")
	if err != nil {
		panic(err)
	}

	m, err := GetMappingsForPlatform("darwin", p)
	if err != nil {
		t.Fatal(err)
	}
	if m["some_file"].String() != "/path/to/some_file" {
		t.Errorf("Mapping value set in mappings_darwin.json is wrong: '%s' in Darwin", m["some_file"])
	}
	if m[".vimrc"].String() != "/override/path/vimrc" {
		t.Errorf("Mapping should be overridden by mappings_darwin.json but actually '%s'", m[".vimrc"])
	}

	m, err = GetMappingsForPlatform("windows", p)
	if err != nil {
		t.Fatal(err)
	}
	if m["some_file"].String() != "" {
		t.Errorf("Different configuration must not be loaded but actually some_file was linked to '%s'", m["some_file"])
	}

	// Note: Consider '~' prefix in JSON path value
	if !strings.HasSuffix(m[".vimrc"].String(), DefaultMappings["windows"][".vimrc"][1:]) {
		t.Errorf("Mapping should not be overridden by mappings_darwin.json on different platform (Windows) but actually '%s'", m[".vimrc"])
	}
}

func TestGetMappingsPlatformSpecificMappingsJsonUnix(t *testing.T) {
	createTestJson("mappings_unixlike.json", `
	{
		"some_file": "/path/to/some_file",
		".vimrc": "/hidden/path/vimrc"
	}
	`)
	createTestJson("mappings_darwin.json", `
	{
		".vimrc": "/override/path/vimrc"
	}
	`)
	defer os.RemoveAll("_test_config")

	p, err := abspath.ExpandFrom("_test_config")
	if err != nil {
		panic(err)
	}

	m, err := GetMappingsForPlatform("darwin", p)
	if err != nil {
		t.Fatal(err)
	}
	if m["some_file"].String() != "/path/to/some_file" {
		t.Errorf("Mapping value set in mappings_unixlike.json is wrong: '%s' in Darwin", m["some_file"])
	}
	if m[".vimrc"].String() != "/override/path/vimrc" {
		t.Errorf("Mapping should be overridden by mappings_darwin.json but actually '%s'", m[".vimrc"])
	}

	m, err = GetMappingsForPlatform("windows", p)
	if err != nil {
		t.Fatal(err)
	}
	if m["some_file"].String() != "" {
		t.Errorf("Different configuration must not be loaded but actually some_file was linked to '%s'", m["some_file"])
	}

	// Note: Consider '~' prefix in JSON path value
	if !strings.HasSuffix(m[".vimrc"].String(), DefaultMappings["windows"][".vimrc"][1:]) {
		t.Errorf("Mapping should not be overridden by mappings_unix.json or mappings_darwin.json on different platform (Windows) but actually '%s'", m[".vimrc"])
	}
}

func TestGetMappingsInvalidJson(t *testing.T) {
	createTestJson("mappings.json", `
	{
		"some_file":
	`)
	defer os.RemoveAll("_test_config")

	p, err := abspath.ExpandFrom("_test_config")
	if err != nil {
		panic(err)
	}

	if _, err := GetMappings(p); err == nil {
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

	p, err := abspath.ExpandFrom("_test_config")
	if err != nil {
		panic(err)
	}

	if _, err := GetMappings(p); err == nil {
		t.Fatalf("Empty key must raise an error")
	}
}

func TestGetMappingsInvalidPathValue(t *testing.T) {
	createTestJson("mappings.json", `
	{
		"some_file": "relative-path"
	}`)
	defer os.RemoveAll("_test_config")

	p, err := abspath.ExpandFrom("_test_config")
	if err != nil {
		panic(err)
	}

	if _, err := GetMappings(p); err == nil {
		t.Fatalf("Relative path must be checked")
	}
}

func mapping(k string, v string) Mappings {
	m := make(Mappings, 1)
	m[k] = getcwd().Join(v)
	return m
}

func openFile(n string) *os.File {
	f, err := os.OpenFile(getcwd().Join(n).String(), os.O_CREATE|os.O_RDWR, 0644)
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
	cwd := getcwd()
	source := cwd.Join(n).String()
	s, err := os.Lstat(source)
	if err != nil {
		return false
	}
	if s.Mode()&os.ModeSymlink != os.ModeSymlink {
		return false
	}
	dist, err := os.Readlink(source)
	if err != nil {
		panic(err)
	}
	return dist == cwd.Join(d).String()
}

func TestLinkNormalFile(t *testing.T) {
	m := mapping("._test_source.conf", "_test.conf")
	f := openFile("._test_source.conf")
	defer func() {
		f.Close()
		defer os.Remove("._test_source.conf")
	}()

	err := m.CreateAllLinks(false)
	if err != nil {
		t.Fatal(err)
	}

	if !isSymlinkTo("_test.conf", "._test_source.conf") {
		t.Fatalf("Symbolic link not found")
	}
	defer os.Remove("_test.conf")

	// Skipping already existing link
	err = m.CreateAllLinks(false)
	if err != nil {
		t.Fatal(err)
	}
}

func TestLinkToNonExistingDir(t *testing.T) {
	m := mapping("._source.conf", "_dist_dir/_dist.conf")
	f := openFile("._source.conf")
	defer func() {
		f.Close()
		defer os.Remove("._source.conf")
	}()

	err := m.CreateAllLinks(false)
	if err != nil {
		t.Fatal(err)
	}

	if !isSymlinkTo("_dist_dir/_dist.conf", "._source.conf") {
		t.Fatalf("Symbolic link not found. Directory was not generated to put symlink into?")
	}
	defer os.RemoveAll("_dist_dir")
}

func TestLinkDirSymlink(t *testing.T) {
	m := mapping("._source_dir", "_dist_dir")
	if err := os.MkdirAll("._source_dir", os.ModeDir|os.ModePerm); err != nil {
		panic(err)
	}
	defer os.Remove("._source_dir")

	err := m.CreateAllLinks(false)
	if err != nil {
		t.Fatal(err)
	}

	if !isSymlinkTo("_dist_dir", "._source_dir") {
		t.Fatalf("Symbolic link to directory not found.")
	}
	defer os.Remove("_dist_dir")
}

func TestLinkSpecifiedMappingOnly(t *testing.T) {
	m := mapping("._source.conf", "_dist.conf")
	m["LICENSE.txt"] = getcwd().Join("_never_created.txt")
	f := openFile("._source.conf")
	defer func() {
		f.Close()
		os.Remove("._source.conf")
	}()

	err := m.CreateSomeLinks([]string{"._source.conf"}, false)
	if err != nil {
		t.Fatal(err)
	}

	if !isSymlinkTo("_dist.conf", "._source.conf") {
		t.Fatalf("Symbolic link not found.")
	}
	defer os.Remove("_dist.conf")

	if isSymlinkTo("_never_created.txt", "LICENSE.txt") {
		t.Fatalf("Symbolic link not found.")
	}
}

func TestLinkSpecifyingNonExistingFile(t *testing.T) {
	m := mapping("LICENSE.txt", "never_created.conf")

	err := m.CreateSomeLinks([]string{}, false)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Lstat("never_created.conf"); err == nil {
		t.Errorf("never_created.conf was created")
		os.Remove("never_created.conf")
	}

	err = m.CreateSomeLinks([]string{"unknown_config.conf"}, false)
	if _, ok := err.(*NothingLinkedError); !ok {
		t.Fatal(err)
	}
	if _, err = os.Lstat("never_created.conf"); err == nil {
		t.Errorf("never_created.conf was created")
		os.Remove("never_created.conf")
	}
}

func TestLinkSourceNotExist(t *testing.T) {
	m := mapping(".unknown.conf", "never_created.conf")
	err := m.CreateAllLinks(false)
	if _, ok := err.(*NothingLinkedError); !ok {
		t.Errorf("Not existing file must be ignored but actually error occurred: %s", err.Error())
	}
	m2 := mapping("unknown.conf", "never_created.conf")
	err = m2.CreateSomeLinks([]string{"unknown.conf"}, false)
	if _, ok := err.(*NothingLinkedError); !ok {
		t.Errorf("Not existing file must be ignored but actually error occurred: %s", err.Error())
	}
}

func TestLinkNullDist(t *testing.T) {
	m := Mappings{"License.txt": abspath.AbsPath{}}
	err := m.CreateAllLinks(false)
	if err == nil {
		t.Errorf("Nothing was linked but error did not occur")
	}
}

func TestLinkDryRun(t *testing.T) {
	m := mapping("._test_source.conf", "_test.conf")
	f := openFile("._test_source.conf")
	defer func() {
		f.Close()
		defer os.Remove("._test_source.conf")
	}()

	err := m.CreateAllLinks(true)
	if err != nil {
		t.Fatal(err)
	}

	if isSymlinkTo("_test.conf", "._test_source.conf") {
		t.Fatalf("Symbolic link not found")
	}
}

func createSymlink(from, to string) {
	cwd := getcwd()
	if err := os.Symlink(cwd.Join(from).String(), cwd.Join(to).String()); err != nil {
		panic(err)
	}
}

func TestUnlinkNoFile(t *testing.T) {
	m := mapping("._source.fonf", "._dist.conf")
	if err := m.UnlinkAll(getcwd()); err != nil {
		t.Error(err)
	}
}

func TestUnlinkFiles(t *testing.T) {
	f := openFile("._source.conf")
	defer func() {
		f.Close()
		os.Remove("._source.conf")
	}()
	createSymlink("._source.conf", "._dist.conf")
	m := mapping("._source.fonf", "._dist.conf")
	if err := m.UnlinkAll(getcwd()); err != nil {
		t.Error(err)
	}

	if _, err := os.Lstat("._dist.conf"); err == nil {
		os.Remove("._dist.conf")
		t.Errorf("Unlinked symlink must be removed")
	}
}

func TestUnlinkAnotherFileAlreadyExist(t *testing.T) {
	openFile("._dummy.conf").Close()
	defer os.Remove("._dummy.conf")
	m := mapping("._source.fonf", "._dummy.conf")
	if err := m.UnlinkAll(getcwd()); err != nil {
		t.Error(err)
	}
}

// e.g.
//	expected: dotfiles/vimrc -> ~/.vimrc
//	actual: another_dir/vimrc -> ~/.vimrc
func TestUnlinkDetectLinkToOutsideRepo(t *testing.T) {
	dir := getcwd().Join("_test_dir")

	if err := os.Mkdir(dir.String(), os.ModePerm|os.ModeDir); err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir.String())

	openFile("_outside.conf").Close()
	defer os.Remove("_outside.conf")

	createSymlink("_outside.conf", "_test.conf")
	m := mapping("_another_test.conf", "_test.conf")
	if err := m.UnlinkAll(dir); err != nil {
		t.Error(err)
	}

	if _, err := os.Lstat(getcwd().Join("_test.conf").String()); err != nil {
		t.Fatalf("When target is already linked to outside dotfiles, error should not occur: %s", err.Error())
	}

	os.Remove("_test.conf")
}

func TestActualLinksEmpty(t *testing.T) {
	m := mapping("._source.conf", "._dest.conf")
	l, err := m.ActualLinks(getcwd())
	if err != nil {
		t.Fatal(err)
	}
	if len(l) > 0 {
		t.Errorf("Link does not exist but actually '%v' was reported", l)
	}
}

func TestActualLinksLinkExists(t *testing.T) {
	openFile("._source.conf").Close()
	defer os.Remove("._source.conf")
	createSymlink("._source.conf", "._dist.conf")
	defer os.Remove("._dist.conf")
	cwd := getcwd()
	m := mapping("._source.fonf", "._dist.conf")

	l, err := m.ActualLinks(cwd)
	if err != nil {
		t.Fatal(err)
	}

	if len(l) != 1 {
		t.Fatalf("Only one mapping is intended to be added but actually %d mappings exist", len(l))
	}

	e, ok := l[cwd.Join("._source.conf").String()]
	if !ok {
		t.Fatalf("._source.conf in current directory must be a source of symlink but actually not: '%v'", l)
	}

	expected := cwd.Join("._dist.conf").String()
	if e != expected {
		t.Fatalf("'%s' is expected as a dist of symlink, but actually '%s'", expected, e)
	}
}

func TestActualLinksNotDotfile(t *testing.T) {
	openFile("._source.conf").Close()
	defer os.Remove("._source.conf")
	openFile("._dist.conf").Close()
	defer os.Remove("._dist.conf")
	cwd := getcwd()
	m := mapping("._source.fonf", "._dist.conf")

	l, err := m.ActualLinks(cwd)
	if err != nil {
		t.Fatal(err)
	}

	if len(l) > 0 {
		t.Fatalf("When a mapping is a hard link, it's not a dotfile and should not considered.  But actually links '%v' are detected", l)
	}
}
