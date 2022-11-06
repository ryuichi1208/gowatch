package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"gitlab.com/kuritayu/logs/lib"
	"log"
	"strconv"
	"time"
)

func main() {

	region := "ap-northeast-1"
	sess := session.Must(
		session.NewSession(&aws.Config{
			Region: aws.String(region),
		}))
	cloudwatch := logs.New(sess)

	// read log group
	log.Println("ALL Log Group")
	log.Println(cloudwatch.LogGroup().FindAll())

	// read log stream
	log.Println("ALL Log Stream")
	s := cloudwatch.LogStream("/aws/lambda/get-sqs")
	log.Println(s.FindAll())

	// read log event
	log.Println("FindAll Log Event")
	e := cloudwatch.LogEvent(
		"/aws/lambda/get-sqs",
		"2021/09/12/[$LATEST]7d7fe71b3aac4a649003f10d9085752b")
	events, _ := e.FindAll()
	for _, event := range events {
		fmt.Print(event.Message)
	}

	// create log stream
	log.Println("create log stream")
	err := s.Save("testStream")
	if err != nil {
		log.Println(err)
	}

	// create log event
	log.Println("create log event")
	e.LogStreamName = "testStream"
	err = e.Save("test message")
	if err != nil {
		log.Println(err)
	}

	// create log event (batch)
	log.Println("create log event (batch)")
	var messages []string
	for i := 0; i < 10; i++ {
		messages = append(messages, fmt.Sprintf("message %v", strconv.Itoa(i)))
	}
	err = e.BatchSave(messages)
	if err != nil {
		log.Println(err)
	}

	// query
	log.Println("query")
	query, err := cloudwatch.LogGroup().GrepByMessage(
		"/aws/lambda/get-sqs",
		"test",
		"2021-09-22T00:00:00+09:00",
		"2021-09-22T23:59:59+09:00",
	)

	if err != nil {
		log.Println(err)
	}
	log.Println(query)

	time.Sleep(5 * time.Second)
	result, err := cloudwatch.LogGroup().Result(query)
	if err != nil {
		log.Println(err)
	}

	for _, line := range result {
		fmt.Printf("%v %v %v\n", logs.Jst(line.Timestamp), line.LogStream, line.Message)
	}

}
