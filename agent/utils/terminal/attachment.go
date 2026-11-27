package terminal

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"sync"
	"time"

	"github.com/1Panel-dev/1Panel/agent/global"
	"github.com/gorilla/websocket"
)

var errAttachmentClosed = errors.New("terminal attachment is closed")

type attachment struct {
	sess    *Session
	ws      *websocket.Conn
	writeMu sync.Mutex
	done    chan struct{}
	closeMu sync.Once
	cursor  uint64
}

func (a *attachment) writeLocked(data []byte) error {
	select {
	case <-a.done:
		return errAttachmentClosed
	default:
		return a.ws.WriteMessage(websocket.TextMessage, data)
	}
}

func (a *attachment) close(code int, reason string) {
	a.closeMu.Do(func() {
		close(a.done)
		_ = a.ws.WriteControl(
			websocket.CloseMessage,
			websocket.FormatCloseMessage(code, reason),
			time.Now().Add(time.Second),
		)
		_ = a.ws.Close()
	})
}

func (a *attachment) Run() {
	defer func() {
		a.close(websocket.CloseNormalClosure, "")
		a.sess.detach(a, false, false, 0)
	}()

	for {
		_, data, err := a.ws.ReadMessage()
		if err != nil {
			return
		}

		var msg WsMsg
		if err := json.Unmarshal(data, &msg); err != nil {
			global.LOG.Errorf("decode terminal websocket message failed: %v", err)
			continue
		}

		switch msg.Type {
		case WsMsgCmd:
			input, err := base64.StdEncoding.DecodeString(msg.Data)
			if err != nil {
				global.LOG.Errorf("decode terminal input failed: %v", err)
				continue
			}
			a.sess.writeInput(input)
		case WsMsgResize:
			if msg.Cols > 0 && msg.Rows > 0 {
				a.sess.resize(msg.Cols, msg.Rows)
			}
		case WsMsgHeartbeat:
			a.writeMu.Lock()
			err := a.writeLocked(data)
			a.writeMu.Unlock()
			if err != nil {
				return
			}
		}
	}
}
