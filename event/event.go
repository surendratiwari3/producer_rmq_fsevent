package event

import "github.com/producer_rmq_fsevent/rabbitmq"

type EventInterface interface {
	ProcessFsEvent(eventData []byte) error
}

type EventHandler struct {
	rmqHandler rabbitmq.RmqInterface
	rmqQueueName   string
}

func NewEventHandler(rmqInterface rabbitmq.RmqInterface, queueName string) EventInterface {
	return &EventHandler{
		rmqHandler:rmqInterface,
		rmqQueueName:queueName,
	}
}

func (eh *EventHandler) ProcessFsEvent(eventData []byte) error {
	return eh.rmqHandler.PublishCallStats(eh.rmqQueueName, eventData)
}