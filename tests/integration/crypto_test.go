//go:build integration
// +build integration

package integration

import (
	"testing"

	"action-tracker/internal/shared/crypto"
)

func TestCrypto_Integration(t *testing.T) {
	t.Run("UUID generation works in real usage", func(t *testing.T) {
		// Test that we can create multiple UUIDs without panic
		for i := 0; i < 10; i++ {
			uuid, err := crypto.NewUUID()
			if err != nil {
				t.Errorf("Failed to generate UUID %d: %v", i, err)
			}
			if uuid.IsNil() {
				t.Errorf("UUID %d should not be nil", i)
			}

			uuidStr := uuid.String()
			if len(uuidStr) != 36 {
				t.Errorf("UUID %d should be 36 characters, got %d", i, len(uuidStr))
			}
		}
	})

	t.Run("UUID parsing works with real data", func(t *testing.T) {
		originalUUID, _ := crypto.NewUUID()
		uuidStr := originalUUID.String()
		t.Logf("Generated UUID: %s", uuidStr)

		parsedUUID, err := crypto.ParseUUID(uuidStr)
		if err != nil {
			t.Logf("Failed to parse UUID: %v", err)
			t.Fatalf("Failed to parse UUID: %v", err)
		}

		parsedStr := parsedUUID.String()
		if parsedStr != uuidStr {
			t.Errorf("UUID round-trip failed: expected %s, got %s",
				uuidStr, parsedStr)
		} else {
			t.Logf("UUID round-trip successful: %s", parsedStr)
		}
	})
}
