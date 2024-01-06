package krakow

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	wheretopark "wheretopark/go"
	"wheretopark/meters"

	geojson "github.com/paulmach/go.geojson"
	"golang.org/x/text/currency"
	"golang.org/x/text/encoding/charmap"
)

type Placemark struct {
	Subzone     string `xml:"name"`
	Card        string `xml:"card"`
	Model       string `xml:"model"`
	Code        string `xml:"parkingmeter"`
	Address     string `xml:"address"`
	Coordinates struct {
		Latitude  float64 `xml:"latitude"`
		Longitude float64 `xml:"longitude"`
	} `xml:"coordinates"`
}

type Folder struct {
	Placemarks []Placemark `xml:"placemark"`
}

var (
	PLACEMARKS_URL = wheretopark.MustParseURL("http://zdmk.krakow.pl/wp-content/themes/justidea_theme/assets/xml/parkomaty.xml")
)

func GetPlacemarks() ([]Placemark, error) {
	folder, err := wheretopark.Get[Folder](PLACEMARKS_URL, nil)
	if err != nil {
		return nil, err
	}
	return folder.Placemarks, nil
}

func CodeMapping(placemarks []Placemark) map[meters.Code]wheretopark.ID {
	ids := make(map[meters.Code]wheretopark.ID, len(placemarks))
	for _, placemark := range placemarks {
		id := wheretopark.CoordinateToID(placemark.Coordinates.Latitude, placemark.Coordinates.Longitude)
		ids[placemark.Code] = id
	}
	return ids
}

func LoadPlacemarks(bytes []byte) ([]Placemark, error) {
	var metadata Folder
	if err := xml.Unmarshal(bytes, &metadata); err != nil {
		return nil, err
	}
	return metadata.Placemarks, nil
}

func LoadPlacemarksFile(path string) ([]Placemark, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	byteValue, err := io.ReadAll(charmap.ISO8859_2.NewDecoder().Reader(file))
	if err != nil {
		return nil, err
	}
	return LoadPlacemarks(byteValue)
}

func newRules(perHourPrice int) []wheretopark.Rule {
	return []wheretopark.Rule{
		{
			Hours: "Mo-Sa 10:00-20:00",
			Pricing: []wheretopark.PricingRule{
				wheretopark.NewIntPricingRule("PT1H", 6).Repeated(),
			},
		},
		{
			Hours: "Mo-Sa 20:00-10:00",
			Pricing: []wheretopark.PricingRule{
				wheretopark.NewIntPricingRule("PT1H", 0).Repeated(),
			},
		},
		{
			Hours: "Su",
			Pricing: []wheretopark.PricingRule{
				wheretopark.NewIntPricingRule("PT1H", 0).Repeated(),
			},
		},
		{
			Hours: "24/7",
			Applies: []string{
				wheretopark.SpotTypeCarDisabled,
			},
			Pricing: []wheretopark.PricingRule{
				wheretopark.NewIntPricingRule("P1M", 0).Repeated(),
			},
		},
	}

}

func (p Placemark) ZoneMatchingRules() []wheretopark.Rule {
	zone := []rune(p.Subzone)[7]
	switch zone {
	case 'A':
		return newRules(6)
	case 'B':
		return newRules(5)
	case 'C':
		return newRules(4)
	default:
		panic(fmt.Sprintf("unknown zone %s", p.Subzone))
	}
}

func (p Placemark) Metadata(totalSpots uint) wheretopark.Metadata {
	metadata := wheretopark.Metadata{
		LastUpdated: &defaultLastUpdated,
		Name:        fmt.Sprintf("Parking nr. %s", p.Code),
		Address:     p.Address,
		Geometry:    geojson.NewPointGeometry([]float64{p.Coordinates.Longitude, p.Coordinates.Latitude}),
		Resources: []string{
			"https://zdmk.krakow.pl/parkowanie/strefa-platnego-parkowania/informacje-ogolne-i-oplaty/",
			"tel:+48-12-616-71-77",
			"mailto:strefa@zdmk.krakow.pl",
		},
		TotalSpots: map[string]uint{
			wheretopark.SpotTypeCar: totalSpots,
		},
		MaxDimensions: nil,
		Features: []string{
			wheretopark.FeatureUncovered,
		},
		PaymentMethods: []string{
			wheretopark.PaymentMethodCash,
			wheretopark.PaymentMethodMobile,
		},
		Comment: map[string]string{
			"pl": `
Abonament można nabyć w biurze Strefy Płatnego Parkowania (ul. W. Reymonta 20) lub w sklepie internetowym eabonament.zdmk.pl.
Za pomocą aplikacji na telefonie można zakupić bilet, umożliwiają to: www.skycash.com, www.anypark.pl, www.mobilet.pl, www.mpay.pl, www.electronicparking.pl, www.cityparkapp.pl, www.flowbird.pl, www.mka.malopolska.pl.
Osoba posiadająca zagraniczny numer telefonu może uiścić opłatę za postój pojazdu w OPP korzystając z następujących aplikacji: www.mobilet.pl (tutaj rejestracja możliwa za pomocą adresu e-mail), www.electronicparking.pl, www.mpay.pl, www.cityparkapp.pl, www.anypark.pl.`,
		},
		Currency: defaultCurrency,
		Timezone: defaultTimezone,
		Rules:    p.ZoneMatchingRules(),
	}
	if p.Card == "zbliżeniowa" {
		metadata.PaymentMethods = append(metadata.PaymentMethods, wheretopark.PaymentMethodCard, wheretopark.PaymentMethodContactless)
	}
	return metadata
}

var (
	defaultLastUpdated = wheretopark.MustParseDate("2024-01-04")
	defaultTimezone    = wheretopark.MustLoadLocation("Europe/Warsaw")
	defaultCurrency    = currency.PLN
)
