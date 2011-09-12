package main

import (
  "fmt"
  "bufio"
  "os"
  "strings"
  "strconv"
  p "./libPuzzle"
  // "runtime"
)

func main() {
  // runtime.GOMAXPROCS(2)
  var line string
  var error os.Error
  var lx, rx, ux, dx int
  
  r := bufio.NewReader(os.Stdin)
  
  line, error = gets(r)
  _, error = fmt.Sscan(line, &lx, &rx, &ux, &dx)
  if error != nil {
    fmt.Println(error)
    os.Exit(1)
  }
  
  line, error = gets(r)
  numProblems, error := strconv.Atoi(line)
  if error != nil {
    fmt.Println(error)
    os.Exit(1)
  }
  
  limit := 300
  // sizeLimit := 16
  // sizeLimit := 36
  // operationLimit := 2 << 21
  operationLimit := 2 << 16
  // operationLimit := 2 << 3
  
  for i:=0;i<numProblems;i++ {
    // runtime.GC()
    isGoal := false
    width, height, board := p.ParseProblem(r)
    prevScore := 100000.0
    
    // if width * height > sizeLimit {
    //   fmt.Println()
    //   continue
    // }
    
    operations := make(p.Operations, 1)
    operations[0] = p.Operation{Board:p.Board{Panels: board, Width: width, Height: height}, History:""}
    
    j:=0
    for ;j<limit;j++ {
      // fmt.Println(j, len(operations), operations[0].Score, operations[0].History)
      for k:=0;k<len(operations);k++ {
        if p.CheckBoard(operations[k].Board) {
          isGoal = true
          fmt.Println(operations[k].History)
          break
        }
      }
      if isGoal {
        break
      }
      
      if len(operations) > operationLimit {
        // // for z:=0;z<len(operations);z++ {
        // //   fmt.Println(operations[z])
        // // }
        fmt.Println()
        break
      }
      
      operations = p.Tick(width, height, operations, prevScore)
      if len(operations) < 1 {
        fmt.Println()
        break
      }
      prevScore = operations[0].Score
    }
    if j==limit {
      // fmt.Println(j, len(operations), p.CalcDistance(operations[0].Board), operations[0].History)
      fmt.Println()
    }
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
