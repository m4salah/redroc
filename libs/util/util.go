package util

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"log/slog"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/nfnt/resize"
	"github.com/spf13/viper"
	"google.golang.org/api/idtoken"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	grpcMetadata "google.golang.org/grpc/metadata"
)

func bindEnvs(iface interface{}, parts ...string) {
	ifv := reflect.ValueOf(iface)
	ift := reflect.TypeOf(iface)
	for i := 0; i < ift.NumField(); i++ {
		v := ifv.Field(i)
		t := ift.Field(i)
		tv, ok := t.Tag.Lookup("mapstructure")
		if !ok {
			continue
		}
		switch v.Kind() {
		case reflect.Struct:
			bindEnvs(v.Interface(), append(parts, tv)...)
		default:
			viper.BindEnv(strings.Join(append(parts, tv), "."))
		}
	}
}

// Builds config - error handling omitted fore brevity
func LoadConfig[Config any](c Config) Config {
	// load from .env file
	viper.AddConfigPath(".")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	viper.ReadInConfig()

	// load from env variables
	bindEnvs(c)
	viper.Unmarshal(&c)
	return c
}

func InitializeSlog(env, release string) {
	logHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}).WithAttrs([]slog.Attr{
		slog.Group("environment", slog.String("release", release), slog.String("env", env)),
	})
	logger := slog.New(logHandler)
	slog.SetDefault(logger)
}

func GetPhoto(r *http.Request) ([]byte, string, error) {
	image, header, err := r.FormFile("file")
	if err != nil {
		return nil, "", fmt.Errorf("failed to get file: %v", err)
	}
	defer image.Close()

	var b bytes.Buffer
	_, err = io.Copy(&b, image)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read file: %v", err)
	}

	return b.Bytes(), header.Filename, nil
}

// GetTags function tries to read and decode photo tags from the request.
func GetTags(r *http.Request) ([]string, error) {
	var tags []string

	data := r.FormValue("hashtags")
	if len(data) == 0 {
		return tags, nil
	}

	err := json.Unmarshal([]byte(data), &tags)
	if err != nil {
		return nil, err
	}

	return tags, nil
}

func MakeThumbnail(photo []byte, width, height uint) ([]byte, error) {
	r := bytes.NewReader(photo)
	img, format, err := image.Decode(r)
	if err != nil {
		return nil, fmt.Errorf("photo decode failed: %v", err)
	}

	thumb := resize.Thumbnail(width, height, img, resize.NearestNeighbor)

	var b bytes.Buffer
	if format == "jpeg" {
		err = jpeg.Encode(&b, thumb, nil)
	} else if format == "png" {
		err = png.Encode(&b, thumb)
	} else if format == "gif" {
		err = gif.Encode(&b, thumb, nil)
	} else {
		err = fmt.Errorf("unsuported image format: %s", format)
	}
	if err != nil {
		return nil, fmt.Errorf("thumbnail encode failed: %v", err)
	}

	return b.Bytes(), nil
}

func GetAuthContext(ctx context.Context, audience string, skipAuth bool) (context.Context, error) {
	if skipAuth {
		return ctx, nil
	}
	// Create an identity token.
	// With a global TokenSource tokens would be reused and auto-refreshed at need.
	// A given TokenSource is specific to the audience.
	tokenSource, err := idtoken.NewTokenSource(ctx, audience)
	if err != nil {
		return nil, fmt.Errorf("idtoken.NewTokenSource: %w", err)
	}
	token, err := tokenSource.Token()
	if err != nil {
		return nil, fmt.Errorf("TokenSource.Token: %w", err)
	}

	// Add token to gRPC Request.
	ctx = grpcMetadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+token.AccessToken)
	return ctx, nil
}

func ExtractServiceURL(addr string) string {
	return "https://" + strings.Split(addr, ":")[0]
}

// CreateTransportCredentials creates a new TLS credentials instance with the system root CA pool.
//
// This is used to create a secure connection to the server.
func CreateTransportCredentials(skipAuth bool) (credentials.TransportCredentials, error) {
	if skipAuth {
		return insecure.NewCredentials(), nil
	}
	systemRoots, err := x509.SystemCertPool()
	if err != nil {
		return nil, fmt.Errorf("failed to load system root CA cert pool")
	}
	creds := credentials.NewTLS(&tls.Config{
		RootCAs: systemRoots,
	})
	return creds, nil
}

// encrypt data using AES algorithm
// this code from https://github.com/purnaresa/bulwark/blob/master/encryption/encrpytion.go
func EncryptAES(plainData, secret []byte) ([]byte, error) {
	block, err := aes.NewCipher(secret)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return nil, err
	}
	return gcm.Seal(
		nonce,
		nonce,
		plainData,
		nil), nil
}

// decrypt data using AES algorithm
// this code from https://github.com/purnaresa/bulwark/blob/master/encryption/encrpytion.go
func DecryptAES(cipherData, secret []byte) ([]byte, error) {
	block, err := aes.NewCipher(secret)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := cipherData[:nonceSize], cipherData[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}

// GetStringOrDefault returns the value of the environment variable named by the key or devaultV if key not found.
func GetStringOrDefault(name, defaultV string) string {
	v, ok := os.LookupEnv(name)
	if !ok {
		return defaultV
	}
	return v
}

func LoadStructFromEnv[Config any](v *Config) (*Config, error) {
	fmt.Println(v)
	val := reflect.ValueOf(v).Elem()
	fmt.Println(val)
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		envKey := fieldType.Tag.Get("env")

		if envKey != "" {
			envValue := os.Getenv(envKey)
			if envValue == "" {
				continue
			}

			switch field.Kind() {
			case reflect.String:
				field.SetString(envValue)
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				intValue, err := strconv.Atoi(envValue)
				if err != nil {
					return nil, err
				}
				field.SetInt(int64(intValue))
				// Add more cases for other types as needed
			}
		}
	}

	return v, nil
}
