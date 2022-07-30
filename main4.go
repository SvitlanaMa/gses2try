package main

import ("errors"
		"fmt"
		"io/ioutil"
		"net/http"
		"os"
		"encoding/json"
		"strconv"
		"strings"
		"bufio"
		"net/smtp")
// структура відповіді про помилку, можливо, щось зайве
type bad_resp struct {
	Description string
	Status_code int 	`json:"status-code"`
	Content_Type string `json:"Content-Type"`
}
// структура відповіді ок для завдання 1, можливо, щось зайве
type good_resp_rate struct {
	Rate float64		`json:"rate"`
	Description string
	Status_code int 	`json:"status-code"`
	Content_Type string	`json:"Content-Type"`
}
// функція для отримання ціни біткоіна
func get_rate() (res_rate float64, err error) {
	// ключ для некомерційного використання, не можна робити більше одного запиту у секунду 
	key := "21fa01a4d3ae91cd5493ede2312131ea32d064b4"
	url := "https://api.nomics.com/v1/currencies/ticker?key=" + key + "&ids=BTC&convert=UAH"
	// робимо гет запит та перевіряємо чи є помилка 
	res, err := http.Get(url)
	if err != nil {
		return 0, fmt.Errorf("Неможливо отримати дані: %v", err)	
	}
	// читаємо відповідь та перевіряємо чи є помилка
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return 0, fmt.Errorf("Не вдалося прочитати відповідь: %v", err)
	}
	// щось по типу шаблона відповіді, щоб вірно спарсити джейсон
	var data []map[string]interface{}
	// декодуємо джейсон
	err1 := json.Unmarshal([]byte(resBody), &data)
	if err1 != nil {
		return 0, fmt.Errorf("Не вдалося отримати дані: %v", err1)
	}
	// беремо тільки ціну, переводимо її в десяткове
	res_rate, err2 := strconv.ParseFloat(data[0]["price"].(string), 64)
	if err2 != nil {
		return 0, fmt.Errorf("Не вдалося отримати дані: %v", err2)
	}
	// повертаємо ціну та відсутність помилки
	return res_rate, nil	
}
//функція для відправки листа на вказану адресу
func send(addr string) (err error) {
	// тестовий сервіс mailtrap не відправляє насправді, можна перевірити результат в акаунті
	from := "svitlana@ttt.ttt"
	user := "16e7d8eb78310e"
	password := "334795e1b21974" 
	// переводимо адресу отримувача у зріз
	to := []string{
	  addr,
	}  
	// дані smtp сервера 
	smtpHost := "smtp.mailtrap.io"
	smtpPort := "2525" 
	// повідомлення із заголовками. спершу отримуємо ціну біткоіна
	btc, err := get_rate()
	if err != nil {
		return fmt.Errorf("Не вдалося отримати дані: %v", err)
	}
	// переводимо її у форматований текст
	btc1 := fmt.Sprintf("%.2f", btc)
	// складаємо повідомлення (лист)
	message := []byte("From: "+ from + "\r\n" +
        "To: " + addr + "\r\n" +
        "Subject: Поточний курс BTC\r\n\r\n" +
        "Поточний курс BTC у гривні: " + btc1 + "\r\n")	
	// аутентифікація
	auth := smtp.PlainAuth("", user, password, smtpHost)	
	// відправляємо лист
	err1 := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err1 != nil {
	  return fmt.Errorf("Лист не відправлено: %v", err1)
	}
	//все гаразд, помилок немає
	return nil
}

