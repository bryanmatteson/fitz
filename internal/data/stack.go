package data

// Stack ...
type Stack struct {
	data   []interface{}
	length int
}

// NewStack ...
func NewStack(capacity int) *Stack {
	return &Stack{data: make([]interface{}, capacity), length: 0}
}

// Len ...
func (s *Stack) Len() int { return s.length }

// Push ...
func (s *Stack) Push(item interface{}) {
	if s.length == cap(s.data) {
		newData := make([]interface{}, cap(s.data)*2)
		for i := 0; i < s.length; i++ {
			newData[i] = s.data[i]
		}
		s.data = newData
	}
	s.data[s.length] = item
	s.length++
}

// Pop ...
func (s *Stack) Pop() interface{} {
	if s.length == 0 {
		panic("underflow")
	}
	s.length--
	item := s.data[s.length]
	s.data[s.length] = nil
	return item
}

// Peek ...
func (s *Stack) Peek(n int) interface{} {
	if s.length == 0 || n < 0 || n >= s.length {
		panic("range error")
	}
	return s.data[s.length-1-n]
}

// Clear ...
func (s *Stack) Clear() {
	for i := 0; i < s.length; i++ {
		s.data[i] = nil
	}
	s.length = 0
}
