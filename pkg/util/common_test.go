package util

import (
	"testing"
)

func TestGetLatestCommitTime(t *testing.T) {
	_, err := GetLatestCommitTime("../../testdata/image-validation/sample_image.png")
	if err != nil {
		t.Error("Test failed:", err)
	}

	t.Log("Test passed")
}
