package helper

// import "regexp"

// func IsCantik(phone string) bool {
// 	// Regex for triple characters (e.g., 111, 222)
// 	tripleRegex := regexp.MustCompile(`(\d)\d*\1\1`)
// 	// Regex for double characters (e.g., 22, 33)
// 	doubleRegex := regexp.MustCompile(`(\d)\d*\1`)

// 	// Check triple characters
// 	if tripleRegex.MatchString(phone) {
// 		return true
// 	}

// 	// Check for at least two double sequences
// 	matches := doubleRegex.FindAllString(phone, -1)
// 	if len(matches) >= 2 {
// 		return true
// 	}

// 	return false
// }
