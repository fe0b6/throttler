package throttler

import (
	"time"
)

// Throttler - объект остановки
type Throttler struct {
	lock         bool
	limiter      chan struct{}
	stopTime     time.Duration
	stopInterval time.Duration
}

// ThrottlerInit - Инициализируем
func (t *Throttler) ThrottlerInit(stopTime, stopInterval time.Duration, limit int) {
	t.stopTime = stopTime
	t.stopInterval = stopInterval

	if limit > 0 {
		t.limiter = make(chan struct{}, limit)
	}

	go t.run()
}

// Throttle - проверяем нужно ли сделать остановку, и если нужно делаем остановку
func (t *Throttler) Throttle() *Throttler {
	if t.lock == true {
		time.Sleep(t.stopTime)
	}

	if t.limiter != nil {
		t.limiter <- struct{}{}
	}

	return t
}

// ThrottlerRelease - освобождаем ресурс
func (t *Throttler) ThrottlerRelease() *Throttler {
	if t.limiter != nil {
		<-t.limiter
	}
	return t
}

func (t *Throttler) run() {
	for {
		time.Sleep(t.stopInterval)

		t.lock = true
		time.Sleep(t.stopTime)
		t.lock = false
	}
}
