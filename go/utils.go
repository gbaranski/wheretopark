package wheretopark

import (
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

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

func WithTimeout(fn func() error, timeout time.Duration) error {
	done := make(chan error)
	go func() {
		done <- fn()
	}()
	select {
	case <-time.After(timeout):
		return fmt.Errorf("timeout")
	case err := <-done:
		return err
	}
}
