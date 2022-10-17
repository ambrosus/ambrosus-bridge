package explorers_clients

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToLower(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		v := ToLower("Hello")
		assert.Equal(t, "hello", v)
	})

	t.Run("slice of strigs", func(t *testing.T) {
		v := ToLower([]string{"Hi", "Hello", "lol"})
		assert.Equal(t, []string{"hi", "hello", "lol"}, v)
	})
}

func TestFilterTxsByFromToAddresses(t *testing.T) {
	t.Run("From as string", func(t *testing.T) {
		inputTxs := []*Transaction{
			{From: "123", To: "111"},
			{From: "456", To: "222"},
		}
		inputFrom := "456"
		inputTo := "222"

		expectedTxs := []*Transaction{
			{From: "456", To: "222"},
		}
		assert.Equal(t, expectedTxs, FilterTxsByFromToAddresses(inputTxs, inputFrom, inputTo))
	})

	t.Run("From as slice of strings", func(t *testing.T) {
		inputTxs := []*Transaction{
			{From: "123", To: "111"},
			{From: "123", To: "222"},
			{From: "456", To: "222"},
		}
		inputFrom := []string{"123", "456"}
		inputTo := "222"

		expectedTxs := []*Transaction{
			{From: "123", To: "222"},
			{From: "456", To: "222"},
		}
		assert.Equal(t, expectedTxs, FilterTxsByFromToAddresses(inputTxs, inputFrom, inputTo))
	})
}
