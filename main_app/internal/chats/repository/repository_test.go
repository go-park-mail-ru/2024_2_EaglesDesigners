package repository

// import (
// 	"context"
// 	"log"
// 	"testing"

// 	model "github.com/go-park-mail-ru/2024_2_EaglesDesigner/main_app/internal/chats/models"
// 	"github.com/google/uuid"
// 	"github.com/jackc/pgx/v4/pgxpool"
// )

// var repository ChatRepository
// var p *pgxpool.Pool

// func TestInit(t *testing.T) {
// 	pool, err := pgxpool.Connect(context.Background(), "postgres://postgres:postgres@localhost:5432/patefon")
// 	if err != nil {
// 		log.Printf("DB not connected: %v", err)
// 		t.Fail()
// 	}

// 	log.Println("DB connected")
// 	p = pool

// 	conn, err := p.Acquire(context.Background())
// 	if err != nil {
// 		t.Fail()
// 	}
// 	conn.QueryRow(context.Background(), `BEGIN;`)

// 	repository, err = NewChatRepository(pool)
// 	if err != nil {
// 		t.Fail()
// 	}
// }

// func TestCreateNewChat_Success(t *testing.T) {
// 	chat := model.Chat{
// 		ChatId:   uuid.New(),
// 		ChatName: "Oleg",
// 		ChatType: "personal",
// 	}
// 	repository.CreateNewChat(chat)
// }

// func TestEnd(t *testing.T) {
// 	conn, err := p.Acquire(context.Background())
// 	if err != nil {
// 		t.Fail()
// 	}

// 	conn.QueryRow(context.Background(), `ROLLBACK;`)
// 	p.Close()
// }
