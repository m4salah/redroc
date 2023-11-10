package types

type Config struct {
	DownloadBackendAddr string `env:"DOWNLOAD_BACKEND_ADDR,notEmpty"`
	UploadBackendAddr   string `env:"UPLOAD_BACKEND_ADDR,notEmpty"`
	SearchBackendAddr   string `env:"SEARCH_BACKEND_ADDR,notEmpty"`
	Port                int    `env:"PORT,notEmpty"`
	Host                string `env:"HOST"`
	ProjectID           string `env:"PROJECT_ID,notEmpty"`
	TopicID             string `env:"TOPIC_ID,notEmpty"`
}
