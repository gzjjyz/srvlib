/**
 * @Author: ChenJunJi
 * @Desc:
 * @Date: 2021/9/3 18:30
 */

package pb3

import (
	"github.com/golang/protobuf/proto"
	"github.com/golang/snappy"
	"github.com/gzjjyz/logger"
)

const PrintError = true

type Message = proto.Message

func Marshal(m proto.Message) (buf []byte, err error) {
	buf, err = proto.Marshal(m)
	if nil != err {
		logger.StackIf(PrintError, "pb3 Marshal error %v", err)
	}
	return
}

func Unmarshal(b []byte, m proto.Message) (err error) {
	err = proto.Unmarshal(b, m)
	if nil != err {
		logger.StackIf(PrintError, "pb3 Unmarshal error %v", err)
	}
	return err
}

// CompressByte 压缩
func CompressByte(pb proto.Message) []byte {
	buff, err := proto.Marshal(pb)
	if err != nil {
		return []byte{}
	}
	return snappy.Encode(nil, buff)
}

// UnCompress 解压缩
func UnCompress(data []byte) []byte {
	buff, _ := snappy.Decode(nil, data)
	return buff
}
