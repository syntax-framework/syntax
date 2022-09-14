package syntax

import (
	"fmt"
	"github.com/syntax-framework/chain"
	"github.com/syntax-framework/chain/middlewares/session"
	"github.com/syntax-framework/shtml/cmn"
	"github.com/syntax-framework/shtml/sht"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// syntaxValidBase64Regex "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_"
var syntaxValidBase64Regex = regexp.MustCompile(`[a-zA-Z0-9\-_]+`)

func isBase64(name string) bool {
	return syntaxValidBase64Regex.MatchString(name)
}

var errorInvalidChannelName = cmn.Err(
	"pubsub.channel.name",
	"Not a valid Channel name.", "Name: %s",
)

var errorInvalidTopicName = cmn.Err(
	"pubsub.topic.name",
	"Not a valid Topic name.", "Name: %s",
)

var errorTopicOnJoinExists = cmn.Err(
	"pubsub.topic.join.exists",
	"There is already an OnJoin callback with the same name.", "Channel: %s", "Topic: %s",
)

// Todos os Channels dentro de um Topic
// @TODO: Mensagens persistentes

// SSEEvent holds all of the event source fields
type SSEEvent struct {
	timestamp time.Time
	ID        []byte
	Data      []byte
	Event     []byte
	Retry     []byte
	Comment   []byte
}

type SSESubscription struct {
	URL         *url.URL
	LastEventID int
	quit        chan *SSESubscription
	removed     chan struct{}
	event       chan *SSEEvent // Send message to client
}

func (s *SSESubscription) close() {
	s.quit <- s
	if s.removed != nil {
		<-s.removed
	}
}

// Socket represents a user's connection to a specific Channel
type Socket struct {
	//Params  Params
	Channel *Channel
	request *http.Request
}

// Channel handle events from clients. Channels are the highest level abstraction for real-time communication components
// in Syntax.
//
// Channels provide a means for bidirectional communication from clients that integrate with the Syntax PubSub layer
// for soft-realtime functionality.
type Channel struct {
	name   string
	onJoin map[string]ChannelOnJoinFunc
}

// ChannelOnJoinFunc  Clients must join a channel to send and receive PubSub events on that channel.
//
// # Authorization
//
// Your channels must register a `OnJoin()` callback that authorizes the socket	for the given topic.
// For example, you could check if the user is allowed to	join that particular room.
//
// To authorize a socket in [Channel.OnJoin()], return `nil`.
//
//	To refuse authorization in [Channel.OnJoin()], return `error`.
type ChannelOnJoinFunc func(topic string, params map[string]interface{}, socket *Socket) error

type ChannelOnMessageFunc func(topic string, params map[string]interface{}, socket *Socket) error

// OnJoin invocado quando o client procura conectar-se a este tópico neste channel
func (c *Channel) OnJoin(topic string, callback ChannelOnJoinFunc) error {
	topic = strings.TrimPrefix(strings.TrimSpace(topic), ":")
	if topic != "*" && (topic == "" || !isBase64(topic)) {
		return errorInvalidTopicName(topic)
	}

	if _, exist := c.onJoin[topic]; exist {
		return errorTopicOnJoinExists(c.name, topic)
	}

	c.onJoin[topic] = callback

	return nil
}

func (c *Channel) On(event string, callback ChannelOnMessageFunc) {

}

// Channel registra um novo channel para comunicação em tempo real
func (s *Syntax) Channel(name string) error {
	name = strings.TrimSpace(name)
	if name == "" || !isBase64(name) {
		return errorInvalidChannelName(name)
	}

	channel := &Channel{
		name: name,
	}

	println(channel)

	return nil
}

// Publish publica em um canal
func (s *Syntax) Publish(topic string, content interface{}) {

}

// Subscribe subscreve em um canal
func (s *Syntax) Subscribe(topic string, cb func(content interface{})) {

}

// initLiveServer iniciliza a conexão viva com esse servidor. Usado para push de eventos e escuta de SSE
func (s *Syntax) initLiveServer() error {

	endpoint := strings.TrimSpace(s.Config.LiveEndpoint)
	if endpoint == "" {
		endpoint = "/live"
	}
	if !strings.HasPrefix(endpoint, "/") {
		endpoint = "/" + endpoint
	}
	endpoint = strings.TrimSuffix(endpoint, "/")

	// add stx.js asset, required on all pages
	asset, err := s.Template.(*sht.TemplateSystem).RegisterAssetJsFilepath("/assets/js/stx.js")
	if err != nil {
		return err
	}
	s.Bundler.AddRequiredAsset(asset)

	asset.Priority = 100
	asset.Attributes = map[string]string{
		"data-endpoint": endpoint,
	}

	// <script src="./../assets/js/stx.js" priority="100"></script>

	// Client conecta-se ao SSE e servidor salva client em uma list, quando cliente desconectar, remove-o
	// Client recebe mensagens via SSE
	// Client informa tópicos que deseja ouvir
	// Client submete PUSH vis POST
	// Uma live controller é um tópico dinamico, existente somente neste server
	// A Live Controller ouve eventos de uma conexão em um tópico

	// lista de usuarios
	// connect e disconnect
	// close conn.channel
	// conn.client = nil

	s.Use(endpoint, &session.Manager{
		Config: session.Config{
			Key:  "_ls",
			Path: endpoint,
			//Domain:     "",
			//Expires:    time.Time{},
			//RawExpires: "",
			//MaxAge:     0,
			//Secure:     false,
			//HttpOnly:   false,
			//SameSite:   0,
			//Raw:        "",
			//Unparsed:   nil,
		},
		Store: &session.Cookie{
			CryptoOptions: session.CryptoOptions{
				SecretKeyBase:  "REMOVER DEPOIS COLOCADO APENAS PARA TESTE",
				EncryptionSalt: "",
				SigningSalt:    "REMOVER DEPOIS COLOCADO APENAS PARA TESTE",
				//Iterations:     0,
				//Length:         0,
				//Digest:         "",
			},
			//Serializer:      nil,
			//Log:             "",
			//RotatingOptions: nil,
		},
	})

	s.POST(endpoint, func(ctx *chain.Context) {
		// @TODO: Parse user command
		filepath := ctx.GetParam("filepath")

		sess, _ := session.Fetch(ctx)
		println(sess.Get("user"))

		//decoder := json.NewDecoder(req.Body)
		//var t test_struct
		//err := decoder.Decode(&t)
		//if err != nil {
		//  panic(err)
		//}
		//log.Println(t.Test)

		println(filepath)
	})

	var messageChan chan string

	go func() {
		for {
			//b := []byte(time.Now().Format(time.RFC3339))
			//if err := nc.Notify(b); err != nil {
			//	log.Fatal(err)
			//}

			if messageChan != nil {
				log.Printf("print message to client")
				// send the message through the available channel
				messageChan <- sht.HashMD5(time.Now().String())
			}

			time.Sleep(1 * time.Second)
		}
	}()

	s.GET(endpoint, func(ctx *chain.Context) {
		w := ctx.Writer.(*chain.ResponseWriterSpy)
		r := ctx.Request

		flusher, ok := w.ResponseWriter.(http.Flusher)
		if !ok {
			http.Error(w, "Connection does not support streaming", http.StatusBadRequest)
			return
		}

		sess, _ := session.Fetch(ctx)
		sess.Put("user", "test")

		lastEventId := 0
		if id := r.Header.Get("Last-Event-ID"); id != "" {
			var err error
			if lastEventId, err = strconv.Atoi(id); err != nil {
				http.Error(w, "Last-Event-ID must be a number!", http.StatusBadRequest)
				return
			}
		}

		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		//w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		//client := &liveReloadClient{
		//	addr:   r.RemoteAddr,
		//	events: make(chan *liveReloadEvent, 10),
		//}
		//go updateClient(client)

		sub := &SSESubscription{
			LastEventID: lastEventId,
			URL:         r.URL,
			//quit:        str.deregister,
			//Connection:  make(chan *Event, 64),
		}

		go func() {
			// Received Browser Disconnection
			<-r.Context().Done()
			sub.close()
		}()

		w.WriteHeader(http.StatusOK)
		flusher.Flush()

		messageChan = make(chan string)
		defer func() {
			close(messageChan)
			messageChan = nil
		}()

		// test
		//timeout := time.After(5 * time.Second)

		// trap the request under loop forever
		for {
			select {
			case message := <-messageChan:
				fmt.Fprintf(w, "data: \"%s\"\n\n", message)
				flusher.Flush()
			case <-r.Context().Done():
				return
			}
		}

		//select {
		//case ev := <-client.events:
		//	var buf bytes.Buffer
		//	enc := json.NewEncoder(&buf)
		//	enc.Encode(ev)
		//	fmt.Fprintf(w, "event: xpto\ndata: %v\n\n", buf.String())
		//case <-timeout:
		//	fmt.Fprintf(w, ": nothing to sent\n\n")
		//}
		//flusher.Flush()

		//ctx := r.Context()
		//
		//ch := make(chan struct{})
		//
		//go func() {
		//  time.Sleep(5 * time.Second)
		//  fmt.Fprintln(w, "Hello World!")
		//  ch <- struct{}{}
		//}()
		//
		//select {
		//case <-ch:
		//case <-ctx.Done():
		//  http.Error(w, ctx.Err().Error(), http.StatusPartialContent)
		//  someCleanUP()
		//}
	})
	return nil
}

func updateClient(client *liveReloadClient) {
	for {
		client.events <- &liveReloadEvent{
			EventType: uint(rand.Uint32()),
		}
	}
}
