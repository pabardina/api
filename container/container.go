package container

import (
	"fmt"
	"github.com/fsouza/go-dockerclient"
	"os"
)

const (
	endpoint = "tcp://127.0.0.1:2375"
)

func ManageContainer() {
	client, err := docker.NewClient(endpoint)

	if err != nil {
		panic(err)
	}

	removeOpts := docker.RemoveContainerOptions{ID: "listener", Force: true}

	if err := client.RemoveContainer(removeOpts); err != nil {
		panic(err)
	}

	cmd := []string{
		"--access-token=" + os.Getenv("ACCESS_TOKEN"),
		"--access-token-secret=" + os.Getenv("ACCESS_TOKEN_SECRET"),
		"--consumer-key=" + os.Getenv("CONSUMER_KEY"),
		"--consumer-secret=" + os.Getenv("CONSUMER_SECRET"),
		"-a=" + os.Getenv("API_ADDRESS"),
		"-q=" + os.Getenv("QUEUE_ADDRESS"),
		"-l=5", "-r=5",
	}

	config := docker.Config{Cmd: cmd, Image: "pabardina/hirondelle-listener"}
	host := docker.HostConfig{NetworkMode: "host"}
	opts := docker.CreateContainerOptions{
		Name:       "listener",
		Config:     &config,
		HostConfig: &host,
	}
	container, err := client.CreateContainer(opts)

	if err != nil {
		panic(err)
	}

	if err := client.StartContainer(container.ID, &host); err != nil {
		panic(err)
	}

	fmt.Println(container.ID)

}
