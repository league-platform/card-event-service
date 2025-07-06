package handlers

import (
    "encoding/json"
    "fmt"
    "net/http"
    "time"

    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/dynamodb"
    "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
    "github.com/gofiber/fiber/v2"
)

type Card struct {
    MatchID   string    `json:"matchId"`
    Player    string    `json:"player"`
    Type      string    `json:"type"` // yellow/red
    Minute    int       `json:"minute"`
    CreatedAt time.Time `json:"createdAt"`
}

func CreateCard(c *fiber.Ctx) error {
    var card Card
    if err := c.BodyParser(&card); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
    }
    card.CreatedAt = time.Now()

    sess := session.Must(session.NewSession(&aws.Config{Region: aws.String("us-east-1")}))
    svc := dynamodb.New(sess)

    av, _ := dynamodbattribute.MarshalMap(card)
    input := &dynamodb.PutItemInput{
        TableName: aws.String("Cards"),
        Item:      av,
    }

    if _, err := svc.PutItem(input); err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save card"})
    }

    fmt.Printf("EVENT: card.issued -> %+v\n", card)
    return c.Status(fiber.StatusCreated).JSON(card)
}

func GetCards(c *fiber.Ctx) error {
    sess := session.Must(session.NewSession(&aws.Config{Region: aws.String("us-east-1")}))
    svc := dynamodb.New(sess)

    input := &dynamodb.ScanInput{
        TableName: aws.String("Cards"),
    }

    result, err := svc.Scan(input)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve cards"})
    }

    var cards []Card
    _ = dynamodbattribute.UnmarshalListOfMaps(result.Items, &cards)
    return c.JSON(cards)
}
