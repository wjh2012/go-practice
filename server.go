package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Get("/", hello)
	app.Post("/upload", uploadData)
	app.Post("/upload2", uploadHandler)

	app.Listen(":3000")
}

// Handler
func hello(c *fiber.Ctx) error {
	return c.SendString("Hello, World!")
}

func uploadData(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		// Handle error
		return err
	}
	destination := fmt.Sprintf("./uploads/%s", file.Filename)
	if err := c.SaveFile(file, destination); err != nil {
		// Handle error
		return err
	}
	return c.SendString("File uploaded successfully with description")

}

// Handler
func uploadHandler(c *fiber.Ctx) error {

	// 파일 받기 (예: 파일 필드명은 "image")
	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Failed to upload image")
	}

	// 이미지 파일을 바로 지정한 경로에 저장 (예: ./uploads/ 디렉토리)
	savePath := "./uploads/" + file.Filename
	err = c.SaveFile(file, savePath)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error saving file")
	}

	// 문자열 필드 받기 (예: 텍스트 필드명은 "description")
	description := c.FormValue("description")
	fmt.Println("Description:", description)

	return c.SendString("File uploaded successfully with description: " + description)
}
