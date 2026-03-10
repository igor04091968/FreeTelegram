package queue

import "errors"

type Options struct {
	MaxSize int
}

type Queue struct {
	maxSize int
	ch      chan any
}

func New(opts Options) *Queue {
	max := opts.MaxSize
	if max <= 0 {
		max = 100
	}
	return &Queue{maxSize: max, ch: make(chan any, max)}
}

func (q *Queue) Enqueue(item any) error {
	select {
	case q.ch <- item:
		return nil
	default:
		return errors.New("queue full")
	}
}

func (q *Queue) Dequeue() (any, bool) {
	item, ok := <-q.ch
	return item, ok
}

func (q *Queue) Close() {
	close(q.ch)
}
