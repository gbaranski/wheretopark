package gdynia

import (
	wheretopark "wheretopark/go"

	_ "embed"

	geojson "github.com/paulmach/go.geojson"
	"github.com/shopspring/decimal"
	"golang.org/x/text/currency"
)

type Configuration struct {
	ParkingLots map[int]wheretopark.Metadata
}

var configuration = Configuration{
	ParkingLots: map[int]wheretopark.Metadata{
		1: {
			Address: "Zawiszy Czarnego 1, 81-374 Gdynia",
			Rules:   typeOneRules,
		},
		2: {
			Address: "Generała Józefa Bema 28, 81-314 Gdynia",
			Rules:   typeOneRules,
		},
		3: {
			Address: "al. Marszałka Piłsudskiego 52/54, 81-382 Gdynia",
			Rules:   typeOneRules,
		},
		4: {
			Address: "Skwer Arki Gdynia, 81-378 Gdynia",
			Rules:   typeTwoRules,
		},
	},
}

func init() {
	for k, v := range configuration.ParkingLots {
		configuration.ParkingLots[k] = wheretopark.Metadata{
			LastUpdated:    &defaultLastUpdated,
			Name:           "",
			Address:        "",
			Geometry:       &geojson.Geometry{},
			Resources:      defaultResources,
			TotalSpots:     v.TotalSpots,
			MaxDimensions:  &wheretopark.Dimensions{},
			Features:       defaultFeatures,
			PaymentMethods: defaultPaymentMethods,
			Comment:        defaultComment,
			Currency:       currency.Unit{},
			Timezone:       defaultTimezone,
			Rules:          v.Rules,
		}
	}
}

var (
	defaultTimezone    = wheretopark.MustLoadLocation("Europe/Warsaw")
	defaultLastUpdated = wheretopark.MustParseDate("2022-12-17")
	defaultResources   = []string{
		"https://www.wtp.waw.pl/parkingi/parking-pr-zeran-pkp/",
		"mailto:ztm@ztm.waw.pl",
		"tel:+48-(22)-56-98-116",
		"tel:+48-(22)-56-98-11",
	}

	defaultFeatures = []string{
		wheretopark.FeatureUncovered,
	}

	defaultPaymentMethods = []string{
		wheretopark.PaymentMethodCash,
		wheretopark.PaymentMethodCard,
		wheretopark.PaymentMethodContactless,
		wheretopark.PaymentMethodMobile,
	}

	defaultComment = map[string]string{
		"pl": `Parking znajduje się w śródmiejskiej strefie płatnego parkowania.
Opłaty za postój pojazdu WD sa określone na https://www.zdiz.gdynia.pl/parkowanie/oplaty/
Bezpłatny postój w niedzielę z identyfikatorem Mieszkańca Gdyni
Opłat nie pobiera się 1 i 6 stycznia, w niedzielę i poniedziałek Świąt Wielkanocnych, 25 i 26 grudnia
Długoterminowe abonamenty są określone na https://www.zdiz.gdynia.pl/parkowanie/oplaty/
		`,
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
)
