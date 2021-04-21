package controller

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"
	
	"github.com/frostyslav/lseg-demo/app/git"
	"github.com/frostyslav/lseg-demo/app/image"
	"github.com/frostyslav/lseg-demo/app/knative"
	"github.com/frostyslav/lseg-demo/app/model"
)

var registryURL string

func init() {
	registryURL = os.Getenv("DOCKER_REGISTRY")
}

func RunFromRepo(hashmap *model.Hash, url, tag, path, id string) error {
	log.Printf("repo url: %s, repo tag: %s, directory: %s, id: %s", url, tag, path, id)

	path, err := git.Clone(url, tag)
	if err != nil {
		hashmap.SetStatus(id, "git failed")
		log.Print(err)
		return err
	}
	log.Print(path)

	hashmap.SetStatus(id, "git done")

	splitted := strings.Split(url, "/")
	repoName := splitted[len(splitted)-1]

	err = buildAndCreate(hashmap, path, repoName, id)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("push and create: %s", err)
	}
	
	imageName, err := image.Build(path, registryURL, repoName)
	if err != nil {
		hashmap.SetStatus(id, "build failed")
		return fmt.Errorf("image build: %s", err)
	}
	
	log.Print(imageName)
	hashmap.SetStatus(id, "build done")
	
	err = image.Push(imageName)
	if err != nil {
		hashmap.SetStatus(id, "push failed")
		return fmt.Errorf("image build: %s", err)
	}
	hashmap.SetStatus(id, "image push done")

	return nil
}

func RunFromCode(hashmap *model.Hash, encodedValue, id, language string) error {
	log.Printf("Run from code")
	log.Printf("encoded value: %s", encodedValue)
	path := fmt.Sprintf("/tmp/%s", id)

	code, err := base64.StdEncoding.DecodeString(encodedValue)
	if err != nil {
		return fmt.Errorf("decode string: %s", err)
	}
	log.Printf("Decode")

	hashmap.SetStatus(id, "decoding done")

	err = os.Mkdir(path, 0755)
	if err != nil {
		return fmt.Errorf("mkdir: %s", err)
	}
	log.Printf("Mkdir")

	err = writeCode(code, path)
	if err != nil {
		return fmt.Errorf("write code: %s", err)
	}

	err = writeDockerfile(fmt.Sprintf("templates/%s/Dockerfile.tmpl", language), path, id)
	if err != nil {
		return fmt.Errorf("write dockerfile: %s", err)
	}

	err = buildAndCreate(hashmap, path, id, id)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("push and create: %s", err)
	}

	return nil
}

func buildAndCreate(hashmap *model.Hash, path, repoName, id string) error {
	imageName, err := image.Build(path, registryURL, repoName)
	if err != nil {
		hashmap.SetStatus(id, "build failed")
		return fmt.Errorf("image build: %s", err)
	}

	log.Print(imageName)
	hashmap.SetStatus(id, "build done")

	err = image.Push(imageName)
	if err != nil {
		hashmap.SetStatus(id, "push failed")
		return fmt.Errorf("image build: %s", err)
	}
	hashmap.SetStatus(id, "image push done")

	url, err := knative.Create(id, imageName)
	if err != nil {
		hashmap.SetStatus(id, "func create failed")
		return fmt.Errorf("knative create: %s", err)
	}

	hashmap.SetStatus(id, "finished")
	hashmap.SetURL(id, url)

	return nil
}

func writeCode(code []byte, directory string) error {
	f, err := os.Create(fmt.Sprintf("%s/main.go", directory))
	if err != nil {
		return fmt.Errorf("create file: %s", err)
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	b, err := w.Write(code)
	if err != nil {
		return fmt.Errorf("write to file: %s", err)
	}
	fmt.Printf("wrote %d bytes\n", b)

	w.Flush()

	return nil
}

func writeDockerfile(tmpl, directory, id string) error {
	f, err := os.Create(fmt.Sprintf("%s/Dockerfile", directory))
	if err != nil {
		return fmt.Errorf("create file: %s", err)
	}
	defer f.Close()

	td := Dockerfile{ID: id}

	fmt.Println("before template")

	t, err := template.ParseFiles(tmpl)
	if err != nil {
		return fmt.Errorf("template file: %s", err)
	}

	err = t.Execute(os.Stdout, td)
	if err != nil {
		return fmt.Errorf("template execute: %s", err)
	}

	err = t.Execute(f, td)
	if err != nil {
		return fmt.Errorf("template execute: %s", err)
	}

	return nil
}

type Dockerfile struct {
	ID string
}
