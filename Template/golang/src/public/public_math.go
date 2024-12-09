package public

import (
	"crypto/sha256"
	"github.com/mr-tron/base58"
	"math/rand"
)

func Rand(max_num int)int{
	return rand.Intn(max_num)
}

func Sha256(encode_str string)[32]byte{

	DBG_LOG("encode str:", encode_str, "=", []byte(encode_str))

	return sha256.Sum256([]byte(encode_str))
}

func Base58(hash [32]byte)string{
	return base58.Encode(hash[:])
}

func Uint8ToUint64(u8array []byte)uint64{

	if len(u8array) < 8{
		return uint64(0)
	}

	ret := uint64(0)

	ret = uint64(u8array[0]) << 56
	ret |= uint64(u8array[1]) << 48
	ret |= uint64(u8array[2]) << 40
	ret |= uint64(u8array[3]) << 32
	ret |= uint64(u8array[4]) << 24
	ret |= uint64(u8array[5]) << 16
	ret |= uint64(u8array[6]) << 8
	ret |= uint64(u8array[7])
	
	return ret
}

func Base58Hash2Uint64(base58Hash string)[]uint64{

	hashBytes, err := base58.Decode(base58Hash)
	if err != nil {
		DBG_ERR("Failed to decode Base58 hash: %v", err)
	}

	if len(hashBytes) != 32 {
		DBG_ERR("Hash must be 32 bytes, but got %d bytes", len(hashBytes))
	}

	ret := []uint64{}
	

	for i := 0; i < 4; i++ {

		chunk := hashBytes[i*8 : (i+1)*8]

		DBG_LOG(chunk)
	
		part := Uint8ToUint64(chunk)

		ret = append(ret, part)
	}
	

	return ret

}

