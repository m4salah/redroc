package pubsub

import "context"

type UploadMessage struct {
	Message struct {
		Data       []byte `json:"data,omitempty"`
		ID         string `json:"messageId"`
		Attributes struct {
			ObjName  string `json:"objName"`
			Hashtags string `json:"hashtags"`
			User     string `json:"user"`
		} `json:"attributes"`
	} `json:"message"`
	Subscription string `json:"subscription"`
}

type PubSub interface {
	Pub(ctx context.Context, message any)
	Sub(ctx context.Context, message any)
}
