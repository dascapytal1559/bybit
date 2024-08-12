package bybit

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/gorilla/websocket"
)

// SubscribeBookticker :
func (s *V5WebsocketPublicService) SubscribeBookticker(
	key V5WebsocketPublicBooktickerParamKey,
	f func(V5WebsocketPublicBooktickerResponse) error,
) (func() error, error) {
	if err := s.addParamBooktickerFunc(key, f); err != nil {
		return nil, err
	}
	param := struct {
		Op   string        `json:"op"`
		Args []interface{} `json:"args"`
	}{
		Op:   "subscribe",
		Args: []interface{}{key.Topic()},
	}
	buf, err := json.Marshal(param)
	if err != nil {
		return nil, err
	}
	if err := s.writeMessage(websocket.TextMessage, buf); err != nil {
		return nil, err
	}
	return func() error {
		param := struct {
			Op   string        `json:"op"`
			Args []interface{} `json:"args"`
		}{
			Op:   "unsubscribe",
			Args: []interface{}{key.Topic()},
		}
		buf, err := json.Marshal(param)
		if err != nil {
			return err
		}
		if err := s.writeMessage(websocket.TextMessage, []byte(buf)); err != nil {
			return err
		}
		s.removeParamBooktickerFunc(key)
		return nil
	}, nil
}

// V5WebsocketPublicBooktickerParamKey :
type V5WebsocketPublicBooktickerParamKey struct {
	Symbol SymbolV5
}

// Topic :
func (k *V5WebsocketPublicBooktickerParamKey) Topic() string {
	return fmt.Sprintf("%s.%s", V5WebsocketPublicTopicBookticker, k.Symbol)
}

// V5WebsocketPublicBooktickerResponse :
type V5WebsocketPublicBooktickerResponse struct {
	Topic     string                          `json:"topic"`
	Type      string                          `json:"type"`
	TimeStamp int64                           `json:"ts"`
	Data      V5WebsocketPublicBooktickerData `json:"data"`
}

// V5WebsocketPublicBooktickerData :
type V5WebsocketPublicBooktickerData struct {
	Symbol      SymbolV5 `json:"s"`
	BidPrice    string   `json:"bp"`
	BidQuantity string   `json:"bq"`
	AskPrice    string   `json:"ap"`
	AskQuantity string   `json:"aq"`
	Timestamp   int64    `json:"t"`
}

// Key :
func (r *V5WebsocketPublicBooktickerResponse) Key() V5WebsocketPublicBooktickerParamKey {
	topic := r.Topic
	arr := strings.Split(topic, ".")
	if arr[0] != V5WebsocketPublicTopicBookticker.String() || len(arr) != 2 {
		return V5WebsocketPublicBooktickerParamKey{}
	}

	return V5WebsocketPublicBooktickerParamKey{
		Symbol: SymbolV5(arr[1]),
	}
}

// addParamTickerFunc :
func (s *V5WebsocketPublicService) addParamBooktickerFunc(
	key V5WebsocketPublicBooktickerParamKey,
	f func(V5WebsocketPublicBooktickerResponse) error,
) error {
	if _, exist := s.paramBooktickerMap[key]; exist {
		return errors.New("already registered for this key")
	}
	s.paramBooktickerMap[key] = f
	return nil
}

// removeParamTickerFunc :
func (s *V5WebsocketPublicService) removeParamBooktickerFunc(key V5WebsocketPublicBooktickerParamKey) {
	delete(s.paramBooktickerMap, key)
}

// retrieveTickerFunc :
func (s *V5WebsocketPublicService) retrieveBooktickerFunc(
	key V5WebsocketPublicBooktickerParamKey,
) (func(V5WebsocketPublicBooktickerResponse) error, error) {
	f, exist := s.paramBooktickerMap[key]
	if !exist {
		return nil, errors.New("func not found")
	}
	return f, nil
}
