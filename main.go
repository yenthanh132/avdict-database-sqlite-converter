package main

import (
  "bufio"
  "os"
  "strings"
  "fmt"
  _ "github.com/mattn/go-sqlite3"
  "database/sql"
)

type Dict struct {
  Word        string `json:"word"`
  Data        string `json:"data"`
  Description string `json:"description"`
  Pronounce   string `json:"pronounce"`
}

type Stack []int

func (s Stack) Push(v int) Stack {
  return append(s, v)
}

func (s Stack) Top() int {
  l := len(s)
  if l == 0 {
    return 0
  }
  return s[l-1]
}

func (s Stack) Pop() (Stack, int) {
  l := len(s)
  return s[:l-1], s[l-1]
}

var html string
var myStack Stack

func GetTag(val int) string {
  if val == 1 {
    return "</ol>"
  } else if val == 2 {
    return "</li>"
  } else if val == 3 {
    return "</ul>"
  } else if val == 4 {
    return "</li>"
  } else if val == 5 {
    return "</ul>"
  } else if val == 6 {
    return "</li>"
  }
  return ""
}

func AddCloseTag(minVal int) {
  for myStack.Top() >= minVal {
    var p int
    myStack, p = myStack.Pop()
    html += GetTag(p)
  }
}

func FormatLine(line string) string {
  if strings.Index(line, "+") != -1 {
    line = strings.Replace(line, "+", ":<i>", 1)
    line += "</i>"
  }
  return line
}

func generate_av() {
  db, _ := sql.Open("sqlite3", "dict_hh.db")
  defer db.Close()
  file, _ := os.Open("anhviet109K.txt")
  defer file.Close()
  scanner := bufio.NewScanner(file)

  result := make([]Dict, 0)
  d := Dict{Word: ""}
  html = ""
  idx := 0
  db.Exec("BEGIN TRANSACTION")
  stmt, _ := db.Prepare("INSERT INTO av(word, html, description, pronounce) values(?,?,?,?)")
  description_state := 0
  for scanner.Scan() {
    line := scanner.Text()
    line = strings.TrimSpace(line)
    if strings.Contains(line, "@") {
      AddCloseTag(1)
      line = strings.Replace(line, "@", "", -1)
      if len(d.Word) > 0 {
        d.Data = html
        result = append(result, d)
        stmt.Exec(d.Word, d.Data, d.Description, d.Pronounce)
        idx += 1
        if idx%10 == 0 {
          fmt.Println("Process: ", idx)
        }
      }
      items := strings.Split(line, "/")
      data_word := strings.TrimSpace(items[0])
      data_pronounce := ""
      if data_word == "1 to 1 relationship" {
        break
      }
      if len(items) > 1 {
        data_pronounce = strings.TrimSpace(items[1])
      }
      description_state = 0
      d = Dict{Word: data_word}
      d.Pronounce = data_pronounce
      myStack = make([]int, 0)
      html = fmt.Sprintf("<h1>%s</h1><h3><i>/%s/</i></h3>", data_word, data_pronounce)
    } else if strings.HasPrefix(line, "*  ") {
      AddCloseTag(1)
      line = FormatLine(line[3:])
      html += fmt.Sprintf("<h2>%s</h2>", line)
      if description_state == 0 {
        description_state = 1
        d.Description += line + ": "
      }
    } else if strings.HasPrefix(line, "!") {
      line = FormatLine(line[1:])
      AddCloseTag(2)
      if myStack.Top() < 1 {
        html += `<h2>thành ngữ</h2><ol>`
        myStack = myStack.Push(1)
      }
      html += "<li>"
      html += line
      myStack = myStack.Push(2)
      if description_state <= 1 {
        description_state = 2
        d.Description += line
      }
    } else if strings.HasPrefix(line, "- ") {
      line = FormatLine(line[2:])
      AddCloseTag(4)
      if myStack.Top() < 3 {
        html += `<ul>`
        myStack = myStack.Push(3)
      }
      html += "<li>"
      html += line
      myStack = myStack.Push(4)
      if description_state <= 1 {
        description_state = 2
        d.Description += line
      }
    } else if strings.HasPrefix(line, "=") {
      line = FormatLine(line[1:])
      AddCloseTag(6)
      if myStack.Top() < 5 {
        html += `<ul style="list-style-type:circle">`
        myStack = myStack.Push(5)
      }
      html += "<li>"
      html += line
      myStack = myStack.Push(6)
      if description_state <= 1 {
        description_state = 2
        d.Description += line
      }
    } else if len(strings.TrimSpace(line)) > 0 {
      html += line + "<br/>"
    }
  }
  if len(d.Word) > 0 {
    d.Data = html
    result = append(result, d)
    stmt.Exec(d.Word, d.Data)
    idx += 1
    if idx%10 == 0 {
      fmt.Println("Process: ", idx)
    }
  }
  db.Exec("END TRANSACTION")
}

