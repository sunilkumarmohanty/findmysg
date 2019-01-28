package crawler

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/aws/aws-sdk-go/service/rds/rdsiface"
)

type mockRDSClient struct {
	resources []testResource
	rdsiface.RDSAPI
}

func (m *mockRDSClient) DescribeDBInstances(*rds.DescribeDBInstancesInput) (*rds.DescribeDBInstancesOutput, error) {
	output := &rds.DescribeDBInstancesOutput{}
	for _, resource := range m.resources {
		instance := &rds.DBInstance{
			DBInstanceArn: aws.String(resource.id),
		}
		for _, sg := range resource.securityGroups {
			instance.VpcSecurityGroups = append(instance.VpcSecurityGroups, &rds.VpcSecurityGroupMembership{VpcSecurityGroupId: aws.String(sg)})
		}

		output.DBInstances = append(output.DBInstances, instance)
	}
	return output, nil
}

func Test_crawlRDS(t *testing.T) {

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
			name:      "No RDS Instances",
			resources: []testResource{},
		},
		{
			name: "RDS Instance with no security groups",
			resources: []testResource{
				testResource{
					id: "dummy_id_1",
				},
			},
		},
		{
			name: "Multiple RDS instances",
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
			rdsConn: &mockRDSClient{
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
			crawlRDS(client, chanResults)
			close(chanResults)
			<-done
			if ok, err := validateResults(results, tt.resources, "RDS"); !ok {
				t.Errorf("Security groups did not match %v", err)
			}
		})
	}
}
