package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

// 定义一个简单的页面模板
const tmpl = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>mie</title>
</head>
<body>
    <h1>Welcome to My Go Website</h1>
    <p>{{.Message}}</p>
</body>
</html>
`

// 页面数据结构体
type PageData struct {
	Message string
}

// 渲染模板的函数
func renderTemplate(w http.ResponseWriter, tmplStr string, data PageData) {
	tmpl, err := template.New("webpage").Parse(tmplStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// 处理根路径请求的函数
func rootHandler(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		Message: "mie",
	}
	renderTemplate(w, tmpl, data)
}

func main() {
	http.HandleFunc("/home", rootHandler)
	fmt.Println("Starting server at :1026")
	if err := http.ListenAndServe(":1026", nil); err != nil {
		log.Fatalf("Could not start server: %s\n", err.Error())
	}
}
