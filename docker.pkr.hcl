variable "name" {
  type    = string
  default = "queue-go"
}

packer {
  required_plugins {
    docker = {
      version = ">= 0.0.7"
      source  = "github.com/hashicorp/docker"
    }
  }
}

source "docker" "ubuntu" {
  image  = "golang:1.21.0-alpine3.18"
  commit = true
  changes = [
    "WORKDIR /go/app/cmd",
    "EXPOSE 6004",
    "ENTRYPOINT /go/app/cmd/cmd"
  ]
}

build {
  name = "chess/${var.name}"
  sources = [
    "source.docker.ubuntu"
  ]
  provisioner "shell" {
    inline = [
      "apk add git",
      "git clone https://github.com/cbotte21/${var.name} app",
      "cd app/cmd",
      "go install",
      "go build",
      "pwd"
    ]
  }
  provisioner "file" {
	source = ".env"
	destination = "/go/app/cmd/.env"
  }
  #provisioner "shell" {
   # inline = [
    #  "cd /go/app/cmd",
     # "echo \"port=6004\" >> .env",
      #"echo \"chess_addr=test\" >> .env"
    #]
  #}
  post-processors {
    post-processor "docker-tag" {
      repository = "chess/queue-go"
      tags       = ["0.1"]
    }
  }
}
