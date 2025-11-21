package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	. "L0_WB/models"
	. "L0_WB/repository"
	. "L0_WB/services"
)

// Кеш и счетчик
var chache = make(map[string]Order)
var i int = 0

func OrderEndPoint() {
	go GetDataFromKafka(chache)
	go Get5Orderd(chache)
	http.HandleFunc("/order/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		idStr := strings.TrimPrefix(r.URL.Path, "/order/")

		idStr = strings.TrimSuffix(idStr, "/")

		if idStr == "" {
			http.Error(w, "ID is required", http.StatusBadRequest)
			return
		}
		fmt.Println(idStr)
		for _, v := range chache {
			fmt.Println(v.OrderUID)
			if v.OrderUID == idStr {
				jsonData, err := json.Marshal(v)
				if err != nil {
					http.Error(w, "Failed to serialize order", http.StatusInternalServerError)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				w.Write(jsonData)
				return
			}
		}

		order, err := GetOrder(idStr)
		if err != nil {
			fmt.Println("Err with get order: ", err)
			http.Error(w, "Failed to get order ", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		js, err := json.Marshal(order)
		if err != nil {
			http.Error(w, "Failed to serialize ", http.StatusInternalServerError)
			return
		}
		w.Write(js)

	})
}
