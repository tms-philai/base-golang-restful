package utils

import (
	"testing"
)

func TestEnrichMetadata_NilMetadata(t *testing.T) {
	enrichMetadata(nil, nil)
}

func TestResponseHelperFunctions(t *testing.T) {
	t.Skip("Skipping tests that require gin.Context - gin.CreateTestContext not available in this version")
}