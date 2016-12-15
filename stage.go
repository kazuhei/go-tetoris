package main

import (
	"sync"
	"time"

	"strconv"

	termbox "github.com/nsf/termbox-go"
)

var mu sync.Mutex

// Direction 移動方向
type Direction int

const (
	// DOWN 下方向
	DOWN = iota
	// LEFT 左方向
	LEFT
	// RIGHT 右方向
	RIGHT
	// ROTATE は右回りに90度回転する
	ROTATE
)

// Stage はゲームのステージです
type Stage interface {
	Start()
}

type stage struct {
	stageWidth        int
	stageHeight       int
	currentTetorimino Tetrimino
	blocks            []point
	score             int
	finished          bool
}

var initialPoint = point{X: 5, Y: 0}

func newStage(width, height int) *stage {
	var tetorimino Tetrimino
	tetorimino = newTetorimino(initialPoint)
	return &stage{
		stageWidth:        width,
		stageHeight:       height,
		currentTetorimino: tetorimino,
		blocks:            []point{},
		score:             0,
		finished:          false,
	}
}

func (s *stage) Start() {
	keyCh := make(chan termbox.Key)
	timerCh := make(chan bool)

	go s.observeUserKeyInput(keyCh)
	go s.observeTimer(timerCh)

	for {
		select {
		case key := <-keyCh: //キーイベント
			mu.Lock()
			switch key {
			case termbox.KeyEsc:
				return
			case termbox.KeyArrowDown:
				s.update(DOWN)
				s.draw()
			case termbox.KeyArrowLeft:
				s.update(LEFT)
				s.draw()
			case termbox.KeyArrowRight:
				s.update(RIGHT)
				s.draw()
			case termbox.KeySpace:
				s.update(ROTATE)
				s.draw()
			}
			mu.Unlock()
		case <-timerCh: //タイマーイベント
			mu.Lock()
			s.update(DOWN)
			s.draw()
			mu.Unlock()
		default:
			break
		}
	}
}

func (s *stage) observeUserKeyInput(kch chan termbox.Key) {
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			kch <- ev.Key
		}
	}
}

func (s *stage) observeTimer(tch chan bool) {
	for {
		tch <- true
		time.Sleep(time.Duration(500) * time.Millisecond)
	}
}

func (s *stage) update(d Direction) {
	currentPoints := s.currentTetorimino.GetCurrentPoints()
	nextPoints := s.currentTetorimino.GetNextPoints(d)
	if s.isInvalidPoints(nextPoints) {
		return // ステージからはみ出ないように
	} else if s.isCollision(nextPoints) {
		s.addBlocks(currentPoints) // 動かす前のtetoriminoをblocksに移動
		s.blocks = s.removeFilledRows(s.blocks, 1)
		if s.judgeFinished(initialPoint) {
			s.finished = true
			return
		}
		s.currentTetorimino = newTetorimino(initialPoint)
	} else {
		s.currentTetorimino.Move(d)
	}
}

func (s *stage) isInvalidPoints(ps []point) bool {
	for _, point := range ps {
		if point.X < 0 || point.X > s.stageWidth-1 {
			return true
		}
	}
	return false
}

func (s *stage) isCollision(ps []point) bool {
	for _, nextPoint := range ps {
		// 床に衝突した場合はoldPointsで固定する
		if nextPoint.Y >= s.stageHeight {
			return true
		}
		// 既存のblockと重なる場合は衝突
		for _, blockPoint := range s.blocks {
			if nextPoint == blockPoint {
				return true
			}
		}
	}
	return false
}

func (s *stage) judgeFinished(p point) bool {
	for _, point := range s.blocks {
		if point.equal(initialPoint) {
			return true
		}
	}
	return false
}

func (s *stage) addBlocks(points []point) {
	for _, point := range points {
		s.blocks = append(s.blocks, point)
	}
}

func (s *stage) removeFilledRows(ps []point, deletedCount int) []point {
	var blocks [][]bool
	// 座標文だけ配列を用意する
	for y := 0; y < s.stageHeight; y++ {
		var row []bool
		for x := 0; x < s.stageWidth; x++ {
			row = append(row, false)
		}
		blocks = append(blocks, row)
	}

	// 既にblockがあるところをtrueにする
	for _, point := range ps {
		blocks[point.Y][point.X] = true
	}

	// 横一列が全てtrueの場合はfalseにする
	removeIndex := -1
	for y, row := range blocks {
		if isFilled(row) {
			removeIndex = y
			for x := range row {
				blocks[y][x] = false
			}
			break
		}
	}

	if removeIndex != -1 {
		s.score = s.score + 10*deletedCount
		// 列を消した場合はそれより上の列を詰める
		for y := removeIndex - 1; y >= 0; y-- {
			for x, block := range blocks[y] {
				if block {
					blocks[y][x] = false
					blocks[y+1][x] = true // 1行消しているのでoutOfIndexしないはず
				}
			}
		}

		// 消されなかったblocks情報からpointsに戻す
		var newPoints []point
		for y, row := range blocks {
			for x, block := range row {
				if block {
					newPoints = append(newPoints, point{X: x, Y: y})
				}
			}
		}
		return s.removeFilledRows(newPoints, deletedCount+1)
	}

	return ps
}

func isFilled(blocks []bool) bool {
	for _, block := range blocks {
		if !block {
			return false
		}
	}
	return true
}

func (s *stage) draw() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	if s.finished {
		drawLine(0, s.stageHeight/2, "       GAME OVER")
		termbox.Flush()
		return
	}
	// 表示座標を初期化
	var displayBlocks [][]bool
	for i := 0; i < s.stageHeight; i++ {
		var row []bool
		for j := 0; j < s.stageWidth; j++ {
			row = append(row, false)
		}
		displayBlocks = append(displayBlocks, row)
	}

	// 落ちてくるブロックの表示
	for _, point := range s.currentTetorimino.GetCurrentPoints() {
		displayBlocks[point.Y][point.X] = true
	}

	// 積もったブロックを表示
	for _, point := range s.blocks {
		displayBlocks[point.Y][point.X] = true
	}

	for y := range displayBlocks {
		drawWall(0, y)
		for x := range displayBlocks[y] {
			if displayBlocks[y][x] {
				drawBlock(x+1, y)
			}
		}
		drawWall(s.stageWidth+1, y)
	}
	drawLine(0, s.stageHeight, "------------")
	drawLine(0, s.stageHeight+1, "score:"+strconv.Itoa(s.score))

	termbox.Flush()
}

func drawWall(x, y int) {
	termbox.SetCell(x, y, rune('|'), termbox.ColorDefault, termbox.ColorDefault)
}

func drawBlock(x, y int) {
	termbox.SetCell(x, y, rune('■'), termbox.ColorDefault, termbox.ColorDefault)
}

func drawLine(x, y int, str string) {
	runes := []rune(str)
	for i := 0; i < len(runes); i++ {
		termbox.SetCell(x+i, y, runes[i], termbox.ColorDefault, termbox.ColorDefault)
	}
}
