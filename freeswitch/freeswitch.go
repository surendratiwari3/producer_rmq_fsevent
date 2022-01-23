package freeswitch

import (
	"encoding/json"
	"github.com/fiorix/go-eventsocket/eventsocket"
	"time"
)

type eventCallbackFunc func([]byte) error

type FreeSwitchInterface interface {
	SubscribeFreeSwitchEvent(evList string, ev eventCallbackFunc) error
}

type freeSwitchHandler struct {
	fsConn              *eventsocket.Connection
	fsAddr              string
	fsPass              string
	fsReconnectInterval int
}

const (
	freeSwitchReconnectTime = 5
)

func NewFreeSwitchHandler(fsAddr string, fsPass string) (FreeSwitchInterface, error) {

	fsConn, err := eventsocket.Dial(fsAddr, fsPass)
	if err != nil {
		return nil, err
	}

	fsHandler := &freeSwitchHandler{
		fsConn: fsConn,
	}

	ticker := time.NewTicker(freeSwitchReconnectTime * time.Second)

	quit := make(chan struct{})

	go func() {
		for {
			select {
			case <-ticker.C:
				fsHandler.reconnect()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	return fsHandler, nil
}

func (fh *freeSwitchHandler) SubscribeFreeSwitchEvent(evList string, ev eventCallbackFunc) error {
	_, err := fh.fsConn.Send(evList)
	go fh.eventProcessor(ev)
	return err
}

func (fh *freeSwitchHandler) reconnect() {
	var err error

	if fh.fsConn != nil {
		return
	}

	fh.fsConn, err = eventsocket.Dial(fh.fsAddr, fh.fsPass)
	if err == nil {
		return
	}
}

func (fh *freeSwitchHandler) eventProcessor(eventCallback eventCallbackFunc) {
	for {
		ev, err := fh.fsConn.ReadEvent()
		if err != nil {
			return
		}

		evHeaderByte, err := json.Marshal(ev.Header)

		if err != nil {
			continue
		}

		go func() {
			_ = eventCallback(evHeaderByte)
		}()
	}
}
