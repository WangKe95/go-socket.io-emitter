package socketio_emitter

import (
	"fmt"
	"strings"
	"testing"

	"gopkg.in/redis.v5"
)

func TestEmitter_Emit(t *testing.T) {
	emitter, _ := NewEmitter(&NewEmitterOpts{
		Host: "localhost",
		Port: 6379,
	})
	if emitter == nil {
		t.Error("emitter is nil")
	}

	conn := redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    "localhost:6379",
	})
	defer conn.Close()

	psc, _ := conn.Subscribe("socket.io#/#")
	emitter.Emit("test", "test emitter")
	for {
		v, _ := psc.Receive()
		switch msg := v.(type) {
		case *redis.Message:
			isContain := strings.Contains(msg.Payload, "test emitter")
			if !isContain {
				t.Error("received msg not right")
				return
			} else {
				return
			}
		}
	}
}

func TestEmitter_In(t *testing.T) {
	emitter, _ := NewEmitter(&NewEmitterOpts{
		Host: "localhost",
		Port: 6379,
	})
	if emitter == nil {
		t.Error("emitter is nil")
	}

	conn := redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    "localhost:6379",
	})
	defer conn.Close()

	room := "test_room"
	psc, _ := conn.Subscribe(fmt.Sprintf("%s%s#", emitter.Channel, room))
	emitter.In(room).Emit("test", "test room")
	for {
		v, _ := psc.Receive()
		switch msg := v.(type) {
		case *redis.Message:
			isContain := strings.Contains(msg.Payload, "test room")
			if !isContain {
				t.Error("received msg not right")
				return
			} else {
				return
			}
		}
	}
}

func TestEmitter_Of(t *testing.T) {
	emitter, _ := NewEmitter(&NewEmitterOpts{
		Host: "localhost",
		Port: 6379,
	})
	if emitter == nil {
		t.Error("emitter is nil")
	}

	conn := redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    "localhost:6379",
	})
	defer conn.Close()

	nsp := "test_nsp"
	psc, _ := conn.Subscribe(fmt.Sprintf("%s#%s#", emitter.Prefix, nsp))
	emitter.Of(nsp).Emit("test", "test nsp")
	for {
		v, _ := psc.Receive()
		switch msg := v.(type) {
		case *redis.Message:
			isContain := strings.Contains(msg.Payload, "test nsp")
			if !isContain {
				t.Error("received msg not right")
				return
			} else {
				return
			}
		}
	}
}
