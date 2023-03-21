package glasgow

import (
	wheretopark "wheretopark/go"

	_ "embed"

	"github.com/shopspring/decimal"
)

type Configuration struct {
	ParkingLots map[string]wheretopark.Metadata
}

var configuration = Configuration{
	ParkingLots: map[string]wheretopark.Metadata{
		"CPG25C_1": {
			Name:    "SEC Car Park",
			Address: "Stobcross Rd, Glasgow G3 8YW, United Kingdom",
			MaxDimensions: &wheretopark.Dimensions{
				Height: 200,
			},
			Comment: map[string]string{
				"en": "Event & concert pre-pay facilities available 12 hours - £11.00, pay at the car park pay station.",
			},
			Rules: []wheretopark.Rule{
				{
					Hours: "24/7",
					Pricing: []wheretopark.PricingRule{
						wheretopark.NewPricingRule("PT1H", decimal.RequireFromString("5.50")),
						wheretopark.NewPricingRule("PT12H", decimal.RequireFromString("11.00")),
						wheretopark.NewPricingRule("PT13H", decimal.RequireFromString("16.00")),
						wheretopark.NewPricingRule("P1D", decimal.RequireFromString("22.00")),
					},
				},
			},
		},

		"CPG24C_1": {
			Name:    "Duke Street Car Park",
			Address: "Unnamed Road, Glasgow G4 0UW, United Kingdom",
			MaxDimensions: &wheretopark.Dimensions{
				Height: 220,
			},
			Comment: map[string]string{
				"en": "All day parking special now only £5.00! Pay at the car park pay machine before you exit, no product required!!",
			},
			Rules: []wheretopark.Rule{
				{
					Hours: "Mo-Sa 08:00-18:00",
					Pricing: []wheretopark.PricingRule{
						wheretopark.NewPricingRule("PT1H", decimal.RequireFromString("2.00")),
						wheretopark.NewPricingRule("PT2H", decimal.RequireFromString("4.00")),
						wheretopark.NewPricingRule("PT10H", decimal.RequireFromString("5.00")),
						// Commented because of the promotion mentioned above in comment field.
						// wheretopark.NewPricingRule("PT3H", decimal.RequireFromString("6.00")),
						// wheretopark.NewPricingRule("PT4H", decimal.RequireFromString("8.00")),
						// wheretopark.NewPricingRule("PT5H", decimal.RequireFromString("10.50")),
						// wheretopark.NewPricingRule("PT6H", decimal.RequireFromString("13.00")),
						// wheretopark.NewPricingRule("PT7H", decimal.RequireFromString("15.50")),
						// wheretopark.NewPricingRule("PT8H", decimal.RequireFromString("18.00")),
						// wheretopark.NewPricingRule("PT9H", decimal.RequireFromString("20.50")),
						// wheretopark.NewPricingRule("PT10H", decimal.RequireFromString("23.00")),
					},
				},
				{
					Hours: "Mo-Su 18:00-08:00",
					Pricing: []wheretopark.PricingRule{
						wheretopark.NewPricingRule("PT14H", decimal.RequireFromString("5.00")),
					},
				},
				{
					Hours: "Su 08:00-18:00",
					Pricing: []wheretopark.PricingRule{
						wheretopark.NewPricingRule("PT10H", decimal.RequireFromString("5.00")),
					},
				},
			},
		},

		"CPG21C_1": {
			Name:    "Dundasvale Car Park II",
			Address: "10 Dundasvale Ct, Glasgow G4 0JS, United Kingdom",
			MaxDimensions: &wheretopark.Dimensions{
				Height: 198,
			},
			Rules: []wheretopark.Rule{
				{
					Hours: "Mo-Sa 08:00-18:00",
					Pricing: []wheretopark.PricingRule{
						wheretopark.NewPricingRule("PT1H", decimal.RequireFromString("2.20")),
						wheretopark.NewPricingRule("PT2H", decimal.RequireFromString("4.40")),
						wheretopark.NewPricingRule("PT3H", decimal.RequireFromString("6.60")),
						wheretopark.NewPricingRule("PT4H", decimal.RequireFromString("8.80")),
						wheretopark.NewPricingRule("PT5H", decimal.RequireFromString("11.00")),
						wheretopark.NewPricingRule("PT6H", decimal.RequireFromString("13.20")),
						wheretopark.NewPricingRule("PT9H", decimal.RequireFromString("15.00")),
						wheretopark.NewPricingRule("PT10H", decimal.RequireFromString("23.00")),
					},
				},
				{
					Hours: "Mo-Su 18:00-08:00",
					Pricing: []wheretopark.PricingRule{
						wheretopark.NewPricingRule("PT14H", decimal.RequireFromString("4.00")),
					},
				},
				{
					Hours: "Su 08:00-18:00",
					Pricing: []wheretopark.PricingRule{
						wheretopark.NewPricingRule("PT10H", decimal.RequireFromString("5.00")),
					},
				},
			},
		},

		"CPG07C_1": {
			Name:    "Charing Cross Car Park",
			Address: "41 Elmbank St, Glasgow G2 4PG, United Kingdom",
			MaxDimensions: &wheretopark.Dimensions{
				Height: 170,
			},
			Rules: []wheretopark.Rule{
				{
					Hours: "Mo-Sa 08:00-18:00",
					Pricing: []wheretopark.PricingRule{
						wheretopark.NewPricingRule("PT1H", decimal.RequireFromString("2.80")),
						wheretopark.NewPricingRule("PT2H", decimal.RequireFromString("6.00")),
						wheretopark.NewPricingRule("PT3H", decimal.RequireFromString("9.00")),
						wheretopark.NewPricingRule("PT4H", decimal.RequireFromString("12.00")),
						wheretopark.NewPricingRule("PT5H", decimal.RequireFromString("15.00")),
						wheretopark.NewPricingRule("PT6H", decimal.RequireFromString("18.00")),
						wheretopark.NewPricingRule("PT7H", decimal.RequireFromString("21.00")),
						wheretopark.NewPricingRule("PT8H", decimal.RequireFromString("24.00")),
						wheretopark.NewPricingRule("PT9H", decimal.RequireFromString("26.50")),
						wheretopark.NewPricingRule("PT10H", decimal.RequireFromString("29.00")),
					},
				},
				{
					Hours: "Mo-Su 18:00-08:00",
					Pricing: []wheretopark.PricingRule{
						wheretopark.NewPricingRule("PT14H", decimal.RequireFromString("4.00")),
					},
				},
				{
					Hours: "Su 08:00-18:00",
					Pricing: []wheretopark.PricingRule{
						wheretopark.NewPricingRule("PT10H", decimal.RequireFromString("5.00")),
					},
				},
			},
		},

		"CPG06C_1": {
			Name:    "Cadogan Square Car Park",
			Address: "Cadogan Square Car Park, 25 Douglas St, Glasgow G2 7BG, United Kingdom",
			MaxDimensions: &wheretopark.Dimensions{
				Height: 200,
			},
			Rules: []wheretopark.Rule{
				{
					Hours: "Mo-Sa 07:00-21:00",
					Pricing: []wheretopark.PricingRule{
						wheretopark.NewPricingRule("PT1H", decimal.RequireFromString("2.80")),
						wheretopark.NewPricingRule("PT2H", decimal.RequireFromString("6.00")),
						wheretopark.NewPricingRule("PT3H", decimal.RequireFromString("9.00")),
						wheretopark.NewPricingRule("PT4H", decimal.RequireFromString("12.00")),
						wheretopark.NewPricingRule("PT5H", decimal.RequireFromString("15.00")),
						wheretopark.NewPricingRule("PT6H", decimal.RequireFromString("18.00")),
						wheretopark.NewPricingRule("PT7H", decimal.RequireFromString("21.00")),
						wheretopark.NewPricingRule("PT8H", decimal.RequireFromString("24.00")),
						wheretopark.NewPricingRule("PT9H", decimal.RequireFromString("26.50")),
						wheretopark.NewPricingRule("PT10H", decimal.RequireFromString("29.00")),
					},
				},
			},
		},

		"CPG04C_1": {
			Name:    "Buchanan Galleries Car",
			Address: "Buchanan galleries, Car Park, Glasgow G1 2FF, United Kingdom",
			MaxDimensions: &wheretopark.Dimensions{
				Height: 210,
			},
			Rules: []wheretopark.Rule{
				{
					Hours: "Mo-Sa 07:00-17:00",
					Pricing: []wheretopark.PricingRule{
						wheretopark.NewPricingRule("PT1H", decimal.RequireFromString("2.80")),
						wheretopark.NewPricingRule("PT2H", decimal.RequireFromString("3.80")),
						wheretopark.NewPricingRule("PT3H", decimal.RequireFromString("5.00")),
						wheretopark.NewPricingRule("PT4H", decimal.RequireFromString("7.50")),
						wheretopark.NewPricingRule("PT5H", decimal.RequireFromString("10.00")),
						wheretopark.NewPricingRule("PT7H", decimal.RequireFromString("13.00")),
						wheretopark.NewPricingRule("PT9H", decimal.RequireFromString("15.00")),
						wheretopark.NewPricingRule("PT24H", decimal.RequireFromString("25.00")),
					},
				},
				{
					Hours: "Mo-Sa 17:00-21:00",
					Pricing: []wheretopark.PricingRule{
						wheretopark.NewPricingRule("PT1H", decimal.RequireFromString("2.00")).Repeated(),
					},
				},
				{
					Hours: "Su 09:00-21:00",
					Pricing: []wheretopark.PricingRule{
						wheretopark.NewPricingRule("PT1H", decimal.RequireFromString("2.00")),
						wheretopark.NewPricingRule("PT2H", decimal.RequireFromString("4.00")).Repeated(),
					},
				},
			},
		},

		"CPG03C_1": {
			Name:    "Cambridge Street Car Park",
			Address: "Cambridge St, Glasgow G3 6RU, United Kingdom",
			MaxDimensions: &wheretopark.Dimensions{
				Height: 198,
			},
			Rules: []wheretopark.Rule{
				{
					Hours: "Mo-Sa 08:00-18:00",
					Pricing: []wheretopark.PricingRule{
						wheretopark.NewPricingRule("PT1H", decimal.RequireFromString("2.80")),
						wheretopark.NewPricingRule("PT2H", decimal.RequireFromString("6.00")),
						wheretopark.NewPricingRule("PT3H", decimal.RequireFromString("9.00")),
						wheretopark.NewPricingRule("PT4H", decimal.RequireFromString("12.00")),
						wheretopark.NewPricingRule("PT5H", decimal.RequireFromString("15.00")),
						wheretopark.NewPricingRule("PT7H", decimal.RequireFromString("18.00")),
						wheretopark.NewPricingRule("PT9H", decimal.RequireFromString("21.00")),
						wheretopark.NewPricingRule("PT10H", decimal.RequireFromString("29.00")),
					},
				},
				{
					Hours: "Mo-Su 18:00-08:00",
					Pricing: []wheretopark.PricingRule{
						wheretopark.NewPricingRule("PT14H", decimal.RequireFromString("2.00")),
					},
				},
				{
					Hours: "Su 08:00-18:00",
					Pricing: []wheretopark.PricingRule{
						wheretopark.NewPricingRule("PT10H", decimal.RequireFromString("4.00")),
					},
				},
			},
		},

		"CPG02C_1": {
			Name:    "Concert Square Car Park",
			Address: "97 Cowcaddens Rd, Glasgow G4 0DE, United Kingdom",
			MaxDimensions: &wheretopark.Dimensions{
				Height: 190,
			},
			Comment: map[string]string{
				"en": "Customer Notice! There is a partial closure of parking bays on the ramp between floors 1 and 2, this is to assist with ongoing maintenance inspection works. The car park remains open to the public.",
			},
			Rules: []wheretopark.Rule{
				{
					Hours: "Mo-Sa 08:00-18:00",
					Pricing: []wheretopark.PricingRule{
						wheretopark.NewPricingRule("PT1H", decimal.RequireFromString("2.80")),
						wheretopark.NewPricingRule("PT2H", decimal.RequireFromString("6.00")),
						wheretopark.NewPricingRule("PT3H", decimal.RequireFromString("9.00")),
						wheretopark.NewPricingRule("PT4H", decimal.RequireFromString("12.00")),
						wheretopark.NewPricingRule("PT5H", decimal.RequireFromString("15.00")),
						wheretopark.NewPricingRule("PT6H", decimal.RequireFromString("18.00")),
						wheretopark.NewPricingRule("PT7H", decimal.RequireFromString("21.00")),
						wheretopark.NewPricingRule("PT8H", decimal.RequireFromString("24.00")),
						wheretopark.NewPricingRule("PT9H", decimal.RequireFromString("26.50")),
						wheretopark.NewPricingRule("PT10H", decimal.RequireFromString("29.00")),
					},
				},
				{
					Hours: "Mo-Su 18:00-08:00",
					Pricing: []wheretopark.PricingRule{
						wheretopark.NewPricingRule("PT14H", decimal.RequireFromString("2.00")),
					},
				},
				{
					Hours: "Su 08:00-18:00",
					Pricing: []wheretopark.PricingRule{
						wheretopark.NewPricingRule("PT10H", decimal.RequireFromString("4.00")),
					},
				},
			},
		},

		"CPG01C_1": {
			Name:    "Dundasvale Car Park I",
			Address: "18 Dundasvale Ct, Glasgow G4 0SY, United Kingdom",
			MaxDimensions: &wheretopark.Dimensions{
				Height: 190,
			},
			Rules: []wheretopark.Rule{
				{
					Hours: "Mo-Sa 08:00-18:00",
					Pricing: []wheretopark.PricingRule{
						wheretopark.NewPricingRule("PT1H", decimal.RequireFromString("2.20")),
						wheretopark.NewPricingRule("PT2H", decimal.RequireFromString("4.40")),
						wheretopark.NewPricingRule("PT3H", decimal.RequireFromString("6.60")),
						wheretopark.NewPricingRule("PT4H", decimal.RequireFromString("8.80")),
						wheretopark.NewPricingRule("PT5H", decimal.RequireFromString("11.00")),
						wheretopark.NewPricingRule("PT6H", decimal.RequireFromString("13.20")),
						wheretopark.NewPricingRule("PT7H", decimal.RequireFromString("13.20")),
						wheretopark.NewPricingRule("PT10H", decimal.RequireFromString("15.00")),
					},
				},
				{
					Hours: "Mo-Su 18:00-08:00",
					Pricing: []wheretopark.PricingRule{
						wheretopark.NewPricingRule("PT14H", decimal.RequireFromString("2.00")),
					},
				},
				{
					Hours: "Su 08:00-18:00",
					Pricing: []wheretopark.PricingRule{
						wheretopark.NewPricingRule("PT10H", decimal.RequireFromString("4.00")),
					},
				},
			},
		},
	},
}

func init() {
	for k, v := range configuration.ParkingLots {
		v.LastUpdated = metadataLastUpdated
		v.Features = defaultFeatures
		v.Resources = defaultResources
		const SOURCE_NOTICE = "Source of data: glasgow.gov.uk."
		if v.Comment == nil {
			v.Comment = make(map[string]string)
			v.Comment["en"] = SOURCE_NOTICE
		} else {
			v.Comment["en"] += "\n"
			v.Comment["en"] += SOURCE_NOTICE
		}

		configuration.ParkingLots[k] = v
	}
}

var (
	metadataLastUpdated = "2023-03-04"
	defaultFeatures     = []string{
		wheretopark.FeatureCovered,
		wheretopark.FeatureUncovered,
		wheretopark.FeatureMonitored,
	}

	defaultResources = []string{
		"https://www.cityparkingglasgow.co.uk",
		"https://glasgow.gov.uk/index.aspx?articleid=27163",
		"tel:+44-141-276-1830",
		"mailto:carparkinfo@cityparkingglasgow.co.uk",
	}
)
