package natsq

import (
	"errors"
	"sync"

	"github.com/nats-io/nats.go"
)

// NatsConn: manage nats connection and Subscriptions.
// Don`t use Subscription.UnSubScribe function instead use NatsConn.UnSubScribe except you known what are you doing, for the reason it would confuse NatsConn
// Don`t use Subscription.Drain function instead use NatsConn.DrainSubscription except you known what are you doing, for the reason it would confuse NatsConn
type NatsConn struct {
	rawConn       *nats.Conn
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
		rawConn:       rawNatsConn,
		subscriptions: sync.Map{},
	}

	rawJs, err := conn.rawConn.JetStream()
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

// Notice!!!!: don`t use rawConn`s Subscribe group functions
// for the reason if would confuse NatsConn
func (c *NatsConn) RawNatsConn() *nats.Conn {
	return c.rawConn
}

// Notice!!!: don`t use rawConn`s JetStream() function Use NatsConn`s JetStream() function instead
// for the reason it would confuse NatsConn
func (c *NatsConn) JetStream() nats.JetStreamContext {
	return js
}

func (c *NatsConn) UnSubscribe(topic string) error {
	foo, ok := c.subscriptions.Load(topic)
	if !ok {
		return ErrTopicNotRegistered
	}

	sub, ok := foo.(*nats.Subscription)
	if !ok {
		return ErrTopicNotRegistered
	}

	err := sub.Drain()
	if err != nil {
		return err
	}

	c.subscriptions.Delete(topic)
	return nil
}

func (c *NatsConn) DrainSubscription(topic string) error {
	return c.UnSubscribe(topic)
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

	sub, err = c.rawConn.Subscribe(subj, cb)
	if err != nil {
		return nil, err
	}

	c.subscriptions.Store(subj, sub)
	return
}

func (c *NatsConn) SubscribeSync(subj string) (sub *nats.Subscription, err error) {
	_, ok := c.subscriptions.Load(subj)
	if ok {
		return nil, ErrTopicAlreadyRegistered
	}

	sub, err = c.rawConn.SubscribeSync(subj)
	if err != nil {
		return nil, err
	}

	c.subscriptions.Store(subj, sub)
	return
}

func (c *NatsConn) QueueSubscribe(subj, queue string, cb nats.MsgHandler) (sub *nats.Subscription, err error) {
	_, ok := c.subscriptions.Load(subj)
	if ok {
		return nil, ErrTopicAlreadyRegistered
	}

	sub, err = c.rawConn.QueueSubscribe(subj, queue, cb)
	if err != nil {
		return nil, err
	}

	c.subscriptions.Store(subj, sub)
	return sub, err
}

func (c *NatsConn) ChanSubscribe(subj string, ch chan *nats.Msg) (sub *nats.Subscription, err error) {
	_, ok := c.subscriptions.Load(subj)
	if ok {
		return nil, ErrTopicAlreadyRegistered
	}

	sub, err = c.rawConn.ChanSubscribe(subj, ch)
	if err != nil {
		return nil, err
	}

	c.subscriptions.Store(subj, sub)
	return
}

func (c *NatsConn) QueueSubscribeSync(subj, queue string) (sub *nats.Subscription, err error) {
	_, ok := c.subscriptions.Load(subj)
	if ok {
		return nil, ErrTopicAlreadyRegistered
	}

	sub, err = c.rawConn.QueueSubscribeSync(subj, queue)
	if err != nil {
		return nil, err
	}

	c.subscriptions.Store(subj, sub)
	return
}

func (c *NatsConn) QueueSubscribeSyncWithChan(subj, queue string, ch chan *nats.Msg) (sub *nats.Subscription, err error) {
	_, ok := c.subscriptions.Load(subj)
	if ok {
		return nil, ErrTopicAlreadyRegistered
	}

	sub, err = c.rawConn.QueueSubscribeSyncWithChan(subj, queue, ch)
	if err != nil {
		return nil, err
	}

	c.subscriptions.Store(subj, sub)
	return
}

