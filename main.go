package main

import (
	"SpiderGo/spider"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

func main() {

	wg := &sync.WaitGroup{}
	fileLock := &sync.Mutex{}
	file, err := os.Open("search_info.txt")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	// 这个方法体执行完成后，关闭文件
	defer file.Close()

	reader := csv.NewReader(file)

	// Read返回的是一个数组，它已经帮我们分割了，
	record, err := reader.ReadAll()
	// 如果读到文件的结尾，EOF的优先级居然比nil还高！
	if err == io.EOF {
		return
	} else if err != nil {
		fmt.Println("记录集错误:", err)
		return
	}
	for i := 0; i < len(record); i++ {
		wg.Add(1)
		go spider.Spider(record[i][0], wg, fileLock)
		time.Sleep(1 * time.Millisecond)
	}
	wg.Wait()
	fmt.Println("main exit...")
}
