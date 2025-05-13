FROM golang:1.24.2

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

CMD ["go", "run", "/app/cmd/NoteVault/main.go"]

EXPOSE 8085