func getRate(w http.ResponseWriter, r *http.Request) {
	// лише гет запити
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(bad_resp{"Невірний метод", 405, "application/json"})
		fmt.Printf("Отримали запит до %s, відповідь - 405\n", r.URL.Path)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	// отримаємо ціну біткоїну
	resp, err := get_rate()
	// якщо помилка, записуємо статус та надсилаємо відповідь у json
	if err != nil {
		bad := bad_resp {"Invalid status value", 400, "application/json"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(bad)
		fmt.Printf("Отримали запит до %s, відповідь - 400 помилка: %s\n", r.URL.Path, err)
		return	
	}
	// якщо все ок, статус ок та джейсон відповідь
	good := good_resp_rate{resp, "Ok", 200, "application/json"}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(good)
	fmt.Printf("Отримали запит до %s, відповідь - 200 ОК\n", r.URL.Path)
}

func subscribe(w http.ResponseWriter, r *http.Request) {
	// лише пост запити
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(bad_resp{"Невірний метод", 405, "application/json"})
		fmt.Printf("Отримали запит до %s, відповідь - 405\n", r.URL.Path)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	// із запиту отримаємо дані, прибираємо пропуски. Вважаємо, що є форма
	email := strings.TrimSpace(r.PostFormValue("email"))
	// читаємо файл
	file, err := os.Open("emails.txt")
    if err != nil{
		// можливо, потрібно конкретизувати, що відповідати в разі таких помилок
		bad := bad_resp {"Invalid status value", 400, "application/json"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(bad)
		fmt.Printf("Отримали запит до %s, відповідь - 400 помилка: %s\n", r.URL.Path, err)
        return 
	}
	// треба закрити файл в кінці виконання фукнції
	// defer file.Close() але не у випадку, коли розробник ще не навчився 
	//відкривати файл для читання та запису одночасно 	
	//скануємо рядки у файлі
	fileScanner := bufio.NewScanner(file)
    fileScanner.Split(bufio.ScanLines)
  
    for fileScanner.Scan() {
		if fileScanner.Text() == email {
			bad := bad_resp {"Ця адреса вже є у базі", 409, "application/json"}
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(bad)
			file.Close()
			fmt.Printf("Отримали запит до %s, відповідь - 409\n", r.URL.Path)
			return
		}   
	}
	file.Close()
	//відкриваємо файл для додавання рядків
	file1, err := os.OpenFile("emails.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
    if err != nil{
		bad := bad_resp {"Щось пойшло не так", 400, "application/json"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(bad)
        fmt.Printf("Отримали запит до %s, відповідь - помилка: %s\n", r.URL.Path, err) 
        return 
	}
	// додаємо адресу
	file1.WriteString(email+"\r\n")
	file1.Close()
	ok := bad_resp {"Адресу додано", 200, "application/json"}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ok)
	fmt.Printf("Отримали запит до %s, відповідь - 200 ОК\n", r.URL.Path)
}

func sendLetters(w http.ResponseWriter, r *http.Request) {
	// лише пост запити
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(bad_resp{"Невірний метод", 405, "application/json"})
		fmt.Printf("Отримали запит до %s, відповідь - 405\n", r.URL.Path)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	// відкриваємо файл для читання
	file, err := os.Open("emails.txt")
    if err != nil{
        bad := bad_resp {"Invalid status value", 400, "application/json"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(bad)
		fmt.Printf("Отримали запит до %s, відповідь - 400 помилка: %s\n", r.URL.Path, err)
        return 
	}
	// треба закрити файл в кінці виконання фукнції
	defer file.Close() 
	//скануємо рядки у файлі	
	fileScanner := bufio.NewScanner(file)
    fileScanner.Split(bufio.ScanLines)
  
    for fileScanner.Scan() {
		// відправляємо лист за кожною адресою із файла
		err := send(fileScanner.Text())
		if err != nil{
			bad := bad_resp {"Invalid status value", 400, "application/json"}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(bad)
			fmt.Printf("Отримали запит до %s, відповідь - 400 помилка: %s\n", r.URL.Path, err)
			return
		}   
	}
	w.WriteHeader(http.StatusOK)
	ok := bad_resp {"E-mailʼи відправлено", 200, "application/json"}
	json.NewEncoder(w).Encode(ok)
	fmt.Printf("Отримали запит до %s, відповідь - 200 ОК\n", r.URL.Path)
}

func main() {
	http.HandleFunc("/api/rate", getRate)
	http.HandleFunc("/api/subscribe", subscribe)
	http.HandleFunc("/api/sendEmails", sendLetters)
	
	fmt.Printf("GSES2 BTC application\nСервер розпочав роботу. Ласкаво просимо!\n\n")

	err := http.ListenAndServe(":3333", nil)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("Помилка серверу, сервер закрито\n")
	} else if err != nil {
		fmt.Printf("Текст помилки: %s\n", err)
		os.Exit(1)
	}
}