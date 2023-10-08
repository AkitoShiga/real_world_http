package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

var image []byte

//　画像ファイルを取得する
func init() {
	var err error
	image, err = ioutil.ReadFile("./image.png")
	if err != nil {
		panic(err)
	}
}

// HTMLをブラウザに送信
// 画像をプッシュする
func handlerHtml(w http.ResponseWriter, r *http.Request) {
	// Pusherにキャスト可能であれば（HTTP/2で接続していたら）プッシュする
	pusher, ok := w.(http.Pusher)

	if ok {
		pusher.Push("/image", nil)
	}
	w.Header().Add("Content-Type", "test/html")
	fmt.Fprintf(w, `<html><body><img src="/image"></body></html>`)
}

func handlerImage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/png")
	w.Write(image)
}

func main() {
	http.HandleFunc("/", handlerHtml)
	http.HandleFunc("/image", handlerImage)
	fmt.Println("start http listening :18443")
	err := http.ListenAndServerTLS(":18443", "server.crt", "server.key", nil)
	fmt.Println(err)
}
