# theater

## Build

```sh
DOCKER_BUILDKIT=1 docker build -t theater .
```

## Run

```sh
docker run --rm -d -p 9501:9501 -v /your/config/dir:/config --name theater theater
```
