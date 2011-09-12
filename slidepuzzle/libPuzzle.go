package libPuzzle

import (
  // "fmt"
  "bufio"
  "os"
  "strings"
  "strconv"
  "math"
  "sort"
)

const (
  Wall = 100
  Blank = 0
  
  NoError = 0
  CantMove = 1000
)

const (
  LeftTop = iota
  RightBottom
  LeftBottom
  RightTop
)

type Board struct {
  Panels []byte
  Width int
  Height int
}

type Operation struct {
  Board Board
  History string
  Score float64
}

type Operations []Operation

func (p Operations) Len() int {
  return len(p)
}

func (p Operations) Less(i,j int) bool {
  if p[i].Score < 0 {
    p[i].Score = CalcDistance(p[i].Board, p[i].Board.Width*2)
  }
  
  if p[j].Score < 0 {
    p[j].Score = CalcDistance(p[j].Board, p[j].Board.Width*2)
  }
  
  return p[i].Score < p[j].Score
}

func (p Operations) Swap(i,j int) {
  tmp := p[i]
  p[i] = p[j]
  p[j] = tmp
}

func ParseProblem(r *bufio.Reader) (int, int, []byte) {
  var width,height int
  
  line, error := r.ReadString('\n')
  if error == nil || error == os.EOF {
    line = strings.TrimRight(line, "\n")
  }
  
  if line == "SKIP" {
    return 100, 100, make([]byte, 1)
  }
  // terms := strings.Split(line, ",") // repository
  terms := strings.Split(line, ",", 3) // macports
  
  width, error = strconv.Atoi(terms[0])
  if error != nil {
    width = 0
  }
  
  height, error = strconv.Atoi(terms[1])
  if error != nil {
    width = 0
  }
  
  board := make([]byte, width * height)
  
  var panel byte
  for i:=0;i<width*height;i++ {
    panel = terms[2][i]
    
    if panel >= '0' && panel <= '9' {
      panel -= '0'
    } else if panel >= 'A' && panel <= 'Z' {
      panel = panel - 'A' + 10
    } else {
      panel = Wall
    }
    
    board[i] = panel
  }
  
  return width, height, board
}

func CheckBoard(b Board) bool {
  i:=0
  for ;i<b.Width*b.Height-1;i++ {
    if b.Panels[i] != uint8(i+1) && b.Panels[i] != Wall {
      return false
    }
  }
  if b.Panels[i] != 0 {
    return false
  }
  
  return true
}

func Tick(w int, h int, operations Operations, prevScore float64) Operations {
  result := make(Operations, 0, len(operations) * 4)
  opChan := make(chan Operation)
  compChan := make(chan int)
  
  for i:=0;i<len(operations);i++ {
    countChan := make(chan bool)
    op := operations[i]
    
    go func() {
      count := 0
      for {
        <-countChan
        count++
        
        if count == 4 {
          compChan <- i
          break
        }
      }
    }()
    
    go func() {
      nextBoard, error := MoveUp(w, h, op.Board.Panels)
      if error == NoError {
        if isEffective(w, h, op, nextBoard) {
          opChan <- Operation{Board: Board{Panels: nextBoard, Width: w, Height: h}, History: op.History + "D", Score:-1}
        }
      }
      countChan <- true
    }()
    
    go func() {
      nextBoard, error := MoveDown(w, h, op.Board.Panels)
      if error == NoError {
        if isEffective(w, h, op, nextBoard) {
          opChan <- Operation{Board: Board{Panels: nextBoard, Width: w, Height: h}, History: op.History + "U", Score:-1}
        }
      }
      countChan <- true
    }()
    
    go func() {
      nextBoard, error := MoveLeft(w, h, op.Board.Panels)
      if error == NoError {
        if isEffective(w, h, op, nextBoard) {
          opChan <- Operation{Board: Board{Panels: nextBoard, Width: w, Height: h}, History: op.History + "R", Score:-1}
        }
      }
      countChan <- true
    }()
    
    go func() {
      nextBoard, error := MoveRight(w, h, op.Board.Panels)
      if error == NoError {
        if isEffective(w, h, op, nextBoard) {
          opChan <- Operation{Board: Board{Panels: nextBoard, Width: w, Height: h}, History: op.History + "L", Score:-1}
        }
      }
      countChan <- true
    }()
  }
  
  tickCount := 0
  isComplete := false
  for {
    select {
    case op := <- opChan:
      result = append(result, op)
    case <- compChan:
      tickCount++
      if tickCount == len(operations) {
        isComplete = true
      }
    }
    
    if isComplete {
      break
    }
  }
  
  limit := 2 << 8
  
  sort.Sort(result)
  if len(result) > limit && prevScore > result[0].Score {
    // numUnit := limit / 4
    ext := make(Operations, 0, limit)
    
    upper := result[:limit]
    // upper := result[:numUnit]
    // lower := result[len(result)-numUnit:]
    ext = append(ext, upper...)
    // ext = append(ext, lower...)
    result = ext
  }
  
  return result
}

