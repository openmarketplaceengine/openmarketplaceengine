package csv

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

const csvFile = "testdata/coopdrive-gps-pings-2022.04.06.csv"

func TestScanner(t *testing.T) {

	file, err := os.Open(csvFile)
	require.NoError(t, err)

	scanner := NewScan(file)

	count := 0

	for {
		l, err := scanner.NextLocation()
		require.NoError(t, err)
		if l == nil {
			break
		}
		count++
	}

	require.Equal(t, 5779, count)
}
