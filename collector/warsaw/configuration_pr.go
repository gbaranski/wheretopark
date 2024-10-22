package warsaw

import (
	wheretopark "wheretopark/go"
)

var (
	prBaseResources = []string{
		"mailto:ztm@ztm.waw.pl",
		"tel:+48-(22)-56-98-116",
		"tel:+48-(22)-56-98-11",
	}

	prRules24h = []wheretopark.Rule{
		{
			Hours:   "24/7",
			Pricing: []wheretopark.PricingRule{},
		},
	}

	prRulesDay = []wheretopark.Rule{
		{
			Hours:   "Mo-Su 04:30-02:30",
			Pricing: []wheretopark.PricingRule{},
		},
	}

	prDefaultComment = map[string]string{
		"pl": `Parkingi przeznaczone są dla samochodów osobowych, motocykli, rowerów i motorowerów.
System parkingów „Parkuj i Jedź” (Park & Ride) umożliwia bezpłatne parkowanie pojazdów osobom, które w chwili wyjazdu z parkingu przedstawią ważny bilet:
  - dobowy,
  - 3-dniowy,
  - weekendowy,
  - weekendowy grupowy,
  - 30-dniowy,
  - 90-dniowy,
  - bilet seniora,
  - bilet dla dzieci z rodzin z trójką dzieci
  - dokument uprawniający do bezpłatnych przejazdów środkami lokalnego transportu zbiorowego organizowanego przez m.st. Warszawę.

W innym przypadku użytkownik jest zobowiązany do uiszczenia przy wyjeździe z parkingu jednorazowej opłaty za wynajem miejsca parkingowego w wysokości 100 zł.

Opłata zaczyna być naliczana po upływie 20 minut od wjazdu na parking.

Opłaty nie pobiera się za:
  - postój rowerów, motorowerów i motocykli,
  - postój pojazdów Zarządcy oraz obsługi parkingu podczas wykonywania czynności eksploatacyjnych,
  - ładowanie pojazdów elektrycznych.

Kontrola uprawnień odbywa się wyrywkowo przez kontrolerów przy wyjeździe na zasadach kontroli w pojazdach komunikacji miejskiej.

Źródło danych: Miasto Stołeczne Warszawa.`,

		"en": `Parking lots are intended for cars, motorcycles, bicycles and mopeds.
The "Park and Ride" (Park & Ride) parking system allows free parking for vehicles that present a valid ticket at the time of departure from the car park:
  - daily,
  - 3-day,
  - weekend,
  - group weekend,
  - 30-day,
  - 90-day,
  - senior ticket,
  - ticket for children from families with three children
  - document entitling to free rides with local public transport organized by the City of Warsaw.

In other cases, the user is obliged to pay a one-time fee for renting a parking space upon departure from the car park in the amount of 100 PLN.

The fee starts to be charged after 20 minutes from entering the car park.

Fees are not charged for:
  - parking bicycles, mopeds and motorcycles,
  - parking vehicles of the Manager and parking staff during the performance of maintenance activities,
  - charging electric vehicles.

The right to use is checked randomly by controllers upon departure on the basis of control in public transport vehicles.

Data source: City of Warsaw.`,
	}

	prParkingLots = map[string]wheretopark.Metadata{
		"u3qckr1s6u": {
			Address:   "ul. Kasprowicza, Wrzeciono, 01-949 Warszawa",
			Resources: append(prBaseResources, "https://www.wtp.waw.pl/parkingi/parking-pr-metro-mlociny-iii/"),
			Features: []wheretopark.Feature{
				wheretopark.FeatureUncovered,
			},

			LastUpdated: &defaultLastUpdated,
			Comment:     prDefaultComment,
			Currency:    defaultCurrency,
			Timezone:    defaultTimezone,
			Rules:       prRulesDay,
		},

		"u3qbmy29nb": {
			Address:   "ul. Karczunkowa 145, 02-873 Warszawa",
			Resources: append(prBaseResources, "https://www.wtp.waw.pl/parkingi/parking-pr-jeziorki-pkp/"),
			Features: []wheretopark.Feature{
				wheretopark.FeatureUncovered,
			},

			LastUpdated: &defaultLastUpdated,
			Comment:     prDefaultComment,
			Currency:    defaultCurrency,
			Timezone:    defaultTimezone,
			Rules:       prRulesDay,
		},

		"u3qbud4y0w": {
			Address:   "Al. Krakowska, Załuski, 02-180 Warszawa",
			Resources: append(prBaseResources, "https://www.wtp.waw.pl/parkingi/parking-pr-al-krakowska"),
			Features: []wheretopark.Feature{
				wheretopark.FeatureCovered,
			},
			MaxDimensions: &wheretopark.Dimensions{
				Height: 190,
			},

			LastUpdated: &defaultLastUpdated,
			Comment:     prDefaultComment,
			Currency:    defaultCurrency,
			Timezone:    defaultTimezone,
			Rules:       prRulesDay,
		},

		"u3qbwwr4dz": {
			Address:   "ul. Ciszewskiego, Ursynów Północny, 02-777 Warszawa",
			Resources: append(prBaseResources, "https://www.wtp.waw.pl/parkingi/parking-pr-metro-stoklosy"),
			Features: []wheretopark.Feature{
				wheretopark.FeatureCovered,
			},
			MaxDimensions: &wheretopark.Dimensions{
				Height: 200,
			},

			LastUpdated: &defaultLastUpdated,
			Comment:     prDefaultComment,
			Currency:    defaultCurrency,
			Timezone:    defaultTimezone,
			Rules:       prRulesDay,
		},

		"u3r1190t8e": {
			Address:   "ul. Szpotańskiego, Anin, 04-610 Warszawa",
			Resources: append(prBaseResources, "https://www.wtp.waw.pl/parkingi/parking-pr-anin-skm"),
			Features: []wheretopark.Feature{
				wheretopark.FeatureUncovered,
			},
			MaxDimensions: &wheretopark.Dimensions{
				Height: 220,
			},

			LastUpdated: &defaultLastUpdated,
			Comment:     prDefaultComment,
			Currency:    defaultCurrency,
			Timezone:    defaultTimezone,
			Rules:       prRulesDay,
		},

		"u3qch4g0rr": {
			Address:   "ul. Półczyńska 8, Bemowo, 01-378 Warszawa",
			Resources: append(prBaseResources, "https://www.wtp.waw.pl/parkingi/parking-pr-polczynska"),
			Features: []wheretopark.Feature{
				wheretopark.FeatureUncovered,
			},

			LastUpdated: &defaultLastUpdated,
			Comment:     prDefaultComment,
			Currency:    defaultCurrency,
			Timezone:    defaultTimezone,
			Rules:       prRulesDay,
		},

		"u3qckr8487": {
			Address:   "ul. Pułku Strzelców Kaniowskich, Huta, 01-949 Warszawa",
			Resources: append(prBaseResources, "https://www.wtp.waw.pl/parkingi/parking-pr-metro-mlociny-ii/"),
			Features: []wheretopark.Feature{
				wheretopark.FeatureUncovered,
			},
			Rules: prRules24h,

			LastUpdated: &defaultLastUpdated,
			Comment:     prDefaultComment,
			Currency:    defaultCurrency,
			Timezone:    defaultTimezone,
		},

		"u3qby7jfb6": {
			Address:   "al. Wilanowska 236, Ksawerów, 02-765 Warszawa",
			Resources: append(prBaseResources, "https://www.wtp.waw.pl/parkingi/parking-pr-metro-wilanowska/"),

			LastUpdated: &defaultLastUpdated,
			Comment:     prDefaultComment,
			Currency:    defaultCurrency,
			Timezone:    defaultTimezone,
			Rules:       prRulesDay,
		},

		"u3qbwrztn7": {
			Address:   "al. Komisji Edukacji Narodowej, Ursynów Północny, 02-787 Warszawa",
			Resources: append(prBaseResources, "https://www.wtp.waw.pl/parkingi/parking-pr-metro-ursynow/"),
			MaxDimensions: &wheretopark.Dimensions{
				Height: 220,
			},
			Features: []wheretopark.Feature{
				wheretopark.FeatureUncovered,
			},

			LastUpdated: &defaultLastUpdated,
			Comment:     prDefaultComment,
			Currency:    defaultCurrency,
			Timezone:    defaultTimezone,
			Rules:       prRulesDay,
		},

		"u3r116bcdu": {
			Address:   "ul. Widoczna 2a, Wawer, 04-647 Warszawa",
			Resources: append(prBaseResources, "https://www.wtp.waw.pl/parkingi/parking-pr-wawer-skm/"),
			Features: []wheretopark.Feature{
				wheretopark.FeatureUncovered,
			},

			LastUpdated: &defaultLastUpdated,
			Comment:     prDefaultComment,
			Currency:    defaultCurrency,
			Timezone:    defaultTimezone,
			Rules:       prRulesDay,
		},

		"u3qctgenxn": {
			Address:   "ul. Marywilska, Białołęka, 03-042 Warszawa",
			Resources: append(prBaseResources, "https://www.wtp.waw.pl/parkingi/parking-pr-zeran-pkp/"),
			Features: []wheretopark.Feature{
				wheretopark.FeatureUncovered,
			},

			LastUpdated: &defaultLastUpdated,
			Comment:     prDefaultComment,
			Currency:    defaultCurrency,
			Timezone:    defaultTimezone,
			Rules:       prRules24h,
		},

		wheretopark.CoordinateToID(52.271140, 20.970460): {
			Address:   "ul. Włościańska 56, 01-710 Warszawa",
			Resources: append(prBaseResources, "https://www.wtp.waw.pl/parkingi/parking-pr-metro-marymont/"),
			MaxDimensions: &wheretopark.Dimensions{
				Height: 180,
			},
			Features: []wheretopark.Feature{
				wheretopark.FeatureCovered,
			},

			LastUpdated: &defaultLastUpdated,
			Comment:     prDefaultComment,
			Currency:    defaultCurrency,
			Timezone:    defaultTimezone,
			Rules:       prRulesDay,
		},

		wheretopark.CoordinateToID(52.191143, 20.867516): {
			Address:   "ul. Orląt Lwowskich 45, 02-495 Warszawa",
			Resources: append(prBaseResources, "https://www.wtp.waw.pl/parkingi/parking-pr-ursus-niedzwiadek/"),
			MaxDimensions: &wheretopark.Dimensions{
				Height: 200,
			},
			Features: []wheretopark.Feature{
				wheretopark.FeatureCovered,
			},

			LastUpdated: &defaultLastUpdated,
			Comment:     prDefaultComment,
			Currency:    defaultCurrency,
			Timezone:    defaultTimezone,
			Rules:       prRulesDay,
		},
	}
)
