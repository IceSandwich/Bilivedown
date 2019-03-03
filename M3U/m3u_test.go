package M3U

import (
	"io/ioutil"
	"testing"
)

func TestParseM3U(t *testing.T) {
	fn := "live_2027557_332_c521e483.m3u8"
	bs, err := ioutil.ReadFile(fn)
	if err != nil {
		t.Error(err)
	}
	st := string(bs)
	data, err := ParseM3u(st)
	if err != nil {
		t.Error(err)
	}
	sf := func(description string, exp interface{}, ans interface{}) {
		if exp != ans {
			t.Errorf("%s: got %v, expected %v\n", description, exp, ans)
		}
	}
	sf("Version", data.Version, 3)
	sf("AllowCache", data.AllowCache, false)
	sf("MediaSequence", data.MediaSequence, 1549602840)
	sf("TargetDuration", data.TargetDuration, 5)
	sf("U1", data.Inf[0].File, "live_2027557_332_c521e483-1549602840.ts")
	sf("U-1", data.Inf[len(data.Inf)-1].File, "live_2027557_332_c521e483-1549602844.ts")
	sf("len", len(data.Inf), 5)
}
