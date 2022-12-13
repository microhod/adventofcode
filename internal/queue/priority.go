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

func (queue *PriorityQueue[T]) ExtractMin() T {
	var item T
	if len(queue.queue) < 1 {
		return item
	}

	queue.sort()

	item = queue.queue[0]
	queue.queue = queue.queue[1:]
	delete(queue.priorities, item)

	return item
}

func (queue *PriorityQueue[T]) Empty() bool {
	return len(queue.queue) == 0
}

func (queue *PriorityQueue[T]) AddWithPriority(item T, priority int) {
	queue.queue = append(queue.queue, item)
	queue.priorities[item] = priority
}

func (queue *PriorityQueue[T]) DecreasePriority(item T, priority int) {
	queue.priorities[item] = priority
}

func (queue *PriorityQueue[T]) sort() {
	// this is a bit crap given we only ever need to move one item, but it works!
	sort.Slice(queue.queue, func(i, j int) bool {
		return queue.priorities[queue.queue[i]] < queue.priorities[queue.queue[j]]
	})
}
