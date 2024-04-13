package utils

import "biocad_internship/internal/model"

func FindUniqueGuids(messages []model.Message) []string {
	uniqueGUIDs := make(map[string]bool)
	var result []string

	for _, msg := range messages {
		if _, exists := uniqueGUIDs[msg.UnitGUID]; !exists {
			uniqueGUIDs[msg.UnitGUID] = true
			result = append(result, msg.UnitGUID)
		}
	}

	return result
}
