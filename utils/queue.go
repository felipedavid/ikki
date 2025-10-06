package utils

type Queue struct {
	start *Node
	end *Node
	length int
}

type Node struct {
	value any
	next *Node
}

func (q *Queue) Enqueue(value any) {
	node := &Node{value: value}

	if q.length == 0 {
		q.start = node
		q.end = node
	} else {
		q.end.next = node
		q.end = node
	}
	q.length++
}

func (q *Queue) Dequeu() any {
	if q.length == 0 {
		return nil
	}
	firstNode := *q.start
	q.start = firstNode.next
	return firstNode.value
}