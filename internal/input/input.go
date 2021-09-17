package input

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/metno/mkharp/internal/harp/obs"
)

type Reader struct {
	r      *csv.Reader
	fields []string
}

func Open(in io.Reader) (*Reader, error) {
	r := csv.NewReader(in)
	header, err := r.Read()
	if err != nil {
		return nil, fmt.Errorf("unable to read header: %w", err)
	}
	if len(header) == 0 {
		return nil, errors.New("empty header")
	}
	for i, h := range header {
		header[i] = strings.TrimSpace(h)
	}
	if header[0] != "time" {
		return nil, errors.New("first element in input header must be \"time\"")
	}
	return &Reader{
		r:      r,
		fields: header,
	}, nil
}

func (r *Reader) Parameters() []string {
	return r.fields[1:]
}

func (r *Reader) Read() ([]obs.Observation, error) {
	var ret []obs.Observation
	for {
		record, err := r.r.Read()
		if err == io.EOF {
			break
		}
		if len(record) != len(r.fields) {
			return nil, errors.New("size mismatch - header/data")
		}
		for i, v := range record {
			record[i] = strings.TrimSpace(v)
		}

		obsdate, err := time.Parse(time.RFC3339, record[0])
		if err != nil {
			return nil, err
		}
		data := make(map[string]float32)
		for i := 1; i < len(record); i++ {
			strValue := record[i]
			if strValue == "" {
				continue
			}
			value, err := strconv.ParseFloat(strValue, 32)
			if err != nil {
				return nil, fmt.Errorf("unable to parse value %s as a float", strValue)
			}
			data[r.fields[i]] = float32(value)
		}

		ret = append(
			ret,
			obs.Observation{
				ValidDate: obsdate,
				Data:      data,
			},
		)
	}

	return ret, nil
}
