package utils
import (
	"crypto/rc4"
	"encoding/base64"
	"errors"
	"sort"
	"strings"
	"unicode"
	"github.com/redis/go-redis/v9"
)

type Step struct {
	order int
	method string
	keys []string
}

type Keys struct {
	steps []Step
}
type ByOrderAscending []Step

func (a ByOrderAscending) Len() int           { return len(a) }
func (a ByOrderAscending) Less(i, j int) bool { return a[i].order < a[j].order }
func (a ByOrderAscending) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

type ByOrderDescending []Step

func (a ByOrderDescending) Len() int           { return len(a) }
func (a ByOrderDescending) Less(i, j int) bool { return a[i].order > a[j].order }
func (a ByOrderDescending) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func (k *Keys) SortByOrderAscending() {
	sort.Sort(ByOrderAscending(k.steps))
}

func (k *Keys) SortByOrderDescending() {
	sort.Sort(ByOrderDescending(k.steps))
}

var keys Keys
func SetupKeys(){
	exchange1 := Step{order: 1,method: "exchange",keys: []string{"5j6Ak1GJaTy8XoC","56kC8jyGoXTAa1J"}}
	reverse1 := Step{order: 2,method: "reverse",keys: []string{}}
	rc4_1 := Step{order: 3,method: "rc4",keys: []string{"hAGMmLFnoa0"}}

	exchange2 := Step{order: 4,method: "exchange",keys: []string{"PUoVzgdK5FLZt", "FVogUPtKzdZL5"}}
	reverse2 := Step{order: 5,method: "reverse",keys: []string{}}
	rc4_2 := Step{order: 6,method: "rc4",keys: []string{"oUHxby23izOI5"}}

	exchange3 := Step{order: 7,method: "exchange",keys: []string{"PEQmieNvWhrOX","OEehvmXQrWiPN"}}
	reverse3 := Step{order: 8,method: "reverse",keys: []string{}}
	rc4_3 := Step{order: 9,method: "rc4",keys: []string{"tX6D4K8mPrq3V"}}
	base64 := Step{order: 10,method: "base64",keys: []string{}}
	keys = Keys{
		steps: []Step{
			exchange1,reverse1,rc4_1,exchange2,reverse2,rc4_2,exchange3,reverse3,rc4_3,base64,
		},
	}
}

func exchange(inputStr, key1, key2 string) string {
	key1Dict := make(map[rune]rune)
	for i, char := range key1 {
		if i < len(key2) {
			key1Dict[char] = rune(key2[i])
		}
	}

	result := make([]rune, len(inputStr))
	for i, char := range inputStr {
		if replacement, ok := key1Dict[char]; ok {
			result[i] = replacement
		} else {
			result[i] = char
		}
	}

	return string(result)
}

func reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

func rc4_encrypt(input_data , input_key string) (error,string) {
	key := []byte(input_key)
	data := []byte(input_data)
	
	cipher, err := rc4.NewCipher(key)
	if err != nil {
		return errors.New("Cipher intilaze error"),""
	}
	ciphertext := make([]byte,len(data))
	cipher.XORKeyStream(ciphertext,data)
	encoded := base64.StdEncoding.EncodeToString(ciphertext)
	return nil,string(encoded)
}

func sanitizeBase64URL(encoded_data string) string {
    var sanitized strings.Builder
    for _, ch := range encoded_data {
        if unicode.IsLetter(ch) || unicode.IsDigit(ch) || ch == '-' || ch == '_' || ch == '=' {
            sanitized.WriteRune(ch)
        }
    }
    return sanitized.String()
}

func rc4_decrypt(encoded_data, input_key string) (error, string) {
	key := []byte(input_key)

	encoded_data= sanitizeBase64URL(encoded_data)
	ciphertext, err := base64.URLEncoding.DecodeString(encoded_data)
	if err != nil {
		println(err.Error())
		return errors.New("Base64 decode error"), ""
	}

	cipher, err := rc4.NewCipher(key)
	if err != nil {
		return errors.New("Cipher initialize error"), ""
	}
	plaintext := make([]byte, len(ciphertext))
	cipher.XORKeyStream(plaintext, ciphertext)
	return nil, string(plaintext)
}

func Vrf_encrypt(data string) string  {
	keys.SortByOrderAscending()

	for _,x := range keys.steps {
		switch x.method {
		case "exchange":
			data = exchange(data,x.keys[0],x.keys[1])
		case "reverse":
			data = reverse(data)
		case "rc4":
			err,tempdata := rc4_encrypt(data,x.keys[0])
			if err != nil {
				return ""
			}
			data = tempdata
		case "base64":
			data = base64.StdEncoding.EncodeToString([]byte(data))
		}
	}
	return data
}

func Vrf_decrypt(data string) string  {
	keys.SortByOrderDescending()
	for _,x := range keys.steps {
		switch x.method {
			case "exchange":
				data = exchange(data,x.keys[1],x.keys[0])
			case "reverse":
				data = reverse(data)
			case "rc4":
				err,tempdata := rc4_decrypt(data,x.keys[0])
				if err != nil {
					return "error on rc4 :" + err.Error()
				}
				data = tempdata
			case "base64":
				tempdata,_ := base64.URLEncoding.DecodeString(data)
				data = string(tempdata)
			}
	}
	return data
}
func RedisClient() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:            "localhost:6379",
		Password:        "",
		DB:              0,
		DisableIndentity: true,
	})

    return rdb
}