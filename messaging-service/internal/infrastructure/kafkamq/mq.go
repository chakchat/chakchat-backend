package kafkamq

import (
	"context"

	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type KafkaMQPublisher struct {
	writer *kafka.Writer
	tracer trace.Tracer
}

func NewKafkaMQPublisher(writer *kafka.Writer) *KafkaMQPublisher {
	return &KafkaMQPublisher{
		writer: writer,
		tracer: otel.GetTracerProvider().Tracer("kafka-mq-publisher"),
	}
}

func (k *KafkaMQPublisher) Publish(ctx context.Context, raw []byte) error {
	var span trace.Span
	ctx, span = k.tracer.Start(ctx, "write")
	defer span.End()

	msg := kafka.Message{
		Value: raw,
	}

	if err := k.writer.WriteMessages(ctx, msg); err != nil {
		span.RecordError(err)
		return err
	}

	return nil
}
