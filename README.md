# Docker ELB

This program registers the current EC2 instance with the specified load balancer while the container is running. When the container exits, it unregisters the instance from the ELB.

Usage:

    docker run -d -p 5000:5000 --name=webapp training/webapp
    docker run -d \
      -v /var/run/docker.sock:/var/run/docker.sock \
      crewjam/dockerelb /dockerelb \
      -container=webapp \
      -elb=WebappLoadBalancer
