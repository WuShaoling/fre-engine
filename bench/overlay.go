package main

/**
overlay fs 创建性能测试
*/

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"syscall"
	"time"
)

const basePath = "overlay_test"
const lowerPath = basePath + "/rootfs"

var count int     // 创建的数量
var parallel bool // 是否并发创建

func init() {
	flag.IntVar(&count, "n", 1, "创建的数量")
	flag.BoolVar(&parallel, "p", false, "并发启动")
}

func createOverlay(id int) []int64 {

	t1 := time.Now().UnixNano()

	// make home path
	homePath := path.Join(basePath, fmt.Sprintf("container_%d", id))
	if err := os.Mkdir(homePath, 0777); err != nil {
		log.Println("mkdir "+homePath, err)
		_ = os.RemoveAll(basePath)
		return nil
	}

	// make upper path
	upperPath := path.Join(homePath, "/upper")
	if err := os.Mkdir(upperPath, 0777); err != nil {
		log.Println("mkdir "+upperPath, err)
		_ = os.RemoveAll(basePath)
		return nil
	}

	// make worker path
	workerPath := path.Join(homePath, "/worker")
	if err := os.Mkdir(workerPath, 0777); err != nil {
		log.Println("mkdir "+workerPath, err)
		_ = os.RemoveAll(basePath)
		return nil
	}

	// make mount path
	mountPath := path.Join(homePath, "/merge")
	if err := os.Mkdir(mountPath, 0777); err != nil {
		log.Println("mkdir "+mountPath, err)
		_ = os.RemoveAll(basePath)
		return nil
	}

	t2 := time.Now().UnixNano()

	data := fmt.Sprintf("lowerdir=%s,upperdir=%s,workdir=%s", lowerPath, upperPath, workerPath)
	if err := syscall.Mount("overlay", mountPath, "overlay", 0, data); err != nil {
		log.Println("syscall.Mount", err)
		_ = os.RemoveAll(basePath)
		return nil
	}

	t3 := time.Now().UnixNano()
	fmt.Printf("%d, %d, %d, %d\n", id, (t2-t1)/1e3, (t3-t2)/1e3, (t3-t1)/1e3)
	return []int64{(t2 - t1) / 1e3, (t3 - t2) / 1e3, (t3 - t1) / 1e3}
}

func cleanOverlay() {
	for i := 0; i < count; i++ {
		homePath := path.Join(basePath, fmt.Sprintf("container_%d", i))
		mountPath := path.Join(homePath + "/merge")
		if err := syscall.Unmount(mountPath, 0); err != nil {
			log.Println(err)
		}
		if err := os.RemoveAll(homePath); err != nil {
			log.Println(err)
		}
	}
}

func mkdirRootIfNotExist(p string) {
	if f, err := os.Stat(p); err == nil { // 如果已经存在了
		if !f.IsDir() { // 如果不是目录，报错
			log.Fatal(err)
		}
	} else if os.IsNotExist(err) { // 如果不存在
		if err := os.MkdirAll(p, 0744); err != nil {
			log.Fatal(err)
		}
	} else {
		log.Fatal(err)
	}
}

func main() {
	flag.Parse()
	fmt.Printf("count=%d, parallel=%v\n", count, parallel)
	mkdirRootIfNotExist(lowerPath)

	if parallel {
		c := make(chan []int64, count)
		for i := 0; i < count; i++ {
			go func(id int) {
				res := createOverlay(id)
				c <- res
			}(i)
		}
		for i := 0; i < count; i++ {
			<-c
		}
	} else {
		for i := 0; i < count; i++ {
			createOverlay(i)
		}
	}

	cleanOverlay()
}
