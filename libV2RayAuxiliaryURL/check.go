package libV2RayAuxiliaryURL

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"text/template"
)

const vmessprefix = "vmess://"

func TryRender(url string) (bool, string, []byte) {
	if strings.HasPrefix(url, vmessprefix) {
		s, b := renderVmess(url)
		return true, s, b
	}
	return false, "", nil
}

func renderVmess(url string) (string, []byte) {
	//First guess vmess favor
	rem := url[len(vmessprefix):]
	data, err := base64.StdEncoding.DecodeString(rem)
	if err != nil {
		fmt.Println(err)
		return "", nil
	}
	datas := string(data)
	//Now try shadowrocket favor

	rocketqr, err := regexp.Compile(`^((?:[^:])*):((?:[^@])*)@((?:[^:])*):((?:\d)*)$`)

	if err != nil {
		fmt.Println(err)
		return "", nil
	}

	match := rocketqr.FindAllStringSubmatch(datas, 1)
	if match != nil && match[0] != nil {
		//shadowrocket favor vmess, now progressing
		s := match[0]
		d := &VmessUrlCtx{}
		d.Chiper = s[1]
		d.Add = s[2]
		d.Id = s[3]
		//d.Port, err = strconv.Atoi(s[4])
		d.Port = s[4]
		if err != nil {
			fmt.Println(err)
			return "", nil
		}

		return ".json", vmessR(d)
	}
	d := &VmessUrlCtx{}
	err = json.Unmarshal(data, d)
	if err != nil {
		fmt.Println(err)
		return "", nil
	}
	return ".json", vmessR(d)
}

func vmessR(vx *VmessUrlCtx) []byte {
	templ := template.New("vmess,json")
	tdata, err := Asset("include/vmess.json")
	if err != nil {
		fmt.Println(err)
		return nil
	}
	templ, err = templ.Parse(string(tdata))
	if err != nil {
		fmt.Println(err)
		return nil
	}
	var wb bytes.Buffer
	templ.Execute(&wb, vx)
	return wb.Bytes()
}
