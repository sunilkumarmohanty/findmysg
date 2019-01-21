package crawler

import (
	"github.com/aws/aws-sdk-go/aws"
	"go.uber.org/zap"
)

func crawlRDS(client *awsClient, results chan *result) {
	res, err := client.rdsConn.DescribeDBInstances(nil)
	if err != nil {
		logger.Error("Unable to connect to RDS", zap.Error(err))
		return
	}
	for _, instance := range res.DBInstances {
		for _, sg := range instance.VpcSecurityGroups {
			results <- &result{
				ID:            instance.DBInstanceArn,
				Type:          aws.String("RDS"),
				SecurityGroup: sg.VpcSecurityGroupId,
			}
		}
	}
}
