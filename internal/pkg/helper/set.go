package helper

type Set map[interface{}]struct{}

func (s *Set) Add(item interface{}) {
	(*s)[item] = struct{}{}
}

func (s *Set) Contains(value interface{}) bool {
	_, exists := (*s)[value]
	return exists
}

func (s *Set) Remove(value interface{}) {
	delete(*s, value)
}

func (s *Set) Size() int {
	return len(*s)
}

func NewSet() Set {
	return make(Set)
}
