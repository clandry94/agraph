[![Build Status](https://travis-ci.org/clandry94/agraph.svg?branch=master)](https://travis-ci.org/clandry94/agraph) [![codecov](https://codecov.io/gh/clandry94/agraph/branch/master/graph/badge.svg)](https://codecov.io/gh/clandry94/agraph)



# AGraph
AGraph is a graph based audio processing library for Go! Inspired from implementing binaural audio virtualization
with [Libavfilter](https://libav.org/documentation/libavfilter.html#Introduction) and frustrated by the severe lack of audio processing libraries for Go,
I decided to make my own.

# Usage

Using this library is simple, all you need to do is initialize filter nodes, connect them, start each node, then
feed your wave pulse code modulated (PCM) data into the first sink. Below is an example of creating an audio file reader,
setting up a graph with a convolution filter and a volume filter, and pumping data through it.

```
reader, err := agraph.NewWaveReader(file)
if err != nil {
	fmt.Println(err)
}

convolutionNode, _ := agraph.NewNode(agraph.ConvolutionFilter, "convoluter")
volumeNode, _ := agraph.NewNode(agraph.VolumeFilter, "volume booster")

convolutionNode.SetSink(volumeNode.Source())
volumeNode.SetSink(make(chan []float64, agraph.SOURCE_SIZE)

go convolutionNode.Process()
go volumeNode.Process()

for {
    data, err := reader.ReadSampleFloat()
    if err ! nil {
        fmt.Error(err)
        break
    }

    convolutionNode.Source() <- data
    filteredData = <- volumeNode.Sink()
}
```

# Todo
- wave file writing
- convolution filter
- volume filter
- leaky integrator filter
- FFT
- some frequency domain filters
