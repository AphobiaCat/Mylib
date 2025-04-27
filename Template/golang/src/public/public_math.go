package public

import (
	"crypto/sha256"
	"github.com/mr-tron/base58"
	"math/rand"
	"time"
	"unsafe"
)

func Rand(max_num int)int{
	return rand.Intn(max_num)
}

func Rand_U64(max_num uint64)uint64{
	return rand.Uint64()
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

type quick_sort_type interface {
	GetNum() int32
}

func QuickSortStruct[T quick_sort_type](sort_array []T) []T {
	now_pos := 0
	i := 0
	k := len(sort_array) - 1

	for i < k {
		for now_pos < k {
			if sort_array[now_pos].GetNum() <= sort_array[k].GetNum() {
				k -= 1
				continue
			} else {
				tmp := sort_array[k]
				sort_array[k] = sort_array[now_pos]
				sort_array[now_pos] = tmp
				now_pos = k
				k -= 1
				break
			}
		}

		for now_pos > i {
			if sort_array[now_pos].GetNum() >= sort_array[i].GetNum() {
				i += 1
				continue
			} else {
				tmp := sort_array[i]
				sort_array[i] = sort_array[now_pos]
				sort_array[now_pos] = tmp
				now_pos = i
				i += 1
				break
			}
		}
	}

	var front, last []T

	if now_pos > 0 {
		front = QuickSortStruct(sort_array[0:now_pos])
	}
	if now_pos+1 < len(sort_array) {
		last = QuickSortStruct(sort_array[now_pos+1 : len(sort_array)])
	}

	ret := front
	ret = append(ret, sort_array[now_pos])
	ret = append(ret, last...)

	return ret
}

func Encode(origin interface{}, encode_key uint64) []uint8 {
	jsonBytes := []uint8(Build_Json(origin))

	encode_key_array := (*[8]uint8)(unsafe.Pointer(&encode_key))

	for i := range jsonBytes {
		jsonBytes[i] ^= encode_key_array[i%8]
	}
	return jsonBytes
}

func Decode(encoded []uint8, decode_key uint64) string {
	decode_key_array := (*[8]uint8)(unsafe.Pointer(&decode_key))

	decoded := make([]uint8, len(encoded))
	for i, b := range encoded {
		decoded[i] = b ^ decode_key_array[i%8]
	}
	return string(decoded)
}

func init(){
	rand.Seed(time.Now().UnixNano())
}

