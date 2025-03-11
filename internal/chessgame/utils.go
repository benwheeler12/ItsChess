package chessgame

func abs(x int) int {
	if x < 0 {
		return x * -1
	}
	return x
}

func contains[T comparable](slice []T, element T) bool {
	for _, item := range slice {
		if item == element {
			return true
		}
	}
	return false
}

func (pc pieceColor) oppositeColor() pieceColor {
	if pc == nocolor {
		panic("Cannot get opposite color of nocolor")
	}
	if pc == white {
		return black
	}
	return white
}

func (v vector2) add(v2 vector2) vector2 {
	return vector2{v.x + v2.x, v.y + v2.y}
}

func filter[T any](slice []T, f func(T) bool) []T {
	filtered := make([]T, 0)
	for _, item := range slice {
		if f(item) {
			filtered = append(filtered, item)
		}
	}
	return filtered
}
