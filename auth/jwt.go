package auth

import (
	_ "embed"
)

// pemファイルを読み込み、バイナリに埋め込むようにする

//go:embed cert/secret.pem
var rawPrivKey []byte

//go:embed cert/public.pem
var rawPubKey []byte
