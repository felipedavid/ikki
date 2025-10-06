package utils

type Queue struct {
	nodes []any
	start int
}


func (q *Queue) Enqueue(value any) {
	q.nodes = append(q.nodes, value)
}

func (q *Queue) Dequeu() any {
	if q.start >= len(q.nodes)/2 {
		q.nodes = q.nodes[q.start:]
		q.start = 0
	}

	val := q.nodes[q.start]
	q.start++
	return val
}

func (q *Queue) Len() int {
	return len(q.nodes) - q.start
}