func generate_va() {
  db, _ := sql.Open("sqlite3", "dict_hh.db")
  defer db.Close()
  file, _ := os.Open("vietanh.txt")
  defer file.Close()
  scanner := bufio.NewScanner(file)

  result := make([]Dict, 0)
  d := Dict{Word: ""}
  html = ""
  idx := 0
  db.Exec("BEGIN TRANSACTION")
  stmt, _ := db.Prepare("INSERT INTO va(word, html, description, pronounce) values(?,?,?,?)")
  description_state := 0
  for scanner.Scan() {
    line := scanner.Text()
    line = strings.TrimSpace(line)
    if strings.Contains(line, "@") {
      AddCloseTag(1)
      line = strings.Replace(line, "@", "", -1)
      if len(d.Word) > 0 {
        d.Data = html
        result = append(result, d)
        stmt.Exec(d.Word, d.Data, d.Description, d.Pronounce)
        idx += 1
        if idx%10 == 0 {
          fmt.Println("Process: ", idx)
        }
      }
      description_state = 0
      items := strings.Split(line, "/")
      data_word := strings.TrimSpace(items[0])
      d = Dict{Word: data_word}
      myStack = make([]int, 0)
      html = fmt.Sprintf("<h1>%s</h1>", data_word)
    } else if strings.HasPrefix(line, "* ") {
      AddCloseTag(1)
      line = FormatLine(line[2:])
      html += fmt.Sprintf("<h2>%s</h2>", line)
      if description_state == 0 {
        description_state = 1
        d.Description += line + ": "
      }
    } else if strings.HasPrefix(line, "!") {
      line = FormatLine(line[1:])
      AddCloseTag(2)
      if myStack.Top() < 1 {
        html += `<h2>idioms</h2><ol>`
        myStack = myStack.Push(1)
      }
      html += "<li>"
      html += line
      myStack = myStack.Push(2)
      if description_state <= 1 {
        description_state = 2
        d.Description += line
      }
    } else if strings.HasPrefix(line, "- ") {
      line = FormatLine(line[2:])
      AddCloseTag(4)
      if myStack.Top() < 3 {
        html += `<ul>`
        myStack = myStack.Push(3)
      }
      html += "<li>"
      html += line
      myStack = myStack.Push(4)
      if description_state <= 1 {
        description_state = 2
        d.Description += line
      }
    } else if strings.HasPrefix(line, "=") {
      line = FormatLine(line[1:])
      AddCloseTag(6)
      if myStack.Top() < 5 {
        html += `<ul style="list-style-type:circle">`
        myStack = myStack.Push(5)
      }
      html += "<li>"
      html += line
      myStack = myStack.Push(6)
      if description_state <= 1 {
        description_state = 2
        d.Description += line
      }
    } else if len(strings.TrimSpace(line)) > 0 {
      html += line + "<br/>"
    }
  }
  if len(d.Word) > 0 {
    d.Data = html
    result = append(result, d)
    stmt.Exec(d.Word, d.Data)
    idx += 1
    if idx%10 == 0 {
      fmt.Println("Process: ", idx)
    }
  }
  db.Exec("END TRANSACTION")
}

func main() {
  generate_av()
  generate_va()
}
