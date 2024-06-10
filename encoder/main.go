package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func encodeVideo(inputPath, outputPath string) error {
	if err := os.MkdirAll(outputPath, 0755); err != nil {
		return err
	}

cmd := exec.Command("ffmpeg", "-i", inputPath,
	"-map", "0:v:0", "-map", "0:a:0",
	"-map", "0:v:0", "-map", "0:a:0",
	"-map", "0:v:0", "-map", "0:a:0",
	"-map", "0:v:0", "-map", "0:a:0",
	"-c:v", "libx264", "-crf", "22", "-c:a", "aac", "-ar", "48000",
	"-filter:v:0", "scale=w=480:h=360", "-maxrate:v:0", "600k", "-b:a:0", "64k",
	"-filter:v:1", "scale=w=640:h=480", "-maxrate:v:1", "900k", "-b:a:1", "128k",
	"-filter:v:2", "scale=w=1280:h=720", "-maxrate:v:2", "1500k", "-b:a:2", "128k",
	"-filter:v:3", "scale=w=1920:h=1080", "-maxrate:v:3", "3000k", "-b:a:3", "192k",
	"-preset", "slow", "-hls_list_size", "0", "-threads", "0", "-f", "hls",
	"-var_stream_map", "v:0,a:0 v:1,a:1 v:2,a:2 v:3,a:3",
	"-master_pl_name", "master.m3u8",
	outputPath+"/video-%v.m3u8",
)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("video")
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to read video file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	filename := strings.ReplaceAll(header.Filename, " ", "_")
	inputPath := filepath.Join("/input", header.Filename)
	outputPath := filepath.Join("/usr/share/nginx/html/hls", strings.ToLower(filename))

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
	http.HandleFunc("POST /upload", uploadHandler)
	fmt.Println("Server started on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
