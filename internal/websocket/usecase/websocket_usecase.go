package usecase


func (u *MessageUsecaseImplm) publishMessageTochat(ctx context.Context, message models.Message) error {
	log := logger.LoggerWithCtx(ctx, logger.Log)
	err := u.redisClient.XAdd(ctx, &redis.XAddArgs{
		Stream: message.ChatId.String(),
		MaxLen: 0,
		ID:     "",
		Values: map[string]interface{}{
			NewMessage: message,
		},
	}).Err()
	log.Printf("publishMessageTochat: %v, %v", NewMessage, message)
	return err
}

// goBroker - брокер, который кинет сообщение из канала в нужный чат
func (u *MessageUsecaseImplm) goBroker(ctx context.Context) {
	log := logger.LoggerWithCtx(ctx, logger.Log)
	for {
		select {
		case message := <-u.messages:
			if _, ok := u.activeChats[message.ChatId]; ok {
				log.Println("Message usecase: добавление сообщения в redis")
				err := u.publishMessageTochat(ctx, message)
				if err != nil {
					log.Errorf("Message usecase -> goBroker: не удалось отправить в redis поток: %v", err)
				} else {
					log.Println("Message usecase -> goBroker: сообщение добавлено в redis поток")
				}
			}
		default:
		}
	}
}

// sendMessagesToUsers - отправляет каждоиму подписчику из activeUsers сообщение
func (u *MessageUsecaseImplm) sendMessagesToUsers(ctx context.Context, message models.Message, activeUsers map[string]bool) {
	log := logger.LoggerWithCtx(ctx, logger.Log)
	for userId := range activeUsers {
		err := u.redisClient.XAdd(ctx, &redis.XAddArgs{
			Stream: userId,
			MaxLen: 0,
			ID:     "",
			Values: map[string]interface{}{
				NewMessage: message,
			},
		}).Err()

		log.Println("Message usecase -> sendMessagesToUsers: сообщение отправлено", message)

		if err != nil {
			log.Errorf("Message usecase -> sendMessagesToUsers: не удалось отправить сообщение: %v", err)
		}
	}
}

// chatBroker кидает новые сообщения подписчикам
func (u *MessageUsecaseImplm) chatBroker(ctx context.Context, chatId uuid.UUID) {
	log := logger.LoggerWithCtx(ctx, logger.Log)
	log.Printf("Message usecase -> chatBroker: брокер создан для чата: %v", chatId)

	defer func() {
		u.activeChats[chatId] = false
		log.Printf("Message usecase -> chatBroker -> defer: брокер закрыт для чата: %v", chatId)
	}()

	duration := 500 * time.Millisecond
	activeUsersInChat := map[string]bool{}

	for {

		messages, err := u.redisClient.XRead(context.Background(), &redis.XReadArgs{
			Streams: []string{chatId.String(), "0"}, // Начинаем с самого начала ("0")
			Count:   5,                              // Получить 5 сообщений
			Block:   0,                              // Блокировать до появления новых сообщений
		}).Result()

		if err != nil {
			fmt.Errorf("Message usecase: Ошибка при чтении сообщений:", err)
			continue
		}

		var msgToDel []string
		// получаем новые сообщения в канал
		for _, message := range messages {
			fmt.Println("Стрим:", message.Stream)
			for _, msg := range message.Messages {
				if f, ok := msg.Values[FeatNewUser]; ok {
					activeUsersInChat[f.(string)] = true
					log.Printf("Message usecase -> chatBroker: добавлен подписчкик %v на чат %v", f.(string), chatId)
				} else if del, ok := msg.Values[DelUser]; ok {
					delete(activeUsersInChat, del.(string))
					log.Printf("Message usecase -> chatBroker: удалён подписчкик %v на чат %v", del.(string), chatId)
				} else if mesInterface, ok := msg.Values[NewMessage]; ok {
					var mes models.Message
					mes.UnmarshalBinary([]byte(mesInterface.(string)))
					log.Println(mes)

					u.sendMessagesToUsers(ctx, mes, activeUsersInChat)
				}
				msgToDel = append(msgToDel, msg.ID)
			}
		}
		// удаляем старые сообщения
		_, err = u.redisClient.XDel(context.Background(), chatId.String(), msgToDel...).Result()
		if err != nil {
			log.Errorf("Message usecase: не удалось удалить сообщения из redis: %v", err)
		}

		// если нет подписанных прользователей, то сворачиваемся
		if len(activeUsersInChat) == 0 {
			return
		}

		time.Sleep(duration)
	}

}

func (u *MessageUsecaseImplm) publishUserIntoChat(ctx context.Context, chatId uuid.UUID, userId uuid.UUID) error {
	log := logger.LoggerWithCtx(ctx, logger.Log)
	err := u.redisClient.XAdd(ctx, &redis.XAddArgs{
		Stream: chatId.String(),
		MaxLen: 0,
		ID:     "",
		Values: map[string]interface{}{
			FeatNewUser: userId.String(),
		},
	}).Err()

	log.Printf("publishUserIntoChat: %v, %v, %v", FeatNewUser, userId.String(), chatId)
	return err
}

