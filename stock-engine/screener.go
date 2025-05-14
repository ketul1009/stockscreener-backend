package engine

import (
	"strconv"
)

func FilterStocks(stocks []Stock, rules map[string]interface{}) []string {
	var matched []string

	for _, stock := range stocks {
		if matchesRules(stock, rules) {
			matched = append(matched, stock.Symbol)
		}
	}

	return matched
}

func matchesRules(stock Stock, rules map[string]interface{}) bool {
	for key, val := range rules {
		switch key {
		case "pe_lt":
			strVal, ok := val.(string)
			if !ok {
				return false
			}
			limit, _ := strconv.ParseFloat(strVal, 64)
			if stock.PE >= limit {
				return false
			}
		case "volume_gt":
			strVal, ok := val.(string)
			if !ok {
				return false
			}
			limit, _ := strconv.Atoi(strVal)
			if stock.Volume <= limit {
				return false
			}
		case "sector":
			if stock.Sector != val {
				return false
			}
			// add more filters as needed
		}
	}
	return true
}
