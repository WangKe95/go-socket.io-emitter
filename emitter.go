package socketio_emitter

import (
	"bytes"
	"fmt"

	"github.com/garyburd/redigo/redis"
	"github.com/vmihailenco/msgpack"
)

const EVENT = 2
const uid = "emitter"

type Emitter struct {
	Redis   redis.Conn
	Prefix  string
	Nsp     string
	Channel string
	rooms   []string
	flags   map[string]interface{}
}

type NewEmitterOpts struct {
	// redis host
	Host     string
	// redis port
	Port     int
	// redis addr
	Addr     string
	// redis connection protocol
	Protocol string
	// prefix of pub/sub key
	Key      string
	// namespace
	Nsp      string
}

// Emitter constructor
func NewEmitter(opts *NewEmitterOpts) (*Emitter, error) {
	var addr string
	if opts.Addr != "" {
		addr = opts.Addr
	} else if opts.Host != "" && opts.Port > 0 {
		addr = fmt.Sprintf("%s:%d", opts.Host, opts.Port)
	} else {
		addr = "localhost:6379"
	}

	protocol := "tcp"
	if opts.Protocol != "" {
		protocol = opts.Protocol
	}

	conn, err := redis.Dial(protocol, addr)
	if err != nil {
		return nil, err
	}

	prefix := "socket.io"
	if opts.Key != "" {
		prefix = opts.Key
	}

	return newEmitter(conn, prefix, "/"), nil
}

func newEmitter(redis redis.Conn, prefix string, nsp string) *Emitter {
	return &Emitter{
		Redis:   redis,
		Prefix:  prefix,
		Nsp:     nsp,
		Channel: fmt.Sprintf("%s#%s#", prefix, nsp),
		rooms:   []string{},
		flags:   make(map[string]interface{}),
	}
}

// Limit emission to a certain `room`.
func (emitter *Emitter) In(room string) *Emitter {
	return emitter.To(room)
}

// Limit emission to a certain `room`.
func (emitter *Emitter) To(room string) *Emitter {
	for _, r := range emitter.rooms {
		if r == room {
			return emitter
		}
	}
	emitter.rooms = append(emitter.rooms, room)
	return emitter
}

// Return a new emitter for the given namespace.
func (emitter *Emitter) Of(nsp string) *Emitter {
	return newEmitter(emitter.Redis, emitter.Prefix, nsp)
}

// Send the packet.
func (emitter *Emitter) Emit(event string, data ...interface{}) (*Emitter, error) {
	pack := make([]interface{}, 0)
	pack = append(pack, uid)

	args := []interface{}{event}
	args = append(args, data...)
	packet := make(map[string]interface{})
	packet["type"] = EVENT
	packet["data"] = args
	packet["nsp"] = emitter.Nsp
	pack = append(pack, packet)

	opts := map[string]interface{} {
		"rooms": emitter.rooms,
		"flags": emitter.flags,
	}
	pack = append(pack, opts)

	channel := emitter.Channel
	if emitter.rooms != nil && len(emitter.rooms) == 1 {
		channel = fmt.Sprintf("%s%s#", channel, emitter.rooms[0])
	}
	buf := &bytes.Buffer{}
	enc := msgpack.NewEncoder(buf)
	err := enc.Encode(pack)
	if err != nil {
		return nil, err
	}
	_, err = emitter.Redis.Do("publish", channel, buf)

	emitter.rooms = []string{}
	emitter.flags = make(map[string]interface{})
	return emitter, err
}
