package meters

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
	"time"
	wheretopark "wheretopark/go"

	"github.com/rs/zerolog/log"
	"golang.org/x/text/encoding/charmap"
)

type Flowbird struct {
	BasePath string
}

func NewFlowbird(basePath string) DataSource {
	return Flowbird{
		BasePath: basePath,
	}
}

func (f Flowbird) Files() ([]string, error) {
	return listFilesByExtension(f.BasePath, ".csv")
}

func (f Flowbird) LoadRecords(file *os.File) (map[Code][]Record, error) {
	reader := csv.NewReader(charmap.ISO8859_2.NewDecoder().Reader(file))
	reader.FieldsPerRecord = 13
	reader.Comma = ';'

	data, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}
	records := make(map[Code][]Record, len(data))
	for i, row := range data {
		if i == 0 {
			continue
		}
		startDate := row[1]
		code := row[2]
		totalTime := row[6]
		zone := row[8]
		subzone := row[9]
		endDate := row[12]

		if totalTime == "-" || zone == "" || subzone == "" {
			log.Debug().Str("row", fmt.Sprintf("%v", row)).Msg("skipping entry")
			continue
		}

		if _, ok := records[code]; !ok {
			records[code] = make([]Record, 0, 8)
		}
		records[code] = append(records[code], Record{
			StartDate: wheretopark.MustParseDateTimeWith(flowbirdDateFormat, startDate),
			EndDate:   wheretopark.MustParseDateTimeWith(flowbirdDateFormat, endDate),
		})
	}
	return records, nil
}

const flowbirdDateFormat = "02/01/2006 15:04"

func parseFlowbirdDuration(s string) (time.Duration, error) {
	s = strings.ReplaceAll(s, " ", "")
	return time.ParseDuration(s)
}
