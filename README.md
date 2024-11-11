# S3 Downloader

## Development

1. Create config file.

Config file for this application located in `/etc/s3-downloader` or `$HOME/.s3-downloader` or the same location with this app. Below is sample of a config file.

Using json

```json
{
  "s3_connection": "http://minio:minio123@localhost:9000/go-binary"
}
```

or using environment variable

```bash
export S3_CONNECTION=http://minio:minio123@localhost:9000/go-binary
```

2. Run makefile for getting golang dependencies

```bash
make tidy
```

3. Run minio with docker container

```bash
make compose-up
```

4. Run application

```bash
go run main.go -h
go run main.go -object dev/s3-downloader -path ./s3-downloader
go run main.go -object dev/s3-downloader -path ./s3-downloader -up
```

5. Build application

```bash
make build
```

## Usage

After we build the application, we can use it using command-line.

- Help Command

```bash
./s3-downloader -h
```

- Download object command

```bash
go run main.go -object dev/s3-downloader -path ./s3-downloader
```

- Upload object command

```bash
go run main.go -object dev/s3-downloader -path ./s3-downloader -up
```
