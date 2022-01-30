package cells

import (
	"bytes"
	"encoding/binary"
	"gooTor/commands"
	"io"
)

type CellAuthChallenge struct {
	Challenge []byte
	Methods   []uint16
}

func (cell CellAuthChallenge) Command() byte {
	return commands.AuthChallenge
}

func (cell CellAuthChallenge) Write(buffer *bytes.Buffer) error {
	_, challengeWriteErr := buffer.Write(cell.Challenge)
	if challengeWriteErr != nil {
		return challengeWriteErr
	}

	methodsCountWriteErr := binary.Write(buffer, binary.BigEndian, uint16(len(cell.Methods)))
	if methodsCountWriteErr != nil {
		return methodsCountWriteErr
	}

	for _, method := range cell.Methods {
		methodWriteErr := binary.Write(buffer, binary.BigEndian, method)
		if methodWriteErr != nil {
			return methodWriteErr
		}
	}

	return nil
}

func ReadCellAuthChallenge(buffer *bytes.Buffer) (*CellAuthChallenge, error) {
	challengeBytes := make([]byte, 32)
	_, challengeReadErr := io.ReadFull(buffer, challengeBytes)
	if challengeReadErr != nil {
		return nil, challengeReadErr
	}

	methodsLengthBytes := make([]byte, 2)
	_, methodsLengthReadErr := io.ReadFull(buffer, methodsLengthBytes)
	if methodsLengthReadErr != nil {
		return nil, methodsLengthReadErr
	}
	methodsLength := binary.BigEndian.Uint16(methodsLengthBytes)

	var methods []uint16

	for i := uint16(0); i < methodsLength; i++ {
		methodBytes := make([]byte, 2)
		_, methodReadErr := io.ReadFull(buffer, methodBytes)
		if methodReadErr != nil {
			return nil, methodReadErr
		}
		method := binary.BigEndian.Uint16(methodBytes)

		methods = append(methods, method)
	}

	return &CellAuthChallenge{Challenge: challengeBytes, Methods: methods}, nil
}
