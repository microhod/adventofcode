package queue

import (
	"container/heap"
)

type PriorityQueue[T any] struct {
	q           *priorityQueue[T]
	LowestFirst bool
}

type PriorityQueueOption[T any] func(*PriorityQueue[T])

func NewPriorityQueue[T any](options ...PriorityQueueOption[T]) PriorityQueue[T] {
	queue := PriorityQueue[T]{
		LowestFirst: true,
	}
	for _, o := range options {
		o(&queue)
	}

	queue.q = &priorityQueue[T]{
		q:           make([]*priorityQueueItem[T], 0),
		lowestFirst: queue.LowestFirst,
	}
	heap.Init(queue.q)
	return queue
}

func (q PriorityQueue[T]) Push(item T, priority int) {
	heap.Push(q.q, &priorityQueueItem[T]{
		value:    item,
		priority: priority,
	})
}

func (q PriorityQueue[T]) Pop() T {
	item := heap.Pop(q.q).(*priorityQueueItem[T])
	return item.value
}

func (q PriorityQueue[T]) Size() int {
	return q.q.Len()
}

// An priorityQueueItem is something we manage in a priority queue.
type priorityQueueItem[T any] struct {
	value    T
	priority int
	index    int
}

// A priorityQueue implements heap.Interface and holds Items.
type priorityQueue[T any] struct {
	q           []*priorityQueueItem[T]
	lowestFirst bool
}

func (pq priorityQueue[T]) Len() int { return len(pq.q) }

func (pq priorityQueue[T]) Less(i, j int) bool {
	if pq.lowestFirst {
		return pq.q[i].priority < pq.q[j].priority
	}
	return pq.q[i].priority > pq.q[j].priority
}

func (pq priorityQueue[T]) Swap(i, j int) {
	pq.q[i], pq.q[j] = pq.q[j], pq.q[i]
	pq.q[i].index = i
	pq.q[j].index = j
}

func (pq *priorityQueue[T]) Push(x any) {
	n := len(pq.q)
	item := x.(*priorityQueueItem[T])
	item.index = n
	pq.q = append(pq.q, item)
}

func (pq *priorityQueue[T]) Pop() any {
	old := pq.q
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	pq.q = old[0 : n-1]
	return item
}
