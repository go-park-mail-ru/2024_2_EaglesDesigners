package usecase

import "log"

// consumeMessages принимает информацию о сообщениях (добавление/изменение/удаление)
func (w *WebsocketUsecase) consumeMessages() {
	for {
		messages, err := w.ch.Consume(
			"message", // queue
			"",        // consumer
			true,      // auto-ack
			false,     // exclusive
			false,     // no-local
			false,     // no-wait
			nil,       // args
		)

		if err != nil {
			log.Fatalf("failed to register a consumer. Error: %s", err)
		}
		for message := range messages {
			log.Printf("received a message: %s", message.Body)
		}

	}
}
