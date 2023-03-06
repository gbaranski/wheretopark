package poznan

import (
	wheretopark "wheretopark/go"

	_ "embed"

	geojson "github.com/paulmach/go.geojson"
)

type Configuration struct {
	ParkingLots map[string]wheretopark.Metadata
}

var configuration = Configuration{
	ParkingLots: map[string]wheretopark.Metadata{
		"swmichala": {
			Name:     "Park & Ride ul. Św. Michała",
			Address:  "ul. Świętego Michała, 61-113 Poznan",
			Geometry: *geojson.NewPointGeometry([]float64{16.962210812192566, 52.408434016240186}),
			Resources: []string{
				"https://www.apcoa.pl/parking-w/poznan/ztm-park-ride-poznan-ul-sw-michala/",
			},
			TotalSpots: map[string]uint{
				wheretopark.SpotTypeCar: 100,
			},
		},

		"biskupinska": {
			Name:     "Park & Ride ul. Biskupińska",
			Address:  "ul. Biskupińska, 60-463 Poznań",
			Geometry: *geojson.NewPointGeometry([]float64{16.864887711571242, 52.46070378280677}),
			Resources: []string{
				"https://www.apcoa.pl/en/parking-in/poznan/ztm-park-ride-poznan-ul-biskupinska/",
			},
			TotalSpots: map[string]uint{
				wheretopark.SpotTypeCar: 50,
			},
		},

		"rondo_staroleka": {
			Name:     "Park & Ride Rondo Starołęka",
			Address:  "ul. Wągrowska, 61-001 Poznań",
			Geometry: *geojson.NewPointGeometry([]float64{16.94430210187254, 52.37916546114334}),
			Resources: []string{
				"https://www.apcoa.pl/en/parking-in/poznan/ztm-park-ride-staroleka-poznan-ul-wagrowska/",
			},
			TotalSpots: map[string]uint{
				wheretopark.SpotTypeCar: 57,
			},
		},

		"szymanowskiego": {
			Name:     "Park & Ride ul. Szymanowskiego",
			Address:  "ul. Szymanowskiego, 61-995 Poznań",
			Geometry: *geojson.NewPointGeometry([]float64{16.916001501951616, 52.46207982767391}),
			Resources: []string{
				"https://www.apcoa.pl/en/parking-in/poznan/ztm-park-ride-poznan-ul-szymanowskiego/",
			},
			TotalSpots: map[string]uint{
				wheretopark.SpotTypeCar: 130,
			},
		},
	},
}

func init() {
	for k, v := range configuration.ParkingLots {
		v.LastUpdated = defaultLastUpdated
		v.MaxDimensions = defaultMaxDimensions
		v.Features = defaultFeatures
		v.Resources = append(baseResources, v.Resources...)
		v.PaymentMethods = defaultPaymentMethods
		v.Comment = defaultComment
		v.Currency = defaultCurrency
		v.Timezone = defaultTimezone
		v.Rules = defaultRules
		configuration.ParkingLots[k] = v
	}
}

var (
	defaultLastUpdated   = "2023-03-06"
	defaultMaxDimensions = &wheretopark.Dimensions{
		Height: -1,
	}
	defaultFeatures = []string{
		wheretopark.FeatureUncovered,
	}
	baseResources = []string{
		"mailto:parkingi@apcoa.pl",
		"tel:+48-22-354-83-80",
		"https://www.ztm.poznan.pl/en/komunikacja/parkuj-i-jedz/",
	}
	defaultPaymentMethods = []string{
		wheretopark.PaymentMethodCash,
		wheretopark.PaymentMethodCard,
		wheretopark.PaymentMethodContactless,
		wheretopark.PaymentMethodMobile,
	}
	defaultComment = map[string]string{
		"pl": `Taryfa Zarządu Transportu Miejskiego w Poznaniu zachęca do łączenia przejazdów samochodem i komunikacją miejską – dzięki biletowi okresowemu na karcie PEKA z parkingów P&R można korzystać bezpłatnie. Osoby nie korzystające z biletu okresowego są zobowiązane wykupić bilet w cenie 10 zł, na podstawie którego można podróżować na liniach ZTM w godzinach działania parkingów.
				W przypadku pozostawienia pojazdu na noc(po godzinie 02:30) obowiązuje opłata dodatkowa 100 PLN.
				`,
	}
	defaultCurrency = "PLN"
	defaultTimezone = "Europe/Warsaw"
	defaultRules    = []wheretopark.Rule{
		{
			Hours: "Mo-Su 04:30-02:30",
			Pricing: []wheretopark.PricingRule{
				wheretopark.NewIntPricingRule("P1D", 10),
			},
		},
	}
)
