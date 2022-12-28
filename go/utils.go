package wheretopark

import "github.com/shopspring/decimal"

func NewIntPricingRule(duration string, price int32) PricingRule {
	return PricingRule{
		Duration:  duration,
		Price:     decimal.NewFromInt32(price),
		Repeating: false,
	}
}

func NewPricingRule(duration string, price decimal.Decimal) PricingRule {
	return PricingRule{
		Duration:  duration,
		Price:     price,
		Repeating: false,
	}
}
