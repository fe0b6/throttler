package throttler

import (
	"sync"
	"testing"
	"time"
)

type testObj struct {
	ID int
	Throttler
}

func TestLimit(t *testing.T) {

	test := testObj{ID: 12}

	test.ThrottlerInit(time.Millisecond*50, time.Millisecond*200, 1)

	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go runWorker(&test, &wg)
	}

	wg.Wait()

	if test.ID != 1000000012 {
		t.Fatal("ID not compare")
	}
}

func runWorker(test *testObj, wg *sync.WaitGroup) {
	defer wg.Done()
	defer test.Throttle().ThrottlerRelease()

	for i := 0; i < 100000000; i++ {
		test.ID++
	}
}

func TestNoLimit(t *testing.T) {

	test := testObj{ID: 12}

	test.ThrottlerInit(time.Millisecond*50, time.Millisecond*200, 0)

	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go runWorker(&test, &wg)
	}

	wg.Wait()
}

func TestBigWait(t *testing.T) {

	test := testObj{ID: 12}

	tn := time.Now()
	test.ThrottlerInit(time.Second*10, time.Nanosecond*50, 4)

	var wg sync.WaitGroup

	for i := 0; i < 20; i++ {
		wg.Add(1)
		go runWorker(&test, &wg)
	}

	wg.Wait()

	if time.Now().Sub(tn) < time.Second*10 {
		t.Fatal("throttle not work")
	}

}
