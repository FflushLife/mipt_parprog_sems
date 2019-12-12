package main

import (
	"fmt"
	"os"
	"pool"
	"psort"
	"runtime"
	"strconv"
	"time"
	"unsafe"
)

func main() {
	var sortInstance *psort.PSort
	var poolInstance *pool.Pool
	var tCount uint64 = 1

	if (len(os.Args) > 1) {
		tCount, _ = strconv.ParseUint(os.Args[1], 10, 64)
	}
	runtime.GOMAXPROCS(int(tCount))

	fmt.Println("Start initializing...")
	target := []int{3, 4, 1, 11, 8, 7, 2, 6, 5, 0, 9}
	for i := 30000000; i > 12; i-- {
		target = append(target, i)
	}
	sortInstance = psort.New(target, tCount)
	// TODO:: try to create with nil
	poolInstance = pool.New(tCount, psort.TSort, unsafe.Pointer(sortInstance))

	start := time.Now()
	poolInstance.Start();
	fmt.Println("time=", time.Now().Sub(start))

	result := make([]int, 0)
	var sliceThreadSize float64 = float64(len(sortInstance.GetTarget())) / float64(tCount)
	for i := uint64(0); i < tCount; i++ {
		l := int(float64(i) * sliceThreadSize)
		r := int(float64(i + 1) * sliceThreadSize)
		result = psort.Merge(result, sortInstance.GetTarget()[l:r])
	}
	fmt.Println("final time=", time.Now().Sub(start))
	fmt.Println(result[:100])

	result = make([]int, 0)
	sortInstance = psort.New([]int{3,2,1,6,5,7,8,14,36,7645,12,65,63}, tCount)
	poolInstance.Lock()
	poolInstance.ChangeTask(unsafe.Pointer(sortInstance))
	fmt.Println("new task started")
	poolInstance.Start()
	//poolInstance.Unlock()
	sliceThreadSize = float64(len(sortInstance.GetTarget())) / float64(tCount)
	for i := uint64(0); i < tCount; i++ {
		l := int(float64(i) * sliceThreadSize)
		r := int(float64(i + 1) * sliceThreadSize)
		result = psort.Merge(result, sortInstance.GetTarget()[l:r])
	}
	fmt.Println("final time=", time.Now().Sub(start))
	fmt.Println(result)
	fmt.Println(target[:100])
}
