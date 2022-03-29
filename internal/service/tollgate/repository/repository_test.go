package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRepository(t *testing.T) {
	t.Run("testFindAll", func(t *testing.T) {
		testFindAll(t)
	})
}

func testFindAll(t *testing.T) {
	tg, err := FindAll()
	require.NoError(t, err)

	assert.Len(t, tg.BboxTollgates, 2)
	assert.Len(t, tg.BboxTollgates[0].Boxes, 5)
	assert.Equal(t, tg.BboxTollgates[0].BoxesRequired, 3)
	assert.NotEmpty(t, tg.BboxTollgates[0].Name)

	assert.Len(t, tg.LineTollgates, 2)
	assert.NotEmpty(t, tg.LineTollgates[0].Name)
}
