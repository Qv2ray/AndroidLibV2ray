package encoding_test

import (
	"bytes"
	"math/rand"
	"reflect"
	"testing"
	"time"

	vlencoding "github.com/xiaokangwang/V2RayConfigureFileUtil/encoding"
)

func TestSeedRand(t *testing.T) {
	rand.Seed(time.Now().UTC().UnixNano())
}

func TestGetEncoder(t *testing.T) {
	_ = &vlencoding.Encoder{}
}

func TestQRConf(t *testing.T) {
	_ = &vlencoding.QRGenConf{}
}

func TestGenerateQRSignle(t *testing.T) {
	qrconf := &vlencoding.QRGenConf{}
	encoder := &vlencoding.Encoder{}
	qrconf.ForceReconstruct = false
	qrconf.MaxQrSize = 1024
	qrconf.ReconsConf = &vlencoding.QRGenConf_ReconstructConf{AtLeastMutiply: 2.0, AtLeastReplacement: 2}
	payloadsize := 512
	randpayload := make([]byte, payloadsize)
	_, err := rand.Read(randpayload)
	if err != nil {
		t.Fatal(err)
	}
	result, err := encoder.EncodeToL2QR(randpayload, *qrconf)
	if err != nil {
		t.Fatal(err)
	}
	qrp, err := result.GetSegmentAsByte(0)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(qrp)
}

func TestDecodeQRSignle(t *testing.T) {
	payloadsize := 512
	randpayload := make([]byte, payloadsize)
	_, err := rand.Read(randpayload)
	qrcode := func() []byte {
		qrconf := &vlencoding.QRGenConf{}
		encoder := &vlencoding.Encoder{}
		qrconf.ForceReconstruct = false
		qrconf.MaxQrSize = 1024
		qrconf.ReconsConf = &vlencoding.QRGenConf_ReconstructConf{AtLeastMutiply: 2.0, AtLeastReplacement: 2}
		if err != nil {
			t.Fatal(err)
		}
		result, err := encoder.EncodeToL2QR(randpayload, *qrconf)
		if err != nil {
			t.Fatal(err)
		}
		qrp, err := result.GetSegmentAsByte(0)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(qrp)
		return qrp
	}()
	encoder := &vlencoding.Encoder{}
	qrdecoder := encoder.StartQRDecode()
	qrdecoder.OnNewDataScanned(qrcode)
	if !qrdecoder.IsDecodeReady() {
		t.Fatal("Asked more segment than necessary")
	}
	output, err := qrdecoder.Finish()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(output)
	if !bytes.Equal(output, randpayload) {
		t.Fatal("Payload Mismatch")
	}
}

func TestDecodeQRMany(t *testing.T) {
	payloadsize := 512
	randpayload := make([]byte, payloadsize)
	_, err := rand.Read(randpayload)
	qrcode := func() [][]byte {
		qrconf := &vlencoding.QRGenConf{}
		encoder := &vlencoding.Encoder{}
		qrconf.ForceReconstruct = true
		qrconf.MaxQrSize = 1024
		qrconf.ReconsConf = &vlencoding.QRGenConf_ReconstructConf{AtLeastMutiply: 2.0, AtLeastReplacement: 2}
		if err != nil {
			t.Fatal(err)
		}
		result, err := encoder.EncodeToL2QR(randpayload, *qrconf)
		if err != nil {
			t.Fatal(err)
		}
		qrcodes := make([][]byte, 0, result.TotalSegment)
		for currentQRcode := 0; currentQRcode < result.TotalSegment; currentQRcode++ {
			qrp, err := result.GetSegmentAsByte(currentQRcode)
			if err != nil {
				t.Fatal(err)
			}
			qrcodes = append(qrcodes, qrp)
		}

		if err != nil {
			t.Fatal(err)
		}

		return qrcodes
	}()
	encoder := &vlencoding.Encoder{}
	qrdecoder := encoder.StartQRDecode()
	for _, currentScan := range qrcode {
		qrdecoder.OnNewDataScanned(currentScan)
	}

	if !qrdecoder.IsDecodeReady() {
		t.Fatal("Asked more segment than necessary")
	}
	output, err := qrdecoder.Finish()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(output)
	if !bytes.Equal(output, randpayload) {
		t.Fatal("Payload Mismatch")
	}
}

