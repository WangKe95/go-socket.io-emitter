package socketio_emitter

import (
	"fmt"
	"strings"
	"testing"

	"github.com/garyburd/redigo/redis"
)

func TestEmitter_Emit(t *testing.T) {
	emitter, _ := NewEmitter(&NewEmitterOpts{
		Host: "localhost",
		Port: 6379,
	})
	if emitter == nil {
		t.Error("emitter is nil")
	}

	conn, _ := redis.Dial("tcp", "localhost:6379")
	defer conn.Close()
	psc := redis.PubSubConn{Conn: conn}

	psc.Subscribe("socket.io#/#")
	emitter.Emit("test", "test emitter")
	for {
		switch v := psc.Receive().(type) {
		case redis.Message:
			isContain := strings.Contains(string(v.Data), "test emitter")
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

	conn, _ := redis.Dial("tcp", "localhost:6379")
	defer conn.Close()
	psc := redis.PubSubConn{Conn: conn}

	room := "test_room"
	psc.Subscribe(fmt.Sprintf("%s%s#", emitter.Channel, room))
	emitter.In(room).Emit("test", "test room")
	for {
		switch v := psc.Receive().(type) {
		case redis.Message:
			isContain := strings.Contains(string(v.Data), "test room")
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

	conn, _ := redis.Dial("tcp", "localhost:6379")
	defer conn.Close()
	psc := redis.PubSubConn{Conn: conn}

	nsp := "test_nsp"
	psc.Subscribe(fmt.Sprintf("%s#%s#", emitter.Prefix, nsp))
	emitter.Of(nsp).Emit("test", "test nsp")
	for {
		switch v := psc.Receive().(type) {
		case redis.Message:
			isContain := strings.Contains(string(v.Data), "test nsp")
			if !isContain {
				t.Error("received msg not right")
				return
			} else {
				return
			}
		}
	}
}
