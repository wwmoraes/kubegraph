package kubegraph

import (
	"fmt"
	"regexp"
)

func stringChunks(str string, chunkSize int) []string {
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

func matchLabels(sourceLabels map[string]string, targetLabels map[string]string) bool {
	for sourceLabelName, sourceLabelValue := range sourceLabels {
		targetLabelValue, exists := targetLabels[sourceLabelName]
		if !exists {
			return false
		}

		if sourceLabelValue != targetLabelValue {
			return false
		}
	}

	return true
}
