package dotfiles

import (
	"os"
	"strings"
	"testing"

	"github.com/rhysd/abspath"
)

func getcwd() abspath.AbsPath {
	cwd, err := abspath.Getwd()
	if err != nil {
		panic(err)
	}
	return cwd
}

func createTestDir() string {
	dir := "_test_config"
	if err := os.MkdirAll(dir, os.ModeDir|os.ModePerm); err != nil {
		panic(err)
	}
	return dir
}

func createTestJSON(fname, contents string) string {
	dir := createTestDir()

	f, err := os.OpenFile(getcwd().Join(dir, fname).String(), os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		os.RemoveAll(dir)
		panic(err)
	}
	defer f.Close()

	_, err = f.WriteString(contents)
	if err != nil {
		os.RemoveAll(dir)
		panic(err)
	}

	return dir
}

func hasOnlyDestination(m Mappings, src string, dest string) bool {
	if len(m[src]) != 1 {
		return false
	}
	return m[src][0].String() == dest
}

func mapping(k string, v string) Mappings {
	m := make(Mappings, 1)
	m[k] = []abspath.AbsPath{getcwd().Join(v)}
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

func createSymlink(from, to string) {
	cwd := getcwd()
	if err := os.Symlink(cwd.Join(from).String(), cwd.Join(to).String()); err != nil {
		panic(err)
	}
}

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
	if len(m[".vimrc"]) == 0 {
		t.Errorf("Any platform default value must have '.vimrc' mapping. %v", m)
	}
}