func TestDecodeQRBig(t *testing.T) {
	t.SkipNow()
	payloadsize := 65536
	randpayload := make([]byte, payloadsize)
	_, err := rand.Read(randpayload)
	qrcode := func() [][]byte {
		qrconf := &vlencoding.QRGenConf{}
		encoder := &vlencoding.Encoder{}
		qrconf.ForceReconstruct = false
		qrconf.MaxQrSize = 1024
		qrconf.ReconsConf = &vlencoding.QRGenConf_ReconstructConf{AtLeastMutiply: 2.0, AtLeastReplacement: 2}
		if err != nil {
			t.Fatal(err)
		}
		result, err := encoder.EncodeToL2QR(randpayload, *qrconf)
		if err != nil {
			t.Fatal(err)
		}
		qrcodes := make([][]byte, 0, result.TotalSegment)
		for currentQRcode := 0; currentQRcode < result.TotalSegment; currentQRcode++ {
			qrp, err := result.GetSegmentAsByte(currentQRcode)
			if err != nil {
				t.Fatal(err)
			}
			qrcodes = append(qrcodes, qrp)
		}

		if err != nil {
			t.Fatal(err)
		}

		return qrcodes
	}()
	encoder := &vlencoding.Encoder{}
	qrdecoder := encoder.StartQRDecode()
	for _, currentScan := range qrcode {
		qrdecoder.OnNewDataScanned(currentScan)
	}

	if !qrdecoder.IsDecodeReady() {
		t.Fatal("Asked more segment than necessary")
	}
	output, err := qrdecoder.Finish()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(output)
	if !bytes.Equal(output, randpayload) {
		t.Fatal("Payload Mismatch")
	}
}

func TestDecodeQRAllSmallSize(t *testing.T) {
	t.SkipNow()
	const MAX_QR_SIZE = 4096

	for QRSizeTesting := 1; QRSizeTesting < MAX_QR_SIZE; QRSizeTesting++ {
		randpayload := make([]byte, QRSizeTesting)
		_, err := rand.Read(randpayload)
		qrcode := func() [][]byte {
			qrconf := &vlencoding.QRGenConf{}
			encoder := &vlencoding.Encoder{}
			qrconf.ForceReconstruct = false
			qrconf.MaxQrSize = 1024
			qrconf.ReconsConf = &vlencoding.QRGenConf_ReconstructConf{AtLeastMutiply: 2.0, AtLeastReplacement: 2}
			if err != nil {
				t.Fatal(err)
			}
			result, err := encoder.EncodeToL2QR(randpayload, *qrconf)
			if err != nil {
				t.Fatal(err)
			}
			qrcodes := make([][]byte, 0, result.TotalSegment)
			for currentQRcode := 0; currentQRcode < result.TotalSegment; currentQRcode++ {
				qrp, err := result.GetSegmentAsByte(currentQRcode)
				if err != nil {
					t.Fatal(err)
				}
				qrcodes = append(qrcodes, qrp)
			}

			if err != nil {
				t.Fatal(err)
			}

			return qrcodes
		}()
		encoder := &vlencoding.Encoder{}
		qrdecoder := encoder.StartQRDecode()
		for _, currentScan := range qrcode {
			qrdecoder.OnNewDataScanned(currentScan)
		}

		if !qrdecoder.IsDecodeReady() {
			t.Fatal("Asked more segment than necessary")
		}
		output, err := qrdecoder.Finish()
		if err != nil {
			t.Fatal(err)
		}
		t.Log(output)
		if !bytes.Equal(output, randpayload) {
			t.Fatal("Payload Mismatch")
		}
	}

}

func TestDecodeQRAllSizeSkip(t *testing.T) {
	t.SkipNow()
	const MAX_QR_SIZE = 65536

	for QRSizeTesting := 1; QRSizeTesting < MAX_QR_SIZE; QRSizeTesting += 1234 {
		randpayload := make([]byte, QRSizeTesting)
		_, err := rand.Read(randpayload)
		qrcode := func() [][]byte {
			qrconf := &vlencoding.QRGenConf{}
			encoder := &vlencoding.Encoder{}
			qrconf.ForceReconstruct = false
			qrconf.MaxQrSize = 1024
			qrconf.ReconsConf = &vlencoding.QRGenConf_ReconstructConf{AtLeastMutiply: 2.0, AtLeastReplacement: 2}
			if err != nil {
				t.Fatal(err)
			}
			result, err := encoder.EncodeToL2QR(randpayload, *qrconf)
			if err != nil {
				t.Fatal(err)
			}
			qrcodes := make([][]byte, 0, result.TotalSegment)
			for currentQRcode := 0; currentQRcode < result.TotalSegment; currentQRcode++ {
				qrp, err := result.GetSegmentAsByte(currentQRcode)
				if err != nil {
					t.Fatal(err)
				}
				qrcodes = append(qrcodes, qrp)
			}

			if err != nil {
				t.Fatal(err)
			}

			return qrcodes
		}()
		encoder := &vlencoding.Encoder{}
		qrdecoder := encoder.StartQRDecode()
		for _, currentScan := range qrcode {
			qrdecoder.OnNewDataScanned(currentScan)
		}

		if !qrdecoder.IsDecodeReady() {
			t.Fatal("Asked more segment than necessary")
		}
		output, err := qrdecoder.Finish()
		if err != nil {
			t.Fatal(err)
		}
		t.Log(output)
		if !bytes.Equal(output, randpayload) {
			t.Fatal("Payload Mismatch")
		}
	}

}

