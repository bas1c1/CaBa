package main

import (
	"sync"
	"sync/atomic"
)

type queueItem struct {
	ID     int64
	tr     transaction
	result chan string
}

type queue struct {
	items    []*queueItem
	mutex    sync.Mutex
	maxSize  int64
	addChan  chan *queueItem
	stopChan chan struct{}
	wg       sync.WaitGroup
	nextID   int64
}

func newQueue(maxSize int64) *queue {
	if maxSize <= 0 {
		maxSize = 8192
	}
	queue := &queue{
		items:    make([]*queueItem, 0, maxSize),
		maxSize:  maxSize,
		addChan:  make(chan *queueItem),
		stopChan: make(chan struct{}),
		nextID:   0,
	}
	queue.wg.Add(1)
	go queue.processItems()
	return queue
}

func (q *queue) Add(tr transaction) (int64, <-chan string, error) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	if int64(len(q.items)) >= q.maxSize {
		removedItem := q.items[0]
		q.items = q.items[1:]
		close(removedItem.result)
	}

	id := atomic.AddInt64(&q.nextID, 1)
	tr.id = id
	resultChan := make(chan string, 1)
	item := &queueItem{
		ID:     id,
		tr:     tr,
		result: resultChan,
	}

	q.items = append(q.items, item)
	q.addChan <- item

	return id, resultChan, nil
}

func (q *queue) Remove(id int64) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	for i, item := range q.items {
		if item.ID == id {
			q.items = append(q.items[:i], q.items[i+1:]...)
			break
		}
	}
}

func (q *queue) processItems() {
	defer q.wg.Done()
	for {
		select {
		case item := <-q.addChan:
			result := item.tr.execute()
			item.result <- result
			close(item.result)
			q.mutex.Lock()
			for i, qItem := range q.items {
				if qItem.ID == item.ID {
					q.items = append(q.items[:i], q.items[i+1:]...)
					break
				}
			}
			q.mutex.Unlock()
		case <-q.stopChan:
			return
		}
	}
}

func (q *queue) Stop() {
	close(q.stopChan)
	q.wg.Wait()
	q.mutex.Lock()
	defer q.mutex.Unlock()
	for _, item := range q.items {
		close(item.result)
	}
}
