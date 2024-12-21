package helper

import (
	"mime/multipart"
	"strings"

	"github.com/go-playground/validator/v10"
)

func ImageFile(fl validator.FieldLevel) bool {
	file, ok := fl.Field().Interface().(*multipart.FileHeader)
	if !ok {
		return false
	}

	allowedExtensions := []string{".jpg", ".jpeg", ".png", ".gif"}
	for _, ext := range allowedExtensions {
		if strings.HasSuffix(strings.ToLower(file.Filename), ext) {
			return true
		}
	}

	return false
}
