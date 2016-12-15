package main

type point struct {
	X int
	Y int
}

func (p point) equal(np point) bool {
	return p.X == np.X && p.Y == np.Y
}