func Sample(orig [][]byte, returnsum int) [][]byte {
	slice := orig
	//shuffle
	for i := range slice {
		j := rand.Intn(i + 1)
		slice[i], slice[j] = slice[j], slice[i]
	}
	return slice[0:returnsum]
}

func TestIngTestSample(t *testing.T) {
	t.Skip()
	r := make([][]byte, 5)
	for re := range r {
		r[re] = make([]byte, 10)
		rand.Read(r[re])
	}
	t.Log(r)
	rs := Sample(r, 1)
	t.Log(rs)
	if reflect.DeepEqual(rs, r) {
		t.FailNow()
	}
}

func TestDecodeQRManyReconstruct(t *testing.T) {
	t.SkipNow()
	payloadsize := 512
	randpayload := make([]byte, payloadsize)
	_, err := rand.Read(randpayload)
	qrcode := func() [][]byte {
		qrconf := &vlencoding.QRGenConf{}
		encoder := &vlencoding.Encoder{}
		qrconf.ForceReconstruct = true
		qrconf.MaxQrSize = 1024
		qrconf.ReconsConf = &vlencoding.QRGenConf_ReconstructConf{AtLeastMutiply: 2.0, AtLeastReplacement: 2}
		if err != nil {
			t.Fatal(err)
		}
		result, err := encoder.EncodeToL2QR(randpayload, *qrconf)
		if err != nil {
			t.Fatal(err)
		}
		qrcodes := make([][]byte, 0, result.TotalSegment)
		for currentQRcode := 0; currentQRcode < result.TotalSegment; currentQRcode++ {
			qrp, err := result.GetSegmentAsByte(currentQRcode)
			if err != nil {
				t.Fatal(err)
			}
			qrcodes = append(qrcodes, qrp)
		}
		if err != nil {
			t.Fatal(err)
		}
		//Remove some of the reconstructable QR
		necesarySegment := result.PayloadSegment
		return Sample(qrcodes, necesarySegment)
	}()
	encoder := &vlencoding.Encoder{}
	qrdecoder := encoder.StartQRDecode()
	for _, currentScan := range qrcode {
		qrdecoder.OnNewDataScanned(currentScan)
	}

	if !qrdecoder.IsDecodeReady() {
		t.Fatal("Asked more segment than necessary")
	}
	output, err := qrdecoder.Finish()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(output)
	if !bytes.Equal(output, randpayload) {
		t.Fatal("Payload Mismatch")
	}
}

func TestDecodeQRBigReconstruct(t *testing.T) {
	t.SkipNow()
	payloadsize := 4097
	randpayload := make([]byte, payloadsize)
	_, err := rand.Read(randpayload)
	qrcode := func() [][]byte {
		qrconf := &vlencoding.QRGenConf{}
		encoder := &vlencoding.Encoder{}
		qrconf.ForceReconstruct = false
		qrconf.MaxQrSize = 1024
		qrconf.ReconsConf = &vlencoding.QRGenConf_ReconstructConf{AtLeastMutiply: 2.0, AtLeastReplacement: 2}
		if err != nil {
			t.Fatal(err)
		}
		result, err := encoder.EncodeToL2QR(randpayload, *qrconf)
		if err != nil {
			t.Fatal(err)
		}
		qrcodes := make([][]byte, 0, result.TotalSegment)
		for currentQRcode := 0; currentQRcode < result.TotalSegment; currentQRcode++ {
			qrp, err := result.GetSegmentAsByte(currentQRcode)
			if err != nil {
				t.Fatal(err)
			}
			qrcodes = append(qrcodes, qrp)
		}
		if err != nil {
			t.Fatal(err)
		}
		//Remove some of the reconstructable QR
		necesarySegment := result.PayloadSegment
		return Sample(qrcodes, necesarySegment)
	}()
	encoder := &vlencoding.Encoder{}
	qrdecoder := encoder.StartQRDecode()
	for _, currentScan := range qrcode {
		qrdecoder.OnNewDataScanned(currentScan)
	}

	if !qrdecoder.IsDecodeReady() {
		t.Fatal("Asked more segment than necessary")
	}
	output, err := qrdecoder.Finish()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(output)
	if !bytes.Equal(output, randpayload) {
		t.Fatal("Payload Mismatch")
	}
}

