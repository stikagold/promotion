package Broker

import (
	"context"
	"errors"
	"fmt"
	"github.com/rabbitmq/amqp091-go"
	"sync"
)

type RabbitMQ struct {
	AutoConnect bool   `json:"auto_connect"`
	Protocol    string `json:"protocol"`
	Auth        bool   `json:"auth"`
	Host        string `json:"host"`
	Port        string `json:"port"`
	User        string `json:"user"`
	Password    string `json:"password"`
	QueueName   string `json:"queue_name"`
	Exchange    string `json:"exchange"`
	RetryCount  int    `json:"retry_count"`

	// Operational fields
	mx         sync.Mutex
	Connection *amqp091.Connection
	Channel    *amqp091.Channel
	Queue      amqp091.Queue
}

func (rmq *RabbitMQ) GetURL() string {
	if rmq.Auth == true {
		return rmq.Protocol + "://" + rmq.User + ":" + rmq.Password + "@" + rmq.Host + ":" + rmq.Port + "/"
	}
	return rmq.Protocol + "://" + rmq.Host + ":" + rmq.Port + "/"
}

func (rmq *RabbitMQ) Initial() error {
	if rmq.IsEmpty() == false {
		var err error
		fmt.Println("Host of rabbitmq is: " + rmq.GetURL())
		rmq.Connection, err = amqp091.Dial(rmq.GetURL())
		if err != nil {
			_ = rmq.Connection.Close()
			return err
		}
		rmq.Channel, err = rmq.Connection.Channel()
		if err != nil {
			_ = rmq.Connection.Close()
			_ = rmq.Channel.Close()
			return err
		}
		fmt.Println("Name of queue should be: " + rmq.QueueName)
		rmq.Queue, err = rmq.Channel.QueueDeclare(
			rmq.QueueName,
			false,
			false,
			false,
			false,
			nil,
		)
		return err
	}
	return nil
}

func (rmq *RabbitMQ) IsEmpty() bool {
	return !rmq.AutoConnect == true
}

func (rmq *RabbitMQ) Close() error {
	_ = rmq.Connection.Close()
	_ = rmq.Channel.Close()
	return nil
}

func (rmq *RabbitMQ) Dispatch(b []byte, qName string) error {
	rmq.mx.Lock()
	defer rmq.mx.Unlock()

	var err error
	if qName == "" && rmq.QueueName == "" {
		return errors.New("no defined queue name")
	}

	if qName != "" && qName != rmq.QueueName {
		_, err = rmq.Channel.QueueDeclare(
			qName,
			false,
			false,
			false,
			false,
			nil,
		)
	}

	if qName == "" {
		qName = rmq.QueueName
	}

	fmt.Println("Start dispatching to rabbitmq")
	err = rmq.Channel.PublishWithContext(context.TODO(), "",
		qName,
		false,
		false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        b,
		})
	return err
}
