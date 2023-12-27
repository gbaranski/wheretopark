package wheretopark_test

import (
	"net/url"
	"testing"
	wheretopark "wheretopark/go"
)

func TestServerClient(t *testing.T) {
	url, err := url.Parse("https://api.wheretopark.app")
	if err != nil {
		t.Fatal(err)
	}
	c := wheretopark.NewServerClient(url)
	providers, err := c.Providers()
	if err != nil {
		t.Fatal(err)
	}
	allParkingLots := make(map[wheretopark.ID]wheretopark.ParkingLot)
	for _, provider := range providers {
		t.Logf("fetching from %s. URL: %s", provider.Name, provider.URL.String())
		parkingLots, err := c.GetFrom(provider)
		if err != nil {
			t.Fatalf("failed to fetch from %s: %s", provider.Name, err)
		}
		t.Logf("fetched %d parking lots from %s", len(parkingLots), provider.Name)
		for id, parkingLot := range parkingLots {
			allParkingLots[id] = parkingLot
		}
	}

	t.Logf("fetched %d parking lots from %d providers", len(allParkingLots), len(providers))

	allParkingLots2, err := c.GetFromMany(providers)
	if err != nil {
		t.Fatal(err)
	}

	equalJson[map[wheretopark.ID]wheretopark.ParkingLot](t, allParkingLots, allParkingLots2, "different parking lot maps returned from GetFrom and GetFromMany")
}
