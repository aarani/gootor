package cells

import (
	"bytes"
	"errors"
	"gooTor/commands"
)

type TorCell interface {
	Write(buffer *bytes.Buffer) error
	Command() byte
}

func Read(command byte, buffer *bytes.Buffer) (*TorCell, error) {
	readCell := func() (*TorCell, error) {
		var cell TorCell
		var readErr error
		switch command {
		case commands.Version:
			cell, readErr = ReadCellVersion(buffer)
		case commands.Certs:
			cell, readErr = ReadCellCerts(buffer)
		case commands.AuthChallenge:
			cell, readErr = ReadCellAuthChallenge(buffer)
		case commands.NetInfo:
			cell, readErr = ReadCellNetInfo(buffer)
		default:
			cell = nil
			readErr = errors.New("NIE")
		}

		return &cell, readErr
	}

	return readCell()
}
