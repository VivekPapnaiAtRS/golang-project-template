package realtimesocketmanager

import (
	"VivekPapnaiAtRS/template/models"
	"encoding/json"
	"github.com/sirupsen/logrus"
)

func (c *RealtimeClient) processChatMessage(message models.Message) {
	messageDataBytes, err := json.Marshal(message.Data)
	if err != nil {
		logrus.Errorf("processChatMessage: unable to marshal chat message: %v", err)
		return
	}

	var chatMessage models.ChatMessageInfo
	err = json.Unmarshal(messageDataBytes, &chatMessage)
	if err != nil {
		logrus.Errorf("processChatMessage: unable to unmarshal chat message: %v", err)
		return
	}

	outboundMessage := map[string]interface{}{
		"messageText": chatMessage.Data,
	}

	outboundMessageBytes, err := json.Marshal(outboundMessage)
	if err != nil {
		logrus.Errorf("processChatMessage: unable to marshal oubound data chat message: %v", err)
		return
	}

	for i := range chatMessage.UserIDs {
		c.hub.getClients <- chatMessage.UserIDs[i]

		out, ok := <-c.hub.outGetClients
		if !ok || out.err != nil {
			logrus.Errorf("processChatMessage: user is not online")
			return
		}

		for i := range out.clients {
			out.clients[i].send <- outboundMessageBytes
		}
	}
}
