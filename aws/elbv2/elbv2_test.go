package elbv2

import (
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/dtan4/esnctl/aws/mock"
	"github.com/golang/mock/gomock"
)

func TestDetachInstance(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	api := mock.NewMockELBV2API(ctrl)
	api.EXPECT().DeregisterTargets(&elbv2.DeregisterTargetsInput{
		TargetGroupArn: aws.String("arn:aws:elasticloadbalancing:ap-northeast-1:012345678901:targetgroup/elasticsearch/0123abcd5678efab"),
		Targets: []*elbv2.TargetDescription{
			{
				Id: aws.String("i-1234abcd"),
			},
		},
	}).Return(&elbv2.DeregisterTargetsOutput{}, nil)

	client := &Client{
		api: api,
	}

	targetGroupARN := "arn:aws:elasticloadbalancing:ap-northeast-1:012345678901:targetgroup/elasticsearch/0123abcd5678efab"
	instanceID := "i-1234abcd"

	if err := client.DetachInstance(targetGroupARN, instanceID); err != nil {
		t.Errorf("error should not be raised: %s", err)
	}
}

func TestListTargetInstances(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	api := mock.NewMockELBV2API(ctrl)
	api.EXPECT().DescribeTargetHealth(&elbv2.DescribeTargetHealthInput{
		TargetGroupArn: aws.String("arn:aws:elasticloadbalancing:ap-northeast-1:012345678901:targetgroup/elasticsearch/0123abcd5678efab"),
	}).Return(&elbv2.DescribeTargetHealthOutput{
		TargetHealthDescriptions: []*elbv2.TargetHealthDescription{
			{
				Target: &elbv2.TargetDescription{
					Id: aws.String("i-1234abcd"),
				},
				TargetHealth: &elbv2.TargetHealth{
					State: aws.String("healthy"),
				},
			},
			{
				Target: &elbv2.TargetDescription{
					Id: aws.String("i-5678efab"),
				},
				TargetHealth: &elbv2.TargetHealth{
					State: aws.String("draining"),
				},
			},
		},
	}, nil)

	client := &Client{
		api: api,
	}

	targetGroupARN := "arn:aws:elasticloadbalancing:ap-northeast-1:012345678901:targetgroup/elasticsearch/0123abcd5678efab"
	expected := []string{
		"i-1234abcd",
		"i-5678efab",
	}

	got, err := client.ListTargetInstances(targetGroupARN)
	if err != nil {
		t.Errorf("error should not be raised: %s", err)
	}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("instance IDs does not match. expected: %q, got: %q", expected, got)
	}
}
