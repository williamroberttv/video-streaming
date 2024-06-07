package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

func encodeVideo(inputPath, outputPath string) error {
	if err := os.MkdirAll(outputPath, 0755); err != nil {
		return err
	}

	cmd := exec.Command("ffmpeg", "-i", inputPath,
		"-vf", "scale=1920:1080", "-c:v", "libx264", "-b:v", "2M", "-c:a", "aac", "-strict", "experimental", "-b:a", "192k", "-f", "hls", "-hls_time", "10", "-hls_list_size", "0", "-preset", "ultrafast", outputPath+"/1080p.m3u8",
		"-vf", "scale=1280:720", "-c:v", "libx264", "-b:v", "1.5M", "-c:a", "aac", "-strict", "experimental", "-b:a", "128k", "-f", "hls", "-hls_time", "10", "-hls_list_size", "0", "-preset", "ultrafast", outputPath+"/720p.m3u8",
		"-vf", "scale=854:480", "-c:v", "libx264", "-b:v", "1M", "-c:a", "aac", "-strict", "experimental", "-b:a", "96k", "-f", "hls", "-hls_time", "10", "-hls_list_size", "0", "-preset", "ultrafast", outputPath+"/480p.m3u8",
		"-vf", "scale=640:360", "-c:v", "libx264", "-b:v", "500k", "-c:a", "aac", "-strict", "experimental", "-b:a", "64k", "-f", "hls", "-hls_time", "10", "-hls_list_size", "0", "-preset", "ultrafast", outputPath+"/360p.m3u8",
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	file, header, err := r.FormFile("video")
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to read video file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// fileNamePath, err := os.ReadDir("./input")
	// if err != nil {
	// 	log.Println(err)
	// }
	// log.Println(fileNamePath)
	inputPath := filepath.Join("/input", header.Filename)
	outputPath := filepath.Join("/output", header.Filename)

	log.Println("Saving video to", inputPath)
	log.Println("Saving encoded video to", outputPath)

		// // Check if the input directory exists
		// if _, err := os.Stat("/input"); os.IsNotExist(err) {
		// 	log.Println("/input directory does not exist")
		// 	http.Error(w, "/input directory does not exist", http.StatusInternalServerError)
		// 	return
		// }
	
		// // Check if the output directory exists
		// if _, err := os.Stat("/output"); os.IsNotExist(err) {
		// 	log.Println("/output directory does not exist")
		// 	http.Error(w, "/output directory does not exist", http.StatusInternalServerError)
		// 	return
		// }
	out, err := os.Create(inputPath)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to create input file", http.StatusInternalServerError)
		return
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to save input file", http.StatusInternalServerError)
		return
		}
		
	err = encodeVideo(inputPath, outputPath)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error encoding video", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Video encoded successfully!")
}

func main() {
	http.HandleFunc("/upload", uploadHandler)
	fmt.Println("Server started on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
