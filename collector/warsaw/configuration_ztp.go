package warsaw

import (
	"fmt"
	wheretopark "wheretopark/go"
)

type ztpPrices struct {
	// 1h
	h int32
	// 2h
	hh int32
	// 3h
	hhh int32
	// repeating 1 hour
	rh int32
	// repeating 1 month
	rm int32
	// repeating 1 month, day tariff
	rmd int32
	// repeating 1 month, night tariff
	rmn int32
}

func ztpRules(car ztpPrices, truck ztpPrices) []wheretopark.Rule {
	rules := []wheretopark.Rule{
		{
			Hours: "24/7",
			Applies: []string{
				wheretopark.SpotTypeCar,
				wheretopark.SpotTypeCarElectric,
				wheretopark.SpotTypeCarDisabled,
				wheretopark.SpotTypeMotorcycle,
			},
			Pricing: []wheretopark.PricingRule{
				wheretopark.NewIntPricingRule("PT1H", car.h),
				wheretopark.NewIntPricingRule("PT2H", car.hh),
				wheretopark.NewIntPricingRule("PT3H", car.hhh),
				wheretopark.NewIntPricingRule("PT1H", car.rh).Repeated(),
				wheretopark.NewIntPricingRule("PT1M", car.rm).Repeated(),
			},
		},
		{
			Hours: "24/7",
			Applies: []string{
				wheretopark.SpotTypeTruck,
			},
			Pricing: []wheretopark.PricingRule{
				wheretopark.NewIntPricingRule("PT1H", truck.h),
				wheretopark.NewIntPricingRule("PT2H", truck.hh),
				wheretopark.NewIntPricingRule("PT3H", truck.hhh),
				wheretopark.NewIntPricingRule("PT1H", truck.rh).Repeated(),
				wheretopark.NewIntPricingRule("PT1M", truck.rm).Repeated(),
			},
		},
	}
	if car.rmd > 0 {
		rules = append(rules, wheretopark.Rule{
			Hours: "Mo-Su 06:00-18:00",
			Applies: []string{
				wheretopark.SpotTypeCar,
				wheretopark.SpotTypeCarElectric,
				wheretopark.SpotTypeCarDisabled,
				wheretopark.SpotTypeMotorcycle,
			},
			Pricing: []wheretopark.PricingRule{
				wheretopark.NewIntPricingRule("PT1M", car.rmd).Repeated(),
			},
		})
	}
	if car.rmn > 0 {
		rules = append(rules, wheretopark.Rule{
			Hours: "Mo-Su 18:00-06:00",
			Applies: []string{
				wheretopark.SpotTypeCar,
				wheretopark.SpotTypeCarElectric,
				wheretopark.SpotTypeCarDisabled,
				wheretopark.SpotTypeMotorcycle,
			},
			Pricing: []wheretopark.PricingRule{
				wheretopark.NewIntPricingRule("PT1M", car.rmn).Repeated(),
			},
		})
	}

	if truck.rmd > 0 {
		rules = append(rules, wheretopark.Rule{
			Hours: "Mo-Su 06:00-18:00",
			Applies: []string{
				wheretopark.SpotTypeTruck,
			},
			Pricing: []wheretopark.PricingRule{
				wheretopark.NewIntPricingRule("PT1M", truck.rmd).Repeated(),
			},
		})
	}
	if truck.rmn > 0 {
		rules = append(rules, wheretopark.Rule{
			Hours: "Mo-Su 18:00-06:00",
			Applies: []string{
				wheretopark.SpotTypeTruck,
			},
			Pricing: []wheretopark.PricingRule{
				wheretopark.NewIntPricingRule("PT1M", truck.rmn).Repeated(),
			},
		})
	}
	return rules
}

func ztpComment(priceForInhabitants int, inhabitantsOnlyDaily bool) map[string]string {
	comment := make(map[string]string)
	if priceForInhabitants > 0 {
		comment["pl"] = fmt.Sprintf("Dla mieszkancow z adresem zameldowania w bezpośrednim sąsiedzstwie parkingu stosuje sie opłatę abonamentową w wys. %d zł miesięcznie", priceForInhabitants)
		comment["en"] = fmt.Sprintf("For inhabitatns with home address in the direct neighborhood of the parking lot, subscription fee of %d PLN monthly is being applied", priceForInhabitants)
		if inhabitantsOnlyDaily {
			comment["pl"] += ". Tyczy się to wyłącznie samochodów osobowych i motocykli w ofercie całodobowej.\n"
			comment["en"] += ". Applied only for cars and motorcycles in the 24-hour offer.\n"
		} else {
			comment["pl"] += ".\n"
			comment["en"] += ".\n"
		}
	}
	comment["pl"] += "Źródło danych: Miasto Stołeczne Warszawa."
	comment["en"] += "Source of data: Capital City of Warsaw."
	return comment
}

