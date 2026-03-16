package utils

import "strings"

func SafeStringPointerValue(stringPointer *string) string {
	var result string

	if stringPointer == nil {
		return result
	}

	return *stringPointer
}

func GetStringPointer(str string) *string {
	return &str
}

func SanitizeDocumentNumber(documentNumber string) string {
	var sanitizedDocNumber strings.Builder

	for _, character := range documentNumber {
		if character >= '0' && character <= '9' {
			sanitizedDocNumber.WriteRune(character)
		}
	}

	return sanitizedDocNumber.String()
}
