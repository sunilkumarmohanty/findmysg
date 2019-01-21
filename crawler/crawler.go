package crawler

import (
	"flag"
	"io"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/rds"
	"go.uber.org/zap"
)

var logger *zap.Logger

func init() {
	var err error
	logger, err = zap.NewDevelopment()
	if err != nil {
		log.Fatalf("cannot initialize logger. Error : %v", err)
	}
}

// Crawl starts crawling AWS Resources
func Crawl() {
	var writer io.Writer
	var output = flag.String("output", "terminal", "Output type")
	var filepath = flag.String("file", "./output.csv", "Path of the file(CSV)")
	flag.Parse()

	// Get writer
	switch strings.ToUpper(*output) {
	case "CSV":
		logFilePath := *filepath
		f, _ := os.Create(logFilePath)
		defer f.Close()
		writer = f
	default:
		writer = os.Stdout
	}

	// TOOD : Do we need a new reporter instance

	// For development use os.Stdout
	results := make(chan *result)

	reporter := newReporter(results, writer)
	client := getAWSClient()
	go reporter.run()

	// Start each crawler one by one
	var wg sync.WaitGroup
	crawlers := []func(*awsClient, chan *result){
		crawlEC2,
		crawlRDS,
	}
	wg.Add(len(crawlers))

	for _, crawler := range crawlers {
		go func(crawl func(*awsClient, chan *result)) {
			crawl(client, results)
			wg.Done()
		}(crawler)
	}

	wg.Wait()
	close(results)
	<-reporter.done
}

func getAWSClient() *awsClient {
	client := &awsClient{}
	sess, _ := session.NewSession(
		&aws.Config{
			Region: aws.String("eu-west-1"),
		})
	client.ec2Conn = ec2.New(sess)
	client.rdsConn = rds.New(sess)
	return client
}

type awsClient struct {
	ec2Conn *ec2.EC2
	rdsConn *rds.RDS
}
