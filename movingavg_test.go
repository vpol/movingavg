package movingavg

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"encoding/json"
)

func TestDefault(t *testing.T) {
	ma := NewMovingAverage()

	for _, f := range [10]float64{31, 37, 47, 9, 31, 18, 25, 40, 6, 0} {
		ma.Add(f)
	}

	assert.Equal(t, 26.963031351466075, ma.Get())
	ma.Set(1.0)
	assert.Equal(t, 1.0, ma.Get())
}

func TestExponentialButDefault(t *testing.T) {
	ma := NewMovingAverage(WithDecay(30))

	for _, f := range [11]float64{19, 31, 14, 8, 48, 20, 41, 25, 39, 17, 17} {
		ma.Add(f)
	}

	assert.Equal(t, 22.37419947001783, ma.Get())

	ma.Set(1.0)
	assert.Equal(t, 1.0, ma.Get())
}

func TestExponentialZero(t *testing.T) {
	ma := NewMovingAverage(WithDecay(5))

	for i, f := range [11]float64{34, 28, 5, 5, 4, 38, 40, 40, 8, 7, 7} {
		ma.Add(f)

		if uint8(i) < MinCounter {
			assert.Equal(t, 0.0, ma.Get())
		}
	}
	ma = NewMovingAverage(WithDecay(5))
	ma.Set(5)
	ma.Add(1)

	assert.True(t, ma.Get() < 5)

}

func TestExponential(t *testing.T) {

	ma := NewMovingAverage(WithDecay(5))
	for _, f := range [11]float64{12, 42, 25, 24, 44, 43, 25, 18, 37, 11, 24} {
		ma.Add(f)
	}

	assert.Equal(t, 26.733333333333338, ma.Get())
}

func TestExponential1(t *testing.T) {

	ma := NewMovingAverage(WithDecay(5))
	for _, f := range [15]float64{10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10} {
		ma.Add(f)
	}

	assert.Equal(t, 10.0, ma.Get())
}

func TestExponentialCustom(t *testing.T) {
	ma := NewMovingAverage(WithDecay(5))

	testSamples := [12]float64{1, 100, 1, 100, 1, 100, 1}
	for i, f := range testSamples {
		ma.Add(f)

		if uint8(i) < MinCounter {
			assert.Equal(t, 0.0, ma.Get())
		}
	}

	assert.True(t, ma.Get() > 1.0)
}

func TestSimpleMarshaling(t *testing.T) {

	ma := NewMovingAverage()
	var ma1 MovingAverage

	for _, f := range [10]float64{31, 37, 47, 9, 31, 18, 25, 40, 6, 0} {
		ma.Add(f)
	}

	assert.Equal(t, 26.963031351466075, ma.Get())
	d, err := json.Marshal(ma)
	assert.Nil(t, err)

	assert.Nil(t, json.Unmarshal(d, &ma1))
	assert.Equal(t, 26.963031351466075, ma1.Get())

}

func TestExponentialMarshaling(t *testing.T) {

	ma := NewMovingAverage(WithDecay(5))
	var ma1 MovingAverage

	for _, f := range [11]float64{12, 42, 25, 24, 44, 43, 25, 18, 37, 11, 24} {
		ma.Add(f)
	}

	assert.Equal(t, 26.733333333333338, ma.Get())
	d, err := json.Marshal(ma)
	assert.Nil(t, err)

	ma.Add(10)
	v := ma.Get()

	assert.Nil(t, json.Unmarshal(d, &ma1))
	assert.Equal(t, 26.733333333333338, ma1.Get())

	assert.EqualValues(t, ma1.Exponential, 1)
	ma1.Add(10)
	assert.Equal(t, ma1.Get(), v)

}
