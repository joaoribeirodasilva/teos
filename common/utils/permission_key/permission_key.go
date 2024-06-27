package permission_key

import (
	"crypto/md5"
	"fmt"
)

func MakePermissionKey(method string, route string) string {
	keyData := fmt.Sprintf("%s_%s", method, route)
	return "perm_" + fmt.Sprintf("%x", md5.Sum([]byte(keyData)))
}
