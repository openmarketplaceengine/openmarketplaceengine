package csv

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

const csvFile = "testdata/coopdrive-gps-pings-2022.04.06.csv"

func TestImportLocations(t *testing.T) {
	file, err := os.Open(csvFile)
	require.NoError(t, err)

	err = ImportLocations(file)
	require.NoError(t, err)
}
