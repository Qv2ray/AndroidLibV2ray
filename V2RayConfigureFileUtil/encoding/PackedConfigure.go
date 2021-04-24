package encoding

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"strings"

	"github.com/golang/protobuf/proto"
	"golang.org/x/crypto/sha3"
)
import "compress/gzip"

func (e *Encoder) PackV2RayConfigureIntoPackedForm(FileName string) ([]byte, error) {

	//Read the file
	f, err := os.Open(FileName)
	if err != nil {
		return nil, err
	}
	ft, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return e.PackV2RayConfigureIntoPackedFormB(FileName, ft)
}

func (e *Encoder) PackV2RayConfigureIntoPackedFormB(FileName string, ft []byte) ([]byte, error) {
	PbRepre := new(LibV2RayPackedConfig)
	//Guess File Type
	var err error
	PbRepre.ConfigType, err = GuessConfigType(FileName)
	if err != nil {
		return nil, err
	}
	//Calc CheckSum
	result := sha3.Sum256(ft)
	PbRepre.CheckSum = result[:]

	PbRepre.GzipCompressed = true

	var bytebuf bytes.Buffer
	r, err := gzip.NewWriterLevel(&bytebuf, gzip.BestCompression)
	if err != nil {
		return nil, err
	}
	_, err = r.Write(ft)
	if err != nil {
		return nil, err
	}
	err = r.Flush()
	if err != nil {
		return nil, err
	}
	err = r.Close()
	if err != nil {
		return nil, err
	}

	PbRepre.Payload = bytebuf.Bytes()

	out, err := proto.Marshal(PbRepre)

	if err != nil {
		return nil, err
	}
	return out, nil
}

func (e *Encoder) UnpackV2RayConfB(payload []byte) (string, []byte, error) {
	PbRepre := new(LibV2RayPackedConfig)
	err := proto.Unmarshal(payload, PbRepre)
	if err != nil {
		return "", nil, err
	}
	var plctx []byte
	//Decompress first
	if PbRepre.GzipCompressed {
		r := bytes.NewReader(PbRepre.Payload)
		gr, err := gzip.NewReader(r)
		if err != nil {
			return "", nil, err
		}
		result, err := ioutil.ReadAll(gr)
		if err != nil {
			return "", nil, err
		}
		plctx = result
	} else {
		plctx = PbRepre.Payload
	}
	psum := sha3.Sum256(plctx)
	if !bytes.Equal(PbRepre.CheckSum, psum[:]) {
		return "", nil, errors.New("Checksum Mismatch")
	}
	var extension string
	switch PbRepre.ConfigType {
	case LibV2RayPackedConfig_LibV2RaySimpleProtoV1:
		extension = ".LibV2RaySimpleProtoV1.pb"
	case LibV2RayPackedConfig_FullProto:
		extension = ".pb"
	case LibV2RayPackedConfig_FullJsonFile:
		extension = ".json"
	}
	return extension, plctx, nil
}

func GuessConfigType(FileName string) (LibV2RayPackedConfig_LibV2RayConfigureType, error) {
	if strings.HasSuffix(FileName, ".LibV2RaySimpleProtoV1.pb") {
		return LibV2RayPackedConfig_LibV2RaySimpleProtoV1, nil
	} else if strings.HasSuffix(FileName, ".pb") {
		return LibV2RayPackedConfig_FullProto, nil
	} else if strings.HasSuffix(FileName, ".json") {
		return LibV2RayPackedConfig_FullJsonFile, nil
	} else {
		return 0, errors.New("Cannot Guess Config Type")
	}
}
