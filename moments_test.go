package moments_test

import (
	"bytes"
	"encoding/gob"
	"io"
	"sync"
	"testing"
	"time"

	"github.com/valerykalashnikov/moments"
)

func TestMoments(t *testing.T) {
	counter := moments.NewMomentsCounter(100 * time.Millisecond)

	assertEqual := func(expected int) {
		val := counter.Count()
		if val != expected {
			t.Error("Expected: ", expected, ",actual ", val)
		}
	}

	var wg sync.WaitGroup

	wg.Add(2)
	go func() {
		defer wg.Done()
		counter.Track()
	}()
	time.Sleep(50 * time.Millisecond)
	go func() {
		defer wg.Done()
		counter.Track()
	}()

	wg.Wait()
	assertEqual(2)

	time.Sleep(50 * time.Millisecond)
	assertEqual(1)

	time.Sleep(50 * time.Millisecond)
	assertEqual(0)
}

func TestSaveMoments(t *testing.T) {
	var b bytes.Buffer

	movingWindow := 100 * time.Millisecond

	counter := moments.NewMomentsCounter(movingWindow)

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		counter.Track()
	}()

	go func() {
		defer wg.Done()
		err := counter.Save(&b)
		if err != nil {
			if err != io.EOF {
				t.Error("Unexpected error: ", err)
			}
		}
	}()

	wg.Wait()

	deserializeObj := &struct {
		Items        []moments.Moment
		MovingWindow time.Duration
	}{}

	savedValue := b.Bytes()
	r := bytes.NewReader(savedValue)

	dataDecoder := gob.NewDecoder(r)
	err := dataDecoder.Decode(&deserializeObj)
	if err != nil {
		t.Error("Unexpected error: ", err)
		return
	}

	if len(deserializeObj.Items) != 1 {
		t.Errorf(
			"Length of saved items doesn't match with original one: expected %v, actual: %v",
			1,
			len(deserializeObj.Items),
		)
		return
	}

	if deserializeObj.MovingWindow != movingWindow {
		t.Errorf(
			"Moving window of saved moments doesn't match with original one: expected %v, actual: %v",
			movingWindow,
			deserializeObj.MovingWindow,
		)
		return
	}
}

func TestNewMomentsCounterFrom_Success(t *testing.T) {
	var b bytes.Buffer

	movingWindow := 200 * time.Millisecond
	counter := moments.NewMomentsCounter(movingWindow)
	counter.Track()

	err := counter.Save(&b)
	if err != nil {
		if err != io.EOF {
			t.Error("Unexpected error: ", err)
			return
		}
	}

	savedValue := b.Bytes()
	r := bytes.NewReader(savedValue)

	newCounter, err := moments.NewMomentsCounterFrom(r)
	if err != nil {
		t.Error("Unexpected error: ", err)
		return
	}
	expectedCount := newCounter.Count()

	if expectedCount != 1 {
		t.Errorf(
			"Count of restored moments doesn't match with original one: expected %v, actual: %v",
			1,
			expectedCount,
		)
	}
}

func TestNewMomentsCounterFrom_ErrNoSavedData(t *testing.T) {
	var b bytes.Buffer

	savedValue := b.Bytes()
	r := bytes.NewReader(savedValue)

	_, err := moments.NewMomentsCounterFrom(r)
	if err != moments.ErrNoSavedMoments {
		t.Error("Unexpected error: ", err)
		return
	}
}

func BenchmarkMomentsAdd(b *testing.B) {
	moments := moments.NewMomentsCounter(1 * time.Minute)

	for n := 0; n < b.N; n++ {
		moments.Track()
		moments.Count()
	}
}
