package sign

import (
	"crypto/md5"
	"encoding/hex"
)

//计算MD5值
func ToMd5(s string) string {
	m := md5.New()
	m.Write([]byte(s))
	return hex.EncodeToString(m.Sum(nil))
}
