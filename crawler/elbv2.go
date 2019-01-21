package crawler

import (
	"github.com/aws/aws-sdk-go/aws"
	"go.uber.org/zap"
)

func crawlELBV2(client *awsClient, results chan *result) {
	res, err := client.elbv2Conn.DescribeLoadBalancers(nil)
	if err != nil {
		logger.Error("Unable to connect to EC2", zap.Error(err))
		return
	}

	for _, lb := range res.LoadBalancers {
		for _, sg := range lb.SecurityGroups {
			results <- &result{
				ID:            lb.LoadBalancerArn,
				Type:          aws.String("ELBV2"),
				SecurityGroup: sg,
			}
		}
	}

}
