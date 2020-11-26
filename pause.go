package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

func main() {
	containerProcessRunAt := time.Now().UnixNano() / 1e3
	id := os.Getenv("id")
	containerCreateAt, _ := strconv.ParseInt(os.Getenv("containerCreateAt"), 10, 64)
	containerProcessStartAt, _ := strconv.ParseInt(os.Getenv("containerProcessStartAt"), 10, 64)
	fmt.Printf("%s, %d, %d, %d\n", id,
		containerProcessStartAt-containerCreateAt,
		containerProcessRunAt-containerProcessStartAt,
		containerProcessRunAt-containerCreateAt)
	time.Sleep(time.Hour)
}
