package id

import "encoding/base64"

func String(id []byte) string {
	return base64.StdEncoding.EncodeToString(id)
}

func Id(id []byte) [16]byte {
	return [16]byte(id)
}
