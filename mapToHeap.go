package main

type Node struct {
	Character string
	Frequency int
	Left      *Node
	Right     *Node
}

type PriorityQueue []*Node

func (pq PriorityQueue) Len() int {
	return len(pq)
}

func (pq PriorityQueue) Less(i int, j int) bool {
	return pq[i].Frequency < pq[j].Frequency
}

func (pq PriorityQueue) Swap(i int, j int) {
	temp := pq[i]
	pq[i] = pq[j]
	pq[j] = temp
}

func (pq *PriorityQueue) Push(x interface{}) {
	item := x.(*Node)
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}

func mapToHeap(letters map[string]int) PriorityQueue {
	pq := make(PriorityQueue, len(letters))
	i := 0

	for char, freq := range letters {
		pq[i] = &Node{Character: char, Frequency: freq}
		i += 1
	}

	return pq
}
