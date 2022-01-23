package main

import (
	"github.com/caarlos0/env/v6"
	"github.com/producer_rmq_fsevent/event"
	"github.com/producer_rmq_fsevent/freeswitch"
	"github.com/producer_rmq_fsevent/rabbitmq"
	"log"
)

type config struct {
	Uri    string `env:"RABBITMQ_URI" envDefault:"amqp://admin:admin@127.0.0.1:5672/"`
	Queue  string `env:"RABBITMQ_QUEUE" envDefault:"call_queue_stats"`
	FsAddr string `env:"FS_ADDR" envDefault:"127.0.0.1:8021"`
	FsPass string `env:"FS_PASS" envDefault:"ClueCon"`
}

var cfg config

const (
	subEventList = "events text CHANNEL_PARK " +
		"CHANNEL_ANSWER " +
		"RECORD_START RECORD_STOP CHANNEL_ORIGINATE " +
		"CHANNEL_PROGRESS_MEDIA " +
		"CHANNEL_HANGUP_COMPLETE DTMF SESSION_HEARTBEAT " +
		"CUSTOM conference::maintenance"
)

func init() {

	if err := env.Parse(&cfg); err != nil {
		log.Printf("%+v\n", err)
	}

}

func main() {

	var rabbitMqHandle rabbitmq.RmqInterface
	var err error

	if rabbitMqHandle, err = rabbitmq.NewRmqAdapter(cfg.Uri); err != nil {
		log.Fatalf("%s", err)
	}

	evHandler := event.NewEventHandler(rabbitMqHandle, cfg.Queue)

	if fsHandler, err := freeswitch.NewFreeSwitchHandler(cfg.FsAddr, cfg.FsPass); err != nil {
		log.Fatalf("%s", err)
	} else if err = fsHandler.SubscribeFreeSwitchEvent(subEventList, evHandler.ProcessFsEvent); err != nil {
		log.Fatalf("%s", err)
	}

	select {}
}
