package utils

import (
	"math/rand"
	"time"
)

const (
	defaultLength = 10
)

var (
	letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	source      = rand.NewSource(time.Now().UnixNano())
	random      = rand.New(source)
)

func RandSlug() string {
	sid := make([]rune, defaultLength)
	for i := range sid {
		sid[i] = letterRunes[random.Intn(len(letterRunes))]
	}
	return string(sid)
}

//func RandSlug() string {
//	token := ""
//	for i := 1; i < defaultLength; i++ {
//		token += strconv.Itoa(random.Intn(9))
//	}
//	return token
//}
