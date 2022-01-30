package network

import (
	"bufio"
	"bytes"
	tls2 "crypto/tls"
	"encoding/binary"
	"fmt"
	"gooTor/cells"
	"gooTor/commands"
	"io"
)

type Message struct {
	circuitId uint16
	cell      *cells.TorCell
}

type GuardConn struct {
	client         *tls2.Conn
	messageChannel chan *Message
}

const network = "tcp"

func connect(endpoint string) (*GuardConn, error) {
	config := tls2.Config{
		InsecureSkipVerify: true,
	}
	tlsConn, connectErr := tls2.Dial(network, endpoint, &config)
	if connectErr != nil {
		return nil, connectErr
	}

	readChannel := make(chan *Message)
	guard := &GuardConn{tlsConn, readChannel}

	go ReadCell(guard)

	return guard, nil
}

func ReadCell(conn *GuardConn) {
	reader := bufio.NewReader(conn.client)
	for {
		circuitIdBytes := make([]byte, 2)
		_, circuitIdErr := io.ReadFull(reader, circuitIdBytes)
		if circuitIdErr != nil {
			return
		}
		circuitId := binary.BigEndian.Uint16(circuitIdBytes)

		command, commandReadErr := reader.ReadByte()
		if commandReadErr != nil {
			break
		}

		if commands.IsVariableLength(command) {
			lengthBytes := make([]byte, 2)
			_, lengthReadErr := io.ReadFull(reader, lengthBytes)
			if lengthReadErr != nil {
				break
			}
			length := binary.BigEndian.Uint16(lengthBytes)

			bodyBytes := make([]byte, length)
			bodyBuffer := bytes.NewBuffer(bodyBytes)
			_, bodyReadErr := io.ReadFull(reader, bodyBytes)
			if bodyReadErr != nil {
				break
			}

			cell, readCellErr := cells.Read(command, bodyBuffer)
			if readCellErr != nil {
				break
			}

			conn.messageChannel <- &Message{circuitId: circuitId, cell: cell}

		} else {
			bodyBytes := make([]byte, 509)
			bodyBuffer := bytes.NewBuffer(bodyBytes)
			_, bodyReadErr := io.ReadFull(reader, bodyBytes)
			if bodyReadErr != nil {
				break
			}

			cell, readCellErr := cells.Read(command, bodyBuffer)
			if readCellErr != nil {
				break
			}

			conn.messageChannel <- &Message{circuitId: circuitId, cell: cell}
		}
	}
	fmt.Println("I died")
}

func (guard *GuardConn) WriteCell(circuitId uint16, cell cells.TorCell) error {
	writeBuffer := new(bytes.Buffer)
	cellBuffer := new(bytes.Buffer)

	// Write circuitId into writeBuffer
	circuitIdWriteErr := binary.Write(writeBuffer, binary.BigEndian, circuitId)
	if circuitIdWriteErr != nil {
		return circuitIdWriteErr
	}

	// Write command into writeBuffer
	writeBuffer.WriteByte(cell.Command())

	// Write cell payload into cellBuffer
	serializeErr := cell.Write(cellBuffer)
	if serializeErr != nil {
		return serializeErr
	}

	if commands.IsVariableLength(cell.Command()) {
		// If cell has variable length write uint16 length into writeBuffer
		variableLengthWriteErr := binary.Write(writeBuffer, binary.BigEndian, uint16(cellBuffer.Len()))
		if variableLengthWriteErr != nil {
			return variableLengthWriteErr
		}
	} else {
		// If cell has fixed-size, create a padding array and write it into cellBuffer
		paddingBytes := make([]byte, 509-cellBuffer.Len())
		cellBuffer.Write(paddingBytes)
	}

	_, cellWriteErr := cellBuffer.WriteTo(writeBuffer)
	if cellWriteErr != nil {
		return cellWriteErr
	}

	_, clientWriteErr := guard.client.Write(writeBuffer.Bytes())
	if clientWriteErr != nil {
		return clientWriteErr
	}

	return nil
}

func (guard *GuardConn) handshake() error {
	versionCell :=
		cells.CellVersion{
			Versions: []uint16{3},
		}

	err := guard.WriteCell(0, versionCell)
	if err != nil {
		return err
	}
	theirVersionCell := (*((<-guard.messageChannel).cell)).(*cells.CellVersion)
	fmt.Printf("Version count: %d", len(theirVersionCell.Versions))
	theirCertsCell := (*((<-guard.messageChannel).cell)).(*cells.CellCerts)
	fmt.Printf("Certificate count: %d", len(theirCertsCell.Certs))
	theirAuthChallengeCell := (*((<-guard.messageChannel).cell)).(*cells.CellAuthChallenge)
	fmt.Printf("Methods count: %d", len(theirAuthChallengeCell.Methods))
	theirNetInfoCell := (*((<-guard.messageChannel).cell)).(*cells.CellNetInfo)
	fmt.Printf("NetInfo time: %d", theirNetInfoCell.Time)
	return nil
}
