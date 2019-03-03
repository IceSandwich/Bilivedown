package Net

import "testing"

func TestGetDomain(t *testing.T) {
	tf := func() func(exp string, ans string) {
		i := 0
		return func(exp string, ans string) {
			b := GetDomain(exp)
			if b != ans {
				t.Errorf("Fail at point %d, got %s, expected %s\n", i, b, ans)
			}
			i++
		}
	}()

	tf("https://www.google.com", "www.google.com")
	tf("https://api.live.bilibili.com/room/v1/Room/playUrl?cid=0&platform=h5&otype=json&quality=0", "api.live.bilibili.com")
}

func TestDownload(t *testing.T) {
	Download("https://www.google.com", "page.html")
}
