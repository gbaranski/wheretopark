package wheretopark_test

import (
	"encoding/json"
	"testing"
	"time"
	wheretopark "wheretopark/go"
	"wheretopark/go/tester"

	"github.com/google/go-cmp/cmp"
	"golang.org/x/text/currency"
)

var ()

func TestEncodeDecodeParkingLot(t *testing.T) {
	data, err := json.Marshal(tester.SampleParkingLot)
	if err != nil {
		t.Fatal(err)
	}
	var decoded wheretopark.ParkingLot
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(tester.SampleParkingLot, decoded, cmp.AllowUnexported(currency.Unit{}, time.Location{})); diff != "" {
		t.Errorf("parking lot mismatch (-want +got):\n%s", diff)
	}
}
