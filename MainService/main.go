package main

import (
	. "L0_WB/controllers"
	"fmt"
	"io"
	"net/http"
	"time"

	_ "github.com/lib/pq"
)

func main() {
	go func() {
		time.Sleep(30 * time.Second)

		client := http.Client{}
		resp, err := client.Get("http://localhost:8081/order/order_888_it")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer resp.Body.Close()
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Ошибка при чтении тела ответа:", err)
			return
		}
		fmt.Println("Ответ сервера:", string(bodyBytes))
	}()

	go OrderEndPoint()

	http.ListenAndServe(":8081", nil)
}
