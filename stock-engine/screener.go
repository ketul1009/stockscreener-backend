package engine

import (
	"strconv"
)

func FilterStocks(stocks []Stock, rules []Rule) []string {
	var matched []string

	for _, stock := range stocks {
		if matchesRules(stock, rules) {
			matched = append(matched, stock.Symbol)
		}
	}

	return matched
}

func matchesRules(stock Stock, rules []Rule) bool {
	// Get the "indicators" sub-map
	indicatorsRaw, ok := stock.Indicators["indicators"]
	if !ok {
		return false
	}
	indicators, ok := indicatorsRaw.(map[string]interface{})
	if !ok {
		return false
	}

	truthValueArray := []bool{}
	conditionArray := []string{}

	for _, rule := range rules {
		if rule.Type == "filter" {
			ruleIndicatorValue := interfaceGet[float64](indicators, rule.TechnicalIndicator, 0)
			var comparisonValue float64

			if rule.ComparisonType == "indicator" {
				comparisonValue = interfaceGet[float64](indicators, rule.ComparisonValue.(string), 0)
			} else if rule.ComparisonType == "number" {
				// Handle string to float64 conversion
				if strVal, ok := rule.ComparisonValue.(string); ok {
					val, err := strconv.ParseFloat(strVal, 64)
					if err != nil {
						return false
					}
					comparisonValue = val
				} else if floatVal, ok := rule.ComparisonValue.(float64); ok {
					comparisonValue = floatVal
				} else {
					return false
				}
			}

			truthValue := evaluateTruthValue(ruleIndicatorValue, comparisonValue, rule.Condition)
			truthValueArray = append(truthValueArray, truthValue)
		} else if rule.Type == "condition" {
			conditionArray = append(conditionArray, rule.Condition)
		}
	}

	if len(truthValueArray) == 0 {
		return false
	}

	truthValue := evaluateRuleTruth(truthValueArray, conditionArray)
	return truthValue
}

func evaluateTruthValue(ruleIndicatorValue float64, ruleComparisonValue float64, ruleCondition string) bool {
	switch ruleCondition {
	case "greater_than":
		return ruleIndicatorValue > ruleComparisonValue
	case "less_than":
		return ruleIndicatorValue < ruleComparisonValue
	case "equal_to":
		return ruleIndicatorValue == ruleComparisonValue
	}
	return false
}

func evaluateRuleTruth(truthValueArray []bool, conditionArray []string) bool {
	truthValue := truthValueArray[0]
	for i, condition := range conditionArray {
		if condition == "AND" {
			truthValue = truthValue && truthValueArray[i+1]
		} else if condition == "OR" {
			truthValue = truthValue || truthValueArray[i+1]
		}
	}
	return truthValue
}
