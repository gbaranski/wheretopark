package meters

import (
	"encoding/csv"
	"os"
	"strconv"
	"time"
	wheretopark "wheretopark/go"

	"github.com/rs/zerolog/log"
	"golang.org/x/text/encoding/charmap"
)

type Flowbird struct {
	basePath string
	mapping  map[uint]wheretopark.ID
}

func NewFlowbird(basePath string, mapping map[uint]wheretopark.ID) DataSource {
	return Flowbird{
		basePath,
		mapping,
	}
}

func (f Flowbird) WorkingScheme() WorkingScheme {
	return WorkingScheme{
		StartHour: 10,
		EndHour:   20,
		Weekdays:  []time.Weekday{time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday, time.Saturday},
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
		strCode := row[2]
		totalTime := row[6]
		zone := row[8]
		subzone := row[9]
		endDate := row[12]

		if totalTime == "-" || zone == "" || subzone == "" {
			// log.Debug().Str("row", fmt.Sprintf("%v", row)).Msg("skipping entry")
			continue
		}

		code, err := strconv.ParseUint(strCode, 10, 32)
		if err != nil {
			log.Warn().Err(err).Str("code", row[2]).Msg("failed to parse code")
			continue
		}
		id, ok := f.mapping[uint(code)]
		if !ok {
			// log.Debug().Uint64("code", code).Msg("missing code mapping")
			continue
		}
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
