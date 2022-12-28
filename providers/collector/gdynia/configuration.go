package gdynia

import (
	wheretopark "wheretopark/go"

	_ "embed"

	"github.com/shopspring/decimal"
)

type Configuration struct {
	ParkingLots map[int]wheretopark.Metadata
}

var configuration = Configuration{
	ParkingLots: map[int]wheretopark.Metadata{
		1: {
			TotalSpots: map[string]uint{
				wheretopark.SpotTypeCar: 78,
			},
			Rules: typeOneRules,
		},
		2: {
			TotalSpots: map[string]uint{
				wheretopark.SpotTypeCar: 50,
			},
			Rules: typeOneRules,
		},
		3: {
			TotalSpots: map[string]uint{
				wheretopark.SpotTypeCar: 50,
			},
			Rules: typeOneRules,
		},
		4: {
			TotalSpots: map[string]uint{
				wheretopark.SpotTypeCar: 75,
			},
			Rules: typeTwoRules,
		},
	},
}

func init() {
	for k, v := range configuration.ParkingLots {
		v.Resources = append(v.Resources, baseResources...)
		v.Features = append(v.Features, baseFeatures...)
		v.PaymentMethods = defaultPaymentMethods
		v.Comment = defaultComment
		configuration.ParkingLots[k] = v
	}
}

var (
	baseResources = []string{
		"https://www.wtp.waw.pl/parkingi/parking-pr-zeran-pkp/",
		"mailto:ztm@ztm.waw.pl",
		"tel:+48-(22)-56-98-116",
		"tel:+48-(22)-56-98-11",
	}

	baseFeatures = []string{
		wheretopark.FeatureUncovered,
	}

	typeOneRules = []wheretopark.Rule{
		{
			Hours: "Mo-Su 08:00-20:00",
			Pricing: []wheretopark.PricingRule{
				wheretopark.NewPricingRule("PT1H", decimal.RequireFromString("5.50")),
				wheretopark.NewPricingRule("PT2H", decimal.RequireFromString("12.10")),
				wheretopark.NewPricingRule("PT3H", decimal.RequireFromString("20.00")),
				wheretopark.NewPricingRule("PT1H", decimal.RequireFromString("5.50")).Repeated(),
			},
		},
		{
			Hours: "Mo-Su 20:00-08:00",
			Pricing: []wheretopark.PricingRule{
				wheretopark.NewPricingRule("PT1H", decimal.Zero),
			},
		},
		{
			Hours: "24/7",
			Applies: []string{
				wheretopark.SpotTypeCarElectric,
				wheretopark.SpotTypeCarDisabled,
			},
			Pricing: []wheretopark.PricingRule{
				wheretopark.NewPricingRule("PT1H", decimal.Zero),
			},
		},
	}

	typeTwoRules = []wheretopark.Rule{
		{
			Hours: "Mo-Su 08:00-20:00",
			Pricing: []wheretopark.PricingRule{
				wheretopark.NewPricingRule("PT1H", decimal.RequireFromString("3.90")),
				wheretopark.NewPricingRule("PT2H", decimal.RequireFromString("8.50")),
				wheretopark.NewPricingRule("PT3H", decimal.RequireFromString("14.00")),
				wheretopark.NewPricingRule("PT1H", decimal.RequireFromString("3.90")).Repeated(),
			},
		},
		{
			Hours: "Mo-Su 20:00-08:00",
			Pricing: []wheretopark.PricingRule{
				wheretopark.NewPricingRule("PT1H", decimal.Zero),
			},
		},
		{
			Hours: "24/7",
			Applies: []string{
				wheretopark.SpotTypeCarElectric,
				wheretopark.SpotTypeCarDisabled,
			},
			Pricing: []wheretopark.PricingRule{
				wheretopark.NewPricingRule("PT1H", decimal.Zero),
			},
		},
	}

	defaultPaymentMethods = []string{
		wheretopark.PaymentMethodCash,
		wheretopark.PaymentMethodCard,
		wheretopark.PaymentMethodContactless,
		wheretopark.PaymentMethodMobile,
	}

	defaultComment = map[string]string{
		"pl": `
			Parking znajduje się w śródmiejskiej strefie płatnego parkowania.
			
			Opłaty za postój pojazdu WD sa określone na https://www.zdiz.gdynia.pl/parkowanie/oplaty/
		  
			Bezpłatny postój w niedzielę z identyfikatorem Mieszkańca Gdyni
		  
			Opłat nie pobiera się 1 i 6 stycznia, w niedzielę i poniedziałek Świąt Wielkanocnych, 25 i 26 grudnia
		  
			Długoterminowe abonamenty są określone na https://www.zdiz.gdynia.pl/parkowanie/oplaty/
		`,
	}
)
