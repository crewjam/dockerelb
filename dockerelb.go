package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/crowdmob/goamz/aws"
	"github.com/crowdmob/goamz/elb"
	docker "github.com/fsouza/go-dockerclient"
)

var containerID = flag.String("container", "",
	"the Docker container ID the container")

var loadBalancerName = flag.String("elb", "", "the ARN of the load balancer")

var dockerEndpoint = flag.String("docker", "unix:///var/run/docker.sock",
	"The path to the Docker endpoint")

func main() {
	flag.Parse()

	dockerClient, err := docker.NewClient(*dockerEndpoint)
	if err != nil {
		fmt.Printf("cannot connect to docker: %s\n", err)
		os.Exit(1)
	}

	instanceID := aws.InstanceId()
	if instanceID == "unknown" {
		fmt.Printf("cannot determine AWS instance ID. not running in EC2?\n")
		os.Exit(1)
	}

	awsAuth, err := aws.GetAuth("", "", "", time.Time{})
	if err != nil {
		fmt.Printf("cannot get AWS auth: %s\n", err)
		os.Exit(1)
	}

	elbConn := elb.New(awsAuth, aws.GetRegion(aws.InstanceRegion()))
	_, err = elbConn.RegisterInstancesWithLoadBalancer(
		[]string{instanceID},
		*loadBalancerName)
	if err != nil {
		fmt.Printf("cannot register instance: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("registered instance %s from ELB %s\n", instanceID,
		*loadBalancerName)

	// Wait for the container to exit
	_, err = dockerClient.WaitContainer(*containerID)
	if err != nil {
		fmt.Printf("docker: %s: waiting: %s\n", *containerID, err)
	} else {
		fmt.Printf("docker: %s exited\n", *containerID)
	}

	_, err = elbConn.DeregisterInstancesFromLoadBalancer(
		[]string{instanceID}, *loadBalancerName)
	if err != nil {
		fmt.Printf("cannot unregister instance: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("unregisterd instance %s from ELB %s\n",
		instanceID, *loadBalancerName)
}
