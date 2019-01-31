package crawler

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elasticache"
	"github.com/aws/aws-sdk-go/service/elasticache/elasticacheiface"
)

type mockElastiCacheClient struct {
	resources []testResource
	elasticacheiface.ElastiCacheAPI
}

func (m *mockElastiCacheClient) DescribeCacheClusters(*elasticache.DescribeCacheClustersInput) (*elasticache.DescribeCacheClustersOutput, error) {
	output := &elasticache.DescribeCacheClustersOutput{}
	for _, resource := range m.resources {
		cluster := &elasticache.CacheCluster{
			CacheClusterId: aws.String(resource.id),
		}
		for _, sg := range resource.securityGroups {
			cluster.SecurityGroups = append(cluster.SecurityGroups, &elasticache.SecurityGroupMembership{SecurityGroupId: aws.String(sg)})
		}
		output.CacheClusters = append(output.CacheClusters, cluster)
	}
	return output, nil
}

func Test_crawlElastiCache(t *testing.T) {

	tests := []struct {
		name      string
		resources []testResource
	}{
		{
			name: "Basic",
			resources: []testResource{
				testResource{
					id:             "dummy_id_1",
					securityGroups: []string{"sg1", "sg2"},
				},
			},
		},
		{
			name:      "No ElastiCache Clusters",
			resources: []testResource{},
		},
		{
			name: "ElastiCache Cluster with no security groups",
			resources: []testResource{
				testResource{
					id: "dummy_id_1",
				},
			},
		},
		{
			name: "Multiple ElastiCache Clusters",
			resources: []testResource{
				testResource{
					id:             "dummy_id_1",
					securityGroups: []string{"sg1", "sg2"},
				},
				testResource{
					id:             "dummy_id_2",
					securityGroups: []string{"sg2", "sg3"},
				},
			},
		},
	}
	for _, tt := range tests {
		chanResults := make(chan *result)
		done := make(chan struct{}, 1)
		client := &awsClient{
			elastiCacheConn: &mockElastiCacheClient{
				resources: tt.resources,
			},
		}
		t.Run(tt.name, func(t *testing.T) {
			var results []*result
			go func() {
				for res := range chanResults {
					results = append(results, res)
				}
				done <- struct{}{}
			}()
			crawlElastiCache(client, chanResults)

			close(chanResults)
			<-done
			if ok, err := validateResults(results, tt.resources, "ElastiCache"); !ok {
				t.Errorf("Security groups did not match %v", err)
			}
		})
	}
}
