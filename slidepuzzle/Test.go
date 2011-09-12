package main

import (
  "fmt"
  "bufio"
  "os"
  "strings"
  p "./libPuzzle"
)

func main() {
  // runtime.GOMAXPROCS(4)
  r := bufio.NewReader(os.Stdin)
  
  width, height, panels := p.ParseProblem(r)
  answer, _ := gets(r)
  
  for i:=0;i<len(answer);i++ {
    fmt.Println(p.CalcDistance(p.Board{Panels: panels, Width: width, Height: height}))
    next := answer[i]
    if (next == 'D') {
      panels, _ = p.MoveUp(width, height, panels)
    } else if (next == 'U') {
      panels, _ = p.MoveDown(width, height, panels)
    } else if (next == 'L') {
      panels, _ = p.MoveRight(width, height, panels)
    } else if (next == 'R') {
      panels, _ = p.MoveLeft(width, height, panels)
    }
  }
  fmt.Println(p.CalcDistance(p.Board{Panels: panels, Width: width, Height: height}))
  
  for y:=0;y<height;y++ {
    for x:=0;x<width;x++ {
      var p byte
      b := panels[x+y*width]
      
      if b < 10 {
        p = b + '0'
      } else if b < 36 {
        p = b + 'A' - 10
      } else {
        p = '='
      }
      
      fmt.Print(string(p))
    }
    fmt.Println()
  }
}

func gets(r *bufio.Reader) (string, os.Error) {
  var s string;
  line, error := r.ReadString('\n')
  if error == nil || error == os.EOF {
    s = strings.TrimRight(line, "\n")
  }
  
  return s, error
}
