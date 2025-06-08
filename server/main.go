package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/joho/godotenv"
)

// User представляет структуру пользователя
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
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
	connString := os.Getenv("MSSQL_CONNECTION_STRING")
	if connString == "" {
		log.Fatal("Ошибка: переменная окружения MSSQL_CONNECTION_STRING не установлена.")
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

		// Получение списка пользователей из базы данных
		users, err := getUsers(connString, input.Query)
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
	port := ":8080"
	fmt.Printf("Сервер запущен на порту %s\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal("Ошибка при запуске сервера:", err)
	}
}

// getUsers выполняет запрос к базе данных и возвращает список пользователей
func getUsers(connString, query string) ([]User, error) {
	var users []User

	// Подключение к базе данных
	db, err := sql.Open("mssql", connString)
	if err != nil {
		return nil, fmt.Errorf("не удалось подключиться к базе данных: %v", err)
	}
	defer db.Close()

	// Выполнение запроса
	rows, err := db.Query("SELECT id, username, email FROM users WHERE username LIKE @query OR email LIKE @query", "%"+query+"%")
	if err != nil {
		return nil, fmt.Errorf("ошибка при выполнении запроса: %v", err)
	}
	defer rows.Close()

	// Чтение результатов
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Username, &user.Email); err != nil {
			return nil, fmt.Errorf("ошибка при чтении данных: %v", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при обработке результатов: %v", err)
	}

	return users, nil
}
