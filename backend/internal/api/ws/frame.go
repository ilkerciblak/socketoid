package ws

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"net"
)

type dataFrameBytes struct {
	fin           byte
	rsv1          byte
	rsv2          byte
	rsv3          byte
	opcode        byte
	isMasked      byte
	payloadLength uint64
	maskingKey    []byte
	payload       []byte
}

const (
	opcodeClose        byte = 0x8
	opcodePing         byte = 0x9
	opcodePong         byte = 0xA
	opcodeBinaryFrame  byte = 0x2
	opcodeUTF8Text     byte = 0x1
	opcodeContinuation byte = 0x0
)

func ReadFrame(buffRW bufio.ReadWriter) (opcode byte, payload []byte, err error) {
	dataFrame, err := readDataFrameBytes(&buffRW)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to read dataframe: %v", err)
	}

	//	if dataFrame.isMasked
	maskedPayload, mask := dataFrame.payload, dataFrame.maskingKey
	var decodedPayload []byte
	for i, payloadByte := range maskedPayload {
		decodedPayload = append(decodedPayload, payloadByte^mask[i%4])
	}

	return dataFrame.opcode, decodedPayload, nil
}

func WriteFrame(conn net.Conn, buff bufio.ReadWriter, opcode byte, payload []byte) error {
	return nil

}

func readDataFrameBytes(buff *bufio.ReadWriter) (*dataFrameBytes, error) {
	// reading first byte
	byte0, err := buff.ReadByte()
	if err != nil {
		return nil, err
	}
	opcode := byte0 & 0x0F
	fin := byte0 & 0x80
	rsv1, rsv2, rsv3 := byte0&0x40, byte0&0x20, byte0&0x10

	byte1, err := buff.ReadByte()
	if err != nil {
		return nil, err
	}
	isMasked := byte1 & 0x80

	// payload length calculations
	var payloadLengthUint uint64

	payloadLengthByte := byte1 & 0x7F

	if l := uint64(payloadLengthByte); l <= 125 {

		payloadLengthUint = l
	} else if l == 126 {
		// read 2 bytes more
		var lengthBytes []byte

		// reading next 2 bytes
		for range 2 {
			lengthByte, err := buff.ReadByte()
			if err != nil {
				return nil, err
			}
			lengthBytes = append(lengthBytes, lengthByte)
		}
		payloadLengthUint = uint64(binary.BigEndian.Uint16(lengthBytes))

	} else if l == 127 {
		var lengthBytes []byte
		for range 8 {
			lengthByte, err := buff.ReadByte()
			if err != nil {
				return nil, err
			}
			lengthBytes = append(lengthBytes, lengthByte)
		}
		payloadLengthUint = uint64(binary.BigEndian.Uint64(lengthBytes))
	}

	// masking-key
	maskingKey := make([]byte, 4)
	if isMasked == 0x80 {
		for i := range 4 {
			maskingByte, err := buff.ReadByte()
			if err != nil {
				return nil, err
			}
			maskingKey[i] = maskingByte
		}
	}

	// reading payload
	var payload []byte
	for range payloadLengthUint {
		payloadFragment, err := buff.ReadByte()
		if err != nil {
			return nil, err
		}
		payload = append(payload, payloadFragment)
	}

	return &dataFrameBytes{
		fin:           fin,
		rsv1:          rsv1,
		rsv2:          rsv2,
		rsv3:          rsv3,
		opcode:        opcode,
		isMasked:      isMasked,
		payloadLength: payloadLengthUint,
		maskingKey:    maskingKey,
		payload:       payload,
	}, nil
}
