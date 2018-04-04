## Golang

## Event-Driven Architecture

## Asynchronous Messaging Queue Protocol

## Prerequisites
### Text Editor
It doesn't matter what you use, use what you feel comfortable in. Personally, I use visual code. It's free and super easy to setup. For installation instructions, go to: https://code.visualstudio.com/

### Docker 
I will be using RabbitMQ, by spinning up a Docker image locally on my machine. You don't need docker to run RabbitMQ, but I would recommend using a local Docker instance, at least for this short tutorial.

### Golang
Installation of Golang is nice and easy. Instructions and binaries can be found at the official Golang site: https://golang.org/doc/install

For this tutorial, I assume some basic understanding of programming and also some very basics of Golang. I will try to explain everything as well as possible, but of course, prior experience with Golang is an advantage.

## Spinning up RabbitMQ
With Docker, this is super simple. Simply type the following command in your terminal:

> docker run --detach --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:3-management

Simply explain we are running a docker image, running the container in --detach mode (as a background process) naming it rabbitmq and exposing ports 5672 and 15672. Finally, we specify the image that we want to pull and eventually run: `rabbitmq:3-management`. Once the docker container has started, you can open a browser and visit http://localhost:15672 to see the management interface. We won't be using the mangement interface, but it's a good way to confirm that everything is working as intended. 

## Writing the Code