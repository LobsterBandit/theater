# theater

## Build

```sh
DOCKER_BUILDKIT=1 docker build -t theater .
```

## Run

```sh
docker run --rm -d \
    -e BRIDGE_IP=your.hue.bridge.IP \
    -e BRIDGE_USER=YourRegisteredHueAppUser \
    -e PLEX_PLAYER=YourPlexPlayerTitle \
    -e PLEX_USER=YourPlexAccountTitle \
    # PORT defaults to 9501 if not set
    -e PORT=9501 \
    -p 9501:9501 \
    -v /your/config/dir:/config \
    --name theater \
    theater
```
