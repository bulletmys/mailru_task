package analysis

import (
	"bytes"
	"fmt"
	"io"
	"sync"
)

func initWrapperParams(word string) *wrapperParams {
	return &wrapperParams{
		mu:       &sync.RWMutex{},
		wg:       &sync.WaitGroup{},
		word:     word,
		goAmount: 0,
		counts:   make([]PathValue, 0),
		queue:    make([]PathReader, 0),
	}
}

//Планировщик обработчиков. Запускает поступившую задачу (через канал DataChan) в горутине,
//если не превышен лимит. Иначе - кладет в очередь. А также ожидает
//команду на завершение работы (через канал Quit)
func (wc *WordCounterParams) WordCounter() {
	p := initWrapperParams(wc.Word)

	for {
		select {
		case reader := <-wc.DataChan:
			//fmt.Println(runtime.NumGoroutine())
			p.mu.Lock()

			if p.goAmount < wc.GoMax {
				p.wg.Add(1)
				go p.counterWrapper(reader)
				p.goAmount++
				p.mu.Unlock()
				continue
			}
			p.queue = append(p.queue, reader)

			p.mu.Unlock()
		case <-wc.Quit:
			p.wg.Wait()
			wc.ResChan <- p.counts
			return
		}
	}
}

//Запускает функцию по подсчету слов в строке, записывает ее результат
//и берет новую задачу из очереди, если она не пустая
func (wp *wrapperParams) counterWrapper(reader PathReader) {
	val, err := counter(*reader.Reader, wp.word)

	wp.mu.Lock()

	if err != nil {
		wp.counts = append(wp.counts, PathValue{Path: reader.Path, Err: err})
	} else {
		wp.counts = append(wp.counts, PathValue{Path: reader.Path, Value: val})
	}

	if len(wp.queue) > 0 {
		reader := wp.queue[0]
		wp.queue = wp.queue[1:]

		wp.mu.Unlock()

		wp.counterWrapper(reader)
		return
	}
	wp.goAmount--

	wp.mu.Unlock()
	wp.wg.Done()
}

//Функция по подсчету вхождений подстроки word в reader
func counter(reader io.ReadCloser, word string) (int64, error) {
	defer reader.Close()

	var count int64

	buf := make([]byte, 4*1024)
	for {
		n, err := reader.Read(buf)
		if err != nil || n == 0 {
			if err == io.EOF {
				break
			}
			return 0, fmt.Errorf("Error while reading from source: %v\n", err)
		}
		count += int64(bytes.Count(buf, []byte(word)))
	}
	return count, nil
}
