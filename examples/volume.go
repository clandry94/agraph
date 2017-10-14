package main
/*
import (
	"fmt"
	"github.com/clandry94/agraph"
	"os"
	"time"
)

func main() {
	file, err := os.OpenFile("long_sample.wav", os.O_RDWR, 066)
	if err != nil {
		fmt.Println(err)
	}

	reader, err := agraph.NewWaveReader(file)
	if err != nil {
		fmt.Println(err)
	}

	f, err := os.Create("volume_increase.wav")
	defer f.Close()
	if err != nil {
		fmt.Println(err)
	}

	writer, err := agraph.NewWaveWriter(f,
		agraph.NumChannels(int(reader.Fmt.Data.NumChannels)),
		agraph.SampleRate(int(reader.Fmt.Data.SampleRate)),
		agraph.BitsPerSample(int(reader.Fmt.Data.BitsPerSample)))


	volumeNode, _ := agraph.NewNode(agraph.NopFilter, "nop1")

	volumeNode.SetSink(make(chan []float64, agraph.SOURCE_SIZE))

	go volumeNode.Process()

	start := time.Now()

	for {
		data, err := reader.ReadSampleFloat()
		if err != nil {
			fmt.Println(err)
			break
		}

		volumeNode.Source() <- data

		//_ = <-secondNode.Sink()
		//filteredData = <-secondNode.Sink()

		//fmt.Println(filtered)
		//fmt.Printf(" %v ", data)



		writer.Write()
	}

	writer.Close()

	end := time.Now()
	fmt.Println(end.Sub(start))

}
*/