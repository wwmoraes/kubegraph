package utils

import (
	"fmt"
	"regexp"
)

func StringChunks(str string, chunkSize int) []string {
	if len(str) <= chunkSize {
		return []string{str}
	}

	re := regexp.MustCompile(fmt.Sprintf(`([\S\r\n]{0,%d})`, chunkSize))

	var chunks []string

	results := re.FindAllStringSubmatch(str, -1)
	for _, match := range results {
		chunks = append(chunks, match[1])
	}

	return chunks
}
