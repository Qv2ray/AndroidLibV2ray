package libv2ray

import (
	"fmt"

	"github.com/xiaokangwang/AndroidLibV2Ray/CoreI"

	core "github.com/v2fly/v2ray-core/v4"
)

/*CheckVersion int
This func will return libv2ray binding version.
*/
func CheckVersion() int {
	return CoreI.CheckVersion()
}

/*CheckVersionX string
This func will return libv2ray binding version and V2Ray version used.
*/
func CheckVersionX() string {
	return fmt.Sprintf("Libv2ray rev. %d, along with V2Ray %s", CheckVersion(), core.Version())
}
