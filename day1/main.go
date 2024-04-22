/*
## Задача № 1
Написать API для указанных маршрутов(endpoints)
"/info"   // Информация об API
"/first"  // Случайное число
"/second" // Случайное число
"/add"    // Сумма двух случайных чисел
"/sub"    // Разность
"/mul"    // Произведение
"/div"    // Деление

*результат вернуть в виде JSON

"math/rand"
number := rand.Intn(100)
! не забудьте про Seed()


GET http://127.0.0.1:1234/first

GET http://127.0.0.1:1234/second

GET http://127.0.0.1:1234/add
GET http://127.0.0.1:1234/sub
GET http://127.0.0.1:1234/mul
GET http://127.0.0.1:1234/div
GET http://127.0.0.1:1234/info
*/

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

var (
	port    string = "8080"
	number1 int
	number2 int
	randGen = rand.New(rand.NewSource(time.Now().UnixNano()))
)

type InfoServer struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type CalcResponse struct {
	Number1   int     `json:"number1"`
	Number2   int     `json:"number2"`
	Operation string  `json:"operation"`
	Result    float64 `json:"result"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func main() {
	fmt.Println("Listen HTTP server on localhost: ", port)
	http.HandleFunc("/info", infoHandler)
	http.HandleFunc("/first", firstHandler)
	http.HandleFunc("/second", secondHandler)
	http.HandleFunc("/add", addHandler)
	http.HandleFunc("/sub", subHandler)
	http.HandleFunc("/mul", mulHandler)
	http.HandleFunc("/div", divHandler)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func infoHandler(w http.ResponseWriter, r *http.Request) {
	info := InfoServer{
		Name:    "Calculator API",
		Version: "1.0",
	}

	sendResponse(info, w)
}

func firstHandler(w http.ResponseWriter, r *http.Request) {
	number1 = randGen.Intn(100)
	calcResp := CalcResponse{
		Number1:   number1,
		Number2:   0,
		Operation: "random",
		Result:    float64(number1),
	}

	sendResponse(calcResp, w)
}

func secondHandler(w http.ResponseWriter, r *http.Request) {
	number2 = randGen.Intn(100)
	calcResp := CalcResponse{
		Number1:   number1,
		Number2:   number2,
		Operation: "random",
		Result:    float64(number2),
	}

	sendResponse(calcResp, w)
}

func addHandler(w http.ResponseWriter, r *http.Request) {
	result := number1 + number2
	calcResp := CalcResponse{
		Number1:   number1,
		Number2:   number2,
		Operation: "add",
		Result:    float64(result),
	}

	sendResponse(calcResp, w)
}

func subHandler(w http.ResponseWriter, r *http.Request) {
	result := float64(number1 - number2)
	calcResp := CalcResponse{
		Number1:   number1,
		Number2:   number2,
		Operation: "sub",
		Result:    result,
	}

	sendResponse(calcResp, w)
}

func mulHandler(w http.ResponseWriter, r *http.Request) {
	result := float64(number1) * float64(number2)
	calcResp := CalcResponse{
		Number1:   number1,
		Number2:   number2,
		Operation: "mul",
		Result:    result,
	}

	sendResponse(calcResp, w)
}

func divHandler(w http.ResponseWriter, r *http.Request) {
	if number2 == 0 {
		errorResp := ErrorResponse{
			Error: "Division by zero",
		}
		sendResponse(errorResp, w)
		return
	}
	result := float64(number1) / float64(number2)
	calcResp := CalcResponse{
		Number1:   number1,
		Number2:   number2,
		Operation: "div",
		Result:    result,
	}

	sendResponse(calcResp, w)
}

func sendResponse(resp interface{}, w http.ResponseWriter) {
	bytes, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(bytes))
}
