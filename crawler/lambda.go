package crawler

import (
	"github.com/aws/aws-sdk-go/aws"
	"go.uber.org/zap"
)

func crawlLambda(client *awsClient, results chan *result) {
	res, err := client.lambdaConn.ListFunctions(nil)
	if err != nil {
		logger.Error("Unable to connect to lambda", zap.Error(err))
		return
	}
	for _, instance := range res.Functions {
		if instance.VpcConfig != nil {
			for _, sg := range instance.VpcConfig.SecurityGroupIds {
				results <- &result{
					ID:            instance.FunctionArn,
					Type:          aws.String("LAMBDA"),
					SecurityGroup: sg,
				}
			}
		}
	}
}
