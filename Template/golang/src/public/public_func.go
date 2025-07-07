package public

import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func Now_Time_S() int64 {
	now := time.Now()
	seconds := now.Unix()
	return seconds
}

func Now_Time_MS() int64 {
	ms := time.Now().UnixMilli()
	return ms
}

func Now_Time_US() int64 {
	us := time.Now().UnixMicro()
	return us
}

func Now_Time_NS() int64 {
	ns := time.Now().UnixNano()
	return ns
}

func Sleep(sleep_ms int) {
	time.Sleep(time.Duration(sleep_ms) * time.Millisecond)
}

func DBG_LOG_VAR(v interface{}) {

	pc, file, line, ok := runtime.Caller(1)
	if !ok {
		fmt.Println("Failed to get caller information")
		return
	}

	fn := runtime.FuncForPC(pc)
	if fn == nil {
		fmt.Println("Failed to get function information")
		return
	}

	path := file
	filename := filepath.Base(path)

	var outputStr string = "[info] file[" + filename + "]\t| func[" + fn.Name() + "]\t| line[" + ConvertToString(line) + "]\t| arg:"

	typ := reflect.TypeOf(v)

	outputStr += ConvertToString(typ)

	for i := 0; i < typ.NumMethod(); i++ {
		method := typ.Method(i)
		outputStr += "\nmethod[" + ConvertToString(method.Name) + "] \t\t\t method type[" + ConvertToString(method.Type) + "]"

	}

	fmt.Println(outputStr)
}

func DBG_LOG(v ...interface{}) {

	if Config.Mode != "pre" {
		return
	}

	pc, file, line, ok := runtime.Caller(1)
	if !ok {
		fmt.Println("Failed to get caller information")
		return
	}

	fn := runtime.FuncForPC(pc)
	if fn == nil {
		fmt.Println("Failed to get function information")
		return
	}

	path := file
	filename := filepath.Base(path)

	var outputStr string = "[info] file[" + filename + "]\t| func[" + fn.Name() + "]\t| line[" + ConvertToString(line) + "]\t| log:"

	for _, val := range v {
		outputStr += ConvertToString(val)
	}

	fmt.Println(outputStr)

}

func DBG_ERR(v ...interface{}) {

	pc, file, line, ok := runtime.Caller(1)
	if !ok {
		fmt.Println("Failed to get caller information")
		return
	}

	fn := runtime.FuncForPC(pc)
	if fn == nil {
		fmt.Println("Failed to get function information")
		return
	}

	path := file
	filename := filepath.Base(path)

	var outputStr string = "[error] file[" + filename + "]\t| func[" + fn.Name() + "]\t| line[" + ConvertToString(line) + "]\t| log:"

	for _, val := range v {
		outputStr += ConvertToString(val)
	}

	fmt.Println("\033[1;31m" + outputStr + "\033[0m")
}

func SplitStrAfterChar(str string, cutAfter rune) (string, string) {

	for i, char_ := range str {
		if char_ == cutAfter {
			return str[:i], str[i+1:]
		}
	}

	return str, ""
}

func ConvertToString(v interface{}) string {
	return fmt.Sprintf("%v", v)
}

func StrMsgToUint32Array(msg string) (int, []uint32) {
	ret_len := len(msg)/4 + 1
	ret := make([]uint32, ret_len)

	for i, _ := range ret {

		msg_len := len(msg)

		//fmt.Println(msg, "--", msg_len)

		if msg_len >= 4 {
			u1 := uint32(msg[0])
			u2 := uint32(msg[1])
			u3 := uint32(msg[2])
			u4 := uint32(msg[3])

			ret[i] = (u1 << 24) | (u2 << 16) | (u3 << 8) | u4
		} else {
			var u1, u2, u3 uint32 = 0, 0, 0

			switch msg_len {
			case 3:
				{
					u1 = uint32(msg[0])
					u2 = uint32(msg[1])
					u3 = uint32(msg[2])

					break
				}
			case 2:
				{
					u1 = uint32(msg[0])
					u2 = uint32(msg[1])

					break
				}
			case 1:
				{
					u1 = uint32(msg[0])

					break
				}
			}

			ret[i] = (u1 << 24) | (u2 << 16) | (u3 << 8)
		}
		if msg_len >= 4 {
			msg = msg[4:]
		}
	}

	return ret_len, ret
}

func HexStrMsgToUint32Array(msg string) (int, []uint32) {
	ret_len := len(msg)/8 + 1
	ret := make([]uint32, ret_len)

	var tmpNum uint32 = 0
	var i int
	var val rune
	for i, val = range msg {

		tmpNum <<= 4

		if val >= '0' && val <= '9' {
			tmpNum += uint32(val - '0')
		} else if val >= 'a' && val <= 'f' {
			tmpNum += uint32(val - 'a' + 10)
		} else if val >= 'A' && val <= 'F' {
			tmpNum += uint32(val - 'A' + 10)
		}

		if (i+1)%8 == 0 {
			ret[i/8] = tmpNum
			tmpNum = 0
		}
	}

	if tmpNum != 0 {
		tmpNum <<= 4 * (7 - i%8)
		ret[ret_len-1] = tmpNum
	}

	return ret_len, ret
}

func RevertUint32ToStr(num_array []uint32) string {

	ret := ""

	for _, val := range num_array {

		u1 := uint8(val >> 24)
		u2 := uint8((val >> 16) & 0xFF)
		u3 := uint8((val >> 8) & 0xFF)
		u4 := uint8(val & 0xFF)

		if u1 == 0 {
			break
		} else if u2 == 0 {
			ret += string(u1)
		} else if u3 == 0 {
			ret += string(u1) + string(u2)
		} else if u4 == 0 {
			ret += string(u1) + string(u2) + string(u3)
		} else {
			ret += string(u1) + string(u2) + string(u3) + string(u4)
		}
	}

	return ret
}

func ConvertNumToStr(num int64) string {
	return strconv.FormatInt(num, 10)
}

func ConvertNumToHexStr(num int64) string {
	return "0x" + strconv.FormatInt(num, 16)
}

func ConvertStrToNum(num string) int64 {
	ret, _ := strconv.ParseInt(num, 10, 64)
	return ret
}

func ConvertHEXStrToNum(num string) int64 {

	if len(num) > 2 && num[0:2] == "0x" {
		num = num[2:]
	}

	ret, _ := strconv.ParseInt(num, 16, 64)
	return ret
}

func next_str(T string) []int {
	next := make([]int, len(T)+1)

	i := 1
	j := 0

	for i < len(T) {

		if j == 0 || T[i-1] == T[j-1] {
			i++
			j++
			next[i] = j

		} else {
			j = next[j]
		}
	}

	return next
}

func Include_Str(S string, T string) int {
	next := next_str(T)

	i := 1
	j := 1

	for i <= len(S) && j <= len(T) {

		if j == 0 || S[i-1] == T[j-1] {
			i++
			j++
		} else {
			j = next[j]
		}
	}

	if j > len(T) {
		return i - len(T)
	}

	return -1
}

func Split(str string, split_by string) []string {
	parts := strings.Split(str, split_by)
	return parts
}
