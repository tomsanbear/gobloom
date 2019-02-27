package bloom

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	tables := []struct {
		capacity      int
		errorRate     float32
		expectedError bool
	}{
		{1, 0.1, false},
		{0, 0.1, true},
		{1, 100.01, true},
	}

	for _, table := range tables {
		_, err := New(1, table.capacity, table.errorRate)

		if table.expectedError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}
