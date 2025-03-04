package time

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

// Duration is a duration that can be Marshal and Unmarshal
type Duration struct {
	time.Duration
}

// UnmarshalJSON unmarshals a JSON duration from format "1h2m3s"
func (d *Duration) UnmarshalJSON(b []byte) (err error) {
	if b[0] == '"' {
		sd := string(b[1 : len(b)-1])
		d.Duration, err = time.ParseDuration(sd)
		return
	}

	var id int64
	id, err = strconv.ParseInt(string(b), 10, 64)
	if err != nil {
		return err
	}
	d.Duration = time.Duration(id)

	return
}

// MarshalJSON marshals a JSON duration to format "1h2m3s"
func (d Duration) MarshalJSON() (b []byte, err error) {
	return []byte(fmt.Sprintf(`%q`, d.Duration.String())), nil
}
