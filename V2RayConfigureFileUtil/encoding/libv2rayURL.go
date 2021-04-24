package encoding

import (
	"strings"

	"github.com/bproctor/base91"
)

const LibV2RayURLSignature = "libv2ray:?"

func (e *Encoder) ByteToV2RayURL(vu []byte) string {
	res := base91.EncodeToString(vu)
	return LibV2RayURLSignature + res
}

func (e *Encoder) V2RayURLToByte(url string) []byte {
	if !strings.HasPrefix(url, LibV2RayURLSignature) {
		return nil
	}
	ctx := []byte(url[len(LibV2RayURLSignature):len(url)])
	ctxd := base91.Decode(ctx)
	return ctxd
}
