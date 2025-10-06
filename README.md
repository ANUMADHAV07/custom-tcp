# Simple TCP-based HTTP Server in Go

This project implements a simple TCP server in Go that listens on port 8080 and supports basic HTTP-like routing with gzip compression support and file creation via POST requests.

## Project Setup

### Prerequisites

- Go (Golang) installed (version 1.19+ recommended)
- Basic command line experience

### Getting Started

1. Clone/download this repository or copy the source code to your local machine.

2. Build the server binary:

```
go build -o tcp-server main.go

```

3. Run the server with a specified directory for files:

```
./tcp-server -directory /absolute/path/to/files

```

## API Endpoints for Testing

### GET `/echo/{text}`

- Returns the `{text}` as response.
- Supports gzip compression if `Accept-Encoding: gzip` header is sent.

Example:

```
curl -v -H "Accept-Encoding: gzip" http://localhost:8080/echo/hello

```

### GET `/files/{filename}`

- Returns the `{filename}` as plain text.

Example:

```
curl http://localhost:8080/files/sample.txt

```

### POST `/files/{filename}`

- Creates a file with the POST request body as content in the server directory.

Example:

```
curl -X POST --data "File content here" http://localhost:8080/files/newfile.txt

```

---

## Notes

- This server is for learning and basic testing purposes.
- Ensure the directory path passed to `-directory` exists and is writable.
- Requests without proper `Content-Length` header may cause hangs on POST.

---

## License

This project is licensed under the MIT License.
