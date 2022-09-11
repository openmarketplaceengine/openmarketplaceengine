package demand

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestLocations(t *testing.T) {
	bytes, err := json.Marshal(greatNeckJobs)
	require.NoError(t, err)
	fmt.Println(string(bytes))
}
