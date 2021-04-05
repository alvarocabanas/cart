package async

import (
	"context"
	"time"

	"go.opencensus.io/trace"

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
	_, span := trace.StartSpan(ctx, "handle_add_item_event")
	defer span.End()

	m.metricsRecorder.Record(time.Now().Unix(), metrics.AddItemEventHandled)
	return nil
}
