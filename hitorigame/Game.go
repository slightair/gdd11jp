package main

import (
  "fmt"
  "bufio"
  "os"
  "strings"
  "strconv"
)

const (
  deleted = -1
)

func main() {
  var line string
  var error os.Error
  var numTest int
  
  r := bufio.NewReader(os.Stdin)
  
  line, error = gets(r)
  _, error = fmt.Sscan(line, &numTest)
  if error != nil {
    fmt.Println(error)
    os.Exit(1)
  }
  
  for i:=0;i<numTest;i++ {
    var numParams int
    var params [10]int
    
    line, error = gets(r)
    _, error = fmt.Sscan(line, &numParams)
    if error != nil {
      fmt.Println(error)
      break
    }
    
    line, error = gets(r)
    if error != nil && error != os.EOF {
      fmt.Println(error)
      break
    }
    
    params = delete(params)
    fields := strings.Fields(line)
    for i:=0;i<numParams;i++ {
      params[i], error = strconv.Atoi(fields[i])
      if error != nil {
        fmt.Println(error)
        break
      }
    }
    
    paramsList := make([][10]int, 1)
    paramsList[0] = params
    
    limit := 21
    n := 0
    isCompleted := false
    
    for n=0;n<limit;n++ {
      for i:=0;i<len(paramsList);i++ {
        if isAllDeleted(paramsList[i]) {
          isCompleted = true
          break
        }
      }
      if isCompleted {
        break;
      }
      
      paramsList = split(paramsList)
    }
    
    fmt.Println(n)
  }
}


func gets(r *bufio.Reader) (string, os.Error) {
  var s string;
  line, error := r.ReadString('\n')
  
  if (error == nil || error == os.EOF) {
    s = strings.TrimRight(line, "\n")
  }
  
  return s, error
}

func divide(params [10]int) [10]int {
  var result [10]int
  
  for i:=0;i<10;i++ {
    if params[i] != deleted {
      result[i] = params[i] / 2
    } else {
      result[i] = params[i]
    }
  }
  
  return result
}

func delete(params [10]int) [10]int {
  var result [10]int
  
  for i:=0;i<10;i++ {
    if params[i] != deleted && params[i] % 5 == 0 {
      result[i] = deleted
    } else {
      result[i] = params[i]
    }
  }
  
  return result
}

func isAllDeleted(params [10]int) bool {
  for i:=0;i<10;i++ {
    if params[i] != deleted {
      return false
    }
  }
  
  return true
}

func split(paramsList [][10]int) [][10]int {
  result := make([][10]int, len(paramsList) * 2)
  
  for i:=0;i<len(paramsList);i++ {
    result[i*2] = divide(paramsList[i])
    result[i*2+1] = delete(paramsList[i])
  }
  
  return result
}