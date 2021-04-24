package Convert

import (
	"fmt"

	"github.com/bitly/go-simplejson"
)

func JsonConvert(input []byte) ([]byte, error) {
	json, err := simplejson.NewJson(input)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	//Load Addititonal resource json
	inbjsona, err := Asset("ShippedBinary/inbound.json")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	inbjsonaj, err := simplejson.NewJson(inbjsona)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	libjsona, err := Asset("ShippedBinary/libv2ray.json")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	libjsonaj, err := simplejson.NewJson(libjsona)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	//Begin json substitute

	json.Set("inbound", inbjsonaj)
	json.Set("#lib2ray", libjsonaj)

	//Begin Postprogressing (if any)
	//Currently, there is none

	return json.MarshalJSON()
}
