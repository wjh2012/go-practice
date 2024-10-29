package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func main() {
	app := fiber.New()

	app.Get("/", hello)
	app.Post("/upload", uploadData)
	app.Post("/upload2", uploadHandler)

	app.Listen(":1010")
}

// Handler
func hello(c *fiber.Ctx) error {
	return c.SendString("Hello, World!")
}

func uploadData(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}

	if _, err := os.Stat("./uploads"); os.IsNotExist(err) {
		err := os.Mkdir("./uploads", os.ModePerm)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to create directory")
		}
	}

	destination := fmt.Sprintf("./uploads/%s", file.Filename)
	if err := c.SaveFile(file, destination); err != nil {
		return err
	}
	return c.SendString("File uploaded successfully with description")
}

// Handler
func uploadHandler(c *fiber.Ctx) error {
	// 파일 받기
	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Failed to upload image")
	}

	// 디렉토리 확인 및 생성
	if _, err := os.Stat("./uploads"); os.IsNotExist(err) {
		err := os.Mkdir("./uploads", os.ModePerm)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to create directory")
		}
	}

	// 파일 저장
	savePath := "./uploads/" + file.Filename
	err = c.SaveFile(file, savePath)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error saving file")
	}

	// 문자열 필드 받기
	description := c.FormValue("description")

	// JSON 데이터 생성
	data := map[string]string{"description": description}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to encode JSON")
	}

	// JSON 파일 저장
	jsonPath := "./uploads/description.json"
	err = os.WriteFile(jsonPath, jsonData, 0644)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to save JSON file")
	}

	fmt.Println("Description:", description)
	return c.SendString("File uploaded successfully with description: " + description)
}
