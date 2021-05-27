package socketio_emitter

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/vmihailenco/msgpack"
	"gopkg.in/redis.v5"
)

const EVENT = 2
const uid = "emitter"

type Emitter struct {
	Redis   *redis.Client
	Prefix  string
	Nsp     string
	Channel string
	rooms   []string
	flags   map[string]interface{}
}

type NewEmitterOpts struct {
	// redis host
	Host string
	// redis port
	Port int
	// redis addr
	Addr string
	// redis connection protocol
	Protocol string
	// prefix of pub/sub key
	Key string
	// namespace
	Nsp string
}

// Emitter constructor
func NewEmitter(opts *NewEmitterOpts) (*Emitter, error) {
	var addr string
	if opts.Addr != "" {
		addr = opts.Addr
		if !strings.HasPrefix(addr, "redis://") {
			addr = "redis://" + addr
		}
	} else if opts.Host != "" && opts.Port > 0 {
		addr = fmt.Sprintf("redis://%s:%d", opts.Host, opts.Port)
	} else {
		addr = "redis://localhost:6379"
	}

	config, err := redis.ParseURL(addr)
	if err != nil {
		fmt.Println("an error occurred when parsing redis dsn:", err)
		return nil, err
	}
	client := redis.NewClient(config)

	prefix := "socket.io"
	if opts.Key != "" {
		prefix = opts.Key
	}

	emitter := newEmitter(client, prefix, "/")
	return emitter, emitter.Redis.Ping().Err()
}

func newEmitter(redis *redis.Client, prefix string, nsp string) *Emitter {
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

	opts := map[string]interface{}{
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
	_ = emitter.Redis.Publish(channel, buf.String())

	emitter.rooms = []string{}
	emitter.flags = make(map[string]interface{})
	return emitter, err
}
