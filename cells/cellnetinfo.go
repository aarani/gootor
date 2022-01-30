package cells

import (
	"bytes"
	"encoding/binary"
	"gooTor/commands"
	"io"
)

type RouterAddress struct {
	Type  byte
	Value []byte
}

func (address RouterAddress) Write(buffer *bytes.Buffer) error {
	typeWriteErr := buffer.WriteByte(address.Type)
	if typeWriteErr != nil {
		return typeWriteErr
	}

	valueLengthWriteErr := buffer.WriteByte(byte(len(address.Value)))
	if valueLengthWriteErr != nil {
		return valueLengthWriteErr
	}

	_, valueWriteErr := buffer.Write(address.Value)
	if valueWriteErr != nil {
		return valueWriteErr
	}

	return nil
}

type CellNetInfo struct {
	Time         uint32
	MyAddresses  []RouterAddress
	OtherAddress RouterAddress
}

func (cell CellNetInfo) Command() byte {
	return commands.NetInfo
}

func (cell CellNetInfo) Write(buffer *bytes.Buffer) error {
	timeWriteErr := binary.Write(buffer, binary.BigEndian, cell.Time)
	if timeWriteErr != nil {
		return timeWriteErr
	}

	otherAddressWriteErr := cell.OtherAddress.Write(buffer)
	if otherAddressWriteErr != nil {
		return otherAddressWriteErr
	}

	myAddressesCountWriteErr := buffer.WriteByte(byte(len(cell.MyAddresses)))
	if myAddressesCountWriteErr != nil {
		return myAddressesCountWriteErr
	}

	for _, address := range cell.MyAddresses {
		myAddressWriteErr := address.Write(buffer)
		if myAddressWriteErr != nil {
			return myAddressWriteErr
		}
	}

	return nil
}

func readRouterAddress(buffer *bytes.Buffer) (*RouterAddress, error) {
	routerType, typeReadErr := buffer.ReadByte()
	if typeReadErr != nil {
		return nil, typeReadErr
	}

	routerValueLength, lengthReadErr := buffer.ReadByte()
	if lengthReadErr != nil {
		return nil, lengthReadErr
	}

	routerValueBytes := make([]byte, routerValueLength)
	_, bodyReadErr := io.ReadFull(buffer, routerValueBytes)
	if bodyReadErr != nil {
		return nil, bodyReadErr
	}

	return &RouterAddress{Type: routerType, Value: routerValueBytes}, nil
}

func ReadCellNetInfo(buffer *bytes.Buffer) (*CellNetInfo, error) {
	timeBytes := make([]byte, 4)
	_, lengthReadErr := io.ReadFull(buffer, timeBytes)
	if lengthReadErr != nil {
		return nil, lengthReadErr
	}
	time := binary.BigEndian.Uint32(timeBytes)

	otherAddress, otherAddressReadErr := readRouterAddress(buffer)
	if otherAddressReadErr != nil {
		return nil, otherAddressReadErr
	}

	myAddressLength, addressLengthErr := buffer.ReadByte()
	if addressLengthErr != nil {
		return nil, addressLengthErr
	}

	var myAddresses []RouterAddress

	for i := byte(0); i < myAddressLength; i++ {
		myAddress, err := readRouterAddress(buffer)
		if err != nil {
			return nil, err
		}
		myAddresses = append(myAddresses, *myAddress)
	}

	return &CellNetInfo{time, myAddresses, *otherAddress}, nil
}
