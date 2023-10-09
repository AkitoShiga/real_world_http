package main

import (
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"time"
)

var html []byte

// HTMLをブラウザに送信
func handlerHtml(w http.ResponseWriter, r *http.Request) {

	// Pusherにキャスト可能であればプッシュする
	w.Header().Add("Content-Type", "text/html")
	w.Write(html)
}

// 素数をブラウザに送信
func handlerPrimeSSE(w http.ResponseWriter, r *http.Request) {
	// writerのチャンク送信をするFlusherを呼び出す
	flusher, ok := w.(http.Flusher)

	// エラーハンドリング
	if !ok {
		http.Error(w, "Streamning unsupported!", http.StatusInternalServerError)
		return
	}

	// 接続断の検知用にコンテキストを取得
	ctx := r.Context()
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "kee-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var num int64 = 1
	for id := 1; id < 100; id++ {

		// 通信が切れても終了
		select {
		case <-ctx.Done():
			fmt.Println("Connection closed from client")
			return
		default:
			// do nothing
		}

		for {
			num++
			// 確率論的に素数を求める
			if big.NewInt(num).ProbablyPrime(20) {
				fmt.Println(num)
				fmt.Fprintf(w, "data: {\"id\": %d, \"number\": %d}\n\n", id, num)
				flusher.Flush()
				time.Sleep(time.Second)
				break
			}
		}
		time.Sleep(time.Second)
	}
	// 100個超えたら放送終了
	fmt.Println("Connection closed from server")
}

func main() {
	var err error

	html, err = ioutil.ReadFile("index.html")
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", handlerHtml)
	http.HandleFunc("/prime", handlerPrimeSSE)
	fmt.Println("start http listening :18888")
	err = http.ListenAndServe(":18888", nil)
	fmt.Println(err)
}