func isEffective(w int, h int, op Operation, next []byte) bool {
  panels := make([]byte, len(op.Board.Panels))
  copy(panels, op.Board.Panels)
  
  for i:=0;i<len(op.History);i++ {
    prev := op.History[len(op.History)-i-1]
    if (prev == 'D') {
      panels, _ = MoveDown(w, h, panels)
    } else if (prev == 'U') {
      panels, _ = MoveUp(w, h, panels)
    } else if (prev == 'L') {
      panels, _ = MoveLeft(w, h, panels)
    } else if (prev == 'R') {
      panels, _ = MoveRight(w, h, panels)
    }
    
    n := 0
    for ;n<w*h;n++ {
      if (panels[n] != next[n]) {
        break
      }
    }
    
    if n == w*h {
      return false
    }
  }
  
  if ChainScore(op.Board) - ChainScore(Board{Panels:next, Width:w, Height:h}) > 1 {
    return false
  }
  
  return true
}

func ChainScore(b Board) int {
  fastest := 0
  latest := len(b.Panels)
  fastest2 := 0
  latest2 := len(b.Panels)
  
  n := 1;
  
  for i:=0;i<b.Width*b.Height;i++ {
    if b.Panels[i] == Wall {
      fastest = i
      break
    }
  }

  for i:=0;i<b.Width*b.Height;i++ {
    idx := b.Width*b.Height-i-1
    if b.Panels[idx] == Wall {
      latest = idx
      break
    }
  }
  
  for i:=0;i<b.Width*b.Height;i++ {
    idx := (b.Height-1-i%b.Height)*b.Width + i/b.Height
    if b.Panels[idx] == Wall {
      fastest2 = idx
      break
    }
  }
  
  for i:=0;i<b.Width*b.Height;i++ {
    idx := (i%b.Height+1)*b.Width - i/b.Height - 1
    if b.Panels[idx] == Wall {
      latest2 = idx
      break
    }
  }
  
  distance := fastest
  start := LeftTop
  if distance > len(b.Panels) - latest - 1 {
    distance = len(b.Panels) - latest - 1
    start = RightBottom
  }
  
  if distance > fastest2 {
    distance = fastest2
    start = LeftBottom
  }
  
  if distance > len(b.Panels) - latest2 - 1 {
    distance = len(b.Panels) - latest2 - 1
    start = RightTop
  }
  
  if start == LeftTop {
    for x:=1;x<b.Width*b.Height-1;x++ {
      i := 0
      for ;i<b.Width*b.Height;i++ {
        p := int(b.Panels[i])
        if p == x {
          if i > 1 && i % b.Width !=0  && int(b.Panels[i-1]) == x+1 {
            n++
          } else if i < len(b.Panels) - 1 && i % b.Width != b.Width-1 && int(b.Panels[i+1]) == x+1 {
            n++
          } else if i > b.Width && int(b.Panels[i-b.Width]) == x+1 {
            n++
          } else if i < len(b.Panels) - 1 - b.Width && int(b.Panels[i+b.Width]) == x+1 {
            n++
          }
          break
        }
      }
    }
  } else if start == RightBottom {
    for x:=b.Width*b.Height-1;x>0;x-- {
      i := 0
      for ;i<b.Width*b.Height;i++ {
        p := int(b.Panels[i])
        if p == x {
          if i > 1 && i % b.Width !=0  && int(b.Panels[i-1]) == x-1 {
            n++
          } else if i < len(b.Panels) - 1 && i % b.Width != b.Width-1 && int(b.Panels[i+1]) == x-1 {
            n++
          } else if i > b.Width && int(b.Panels[i-b.Width]) == x-1 {
            n++
          } else if i < len(b.Panels) - 1 - b.Width && int(b.Panels[i+b.Width]) == x-1 {
            n++
          }
          break
        }
      }
    }
  } else if start == LeftBottom {
    for x:=1;x<b.Width*b.Height-1;x++ {
      i := 0
      for ;i<b.Width*b.Height;i++ {
        p := int(b.Panels[i])
        if p == x {
          if i > 1 && i % b.Width !=0  && int(b.Panels[i-1]) == x-b.Width {
            n++
          } else if i < len(b.Panels) - 1 && i % b.Width != b.Width-1 && int(b.Panels[i+1]) == x-b.Width {
            n++
          } else if i > b.Width && int(b.Panels[i-b.Width]) == x-b.Width {
            n++
          } else if i < len(b.Panels) - 1 - b.Width && int(b.Panels[i+b.Width]) == x-b.Width {
            n++
          }
          break
        }
      }
    }
  } else {
    for x:=b.Width*b.Height-1;x>0;x-- {
      i := 0
      for ;i<b.Width*b.Height;i++ {
        p := int(b.Panels[i])
        if p == x {
          if i > 1 && i % b.Width !=0  && int(b.Panels[i-1]) == x+b.Width {
            n++
          } else if i < len(b.Panels) - 1 && i % b.Width != b.Width-1 && int(b.Panels[i+1]) == x+b.Width {
            n++
          } else if i > b.Width && int(b.Panels[i-b.Width]) == x+b.Width {
            n++
          } else if i < len(b.Panels) - 1 - b.Width && int(b.Panels[i+b.Width]) == x+b.Width {
            n++
          }
          break
        }
      }
    }
  }
  
  return n
}

