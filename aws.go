package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	log "github.com/sirupsen/logrus"
)

func newEC2() *EC2 {
	sess := session.Must(session.NewSession())
	a := EC2{
		client: ec2.New(sess),
	}
	return &a
}

// EC2 struct
type EC2 struct {
	client *ec2.EC2
}

func (e *EC2) awsGetInstanceID(nodename string) (string, error) {
	describeInput := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("network-interface.private-dns-name"),
				Values: []*string{aws.String(nodename)},
			},
		},
	}

	resp, err := e.client.DescribeInstances(describeInput)
	if err != nil {
		return "", fmt.Errorf("failed to get Instances %v", err.Error())
	}
	for index, resarvation := range resp.Reservations {
		log.Debugf("found instances in resp %v", len(resarvation.Instances))
		for _, instance := range resp.Reservations[index].Instances {
			if *instance.PrivateDnsName == nodename {
				return *instance.InstanceId, nil
			}
		}
	}
	return "", fmt.Errorf("failed to find instance with name %v", nodename)
}

func (e *EC2) awsTerminateInstance(nodename string) error {
	instanceID, err := e.awsGetInstanceID(nodename)
	if err != nil {
		return err
	}

	terminationInput := &ec2.TerminateInstancesInput{
		InstanceIds: []*string{aws.String(instanceID)},
	}
	_, err = e.client.TerminateInstances(terminationInput)
	if err != nil {
		return fmt.Errorf("failed to terminate instance %v %v %v", nodename, instanceID, err.Error())
	}
	return nil

}
