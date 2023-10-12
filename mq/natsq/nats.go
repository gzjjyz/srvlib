package natsq

import (
	"errors"
	"sync"

	"github.com/nats-io/nats.go"
)

// NatsConn: manage nats connection and subscriptions.
type NatsConn struct {
	*nats.Conn
	subscriptions sync.Map
}

var (
	Url  = nats.DefaultURL
	conn *NatsConn

	js *JetStreamContext
)

type DrainAllSubscriptionsErrorCb func(topic string, err error)

var (
	ErrTopicAlreadyRegistered = errors.New("topic already registered")
	ErrTopicNotRegistered     = errors.New("topic not registered")
	ErrAlreadyInited          = errors.New("natsq already inited")
)

func Init(url string, opts ...nats.Option) error {
	if conn != nil || js != nil {
		return ErrAlreadyInited
	}

	if url != "" {
		Url = url
	}

	var err error
	rawNatsConn, err := nats.Connect(Url, opts...)
	if err != nil {
		return err
	}

	conn = &NatsConn{
		Conn:          rawNatsConn,
		subscriptions: sync.Map{},
	}

	rawJs, err := conn.Conn.JetStream()
	if err != nil {
		return err
	}

	js = &JetStreamContext{
		JetStreamContext: rawJs,
	}

	return nil
}

func GetConn() *NatsConn {
	return conn
}

func (c *NatsConn) JetStream() nats.JetStreamContext {
	return js
}

// generally you don`t need to use this method, the "raw nats conn" will manage all subscriptions.
// only when you have further manipulation on subscription you should registe it.
func (c *NatsConn) RegisterSubscribtion(topic string, subObj *nats.Subscription) error {
	_, ok := c.subscriptions.Load(topic)

	if ok {
		return ErrTopicAlreadyRegistered
	}

	c.subscriptions.Store(topic, subObj)
	return nil
}

func (c *NatsConn) GetRegistedSub(topic string) *nats.Subscription {
	foo, ok := c.subscriptions.Load(topic)
	if !ok {
		return nil
	}

	sub, ok := foo.(*nats.Subscription)
	if !ok {
		return nil
	}

	return sub
}

func (c *NatsConn) Subscribe(subj string, cb nats.MsgHandler) (sub *nats.Subscription, err error) {
	_, ok := c.subscriptions.Load(subj)
	if ok {
		return nil, ErrTopicAlreadyRegistered
	}

	sub, err = c.Conn.Subscribe(subj, cb)
	if err != nil {
		return nil, err
	}

	err = c.RegisterSubscribtion(subj, sub)
	return sub, err
}

// NOTICE: Gnerally conn will manage all subscriptions, so you don`t need to Unsubscribe manually.
// in special case, you can use this method to Unsubscribe as you whish.
func (c *NatsConn) UnSubscribeRegsitedSub(topic string) error {
	subObj := c.GetRegistedSub(topic)
	if subObj == nil {
		return ErrTopicNotRegistered
	}

	err := subObj.Drain()
	if err != nil {
		return err
	}

	c.subscriptions.Delete(topic)
	return nil
}

func (c *NatsConn) DrainAllSubscription(errCb DrainAllSubscriptionsErrorCb) {
	c.subscriptions.Range(func(key, value any) bool {
		sub, ok := value.(*nats.Subscription)
		if !ok {
			errCb(key.(string), ErrTopicNotRegistered)
		}

		err := sub.Drain()
		if err != nil {
			errCb(key.(string), err)
		}
		return true
	})
}

// Close will wait for all msg processed and Unsubscribe all  conn/JetStream`s subscriptions.
func (c *NatsConn) Close() error {
	if conn == nil {
		return nil
	}

	return c.Drain()
}

// ForceClose is unsafe close, it will close connection immediately without process left msg
// and  unscript conn/JetStream`s subscriptions.
func (c *NatsConn) ForceClose() {
	c.Close()
}

type JetStreamContext struct {
	nats.JetStreamContext
}

func (js *JetStreamContext) Subscribe(topic string, cb nats.MsgHandler, opts ...nats.SubOpt) (sub *nats.Subscription, err error) {
	_, ok := conn.subscriptions.Load(topic)
	if ok {
		return nil, ErrTopicAlreadyRegistered
	}

	sub, err = js.JetStreamContext.Subscribe(topic, cb)
	if err != nil {
		return nil, err
	}

	err = conn.RegisterSubscribtion(topic, sub)
	return sub, err
}
