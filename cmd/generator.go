package main

import (
	"fmt"
	"math/rand"
	"net"
	"os"
	"time"

	"github.com/camathieu/accessmon"
)

var ips = []net.IP{net.ParseIP("127.0.0.1"), net.ParseIP("::1"), net.ParseIP("1.1.1.1")}
var users = []string{"user1", "user2", "user3"}
var methods = []string{"GET", "POST", "DELETE"}
var paths = []string{"/api", "/www", "/static"}
var versions = []string{"HTTP/1.1", "HTTP/2.2"}
var codes = []int{200, 301, 400, 404, 500}

func generateRandomRequest() (req *accessmon.Request) {

	req = &accessmon.Request{
		SourceIP:    ips[rand.Intn(len(ips))],
		User:        users[rand.Intn(len(users))],
		Time:        time.Now(),
		Method:      methods[rand.Intn(len(methods))],
		Path:        paths[rand.Intn(len(paths))],
		HTTPVersion: versions[rand.Intn(len(versions))],
		Code:        codes[rand.Intn(len(codes))],
		Size:        rand.Intn(100),
	}

	return req
}

func generator(path string) (err error) {
	rand.Seed(time.Now().Unix()) // initialize global pseudo random generator

	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}

	ticker := time.Tick(time.Second)

	maxSpeed := 100
	speed := 50
	raise := true
	for {
		<-ticker

		fmt.Printf("writing to %s at %d requests per seconds\n", path, speed)

		for i := 0; i < speed; i++ {

			req := generateRandomRequest()

			_, err = file.WriteString(req.String() + "\n")
			if err != nil {
				return err
			}
		}

		if raise {
			speed++
			if speed >= maxSpeed {
				raise = false
			}
		} else {
			speed--
			if speed <= 0 {
				raise = true
			}
		}
	}
}
