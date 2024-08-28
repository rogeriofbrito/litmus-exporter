package util

func SliceMap[T any, R any](input []T, fn func(T) R) []R {
	output := make([]R, len(input))
	for i, v := range input {
		output[i] = fn(v)
	}
	return output
}
