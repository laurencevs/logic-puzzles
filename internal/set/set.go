package set

type Set[T comparable] map[T]struct{}

func New[T comparable](values ...T) Set[T] {
	s := make(Set[T], len(values))
	for _, v := range values {
		s.Add(v)
	}
	return s
}

func (s Set[T]) Add(value T) {
	s[value] = struct{}{}
}

func (s Set[T]) Contains(value T) bool {
	_, ok := s[value]
	return ok
}
