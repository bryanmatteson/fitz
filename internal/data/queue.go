package data

// Queue ...
type Queue struct {
	data []interface{}
}

// NewQueueWithSize ...
func NewQueueWithSize(capacity int) *Queue {
	return &Queue{data: make([]interface{}, 0, capacity)}
}

// NewQueue ...
func NewQueue(data []interface{}) *Queue {
	return &Queue{data: data}
}

// Len ...
func (s *Queue) Len() int { return len(s.data) }

// Enqueue ...
func (s *Queue) Enqueue(item interface{}) {
	s.data = append(s.data, item)
}

// Dequeue ...
func (s *Queue) Dequeue() (result interface{}) {
	if len(s.data) == 0 {
		return nil
	}

	result = s.data[0]
	s.data = s.data[1:]
	return
}
