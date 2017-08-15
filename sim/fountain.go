package main

import (
	"fmt"
	"math/rand"
	"time"
)

const DB_SIZE = 1000000
const FAIL_PRC = 70.0

var db [DB_SIZE]int
var hop_cnt [DB_SIZE]int
var total_msg_cnt = 0
var total_recived = 0
var last_change = time.Now()

func run(index int, done chan bool) {
	if rand.Float64()*100 < FAIL_PRC {
		return
	}

	left1 := (index - 1 + DB_SIZE) % DB_SIZE
	right1 := (index + 1 + DB_SIZE) % DB_SIZE

	for db[index] == -1 {
		time.Sleep(5 * time.Millisecond)
	}
	done <- true
	hop := hop_cnt[index] + 1
	//time to send message
	//sleep_rand := time.Duration(20 + rand.Intn(180))
	sleep_rand := time.Duration(50)
	time.Sleep(sleep_rand * time.Millisecond)

	msg_from_neighbor := db[index] == left1 || db[index] == right1 //not a nighbur

	if msg_from_neighbor {
		if db[index] == left1 { //coming from left pass it to right
			db[right1] = index
			hop_cnt[right1] = hop
		} else { //pass it to left
			db[left1] = index
			hop_cnt[left1] = hop
		}
		total_msg_cnt += 1
	} else { //create a fountain
		for i := 0; i < 4; i++ {
			left := (DB_SIZE - rand.Intn(DB_SIZE/2-2) + DB_SIZE) % DB_SIZE
			right := (DB_SIZE + rand.Intn(DB_SIZE/2+2) + DB_SIZE) % DB_SIZE
			db[left] = index
			hop_cnt[left] = hop
			db[right] = index
			hop_cnt[right] = hop
			total_msg_cnt += 2
		}
		db[left1] = index
		hop_cnt[left1] = hop
		db[right1] = index
		hop_cnt[right1] = hop
		total_msg_cnt += 2
	}
	last_change = time.Now()
	//fmt.Println("Done", index)
}

func test_fountain() {
	//init to empty
	for i := 0; i < len(db); i++ {
		db[i] = -1
	}
	//first fountain
	for i := 0; i < 100; i += 10 {
		db[i] = 1
	}
	//db[3] = 30 // first message
	done := make(chan bool)
	for i := 0; i < DB_SIZE; i++ {
		go run(i, done)
	}

	time.Sleep(500 * time.Millisecond)

	fmt.Println("starting...")

	start := time.Now()
	last_change = time.Now()

	timeout := make(chan bool, 1)
	go func() {
		for i := 0; i < DB_SIZE; i++ {
			<-done
			last_change = time.Now()
			total_recived += 1
			//fmt.Println("So Far", i)
		}
		timeout <- false
	}()

	go func() {
		for {

			if time.Now().Sub(last_change) > time.Duration(500*time.Millisecond) {
				timeout <- true
			}
		}
	}()

	<-timeout

	elapsed := time.Since(start)
	fmt.Println("Total msgs", total_msg_cnt, "DB_DIZE", DB_SIZE, "total_recived=", total_recived)
	fmt.Println("Took", elapsed)
	max := hop_cnt[0]
	for _, e := range hop_cnt {
		if e > max {
			max = e
		}
	}
	fmt.Println("max hops", max)
	//time.Sleep(1000 * time.Millisecond)
}
