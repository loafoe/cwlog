package main

import (
    "flag"
    "fmt"
    "time"

    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/cloudwatchlogs"
    "github.com/philips-software/cwlog/hsdp"
    "github.com/spf13/pflag"
    "github.com/spf13/viper"
)

func GetLogEvents(sess *session.Session, limit int64, nextToken string, logGroupName string, logStreamName string) (*cloudwatchlogs.GetLogEventsOutput, error) {
    svc := cloudwatchlogs.New(sess)

    input := &cloudwatchlogs.GetLogEventsInput{
        Limit:         &limit,
        LogGroupName:  &logGroupName,
        LogStreamName: &logStreamName,
    }
    if nextToken != "" {
        input.NextToken = &nextToken
    }

    resp, err := svc.GetLogEvents(input)
    if err != nil {
        return nil, err
    }

    return resp, nil
}

func ListLogStreams(sess *session.Session, limit int64, logGroupName string) (*cloudwatchlogs.DescribeLogStreamsOutput, error) {
    svc := cloudwatchlogs.New(sess)
    resp, err := svc.DescribeLogStreams(&cloudwatchlogs.DescribeLogStreamsInput{
        Limit: &limit,
        LogGroupName: &logGroupName,
    })
    if err != nil {
        return nil, err
    }
    return resp, nil
}

func main() {
    viper.SetEnvPrefix("cwlog")
    viper.AutomaticEnv()
    flag.Int64("limit", 50, "The maximum number of events to retrieve")
    flag.String("group", "", "The name of the log group")
    flag.String("stream", "", "The name of the log stream")
    pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
    pflag.Parse()
    _ = viper.BindPFlags(pflag.CommandLine)

    limit := viper.GetInt64("limit")
    logGroupName := viper.GetString("group")
    logStreamName := viper.GetString("stream")

    listStreams := false

    if logGroupName == ""  {
        fmt.Println("You must supply a log group name (-g LOG-GROUP) and log stream name (-s LOG-STREAM)")
        return
    }
    if logStreamName == "" {
        fmt.Println("Listing stream names")
        listStreams = true
    }

    sess := session.Must(session.NewSessionWithOptions(session.Options{
        SharedConfigState: session.SharedConfigEnable,
    }))

    if listStreams {
        resp, err :=ListLogStreams(sess, limit, logGroupName)
        if err != nil {
            fmt.Println("Got error getting log streams:")
            fmt.Println(err)
            return
        }
        for _, stream := range resp.LogStreams {

            fmt.Println("  ", *stream.LogStreamName)
        }
        return
    }

    gotToken := ""
    nextToken := ""

    fmt.Println("Event messages for stream " + logStreamName + " in log group  " + logGroupName)

    queue, done, err := hsdp.StartLogger()
    if err != nil {
        fmt.Printf("Error setting up HSDP logger: %v\n", err)
        return
    }

    for {
        resp, err := GetLogEvents(sess, limit, nextToken, logGroupName, logStreamName)
        if err != nil {
            fmt.Println("Got error getting log events:")
            fmt.Println(err)
            break
        }
        for _, event := range resp.Events {
            queue <- *event
        }
        gotToken = nextToken
        nextToken = *resp.NextForwardToken
        if gotToken == nextToken {
            time.Sleep(1 * time.Second)
        }
    }
    done <- true
}
