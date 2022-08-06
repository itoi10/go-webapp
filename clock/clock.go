package clock

import (
	"time"
)

type Clocker interface {
	Now() time.Time
}

// アプリケーションで実際に使うtime.Now関数のラッパー
type RealClocker struct{}

func (r RealClocker) Now() time.Time {
	return time.Now()
}

// テスト用の固定時間を返す
type FixedClocker struct{}

func (fc FixedClocker) Now() time.Time {
	return time.Date(2022, 5, 10, 12, 34, 56, 0, time.UTC)
}
