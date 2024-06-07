FROM golang:1.22.4-alpine3.20 AS builder

WORKDIR /app

COPY . .

# RUN go mod init encoder
RUN go mod tidy
RUN go build -o video_encoder

# Final stage
FROM alpine:latest


# Install FFmpeg
RUN apk add --no-cache ffmpeg
# Copy the binary and ffmpeg from the builder stage
# COPY --from=builder /usr/bin/ffmpeg /usr/bin/ffmpeg
COPY --from=builder /app/video_encoder /video_encoder

# Create necessary directories

WORKDIR /app

RUN mkdir /input /output

# Expose port 8080
EXPOSE 8080

ENTRYPOINT ["/video_encoder"]
