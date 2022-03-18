package dao

import (
	"embed"
	"testing"

	"github.com/openmarketplaceengine/openmarketplaceengine/cfg"
	"github.com/stretchr/testify/require"
)

const fsysPath = "testdata/fsys/"

//go:embed testdata/fsys/*.yaml
//go:embed testdata/fsys/*.sql
var testFsys embed.FS

//-----------------------------------------------------------------------------

func TestFsysExec(t *testing.T) {
	WillTest(t, "test")
	fse := NewFsysExec(testFsys, fsysPath, "index.yaml")
	require.NoError(t, ExecTX(cfg.Context(), fse))
}

//-----------------------------------------------------------------------------

func TestFsysIndex(t *testing.T) {
	expect := []string{"table1.sql", "table2.sql", "table3.sql"}
	fse := NewFsysExec(testFsys, fsysPath, "index.yaml")
	ary, err := fse.Names()
	require.NoError(t, err)
	require.Equal(t, expect, ary)
}

//-----------------------------------------------------------------------------

func TestIndexYAML(t *testing.T) {
	const data = `
- a
- b
- c
`
	var index Index
	testIndexReader(t, "YAML", index.readYAML, data, []string{"a", "b", "c"})
}

//-----------------------------------------------------------------------------

func TestIndexJSON(t *testing.T) {
	const data = `["a", "b", "c"]`
	var index Index
	testIndexReader(t, "JSON", index.readJSON, data, []string{"a", "b", "c"})
}

//-----------------------------------------------------------------------------
// Helpers
//-----------------------------------------------------------------------------

func testIndexReader(t *testing.T, name string, reader indexReader, data string, expect []string) {
	ary, err := reader([]byte(data))
	require.NoError(t, err, "index reader %q failed: %s", name)
	require.Equal(t, expect, ary, "index reader %q mismatch", name)
}
