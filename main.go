package main

import (
	"net/http"
	"strconv" //  - Конвертация типов данных
	"github.com/labstack/echo/v4"
)

// JSON - Стандарт передачи данных в RESTAPI
// Context содержит в себе и Writer и Responce

// c.JSON(http.StatosOK, Message) // Эта функция отдает вписанные данные клиенту, а Ошибку возвращает серверу.
// Содержит ошибку внутри себя.

// Парсить (от англ. parse) — значит разбирать данные из строки (или потока) в удобный для программы вид.
// Распарсить (превратить в структуру).

// c.Bind(&..) - Возвращает ошибку. Является приемником JSON структуры от Фронтенда и декодирует ее в структуру Go по указателю.
// c.Bind(&..) – парсит JSON из тела запроса в структуру.

// CRUD - create, read, update, delete.

//=============================================================================================================

// Структура - Запрос 
type Message struct {
	ID 		int 	`json:"id"`
	Text 	string 	`json:"text"`
}

// Структура - Ответ // Что вернет сервер клиенту
type Responce struct {
	Status 	string 	`json:"status"`
	Message string 	`json:"message"`
}

// Массив структур Message
var messages = make(map[int]Message)

// ID 
var nextID int = 1

//==============================================================================================================

// Просмотр всех имеющихся JSON Структур
func GetJSON(c echo.Context) error { 
	var msgSlice []Message

	for _, msg := range messages {
		msgSlice = append(msgSlice, msg)
	}

	// Статус-Код и структуру
	return c.JSON(http.StatusOK, &msgSlice) // Клиент получит массив JSON-объектов 
}

// Создание JSON структуры по данным, полученным фронтом.
func PostJSON(c echo.Context) error {
	var mes Message // Объект класса Message

	// Проверка на ошибку в c.Bind()
	if err := c.Bind(&mes); err != nil { // Если ошибка есть, ТО ..
		return c.JSON(http.StatusBadRequest, Responce{
			Status: "Error",
			Message: "Could not add the Message",
		})
	}

	mes.ID = nextID // Заполнение структуры
	nextID++

	messages[mes.ID] = mes // Добавили в мапу Структуру c текстом из пришедшего JSON от фронтенда

	return c.JSON(http.StatusOK, Responce{
		Status: "Success",
		Message: mes.Text,
	})
}

// Обновление данных
func PatchJSON(c echo.Context) error {
	idStr := c.Param("id") // Параметр из URL // приходит в стринге
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Responce{
			Status: "Error",
			Message: "Bad Parameter URL",
		})
	}

	// ЗАполняем стурктуру 
	var UpdateMes Message
	if err := c.Bind(&UpdateMes); err != nil { // Если ошибка есть, ТО ..
		return c.JSON(http.StatusBadRequest, Responce{
			Status: "Error",
			Message: "Could not UPDATE the MESSAGE",
		})
	}

	// Проверяет, есть ли что-то по ID
	if _, exist := messages[id]; !exist {
		return c.JSON(http.StatusBadRequest, Responce{
			Status: "Error",
			Message: "MESSAGE doesn`t exists",
		})
	}

	UpdateMes.ID = id 	// Заполнили структуру
	messages[id] = UpdateMes // Отдали структуру в мапу

	return c.JSON(http.StatusOK, Responce{
		Status: "Success",
		Message: "Message was updated",
	})


}

func DeleteJSON(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Responce{
			Status: "Error",
			Message: "Bad Parameter URL",
		})
	}

	if _, exist := messages[id]; !exist {
		return c.JSON(http.StatusBadRequest, Responce{
			Status: "Error",
			Message: "MESSAGE doesn`t exists",
		})
	}

	delete(messages, id) // Удаление

	return c.JSON(http.StatusOK, Responce{
		Status: "Success",
		Message: "Message was Deleted",
	})

}

func main() {
	e := echo.New()

	// CRUD
	e.GET("/get", GetJSON)
	e.POST("/post", PostJSON)
	e.PATCH("/patch/:id", PatchJSON)
	e.DELETE("/delete/:id", DeleteJSON)

	e.Start(":8080")
}
