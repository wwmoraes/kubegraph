package utils

func MatchLabels(sourceLabels map[string]string, targetLabels map[string]string) bool {
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
