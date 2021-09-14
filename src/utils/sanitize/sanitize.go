package sanitize

// Sanitize string
func Sanitize(str string) string {
	var sanStr string
	for _, char := range str {
		if len(string(char)) < 4 {
			sanStr = sanStr + string(char)
		}
	}
	return sanStr
}
