package controller

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	models "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/chat_list/models"
)

// Mock service to imitate the behavior of the service layer
type mockUsecase struct {
	returnChats []models.Chat // Пример чата
	err         error         // Ошибка, которую нужно вернуть
}

func (m *mockUsecase) GetChats(cookies []*http.Cookie) ([]models.Chat, error) {
	return m.returnChats, m.err
}

func TestHandler_Success(t *testing.T) {
	// Создаём моки данных
	expectedChats := []models.Chat{
		models.Chat{
			ChatId:      1,
			ChatName:    "Чат с пользователем 2",
			ChatType:    "personalMessages",
			UsersId:     []int{1, 2},
			LastMessage: "Когда за кофе?",
			AvatarURL:   "https://yandex-images.clstorage.net/bVLC53139/667e899dbzgI/Lec3og97oM2J8jgAbwmbs1UEQ_j2WQe6H7Tz0tGHlNUDiLp06xNO9LooehtZCLyucrVfOV3bNS1vNvr_fMoMLbniE8frC6CczUKcwc_ImueU0HKs18lHz490gERWwAOWtD4IttmRuiGPuG9PrfwYeJTUCT5PeyM6mMdYuXvvucreJwTwaprjvy1RSHHf2XlUxVagbjT_Z3s54KP1tiFyt1ZNSQbbE3rzTsqefIsOGIsUYXo-bNgSKZq1WSlJDWhoz9XEo5uL4K6ts3gATet5lVTXMkxWjOY7KIFDVJZywRQU_5pGKzNLwd44T06oTSlJF8LpfEw5cni4VlkKqL3pur7HtGJr6fJM_7N44i2JKwV39mKvRF_0WggAkvaH0qM3dF3rRmhTySM-KU0umyyYWiVw2A7POeGIuvaYeurMKuh91cfSm1uT3gyRmwOOe-pFVEUzzFSupwioEVOHBbBzZ2Tv-6XKcnrQn8kvjEot-5skUlksz1lzKIsGKytJHes6zpeFgbqbE_-fI5hA_TpqVeemkPzm7qcL-UCy5NVyUTYUrfi2K2JIwT85fHzLzVhIFiBIH-9b0Xt45-m5aDwpWu_VtENpCtPcbhE7gL9aisTlNiP8hE2kW1qwEiTkwILXxe9adbgRuvDt6U59ag4qWbVyOm7eG4Nq2QW5mTveSXpcdzTz-TmS_F8iiCEc6Or0JvZy3ldd9-q7MMEWl-FS5YT_GUTZQQqiztuNzust-xr0cJjMzfuSmqqlyvuIPKlJbadlI-l7IS6vUWixLdt6p1dG0dy3bMdYunPBRcWioYYX3BsV2fP7oq26jA6YvgmLxQKqnr_IU0tJJ4greQw5Sp0HhdFYu4DeXGFbEF3JOWb3ZIJvR851WQnhM2UX8VKF983KBAoDOOK9qp1_aC9YK0WCa_8P-UGZ2Eaqy5s9C4qsNofD6Rkg3qyRGCPveNplxMSyzjbctEk5w6C2k",
		},
		models.Chat{
			ChatId:      2,
			ChatName:    "МГТУ",
			ChatType:    "group",
			UsersId:     []int{1, 2, 3},
			LastMessage: "У нас еще вечером одна пара",
			AvatarURL:   "https://polymerbranch.com/wp-content/uploads/2023/11/mgtu.webp",
		},
	}

	mockSvc := &mockUsecase{returnChats: expectedChats, err: nil}
	chatController := NewChatController(mockSvc)

	// Создаём запрос
	req := httptest.NewRequest("GET", "/chats", nil)
	rec := httptest.NewRecorder()

	// Вызываем обработчик
	chatController.Handler(rec, req)

	// Проверяем статус-код
	res := rec.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, res.StatusCode)
	}

	// Проверяем содержимое JSON
	var chatsDTO models.ChatsDTO
	if err := json.NewDecoder(rec.Body).Decode(&chatsDTO); err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	if len(chatsDTO.Chats) != len(expectedChats) {
		t.Errorf("Expected %d chats, got %d", len(expectedChats), len(chatsDTO.Chats))
	}
}

func TestHandler_Unauthorized(t *testing.T) {
	mockSvc := &mockUsecase{returnChats: nil, err: errors.New("unauthorized")}
	chatController := NewChatController(mockSvc)

	// Создаём запрос
	req := httptest.NewRequest("GET", "/chats", nil)
	rec := httptest.NewRecorder()

	// Вызываем обработчик
	chatController.Handler(rec, req)

	// Проверяем статус-код
	res := rec.Result()
	if res.StatusCode != http.StatusUnauthorized {
		t.Errorf("Expected status code %d, got %d", http.StatusUnauthorized, res.StatusCode)
	}
}
