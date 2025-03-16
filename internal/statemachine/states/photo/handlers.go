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
		if update.Message.Photo != nil && len(update.Message.Photo) > 0 {
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
	if url == "" {
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

	// Валидация формата
	if !state.validatePhoto(resp, url, chatID) {
		return
	}

	// Сохранение файла
	uniqueFileName, err := state.savePhoto(resp, chatID)
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

func (state *PhotoState) validatePhoto(resp *http.Response, url string, chatID int64) bool {
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

	// Теперь валидируем по магическим числам
	// Читаем первые 2 байта для проверки магических чисел JPEG
	buffer := make([]byte, 2)
	_, err := io.ReadFull(resp.Body, buffer)

	log.Debug("validatePhoto", "url", url, "buffer", buffer)

	if err != nil {
		log.Error("Error reading file bytes", "url", url, "err", err)

		manul.SendMessage(chatID, "Ошибка: не удалось прочитать файл 🤕")
		return false
	}

	// JPEG начинается с FF D8
	isJPEG := bytes.Equal(buffer, []byte{0xFF, 0xD8})
	log.Debug("validatePhoto", "url", url, "bytes.Equal(buffer, []byte{0xFF, 0xD8})", isJPEG)

	if !isJPEG {
		log.Error("Incorrect file format", "url", url, "first_bytes", buffer)

		manul.SendMessage(chatID, "Ошибка: формат файла должен быть JPEG 🤓")
		return false
	}

	return true
}

func (state *PhotoState) savePhoto(resp *http.Response, chatID int64) (string, error) {
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
	out, err := os.Create(savePath)
	if err != nil {
		log.Error("Error when creating a file", "err", err)

		manul.SendMessage(chatID, "Внутренняя ошибка: не удалось сохранить файл 🤕")
		return "", err
	}
	defer out.Close()

	// Копируем содержимое файла
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		log.Error("Error when saving a file", "err", err)

		manul.SendMessage(chatID, "Внутренняя ошибка: не удалось сохранить файл 🤕")
		// Удаляем файл, если копирование не удалось
		if removeErr := os.Remove(savePath); removeErr != nil {
			log.Error("Error deleting a file after a failed save", "err", removeErr)
		}
		return "", err
	}

	return uniqueFileName, nil
}
