package utils

func SafeStringPointerValue(stringPointer *string) string {
	var result string

	if stringPointer == nil {
		return result
	}

	return *stringPointer
}
