package ws

import (
	"math/rand"
	"time"
)

func GetRandomStri(l int) string {

	rand.Seed(time.Now().UnixNano())
	result := make([]byte, l)

	for i := 0; i < l; i++ {
		rand := rand.Intn(62)
		if rand < 10 {
			result[i] = byte(48 + rand)
		} else if rand < 36 {
			result[i] = byte(55 + rand)
		} else {
			result[i] = byte(61 + rand)
		}
	}
	return string(result)
}
