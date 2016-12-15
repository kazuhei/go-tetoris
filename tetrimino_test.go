package main

import (
	"testing"
)

type pointTestCase struct {
	point  point
	expect point
}

type directionTestCase struct {
	direction Direction
	expect    point
}

func TestRotatePoint(t *testing.T) {
	testCases := []pointTestCase{
		{point{X: 2, Y: 1}, point{X: 1, Y: -2}},
		{point{X: 1, Y: -2}, point{X: -2, Y: -1}},
		{point{X: -2, Y: -1}, point{X: -1, Y: 2}},
		{point{X: -1, Y: 2}, point{X: 2, Y: 1}},
	}
	for i, testCase := range testCases {
		actual := rotatePoint(testCase.point)
		if actual != testCase.expect {
			t.Errorf("#%d: got:%v expect:%v", i, actual, testCase.expect)
		}
	}
}

func TestGetNextPoint(t *testing.T) {
	prevPoint := point{X: 3, Y: 3}
	testCases := []directionTestCase{
		{LEFT, point{X: 2, Y: 3}},
		{DOWN, point{X: 3, Y: 4}},
		{RIGHT, point{X: 4, Y: 3}},
	}
	for i, testCase := range testCases {
		tetorimino := newTetorimino(prevPoint)
		actual := tetorimino.getNextPoint(testCase.direction)
		if actual != testCase.expect {
			t.Errorf("#%d: got:%v expect:%v", i, actual, testCase.expect)
		}
	}
}
