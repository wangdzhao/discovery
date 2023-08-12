package util

import "regexp"

func ProcessString(input string) string {
	// 用于匹配多行的空行，并将其替换为一个空行
	re := regexp.MustCompile(`(\n\s*){2,}`)
	return re.ReplaceAllString(input, "\n\n")
}
