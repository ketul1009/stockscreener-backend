package engine

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
	// Get the "indicators" sub-map
	indicatorsRaw, ok := stock.Indicators["indicators"]
	if !ok {
		return false
	}
	indicators, ok := indicatorsRaw.(map[string]interface{})
	if !ok {
		return false
	}

	// Get EMA5 and EMA20 as float64
	ema5Raw, ok1 := indicators["EMA5"]
	ema20Raw, ok2 := indicators["EMA20"]
	if !ok1 || !ok2 {
		return false
	}
	ema5, ok1 := ema5Raw.(float64)
	ema20, ok2 := ema20Raw.(float64)
	if ok1 && ok2 && ema5 > ema20 {
		return true
	}
	return false
}
