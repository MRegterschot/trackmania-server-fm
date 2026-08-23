# Trackmania Server File Manager

A lightweight HTTP API for managing the `UserData` directory of a Trackmania dedicated server (maps, match settings, plugin scripts, and general file operations).

## Deploying with Docker

The image is available at [`marijnregterschot/trackmania-server-fm`](https://hub.docker.com/r/marijnregterschot/trackmania-server-fm) on Docker Hub:

```bash
docker pull marijnregterschot/trackmania-server-fm:latest
```

### 1. Configure environment variables

| Variable            | Description                                      | Default          |
| ------------------- | ------------------------------------------------- | ---------------- |
| `FM_PORT`           | Port the API listens on                            | `3300`           |
| `FM_USERDATA_PATH`  | Path (inside the container) to the `UserData` dir  | `/app/UserData`  |
| `FM_LOG_LEVEL`      | Log level: `DEBUG`, `INFO`, `WARN`, `ERROR`        | `INFO`           |
| `FM_PASSWORD`       | Password required to authenticate API requests. Leave empty to disable auth | *(empty)* |

Copy the template if you want to run outside Docker or keep values in a file:

```bash
cp .env.template .env
```

### 2. Run the container

Mount your Trackmania dedicated server's `UserData` directory into the container at the path set by `FM_USERDATA_PATH` (default `/app/UserData`), and expose the API port:

```bash
docker run -d \
  --name trackmania-server-fm \
  -p 3300:3300 \
  -e FM_PASSWORD=changeme \
  -v /path/to/your/UserData:/app/UserData \
  marijnregterschot/trackmania-server-fm:latest
```

The API is now available at `http://localhost:3300`.

### Using docker-compose

```yaml
services:
  trackmania-server-fm:
    image: marijnregterschot/trackmania-server-fm:latest
    container_name: trackmania-server-fm
    restart: unless-stopped
    ports:
      - "3300:3300"
    environment:
      FM_PORT: 3300
      FM_USERDATA_PATH: /app/UserData
      FM_LOG_LEVEL: INFO
      FM_PASSWORD: changeme # if exposed to the internet, make sure to set a strong password
    volumes:
      - /path/to/your/UserData:/app/UserData
```

Start it with:

```bash
docker compose up -d
```

### Authentication

If `FM_PASSWORD` is set, requests must include it via the auth middleware (e.g. as a bearer token/header, depending on your client setup). Leaving `FM_PASSWORD` empty disables authentication entirely — only do this on a trusted network.
