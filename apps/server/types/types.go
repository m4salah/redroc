package types

import "net"

type Config struct {
	DownloadBackendAddr net.Addr `env:"DOWNLOAD_BACKEND_ADDR,notEmpty"`
	UploadBackendAddr   net.Addr `env:"UPLOAD_BACKEND_ADDR,notEmpty"`
	SearchBackendAddr   net.Addr `env:"SEARCH_BACKEND_ADDR,notEmpty"`
	Port                int      `env:"PORT,notEmpty"`
	Host                string   `env:"HOST"`
	ProjectID           string   `env:"PROJECT_ID,notEmpty"`
	TopicID             string   `env:"TOPIC_ID,notEmpty"`
}