func (c *NatsConn) ChanQueueSubscribe(subj, queue string, ch chan *nats.Msg) (sub *nats.Subscription, err error) {
	_, ok := c.subscriptions.Load(subj)
	if ok {
		return nil, ErrTopicAlreadyRegistered
	}

	sub, err = c.rawConn.ChanQueueSubscribe(subj, queue, ch)
	if err != nil {
		return nil, err
	}

	c.subscriptions.Store(subj, sub)
	return
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

func (c *NatsConn) Close() error {
	if conn == nil {
		return nil
	}

	return c.rawConn.Drain()
}

// ForceClose is unsafe close, it will close connection immediately without process left msg
// and  unscript conn/JetStream`s subscriptions.
func (c *NatsConn) ForceClose() {
	c.rawConn.Close()
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

	conn.subscriptions.Store(topic, sub)
	return
}

func (js *JetStreamContext) SubscribeSync(subj string, opts ...nats.SubOpt) (sub *nats.Subscription, err error) {
	_, ok := conn.subscriptions.Load(subj)
	if ok {
		return nil, ErrTopicAlreadyRegistered
	}

	sub, err = js.JetStreamContext.SubscribeSync(subj, opts...)
	if err != nil {
		return nil, err
	}

	conn.subscriptions.Store(subj, sub)
	return
}

func (js JetStreamContext) QueueSubscribe(topic, queue string, cb nats.MsgHandler, opts ...nats.SubOpt) (sub *nats.Subscription, err error) {
	_, ok := conn.subscriptions.Load(topic)
	if ok {
		return nil, ErrTopicAlreadyRegistered
	}

	sub, err = js.JetStreamContext.QueueSubscribe(topic, queue, cb)
	if err != nil {
		return nil, err
	}

	conn.subscriptions.Store(topic, sub)
	return
}

func (js JetStreamContext) QueueSubscribeSync(subj, queue string, opts ...nats.SubOpt) (sub *nats.Subscription, err error) {
	_, ok := conn.subscriptions.Load(subj)
	if ok {
		return nil, ErrTopicAlreadyRegistered
	}

	sub, err = js.JetStreamContext.QueueSubscribeSync(subj, queue, opts...)
	if err != nil {
		return nil, err
	}

	conn.subscriptions.Store(subj, sub)
	return
}

func (js *JetStreamContext) ChanSubscribe(topic string, ch chan *nats.Msg, opts ...nats.SubOpt) (sub *nats.Subscription, err error) {
	_, ok := conn.subscriptions.Load(topic)
	if ok {
		return nil, ErrTopicAlreadyRegistered
	}

	sub, err = js.JetStreamContext.ChanSubscribe(topic, ch)
	if err != nil {
		return nil, err
	}

	conn.subscriptions.Store(topic, sub)
	return
}

func (js JetStreamContext) ChanQueueSubscribe(topic, queue string, ch chan *nats.Msg, opts ...nats.SubOpt) (sub *nats.Subscription, err error) {
	_, ok := conn.subscriptions.Load(topic)
	if ok {
		return nil, ErrTopicAlreadyRegistered
	}

	sub, err = js.JetStreamContext.ChanQueueSubscribe(topic, queue, ch)
	if err != nil {
		return nil, err
	}

	conn.subscriptions.Store(topic, sub)
	return sub, err
}

func (js JetStreamContext) PullSubscribe(subject, durable string, opts ...nats.SubOpt) (sub *nats.Subscription, err error) {
	_, ok := conn.subscriptions.Load(subject)
	if ok {
		return nil, ErrTopicAlreadyRegistered
	}

	sub, err = js.JetStreamContext.PullSubscribe(subject, durable)
	if err != nil {
		return nil, err
	}

	conn.subscriptions.Store(subject, sub)
	return
}
