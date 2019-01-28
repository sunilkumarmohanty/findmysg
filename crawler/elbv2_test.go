package crawler

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/aws/aws-sdk-go/service/elbv2/elbv2iface"

	"github.com/aws/aws-sdk-go/aws"
)

type mockELBV2CLient struct {
	resources []testResource
	elbv2iface.ELBV2API
}

func (m *mockELBV2CLient) DescribeLoadBalancers(*elbv2.DescribeLoadBalancersInput) (*elbv2.DescribeLoadBalancersOutput, error) {
	output := &elbv2.DescribeLoadBalancersOutput{}
	for _, resource := range m.resources {
		instance := &elbv2.LoadBalancer{
			LoadBalancerArn: aws.String(resource.id),
		}
		for _, sg := range resource.securityGroups {
			instance.SecurityGroups = append(instance.SecurityGroups, aws.String(sg))
		}

		output.LoadBalancers = append(output.LoadBalancers, instance)
	}
	return output, nil
}

func Test_crawlELBV2(t *testing.T) {

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
			name:      "No ELBV2 Instances",
			resources: []testResource{},
		},
		{
			name: "ELBV2 Instance with no security groups",
			resources: []testResource{
				testResource{
					id: "dummy_id_1",
				},
			},
		},
		{
			name: "Multiple ELBV2 instances",
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
			elbv2Conn: &mockELBV2CLient{
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
			crawlELBV2(client, chanResults)
			close(chanResults)
			<-done
			if ok, err := validateResults(results, tt.resources, "ELBV2"); !ok {
				t.Errorf("Security groups did not match %v", err)
			}
		})
	}
}
