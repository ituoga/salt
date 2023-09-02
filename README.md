# Salt - tiny wrapper around NATS

still in early alfa so API will change!

todo:
- pub/sub
- todo

`go get -u -v github.com/ituoga/salt`

```golang
package main

import (
	"log"
	"time"

	"github.com/ituoga/salt"
)

type permissionsChecker struct {
	token    string
	resource string
	db       *struct{}
}

func (p *permissionsChecker) WithToken(token string) salt.RBACContext {
	p.token = token
	return p
}

func (p permissionsChecker) Can(action string) bool {
	// TODO: get User from DB
	// TODO: get user permissions from DB
	// TODO: check if user has action and/or for resource
	return true
}

func (p *permissionsChecker) WithResource(resource string) salt.RBACContext {
	p.resource = resource
	return p
}

var _ salt.RBACContext = (*permissionsChecker)(nil)

func main() {
	router := salt.NewRouter()
	router.WithPermission(&permissionsChecker{
		// TODO: implement this
		// TODO: db connection and other
	})

	router.Handle("topic.example-1", func(ctx *salt.Context) {
		ctx.Response().Reply("topic")
	})

	router.Handle("topic.example-2", func(ctx *salt.Context) {
		if ctx.Can("view-own") {
			ctx.Response().Reply("owns")
			return
		}
		ctx.Response().Reply("default resonse")
	})

	go func() {
		c, _ := salt.NewClient("nats://localhost:4222")
		for {
			time.Sleep(1 * time.Second)
			response, err := c.Request("topic.example-1",
				salt.WithPayloadClient("hello world"),
				salt.WithTokenClient("token"),
			)
			if err != nil {
				log.Printf("%v", err)
				continue
			}
			log.Printf("%v", response)
		}
	}()

	router.Run("nats://localhost:4222")
	select {}
}

```