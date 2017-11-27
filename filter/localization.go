package filter

import (
	"time"
	"math"
	"strconv"
)

/*
	Sound Localization
	Assuming a simple, spherical head and that the sound sources are
	infinitely far away s.t. the sound reaches the ears
    in a straight line

	∆t = r/c(θ + sin(θ))

	where r is the radius of the head and c = 345 m/s (speed of sound at
	23C)

	Suppose 'o' are the eyes of the subject and Y is the pinnae, then

		  		  0°
				o  o
			 x        x
			x		   x
	   -90° Y          Y  90°
			x          x
			 x        x
				x  x
		      180°/-180°
*/

const (
	R = 9   // radius of the head in cm
	C = 345 // speed of sound in m/s
)

type Localization struct {
	source chan []uint16
	sink   chan []uint16
	Name   string
	meta   MetaData
	Angle  float64
	itd    time.Duration
	isd    int
	buffer chan []uint16
}

func newLocalization(name string, meta MetaData, angle float64) (Node, error) {
	itd, err := itd(angle, R, C)
	if err != nil {
		return nil, err
	}

	// interaural sample delay = time (ms) * sampleRate
	// NEED SAMPLE RATE!!
	//isd := (itd.Seconds() / 1000) *

	return &Localization{
		source: make(chan []uint16, SOURCE_SIZE),
		sink:   nil,
		Name:   name,
		meta:   meta,
		Angle: angle,
		itd: 	itd,
		isd: 	0,
	}, nil
}

func (n *Localization) SetSink(c chan []uint16) {
	n.sink = c
}

func (n *Localization) SetSource(c chan []uint16) {
	n.source = c
}

func (n *Localization) Source() chan []uint16 {
	return n.source
}

func (n *Localization) Sink() chan []uint16 {
	return n.sink
}

func (n *Localization) Process() error {
	for {
		select {
		case data := <-n.source:
			var filteredData, err = n.do(data)
			//fmt.Printf("Data processed from %v, here it is: %v\n", n.Name, filteredData)
			if err != nil {
				panic("Could not filter!")
			}
			n.sink <- filteredData
		}
	}
	return nil
}

func (n *Localization) do(data []uint16) ([]uint16, error) {
	return data, nil
}

// interaural time difference (ITD)
// ∆t = r/c(θ + sin(θ))
func itd(angle float64, radius float64, c float64) (time.Duration, error) {
	// calculate the itd
	itd := (radius / c) * (angle + math.Sin(angle))

	// prepare the itd to be parsed as a duration
	itdStrConv := strconv.FormatFloat(itd, 'f', -1, 64)

	// duration
	return time.ParseDuration(itdStrConv)
}

