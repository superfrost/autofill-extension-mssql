package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/joho/godotenv"
)

// User представляет структуру пользователя
type User struct {
	MKABID        int    `json:"id"`
	NAME          string `json:"name"`
	FAMILY        string `json:"family"`
	OT            string `json:"ot"`
	SS            string `json:"ss"`
	S_DOC         string `json:"s_doc"`
	N_DOC         string `json:"n_doc"`
	DATE_BD       string `json:"date_bd"`
	ADRES         string `json:"adres"`
	Rf_kl_SexID   string `json:"sexid"`
	ContactEmail  string `json:"contactEmail"`
	ContactMPhone string `json:"contactMPhone"`
}

// Response представляет структуру ответа
type Response struct {
	Users []User `json:"users"`
	Error string `json:"error,omitempty"`
}

func main() {
	// Загрузка переменных окружения из .env файла
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка загрузки .env файла:", err)
	}

	// Получение данных для подключения из переменных окружения

	query := url.Values{}
	query.Add("database", os.Getenv("DB_NAME"))

	u := &url.URL{
		Scheme:   "sqlserver",
		User:     url.UserPassword(os.Getenv("DB_USER"), os.Getenv("DB_PASS")),
		Host:     net.JoinHostPort(os.Getenv("DB_HOST"), os.Getenv("DB_PORT")),
		RawQuery: query.Encode(),
	}

	// Обработка маршрута /list
	http.HandleFunc("/list", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
			return
		}

		var input struct {
			Query string `json:"query"`
		}

		// Декодирование JSON из тела запроса
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
			return
		}

		log.Println("query:", input.Query)

		// Получение списка пользователей из базы данных
		users, err := getUsers(u.String(), input.Query)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Формирование ответа
		response := Response{Users: users}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	// Запуск сервера
	port := ":8070"
	fmt.Printf("Сервер запущен на порту %s\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal("Ошибка при запуске сервера:", err)
	}
}

// getUsers выполняет запрос к базе данных и возвращает список пользователей
func getUsers(connString, query string) ([]User, error) {
	var users []User

	// Подключение к базе данных
	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		return nil, fmt.Errorf("не удалось подключиться к базе данных: %v", err)
	}
	defer db.Close()

	// Выполнение запроса
	rows, err := db.Query("SELECT MKABID, FAMILY, NAME, OT, SS, S_DOC, N_DOC, DATE_BD, ADRES, rf_kl_SexID, contactEmail, contactMPhone FROM hlt_MKAB WHERE SS LIKE @p1", "%"+query+"%")
	if err != nil {
		return nil, fmt.Errorf("ошибка при выполнении запроса: %v", err)
	}
	defer rows.Close()

	// Чтение результатов
	for rows.Next() {
		var user User
		if err := rows.Scan(
			&user.MKABID,
			&user.NAME,
			&user.FAMILY,
			&user.OT,
			&user.SS,
			&user.S_DOC,
			&user.N_DOC,
			&user.DATE_BD,
			&user.ADRES,
			&user.Rf_kl_SexID,
			&user.ContactEmail,
			&user.ContactMPhone); err != nil {
			return nil, fmt.Errorf("ошибка при чтении данных: %v", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при обработке результатов: %v", err)
	}

	return users, nil
}
