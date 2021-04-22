# Serverless demo

This repository contains a demo of Kubernetes serverless implementation done with minikube, knative, kourier and ace editor.

## Prerequisites

* Docker
* Minikube

## Installation 

Clone repo
```bash
  git clone https://github.com/frostyslav/lseg-demo
  cd lseg-demo
```

Start minikube
```bash
  make start-minikube
```

Install local registry and create necessary port forwards
```bash
  make install-registry
```

Install Knative and Kourier and modify Knative config to be able to read from insecure local registry 
```bash
  make install-knative
```

Build frontend and backend images and push them to local registry 
```bash
  make build-dockerfiles
```

Install frontend
```bash
  make install-frontend
```

Install backend
```bash
  make install-backend
```

Make frontend accessible from the local machine
```bash
  kubectl port-forward --address 0.0.0.0 service/frontend 8090:8090
```

Make backend accessible from the local machine
```bash
  kubectl port-forward --address 0.0.0.0 service/backend 8080:8080
```

Make Kourier accessible from the local machine
```bash
  kubectl port-forward --address 0.0.0.0 --namespace kourier-system service/kourier 8000:80
```

## Usage/Examples

Open UI at http://127.0.0.1:8090.

Paste the code into text box, or provide the path to repository with existing Dockerfile and hit `Send Code` button.

helloworld-go is a good starting point: https://raw.githubusercontent.com/knative-sample/helloworld-go/master/helloworld.go

Wait for the page to return the URL.

curl the function to see if it is working
```bash
    curl -H "Host: <replace_me>.default.example.com" 127.0.0.1:8000
```

Wait a minute and check that function got destroyed
```bash
    kubectl get pods
```

curl the function again to recreate it
```bash
    curl -H "Host: <replace_me>.default.example.com" 127.0.0.1:8000
```

## API Reference

#### Create function

```http
  POST /func_create
```
Requires either `repo` dictionary with repo details and Dockefile present or `code` in base64.
| Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `repo` | `dictionary` | Repo details |
| `repo.url` | `string` | Repo URL |
| `repo.tag` | `string` | Repo tag or hash |
| `repo.path` | `string` | Path to the Dockerfile |
| `language` | `string` | Language of the codebase |
| `code` | `base64` | Code to create a function with |

#### Get item

```http
  GET /func_state
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `id`      | `string` | **Required**. Id of function to check |
