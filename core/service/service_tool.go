package service

import (
	"time"
	"turboengine/core/api"
)

var (
	epoch, _ = time.Parse("2006-01-02 15:04:05", "2019-01-01 00:00:00")
)

// 生成GUID
// |64    24|23      13|12      1|
// |time(ms)|id(0~3FF) |id(0~FFF)|
func (s *service) GenGUID() uint64 {
	ts := uint64(time.Now().Sub(epoch).Nanoseconds()) / uint64(time.Millisecond) // ms
	if ts > 0x1FFFFFFFFFF {                                                      // 69years
		ts = ts & 0x1FFFFFFFFFF
	}
	if s.uuidTs == ts {
		s.uuid++
		if s.uuid > 0xFFF { // 用完了
			time.Sleep(time.Millisecond)
			ts++
			s.uuid = 0
		}
	} else {
		s.uuid = 0
	}
	s.uuidTs = ts
	return ts<<22 | uint64(s.id&uint16(api.MAX_SID))<<12 | s.uuid
}
