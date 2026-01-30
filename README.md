# 🍋 Lemon

(LLM generated README apart from the last section.)

**Lemon** is a tiny Go-based launcher for local web apps (HTML / CSS / WASM).

It scans an `apps/` directory, exposes a small JSON API, and serves each app
securely — no frameworks, no build steps, no magic.

Perfect for:
- WASM experiments
- Internal tools
- Small self-contained web apps
- Local dashboards

## Features

- Simple Go HTTP server
- Secure file serving (directory traversal protection)
- `/api/apps` endpoint for app discovery
- Modern dark UI (said ChatGPT)
- Zero dependencies (Apart from Go obviously)

## System Structure

```text
.
├── main.go        # Go server
├── static/index.html     # Launcher UI
├── static/style.css      # Launcher CSS
└── apps/          # Local apps
    └── my-app/
        ├── index.html / style.css / main.wasm 
        └── lemon.json
```

## App Format
Each app lives in its own subdirectory inside apps/ and must include
a lemon.json file.

Example:
```json
{
  "title": "Example App",
  "description": "A small demo app"
}
```
The app is then available at:

http://localhost:8080/apps/my-app/

## Running Lemon
```bash
go run main.go
```
Then open:

http://localhost:8080

## Security Notes
- Directory listing is disabled
- Only files inside apps/ can be served
- (This is not LLM writing) https support is not implemented. I'm using Traefik in front of this which takes care of TLS cert and https stuff.

## Other Notes
Lemon is intentionally minimal.

I'm tired of the npm and web-related libraries having too many CVE. Hence, I've been only creating web app for myself just with simple html file and WASM. This is basically the somewhat extream way of achieving "Web App Security". Maybe. I don't know. I'm only using this on my home server so who cares really.
