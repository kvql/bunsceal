package visualise

import (
	"testing"
)

func TestValidateImageVersion(t *testing.T) {
	t.Run("Updated image, validation success", func(t *testing.T) {
		// Test case: Image is up to date
		if ok, err := ValidateImageVersion("../../testdata/image-validation/old-directory", "../../testdata/image-validation/sample_image.png"); !ok && err == nil {
			t.Error("Expected ValidateImageVersion to return true")
		} else if err != nil {
			t.Error("Expected ValidateImageVersion to return true, returned error: ", err)
		}
	})
	// Test case: Image is not up to date
	t.Run("Test validation failure, image not updated", func(t *testing.T) {
		if ok, err := ValidateImageVersion("../../testdata/image-validation/newer-directory", "../../testdata/image-validation/sample_image.png"); ok && err == nil {
			t.Error("Expected ValidateImageVersion to return false")
		} else if err != nil {
			t.Error("Expected ValidateImageVersion to return false, returned error: ", err)
		}
	})
}
