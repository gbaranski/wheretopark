package meters

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
	wheretopark "wheretopark/go"

	"github.com/rs/zerolog/log"
	"github.com/xuri/excelize/v2"
)

type Solari struct {
	basePath string
	mapping  map[uint]wheretopark.ID
}

func NewSolari(basePath string, mapping map[uint]wheretopark.ID) DataSource {
	return Solari{
		basePath,
		mapping,
	}
}

func (s Solari) WorkingScheme() WorkingScheme {
	return WorkingScheme{
		StartHour: 10,
		EndHour:   20,
		Weekdays:  []time.Weekday{time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday, time.Saturday},
	}
}
func (s Solari) Files() ([]string, error) {
	return listFilesWithExtension(s.basePath, "xlsx")
}

func (f Solari) LoadRecords(file *os.File) (map[wheretopark.ID][]Record, error) {
	spreadsheet, err := excelize.OpenReader(file)
	if err != nil {
		return nil, fmt.Errorf("error reading spreadsheet: %w", err)
	}
	defer spreadsheet.Close()

	rows, err := spreadsheet.GetRows("Transactions")
	if err != nil {
		return nil, err
	}

	records := make(map[wheretopark.ID][]Record, len(rows))
	for i, row := range rows {
		if i == 0 {
			continue
		}
		strCode := row[1]
		dateStr := row[4]
		durationStr := row[5]
		if durationStr == "" || !strings.Contains(dateStr, "/") {
			// log.Debug().Str("row", fmt.Sprintf("%v", row)).Msg("skipping entry")
			continue
		}
		startDate := wheretopark.MustParseDateTimeWith(solariDateFormat, dateStr)
		duration := wheretopark.Must(parseSolariDuration(durationStr))

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
			StartDate: startDate,
			EndDate:   startDate.Add(duration),
		})
	}

	return records, nil
}

const solariDateFormat = "1/2/06 15:04"

func parseSolariDuration(str string) (time.Duration, error) {
	a := strings.Split(str, ":")
	h := a[0]
	m := a[1]
	s := a[2]
	if s != "00" {
		return 0, fmt.Errorf("invalid duration: %s", s)
	}

	strDuration := fmt.Sprintf("%sh%sm", h, m)
	return time.ParseDuration(strDuration)
}
