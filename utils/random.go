package utils

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
)

func RandomInt(min, max int) int {
	if min == max {
		return min
	}
	var randInt int64
	for {
		err := binary.Read(rand.Reader, binary.LittleEndian, &randInt)
		if err != nil {
			continue
		}
		randInt = randInt%(int64(max-min+1)) + int64(min)
		if randInt >= int64(min) && randInt <= int64(max) {
			return int(randInt)
		}
	}
}

// 生成指定长度的随机字符串，包含小写字母、大写字母、短横线和下划线
const STR_azAZ09_ = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ-_"
const STR_azAZ09 = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const STR_af19 = "1234567890abcdef"

func RandomString(length int, charset string) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	var sb strings.Builder
	for _, byteValue := range b {
		// 将每个随机字节转换为字符集索引
		index := int(byteValue) % len(charset)
		sb.WriteByte(charset[index])
	}
	return sb.String()
}
func RandomBytes(length int) [32]byte {
	var randomData [32]byte
	for i := 0; i < 32; i++ {
		randomData[i] = byte(RandomInt(0, 255))
	}
	return randomData
}

func RandomHex(length int) string {
	len2 := length
	if len2%2 != 0 {
		len2++
	}
	len2 = len2 / 2
	b := make([]byte, len2) // 长度除以2因为每对16进制字符表示1个字节
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return hex.EncodeToString(b)[:length]
}

func UUIDv4() string {
	uuid := make([]byte, 16)
	// 生成随机UUID
	n, err := rand.Read(uuid)
	if n != len(uuid) || err != nil {
		return ""
	}
	uuid[6] = (uuid[6] & 0x0f) | 0x40
	uuid[8] = (uuid[8] & 0x3f) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:])
}

func RandomArray[V1 any](arr []V1) V1 {
	return arr[RandomInt(0, len(arr)-1)]
}

func RandomPhone(len int) string {
	var ls string = strconv.Itoa(RandomInt(1, 9))
	for i := 0; i < len-1; i++ {
		ls += strconv.Itoa(RandomInt(0, 9))
	}
	return ls
}

func RandomSplitStringArr(s string, p string) string {
	arr := strings.Split(s, p)
	arr = ShuffleStrings(arr)
	arr = ShuffleStrings(arr)
	return strings.Join(arr, p)
}
func ShuffleStrings[V1 any](strs []V1) []V1 {
	// 创建一个副本以避免修改原始数组
	shuffled := make([]V1, len(strs))
	copy(shuffled, strs)

	for i := len(shuffled) - 1; i > 0; i-- {
		// 随机选择一个小于当前索引的数j
		j := RandomInt(0, i)
		// 交换当前索引i和随机索引j的元素
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	}
	return shuffled
}
func RandomStringJoin(j string, j2 string, min, max int) string {
	size := RandomInt(min, max)
	arr := make([]string, size)
	for i := 0; i < size; i++ {
		name := RandomString(RandomInt(2, 8), STR_azAZ09)
		val := RandomString(RandomInt(4, 32), STR_azAZ09_)
		arr = append(arr, name+j+val)
	}
	return strings.Join(arr, j2)
}
