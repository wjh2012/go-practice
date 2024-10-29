package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", hello)
	http.HandleFunc("/upload", uploadData)
	http.HandleFunc("/upload2", uploadHandler)

	fmt.Println("Server is running on port 1010...")
	http.ListenAndServe(":1010", nil)
}

// 기본 핸들러
func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, World!"))
}

// 디렉토리 생성 함수
func ensureUploadDir() error {
	if _, err := os.Stat("./uploads"); os.IsNotExist(err) {
		return os.Mkdir("./uploads", os.ModePerm)
	}
	return nil
}

// 파일 업로드 핸들러
func uploadData(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// 디렉토리 확인 및 생성
	if err := ensureUploadDir(); err != nil {
		http.Error(w, "Failed to create directory", http.StatusInternalServerError)
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to upload file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// 파일 저장
	destination := "./uploads/uploaded_file"
	out, err := os.Create(destination)
	if err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}
	defer out.Close()
	io.Copy(out, file)

	w.Write([]byte("File uploaded successfully"))
}

// 파일 및 JSON 저장 핸들러
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// 디렉토리 확인 및 생성
	if err := ensureUploadDir(); err != nil {
		http.Error(w, "Failed to create directory", http.StatusInternalServerError)
		return
	}

	// 파일 받기
	file, header, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Failed to upload image", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// 이미지 파일 저장
	savePath := "./uploads/" + header.Filename
	out, err := os.Create(savePath)
	if err != nil {
		http.Error(w, "Error saving file", http.StatusInternalServerError)
		return
	}
	defer out.Close()
	io.Copy(out, file)

	// 문자열 필드 받기
	description := r.FormValue("description")

	// JSON 데이터 생성
	data := map[string]string{"description": description}
	jsonData, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
		return
	}

	// JSON 파일 저장
	jsonPath := "./uploads/description.json"
	err = os.WriteFile(jsonPath, jsonData, 0644)
	if err != nil {
		http.Error(w, "Failed to save JSON file", http.StatusInternalServerError)
		return
	}

	fmt.Println("Description:", description)
	w.Write([]byte("File uploaded successfully with description: " + description))
}
