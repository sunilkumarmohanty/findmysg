package crawler

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws"
)

type testResource struct {
	id             string
	securityGroups []string
}

func validateResults(results []*result, resources []testResource, resourceType string) (bool, error) {
	totalResSG := 0
	for _, resource := range resources {
		totalResSG += len(resource.securityGroups)
	}

	if len(results) != totalResSG {
		return false, errors.New("Length did not match")
	}
	// The results will be in the same order as the resources
	var i = 0
	for _, resource := range resources {
		for _, sg := range resource.securityGroups {
			if aws.StringValue(results[i].ID) != resource.id || aws.StringValue(results[i].SecurityGroup) != sg || aws.StringValue(results[i].Type) != resourceType {
				return false, errors.New("Value did not match")
			}
			i++
		}
	}

	return true, nil
}
