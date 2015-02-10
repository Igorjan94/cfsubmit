package cfsubmit

import (
	"strconv"
	"time"
)

type EpochTime time.Time

func (t *EpochTime) UnmarshalJSON(b []byte) error {
	result, err := strconv.ParseInt(string(b), 0, 64)
	if err != nil {
		return err
	}
	*t = EpochTime(time.Unix(result/1000, 0))
	return nil
}
