package trilib

type Tri[T any] struct {
	Previous T
	Current  T
	Next     T
}

func Make[T any](list []T) []Tri[T] {
	result := []Tri[T]{}
	count := len(list)
	for i, current := range list {
		var previous T
		previousIndex := i - 1
		if 0 <= previousIndex {
			previous = list[previousIndex]
		}

		var next T
		nextIndex := i + 1
		if nextIndex < count {
			next = list[nextIndex]
		}

		result = append(result, Tri[T]{
			Previous: previous,
			Current:  current,
			Next:     next,
		})
	}
	return result
}
