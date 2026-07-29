package utils

import (
	"GameManager/domain/entity"
	"bytes"
	"encoding/binary"
	"io"
)

// Pack 将多个 Img 打包成一个 []byte
func Pack(imgs ...entity.Img) []byte {
	buf := new(bytes.Buffer)
	for _, img := range imgs {
		// 写入长度（使用固定字节数，例如 4 字节）
		length := uint32(len(img))
		binary.Write(buf, binary.LittleEndian, length)
		// 写入数据
		buf.Write(img)
	}
	return buf.Bytes()
}

// Unpack 将打包后的二进制数据解包为多个 Img
func Unpack(data []byte) ([]entity.Img, error) {
	reader := bytes.NewReader(data)
	var imgs []entity.Img
	for reader.Len() > 0 {
		var length uint32
		err := binary.Read(reader, binary.LittleEndian, &length)
		if err != nil {
			return nil, err
		}
		img := make(entity.Img, length)
		_, err = io.ReadFull(reader, img)
		if err != nil {
			return nil, err
		}
		imgs = append(imgs, img)
	}
	return imgs, nil
}
