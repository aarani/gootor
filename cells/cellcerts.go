package cells

import (
	"bytes"
	"encoding/binary"
	"gooTor/commands"
	"io"
)

type Cert struct {
	Type        byte
	Certificate []byte
}
type CellCerts struct {
	Certs []Cert
}

func (cell CellCerts) Command() byte {
	return commands.Certs
}

func (cell CellCerts) Write(buffer *bytes.Buffer) error {
	buffer.WriteByte(byte(len(cell.Certs)))

	for _, cert := range cell.Certs {
		typeWriteErr := binary.Write(buffer, binary.BigEndian, cert.Type)
		if typeWriteErr != nil {
			return typeWriteErr
		}

		certLengthErr := binary.Write(buffer, binary.BigEndian, uint16(len(cert.Certificate)))
		if certLengthErr != nil {
			return certLengthErr
		}

		buffer.Write(cert.Certificate)
	}
	return nil
}

func ReadCellCerts(buffer *bytes.Buffer) (*CellCerts, error) {
	certificatesCount, countReadError := buffer.ReadByte()
	if countReadError != nil {
		return nil, countReadError
	}

	var certificates []Cert

	for i := byte(0); i < certificatesCount; i++ {
		typeByte, typeReadErr := buffer.ReadByte()
		if typeReadErr != nil {
			return nil, typeReadErr
		}

		certificateLengthBytes := make([]byte, 2)
		_, lengthReadErr := io.ReadFull(buffer, certificateLengthBytes)
		if lengthReadErr != nil {
			return nil, lengthReadErr
		}
		certificateLength := binary.BigEndian.Uint16(certificateLengthBytes)

		certificateBodyBytes := make([]byte, certificateLength)
		_, bodyReadErr := io.ReadFull(buffer, certificateBodyBytes)
		if bodyReadErr != nil {
			return nil, bodyReadErr
		}

		certificates = append(certificates, Cert{typeByte, certificateBodyBytes})
	}

	return &CellCerts{certificates}, nil
}
