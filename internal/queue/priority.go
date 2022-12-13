package queue

import "sort"

type PriorityQueue[T comparable] struct {
	priorities map[T]int
	queue      []T
}

func NewPriorityQueue[T comparable]() *PriorityQueue[T] {
	return &PriorityQueue[T]{
		priorities: map[T]int{},
	}
}

func (queue *PriorityQueue[T]) Get() T {
	var item T
	if len(queue.queue) < 1 {
		return item
	}

	queue.sort()
	item = queue.queue[0]

	// remove item from queue
	delete(queue.priorities, item)
	queue.queue = queue.queue[1:]

	return item
}

func (queue *PriorityQueue[T]) Put(item T, priority int) {
	if _, exists := queue.priorities[item]; !exists {
		queue.queue = append(queue.queue, item)
	}
	queue.priorities[item] = priority
}

func (queue *PriorityQueue[T]) Size() int {
	return len(queue.queue)
}

func (queue *PriorityQueue[T]) sort() {
	// this is a bit crap given we only ever need to move one item, but it works!
	sort.Slice(queue.queue, func(i, j int) bool {
		return queue.priorities[queue.queue[i]] < queue.priorities[queue.queue[j]]
	})
}
