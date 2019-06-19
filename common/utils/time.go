package utils

import (
	"time"
)

type FPS struct {
	fps        int
	frameCount int
	timeLeft   time.Duration
}

func (f *FPS) Reset() {
	f.fps = 0
	f.frameCount = 0
	f.timeLeft = time.Second
}

func (f *FPS) Update(delta time.Duration) {
	f.frameCount++
	f.timeLeft = f.timeLeft - delta
	if f.timeLeft < 0 {
		f.fps = f.frameCount
		f.timeLeft = time.Second
		f.frameCount = 0
	}
}

func (f *FPS) FixedFPS(fps int) time.Duration {
	remain := fps - f.frameCount - 1
	if remain > 0 {
		return f.timeLeft / time.Duration(remain)
	}

	ps := time.Second / time.Duration(fps)
	offset := f.timeLeft
	if offset < ps { // 防止间隔过小
		offset = ps
	}
	return offset
}

// 服务的时间集合
type Time struct {
	time       time.Time     // 启动时间
	updateTime time.Time     // 帧更新时间
	deltaTime  time.Duration // 帧间隔时间
	frameCount int           // 总帧数
	fps        *FPS          // fps
	fixedFPS   int           // fps 固定帧数
}

func NewTime(fps int) *Time {
	t := &Time{}
	t.time = time.Now()
	t.updateTime = t.time
	t.fps = &FPS{}
	t.fps.Reset()
	t.fixedFPS = fps
	return t
}

func (t *Time) FrameCount() int {
	return t.frameCount
}

func (t *Time) FPS() int {
	if t.fps.fps == 0 {
		return t.fps.frameCount
	}
	return t.fps.fps
}

func (t *Time) FixedFPS() int {
	return t.fixedFPS
}

func (t *Time) NextFrame() time.Duration {
	return t.fps.FixedFPS(t.fixedFPS)
}

func (t *Time) DeltaTime() time.Duration {
	return t.deltaTime
}

// 更新时间
func (t *Time) Update() {
	now := time.Now()
	t.deltaTime = now.Sub(t.updateTime)
	t.updateTime = now
	t.frameCount++
	t.fps.Update(t.deltaTime)
}

// 获取服务运行时间
func (t *Time) Time() time.Duration {
	return time.Now().Sub(t.time)
}