// initChatsForUser - отправляет в чаты (брокерам сообщений по определенным чатам), что пользователь пришел и ему надо кидать сообщения
func (u *MessageUsecaseImplm) initChatsForUser(ctx context.Context, userId uuid.UUID) error {
	log := logger.LoggerWithCtx(ctx, logger.Log)
	log.Printf("Message usecase -> initChatsForUser: инициируем брокеров для пользователя: %v", userId)
	chats, err := u.chatRepository.GetUserChats(ctx, userId, 0)
	if err != nil {
		return fmt.Errorf("Не удалось получить чаты пользователя: %v", err)
	}
	for _, chat := range chats {
		err := u.publishUserIntoChat(ctx, chat.ChatId, userId)
		if err != nil {
			log.Errorf("Message usecase -> publishUserIntoChat: Не удалось добавить пользователя в чат: %v", err)
		}
		if _, ok := u.activeChats[chat.ChatId]; !ok {
			u.activeChats[chat.ChatId] = true
			log.Printf("Chat usecase -> initChatsForUser: создание брокера для чата: %v", chat.ChatId)
			go u.chatBroker(context.Background(), chat.ChatId)
		}
	}
	return nil
}

func (u *MessageUsecaseImplm) deleteUserFromChat(ctx context.Context, chatId uuid.UUID, userId uuid.UUID) error {
	err := u.redisClient.XAdd(ctx, &redis.XAddArgs{
		Stream: chatId.String(),
		MaxLen: 0,
		ID:     "",
		Values: map[string]interface{}{
			DelUser: userId.String(),
		},
	}).Err()

	return err
}

// deactivateUsrChats - отправляет в чаты (брокерам сообщений по определенным чатам), что пользователь ушёл
func (u *MessageUsecaseImplm) deactivateUsrChats(ctx context.Context, userId uuid.UUID) {
	log := logger.LoggerWithCtx(ctx, logger.Log)
	chats, err := u.chatRepository.GetUserChats(ctx, userId, 0)
	if err != nil {
		log.Println("Не удалось получить чаты пользователя")
		return
	}

	for _, chat := range chats {
		u.deleteUserFromChat(ctx, chat.ChatId, userId)
	}
}

// ScanForNewMessages сканирует redis stream с именем равным id пользователя
func (u *MessageUsecaseImplm) ScanForNewMessages(ctx context.Context, channel chan<- models.WebScoketDTO, res chan<- error, closeChannel <-chan bool) {
	log := logger.LoggerWithCtx(ctx, logger.Log)
	defer func() {
		close(channel)
		close(res)
	}()
	user, ok := ctx.Value(auth.UserKey).(jwt.User)
	log.Println(user)
	if !ok {
		return
	}
	u.activeUsers[user.ID] = true

	// убираем пользователя из списка активных юзеров
	defer func() {
		u.activeUsers[user.ID] = false
		u.deactivateUsrChats(context.Background(), user.ID)
	}()

	// создаем рутины на чаты, если еще не существуют
	err := u.initChatsForUser(ctx, user.ID)
	if err != nil {
		log.Errorf("Message usecase: не удалось инициировать чаты пользоватея: %v", err)
		res <- err
		return
	}

	log.Println("Message usecase: начат поиск новых сообщений")

	duration := 500 * time.Millisecond

	for {
		select {
		case <-closeChannel:
			log.Println("Message usecase: scanning stoped")
			return
		default:
			time.Sleep(duration)
			// Чтение сообщений из Stream
			messages, err := u.redisClient.XRead(context.Background(), &redis.XReadArgs{
				Streams: []string{user.ID.String(), "0"}, // Начинаем с самого начала ("0")
				Count:   5,                               // Получить 5 сообщений
				Block:   0,                               // Блокировать до появления новых сообщений
			}).Result()

			if err != nil {
				fmt.Errorf("Message usecase: Ошибка при чтении сообщений: %v", err)
				continue
			}
			var msgIds []string
			// получаем новые сообщения в канал

			for _, message := range messages {
				for _, msg := range message.Messages {
					mes, err := u.makeWebSocketDTO(msg.Values)
					if err != nil {
						log.Errorf("Message usecase -> websocket: не удалось получить данные и канала: %v", err)
						continue
					}
					if mes.MsgType == models.FeatUserInChat {
						// создаем рутины на чаты, если еще не существуют
						err := u.initChatsForUser(ctx, user.ID)
						if err != nil {
							log.Errorf("Message usecase: не удалось инициировать чаты пользоватея: %v", err)
							res <- err
							return
						}
					}

					channel <- mes

					msgIds = append(msgIds, msg.ID)
				}
			}
			_, err = u.redisClient.XDel(context.Background(), user.ID.String(), msgIds...).Result()
			if err != nil {
				log.Errorf("Message usecase: не удалось удалить сообщения из redis: %v", err)
			}
		}
	}
}

func (u *MessageUsecaseImplm) makeWebSocketDTO(newInfo map[string]interface{}) (models.WebScoketDTO, error) {
	log.Println("че-то читаем")
	if chatInterface, ok := newInfo[FeatNewUser]; ok {
		var chat chatModel.ChatDTOOutput
		chat.UnmarshalBinary([]byte(chatInterface.(string)))
		return models.WebScoketDTO{
			MsgType: models.FeatUserInChat,
			Payload: chat,
		}, nil
	} else if del, ok := newInfo[DelUser]; ok {
		return models.WebScoketDTO{
			MsgType: models.FeatUserInChat,
			Payload: del,
		}, nil
	} else if messageInterface, ok := newInfo[NewMessage]; ok {
		var message models.Message
		message.UnmarshalBinary([]byte(messageInterface.(string)))
		return models.WebScoketDTO{
			MsgType: models.NewMessage,
			Payload: message,
		}, nil
	}

	return models.WebScoketDTO{}, errors.New("No new messages")
}

func (u *MessageUsecaseImplm) GetOnlineUsers() map[uuid.UUID]bool {
	return u.activeUsers
}
