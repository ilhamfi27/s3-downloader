package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/spf13/viper"
)

type Config struct {
	S3Connection string `mapstructure:"s3_connection"`
	S3ObjectKey  string `mapstructure:"s3_object_key"`
	DownloadPath string `mapstructure:"download_path"`

	S3Endpoint  string `mapstructure:"-"`
	S3AccessKey string `mapstructure:"-"`
	S3SecretKey string `mapstructure:"-"`
	S3Bucket    string `mapstructure:"-"`
	S3Secure    bool   `mapstructure:"-"`

	Uploading bool `mapstructure:"-"`
}

// LoadConfig loads configuration from a JSON or YAML file and overrides with environment variables
func LoadConfig() (*Config, error) {
	// Set the file name and type for config
	viper.SetConfigName("config")               // File name without extension
	viper.SetConfigType("json")                 // or "json" if using JSON config file
	viper.AddConfigPath("/etc/s3-downloader")   // Path to look for the config file
	viper.AddConfigPath("$HOME/.s3-downloader") // Path to look for the config file
	viper.AddConfigPath(".")                    // Path to look for the config file

	// Enable environment variable overrides
	viper.AutomaticEnv()

	// Bind environment variables to config fields
	viper.BindEnv("s3_connection", "S3_CONNECTION")

	// Set default values (optional)
	viper.SetDefault("download_path", "./downloaded-object")

	// Read in the config file, if it exists
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("No config file found, using environment variables")
	}

	// Map the configuration to the Config struct
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("could not decode config into struct: %w", err)
	}

	parsedUrl, err := url.Parse(config.S3Connection)
	if err != nil {
		return nil, fmt.Errorf("failed to parse S3 connection URL: %w", err)
	}

	proto := parsedUrl.Scheme

	if proto == "http" {
		config.S3Endpoint = parsedUrl.Hostname() + ":" + parsedUrl.Port()
		config.S3Secure = false
	} else {
		config.S3Endpoint = parsedUrl.Hostname()
		config.S3Secure = true
	}
	config.S3AccessKey = parsedUrl.User.Username()
	config.S3SecretKey, _ = parsedUrl.User.Password()
	config.S3Bucket = strings.Trim(parsedUrl.Path, "/")

	objectKey := flag.String("object", "", "S3 object key to download")
	downloadPath := flag.String("path", config.DownloadPath, "Path to download the object")
	uploading := flag.Bool("up", false, "Flag to indicate that the program is being run by the downloading service")
	flag.Parse()

	return &Config{
		S3Endpoint:   config.S3Endpoint,
		S3AccessKey:  config.S3AccessKey,
		S3SecretKey:  config.S3SecretKey,
		S3Bucket:     config.S3Bucket,
		S3Connection: config.S3Connection,
		S3Secure:     config.S3Secure,
		S3ObjectKey:  *objectKey,
		DownloadPath: *downloadPath,
		Uploading:    *uploading,
	}, nil
}

func DownloadObject(config *Config) error {
	minioClient, err := minio.New(config.S3Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.S3AccessKey, config.S3SecretKey, ""),
		Secure: config.S3Secure,
	})
	if err != nil {
		return fmt.Errorf("failed to initialize MinIO client: %w", err)
	}

	ctx := context.Background()
	if config.Uploading {
		_, err := minioClient.FPutObject(ctx, config.S3Bucket, config.S3ObjectKey, config.DownloadPath, minio.PutObjectOptions{})
		if err != nil {
			return fmt.Errorf("failed to upload object: %w", err)
		}
		fmt.Printf("Successfully uploaded %s to %s\n", config.DownloadPath, config.S3ObjectKey)
		return nil
	}

	err = minioClient.FGetObject(ctx, config.S3Bucket, config.S3ObjectKey, config.DownloadPath, minio.GetObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to download object: %w", err)
	}

	fmt.Printf("Successfully downloaded %s to %s\n", config.S3ObjectKey, config.DownloadPath)
	return nil
}

func main() {
	config, err := LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	err = DownloadObject(config)
	if err != nil {
		log.Fatalf("Failed to download object: %v", err)
	}
}
