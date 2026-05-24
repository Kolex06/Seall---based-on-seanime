<p align="center">
<img src="../docs/images/seall-logo.png" alt="preview" width="70px"/>
</p>

<h2 align="center"><b>Seall Denshi</b></h2>

<p align="center">
Electron-based desktop client for Seall. Embeds server and web interface. Successor to Seall Desktop.
</p>

## Prerequisites

- Go 1.24+
- Node.js 20+ and npm

---

## Development

### Web Interface

```shell
# Working dir: ./seall-web
npm run dev:denshi
```
 
### Sidecar

1. Build the server

	```shell
	# Working dir: .
 
	# Windows
	go build -o seall.exe -trimpath -ldflags="-s -w" -tags=nosystray
 
	# Linux, macOS
	go build -o seall -trimpath -ldflags="-s -w"
	```
 
2. Move the binary to `./seall-denshi/binaries`

3. Rename the binary:

   - For Windows: `seall-server-windows.exe`
   - For macOS/Intel: `seall-server-darwin-amd64`
   - For macOS/ARM: `seall-server-darwin-arm64`
   - For Linux/x86_64: `seall-server-linux-amd64`
   - For Linux/ARM64: `seall-server-linux-arm64`

### Electron

1. Setup

	```shell
	# Working dir: ./seall-denshi
	npm install
	```

2. Run

    `TEST_DATADIR` can be used in development mode, it should point to a dummy data directory for testing purposes.

    ```shell
    # Working dir: ./seall-desktop
    TEST_DATADIR="/path/to/data/dir" npm run dev
   ```

---

## Build

### Web Interface
   
```shell
# Working dir: ./seall-web
npm run build
npm run build:denshi
```

Move the output `./seall-web/out` to `./web`
Move the output `./seall-web/out-denshi` to `./seall-denshi/web-denshi`

```shell
# UNIX command
mv ./seall-web/out ./web
mv ./seall-web/out-denshi ./seall-denshi/web-denshi
```

### Sidecar

1. Build the server

	```shell
	# Working dir: .
 
	# Windows
	go build -o seall.exe -trimpath -ldflags="-s -w" -tags=nosystray
 
	# Linux, macOS
	go build -o seall -trimpath -ldflags="-s -w"
	```
 
2. Move the binary to `./seall-denshi/binaries`

3. Rename the binary:

   - For Windows: `seall-server-windows.exe`
   - For macOS/Intel: `seall-server-darwin-amd64`
   - For macOS/ARM: `seall-server-darwin-arm64`
   - For Linux/x86_64: `seall-server-linux-amd64`
   - For Linux/ARM64: `seall-server-linux-arm64`

### Electron

To build the desktop client for all platforms:

```
npm run build
```

To build for specific platforms:

```
npm run build:mac
npm run build:win
npm run build:linux
```

Output is in `./seall-denshi/dist/...`
