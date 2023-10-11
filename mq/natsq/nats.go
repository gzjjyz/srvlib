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

	js nats.JetStreamContext
)

var (
	ErrTopicAlreadyRegistered = errors.New("topic already registered")
	ErrTopicNotRegistered     = errors.New("topic not registered")
	ErrAlreadyInited          = errors.New("natsq already inited")
)

func Init(url string) error {
	if conn != nil || js != nil {
		return ErrAlreadyInited
	}

	if url != "" {
		Url = url
	}

	var err error
	rawNatsConn, err := nats.Connect(Url)
	if err != nil {
		return err
	}

	conn = &NatsConn{
		Conn:          rawNatsConn,
		subscriptions: sync.Map{},
	}

	js, err = conn.Conn.JetStream()
	if err != nil {
		return err
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
func (c *NatsConn) RegSubForFurtherMannipulate(topic string, subObj *nats.Subscription) error {
	_, ok := c.subscriptions.Load(topic)

	if !ok {
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