var (
	ztpDefaultFeatures = []wheretopark.Feature{
		wheretopark.FeatureUncovered,
		wheretopark.FeatureGuarded,
	}

	ztpDefaultPaymentMethods = []wheretopark.Feature{
		wheretopark.PaymentMethodCash,
		wheretopark.PaymentMethodCard,
		wheretopark.PaymentMethodContactless,
		wheretopark.PaymentMethodMobile,
	}

	ztpBaseResources = []string{
		"https://www.ztp.waw.pl/strona-111-wykaz_parkingow_strzezonych.html",
		"mailto:nopar@wp.pl",
		"tel:+48-601-248-664",
	}

	ztpParkingLots = map[string]wheretopark.Metadata{
		"u3qcnqcwne": {
			Address: "ul. Bednarska 9, Powiśle, 00-310 Warszawa",
			Resources: append(
				ztpBaseResources,
				"https://www.ztp.waw.pl/portal/download/file_id/164.html",
				"https://www.ztp.waw.pl/portal/download/file_id/165.html",
			),
			Comment: ztpComment(200, false),
			Rules: ztpRules(
				ztpPrices{
					h:   5,
					hh:  11,
					hhh: 14,
					rh:  4,
					rm:  300,
					rmd: 180,
					rmn: 140,
				}, ztpPrices{
					h:   6,
					hh:  13,
					hhh: 18,
					rh:  4,
					rm:  330,
					rmd: 240,
					rmn: 180,
				},
			),

			LastUpdated:    &defaultLastUpdated,
			Features:       ztpDefaultFeatures,
			PaymentMethods: ztpDefaultPaymentMethods,
			Currency:       defaultCurrency,
			Timezone:       defaultTimezone,
		},

		wheretopark.CoordinateToID(52.252094, 21.012700): {
			Address: "ul. Boleść 6, Nowe Miasto, 00-259 Warszawa",
			Resources: append(
				ztpBaseResources,
				"https://www.ztp.waw.pl/portal/download/file_id/166.html",
				"https://www.ztp.waw.pl/portal/download/file_id/167.html",
			),
			Comment: ztpComment(200, false),
			Rules: ztpRules(
				ztpPrices{
					h:   5,
					hh:  11,
					hhh: 14,
					rh:  4,
					rm:  300,
					rmd: 180,
					rmn: 140,
				}, ztpPrices{
					h:   6,
					hh:  13,
					hhh: 18,
					rh:  4,
					rm:  330,
					rmd: 240,
					rmn: 180,
				},
			),

			LastUpdated:    &defaultLastUpdated,
			Features:       ztpDefaultFeatures,
			PaymentMethods: ztpDefaultPaymentMethods,
			Currency:       defaultCurrency,
			Timezone:       defaultTimezone,
		},

		"u3qcq0n7c9": {
			Address: "ul. Bugaj, Śródmieście, 00-284 Warszawa",
			Resources: append(
				ztpBaseResources,
				"https://www.ztp.waw.pl/portal/download/file_id/168.html",
				"https://www.ztp.waw.pl/portal/download/file_id/169.html",
			),
			Comment: ztpComment(0, false),
			Rules: ztpRules(
				ztpPrices{
					h:   5,
					hh:  10,
					hhh: 14,
					rh:  4,
					rm:  200,
					rmd: 150,
				}, ztpPrices{
					h:   7,
					hh:  12,
					hhh: 17,
					rh:  5,
					rm:  400,
				},
			),

			LastUpdated:    &defaultLastUpdated,
			Features:       ztpDefaultFeatures,
			PaymentMethods: ztpDefaultPaymentMethods,
			Currency:       defaultCurrency,
			Timezone:       defaultTimezone,
		},

		"u3qcn41g7d": {
			Address: "ul. Filtrowa 1, Śródmieście Południowe, 00-611 Warszawa",
			Resources: append(
				ztpBaseResources,
				"https://www.ztp.waw.pl/portal/download/file_id/170.html",
				"https://www.ztp.waw.pl/portal/download/file_id/171.html",
			),
			Comment: ztpComment(250, true),
			Rules: ztpRules(
				ztpPrices{
					h:   5,
					hh:  11,
					hhh: 14,
					rh:  4,
					rm:  370,
					rmd: 270,
					rmn: 200,
				}, ztpPrices{
					h:   6,
					hh:  13,
					hhh: 18,
					rh:  4,
					rm:  460,
					rmd: 340,
					rmn: 250,
				},
			),

			LastUpdated:    &defaultLastUpdated,
			Features:       ztpDefaultFeatures,
			PaymentMethods: ztpDefaultPaymentMethods,
			Currency:       defaultCurrency,
			Timezone:       defaultTimezone,
		},

		"u3qcnm9txg": {
			Address: "ul. Karasia, Śródmieście Północne, 00-327 Warszawa",
			Resources: append(
				ztpBaseResources,
				"https://www.ztp.waw.pl/portal/download/file_id/174.html",
				"https://www.ztp.waw.pl/portal/download/file_id/175.html",
			),
			Comment: ztpComment(350, false),
			Rules: ztpRules(
				ztpPrices{
					h:   5,
					hh:  11,
					hhh: 14,
					rh:  4,
					rm:  470,
					rmd: 370,
					rmn: 230,
				}, ztpPrices{
					h:   6,
					hh:  13,
					hhh: 18,
					rh:  4,
					rm:  560,
					rmd: 420,
					rmn: 300,
				},
			),

			LastUpdated:    &defaultLastUpdated,
			Features:       ztpDefaultFeatures,
			PaymentMethods: ztpDefaultPaymentMethods,
			Currency:       defaultCurrency,
			Timezone:       defaultTimezone,
		},

		"u3qcm8j3ep": {
			Address: "ul. Miła, Muranów, 03-402 Warszawa",
			Resources: append(
				ztpBaseResources,
				"https://www.ztp.waw.pl/portal/download/file_id/176.html",
				"https://www.ztp.waw.pl/portal/download/file_id/177.html",
			),
			Comment: ztpComment(0, false),
			Rules: ztpRules(
				ztpPrices{
					h:   4,
					hh:  9,
					hhh: 12,
					rh:  3,
					rm:  200,
					rmd: 160,
					rmn: 140,
				}, ztpPrices{
					h:   5,
					hh:  11,
					hhh: 15,
					rh:  4,
					rm:  230,
					rmd: 200,
					rmn: 160,
				},
			),

			LastUpdated:    &defaultLastUpdated,
			Features:       ztpDefaultFeatures,
			PaymentMethods: ztpDefaultPaymentMethods,
			Currency:       defaultCurrency,
			Timezone:       defaultTimezone,
		},

		"u3qcndw3tr": {
			// also called Myśliwiecka
			Address: "ul. Łazienkowska 7, Śródmieście, 00-459 Warszawa",
			Resources: append(
				ztpBaseResources,
				"https://www.ztp.waw.pl/portal/download/file_id/178.html",
				"https://www.ztp.waw.pl/portal/download/file_id/179.html",
			),
			Comment: ztpComment(0, false),
			Rules: ztpRules(
				ztpPrices{
					h:   4,
					hh:  9,
					hhh: 12,
					rh:  3,
					rm:  250,
					rmd: 200,
				}, ztpPrices{
					h:   5,
					hh:  11,
					hhh: 15,
					rh:  4,
					rm:  300,
					rmd: 250,
				},
			),

			LastUpdated:    &defaultLastUpdated,
			Features:       ztpDefaultFeatures,
			PaymentMethods: ztpDefaultPaymentMethods,
			Currency:       defaultCurrency,
			Timezone:       defaultTimezone,
		},

		"u3qcnhgzws": {
			Address: "Plac Młynarskiego, Śródmieście Północne, 00-281 Warszawa",
			Resources: append(
				ztpBaseResources,
				"https://www.ztp.waw.pl/portal/download/file_id/180.html",
				"https://www.ztp.waw.pl/portal/download/file_id/181.html",
			),
			Comment: ztpComment(250, false),
			Rules: ztpRules(
				ztpPrices{
					h:   5,
					hh:  11,
					hhh: 14,
					rh:  4,
					rm:  470,
					rmd: 370,
					rmn: 230,
				}, ztpPrices{
					h:   6,
					hh:  13,
					hhh: 18,
					rh:  4,
					rm:  560,
					rmd: 420,
					rmn: 310,
				},
			),

			LastUpdated:    &defaultLastUpdated,
			Features:       ztpDefaultFeatures,
			PaymentMethods: ztpDefaultPaymentMethods,
			Currency:       defaultCurrency,
			Timezone:       defaultTimezone,
		},

		"u3qcjvvj2z": {
			Address: "Plac Żelaznej Bramy 10, Śródmieście Północne, 01-136 Warszawa",
			Resources: append(
				ztpBaseResources,
				"https://www.ztp.waw.pl/portal/download/file_id/182.html",
				"https://www.ztp.waw.pl/portal/download/file_id/183.html",
			),
			Comment: ztpComment(250, true),
			Rules: ztpRules(
				ztpPrices{
					h:   5,
					hh:  11,
					hhh: 14,
					rh:  4,
					rm:  470,
					rmd: 200,
					rmn: 140,
				}, ztpPrices{
					h:   6,
					hh:  13,
					hhh: 18,
					rh:  4,
					rm:  560,
					rmd: 420,
					rmn: 310,
				},
			),

			LastUpdated:    &defaultLastUpdated,
			Features:       ztpDefaultFeatures,
			PaymentMethods: ztpDefaultPaymentMethods,
			Currency:       defaultCurrency,
			Timezone:       defaultTimezone,
		},

		"u3qcntjc55": {
			Address: "aleja 3 Maja, 00-401 Warszawa",
			Resources: append(
				ztpBaseResources,
				"https://www.ztp.waw.pl/portal/download/file_id/185.html",
				"https://www.ztp.waw.pl/portal/download/file_id/186.html",
			),
			// TODO: add info about guaranteed 25 parking slots
			Comment: ztpComment(200, true),
			Rules: ztpRules(
				ztpPrices{
					h:   5,
					hh:  11,
					hhh: 14,
					rh:  4,
					rm:  370,
				}, ztpPrices{
					h:   6,
					hh:  13,
					hhh: 18,
					rh:  4,
					rm:  560,
				},
			),

			LastUpdated:    &defaultLastUpdated,
			Features:       ztpDefaultFeatures,
			PaymentMethods: ztpDefaultPaymentMethods,
			Currency:       defaultCurrency,
			Timezone:       defaultTimezone,
		},

		"u3qcnkxqwd": {
			// przy stacji PKP Powiśle
			Address: "Smolna, Śródmieście, 00-375 Warszawa",
			Resources: append(
				ztpBaseResources,
				"https://www.ztp.waw.pl/portal/download/file_id/187.html",
				"https://www.ztp.waw.pl/portal/download/file_id/188.html",
			),
			Comment: ztpComment(250, false),
			Rules: ztpRules(
				ztpPrices{
					h:   5,
					hh:  11,
					hhh: 14,
					rh:  4,
					rm:  370,
				}, ztpPrices{
					h:   6,
					hh:  13,
					hhh: 18,
					rh:  4,
					rm:  560,
				},
			),

			LastUpdated:    &defaultLastUpdated,
			Features:       ztpDefaultFeatures,
			PaymentMethods: ztpDefaultPaymentMethods,
			Currency:       defaultCurrency,
			Timezone:       defaultTimezone,
		},

		"u3qcn5wfrm": {
			Address: "ul. Hoża, Śródmieście Południowe, 00-682 Warszawa",
			Resources: append(
				ztpBaseResources,
				"https://www.ztp.waw.pl/portal/download/file_id/189.html",
				"https://www.ztp.waw.pl/portal/download/file_id/190.html",
			),
			Comment: ztpComment(250, false),
			Rules: ztpRules(
				ztpPrices{
					h:   5,
					hh:  11,
					hhh: 14,
					rh:  4,
					rm:  470,
					rmd: 370,
					rmn: 230,
				}, ztpPrices{
					h:   6,
					hh:  13,
					hhh: 18,
					rh:  4,
					rm:  560,
					rmd: 420,
					rmn: 310,
				},
			),

			LastUpdated:    &defaultLastUpdated,
			Features:       ztpDefaultFeatures,
			PaymentMethods: ztpDefaultPaymentMethods,
			Currency:       defaultCurrency,
			Timezone:       defaultTimezone,
		},

		"u3qcjvuk6p": {
			Address: "ul. Ptasia, 03-402 Warszaw",
			Resources: append(
				ztpBaseResources,
				"https://www.ztp.waw.pl/portal/download/file_id/621.html",
				"https://www.ztp.waw.pl/portal/download/file_id/622.html",
			),
			Comment: ztpComment(250, false),
			Rules: ztpRules(
				ztpPrices{
					h:   5,
					hh:  11,
					hhh: 14,
					rh:  4,
					rm:  470,
					rmd: 370,
					rmn: 230,
				}, ztpPrices{
					h:   6,
					hh:  13,
					hhh: 18,
					rh:  4,
					rm:  560,
					rmd: 420,
					rmn: 310,
				},
			),

			LastUpdated:    &defaultLastUpdated,
			Features:       ztpDefaultFeatures,
			PaymentMethods: ztpDefaultPaymentMethods,
			Currency:       defaultCurrency,
			Timezone:       defaultTimezone,
		},
	}
)
