package utils

import (
	"bytes"
	"encoding/gob"
	"sync"
)

type GobBuff struct{}

var bufferPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

// NewBuffer 从池中获取新 bytes.Buffer
func (gb GobBuff) NewBuffer() *bytes.Buffer {
	return bufferPool.Get().(*bytes.Buffer)
}

// PutBuffer 将 Buffer放入池中
func (gb GobBuff) PutBuffer(buf *bytes.Buffer) {
	// See https://golang.org/issue/23199
	const maxSize = 1 << 16
	if buf != nil && buf.Cap() < maxSize { // 对于大Buffer直接丢弃
		buf.Reset()
		bufferPool.Put(buf)
	}
}

// BytesToStruct
// 从字节转换成结构体
func (gb GobBuff) BytesToStruct(data []byte, obj interface{}) error {
	buf := gb.NewBuffer()
	defer gb.PutBuffer(buf)
	buf.Write(data)
	return gob.NewDecoder(buf).Decode(obj)
}

// StructToBytes
// 将一个接口转换成字节
func (gb GobBuff) StructToBytes(obj interface{}) ([]byte, error) {
	buf := gb.NewBuffer()
	defer gb.PutBuffer(buf)
	err := gob.NewEncoder(buf).Encode(obj)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
