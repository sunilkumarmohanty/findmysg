package crawler

import (
	"github.com/aws/aws-sdk-go/aws"
	"go.uber.org/zap"
)

func crawlEC2(client *awsClient, results chan *result) {
	res, err := client.ec2Conn.DescribeInstances(nil)
	if err != nil {
		logger.Error("Unable to connect to EC2", zap.Error(err))
		return
	}
	for _, rsv := range res.Reservations {
		for _, instance := range rsv.Instances {
			if instance.State != nil && *instance.State.Name != "terminated" {
				for _, sg := range instance.SecurityGroups {
					results <- &result{
						ID:            instance.InstanceId,
						Type:          aws.String("EC2"),
						SecurityGroup: sg.GroupId,
					}
				}
			}
		}
	}
}