func TestErrHandleingQRMisMatch(t *testing.T) {
	t.SkipNow()
	payloadsize := 4097
	randpayload := make([]byte, payloadsize)
	_, err := rand.Read(randpayload)
	qrcodeg := func() [][]byte {
		qrconf := &vlencoding.QRGenConf{}
		encoder := &vlencoding.Encoder{}
		qrconf.ForceReconstruct = false
		qrconf.MaxQrSize = 1024
		qrconf.ReconsConf = &vlencoding.QRGenConf_ReconstructConf{AtLeastMutiply: 2.0, AtLeastReplacement: 2}
		if err != nil {
			t.Fatal(err)
		}
		result, err := encoder.EncodeToL2QR(randpayload, *qrconf)
		if err != nil {
			t.Fatal(err)
		}
		qrcodes := make([][]byte, 0, result.TotalSegment)
		for currentQRcode := 0; currentQRcode < result.TotalSegment; currentQRcode++ {
			qrp, err := result.GetSegmentAsByte(currentQRcode)
			if err != nil {
				t.Fatal(err)
			}
			qrcodes = append(qrcodes, qrp)
		}
		if err != nil {
			t.Fatal(err)
		}
		//Remove some of the reconstructable QR
		necesarySegment := result.PayloadSegment
		return Sample(qrcodes, necesarySegment)
	}
	qrcode := qrcodeg()
	randpayloado := randpayload
	randpayload = make([]byte, payloadsize)
	_, err = rand.Read(randpayload)
	qrcodm := qrcodeg()
	encoder := &vlencoding.Encoder{}
	qrdecoder := encoder.StartQRDecode()
	for _, currentScan := range qrcode[0:2] {
		qrdecoder.OnNewDataScanned(currentScan)
	}

	err = qrdecoder.OnNewDataScanned(qrcodm[0])
	if err == nil {
		t.FailNow()
	}

	for _, currentScan := range qrcode[2:len(qrcode)] {
		qrdecoder.OnNewDataScanned(currentScan)
	}

	if !qrdecoder.IsDecodeReady() {
		t.Fatal("Asked more segment than necessary")
	}
	output, err := qrdecoder.Finish()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(output)
	if !bytes.Equal(output, randpayloado) {
		t.Fatal("Payload Mismatch")
	}
}

func TestErrHandleingQRInsuffent(t *testing.T) {
	t.SkipNow()
	payloadsize := 4097
	randpayload := make([]byte, payloadsize)
	_, err := rand.Read(randpayload)
	qrcode := func() [][]byte {
		qrconf := &vlencoding.QRGenConf{}
		encoder := &vlencoding.Encoder{}
		qrconf.ForceReconstruct = false
		qrconf.MaxQrSize = 1024
		qrconf.ReconsConf = &vlencoding.QRGenConf_ReconstructConf{AtLeastMutiply: 2.0, AtLeastReplacement: 2}
		if err != nil {
			t.Fatal(err)
		}
		result, err := encoder.EncodeToL2QR(randpayload, *qrconf)
		if err != nil {
			t.Fatal(err)
		}
		qrcodes := make([][]byte, 0, result.TotalSegment)
		for currentQRcode := 0; currentQRcode < result.TotalSegment; currentQRcode++ {
			qrp, err := result.GetSegmentAsByte(currentQRcode)
			if err != nil {
				t.Fatal(err)
			}
			qrcodes = append(qrcodes, qrp)
		}
		if err != nil {
			t.Fatal(err)
		}
		//Remove some of the reconstructable QR
		necesarySegment := result.PayloadSegment
		return Sample(qrcodes, necesarySegment-1)
	}()
	encoder := &vlencoding.Encoder{}
	qrdecoder := encoder.StartQRDecode()
	for _, currentScan := range qrcode {
		qrdecoder.OnNewDataScanned(currentScan)
	}

	if qrdecoder.IsDecodeReady() {
		t.Fatal("Not Complain Insuffent segment")
	}
	output, err := qrdecoder.Finish()
	if err == nil {
		t.Fatal("Not Complain Insuffent segment")
	}
	t.Log(output)

}
