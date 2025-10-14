// generated using: go run ./gen/ algorithms.go algorithms.gen.go
package httpsig

import "crypto"

func stringToHash(name string) crypto.Hash {
	switch name {
	case md4String:
		return crypto.MD4
	case md5String:
		return crypto.MD5
	case ripemd160String:
		return crypto.RIPEMD160
	case md5sha1String:
		return crypto.MD5SHA1
	case sha1String:
		return crypto.SHA1
	case sha224String:
		return crypto.SHA224
	case sha256String:
		return crypto.SHA256
	case sha384String:
		return crypto.SHA384
	case sha512String:
		return crypto.SHA512
	case sha3_224String:
		return crypto.SHA3_224
	case sha3_256String:
		return crypto.SHA3_256
	case sha3_384String:
		return crypto.SHA3_384
	case sha3_512String:
		return crypto.SHA3_512
	case sha512_224String:
		return crypto.SHA512_224
	case sha512_256String:
		return crypto.SHA512_256
	case blake2s_256String:
		return crypto.BLAKE2s_256
	case blake2b_256String:
		return crypto.BLAKE2b_256
	case blake2b_384String:
		return crypto.BLAKE2b_384
	case blake2b_512String:
		return crypto.BLAKE2b_512
	default:
		return 0
	}
}
