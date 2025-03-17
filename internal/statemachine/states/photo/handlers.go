package photo

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/epistax1s/photo-manul/internal/log"
	. "github.com/epistax1s/photo-manul/internal/statemachine/core"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
)

func (state *PhotoState) Init(update *tgbotapi.Update) {
	manul := state.server.Manul
	chatID := update.FromChat().ID

	employee := state.context.Employee

	msg := fmt.Sprintf("👨🏻‍💼 Сотрудник:\n\n"+
		"ID = %d\n"+
		"ФИО = %s\n", employee.EmployeeID, employee.EmployeeName)

	manul.SendMessage(chatID, msg)
	manul.SendMessage(chatID, "Пришлите фото 📸")
}

func (state *PhotoState) Handle(update *tgbotapi.Update) {
	manul := state.server.Manul
	employeeService := state.server.EmployeeService

	chatID := update.FromChat().ID

	// Извлекаем ID файла
	var fileID string
	if update.Message != nil {
		if len(update.Message.Photo) > 0 {
			// Если это фото, берём последнюю (самую большую) версию
			photo := update.Message.Photo[len(update.Message.Photo)-1]
			fileID = photo.FileID
		} else if update.Message.Document != nil {
			// Если это документ, берём его fileID
			fileID = update.Message.Document.FileID
		} else {
			manul.SendMessage(chatID, "Ошибка: необходимо прислать фото в формате JPEG (JPG) (можно файлом) 🚨")
			return
		}
	} else {
		manul.SendMessage(chatID, "Ошибка: необходимо прислать фото в формате JPEG (JPG) (можно файлом) 🚨")
		return
	}

	// Получаем прямой URL
	url, err := manul.GetFileDirectURL(fileID)
	if err != nil {
		manul.SendMessage(chatID, "Ошибка: не удалось получить URL файла 🤕")
		return
	}

	// Скачиваем файл
	resp, err := http.Get(url)
	if err != nil {
		log.Error("Downloading a file", "url", url, "err", err)

		manul.SendMessage(chatID, "Ошибка: не удалось скачать файл 🤕")
		return
	}
	defer resp.Body.Close()

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		log.Error("Downloading a file", "url", url, "code", resp.StatusCode)

		manul.SendMessage(chatID, "Ошибка: не удалось скачать файл 🤕")
		return
	}

	// Читаем всё содержимое в буфер
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error("Error reading response body", "err", err)
		manul.SendMessage(chatID, "Ошибка: не удалось прочитать файл 🤕")
		return
	}

	// Валидация формата
	if !state.validatePhoto(data, url, chatID) {
		return
	}

	// Сохранение файла
	uniqueFileName, err := state.savePhoto(data, chatID)
	if err != nil {
		return
	}

	employee := state.context.Employee
	employee.ImagePath = uniqueFileName

	if err := employeeService.UpdateEmployee(employee); err != nil {
		log.Error("Error when trying to update an employee in the database", "employee", employee, "err", err)

		manul.SendMessage(chatID, "Ошибка: не удалось обновить сотрудника 🤕")
		return
	}

	manul.SendMessage(chatID, "Фотография была успешно загружена ✅")

	state.stateMachine.
		Set(Idle, chatID, &StateContext{}).
		Init(update)
}

func (state *PhotoState) validatePhoto(data []byte, url string, chatID int64) bool {
	manul := state.server.Manul

	// Валидируем по расширению
	// Извлекаем расширение из URL
	extension := filepath.Ext(url)

	log.Debug("validatePhoto", "url", url, "extension", extension)

	if extension == "" {
		log.Error("No file extension found", "url", url)

		manul.SendMessage(chatID, "Ошибка: с файлом что-то не так 🤕")
		return false
	}

	// Приводим расширение к нижнему регистру для проверки
	lowerExt := strings.ToLower(extension)

	// Проверяем, является ли расширение .jpg или .jpeg
	if lowerExt != ".jpg" && lowerExt != ".jpeg" {
		log.Error("Incorrect file extension", "url", url, "extension", extension)

		manul.SendMessage(chatID, "Ошибка: файл должен иметь расширение .jpg или .jpeg 🤓")
		return false
	}

	// Теперь валидируем по магическим числам (первые 2 байта)
	if len(data) < 2 {
		log.Error("File too small", "url", url)
		manul.SendMessage(chatID, "Ошибка: формат файла должен быть JPEG 🤓")
		return false
	}
	isJPEG := bytes.Equal(data[:2], []byte{0xFF, 0xD8})
	log.Debug("validatePhoto", "url", url, "isJPEG", isJPEG)
	if !isJPEG {
		log.Error("Incorrect file format", "url", url, "first_bytes", data[:2])
		manul.SendMessage(chatID, "Ошибка: формат файла должен быть JPEG 🤓")
		return false
	}

	return true
}

func (state *PhotoState) savePhoto(data []byte, chatID int64) (string, error) {
	manul := state.server.Manul

	// Генерируем уникальное имя файла
	uniqueFileName := uuid.New().String() + ".jpg"
	saveDir := "/app/photos/"
	savePath := filepath.Join(saveDir, uniqueFileName)

	// Создаём директорию, если её нет
	if _, err := os.Stat(saveDir); os.IsNotExist(err) {
		if err := os.MkdirAll(saveDir, 0755); err != nil {
			log.Error("Error when creating a directory", "err", err)

			manul.SendMessage(chatID, "Внутренняя ошибка: не удалось создать директорию 🤕")
			return "", err
		}
	}

	// Создаём файл
	err := os.WriteFile(savePath, data, 0644)
	if err != nil {
		log.Error("Error when saving a file", "err", err)
		manul.SendMessage(chatID, "Внутренняя ошибка: не удалось сохранить файл 🤕")
		return "", err
	}

	return uniqueFileName, nil
}
