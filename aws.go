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

// awsGetInstanceID by aws private dns name
func (e *EC2) awsGetInstanceID(privateDnsName string) (string, error) {
	describeInput := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("network-interface.private-dns-name"),
				Values: []*string{aws.String(privateDnsName)},
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
			if *instance.PrivateDnsName == privateDnsName {
				log.Debugf("found instance %v with id: %v", privateDnsName, *instance.InstanceId)
				return *instance.InstanceId, nil
			}
		}
	}
	return "", fmt.Errorf("failed to find instance with private dns name %v", privateDnsName)
}

// awsTerminateInstance terminate instance with private dns name
func (e *EC2) awsTerminateInstance(privateDnsName string) error {
	log.Infof("starting to terminate %v", privateDnsName)
	instanceID, err := e.awsGetInstanceID(privateDnsName)
	if err != nil {
		return err
	}

	terminationInput := &ec2.TerminateInstancesInput{
		InstanceIds: []*string{aws.String(instanceID)},
	}
	_, err = e.client.TerminateInstances(terminationInput)
	if err != nil {
		return fmt.Errorf("failed to terminate instance %v %v %v", privateDnsName, instanceID, err.Error())
	}
	log.Debugf("terminated instance %v", privateDnsName)
	return nil

}
