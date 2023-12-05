package main

import (
	"fmt"
	"log"
	"net"
	"sync"
)

type ReportEntry struct {
	ID           int    `json:"Id"`
	PID          *int   `json:"Pid"`
	OriginalURL  string `json:"FullURL,omitempty"`
	ShortURL     string `json:"ShortenURL,omitempty"`
	SourceIP     string `json:"SourceIP"`
	TimeInterval string `json:"TimeInterval"`
	Count        int    `json:"Count"`
}

type ReportData struct {
	Entries []ReportEntry `json:"entries"`
}

type DetailReport struct {
	Count   int                      `json:"Count,omitempty"`
	Details map[string]*DetailReport `json:"Details,omitempty"`
}

var mu sync.Mutex
var reportData ReportData

func startStatisticService() {
	listener, err := net.Listen("tcp", ":5252")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	fmt.Println("Сервис статистики запущен. Ожидание подключений...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Ошибка при подключении клиента:", err)
			continue
		}

		fmt.Println("Подключение от", conn.RemoteAddr().String())

		go handleConnection(conn)
	}
}

func startUserInterface() {
	for {
		displayMenu()

		choice := getUserChoice()

		switch choice {
		case 1:
			existingEntries := reportData.Entries
			sendJSONRequest(existingEntries)
		case 2:
			fmt.Println("Отчёт генерируется. Введите порядок детализаций (например, \"SourceIP TimeInterval URL\"): ")

			// Запрос у пользователя каждого элемента последовательности
			var detailsOrder []string
			for i := 0; i < 3; i++ {
				var detail string
				fmt.Scan(&detail)
				detailsOrder = append(detailsOrder, detail)
			}

			report := generateReport(detailsOrder)
			fmt.Println("Отчет:")
			printJSON(report)
			// Сохранение отчета в файл
			saveReportToFile(report, "rep.json")
		case 3:
			fmt.Println("Выход из программы.")
			return
		default:
			fmt.Println("Некорректный выбор. Пожалуйста, попробуйте ещё раз и введите число от 1 до 3.")
		}
	}
}

func displayMenu() {
	fmt.Println("\nМеню:")
	fmt.Println("1. Подгрузить данные с сервера (SENDJSON)")
	fmt.Println("2. Выполнить отчет")
	fmt.Println("3. exit")
	fmt.Print("Выберите операцию: ")
}

func getUserChoice() int {
	var choice int
	fmt.Scan(&choice)
	return choice
}

func main() {
	go startStatisticService()
	startUserInterface()
}
