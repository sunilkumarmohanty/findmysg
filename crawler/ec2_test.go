package crawler

import (
	"log"
	"testing"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"

	"github.com/aws/aws-sdk-go/aws"
)

type mockEC2CLient struct {
	resources []testResource
	ec2iface.EC2API
}

func (m *mockEC2CLient) DescribeInstances(*ec2.DescribeInstancesInput) (*ec2.DescribeInstancesOutput, error) {
	output := &ec2.DescribeInstancesOutput{}
	output.Reservations = []*ec2.Reservation{
		&ec2.Reservation{},
	}
	for _, resource := range m.resources {
		instance := &ec2.Instance{
			State: &ec2.InstanceState{
				Name: aws.String("runnin"),
			},
			InstanceId: aws.String(resource.id),
		}
		for _, sg := range resource.securityGroups {
			instance.SecurityGroups = append(instance.SecurityGroups, &ec2.GroupIdentifier{GroupId: aws.String(sg)})
		}

		output.Reservations[0].Instances = append(output.Reservations[0].Instances, instance)
		log.Println(output)

	}
	return output, nil
}

func Test_crawlEC2(t *testing.T) {

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
			name:      "No EC2 Instances",
			resources: []testResource{},
		},
		{
			name: "EC2 Instance with no security groups",
			resources: []testResource{
				testResource{
					id: "dummy_id_1",
				},
			},
		},
		{
			name: "Multiple EC2 instances",
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
			ec2Conn: &mockEC2CLient{
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
			crawlEC2(client, chanResults)
			close(chanResults)
			<-done
			if ok, err := validateResults(results, tt.resources, "EC2"); !ok {
				t.Errorf("Security groups did not match %v", err)
			}
		})
	}
}
