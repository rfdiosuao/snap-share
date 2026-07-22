# Snap & Share ⚡️

[中文](README.md) | **English**

> Upload instantly, scan to download, and let temporary files expire automatically.

Snap & Share is a lightweight file-transfer service built with Go and Vue 3. It is designed for quick transfers between a computer and a phone, with a single backend port serving both the web interface and file APIs.

## Features

- Drag-and-drop uploads, including large files
- QR codes for convenient mobile downloads
- Automatic expiration after one hour or five downloads by default
- A Vue 3 frontend and a compact Go backend
- Simple single-port deployment

## Requirements

- Go 1.23 or newer
- Node.js 18 or newer with npm

## Configure the server

Edit `backend/config.json` before exposing the service:

```json
{
  "base_url": "http://localhost:8080",
  "port": ":8080"
}
```

Set `base_url` to an address that download devices can reach. Review the remaining upload limits and retention settings in the same file.

## Build and run

Build the frontend:

```bash
cd frontend
npm install
npm run build
```

Copy the generated `frontend/dist` directory to `backend/dist`, then start the Go service:

```bash
cd backend
go mod tidy
go run main.go
```

Open `http://localhost:8080` unless you changed the port.

## Security notes

- Treat uploaded files and generated links as sensitive.
- Use HTTPS and an authenticated reverse proxy for public deployments.
- Keep retention limits short and monitor available disk space.
- Do not rely on automatic deletion as a substitute for access control.

## License status

The original project documentation identifies this project as MIT licensed, but the repository currently does not include a standalone license file. Add or verify the license text before redistribution.