func CalcDistance(b Board, checkNext int) float64 {
  distance := 0.0
  
  // // 1.
  // for i:=0;i<b.Width*b.Height;i++ {
  //   p := int(b.Panels[i])
  //   if p == Wall {
  //     continue
  //   }
  //   
  //   var x, y int
  //   if p != Blank {
  //     x = p%b.Width-(i+1)%b.Width
  //     y = p/b.Width-i/b.Width
  //   } else {
  //     x = (i+1)%b.Width
  //     y = i/b.Width
  //   }
  //   distance += math.Hypot(float64(x), float64(y))
  // }
  
  // 2.
  // for i:=0;i<b.Width*b.Height;i++ {
  //   p := int(b.Panels[i])
  //   if p == Wall || p == Blank{
  //     continue
  //   }
  //   
  //   x := p%b.Width-(i+1)%b.Width
  //   y := p/b.Width-i/b.Width
  //   distance += math.Hypot(float64(x), float64(y))
  // }
  
  // // 3.
  // for i:=0;i<b.Width*b.Height;i++ {
  //   p := int(b.Panels[i])
  //   if p == Wall || p == Blank {
  //     continue
  //   }
  //   
  //   x := p%b.Width-(i+1)%b.Width
  //   y := p/b.Width-(i+1)/b.Width
  //   distance += math.Hypot(float64(x), float64(y))
  // }
  
  // // 4.
  // for i:=0;i<b.Width*b.Height;i++ {
  //   p := int(b.Panels[i])
  //   if p == Wall || p == Blank {
  //     continue
  //   }
  //   
  //   x := p%b.Width-(i+1)%b.Width
  //   y := p/b.Width-(i+1)/b.Width
  //   distance += math.Hypot(float64(x), float64(y))
  // }
  // 
  // combo := 1
  // n := 0;
  // for x:=1;x<b.Width*b.Height-1;x++ {
  //   i := 0
  //   for ;i<b.Width*b.Height;i++ {
  //     p := int(b.Panels[i])
  //     if p == x {
  //       if i > 1 && i % b.Width !=0  && int(b.Panels[i-1]) == x+1 {
  //         n += combo
  //         combo++
  //       } else if i < len(b.Panels) - 1 && i % b.Width != b.Width-1 && int(b.Panels[i+1]) == x+1 {
  //         n += combo
  //         combo++
  //       } else if i > b.Width && int(b.Panels[i-b.Width]) == x+1 {
  //         n += combo
  //         combo++
  //       } else if i < len(b.Panels) - 1 - b.Width && int(b.Panels[i+b.Width]) == x+1 {
  //         n += combo
  //         combo++
  //       } else {
  //         combo = 1
  //       }
  //       break
  //     }
  //   }
  //   if i == b.Width*b.Height {
  //     combo = 1
  //   }
  // }
  // distance += 10.0 / float64(n)
  
  // 5.
  // combo := 1
  // n := 0;
  // for x:=1;x<b.Width*b.Height-1;x++ {
  //   i := 0
  //   for ;i<b.Width*b.Height;i++ {
  //     p := int(b.Panels[i])
  //     if p == x {
  //       if i > 1 && i % b.Width !=0  && int(b.Panels[i-1]) == x+1 {
  //         n += combo
  //         combo++
  //       } else if i < len(b.Panels) - 1 && i % b.Width != b.Width-1 && int(b.Panels[i+1]) == x+1 {
  //         n += combo
  //         combo++
  //       } else if i > b.Width && int(b.Panels[i-b.Width]) == x+1 {
  //         n += combo
  //         combo++
  //       } else if i < len(b.Panels) - 1 - b.Width && int(b.Panels[i+b.Width]) == x+1 {
  //         n += combo
  //         combo++
  //       } else {
  //         combo = 1
  //       }
  //       break
  //     }
  //   }
  //   if i == b.Width*b.Height {
  //     combo = 1
  //   }
  // }
  // distance = 1.0 / float64(n)
  
  // // 6.
  // 
  // fastest := 0
  // latest := len(b.Panels)
  // 
  // for i:=0;i<b.Width*b.Height;i++ {
  //   if b.Panels[i] == Wall {
  //     fastest = i
  //     break
  //   }
  // }
  // 
  // for i:=0;i<b.Width*b.Height;i++ {
  //   idx := b.Width*b.Height-i-1
  //   if b.Panels[idx] == Wall {
  //     latest = idx
  //     break
  //   }
  // }
  // 
  // combo := 1
  // n := 0;
  // if fastest > len(b.Panels) - latest - 1 {
  //   for x:=1;x<b.Width*b.Height-1;x++ {
  //     if combo > b.Width {
  //       combo = 1
  //     }
  //     
  //     i := 0
  //     for ;i<b.Width*b.Height;i++ {
  //       p := int(b.Panels[i])
  //       if p == x {
  //         if i > 1 && i % b.Width !=0  && int(b.Panels[i-1]) == x+1 {
  //           n += combo
  //           combo++
  //         } else if i < len(b.Panels) - 1 && i % b.Width != b.Width-1 && int(b.Panels[i+1]) == x+1 {
  //           n += combo
  //           combo++
  //         } else if i > b.Width && int(b.Panels[i-b.Width]) == x+1 {
  //           n += combo
  //           combo++
  //         } else if i < len(b.Panels) - 1 - b.Width && int(b.Panels[i+b.Width]) == x+1 {
  //           n += combo
  //           combo++
  //         } else {
  //           combo = 1
  //         }
  //         break
  //       }
  //     }
  //     if i == b.Width*b.Height {
  //       combo = 1
  //     }
  //   }
  //   
  //   for i:=0;i<b.Width*b.Height-1;i++ {
  //     if b.Panels[i] == uint8(i+1) || b.Panels[i] == Wall {
  //       n++
  //     } else {
  //       break;
  //     }
  //   }
  // } else {
  //   for x:=b.Width*b.Height-1;x>0;x-- {
  //     if combo > b.Width {
  //       combo = 1
  //     }
  //     
  //     i := 0
  //     for ;i<b.Width*b.Height;i++ {
  //       p := int(b.Panels[i])
  //       if p == x {
  //         if i > 1 && i % b.Width !=0  && int(b.Panels[i-1]) == x-1 {
  //           n += combo
  //           combo++
  //         } else if i < len(b.Panels) - 1 && i % b.Width != b.Width-1 && int(b.Panels[i+1]) == x-1 {
  //           n += combo
  //           combo++
  //         } else if i > b.Width && int(b.Panels[i-b.Width]) == x-1 {
  //           n += combo
  //           combo++
  //         } else if i < len(b.Panels) - 1 - b.Width && int(b.Panels[i+b.Width]) == x-1 {
  //           n += combo
  //           combo++
  //         } else {
  //           combo = 1
  //         }
  //         break
  //       }
  //     }
  //     if i == b.Width*b.Height {
  //       combo = 1
  //     }
  //   }
  //   
  //   for i:=b.Width*b.Height-2;i>0;i++ {
  //     if b.Panels[i] == uint8(i+1) || b.Panels[i] == Wall {
  //       n++
  //     } else {
  //       break;
  //     }
  //   }
  // }
  // 
  // distance += 10.0 / float64(n)
  
  // // 7.
  // for i:=0;i<b.Width*b.Height;i++ {
  //   p := int(b.Panels[i])
  //   if p == Wall || p == Blank {
  //     continue
  //   }
  //   
  //   x := p%b.Width-(i+1)%b.Width
  //   y := p/b.Width-(i+1)/b.Width
  //   distance += math.Hypot(float64(x), float64(y))
  // }
  // 
  // fastest := 0
  // latest := len(b.Panels)
  // 
  // for i:=0;i<b.Width*b.Height;i++ {
  //   if b.Panels[i] == Wall {
  //     fastest = i
  //     break
  //   }
  // }
  // 
  // for i:=0;i<b.Width*b.Height;i++ {
  //   idx := b.Width*b.Height-i-1
  //   if b.Panels[idx] == Wall {
  //     latest = idx
  //     break
  //   }
  // }
  // 
  // combo := 1
  // n := 1;
  // if fastest > len(b.Panels) - latest - 1 {
  //   for x:=1;x<b.Width*b.Height-1;x++ {
  //     if combo > b.Width {
  //       combo = 1
  //     }
  //     
  //     i := 0
  //     for ;i<b.Width*b.Height;i++ {
  //       p := int(b.Panels[i])
  //       if p == x {
  //         if i > 1 && i % b.Width !=0  && int(b.Panels[i-1]) == x+1 {
  //           n += combo
  //           combo++
  //         } else if i < len(b.Panels) - 1 && i % b.Width != b.Width-1 && int(b.Panels[i+1]) == x+1 {
  //           n += combo
  //           combo++
  //         } else if i > b.Width && int(b.Panels[i-b.Width]) == x+1 {
  //           n += combo
  //           combo++
  //         } else if i < len(b.Panels) - 1 - b.Width && int(b.Panels[i+b.Width]) == x+1 {
  //           n += combo
  //           combo++
  //         } else {
  //           combo = 1
  //         }
  //         break
  //       }
  //     }
  //     if i == b.Width*b.Height {
  //       combo = 1
  //     }
  //   }
  //   
  //   for i:=0;i<b.Width*b.Height-1;i++ {
  //     if b.Panels[i] == uint8(i+1) || b.Panels[i] == Wall {
  //       n++
  //     } else {
  //       break;
  //     }
  //   }
  // } else {
  //   for x:=b.Width*b.Height-1;x>0;x-- {
  //     if combo > b.Width {
  //       combo = 1
  //     }
  //     
  //     i := 0
  //     for ;i<b.Width*b.Height;i++ {
  //       p := int(b.Panels[i])
  //       if p == x {
  //         if i > 1 && i % b.Width !=0  && int(b.Panels[i-1]) == x-1 {
  //           n += combo
  //           combo++
  //         } else if i < len(b.Panels) - 1 && i % b.Width != b.Width-1 && int(b.Panels[i+1]) == x-1 {
  //           n += combo
  //           combo++
  //         } else if i > b.Width && int(b.Panels[i-b.Width]) == x-1 {
  //           n += combo
  //           combo++
  //         } else if i < len(b.Panels) - 1 - b.Width && int(b.Panels[i+b.Width]) == x-1 {
  //           n += combo
  //           combo++
  //         } else {
  //           combo = 1
  //         }
  //         break
  //       }
  //     }
  //     if i == b.Width*b.Height {
  //       combo = 1
  //     }
  //   }
  //   
  //   for i:=b.Width*b.Height-2;i>0;i-- {
  //     if b.Panels[i] == uint8(i+1) || b.Panels[i] == Wall {
  //       n++
  //     } else {
  //       break;
  //     }
  //   }
  // }
  // 
  // distance += 10.0 / float64(n)
  // 
  // if (checkNext > 0) {
  //   score := distance * float64(checkNext)
  //   
  //   var next []byte
  //   var err int
  //   next, err = MoveUp(b.Width, b.Height, b.Panels)
  //   
  //   if err == NoError {
  //     nBoard := Board{Panels:next, Width:b.Width, Height:b.Height}
  //     nScore := CalcDistance(nBoard, 0)
  //     if CheckBoard(nBoard) {
  //       return 0.0
  //     } else if distance > nScore {
  //       gScore := CalcDistance(nBoard, checkNext - 1)
  //       if score > gScore {
  //         score = gScore
  //       }
  //     }
  //   }
  //   next, err = MoveDown(b.Width, b.Height, b.Panels)
  //   if err == NoError {
  //     nBoard := Board{Panels:next, Width:b.Width, Height:b.Height}
  //     nScore := CalcDistance(nBoard, 0)
  //     if CheckBoard(nBoard) {
  //       return 0.0
  //     } else if distance > nScore {
  //       gScore := CalcDistance(nBoard, checkNext - 1)
  //       if score > gScore {
  //         score = gScore
  //       }
  //     }
  //   }
  //   next, err = MoveLeft(b.Width, b.Height, b.Panels)
  //   if err == NoError {
  //     nBoard := Board{Panels:next, Width:b.Width, Height:b.Height}
  //     nScore := CalcDistance(nBoard, 0)
  //     if CheckBoard(nBoard) {
  //       return 0.0
  //     } else if distance > nScore {
  //       gScore := CalcDistance(nBoard, checkNext - 1)
  //       if score > gScore {
  //         score = gScore
  //       }
  //     }
  //   }
  //   next, err = MoveRight(b.Width, b.Height, b.Panels)
  //   if err == NoError {
  //     nBoard := Board{Panels:next, Width:b.Width, Height:b.Height}
  //     nScore := CalcDistance(nBoard, 0)
  //     if CheckBoard(nBoard) {
  //       return 0.0
  //     } else if distance > nScore {
  //       gScore := CalcDistance(nBoard, checkNext - 1)
  //       if score > gScore {
  //         score = gScore
  //       }
  //     }
  //   }
  //   
  //   distance += score
  // }
  
  // 8.
  for i:=0;i<b.Width*b.Height;i++ {
    p := int(b.Panels[i])
    if p == Wall || p == Blank {
      continue
    }
    
    x := p%b.Width-(i+1)%b.Width
    y := p/b.Width-(i+1)/b.Width
    distance += math.Hypot(float64(x), float64(y))
  }
  
  fastest := 0
  latest := len(b.Panels)
  fastest2 := 0
  latest2 := len(b.Panels)
  
  n := 1;
  combo := 1
  chain := 0
  
  for i:=0;i<b.Width*b.Height;i++ {
    if b.Panels[i] == Wall {
      fastest = i
      break
    }
  }
  
  for i:=0;i<b.Width*b.Height;i++ {
    idx := b.Width*b.Height-i-1
    if b.Panels[idx] == Wall {
      latest = idx
      break
    }
  }
  
  for i:=0;i<b.Width*b.Height;i++ {
    idx := (b.Height-1-i%b.Height)*b.Width + i/b.Height
    if b.Panels[idx] == Wall {
      fastest2 = idx
      break
    }
  }

  for i:=0;i<b.Width*b.Height;i++ {
    idx := (i%b.Height+1)*b.Width - i/b.Height - 1
    if b.Panels[idx] == Wall {
      latest2 = idx
      break
    }
  }

  shortestDistance := fastest
  start := LeftTop
  if shortestDistance > len(b.Panels) - latest - 1 {
    shortestDistance = len(b.Panels) - latest - 1
    start = RightBottom
  }

  if shortestDistance > fastest2 {
    shortestDistance = fastest2
    start = LeftBottom
  }

  if shortestDistance > len(b.Panels) - latest2 - 1 {
    shortestDistance = len(b.Panels) - latest2 - 1
    start = RightTop
  }

  if start == LeftTop {
    for x:=1;x<b.Width*b.Height-1;x++ {
      if chain > 2 {
        break
      }
      
      if combo > b.Width {
        combo = 1
        chain++
      }
      
      i := 0
      for ;i<b.Width*b.Height;i++ {
        p := int(b.Panels[i])
        if p == x {
          if i > 1 && i % b.Width !=0  && int(b.Panels[i-1]) == x+1 {
            n += combo
            combo++
          } else if i < len(b.Panels) - 1 && i % b.Width != b.Width-1 && int(b.Panels[i+1]) == x+1 {
            n += combo
            combo++
          } else if i > b.Width && int(b.Panels[i-b.Width]) == x+1 {
            n += combo
            combo++
          } else if i < len(b.Panels) - 1 - b.Width && int(b.Panels[i+b.Width]) == x+1 {
            n += combo
            combo++
          } else {
            combo = 1
            chain++
          }
          break
        }
      }
      if i == b.Width*b.Height {
        combo = 1
        chain++
      }
    }
    
    for i:=0;i<b.Width*b.Height-1;i++ {
      if b.Panels[i] == uint8(i+1) || b.Panels[i] == Wall {
        n++
      } else {
        break;
      }
    }
  } else if start == RightBottom {
    for x:=b.Width*b.Height-1;x>0;x-- {
      if chain > 2 {
        break
      }
      
      if combo > b.Width {
        combo = 1
        chain++
      }
      
      i := 0
      for ;i<b.Width*b.Height;i++ {
        p := int(b.Panels[i])
        if p == x {
          if i > 1 && i % b.Width !=0  && int(b.Panels[i-1]) == x-1 {
            n += combo
            combo++
          } else if i < len(b.Panels) - 1 && i % b.Width != b.Width-1 && int(b.Panels[i+1]) == x-1 {
            n += combo
            combo++
          } else if i > b.Width && int(b.Panels[i-b.Width]) == x-1 {
            n += combo
            combo++
          } else if i < len(b.Panels) - 1 - b.Width && int(b.Panels[i+b.Width]) == x-1 {
            n += combo
            combo++
          } else {
            combo = 1
            chain++
          }
          break
        }
      }
      if i == b.Width*b.Height {
        combo = 1
        chain++
      }
    }
    
    for i:=b.Width*b.Height-2;i>0;i-- {
      if b.Panels[i] == uint8(i+1) || b.Panels[i] == Wall {
        n++
      } else {
        break;
      }
    }
  } else if start == LeftBottom {
    for x:=1;x<b.Width*b.Height-1;x++ {
      if chain > 2 {
        break
      }
      
      if combo > b.Width {
        combo = 1
        chain++
      }
      
      i := 0
      for ;i<b.Width*b.Height;i++ {
        p := int(b.Panels[i])
        if p == x {
          if i > 1 && i % b.Width !=0  && int(b.Panels[i-1]) == x-b.Width {
            n += combo
            combo++
          } else if i < len(b.Panels) - 1 && i % b.Width != b.Width-1 && int(b.Panels[i+1]) == x-b.Width {
            n += combo
            combo++
          } else if i > b.Width && int(b.Panels[i-b.Width]) == x-b.Width {
            n += combo
            combo++
          } else if i < len(b.Panels) - 1 - b.Width && int(b.Panels[i+b.Width]) == x-b.Width {
            n += combo
            combo++
          } else {
            combo = 1
            chain++
          }
          break
        }
      }
      if i == b.Width*b.Height {
        combo = 1
        chain++
      }
    }
    
    for i:=0;i<b.Width*b.Height;i++ {
      idx := (b.Height-1-i%b.Height)*b.Width + i/b.Height
      
      if b.Panels[idx] == uint8(idx+1) || b.Panels[idx] == Wall {
        n++
      } else if b.Panels[idx] == Blank {
        continue
      } else {
        break;
      }
    }
  } else {
    for x:=b.Width*b.Height-1;x>0;x-- {
      if chain > 2 {
        break
      }
      
      if combo > b.Width {
        combo = 1
        chain++
      }
      
      i := 0
      for ;i<b.Width*b.Height;i++ {
        p := int(b.Panels[i])
        if p == x {
          if i > 1 && i % b.Width !=0  && int(b.Panels[i-1]) == x+b.Width {
            n += combo
            combo++
          } else if i < len(b.Panels) - 1 && i % b.Width != b.Width-1 && int(b.Panels[i+1]) == x+b.Width {
            n += combo
            combo++
          } else if i > b.Width && int(b.Panels[i-b.Width]) == x+b.Width {
            n += combo
            combo++
          } else if i < len(b.Panels) - 1 - b.Width && int(b.Panels[i+b.Width]) == x+b.Width {
            n += combo
            combo++
          } else {
            combo = 1
            chain++
          }
          break
        }
      }
      if i == b.Width*b.Height {
        combo = 1
        chain++
      }
    }
    
    for i:=0;i<b.Width*b.Height;i++ {
      idx := (i%b.Height+1)*b.Width - i/b.Height - 1
      
      if b.Panels[idx] == uint8(idx+1) || b.Panels[idx] == Wall {
        n++
      } else if b.Panels[idx] == Blank {
        continue
      } else {
        break;
      }
    }
  }
  
  distance += 10.0 / float64(n)
  
  if (checkNext > 0) {
    score := distance * float64(checkNext)
    
    var next []byte
    var err int
    next, err = MoveUp(b.Width, b.Height, b.Panels)
    
    if err == NoError {
      nBoard := Board{Panels:next, Width:b.Width, Height:b.Height}
      nScore := CalcDistance(nBoard, 0)
      if CheckBoard(nBoard) {
        return 0.0
      } else if distance > nScore {
        gScore := CalcDistance(nBoard, checkNext - 1)
        if score > gScore {
          score = gScore
        }
      }
    }
    next, err = MoveDown(b.Width, b.Height, b.Panels)
    if err == NoError {
      nBoard := Board{Panels:next, Width:b.Width, Height:b.Height}
      nScore := CalcDistance(nBoard, 0)
      if CheckBoard(nBoard) {
        return 0.0
      } else if distance > nScore {
        gScore := CalcDistance(nBoard, checkNext - 1)
        if score > gScore {
          score = gScore
        }
      }
    }
    next, err = MoveLeft(b.Width, b.Height, b.Panels)
    if err == NoError {
      nBoard := Board{Panels:next, Width:b.Width, Height:b.Height}
      nScore := CalcDistance(nBoard, 0)
      if CheckBoard(nBoard) {
        return 0.0
      } else if distance > nScore {
        gScore := CalcDistance(nBoard, checkNext - 1)
        if score > gScore {
          score = gScore
        }
      }
    }
    next, err = MoveRight(b.Width, b.Height, b.Panels)
    if err == NoError {
      nBoard := Board{Panels:next, Width:b.Width, Height:b.Height}
      nScore := CalcDistance(nBoard, 0)
      if CheckBoard(nBoard) {
        return 0.0
      } else if distance > nScore {
        gScore := CalcDistance(nBoard, checkNext - 1)
        if score > gScore {
          score = gScore
        }
      }
    }
    
    distance += score
  }
  
  return distance
}

