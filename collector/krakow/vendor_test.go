package krakow_test

import (
	"testing"
	"wheretopark/collector/krakow"
)

const EXAMPLE_PLACEMARKS = `
<?xml version="1.0" encoding="utf-8"?>
<folder>
<placemark>
        <name>Sektor A2</name>
        <card>zbliżeniowa</card>
        <model>Solari SPAZIO EVO</model>
        <parkingmeter>2003</parkingmeter>
        <address>ul. Zwierzyniecka 21-23</address>
        <coordinates>
            <latitude>50.05756161111111</latitude>
            <longitude>19.930538805555557</longitude>
        </coordinates>
</placemark>
<placemark>
        <name>Sektor A2</name>
        <card>zbliżeniowa</card>
        <model>Solari SPAZIO EVO</model>
        <parkingmeter>2004</parkingmeter>
        <address>ul. Zwierzyniecka 27</address>
        <coordinates>
            <latitude>50.056936111111106</latitude>
            <longitude>19.9297</longitude>
        </coordinates>
</placemark>
</folder>
`

func TestLoadPlacemarks(t *testing.T) {
	_, err := krakow.LoadPlacemarks([]byte(EXAMPLE_PLACEMARKS))
	if err != nil {
		t.Fatal(err)
	}
}

func TestMetadata(t *testing.T) {
	placemarks, err := krakow.LoadPlacemarks([]byte(EXAMPLE_PLACEMARKS))
	if err != nil {
		t.Fatal(err)
	}
	for _, placemark := range placemarks {
		metadata := placemark.Metadata(20)
		err := metadata.Validate()
		if err != nil {
			t.Fatal(err)
		}
	}
}
