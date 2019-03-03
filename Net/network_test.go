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

	tf("https://www.baidu.com", "www.baidu.com")
	tf("https://api.live.bilibili.com/room/v1/Room/playUrl?cid=%d&platform=h5&otype=json&quality=0", "api.live.bilibili.com")
}

func TestDownload(t *testing.T) {
	//Download("https://api.live.bilibili.com/room/v1/Room/playUrl?cid=37702&platform=h5&otype=json&quality=0", "page.html")
	Download("https://www.baidu.com", "page2.html")
}
