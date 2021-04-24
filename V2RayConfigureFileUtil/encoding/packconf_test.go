package encoding_test

import (
	"bytes"
	"math/rand"
	"testing"
	"time"

	vlencoding "github.com/xiaokangwang/V2RayConfigureFileUtil/encoding"
)

func TestSeedRand2(t *testing.T) {
	rand.Seed(time.Now().UTC().UnixNano())
}

func TestGetEncoder2(t *testing.T) {
	_ = &vlencoding.Encoder{}
}
func TestEncodeConf(t *testing.T) {
	ec := &vlencoding.Encoder{}
	payloadsize := 512
	randpayload := make([]byte, payloadsize)
	out, err := ec.PackV2RayConfigureIntoPackedFormB("a.pb", randpayload)
	if err != nil {
		t.Fatal(err)
	}
	a, b, err := ec.UnpackV2RayConfB(out)
	if err != nil {
		t.Fatal(err)
	}
	if a != ".pb" {
		t.Fatal("ext err", a)
	}
	if !bytes.Equal(randpayload, b) {
		t.Fatal("ctx err")
	}
}

func TestEncodeEx(t *testing.T) {
	ec := &vlencoding.Encoder{}
	payloadsize := 512
	randpayload := make([]byte, payloadsize)
	out, err := ec.PackV2RayConfigureIntoPackedFormB("a.json", randpayload)
	if err != nil {
		t.Fatal(err)
	}
	a, b, err := ec.UnpackV2RayConfB(out)
	if err != nil {
		t.Fatal(err)
	}
	if a != ".json" {
		t.Fatal("ext err", a)
	}
	if !bytes.Equal(randpayload, b) {
		t.Fatal("ctx err")
	}
}

func TestEncodeEx2(t *testing.T) {
	ec := &vlencoding.Encoder{}
	payloadsize := 512
	randpayload := make([]byte, payloadsize)
	out, err := ec.PackV2RayConfigureIntoPackedFormB("a.LibV2RaySimpleProtoV1.pb", randpayload)
	if err != nil {
		t.Fatal(err)
	}
	a, b, err := ec.UnpackV2RayConfB(out)
	if err != nil {
		t.Fatal(err)
	}
	if a != ".LibV2RaySimpleProtoV1.pb" {
		t.Fatal("ext err", a)
	}
	if !bytes.Equal(randpayload, b) {
		t.Fatal("ctx err")
	}
}
