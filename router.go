package salt

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/nats-io/nats.go"
)

// Router struct
type Router struct {
	routes      map[string]HandlerFunc
	eventRoutes map[string]HandlerFunc
	prefix      string
	rbac        RBACContext
	resource    string
}

// NewRouter creates a new router
func NewRouter() *Router {
	return &Router{
		routes: make(map[string]HandlerFunc),
	}
}

func (r *Router) WithPermission(rbac RBACContext) *Router {
	r.rbac = rbac
	return r
}

// Handle registers a new route with a handler
func (r *Router) Handle(topic string, handler HandlerFunc) {
	fullTopic := strings.TrimPrefix(r.prefix+"."+strings.TrimPrefix(topic, "."), ".")
	if r.resource != "" {
		fullTopic = r.resource + "." + fullTopic
	}
	rbac := RBAC(topic)
	r.routes[fullTopic] = rbac(handler)
}

func (r *Router) HandleEvent(topic string, handler HandlerFunc) {
	fullTopic := strings.TrimPrefix(r.prefix+"."+strings.TrimPrefix(topic, "."), ".")
	r.eventRoutes[fullTopic] = handler
}

// Dispatch dispatches the message to the right handler
func (r *Router) Dispatch(topic string, ctx *Context) {
	if handler, exists := r.routes[topic]; exists {
		handler(ctx)
	}
}

func (r *Router) WithResource(resource string) *Router {
	r.resource = resource
	return r
}

// Run runs the router, connects to nats and subscribes to all routes
func (r *Router) Run(natsURL string) error {
	nc, err := nats.Connect(natsURL)
	if err != nil {
		return err
	}

	for topic, _ := range r.routes {
		nc.QueueSubscribe(topic, topic, func(msg *nats.Msg) {
			var systemMessage SystemMessage
			err := json.Unmarshal(msg.Data, &systemMessage)
			if err != nil {
				log.Printf("Error parsing system message: %v", err)
				nc.Publish(msg.Reply, []byte(NewSystemResponse(http.StatusInternalServerError, "Invalid json").ToJSON()))
				return
			}
			ctx := &Context{nc: nc, ncMsg: msg, req: systemMessage, rbac: r.rbac, router: r}
			r.Dispatch(msg.Subject, ctx)
		})
	}

	for topic, _ := range r.eventRoutes {
		nc.Subscribe(topic, func(msg *nats.Msg) {
			var systemMessage SystemMessage
			err := json.Unmarshal(msg.Data, &systemMessage)
			if err != nil {
				log.Printf("Error parsing system message: %v", err)
				nc.Publish(msg.Reply, []byte(NewSystemResponse(http.StatusInternalServerError, "Invalid json").ToJSON()))
				return
			}
			ctx := &Context{nc: nc, ncMsg: msg, req: systemMessage, rbac: r.rbac, router: r}
			r.Dispatch(msg.Subject, ctx)
		})
	}
	return nil
}
