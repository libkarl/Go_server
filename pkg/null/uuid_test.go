package null

import (
	"encoding/json"
	"fmt"

	"github.com/stretchr/testify/assert"

	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"

	"testing"
)

// TestNullUUIDMarshal tests the marshal function of the NullUUID
func TestNullUUIDMarshal(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
		id := uuid.NewV4()

		nid := NewUUIDValid(id)
		data, err := json.Marshal(nid)
		require.NoError(t, err)

		assert.Equal(t, fmt.Sprintf("\"%s\"", id.String()), string(data))
	})

	t.Run("Invalid", func(t *testing.T) {
		nid := NewUUIDInvalid()
		data, err := json.Marshal(nid)
		require.NoError(t, err)

		assert.Equal(t, "null", string(data))
	})
}

// TestNullUUIDUnmarshal tests the unmarshal function of the nullUUID.
func TestNullUUIDUnmarshal(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
		id := uuid.NewV4()
		data := fmt.Sprintf("\"%s\"", id.String())

		nid := UUID{}
		err := json.Unmarshal([]byte(data), &nid)
		require.NoError(t, err)

		assert.True(t, nid.Valid)
		assert.Equal(t, nid.UUID, id)
	})

	t.Run("Null", func(t *testing.T) {
		data := []byte("null")

		// having set a non null empty uuid
		nid := NewUUIDValid(empty)
		assert.True(t, nid.Valid)
		// unmarshaling a null will change the value of valid
		err := json.Unmarshal(data, &nid)
		require.NoError(t, err)

		assert.False(t, nid.Valid)
		assert.Equal(t, nid.UUID, empty)
	})
}
