package main

import (
	"encoding/json"
	"fmt"
)

// ConditionMatch matches conditions in policy to attributes on zone
func ConditionMatch(attributes *Zone, condition Condition) bool {
	fmt.Printf("Attributes: %v condition: %s", attributes, condition)
	if condition.Type == "matchSuffix" {
		return attributes.SuffixMatch(condition.Value.(string))
	}
	return false
}

func ConditionMatchFunc(args ...interface{}) (interface{}, error) {
	attributes := args[0].(*Zone)
	var condition Condition
	err := json.Unmarshal([]byte(args[1].(string)), &condition)
	if err != nil {
		return nil, err
	}

	return (bool)(ConditionMatch(attributes, condition)), nil
}