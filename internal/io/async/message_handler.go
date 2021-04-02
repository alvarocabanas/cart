package async

import (
	"context"

	"github.com/alvarocabanas/cart/internal/metrics"
)

type MessageHandler struct {
	metricsRecorder metrics.Recorder
}

func NewMessageHandler(metricsRecorder metrics.Recorder) MessageHandler {
	return MessageHandler{
		metricsRecorder: metricsRecorder,
	}
}

func (m MessageHandler) Handle(ctx context.Context, message []byte) error {
	return nil
}
