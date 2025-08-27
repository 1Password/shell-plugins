package sdk

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testStruct struct {
	Text   string
	Bytes  []byte
	Ok     bool
	Number int
}

func TestCacheStateGetBytes(t *testing.T) {
	byteData := []byte("I am but a humble piece of data")

	cacheEntry := CacheEntry{
		Data:      byteData,
		ExpiresAt: time.Now().Add(1 * time.Minute),
	}
	cache := CacheState{}
	cache["myKey"] = cacheEntry

	byteResult := make([]byte, len(byteData))
	ok := cache.Get("myKey", &byteResult)

	assert.True(t, ok)
	assert.Equal(t, byteData, byteResult)
}

func TestCacheStateGetStruct(t *testing.T) {
	structData := testStruct{
		Text:   "I am but a humble piece of data",
		Bytes:  []byte("I am but a humble piece of data"),
		Ok:     true,
		Number: 42,
	}

	byteData, err := json.Marshal(structData)
	require.NoError(t, err)

	cacheEntry := CacheEntry{
		Data:      byteData,
		ExpiresAt: time.Now().Add(1 * time.Minute),
	}
	cache := CacheState{}
	cache["myKey"] = cacheEntry

	var structResult testStruct
	ok := cache.Get("myKey", &structResult)

	assert.True(t, ok)
	assert.Equal(t, structData, structResult)
}

func TestCacheStateGetBadInput(t *testing.T) {
	cacheEntry := CacheEntry{
		Data:      []byte("some data"),
		ExpiresAt: time.Now().Add(1 * time.Minute),
	}
	cache := CacheState{}
	cache["myKey"] = cacheEntry
	correctOutput := make([]byte, len(cacheEntry.Data))

	var structResult testStruct
	ok := cache.Get("myKey", &structResult)
	assert.False(t, ok)

	ok = cache.Get("wrongKey", &correctOutput)
	assert.False(t, ok)

	ok = cache.Get("myKey", nil)
	assert.False(t, ok)

	ok = cache.Get("myKey", &correctOutput)
	assert.True(t, ok)
}

func TestCacheOperationsPutBytes(t *testing.T) {
	cacheOps := CacheOperations{
		Puts:    make(CacheState),
		Removes: nil,
	}

	byteData := []byte("I am but a humble piece of data")

	err := cacheOps.Put("session_token", byteData, time.Now().Add(1*time.Minute))
	require.NoError(t, err)

	entry, ok := cacheOps.Puts["session_token"]
	assert.True(t, ok)

	assert.Equal(t, byteData, entry.Data)
}

func TestCacheOperationsPutStruct(t *testing.T) {
	cacheOps := CacheOperations{
		Puts:    make(CacheState),
		Removes: nil,
	}

	structData := testStruct{
		Text:   "I am but a humble piece of data",
		Bytes:  []byte("I am but a humble piece of data"),
		Ok:     true,
		Number: 42,
	}
	err := cacheOps.Put("session_token", structData, time.Now().Add(1*time.Minute))
	require.NoError(t, err)

	entry, ok := cacheOps.Puts["session_token"]
	assert.True(t, ok)

	var structResult testStruct
	err = json.Unmarshal(entry.Data, &structResult)
	require.NoError(t, err)

	assert.Equal(t, structData, structResult)
}

func TestProvisionOutputAddArgsAtIndex(t *testing.T) {
	tc := []struct {
		name     string
		initial  []string
		position int
		args     []string
		expected []string
	}{
		{
			name:     "Insert at the beginning",
			initial:  []string{"arg2", "arg3"},
			position: 0,
			args:     []string{"arg1"},
			expected: []string{"arg1", "arg2", "arg3"},
		},
		{
			name:     "Insert in the middle",
			initial:  []string{"arg1", "arg3"},
			position: 1,
			args:     []string{"arg2"},
			expected: []string{"arg1", "arg2", "arg3"},
		},
		{
			name:     "Insert at the end",
			initial:  []string{"arg1", "arg2"},
			position: -1,
			args:     []string{"arg3"},
			expected: []string{"arg1", "arg2", "arg3"},
		},
		{
			name:     "Append at out-of-range index",
			initial:  []string{"arg1", "arg2"},
			position: 5,
			args:     []string{"arg3"},
			expected: []string{"arg1", "arg2", "arg3"},
		},
		{
			name:     "Insert at negative index (should prepend)",
			initial:  []string{"arg2", "arg3"},
			position: -5,
			args:     []string{"arg1"},
			expected: []string{"arg1", "arg2", "arg3"},
		},
	}

	for _, tc := range tc {
		t.Run(tc.name, func(t *testing.T) {
			out := ProvisionOutput{CommandLine: append([]string{}, tc.initial...)}
			out.AddArgsAtIndex(tc.position, tc.args...)

			if !reflect.DeepEqual(out.CommandLine, tc.expected) {
				t.Errorf("expected %v, got %v", tc.expected, out.CommandLine)
			}
		})
	}
}
