package agraph

import "os"

type Codec interface {
	decode() []byte
	encode() []byte
}

type Mp3 struct {
	file *os.File
}

func (mp3 Mp3) decode() []byte {
	return []byte("Hello, world!")
}

func (mp3 Mp3) encode() []byte {
	return []byte("Hello, world!")
}
