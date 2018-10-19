package proto

import (
	"encoding/binary"
	"bytes"
)

type message struct {
	len int
	buf byte[]
}

func IntToBytes(n int) byte[]{
	
}

func Pack(msg byte[]) []byte{
	len := len(msg)
	// len转byte需要考虑大小端
	return append(len, msg...)
}

func UnPack() {
}
