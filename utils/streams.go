package utils

type filterFunc[A any] func(A) bool
type mapFunc[A, B any] func(A) B
type reduceFunc[A any] func(A, A) A

func Filter[A any](input []A, f filterFunc[A]) []A {
	var output []A
	for _, element := range input {
		if f(element) {
			output = append(output, element)
		}
	}
	return output
}

func Map[A, B any](input []A, m mapFunc[A, B]) []B {
	output := make([]B, len(input))
	for i, element := range input {
		output[i] = m(element)
	}
	return output
}

func Reduce[A any](input []A, r reduceFunc[A], initial A) A {
	acc := initial
	for _, v := range input {
		acc = r(acc, v)
	}
	return acc
}

type DescendingSort []string

func (d DescendingSort) Len() int           { return len(d) }
func (d DescendingSort) Less(i, j int) bool { return d[i] > d[j] }
func (d DescendingSort) Swap(i, j int)      { d[i], d[j] = d[j], d[i] }
