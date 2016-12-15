package main

import (
	"reflect"
	"testing"
)

type blocksTestCase struct {
	blocks []bool
	expect bool
}

type pointsTestCase struct {
	points []point
	expect bool
}

type transformPointsTestCase struct {
	points []point
	expect []point
}

func TestIsInvalidPoints(t *testing.T) {
	var stage = newStage(10, 10)
	testCases := []pointsTestCase{
		{[]point{point{X: 0, Y: 1}}, false}, // 正常系
		{[]point{point{X: -1, Y: 1}}, true}, // -x方向にstageからはみ出ている
		{[]point{point{X: 10, Y: 1}}, true}, // x方向にstageからはみ出ている
	}
	for i, testCase := range testCases {
		actual := stage.isInvalidPoints(testCase.points)
		if actual != testCase.expect {
			t.Errorf("#%d: got:%v expect:%v", i, actual, testCase.expect)
		}
	}
}

func TestIsCollision(t *testing.T) {
	var stage = newStage(10, 10)
	stage.blocks = []point{point{X: 8, Y: 9}}
	testCases := []pointsTestCase{
		{[]point{point{X: 0, Y: 0}}, false}, // 正常系
		{[]point{point{X: 0, Y: 10}}, true}, // 最下部に衝突
		{[]point{point{X: 8, Y: 9}}, true},  // 他のブロックに衝突
	}
	for i, testCase := range testCases {
		actual := stage.isCollision(testCase.points)
		if actual != testCase.expect {
			t.Errorf("#%d: got:%v expect:%v", i, actual, testCase.expect)
		}
	}
}

func TestAddBlock(t *testing.T) {
	var stage = newStage(10, 10)
	p := point{X: 1, Y: 1}
	testCase := transformPointsTestCase{[]point{p}, []point{p}}
	stage.addBlocks(testCase.points)
	if !reflect.DeepEqual(stage.blocks, testCase.expect) {
		t.Errorf("got:%v expect:%v", stage.blocks, testCase.expect)
	}
}

func TestIsFilled(t *testing.T) {
	testCases := []blocksTestCase{
		{[]bool{false, false, false}, false},
		{[]bool{true, false, false}, false},
		{[]bool{true, true, true}, true},
	}
	for i, testCase := range testCases {
		actual := isFilled(testCase.blocks)
		if actual != testCase.expect {
			t.Errorf("#%d: got:%v expect:%v", i, actual, testCase.expect)
		}
	}
}

func TestRemoveFilledRows(t *testing.T) {
	var stage = newStage(3, 3)
	testCases := []transformPointsTestCase{
		{ // 列を埋めていないときは何も起きない
			[]point{{X: 0, Y: 0}},
			[]point{{X: 0, Y: 0}},
		},
		{ // 列を埋めていて、さらに上に載っているときはそれが落ちてくる
			[]point{{X: 0, Y: 0}, {X: 0, Y: 1}, {X: 1, Y: 1}, {X: 2, Y: 1}},
			[]point{{X: 0, Y: 1}},
		},
		{ // 2列を埋めているときは2段消えて、上のブロックは2段ずれる
			[]point{{X: 0, Y: 0}, {X: 0, Y: 1}, {X: 1, Y: 1}, {X: 2, Y: 1}, {X: 0, Y: 2}, {X: 1, Y: 2}, {X: 2, Y: 2}},
			[]point{{X: 0, Y: 2}},
		},
	}
	for i, testCase := range testCases {
		actual := stage.removeFilledRows(testCase.points, 1)
		if !reflect.DeepEqual(testCase.expect, actual) {
			t.Errorf("#%d: got:%v expect:%v", i, actual, testCase.expect)
		}
	}
}
