
// import (
// 	"context"
// 	"fmt"
// 	"log"
// 	"github.com/aws/aws-lambda-go/lambda"
// )

// type MyEvent struct {
// 	Name string `json:"name"`
// }

// func HandleRequest(ctx context.Context, event *MyEvent) (*string, error) {
// 	if event == nil {
// 		log.Println("received nil event")
// 		return nil, fmt.Errorf("received nil event")
// 	}
// 	message := fmt.Sprintf("Hello %s!", event.Name)
// 	log.Printf("message: %s", message)
// 	return &message, nil
// }

package main

import (
	"context"
	"fmt"
	"os"
	"github.com/jmoiron/sqlx"
	"github.com/olivere/env"
	"log/slog"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/tkhrk1010/go-samples/actor-model/event-sourcing/rmu/rmu"
)

func HandleRequest(ctx context.Context, event events.DynamoDBEvent) error {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	dbUrl := env.String("", "DATABASE_URL")
	if dbUrl == "" {
		panic("DATABASE_URL is required")
	}

	dataSourceName := fmt.Sprintf("%s?parseTime=true", dbUrl)
	db, err := sqlx.Connect("mysql", dataSourceName)
	if err != nil {
		panic(err.Error())
	}
	defer func(db *sqlx.DB) {
		if db != nil {
			err := db.Close()
			if err != nil {
				panic(err.Error())
			}
		}
	}(db)
	dao := rmu.NewWindSpeedDaoImpl(db)
	readModelUpdater := rmu.NewReadModelUpdater(&dao)
	// FIXME: なんか二重に呼んでて変
	lambda.Start(readModelUpdater.UpdateReadModel)
	return nil
}

func main() {
	lambda.Start(HandleRequest)
}
