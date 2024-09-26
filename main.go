package main

import (
	"encoding/json"
	"github.com/google/uuid"
	"io"
	"net/http"
	"os"
	"os/exec"
)

func main() {
	http.HandleFunc("/ping", defaultRoute)
	http.HandleFunc("/api/ping", defaultRoute)
	http.HandleFunc("/api/ocr", ocr)
	server := &http.Server{Addr: ":8080", Handler: nil}
	err := server.ListenAndServe()
	if err != nil {
		return
	}
}

func ocr(response http.ResponseWriter, request *http.Request) {
	img := request.URL.Query().Get("img")
	if img == "" {
		res(response, "101", "param error", "")
		return
	}
	id := uuid.New().String()
	err := downloadFile(id, img)
	if err != nil {
		res(response, "102", err.Error(), "")
		return
	}
	defer os.Remove(id + ".png")
	cmd := exec.Command("tesseract", id+".png", id, "-l", "chi_sim")
	err = cmd.Run()
	if err != nil {
		res(response, "103", err.Error(), "")
		return
	}
	content, err := os.ReadFile(id + ".txt")
	if err != nil {
		res(response, "104", err.Error(), "")
		return
	}
	defer os.Remove(id + ".txt")
	res(response, "0", "", (string)(content))
}

func res(response http.ResponseWriter, code, msg, data string) {
	response.Header().Set("Content-Type", "application/json;charset=UTF-8")
	strByte, err := json.Marshal(Response{
		Code: code,
		Msg:  msg,
		Data: data,
	})
	if err != nil {
		response.Write([]byte("{\"code\":\"100\"}"))
	} else {
		response.Write(strByte)
	}
}

func downloadFile(id, img string) error {
	out, err := os.Create(id + ".png")
	if err != nil {
		return err
	}
	defer out.Close()
	resp, err := http.Get(img)
	if err != nil {
		return err
	}
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	return nil
}

func defaultRoute(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json;charset=UTF-8")
	response.Write([]byte("{\"code\":\"0\"}"))
}
