package crawler

import (
	"github.com/aws/aws-sdk-go/aws"
	"go.uber.org/zap"
)

func crawlElastiCache(client *awsClient, results chan *result) {
	res, err := client.elastiCacheConn.DescribeCacheClusters(nil)
	if err != nil {
		logger.Error("Unable to connect to EC2", zap.Error(err))
		return
	}
	for _, cluster := range res.CacheClusters {
		for _, sg := range cluster.SecurityGroups {
			results <- &result{
				ID:            cluster.CacheClusterId,
				Type:          aws.String("ElastiCache"),
				SecurityGroup: sg.SecurityGroupId,
			}
		}
	}
}
