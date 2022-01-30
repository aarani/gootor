package cells

import (
	"bytes"
	"encoding/binary"
	"errors"
	"gooTor/commands"
	"io"
)

type CellVersion struct {
	Versions []uint16
}

func (cell CellVersion) Command() byte {
	return commands.Version
}

func (cell CellVersion) Write(buffer *bytes.Buffer) error {
	for _, version := range cell.Versions {
		err := binary.Write(buffer, binary.BigEndian, version)
		if err != nil {
			return err
		}
	}
	return nil
}

func ReadCellVersion(buffer *bytes.Buffer) (*CellVersion, error) {
	if buffer.Len()%2 != 0 {
		return nil, errors.New("version cell length is not divisible by 2")
	}

	var versions []uint16

	for {
		numberBytes := make([]byte, 2)
		n, err := buffer.Read(numberBytes)
		if err != nil && err != io.EOF {
			return nil, err
		}
		if n == 0 {
			break
		}

		versions = append(versions, binary.BigEndian.Uint16(numberBytes))
	}

	return &CellVersion{Versions: versions}, nil
}
