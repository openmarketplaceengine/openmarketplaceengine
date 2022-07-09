package dao

import (
	"embed"
	"testing"

	"github.com/openmarketplaceengine/openmarketplaceengine/cfg"
	"github.com/stretchr/testify/require"
)

const testFsysPath = "testdata/fsys/"

//go:embed testdata/fsys/*.yaml
//go:embed testdata/fsys/*.sql
var testFsys embed.FS

//-----------------------------------------------------------------------------

func TestFsysExec(t *testing.T) {
	WillTest(t, "test")
	ctx := cfg.Context()
	fse := NewFsysExec(testFsys, testFsysPath, "index.yaml")
	t.Cleanup(func() {
		var drop Drop
		err := ExecTX(ctx, drop.AppendTable("fsys_1", "fsys_2", "fsys_3"))
		if err != nil {
			t.Error(err)
		}
	})
	require.NoError(t, ExecTX(ctx, fse))
}

//-----------------------------------------------------------------------------

func TestFsysIndex(t *testing.T) {
	expect := []string{"table1.sql", "table2.sql", "table3.sql"}
	fse := NewFsysExec(testFsys, testFsysPath, "index.yaml")
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
