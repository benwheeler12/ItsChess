package itschess

import "image/color"

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

func (v vector2) equals(v2 vector2) bool {
	return v.x == v2.x && v.y == v2.y
}

func (p point) add(p2 point) point {
	return point{p.x + p2.x, p.y + p2.y}
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

func rotateChessCoord(square vector2, rotationFactor int) vector2 {
	// normalize between 0 and 3
	rotationFactor = ((rotationFactor % 4) + 4) % 4

	x, y := 7-square.x, square.y

	switch rotationFactor {
	case 0:
		return vector2{x, y}
	case 1:
		return vector2{7 - y, x}
	case 2:
		return vector2{7 - x, 7 - y}
	case 3:
		return vector2{y, 7 - x}
	}
	return vector2{x, y} // default case, shouldn't be reached
}

func tintColor(originalColor color.RGBA, tint color.RGBA, alpha float64) color.RGBA {
	// Ensure alpha is in range [0,1]
	if alpha < 0 || alpha > 1 {
		panic("alpha value not between 0 and 1")
	}

	return color.RGBA{
		uint8((float64(originalColor.R) * (1 - alpha)) + (float64(tint.R) * alpha)),
		uint8((float64(originalColor.G) * (1 - alpha)) + (float64(tint.G) * alpha)),
		uint8((float64(originalColor.B) * (1 - alpha)) + (float64(tint.B) * alpha)),
		originalColor.A,
	}

}
