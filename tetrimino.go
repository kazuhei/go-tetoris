package main

import "math/rand"

// Shape Tetoriminoの形の種類
type Shape int

const (
	IShape = iota
	OShape
	SShape
	ZShape
	JShape
	LShape
	TShape
)

// Tetrimino ランダム生成されるブロックのインターフェース
type Tetrimino interface {
	Move(d Direction) []point
	GetCurrentPoints() []point
	GetNextPoints(d Direction) []point
}

type tetorimino struct {
	rotateCount int
	shape       Shape
	point       point
}

func newTetorimino(p point) *tetorimino {
	shapes := []Shape{IShape, OShape, SShape, ZShape, JShape, LShape, TShape}
	shape := shapes[rand.Intn(len(shapes))]
	return &tetorimino{
		rotateCount: 0,
		shape:       shape,
		point:       p,
	}
}

func (t *tetorimino) Move(d Direction) []point {
	t.point = t.getNextPoint(d)
	if d == ROTATE {
		t.rotateCount++
		return t.getPoints(t.point, t.rotateCount)
	}
	return t.getPoints(t.point, t.rotateCount)
}

func (t *tetorimino) GetCurrentPoints() []point {
	return t.getPoints(t.point, t.rotateCount)
}

func (t *tetorimino) GetNextPoints(d Direction) []point {
	point := t.getNextPoint(d)
	if d == ROTATE {
		return t.getPoints(point, t.rotateCount+1)
	}
	return t.getPoints(point, t.rotateCount)
}

func (t *tetorimino) getNextPoint(d Direction) point {
	newPoint := point{X: t.point.X, Y: t.point.Y}
	switch d {
	case DOWN:
		newPoint.Y = newPoint.Y + 1
	case LEFT:
		newPoint.X = newPoint.X - 1
	case RIGHT:
		newPoint.X = newPoint.X + 1
	case ROTATE:
		//
	}
	return newPoint
}

func (t *tetorimino) getPoints(p point, rotateCount int) []point {
	originpoints := getOriginPoints(t.shape)
	points := []point{}
	for _, origin := range originpoints {
		rotatedPoint := origin
		mod4 := rotateCount % 4
		for i := 0; i < mod4; i++ {
			rotatedPoint = rotatePoint(rotatedPoint)
		}
		points = append(points, point{X: p.X + rotatedPoint.X, Y: p.Y + rotatedPoint.Y})
	}
	return points
}

func rotatePoint(p point) point {
	return point{X: p.Y, Y: -p.X}
}

func getOriginPoints(s Shape) []point {
	switch s {
	case IShape:
		return []point{
			{X: 0, Y: 0},
			{X: 0, Y: 1},
			{X: 0, Y: 2},
			{X: 0, Y: 3},
		}
	case OShape:
		return []point{
			{X: 1, Y: 0},
			{X: 0, Y: 0},
			{X: 1, Y: 1},
			{X: 0, Y: 1},
		}
	case SShape:
		return []point{
			{X: -1, Y: 1},
			{X: 0, Y: 1},
			{X: 0, Y: 0},
			{X: 1, Y: 0},
		}
	case ZShape:
		return []point{
			{X: -1, Y: 0},
			{X: 0, Y: 1},
			{X: 0, Y: 0},
			{X: 1, Y: 1},
		}
	case JShape:
		return []point{
			{X: -1, Y: 0},
			{X: -1, Y: 1},
			{X: 0, Y: 1},
			{X: 1, Y: 1},
		}
	case LShape:
		return []point{
			{X: -1, Y: 1},
			{X: 0, Y: 1},
			{X: 1, Y: 1},
			{X: 1, Y: 0},
		}
	case TShape:
		return []point{
			{X: -1, Y: 1},
			{X: 0, Y: 1},
			{X: 1, Y: 1},
			{X: 0, Y: 0},
		}
	}
	return []point{
		{X: 0, Y: 0},
		{X: 0, Y: 1},
		{X: 0, Y: 2},
		{X: 0, Y: 3},
	}
}
