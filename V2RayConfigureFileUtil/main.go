package main

//go:generate make all
import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"image/png"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/xiaokangwang/V2RayConfigureFileUtil/Convert"
	vlencoding "github.com/xiaokangwang/V2RayConfigureFileUtil/encoding"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/aztec"
	qrenc "github.com/boombuler/barcode/qr"
)

func main() {
	var target = flag.String("t", "QR", "output target")
	var input = flag.String("i", "i.json", "input file/dir")
	var output = flag.String("o", "o", "output file/dir")
	flag.Parse()
	switch *target {
	case "QR":
		qr(*input, *output)
	case "seg":
		seg(*input, *output)
	case "AZ":
		az(*input, *output)
	case "rseg":
		rseg(*input, *output)
		//case "fpb":
	case "lv2json":
		lv2json(*input, *output)
	case "jsongrace":
		jsongrace(*input, *output)

	}
}

func seg(input, output string) {
	ec := &vlencoding.Encoder{}
	qrconf := &vlencoding.QRGenConf{}
	qrconf.ForceReconstruct = false
	qrconf.MaxQrSize = 1024
	qrconf.ReconsConf = &vlencoding.QRGenConf_ReconstructConf{AtLeastMutiply: 2.0, AtLeastReplacement: 4}
	data, err := ec.PackV2RayConfigureIntoPackedForm(input)
	if err != nil {
		panic(err)
	}
	qrd, err := ec.EncodeToL2QR(data, *qrconf)
	if err != nil {
		panic(err)
	}
	for currOutput := 0; currOutput < qrd.TotalSegment; currOutput++ {
		qd, err := qrd.GetSegmentAsByte(currOutput)
		if err != nil {
			panic(err)
		}
		err = ioutil.WriteFile(fmt.Sprintf("%v/%v.seg.libv2ray.pb", output, currOutput), qd, 0600)
		if err != nil {
			panic(err)
		}
	}
}

func qr(input, output string) {
	ec := &vlencoding.Encoder{}
	payload, err := ioutil.ReadFile(input)
	if err != nil {
		panic(err)
	}
	url := ec.ByteToV2RayURL(payload)

	qrCode, err := qrenc.Encode(url, qrenc.L, qrenc.Auto)
	if err != nil {
		panic(err)
	}
	qrCode, err = barcode.Scale(qrCode, 1024, 1024)
	if err != nil {
		panic(err)
	}
	var pngbuf bytes.Buffer
	png.Encode(&pngbuf, qrCode)
	ioutil.WriteFile(output, pngbuf.Bytes(), 0600)
}

func az(input, output string) {
	ec := &vlencoding.Encoder{}
	payload, err := ioutil.ReadFile(input)
	if err != nil {
		panic(err)
	}
	url := ec.ByteToV2RayURL(payload)

	qrCode, err := aztec.Encode([]byte(url), 3, 0)
	if err != nil {
		panic(err)
	}
	qrCode, err = barcode.Scale(qrCode, 1024, 1024)
	if err != nil {
		panic(err)
	}
	var pngbuf bytes.Buffer
	png.Encode(&pngbuf, qrCode)
	ioutil.WriteFile(output, pngbuf.Bytes(), 0600)
}

func rseg(input, output string) {
	ec := &vlencoding.Encoder{}
	qd := ec.StartQRDecode()
	filepath.Walk(input, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			payload, err := ioutil.ReadFile(input)
			if err != nil {
				fmt.Println(err)
			}
			err = qd.OnNewDataScanned(payload)
			if err != nil {
				fmt.Println(err)
			}
		}
		return nil
	})
	if !qd.IsDecodeReady() {
		panic(fmt.Sprintf("Decode not ready %v / %v", qd.PieceReceived, qd.PieceNeeded))
	}
	out, err := qd.Finish()
	if err != nil {
		panic(err)
	}
	//Now unwrap file
	ext, pay, err := ec.UnpackV2RayConfB(out)
	if err != nil {
		panic(err)
	}
	//WriteFile
	err = ioutil.WriteFile(output+ext, pay, 0600)
	if err != nil {
		panic(err)
	}
}

func lv2json(input, output string) {
	inputb, err := ioutil.ReadFile(input)
	if err != nil {
		panic(err)
	}
	bin, err := Convert.JsonConvert(inputb)
	if err != nil {
		panic(err)
	}
	ioutil.WriteFile(output, bin, 0600)
}

func jsonprettyprint(b []byte) ([]byte, error) {
	var out bytes.Buffer
	err := json.Indent(&out, b, "", "  ")
	return out.Bytes(), err
}

func jsonmini(b []byte) ([]byte, error) {
	var out bytes.Buffer
	err := json.Indent(&out, b, "", "")
	return out.Bytes(), err
}

func jsongrace(input, output string) {
	inputb, err := ioutil.ReadFile(input)
	if err != nil {
		panic(err)
	}
	bin, err := jsonprettyprint(inputb)
	if err != nil {
		panic(err)
	}
	ioutil.WriteFile(output, bin, 0600)
}
