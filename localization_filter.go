package agraph

import (
	"time"
	"math"
	"strconv"
	"container/list"
	"fmt"
	"strings"
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
		     -180°/180°
*/

const (
	R = 9   // radius of the head in cm
	C = 345 // speed of sound in m/s
	sampleRate = 22050  // won't work with audio that isn't 8000 (won't give the proper delay length)
)

type Localization struct {
	source chan []uint16
	sink   chan []uint16
	Name   string
	Angle  float64
	itd    time.Duration
	sampleDelay int
	delayChannel int
	delayBuf *list.List
}

func newLocalization(name string, angle float64) (Node, error) {
	itd, err := itd(angle, R, C)
	if err != nil {
		return nil, err
	}

	delayChannel := 0

	if angle > 0 {
		delayChannel = 1
	}

	// number of samples to buffer
	sampleDelay := int(math.Floor(itd.Seconds() * sampleRate))

	fmt.Printf("sample delay info: \n" +
		"  - itd seconds: %v \n" +
		"  - sampleRate: %v \n" +
		"  - sample delay %v \n", itd.Seconds(), sampleRate, sampleDelay)

	return &Localization{
		source: make(chan []uint16, SOURCE_SIZE),
		sink:   nil,
		Name:   name,
		Angle: angle,
		itd: 	itd,
		sampleDelay: sampleDelay,
		delayChannel: delayChannel,
		delayBuf: list.New(),
	}, nil
}

func (n *Localization) SetSink(c chan []uint16) {
	n.sink = c
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
	// add samples to the buffer if it isn't full
	if n.delayBuf.Len() < n.sampleDelay {
		n.delayBuf.PushFront(data[n.delayChannel])
	}
	//fmt.Println(n.sampleDelay)
	//fmt.Printf("LENGH OF DELAY BUF: %v\n", n.delayBuf.Len())

	// if the buffer is full, replace the sample with the value
	if n.delayBuf.Len() == n.sampleDelay {
		samp := n.delayBuf.Back().Value.(uint16)
		//fmt.Printf("Reading from buffer! \n  - samp uint16: %v \n", samp)
		n.delayBuf.Remove(n.delayBuf.Back())
		data[n.delayChannel] = samp
	} else {
		// otherwise, set that sample in the data stream to 0
		data[n.delayChannel] = 0
	}

	return data, nil
}

// interaural time difference (ITD)
// ∆t = r/c(θ + sin(θ))
// angle is in degrees coming in
func itd(angle float64, radius float64, c float64) (time.Duration, error) {

	angRad := angle * 0.0174533

	// calculate the itd
	itd := math.Abs((radius / c) * (angRad + math.Sin(angle)))

	// prepare the itd to be parsed as a duration
	itdStrConv := strconv.FormatFloat(itd, 'f', -1, 64)
	fmt.Println(itdStrConv)
	unitIdent := "s"
	itdStrArr := []string{itdStrConv, unitIdent}


	itdStr := strings.Join(itdStrArr, "")

	// duration
	return time.ParseDuration(itdStr)
}

