package queue

import (
	"SocialNetHL/models"
	"context"
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"sync"
	"time"
)

type Rabbit struct {
	postQ    *amqp.Queue
	friendsQ *amqp.Queue
	ch       *amqp.Channel
}

var (
	rabbitInstance *Rabbit
	rabbitOnce     sync.Once
)

func NewFeedQueue(connString string, postQueueName string, friendsQueueName string) (*Rabbit, error) {
	rabbitOnce.Do(func() {
		//conn, err := amqp.Dial("amqp://user:password@localhost:5672/")
		conn, err := amqp.Dial(connString)
		failOnError(err, "Failed to connect to RabbitMQ")
		//defer conn.Close()

		ch, err := conn.Channel()
		failOnError(err, "Failed to open a channel")
		//defer ch.Close()
		pq, err := ch.QueueDeclare(
			postQueueName, // name
			false,         // durable
			false,         // delete when unused
			false,         // exclusive
			false,         // no-wait
			nil,           // arguments
		)
		fq, err := ch.QueueDeclare(
			friendsQueueName, // name
			false,            // durable
			false,            // delete when unused
			false,            // exclusive
			false,            // no-wait
			nil,              // arguments
		)
		failOnError(err, "Failed to declare a queue")
		rabbitInstance = &Rabbit{postQ: &pq, friendsQ: &fq, ch: ch}
	})

	return rabbitInstance, nil
}

func (r *Rabbit) SendPostToFeed(ctx context.Context, post models.Post) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	body, _ := json.Marshal(post)
	err := r.ch.PublishWithContext(ctx,
		"",           // exchange
		r.postQ.Name, // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			//ContentType: "application/json",
			ContentType: "text/plain",
			Body:        body,
		})
	failOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s\n", body)
	return err
}

func (r *Rabbit) GetPostForFeed(ch chan models.Post) {
	msgs, err := r.ch.Consume(
		r.postQ.Name, // queue
		"",           // consumer
		true,         // auto-ack
		false,        // exclusive
		false,        // no-local
		false,        // no-wait
		nil,          // args
	)
	failOnError(err, "Failed to consume posts")

	for {
		select {
		case d := <-msgs:
			log.Printf("Received a message: %s", d.Body)
			var postMsg models.Post
			err = json.Unmarshal(d.Body, &postMsg)
			if err != nil {
				log.Printf("Cannot proceess message from posts queue, err: %v\n", err)
			}
			ch <- postMsg
		}
	}
}

func (r *Rabbit) SendFriendToUpdateFeed(ctx context.Context, req models.UpdateFeedRequest) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	body, err := json.Marshal(req)
	err = r.ch.PublishWithContext(ctx,
		"",              // exchange
		r.friendsQ.Name, // routing key
		false,           // mandatory
		false,           // immediate
		amqp.Publishing{
			//ContentType: "application/json",
			ContentType: "text/plain",
			Body:        body,
		})
	failOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s\n", body)
	return err
}

func (r *Rabbit) GetFriendsForUpdateFeed(ch chan models.UpdateFeedRequest) {
	msgs, err := r.ch.Consume(
		r.friendsQ.Name, // queue
		"",              // consumer
		true,            // auto-ack
		false,           // exclusive
		false,           // no-local
		false,           // no-wait
		nil,             // args
	)
	failOnError(err, "Failed to update feed")

	for {
		select {
		case d := <-msgs:
			log.Printf("Received a message: %s", d.Body)
			var friendId models.UpdateFeedRequest
			err = json.Unmarshal(d.Body, &friendId)
			if err != nil {
				log.Printf("Cannot proceess message from friends queue, err: %v\n", err)
			}
			ch <- friendId
		}
	}
}

func (r *Rabbit) Close() error {
	err := r.ch.Close()
	return err
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
