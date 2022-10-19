package events

import (
	"bytes"
	"context"
	"cqrs/models"
	"encoding/gob"

	"github.com/nats-io/nats.go"
)

type NatsEventsStore struct {
	conn            *nats.Conn
	feedCreatedSub  *nats.Subscription
	feedCreatedChan chan CreatedFeedMessage
}

// construtor for the line 7
func NewNats(url string) (*NatsEventsStore, error) {
	conn, err := nats.Connect(url)

	if err != nil {
		return nil, err
	}

	return &NatsEventsStore{
		conn: conn,
	}, nil

}

// implemation of event close
func (n *NatsEventsStore) Close() {
	if n.conn != nil {
		n.conn.Close()
	}

	if n.feedCreatedSub != nil {
		n.feedCreatedSub.Unsubscribe()
	}

	close(n.feedCreatedChan)
}

// encode data msg
func (n *NatsEventsStore) encodeMessage(m Message) ([]byte, error) {
	b := bytes.Buffer{}
	//codificar m a bytes
	err := gob.NewEncoder(&b).Encode(m)

	if err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

// recive msg and encode msg
func (n *NatsEventsStore) PublishCreatedFeed(ctx context.Context, feed *models.Feed) error {

	msg := CreatedFeedMessage{
		Id:          feed.ID,
		Title:       feed.Title,
		Description: feed.Description,
		CreatedAt:   feed.CreatedAt,
	}

	data, err := n.encodeMessage(msg)

	if err != nil {
		return err
	}
	return n.conn.Publish(msg.Type(), data)
}

func (n *NatsEventsStore) decodeMsg(data []byte, m interface{}) error {
	b := bytes.Buffer{}
	b.Write(data)

	return gob.NewDecoder(&b).Decode(m)
}

func (n *NatsEventsStore) OnCreatedFeed(f func(CreatedFeedMessage)) (err error) {
	msg := CreatedFeedMessage{}

	n.feedCreatedSub, err = n.conn.Subscribe(msg.Type(), func(m *nats.Msg) {
		n.decodeMsg(m.Data, &msg)
		f(msg)
	})

	return
}

func (n *NatsEventsStore) SubscribeCreatedFeed(ctx context.Context) (<-chan CreatedFeedMessage, error) {
	msg := CreatedFeedMessage{}
	n.feedCreatedChan = make(chan CreatedFeedMessage, 64)
	ch := make(chan *nats.Msg, 64)
	var err error
	n.feedCreatedSub, err = n.conn.ChanSubscribe(msg.Type(), ch)

	if err != nil {
		return nil, err
	}

	go func() {
		for {
			select {
			case ms := <-ch:
				n.decodeMsg(ms.Data, &msg)
				n.feedCreatedChan <- msg
			}
		}

	}()
	/*return conection chan, chan with data, nil */
	return (<-chan CreatedFeedMessage)(n.feedCreatedChan), nil
}
