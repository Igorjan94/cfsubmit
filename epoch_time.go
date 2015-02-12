package cfsubmit

import (
	"strconv"
	"time"
)

//time.Time object; can be unmarshalled from JSON given an unix epoch time
//See https://groups.google.com/forum/#!topic/golang-nuts/FozkbHiSP6M
type EpochTime time.Time

func (t *EpochTime) UnmarshalJSON(b []byte) error {
	result, err := strconv.ParseInt(string(b), 0, 64)
	if err != nil {
		return err
	}
	*t = EpochTime(time.Unix(result, 0))
	return nil
}
