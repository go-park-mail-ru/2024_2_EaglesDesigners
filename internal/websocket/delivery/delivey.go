package delivery



func (h *MessageController) HandleConnection(w http.ResponseWriter, r *http.Request) {
	log := logger.LoggerWithCtx(r.Context(), logger.Log)
	// начало

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Delivery: error during connection upgrade:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer log.Println("Message delivery: websocket is closing")
	defer conn.Close()

	// Здесь можно хранить список старых сообщений (например, в массиве или в базе данных)
	messageChannel := make(chan models.WebScoketDTO, 10)
	errChannel := make(chan error, 10)
	closeChannel := make(chan bool, 1)

	defer func() {
		closeChannel <- true
		close(closeChannel)
	}()

	if err != nil {
		log.Println("Error reading message:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	go h.usecase.ScanForNewMessages(r.Context(), messageChannel, errChannel, closeChannel)

	// пока соеденено
	duration := 500 * time.Millisecond

	for {
		select {
		case err = <-errChannel:

			if err != nil {
				log.Printf("Delivery: ошибка в поиске новых сообщений: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		case message := <-messageChannel:
			// запись новых сообщений
			log.Println("Message delivery websocket: получены новые сообщения")

			conn.WriteJSON(message)

		default:
			time.Sleep(duration)
		}

	}
}
