package crawler

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/lambda/lambdaiface"
)

type mockLambdaClient struct {
	resources []testResource
	lambdaiface.LambdaAPI
}

func (m *mockLambdaClient) ListFunctions(*lambda.ListFunctionsInput) (*lambda.ListFunctionsOutput, error) {
	output := &lambda.ListFunctionsOutput{}
	for _, resource := range m.resources {
		instance := &lambda.FunctionConfiguration{
			FunctionArn: aws.String(resource.id),
			VpcConfig:   &lambda.VpcConfigResponse{},
		}
		for _, sg := range resource.securityGroups {
			instance.VpcConfig.SecurityGroupIds = append(instance.VpcConfig.SecurityGroupIds, aws.String(sg))
		}

		output.Functions = append(output.Functions, instance)
	}
	return output, nil
}

func Test_crawlLambda(t *testing.T) {
	type args struct {
		client  *awsClient
		results chan *result
	}
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
			name:      "No Lambda Functions",
			resources: []testResource{},
		},
		{
			name: "Lambda Function with no security groups",
			resources: []testResource{
				testResource{
					id: "dummy_id_1",
				},
			},
		},
		{
			name: "Multiple Lambda Functions",
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
			lambdaConn: &mockLambdaClient{
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
			crawlLambda(client, chanResults)
			close(chanResults)
			<-done
			if ok, err := validateResults(results, tt.resources, "LAMBDA"); !ok {
				t.Errorf("Security groups did not match %v", err)
			}
		})
	}
}
