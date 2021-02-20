package main

import (
  "fmt"
  "net/http"
  "html/template"
  "encoding/json"
  "time"
  "bytes"

)

type FormData struct {
  Loaded bool
  Result string
}

type RespponsePostData struct {
  Id int64         `json:"id"`
}

type RespponseGetData struct {
  Answer string   `json:"answer"`
}

var api_url = "http://192.168.1.238:8080/"


// Display the named template
func display(w http.ResponseWriter, page string, is_load bool, result string) {
  t, err := template.ParseFiles("templates/index.html")

  if err != nil {
    fmt.Println(w, err.Error())
  }

  Data := FormData{is_load, result}
  t.Execute(w, Data)
}

func post_data(url string, text string) (session_id int64) {
  postBody, _ := json.Marshal(map[string]string{
    "text":text,
  })

  responseBody := bytes.NewBuffer(postBody)

  resp, err :=  http.Post(url, "application/json", responseBody)
  if err != nil {
    fmt.Println("An Error occured %v", err)
  }
  defer resp.Body.Close()

  post_data := RespponsePostData{}
  json.NewDecoder(resp.Body).Decode(&post_data)
  return post_data.Id
}

func get_data_from_session(session int64) <- chan string {
  url_session := fmt.Sprintf("%s%d", api_url, session)
  fmt.Println(url_session)

  r := make(chan string)
  go func() {
    defer close(r)
    for {
      time.Sleep(time.Second)
      resp, err := http.Get(url_session)
      if err != nil {
        panic(err)
      }
      defer resp.Body.Close()

      get_data := RespponseGetData{}
      json.NewDecoder(resp.Body).Decode(&get_data)

      if get_data.Answer != "" && get_data.Answer != "In work" {
        r <- get_data.Answer
        break
      }
    }
  }()
  return r
}

func upload_text(w http.ResponseWriter, r *http.Request) {
  key_sentence := r.FormValue("key_sentence")

  fmt.Println("key sentence = ", key_sentence)
  
  session := post_data(api_url, key_sentence)
  fmt.Printf("session = %d ", session)

  result_txt := <- get_data_from_session(session)
  fmt.Println(result_txt)

  display(w, "index", true, result_txt)
}

func view_handler(w http.ResponseWriter, r *http.Request) {
  switch r.Method {
  	case "GET":
  		display(w, "index", false, "")
  	case "POST":
  		upload_text(w, r)
	}
}

func main() {
	http.HandleFunc("/", view_handler)
  http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	http.ListenAndServe(":8000", nil)
}
