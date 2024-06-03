package fv

import (
	"strings"
	"time"
)

type FVTime time.Time

func (t *FVTime) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var buf string
	err := unmarshal(&buf)
	if err != nil {
		return err
	}

	tt, err := time.Parse("2006-01-02", strings.TrimSpace(buf))
	if err != nil {
		return err
	}
	*t = FVTime(tt)
	return nil
}

func (t FVTime) Time() time.Time {
	return time.Time(t)
}
