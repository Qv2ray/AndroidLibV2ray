package encoding

import (
	"errors"
	"math"

	"github.com/golang/protobuf/proto"
	"github.com/klauspost/reedsolomon"

	"golang.org/x/crypto/sha3"
)

func (e *Encoder) EncodeToL2QR(payload []byte, conf QRGenConf) (*L2QRGen, error) {
	//First Check the size of data being encoded
	size := len(payload)
	Gener := new(L2QRGen)
	Gener.templ.PayloadTotalSize = int32(size)
	if size < int(conf.MaxQrSize) && !conf.ForceReconstruct {
		Gener.TotalSegment = 1
		Gener.templ.OriginalSegmentNum = 1
		Gener.templ.GroupCount = 1
		Gener.templ.MessageSegType = LibV2RayQRCode_NoSegmention
		Gener.payload = make([][]byte, 1)
		Gener.payload[0] = payload
	} else {
		Gener.PayloadSegment = (size / int(conf.MaxQrSize)) + 1
		Gener.ReconstructSegment = int(math.Max((float64)(conf.ReconsConf.AtLeastReplacement), conf.ReconsConf.AtLeastMutiply*((float64)(Gener.PayloadSegment))))
		Gener.TotalSegment = Gener.ReconstructSegment + Gener.PayloadSegment
		Gener.templ.GroupCount = int32(Gener.TotalSegment)
		Gener.templ.MessageSegType = LibV2RayQRCode_Reconstruct
		//Gener.payload = make([][]byte, Gener.TotalSegment)
		//BytePerSegment := int(math.Floor(float64(size)/float64(Gener.PayloadSegment)) + 1)
		//Gener.templ.ReconstructPadTo = int32(BytePerSegment)
		Gener.templ.OriginalSegmentNum = int32(Gener.PayloadSegment)
		/*
			//Copy First Few Segments

			for CurrentLinking := 0; CurrentLinking < Gener.PayloadSegment-1; CurrentLinking++ {
				if CurrentLinking == Gener.PayloadSegment-1 {
					//First find out if input data is even to payload
					RecSize := Gener.PayloadSegment * BytePerSegment
					padding := RecSize - size
					lseg := BytePerSegment - padding
					if padding == 0 {
						Gener.payload[CurrentLinking] = payload[CurrentLinking*BytePerSegment : (CurrentLinking+1)*BytePerSegment-1]
					} else {
						Gener.payload[CurrentLinking] = make([]byte, BytePerSegment)
						copy(Gener.payload[CurrentLinking][0:lseg], payload[CurrentLinking*BytePerSegment:(CurrentLinking+1)*BytePerSegment-1-padding])
					}
				} else {
					Gener.payload[CurrentLinking] = payload[CurrentLinking*BytePerSegment : (CurrentLinking+1)*BytePerSegment-1]
				}
			}

			//alloc memory for shards

			for CurrentAllocShard := 0; CurrentAllocShard < Gener.ReconstructSegment; CurrentAllocShard++ {
				Gener.payload[Gener.PayloadSegment+CurrentAllocShard] = make([]byte, BytePerSegment)
			}
		*/
		//now calc reconstruction
		enc, err := reedsolomon.New(Gener.PayloadSegment, Gener.ReconstructSegment)
		if err != nil {
			return nil, err
		}
		Gener.payload, err = enc.Split(payload)
		if err != nil {
			return nil, err
		}
		enc.Encode(Gener.payload)
		if err != nil {
			return nil, err
		}
		ok, err := enc.Verify(Gener.payload)
		if err != nil {
			return nil, err
		}
		if !ok {
			return nil, errors.New("RS Construct Failed")
		}

	}

	MsgCheckSum := sha3.Sum256(payload)

	Gener.templ.MsgChecksum = MsgCheckSum[:]

	return Gener, nil
}

type L2QRGen struct {
	PayloadSegment     int
	ReconstructSegment int
	TotalSegment       int
	templ              LibV2RayQRCode
	payload            [][]byte
}

func (qg *L2QRGen) GetSegmentAsByte(seg int) ([]byte, error) {
	objectPayload := qg.payload[seg]
	PayloadCheckSum := sha3.Sum256(objectPayload)
	qg.templ.Payload = objectPayload
	qg.templ.PayloadChecksum = PayloadCheckSum[:]
	qg.templ.MessgaeNumber = int32(seg)
	return proto.Marshal(&qg.templ)
}