func TestGetMappingsConfigFileNotExist(t *testing.T) {
	testDir := createTestDir()
	defer os.Remove(testDir)

	p, err := abspath.ExpandFrom(testDir)
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
	if len(m[".vimrc"]) == 0 {
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

func TestGetMappingsMappingsJson(t *testing.T) {
	testDir := createTestJSON("mappings.json", `
	{
		"some_file": "/path/to/some_file",
		".vimrc": "/override/path/vimrc",
		".conf": "~/path/in/home",
		"multi_dest": ["/dest1", "/dest2"]
	}
	`)
	defer os.RemoveAll(testDir)

	p, err := abspath.ExpandFrom(testDir)
	if err != nil {
		panic(err)
	}

	m, err := GetMappingsForPlatform("unknown", p)
	if err != nil {
		t.Fatal(err)
	}

	m, err = GetMappingsForPlatform("darwin", p)
	if err != nil {
		t.Fatal(err)
	}
	if !hasOnlyDestination(m, "some_file", "/path/to/some_file") {
		t.Errorf("Mapping value set in mappings.json is wrong: '%s' in Darwin", m["some_file"])
	}
	if !hasOnlyDestination(m, ".vimrc", "/override/path/vimrc") {
		t.Errorf("Mapping should be overridden but actually '%s' for Darwin platform", m[".vimrc"])
	}
	if p := m["multi_dest"]; len(p) != 2 || p[0].String() != "/dest1" || p[1].String() != "/dest2" {
		t.Errorf("Expected two mappings but got '%s' in Darwin", p)
	}
}

func TestGetMappingsPlatformSpecificMappingsJson(t *testing.T) {
	testDir := createTestJSON("mappings_darwin.json", `
	{
		"some_file": "/path/to/some_file",
		".vimrc": "/override/path/vimrc"
	}
	`)
	defer os.RemoveAll(testDir)

	p, err := abspath.ExpandFrom(testDir)
	if err != nil {
		panic(err)
	}

	m, err := GetMappingsForPlatform("darwin", p)
	if err != nil {
		t.Fatal(err)
	}
	if !hasOnlyDestination(m, "some_file", "/path/to/some_file") {
		t.Errorf("Mapping value set in mappings_darwin.json is wrong: '%s' in Darwin", m["some_file"])
	}
	if !hasOnlyDestination(m, ".vimrc", "/override/path/vimrc") {
		t.Errorf("Mapping should be overridden by mappings_darwin.json but actually '%s'", m[".vimrc"])
	}

	m, err = GetMappingsForPlatform("windows", p)
	if err != nil {
		t.Fatal(err)
	}
	if len(m["some_file"]) != 0 {
		t.Errorf("Different configuration must not be loaded but actually some_file was linked to '%s'", m["some_file"])
	}

	// Note: Consider '~' prefix in JSON path value
	if !strings.HasSuffix(m[".vimrc"][0].String(), defaultMappings["windows"][".vimrc"][0][1:]) {
		t.Errorf("Mapping should not be overridden by mappings_darwin.json on different platform (Windows) but actually '%s'", m[".vimrc"][0])
	}
}

func TestGetMappingsPlatformSpecificMappingsJsonUnix(t *testing.T) {
	testDir := createTestJSON("mappings_unixlike.json", `
	{
		"some_file": "/path/to/some_file",
		".vimrc": "/hidden/path/vimrc"
	}
	`)
	createTestJSON("mappings_darwin.json", `
	{
		".vimrc": "/override/path/vimrc"
	}
	`)
	defer os.RemoveAll(testDir)

	p, err := abspath.ExpandFrom(testDir)
	if err != nil {
		panic(err)
	}

	m, err := GetMappingsForPlatform("darwin", p)
	if err != nil {
		t.Fatal(err)
	}
	if !hasOnlyDestination(m, "some_file", "/path/to/some_file") {
		t.Errorf("Mapping value set in mappings_unixlike.json is wrong: '%s' in Darwin", m["some_file"])
	}
	if !hasOnlyDestination(m, ".vimrc", "/override/path/vimrc") {
		t.Errorf("Mapping should be overridden by mappings_darwin.json but actually '%s'", m[".vimrc"])
	}

	m, err = GetMappingsForPlatform("windows", p)
	if err != nil {
		t.Fatal(err)
	}
	if len(m["some_file"]) != 0 {
		t.Errorf("Different configuration must not be loaded but actually some_file was linked to '%s'", m["some_file"])
	}

	// Note: Consider '~' prefix in JSON path value
	if !strings.HasSuffix(m[".vimrc"][0].String(), defaultMappings["windows"][".vimrc"][0][1:]) {
		t.Errorf("Mapping should not be overridden by mappings_unix.json or mappings_darwin.json on different platform (Windows) but actually '%s'", m[".vimrc"][0])
	}
}

func TestGetMappingsInvalidJson(t *testing.T) {
	testDir := createTestJSON("mappings.json", `
	{
		"some_file":
	`)
	defer os.RemoveAll(testDir)

	p, err := abspath.ExpandFrom(testDir)
	if err != nil {
		panic(err)
	}

	if _, err := GetMappings(p); err == nil {
		t.Fatalf("Invalid Json configuration must raise a parse error")
	}
}

func TestGetMappingsEmptyKey(t *testing.T) {
	testDir := createTestJSON("mappings.json", `
	{
		"": "/path/to/somewhere"
	}
	`)
	defer os.RemoveAll(testDir)

	p, err := abspath.ExpandFrom(testDir)
	if err != nil {
		panic(err)
	}

	if _, err := GetMappings(p); err == nil {
		t.Fatalf("Empty key must raise an error")
	}
}

func TestGetMappingsInvalidPathValue(t *testing.T) {
	testDir := createTestJSON("mappings.json", `
	{
		"some_file": "relative-path"
	}`)
	defer os.RemoveAll(testDir)

	p, err := abspath.ExpandFrom(testDir)
	if err != nil {
		panic(err)
	}

	if _, err := GetMappings(p); err == nil {
		t.Fatalf("Relative path must be checked")
	}
}

func TestLinkNormalFile(t *testing.T) {
	cwd := getcwd()
	m := mapping("._test_source.conf", "_test.conf")
	f := openFile("._test_source.conf")
	defer func() {
		f.Close()
		defer os.Remove("._test_source.conf")
	}()

	err := m.CreateAllLinks(cwd, false)
	if err != nil {
		t.Fatal(err)
	}

	if !isSymlinkTo("_test.conf", "._test_source.conf") {
		t.Fatalf("Symbolic link not found")
	}
	defer os.Remove("_test.conf")

	// Skipping already existing link
	err = m.CreateAllLinks(cwd, false)
	if err != nil {
		t.Fatal(err)
	}
}

func TestLinkToNonExistingDir(t *testing.T) {
	cwd := getcwd()
	m := mapping("._source.conf", "_dist_dir/_dist.conf")
	f := openFile("._source.conf")
	defer func() {
		f.Close()
		defer os.Remove("._source.conf")
	}()

	err := m.CreateAllLinks(cwd, false)
	if err != nil {
		t.Fatal(err)
	}

	if !isSymlinkTo("_dist_dir/_dist.conf", "._source.conf") {
		t.Fatalf("Symbolic link not found. Directory was not generated to put symlink into?")
	}
	defer os.RemoveAll("_dist_dir")
}

func TestLinkDirSymlink(t *testing.T) {
	cwd := getcwd()
	m := mapping("._source_dir", "_dist_dir")
	if err := os.MkdirAll("._source_dir", os.ModeDir|os.ModePerm); err != nil {
		panic(err)
	}
	defer os.Remove("._source_dir")

	err := m.CreateAllLinks(cwd, false)
	if err != nil {
		t.Fatal(err)
	}

	if !isSymlinkTo("_dist_dir", "._source_dir") {
		t.Fatalf("Symbolic link to directory not found.")
	}
	defer os.Remove("_dist_dir")
}

func TestLinkSpecifiedMappingOnly(t *testing.T) {
	cwd := getcwd()
	m := mapping("._source.conf", "_dist.conf")
	m["LICENSE.txt"] = []abspath.AbsPath{
		getcwd().Join("_never_created.txt"),
	}
	f := openFile("._source.conf")
	defer func() {
		f.Close()
		os.Remove("._source.conf")
	}()

	err := m.CreateSomeLinks([]string{"._source.conf"}, cwd, false)
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
	cwd := getcwd()
	m := mapping("LICENSE.txt", "never_created.conf")

	err := m.CreateSomeLinks([]string{}, cwd, false)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Lstat("never_created.conf"); err == nil {
		t.Errorf("never_created.conf was created")
		os.Remove("never_created.conf")
	}

	err = m.CreateSomeLinks([]string{"unknown_config.conf"}, cwd, false)
	if _, ok := err.(*NothingLinkedError); !ok {
		t.Fatal(err)
	}
	if _, err = os.Lstat("never_created.conf"); err == nil {
		t.Errorf("never_created.conf was created")
		os.Remove("never_created.conf")
	}
}

func TestLinkSourceNotExist(t *testing.T) {
	cwd := getcwd()
	m := mapping(".unknown.conf", "never_created.conf")
	err := m.CreateAllLinks(cwd, false)
	if _, ok := err.(*NothingLinkedError); !ok {
		t.Errorf("Not existing file must be ignored but actually error occurred: %s", err.Error())
	}
	m2 := mapping("unknown.conf", "never_created.conf")
	err = m2.CreateSomeLinks([]string{"unknown.conf"}, cwd, false)
	if _, ok := err.(*NothingLinkedError); !ok {
		t.Errorf("Not existing file must be ignored but actually error occurred: %s", err.Error())
	}
}

func TestLinkNullDest(t *testing.T) {
	cwd := getcwd()
	m := Mappings{
		"empty":     []abspath.AbsPath{},
		"null_only": []abspath.AbsPath{abspath.AbsPath{}},
	}
	err := m.CreateAllLinks(cwd, false)
	if err == nil {
		t.Errorf("Nothing was linked but error did not occur")
	}
}

func TestLinkDryRun(t *testing.T) {
	cwd := getcwd()
	m := mapping("._test_source.conf", "_test.conf")
	f := openFile("._test_source.conf")
	defer func() {
		f.Close()
		defer os.Remove("._test_source.conf")
	}()

	err := m.CreateAllLinks(cwd, true)
	if err != nil {
		t.Fatal(err)
	}

	if isSymlinkTo("_test.conf", "._test_source.conf") {
		t.Fatalf("Symbolic link should not be found")
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
	defer os.Remove("_test.conf")

	m := mapping("_another_test.conf", "_test.conf")
	if err := m.UnlinkAll(dir); err != nil {
		t.Error(err)
	}

	if _, err := os.Lstat(getcwd().Join("_test.conf").String()); err != nil {
		t.Fatalf("When target is already linked to outside dotfiles, error should not occur: %s", err.Error())
	}
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

	if l[0].src != cwd.Join("._source.conf").String() {
		t.Fatalf("._source.conf in current directory must be a source of symlink but actually not: '%v'", l)
	}

	expected := cwd.Join("._dist.conf").String()
	if l[0].dst != expected {
		t.Fatalf("'%s' is expected as a dist of symlink, but actually '%s'", expected, l[0].dst)
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

func TestActualLinksTwoDestsFromOneSource(t *testing.T) {
	openFile("._source.conf").Close()
	defer os.Remove("._source.conf")
	createSymlink("._source.conf", "._dest1.conf")
	defer os.Remove("._dest1.conf")
	createSymlink("._source.conf", "._dest2.conf")
	defer os.Remove("._dest2.conf")
	cwd := getcwd()
	m := Mappings{
		"._source.conf": []abspath.AbsPath{getcwd().Join("._dest1.conf"), getcwd().Join("._dest2.conf")},
	}

	links, err := m.ActualLinks(cwd)
	if err != nil {
		t.Fatal(err)
	}

	if len(links) != 2 {
		t.Fatalf("Two mappings are intended to be added but actually %d mappings exist", len(links))
	}

	src := cwd.Join("._source.conf").String()
	for i, c := range []string{"._dest1.conf", "._dest2.conf"} {
		l := links[i]
		if l.src != src {
			t.Fatalf("Wanted %+v but got %+v for source (index=%d)", src, l.src, i)
		}
		dst := cwd.Join(c).String()
		if l.dst != dst {
			t.Fatalf("Wanted %+v but got %+v for source (index=%d)", dst, l.dst, i)
		}
	}
}

func TestConvertMappingsJSONToMappings(t *testing.T) {
	json := mappingsJSON{
		"empty":     []string{},
		"null_only": []string{""},
	}
	m, err := convertMappingsJSONToMappings(json)
	if err != nil {
		t.Fatal(err)
	}
	if len(m["empty"]) != 0 {
		t.Fatalf("Converted mapping value for `empty` is wrong: '%v'", m["empty"])
	}
	// Expected value for `null_only` is also an empty slice,
	// because the empty string is ignored when converting.
	if len(m["null_only"]) != 0 {
		t.Fatalf("Converted mapping value for `null_only` is wrong: '%v'", m["null_only"])
	}
}
