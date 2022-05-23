/*
1、 go mod init mobileCode
2、go get -u github.com/go-redis/redis
相当与6版本。。

实现功能描述
1、5分钟可以要求3次验证码
2、每次新验证码会覆盖前面一个，验证码有效期1分钟
*/
package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

var client *redis.Client

func initClient() {
	client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)
}

func generateRandStr() string {
	var ret string = ""
	for i := 0; i < 6; i++ {
		ret += strconv.Itoa(rand.Intn(10))
	}
	return ret
}

func requestCode(mobile string) {
	val, err := client.Get(mobile + "_code").Result()

	if err == redis.Nil {
		err1 := client.Set(mobile+"_code", 0, time.Minute*5).Err()
		if err1 != nil {
			fmt.Println("see me2?")
			panic(err1)
		}
		client.Set(mobile+"_code_num", generateRandStr(), time.Minute*1)

	} else if err != nil {
		fmt.Println("see me?")
		panic(err)
	} else {
		intVal, err1 := strconv.Atoi(val)
		if err1 != nil {
			panic(err1)
		}
		if intVal < 2 {
			client.Incr(mobile + "_code")
			client.Set(mobile+"_code_num", generateRandStr(), time.Minute*1)
		}
	}
}

func checkCodeNum(mobile string, input string) bool {
	val, err := client.Get(mobile + "_code_num").Result()
	if err == redis.Nil {
		return false
	} else if err != nil {
		fmt.Println("see me?")
		panic(err)
	} else {
		return val == input
	}
}

func main() {
	initClient()
	rand.Seed(time.Now().UnixNano())
	//fmt.Println(generateRandStr())
	//fmt.Println(generateRandStr())
	//fmt.Println(generateRandStr())
	requestCode("13610088588")
	//requestCode("13610088588")
	//requestCode("13610088588")
	//requestCode("13610088588")

	fmt.Println(checkCodeNum("13610088588", "53823"))
	fmt.Println(checkCodeNum("13610088588", "43423"))
	fmt.Println(checkCodeNum("13610088588", "439050"))
	fmt.Println(checkCodeNum("13610088588", "083308"))
}
