# AutoDeploy

**AutoDeploy** is a Go-based HTTP service designed for seamless file deployment. The service accepts a ZIP file, extracts it to a target directory, generates an Nginx configuration to serve the content, and reloads Nginx automatically to apply the new configuration. This tool is ideal for automating the deployment of static sites or other file-based web applications.
( it just for deploy static sites that has index.html )
## Features

- **File Upload**: Accepts ZIP files through an HTTP POST request.
- **Automatic Extraction**: Unpacks the ZIP file to a designated target directory.
- **Nginx Configuration**: Creates a custom Nginx configuration to serve the extracted content.
- **Nginx Reload**: Reloads Nginx to apply the new configuration, making the files immediately accessible.

## Requirements

- **Go** 1.13 or higher
- **Nginx** installed on the server
- Permissions to modify Nginx configurations and reload the Nginx service

## Installation

1. **Clone the repository**:
   ```bash
   git clone https://github.com/mohammadtafakori01/AutoDeploy.git
   cd AutoDeploy
   ```

2. **Build the application**:
   ```bash
   go build
   ```

3. **Run the application**:
   ```bash
   ./AutoDeploy
   ```

   The server will start on `http://localhost:8080`.

## Usage

### Uploading a ZIP file

To deploy a ZIP file, send a POST request to the `/upload` endpoint:

- **URL**: `http://localhost:8080/upload`
- **Method**: `POST`
- **Form Data**:
  - `metadata`: JSON string containing the `targetDir`, e.g., `{"targetDir": "./myApp"}`.
  - `zipfile`: The ZIP file to upload.

### Example cURL Command

```bash
curl -X POST http://localhost:8080/upload \
  -F "metadata={\"targetDir\": \"./myApp\"}" \
  -F "zipfile=@/path/to/yourfile.zip"
```

If successful, the response will include the URL where the files are served.

## Directory Structure

The project is organized as follows:

```
AutoDeploy/
├── main.go                # Main application file
├── extract/               # Extraction logic for ZIP files
│   └── extract.go
├── nginxHandler/          # Nginx configuration and reload handler
│   └── nginxHandler.go
└── upload/                # File upload handling
    └── upload.go
```

## Configuration Files

AutoDeploy creates configuration files in `/etc/nginx/sites-available` and `/etc/nginx/sites-enabled` with symlinks. Nginx must be installed and configured to serve these sites.

## Permissions

Ensure that the application has permission to:
- Create and modify files in `/etc/nginx/sites-available/` and `/etc/nginx/sites-enabled/`
- Reload the Nginx service

## Security Considerations

- **Input Validation**: The service includes basic validation. You may add additional checks as necessary.
- **Permissions**: Running this service with elevated privileges should be handled carefully to prevent security risks.

## License

This project is licensed under the MIT License.