func MoveUp(w int, h int, d []byte) ([]byte, int) {
  error := 0
  isMoved := false
  result := make([]byte, len(d))
  copy(result, d)
  
  for y:=0;y<h;y++ {
    for x:=0;x<w;x++ {
      if d[y*w+x] == Blank {
        if y == h-1 {
          error = CantMove
        } else if d[(y+1)*w+x] == Wall {
          error = CantMove
        } else {
          result[y*w+x] = d[(y+1)*w+x]
          result[(y+1)*w+x] = Blank
        }
        
        isMoved = true
        break
      }
    }
    
    if isMoved {
      break
    }
  }
  
  return result, error
}

func MoveDown(w int, h int, d []byte) ([]byte, int) {
  error := 0
  isMoved := false
  result := make([]byte, len(d))
  copy(result, d)
  
  for y:=0;y<h;y++ {
    for x:=0;x<w;x++ {
      if d[y*w+x] == Blank {
        if y == 0 {
          error = CantMove
        } else if d[(y-1)*w+x] == Wall {
          error = CantMove
        } else {
          result[y*w+x] = d[(y-1)*w+x]
          result[(y-1)*w+x] = Blank
        }
        
        isMoved = true
        break
      }
    }
    
    if isMoved {
      break
    }
  }
  
  return result, error
}

func MoveLeft(w int, h int, d []byte) ([]byte, int) {
  error := 0
  isMoved := false
  result := make([]byte, len(d))
  copy(result, d)
  
  for y:=0;y<h;y++ {
    for x:=0;x<w;x++ {
      if d[y*w+x] == Blank {
        if x == w-1 {
          error = CantMove
        } else if d[y*w+x+1] == Wall {
          error = CantMove
        } else {
          result[y*w+x] = d[y*w+x+1]
          result[y*w+x+1] = Blank
        }
        
        isMoved = true
        break
      }
    }
    
    if isMoved {
      break
    }
  }
  
  return result, error
}

func MoveRight(w int, h int, d []byte) ([]byte, int) {
  error := 0
  isMoved := false
  result := make([]byte, len(d))
  copy(result, d)
  
  for y:=0;y<h;y++ {
    for x:=0;x<w;x++ {
      if d[y*w+x] == Blank {
        if x == 0 {
          error = CantMove
        } else if d[y*w+x-1] == Wall {
          error = CantMove
        } else {
          result[y*w+x] = d[y*w+x-1]
          result[y*w+x-1] = Blank
        }
        
        isMoved = true
        break
      }
    }
    
    if isMoved {
      break
    }
  }
  
  return result, error
}
