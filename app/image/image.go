package image

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
)

type ErrorLine struct {
	Error       string      `json:"error"`
	ErrorDetail ErrorDetail `json:"errorDetail"`
}

type ErrorDetail struct {
	Message string `json:"message"`
}

func Build(path, registryURL, repoName string) (string, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return "", fmt.Errorf("docker client: %s", err)
	}

	imageName := ""
	if registryURL != "" {
		imageName = fmt.Sprintf("%s/%s", registryURL, repoName)
	} else {
		imageName = fmt.Sprintf("%s", repoName)
	}

	err = imageBuild(cli, path, imageName)
	if err != nil {
		return "", fmt.Errorf("internal image build: %s", err)
	}

	return imageName, nil
}

func imageBuild(dockerClient *client.Client, path, imageName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*600)
	defer cancel()

	tar, err := archive.TarWithOptions(path, &archive.TarOptions{})
	if err != nil {
		return fmt.Errorf("tar: %s", err)
	}

	opts := types.ImageBuildOptions{
		Dockerfile: "Dockerfile",
		Tags:       []string{imageName},
		Remove:     true,
	}
	res, err := dockerClient.ImageBuild(ctx, tar, opts)
	if err != nil {
		return fmt.Errorf("image build: %s", err)
	}

	defer res.Body.Close()

	err = print(res.Body)
	if err != nil {
		return fmt.Errorf("print: %s", err)
	}

	return nil
}

func print(rd io.Reader) error {
	var lastLine string

	scanner := bufio.NewScanner(rd)
	for scanner.Scan() {
		lastLine = scanner.Text()
		fmt.Println(scanner.Text())
	}

	errLine := &ErrorLine{}
	json.Unmarshal([]byte(lastLine), errLine)
	if errLine.Error != "" {
		return fmt.Errorf("json unmarshal: %s", errLine.Error)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scanner: %s", err)
	}

	return nil
}

func Push(imageName string) error {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return fmt.Errorf("docker client: %s", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*600)
	defer cancel()

	//authConfigBytes, err := json.Marshal(authConfig)
	//if err != nil {
	//	return "", fmt.Errorf("json marshal: %s", err)
	//}
	//authConfigEncoded := base64.URLEncoding.EncodeToString(authConfigBytes)

	tag := imageName
	opts := types.ImagePushOptions{RegistryAuth: "MTIz"} //RegistryAuth: authConfigEncoded}
	rd, err := cli.ImagePush(ctx, tag, opts)
	if err != nil {
		return fmt.Errorf("image push: %s", err)
	}

	defer rd.Close()

	err = print(rd)
	if err != nil {
		return fmt.Errorf("print: %s", err)
	}

	return nil
}
