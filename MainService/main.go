package main

import (
	. "L0_WB/controllers"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}
func main() {
	go func() {
		time.Sleep(30 * time.Second)
		APP_PORT, e := os.LookupEnv("APP_PORT")
		if e != nil {
			fmt.Println(err)
			return
		}
		client := http.Client{}
		resp, err := client.Get("http://localhost:" + APP_PORT + "/order/order_888_it")
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
