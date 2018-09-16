package moments

import (
	"encoding/gob"
	"errors"
	"io"
	"sync"
	"time"
)

var ErrNoSavedMoments = errors.New("no saved data")

type Moment struct {
	Tm           *time.Time
	MovingWindow time.Duration
}

func (m *Moment) Expired() bool {
	if m.Tm.Before(time.Now().Add(-1 * m.MovingWindow)) {
		return true
	}
	return false
}

type MomentsCounter struct {
	items []Moment

	movingWindow time.Duration

	mu sync.Mutex
}

func NewMomentsCounter(movingWindow time.Duration) *MomentsCounter {
	return &MomentsCounter{
		items:        make([]Moment, 0),
		movingWindow: movingWindow,
	}
}

func NewMomentsCounterFrom(r io.Reader) (*MomentsCounter, error) {
	deserializeObj := &struct {
		Items        []Moment
		MovingWindow time.Duration
	}{}

	dataDecoder := gob.NewDecoder(r)
	err := dataDecoder.Decode(&deserializeObj)
	if err != nil {
		if err == io.EOF {
			return nil, ErrNoSavedMoments
		}
		return nil, err
	}

	return &MomentsCounter{
		items:        deserializeObj.Items,
		movingWindow: deserializeObj.MovingWindow,
	}, nil
}

func (moments *MomentsCounter) Track() {
	moments.mu.Lock()
	defer moments.mu.Unlock()

	now := time.Now()
	m := Moment{
		Tm:           &now,
		MovingWindow: moments.movingWindow,
	}

	moments.items = append(moments.items, m)
}

func (moments *MomentsCounter) Count() int {
	var counter int

	moments.mu.Lock()
	defer moments.mu.Unlock()

	for _, m := range moments.items {
		if m.Expired() {
			counter++
		} else {
			break
		}
	}

	moments.items = moments.items[counter:]

	return len(moments.items)
}

func (moments *MomentsCounter) Save(w io.Writer) error {
	dataEncoder := gob.NewEncoder(w)

	moments.mu.Lock()
	defer moments.mu.Unlock()

	serializeObj := &struct {
		Items        []Moment
		MovingWindow time.Duration
	}{moments.items, moments.movingWindow}

	return dataEncoder.Encode(serializeObj)
}
