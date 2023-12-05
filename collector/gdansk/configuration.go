package gdansk

import (
	wheretopark "wheretopark/go"

	_ "embed"

	"golang.org/x/text/currency"
)

type Configuration struct {
	ParkingLots map[string]wheretopark.Metadata
}

var configuration = Configuration{
	ParkingLots: map[string]wheretopark.Metadata{
		"1": {
			Resources: []string{
				"mailto:galeria@galeriabaltycka.pl",
				"tel:+48-58-521-85-52",
				"https://www.galeriabaltycka.pl/o-centrum/dojazd-parkingi/parkingi",
			},
			TotalSpots: map[string]uint{
				wheretopark.SpotTypeCar: 1000,
			},
			MaxDimensions: &wheretopark.Dimensions{
				Height: 190,
			},
			Features: []string{
				wheretopark.FeatureCovered,
				wheretopark.FeatureUncovered,
			},
			Comment: map[string]string{
				"pl": `Na dwóch najwyższych kondygnacjach budynku centrum handlowego oferujemy dwupoziomowy parking i 1100 miejsc postojowych. 
Wjazd do centrum handlowego odbywa się z ronda od strony ulicy Dmowskiego w Gdańsku. 
Komunikację między poziomami parkingowymi a poziomami handlowymi centrum handlowego zapewniają schody ruchome i windy szybkobieżne.
Prosimy o zachowanie biletu parkingowego i opłacenie należności za postój w kasie automatycznej, znajdującej się przy wyjściu z parkingu.`,

				"en": `We have prepared a two-level car park with 1,100 parking spaces (including those for disabled people) for our clients.
It is situated on the two top floors of the building. 
You can get there driving from the roundabout from the direction of Dmowskiego Street. 
Both levels of the car park can be reached by a spiral parking ramp. 
Escalators and high-speed lifts will take you from the car park decks to the Gallery's floors and back.`,

				"ru": `На двух верхних этажах здания торгового центра расположен двухуровневый паркинг на 1100 парковочных мест.
В торговый центр можно попасть с кругового перекрестка со стороны улицы Дмовскего в Гданьске.
Сообщение между уровнями паркинга и торговыми уровнями центра обеспечивают эскалаторы и скоростные лифты.`,
			},
			Rules: []wheretopark.Rule{
				{
					Hours: "Mo-Sa 08:00-22:00; Su 09:00-21:00",
					Pricing: []wheretopark.PricingRule{
						wheretopark.NewIntPricingRule("PT1H", 0),
						wheretopark.NewIntPricingRule("PT2H", 2),
						wheretopark.NewIntPricingRule("PT3H", 3),
						wheretopark.NewIntPricingRule("P1D", 25),
						wheretopark.NewIntPricingRule("PT1H", 4).Repeated(),
					},
				},
			},
		},

		"2": {
			Resources: []string{
				"mailto:biuro@gchmanhattan.pl",
				"tel:+48-58-767-70-16",
				"https://gchmanhattan.pl/o-centrum/",
			},
			TotalSpots: map[string]uint{
				wheretopark.SpotTypeCar: 360,
			},
			Features: []string{
				wheretopark.FeatureUnderground,
			},
			PaymentMethods: []string{
				wheretopark.PaymentMethodCash,
				wheretopark.PaymentMethodCard,
				wheretopark.PaymentMethodContactless,
			},
			Comment: map[string]string{
				"pl": "Dwupoziomowy parking rotacyjny",
			},
			Rules: []wheretopark.Rule{
				{
					Hours: "Mo-Th 06:00-21:00",
					Pricing: []wheretopark.PricingRule{
						wheretopark.NewIntPricingRule("PT30M", 0),
						wheretopark.NewIntPricingRule("PT1H", 4).Repeated(),
					},
				},
				{
					Hours: "Fr 06:00-24:00",
					Pricing: []wheretopark.PricingRule{
						wheretopark.NewIntPricingRule("PT30M", 0),
						wheretopark.NewIntPricingRule("PT1H", 4).Repeated(),
					},
				},
				{
					Hours: "Sa 08:00-04:00",
					Pricing: []wheretopark.PricingRule{
						wheretopark.NewIntPricingRule("PT1H", 0),
						wheretopark.NewIntPricingRule("PT1H", 4).Repeated(),
					},
				},
				{
					Hours: "Su 08:00-22:00",
					Pricing: []wheretopark.PricingRule{
						wheretopark.NewIntPricingRule("PT1H", 0).Repeated(),
					},
				},
			},
		},

		"3": {
			Resources: []string{
				"mailto:parkingi@apcoa.pl",
				"tel:+48-22-354-83-80",
				"https://www.apcoa.pl/parking-w/gdansk-1/polsat-plus-arena-gdansk-p4-ul-pok-lechii-gdansk/",
				"https://www.apcoa.pl/parking-w/gdansk-1/polsat-plus-arena-gdansk-p3-ul-pok-lechii-gdansk/",
			},
			TotalSpots: map[string]uint{
				wheretopark.SpotTypeCar:         2440,
				wheretopark.SpotTypeCarDisabled: 60,
			},
			MaxDimensions: &wheretopark.Dimensions{
				Height: -1,
			},
			Features: []string{
				wheretopark.FeatureUncovered,
			},
			Comment: map[string]string{
				"pl": `Parking naziemny zlokalizowany przy stadionie Energa w Gdańsku.
Do dyspozycji klientów jest dostępnych 2500 miejsc parkingowych, w tym 40 bezpłatnych dla osób niepełnosprawnych.
Parking jest podzielony na 3 segmenty: 
  - P5 - od strony Fun Areny, 
  - P3 i P4 - od ulicy Pokoleń Lechii Gdańsk oraz
  - P1 - parking szlabanowy VIP.
Segmenty P3, P4 i P5 są wyposażone w 12 parkometrów nowej generacji, w których można dokonać płatności nie tylko bilonem, ale również kartą.
Parking P5 jest wyposażony w 1 wjazd i 1 wyjazd oraz w 2 kasy automatyczne.

W dniach imprez masowych można płacić aplikacją FLOW. Ceny zależne od rodzaju imprezy, np. 5/10/20 PLN za czas trwania imprezy.`,
				"en": `Above-ground car park located by the Energa stadium in Gdańsk.
There are about 2500 parking spaces for clients, including 40 free spaces for disabled persons. 
The car park is divided into 3 sections: 
  - P5 - from the Fun Arena side, 
  - P3 and P4 - from Pokoleń Lechii Gdańsk street side 
  - P1 - VIP barrier car park. 
Sections P3, P4 and P5 are equipped with 12 new-generation parking meters, which accept payments not only in coins, but also by card. 
P5 section has 1 entry and 1 exit as well as 2 automatic cash registers.

On the days of mass events, you can pay with the FLOW application. Prices depend on the type of event, e.g. 5/10/20 PLN for the duration of the event.`,
			},
			Rules: []wheretopark.Rule{
				{
					Hours: "24/7",
					Pricing: []wheretopark.PricingRule{
						wheretopark.NewIntPricingRule("PT30M", 1),
						wheretopark.NewIntPricingRule("PT1H", 4).Repeated(),
					},
				},
				{
					Hours: "24/7",
					Applies: []wheretopark.SpotType{
						wheretopark.SpotTypeCarDisabled,
					},
					Pricing: []wheretopark.PricingRule{
						wheretopark.NewIntPricingRule("PT1H", 0).Repeated(),
					},
				},
			},
		},

		"4": {
			Resources: []string{
				"mailto:parkingi@apcoa.pl",
				"tel:+48-22-354-83-80",
				"https://www.apcoa.pl/parking-w/gdansk-1/amber-expo-gdansk-ul-zaglowa-11/",
			},
			TotalSpots: map[string]uint{
				wheretopark.SpotTypeCar: 510,
			},
			MaxDimensions: &wheretopark.Dimensions{
				Height: -1,
			},
			Features: []string{
				wheretopark.FeatureUncovered,
			},
			PaymentMethods: []string{
				wheretopark.PaymentMethodCash,
				wheretopark.PaymentMethodCard,
				wheretopark.PaymentMethodContactless,
			},
			Rules: []wheretopark.Rule{
				{
					Hours: "24/7",
					Pricing: []wheretopark.PricingRule{
						wheretopark.NewIntPricingRule("PT30M", 1),
						wheretopark.NewIntPricingRule("PT1H", 4).Repeated(),
					},
				},
			},
		},

		"5": {
			Resources: []string{
				"mailto:info@madison.gda.pl",
				"tel:+48-58-766-75-39",
				"https://madison.gda.pl/pl/parking",
				"https://madison.gda.pl/pl/parking/regulamin-parkingu",
			},
			TotalSpots: map[string]uint{
				wheretopark.SpotTypeCar:         111,
				wheretopark.SpotTypeCarDisabled: 4,
				wheretopark.SpotTypeMotorcycle:  2,
			},
			Features: []string{
				wheretopark.FeatureUnderground,
			},
			Comment: map[string]string{
				"pl": `DODATKOWA GODZINA PARKOWANIA GRATIS - za zakupy na kwotę min. 50 zł (szczegóły w pkt. III ust. 4 [Regulaminu Parkingu](https://madison.gda.pl/pl/parking/regulamin-parkingu)).
W niedziele i świeta wyjście z parkingu przez Centrum Medyczne Rajska`,
			},
			Rules: []wheretopark.Rule{
				{
					Hours: "Mo-Fr 08:00-22:00",
					Pricing: []wheretopark.PricingRule{
						wheretopark.NewIntPricingRule("PT1H", 0),
						wheretopark.NewIntPricingRule("PT2H", 5),
						wheretopark.NewIntPricingRule("PT3H", 10),
						wheretopark.NewIntPricingRule("PT1H", 8).Repeated(),
					},
				},
				{
					Hours: "Sa 08:00-22:00",
					Pricing: []wheretopark.PricingRule{
						wheretopark.NewIntPricingRule("PT1H", 0),
						wheretopark.NewIntPricingRule("PT2H", 0),
						wheretopark.NewIntPricingRule("PT3H", 5),
						wheretopark.NewIntPricingRule("PT1H", 8).Repeated(),
					},
				},
				{
					Hours: "Su 08:00-22:00",
					Pricing: []wheretopark.PricingRule{
						wheretopark.NewIntPricingRule("PT1H", 5).Repeated(),
					},
				},

				{
					Hours: "Mo-Su 22:00-08:00",
					Pricing: []wheretopark.PricingRule{
						wheretopark.NewIntPricingRule("PT1H", 1).Repeated(),
					},
				},
			},
		},

		"7": {
			Resources: []string{
				"mailto:e.gazda@amuz.gda.pl",
				"tel:+48-58-300-92-06",
				"https://www.amuz.gda.pl/akademia/infrastruktura/parking,108",
			},
			TotalSpots: map[string]uint{
				wheretopark.SpotTypeCar: 215,
			},
			Features: []string{
				wheretopark.FeatureUncovered,
			},
			Comment: map[string]string{
				"pl": `Goście nocujący w Domu Muzyka oraz goście hotelowi nocujący w Domu Sonata, klienci Restauracji Domu Muzyka - bezpłatnie
autokar z wycieczką nocującą w Domu Muzyka lub w Domu Sonata lub nocujący tylko kierowca autokaru - opłata indywidualna ustalana przez kierownika Domu Muzyka/Domu Sonata
					
Dla studentów i pracowników aMuz są dostępne specjalne abonamenty, więcej informacji na https://www.amuz.gda.pl/akademia/infrastruktura/parking,108
				`,
			},
			Rules: []wheretopark.Rule{
				{
					Hours: "24/7",
					Pricing: []wheretopark.PricingRule{
						wheretopark.NewIntPricingRule("PT30M", 0),
						wheretopark.NewIntPricingRule("PT1H", 6).Repeated(),
					},
				},
			},
		},

		"14": {
			Resources: []string{
				"https://ecs.gda.pl/title,Parking_,pid,601.htm",
				"tel:+48-58-772-40-00",
				"mailto:ecs@ecs.gda.pl",
			},
			TotalSpots: map[string]uint{
				wheretopark.SpotTypeCar: 287,
			},
			Features: []string{
				wheretopark.FeatureUnderground,
			},
			Comment: map[string]string{
				"pl": `Zakup biletu wstępu na wystawę stałą: 4 godziny parkowania w cenie biletu / zabierz ze sobą bilet parkingowy, a kupując bilet na wystawę okaż bilet parkingowy w kasie (zrób to nim upłynie 30 minut od chwili zaparkowania)`,
				"en": `ATTENTION! Ticket to the permanent exhibition | 4 hours parking included in ticket price. Take your parking ticket with you, and when buying your ticket for the permanent exhibition, present your parking ticket at the ticket office (do it before the 30 minutes have passed from the moment you parked)`,
			},
			Rules: []wheretopark.Rule{
				{
					Hours: "Mo-Su 06:00-22:00",
					Pricing: []wheretopark.PricingRule{
						wheretopark.NewIntPricingRule("PT30M", 0),
						wheretopark.NewIntPricingRule("PT1H", 4),
						wheretopark.NewIntPricingRule("PT2H", 8),
						wheretopark.NewIntPricingRule("PT3H", 9),
						wheretopark.NewIntPricingRule("PT1H", 10).Repeated(),
					},
				},
			},
		},

		"15": {
			Resources: []string{
				"mailto:bartlomiej.jablonski@carpark.com.pl",
				"tel:+48-661-552-882",
				"https://galeriametropolia.pl/parking/",
			},
			TotalSpots: map[string]uint{
				wheretopark.SpotTypeCar: 1500,
			},
			Features: []string{
				wheretopark.FeatureUncovered,
				wheretopark.FeatureCovered,
			},
			Comment: map[string]string{
				"pl": `Klienci kina Helios - pierwsze 3h bezpłatnie (bilet parkingowy należy okazać w kasie kina)
Klienci siłowni Tone Zone - pierwsze 2h bezpłatnie(po zeskanowaniu biletu w recepcji). Po przekroczeniu przysługującego czasu opłata zostaje naliczona za cały okres parkowania.
Miejsca abonamentowe znajduję się na poziomie +3 (zadaszone) oraz na parkingu zewnętrznym.
Parking na poziomie 0 nie oferuje możliwości miejsc abonamentowych.
				`,
			},
			Rules: []wheretopark.Rule{
				{
					Hours: "24/7",
					Pricing: []wheretopark.PricingRule{
						wheretopark.NewIntPricingRule("PT1H", 0),
						wheretopark.NewIntPricingRule("PT2H", 3),
						wheretopark.NewIntPricingRule("PT3H", 7),
						wheretopark.NewIntPricingRule("P1D", 17),
						wheretopark.NewIntPricingRule("P1M", 150),
						wheretopark.NewIntPricingRule("PT1H", 5).Repeated(),
					},
				},
			},
		},

		"17": {
			Resources: []string{
				"mailto:parking@forumgdansk.pl",
				"tel:+48-661-551-882",
				"https://forumgdansk.pl/pl/przydatne-informacje/parking",
			},
			TotalSpots: map[string]uint{
				wheretopark.SpotTypeCar: 1008,
			},
			MaxDimensions: &wheretopark.Dimensions{
				Height: 200,
			},
			Features: []string{
				wheretopark.FeatureUncovered,
				wheretopark.FeatureCovered,
			},
			Comment: map[string]string{
				"pl": `Dla klientow
- kina Helios, 3 godziny parkowania bezpłatne.
- City Fit i Media Markt, 2 godziny parkowania bezpłatne. 

Abonament miesięczny jest dostępny, po więcej informacji: https://forumgdansk.pl/pl/przydatne-informacje/parking
					
Biuro parkingu znajduje się na poziomie +2.
Na parkingu znajdują się miejsca parkingowe dla osób niepełnosprawnych oraz rodzin z dziećmi.

W godzinach 22:00 - 8:00 wejście na parking możliwe jest wejściem „nocnym” od strony ul. 3 Maja

KASY ZNAJDUJĄ SIĘ NA KAŻDYM POZIOMIE PARKIGU PRZY WINDACH

Płatności mobilne są realizowane bezdotykowo poprzez aplikację NaviPay (pobierz na Android lub IOS).

W przypadku zgubienia biletu parkingowego kopię biletu można wydrukować w kasie parkingowej.
				`,
			},
			Rules: []wheretopark.Rule{
				{
					Hours: "24/7",
					Pricing: []wheretopark.PricingRule{
						wheretopark.NewIntPricingRule("PT30M", 0),
						wheretopark.NewIntPricingRule("PT1H", 3),
						wheretopark.NewIntPricingRule("PT2H", 8),
						wheretopark.NewIntPricingRule("PT3H", 13),
						wheretopark.NewIntPricingRule("PT1H", 4).Repeated(),
					},
				},
			},
		},

		"18": {
			Resources: []string{
				"mailto:parkingi@apcoa.pl",
				"tel:+48-22-354-83-80",
				"https://www.apcoa.pl/parking-w/gdansk-1/gdansk-ul-blekitna/",
			},
			TotalSpots: map[string]uint{
				wheretopark.SpotTypeCar: 197,
			},
			MaxDimensions: &wheretopark.Dimensions{
				Height: -1,
			},
			Features: []string{
				wheretopark.FeatureUncovered,
			},
			PaymentMethods: []string{
				wheretopark.PaymentMethodCash,
				wheretopark.PaymentMethodCard,
				wheretopark.PaymentMethodContactless,
			},
			Comment: map[string]string{
				"pl": `Po wcześniejszym zgłoszeniu obsłudze przez kierującego takiej potrzeby pierwszy kwadrans będzie bezpłatny. 
Po godzinie 17:00 do dnia następnego do godziny 9:00 za pozostawienie samochodu na parkingu nie trzeba będzie płacić.
Z konieczności wniesienia opłaty za wjazd na parkingi nadmorskie zwolnione będą osoby niepełnosprawne, posiadające  identyfikatora „N+” wydany przez Gdański Zarząd Dróg i Zieleni lub odpowiadający mu dokument wydany przez innego zarządcę drogi. 
Bezpłatny będzie również postój pojazdów elektrycznych i PHEV w rozumieniu Ustawy o elektromobilności i paliwach alternatywnych. 
W tym przypadku podstawą będzie niebieska naklejka na szybie pojazdu z symbolem EE, karta zerowej stawki opłat na pojazdy PHEV a także zielona tablica rejestracyjna lub wskazanie gniazdka do ładowania samochodu pracownikowi obsługi przy wjeździe na parking.
				`,
			},
			Rules: []wheretopark.Rule{
				{
					Hours: "Mo-Su 09:00-19:00",
					Pricing: []wheretopark.PricingRule{
						wheretopark.NewIntPricingRule("P1D", 10),
					},
				},
				{
					Hours: "Mo-Su 19:00-09:00",
					Pricing: []wheretopark.PricingRule{
						wheretopark.NewIntPricingRule("PT1H", 0).Repeated(),
					},
				},
				{
					Hours: "24/7",
					Applies: []string{
						wheretopark.SpotTypeCarElectric,
						wheretopark.SpotTypeCarDisabled,
					},
					Pricing: []wheretopark.PricingRule{
						wheretopark.NewIntPricingRule("PT1H", 0).Repeated(),
					},
				},
			},
		},

		"19": {
			Resources: []string{
				"https://www.apcoa.pl/parking-w/gdansk-1/gdansk-ul-kapliczna/",
				"https://gzdiz.gda.pl/aktualnosci/parkingi-nadmorskie-w-sezonie-letnim-2021,a,4773",
			},
			TotalSpots: map[string]uint{
				wheretopark.SpotTypeCar: 60,
			},
			MaxDimensions: &wheretopark.Dimensions{
				Height: -1,
			},
			Features: []string{
				wheretopark.FeatureUncovered,
			},
			PaymentMethods: []string{
				wheretopark.PaymentMethodCash,
				wheretopark.PaymentMethodCard,
				wheretopark.PaymentMethodContactless,
			},
			Comment: map[string]string{
				"pl": `Po wcześniejszym zgłoszeniu obsłudze przez kierującego takiej potrzeby pierwszy kwadrans będzie bezpłatny. 
Po godzinie 17:00 do dnia następnego do godziny 9:00 za pozostawienie samochodu na parkingu nie trzeba będzie płacić.
Z konieczności wniesienia opłaty za wjazd na parkingi nadmorskie zwolnione będą osoby niepełnosprawne, posiadające  identyfikatora „N+” wydany przez Gdański Zarząd Dróg i Zieleni lub odpowiadający mu dokument wydany przez innego zarządcę drogi. 
Bezpłatny będzie również postój pojazdów elektrycznych i PHEV w rozumieniu Ustawy o elektromobilności i paliwach alternatywnych. 
W tym przypadku podstawą będzie niebieska naklejka na szybie pojazdu z symbolem EE, karta zerowej stawki opłat na pojazdy PHEV a także zielona tablica rejestracyjna lub wskazanie gniazdka do ładowania samochodu pracownikowi obsługi przy wjeździe na parking.
				`,
			},
			Rules: []wheretopark.Rule{
				{
					Hours: "Mo-Su 09:00-19:00",
					Pricing: []wheretopark.PricingRule{
						wheretopark.NewIntPricingRule("P1D", 10),
					},
				},
				{
					Hours: "Mo-Su 19:00-09:00",
					Pricing: []wheretopark.PricingRule{
						wheretopark.NewIntPricingRule("PT1H", 0).Repeated(),
					},
				},
				{
					Hours: "24/7",
					Applies: []string{
						wheretopark.SpotTypeCarElectric,
						wheretopark.SpotTypeCarDisabled,
					},
					Pricing: []wheretopark.PricingRule{
						wheretopark.NewIntPricingRule("PT1H", 0).Repeated(),
					},
				},
			},
		},
		"20": {
			Resources: []string{
				"https://gzdiz.gda.pl/aktualnosci/parkingi-nadmorskie-w-sezonie-letnim-2021,a,4773",
			},
			TotalSpots: map[string]uint{
				wheretopark.SpotTypeCar: 123,
			},
			MaxDimensions: &wheretopark.Dimensions{
				Height: -1,
			},
			Features: []string{
				wheretopark.FeatureUncovered,
			},
			PaymentMethods: []string{},
			Comment:        map[string]string{},
			Rules: []wheretopark.Rule{
				{
					Hours: "24/7",
					Pricing: []wheretopark.PricingRule{
						wheretopark.NewIntPricingRule("PT1H", 0).Repeated(),
					},
				},
			},
		},
		"21": {
			Resources: []string{
				"https://gzdiz.gda.pl/aktualnosci/parkingi-nadmorskie-w-sezonie-letnim-2021,a,4773",
			},
			TotalSpots: map[string]uint{
				wheretopark.SpotTypeCar: 340,
			},
			MaxDimensions: &wheretopark.Dimensions{
				Height: -1,
			},
			Features: []string{
				wheretopark.FeatureUncovered,
			},
			PaymentMethods: []string{
				wheretopark.PaymentMethodCash,
				wheretopark.PaymentMethodCard,
				wheretopark.PaymentMethodContactless,
			},
			Comment: map[string]string{
				"pl": `Po wcześniejszym zgłoszeniu obsłudze przez kierującego takiej potrzeby pierwszy kwadrans będzie bezpłatny. 
Po godzinie 17:00 do dnia następnego do godziny 9:00 za pozostawienie samochodu na parkingu nie trzeba będzie płacić.
Z konieczności wniesienia opłaty za wjazd na parkingi nadmorskie zwolnione będą osoby niepełnosprawne, posiadające  identyfikatora „N+” wydany przez Gdański Zarząd Dróg i Zieleni lub odpowiadający mu dokument wydany przez innego zarządcę drogi. 
Bezpłatny będzie również postój pojazdów elektrycznych i PHEV w rozumieniu Ustawy o elektromobilności i paliwach alternatywnych. 
W tym przypadku podstawą będzie niebieska naklejka na szybie pojazdu z symbolem EE, karta zerowej stawki opłat na pojazdy PHEV a także zielona tablica rejestracyjna lub wskazanie gniazdka do ładowania samochodu pracownikowi obsługi przy wjeździe na parking.
				`,
			},
			Rules: []wheretopark.Rule{
				{
					Hours: "Mo-Su 09:00-19:00",
					Pricing: []wheretopark.PricingRule{
						wheretopark.NewIntPricingRule("P1D", 10),
					},
				},
				{
					Hours: "Mo-Su 19:00-09:00",
					Pricing: []wheretopark.PricingRule{
						wheretopark.NewIntPricingRule("PT1H", 0).Repeated(),
					},
				},
				{
					Hours: "24/7",
					Applies: []string{
						wheretopark.SpotTypeCarElectric,
						wheretopark.SpotTypeCarDisabled,
					},
					Pricing: []wheretopark.PricingRule{
						wheretopark.NewIntPricingRule("PT1H", 0).Repeated(),
					},
				},
			},
		},
	},
}

func init() {
	for k, v := range configuration.ParkingLots {
		configuration.ParkingLots[k] = wheretopark.Metadata{
			LastUpdated:    &defaultLastUpdated,
			Name:           v.Name,
			Address:        v.Address,
			Geometry:       v.Geometry,
			Resources:      v.Resources,
			TotalSpots:     v.TotalSpots,
			MaxDimensions:  v.MaxDimensions,
			Features:       v.Features,
			PaymentMethods: v.PaymentMethods,
			Comment:        v.Comment,
			Currency:       defaultCurrency,
			Timezone:       defaultTimezone,
			Rules:          v.Rules,
		}
	}
}

var (
	defaultLastUpdated = wheretopark.MustParseDate("2022-12-17")
	defaultTimezone    = wheretopark.MustLoadLocation("Europe/Warsaw")
	defaultCurrency    = currency.PLN
)
