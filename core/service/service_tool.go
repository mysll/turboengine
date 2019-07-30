package service

import "time"

var (
	magic_time, _ = time.Parse("2006-01-02 15:04:05", "2019-01-01 00:00:00")
)

// 生成GUID
// |63 48|47 16|15       4|3   0|
// |sid  |time |id(0~FFF) |ms   |
func (s *service) GenerateGUID() int64 {
	s.uuid++
	dur := time.Now().Sub(magic_time).Seconds()
	ms := int64(dur*10) - int64(dur)*10
	if ms == 0 {
		ms = 1
	}
	return (int64(s.id)&0xFFFF)<<48 |
		(int64(dur)&0xFFFFFFFF)<<16 |
		(int64(s.uuid%0xFFF)&0xFFF)<<4 |
		int64(0xF/ms)&0xF
}
