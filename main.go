package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"time"

	"launchpad.net/gommap"
)

func main() {
	f, err := ioutil.TempFile("", "mmap-test-")
	if err != nil {
		panic(err)
	}
	defer os.Remove(f.Name())
	defer f.Close()

	fmt.Println(f.Name())

	err = f.Truncate(1024 * 1024 * 1024)
	if err != nil {
		panic(err)
	}

	mmap, err := gommap.Map(f.Fd(), gommap.PROT_WRITE,
		gommap.MAP_SHARED)
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			<-time.Tick(10 * time.Second)
			err := mmap.Sync(gommap.MS_SYNC)
			if err != nil {
				panic(err)
			}
			fmt.Println("flushed")
		}
	}()

	fmt.Println(len(mmap))

	rand.Seed(42)

	for {
		i := rand.Intn(len(mmap))
		mmap[i] = byte(rand.Int())
	}

}
