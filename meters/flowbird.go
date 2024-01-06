package meters

import (
	"encoding/csv"
	"fmt"
	"os"
	wheretopark "wheretopark/go"

	"github.com/rs/zerolog/log"
	"golang.org/x/text/encoding/charmap"
)

type Flowbird struct {
	basePath string
	mapping  map[Code]wheretopark.ID
}

func NewFlowbird(basePath string, mapping map[Code]wheretopark.ID) DataSource {
	return Flowbird{
		basePath,
		mapping,
	}
}

func (f Flowbird) Files() ([]string, error) {
	return listFilesWithExtension(f.basePath, "csv")
}

func (f Flowbird) LoadRecords(file *os.File) (map[wheretopark.ID][]Record, error) {
	reader := csv.NewReader(charmap.ISO8859_2.NewDecoder().Reader(file))
	reader.FieldsPerRecord = 13
	reader.Comma = ';'

	data, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}
	records := make(map[wheretopark.ID][]Record, len(data))
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

		id := f.mapping[code]
		if _, ok := records[id]; !ok {
			records[id] = make([]Record, 0, 8)
		}
		records[id] = append(records[id], Record{
			StartDate: wheretopark.MustParseDateTimeWith(flowbirdDateFormat, startDate),
			EndDate:   wheretopark.MustParseDateTimeWith(flowbirdDateFormat, endDate),
		})
	}
	return records, nil
}

const flowbirdDateFormat = "02/01/2006 15:04"
