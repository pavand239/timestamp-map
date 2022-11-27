package timestampmap

import (
	"fmt"
	"time"
)

type TimestampMap struct {
	keys        []int64
	values      map[int64]interface{}
	maxSize     int64
	currentSize int64
	period      int64
}

type TimestampMapValue struct {
	Timestamp 	int64 		`json:"timestamp"`
	Value 		interface{} `json:"value"`
}

func NewTimestampMap(lifetimeSec int64, period int64) *TimestampMap {
	size := lifetimeSec/period + 1
	return &TimestampMap{
		keys:    make([]int64, 0, size),
		values:  make(map[int64]interface{}, size),
		maxSize: size,
		period:  period,
	}
}

func (t *TimestampMap) Add(val interface{}) {
	timestamp := time.Now().Unix()
	timestampRounded := t.roundTimestamp(timestamp)

	if t.currentSize == t.maxSize {
		var keyForDelete int64
		keyForDelete, t.keys = t.keys[0], t.keys[1:]
		delete(t.values, keyForDelete)
	} else {
		t.currentSize++
	}

	t.keys = append(t.keys, timestampRounded)

	t.values[timestampRounded] = val
}

func (t *TimestampMap) roundTimestamp(timestamp int64) int64 {
	return timestamp - (timestamp % t.period)
}

func (t *TimestampMap) GetValue(timestamp int64) (*TimestampMapValue, error) {
	timestampRounded := t.roundTimestamp(timestamp)
	val, ok := t.values[t.roundTimestamp(timestampRounded)]
	if ok {
		return &TimestampMapValue{
			timestampRounded,
			val,
		}, nil
	}

	key, err := searchClosestKey(t.keys, timestampRounded)
	if err != nil {
		if err.Error() == "empty slice" {
			return nil, fmt.Errorf("no saved values")
		}
		return nil, err
	}

	return &TimestampMapValue{
		key,
		t.values[key],
	}, nil
}

// keys have to be sorted ASC
// if timestamp less than first or greater than last return first or last respectivly
// if timestamp is equally close to 2 keys, return smallest of them
func searchClosestKey(keys []int64, timestamp int64) (int64, error) {
	if len(keys) == 0 {
		return 0, fmt.Errorf("empty slice")
	}
	if len(keys) == 1 {
		return keys[0], nil
	}

	if timestamp >= keys[len(keys)-1] {
		return keys[len(keys)-1], nil
	}
	if timestamp <= keys[0] {
		return keys[0], nil
	}

	if len(keys) == 2 {
		if timestamp-keys[0] <= keys[1]-timestamp {
			return keys[0], nil
		}
		return keys[1], nil
	}

	leftEdgeValue, rightEdgeValue := keys[len(keys)/2-1], keys[len(keys)/2]
	leftEdgeDiff, rightEdgeDiff := abs(timestamp-leftEdgeValue), abs(timestamp-rightEdgeValue)

	if rightEdgeDiff < leftEdgeDiff {
		return searchClosestKey(keys[len(keys)/2:], timestamp)
	}

	return searchClosestKey(keys[:len(keys)/2], timestamp)
}

func abs(n int64) int64 {
	if n < 0 {
		return -n
	}

	return n
}
