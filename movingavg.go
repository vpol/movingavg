package movingavg

import (
	"encoding/json"
)

const (
	MinCounter = 10
	MinAge     = 30.0
)

type MovingAverage struct {
	Exponential uint8   `json:"e,omitempty"`
	Decay       float64 `json:"d"`
	Value       float64 `json:"v"`
	Counter     uint8   `json:"c"`
	add         AddFunc
	set         SetFunc
	get         GetFunc
}

func defaultMovingAverage() *MovingAverage {
	return &MovingAverage{
		Decay:   2 / (float64(MinAge) + 1),
		Value:   0,
		Counter: 0,
		add:     defaultAdd,
		get:     defaultGet,
		set:     defaultSet,
	}
}

type Option func(*MovingAverage)

func Value(d float64) Option {
	return func(m *MovingAverage) {
		m.Value = d
	}
}

func Counter(d uint8) Option {
	return func(m *MovingAverage) {
		m.Counter = d
	}
}

func decay(d float64) Option {
	return func(m *MovingAverage) {
		m.Exponential = 1
		m.Decay = d
	}
}


func WithDecay(d float64) Option {
	return func(m *MovingAverage) {
		if d == 0 || d == MinAge {
			return
		}
		// if decay is set, it's exponential ma
		m.Decay = 2 / (d + 1)
		m.Exponential = 1
		m.add = exponentialAdd
		m.get = exponentialGet
		m.set = exponentialSet
	}
}

type AddFunc func(ma *MovingAverage, value float64)
type SetFunc func(ma *MovingAverage, value float64)
type GetFunc func(ma *MovingAverage) (value float64)

func defaultAdd(ma *MovingAverage, value float64) {
	if ma.Value == 0 { // this is a proxy for "uninitialized"
		ma.Value = value
	} else {
		ma.Value = (value * ma.Decay) + (ma.Value * (1 - ma.Decay))
	}
}

func defaultGet(ma *MovingAverage) float64 {
	return ma.Value
}

func defaultSet(ma *MovingAverage, value float64) {
	ma.Value = value
}

func exponentialAdd(ma *MovingAverage, value float64) {
	if ma.Counter < MinCounter {
		ma.Counter++
		ma.Value += value
	} else if ma.Counter == MinCounter {
		ma.Counter++
		ma.Value = ma.Value / float64(MinCounter)
		ma.Value = (value * ma.Decay) + (ma.Value * (1 - ma.Decay))
	} else {
		ma.Value = (value * ma.Decay) + (ma.Value * (1 - ma.Decay))
	}
}

func exponentialGet(ma *MovingAverage) float64 {
	if ma.Counter <= MinCounter {
		return 0.0
	}

	return ma.Value
}

func exponentialSet(ma *MovingAverage, value float64) {
	ma.Value = value
	if ma.Counter <= MinCounter {
		ma.Counter = MinCounter + 1
	}
}

func weightedAdd(ma *MovingAverage, value float64) {
	if ma.Counter < MinCounter {
		ma.Counter++
		ma.Value += value
	} else if ma.Counter == MinCounter {
		ma.Counter++
		ma.Value = ma.Value / float64(MinCounter)
		ma.Value = (value * ma.Decay) + (ma.Value * (1 - ma.Decay))
	} else {
		ma.Value = (value * ma.Decay) + (ma.Value * (1 - ma.Decay))
	}
}

func weightedGet(ma *MovingAverage) float64 {
	if ma.Counter <= MinCounter {
		return 0.0
	}

	return ma.Value
}

func weightedSet(ma *MovingAverage, value float64) {
	ma.Value = value
	if ma.Counter <= MinCounter {
		ma.Counter = MinCounter + 1
	}
}

func (ma *MovingAverage) Add(value float64) {
	ma.add(ma, value)
}

func (ma *MovingAverage) Get() float64 {
	return ma.get(ma)
}

func (ma *MovingAverage) Set(value float64) {
	ma.set(ma, value)
}

func NewMovingAverage(options ...Option) (ma *MovingAverage) {
	ma = defaultMovingAverage()
	for _, opt := range options {
		opt(ma)
	}

	return
}

func (ma *MovingAverage) UnmarshalJSON(data []byte) (err error) {
	var s = struct {
		Exponential uint8   `json:"e,omitempty"`
		Decay       float64 `json:"d"`
		Value       float64 `json:"v"`
		Counter     uint8   `json:"c"`
	}{}
	err = json.Unmarshal(data, &s)
	if err != nil {
		return
	}

	switch s.Exponential {
	case 1:
		*ma = *NewMovingAverage(Value(s.Value), Counter(s.Counter), decay(s.Decay))
	default:
		*ma = *NewMovingAverage(Value(s.Value), Counter(s.Counter))
	}

	return
}
