package game

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

type StrFloat float64

func (s *StrFloat) UnmarshalJSON(data []byte) error {
	if string(data) == "-1" {
		*s = -1
		return nil
	}

	if len(data) < 2 || data[0] != '"' || data[len(data)-1] != '"' {
		return fmt.Errorf("malformed data")
	}

	v, err := strconv.ParseFloat(string(data[1:len(data)-1]), 64)
	if err != nil {
		return err
	}
	*s = StrFloat(v)
	return nil
}

type Event struct {
	Code      int             `json:"code"`
	Timestamp time.Time       `json:"timestamp"`
	Data      json.RawMessage `json:"data"`
}

type EventMove struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}
