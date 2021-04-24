package encoding

import (
	"bytes"
	"errors"

	"github.com/golang/protobuf/proto"
	"github.com/klauspost/reedsolomon"
	"golang.org/x/crypto/sha3"
)
import "github.com/davecgh/go-spew/spew"

type QRDecoder struct {
	notfirstQR    bool
	firstMessage  *LibV2RayQRCode
	payload       [][]byte
	PieceNeeded   int
	PieceReceived int
}

func (e *Encoder) StartQRDecode() *QRDecoder {
	return &QRDecoder{}
}

func (qd *QRDecoder) OnNewDataScanned(data []byte) error {
	newdata := new(LibV2RayQRCode)
	err := proto.Unmarshal(data, newdata)
	spew.Dump(newdata)
	if err != nil {
		return err
	}
	//Is this first QR being progressed? if so, do init
	if !qd.notfirstQR {
		qd.memalloc(newdata)
		qd.notfirstQR = true
		qd.firstMessage = newdata
		qd.PieceNeeded = int(newdata.OriginalSegmentNum)
	} else {
		//Or Check if is same message
		if !bytes.Equal(newdata.MsgChecksum, qd.firstMessage.MsgChecksum) {
			return errors.New("inconsistent message group")
		}
	}
	//Check if payload Checksum is correct
	PayloadCheckSum := sha3.Sum256(newdata.Payload)

	if !bytes.Equal(PayloadCheckSum[:], newdata.PayloadChecksum) {
		return errors.New("Bad Checksum")
	}
	//Assign payload if of differnet Message Number
	if qd.payload[newdata.MessgaeNumber] == nil {
		qd.payload[newdata.MessgaeNumber] = newdata.Payload
		//Update Stat
		qd.PieceReceived += 1
	}
	return nil
}

func (qd *QRDecoder) memalloc(firstQR *LibV2RayQRCode) {
	//Check if it is only consist of one segment
	if firstQR.MessageSegType == LibV2RayQRCode_NoSegmention {
		//No additional alloc is needed, as there is only one message
		qd.payload = make([][]byte, 1)
	} else {
		qd.payload = make([][]byte, firstQR.GroupCount)
	}
}

func (qd *QRDecoder) IsDecodeReady() bool {
	return qd.PieceNeeded <= qd.PieceReceived
}

func (qd *QRDecoder) Finish() ([]byte, error) {
	if !qd.IsDecodeReady() {
		return nil, errors.New("Not Ready")
	}
	//if there is only one piece, just return it
	if qd.firstMessage.MessageSegType == LibV2RayQRCode_NoSegmention {
		return qd.payload[0], nil
	}
	//then we do the reconstruction
	enc, err := reedsolomon.New(int(qd.firstMessage.OriginalSegmentNum), int(qd.firstMessage.GroupCount-qd.firstMessage.OriginalSegmentNum))
	if err != nil {
		return nil, err
	}
	err = enc.Reconstruct(qd.payload)
	if err != nil {
		return nil, err
	}
	var bytebuf bytes.Buffer
	enc.Join(&bytebuf, qd.payload, int(qd.firstMessage.PayloadTotalSize))
	return bytebuf.Bytes(), nil
}
