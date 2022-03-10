package utils

import (
	"log"
	"sync"
)

func (f fileInfo) Run(thread int) error {
	err := f.getHeader()
	if err != nil {
		return err
	}
	jobs := make([]filePart, thread)
	eachSize := f.Size / thread
	for i := range jobs {
		jobs[i].Index = i
		if i == 0 {
			jobs[i].Start = 0
		} else {
			jobs[i].Start = jobs[i-1].End + 1
		}
		if i < thread - 1 {
			jobs[i].End = jobs[i].Start + eachSize
		} else {
			jobs[i].End = f.Size - 1
		}
	}

	var wg sync.WaitGroup
	for _, job := range jobs {
		wg.Add(1)
		go func(job filePart) {
			defer wg.Done()
			err := f.downloadPart(job)
			if err != nil {
				log.Println("下载分片失败:", err, job)
			}
		}(job)
	}
	wg.Wait()
	return f.mergeFileParts()
}
