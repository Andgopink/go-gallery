# Gallery - Web Image Gallery

A simple web-based image gallery built with Go using only the standard library.

## Features

- Upload images (supported formats: .jpg, .jpeg, .png, .webp)
- Display uploaded images in a responsive grid
- Delete images
- Modern, clean UI with gradient background

## Requirements

- Go 1.21 or later

## Running the Server

1. Make sure Go is installed and on your PATH
2. Run the following command:

```bash
cd E:\gallery
go run main.go
```

The server will start on http://localhost:8080

## API Endpoints

### GET /
Serves the main gallery page with HTML template.

### GET /api/images
Returns a JSON array of uploaded image filenames.

```json
[
  "1234567890_image.jpg",
  "7890123456_photo.png"
]
```

### POST /api/images/upload
Uploads an image via multipart/form-data with a 'file' field.

Returns 200 OK on success.

### GET /uploads/
Serves uploaded files.

### DELETE /uploads/{filename}
Deletes an uploaded image.

## File Structure

```
E:\gallery/
├── main.go          # Main application code
├── go.mod           # Go module file
├── templates/
│   └── index.html   # HTML template for the gallery
├── uploads/         # Directory for uploaded images
└── README.md        # This file
```

## Notes

- Images are automatically renamed with a timestamp to avoid conflicts
- The server is designed to work with the standard library only (net/http, html/template)