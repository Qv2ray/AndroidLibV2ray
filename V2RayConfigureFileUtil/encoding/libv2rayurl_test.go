package encoding_test

import (
	"bytes"
	"crypto/rand"
	"testing"

	vlencoding "github.com/xiaokangwang/V2RayConfigureFileUtil/encoding"
)

func TestLibV2RayURL(t *testing.T) {
	payloadsize := 16
	randpayload := make([]byte, payloadsize)
	rand.Read(randpayload)
	ec := &vlencoding.Encoder{}
	vstr := ec.ByteToV2RayURL(randpayload)
	dec := ec.V2RayURLToByte(vstr)
	if !bytes.Equal(randpayload, dec) {
		t.FailNow()
	}
}
