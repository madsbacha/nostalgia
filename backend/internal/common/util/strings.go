package util

func Contains[T comparable](s []T, str T) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func Remove[T comparable](s []T, str T) []T {
	res := make([]T, 0)
	for _, v := range s {
		if v != str {
			res = append(res, v)
		}
	}
	return res
}

func Intersection[T comparable](a []T, b []T) []T {
	res := make([]T, 0)
	for _, v := range a {
		if Contains(b, v) {
			res = append(res, v)
		}
	}
	return res
}

func Unique[T comparable](a []T) []T {
	res := make([]T, 0)
	seen := make(map[T]bool)
	for _, v := range a {
		if !seen[v] {
			seen[v] = true
			res = append(res, v)
		}
	}
	return res
}

func Union[T comparable](a []T, b []T) []T {
	res := make([]T, 0)
	seen := make(map[T]bool)
	for _, v := range a {
		if !seen[v] {
			seen[v] = true
			res = append(res, v)
		}
	}
	for _, v := range b {
		if !seen[v] {
			seen[v] = true
			res = append(res, v)
		}
	}
	return res
}

func Except[T comparable](list []T, except []T) []T {
	res := make([]T, 0)
	seen := make(map[T]bool)
	for _, v := range except {
		seen[v] = true
	}
	for _, v := range list {
		if !seen[v] {
			res = append(res, v)
		}
	}
	return res
}
