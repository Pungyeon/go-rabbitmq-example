# Using RabbitMQ with Golang

## Golang
So, why Golang? Why not Java or Python? Java has a much larger user base and Python is a much simpler language to write. So why choose Golang? Well. Golang is a modern language, which has great support for modern software architecture. Golang is a very small language, it's compiled (not transpiled and not run by the JVM) so, building Docker containers with Golang is a match made in heaven. Small Docker images, with performance similar to languages like C. So, this also makes Golang a great language for writing microservices... And, you can't say microservices without saying event-driven architecture! So, let's write a simple Golang program, to exemplify using Golang with RabbitMQ to support and event driven architecture.

There are many other fantastic features of Golang, but I won't go too much into detail. If you are interested, I would recommend watching this short interview with Nic Jackson from Hashicorp: https://www.youtube.com/watch?v=qlwp0mHFLHU

## Event-Driven Architecture
Event Driven Architecture has been popular looooooooooong before Microservices, but now that Microservices are all the talk, so is EDA. Essentially, EDA is a pattern for communication of state. It's been immensely popular in the financial industry for decades, as the pattern is particularly suited for handling transaction state. The reason why it has become so attached to the conversation of Microservices, is that in Microservice Architecture, you want everything to be loosely coupled. Essentially, you don't want one service to be attached to another. You want to avoid situations in which you change something in one service and then must make a corresponding change to one or all other services.

Let's think of an HTTP service, in which we are communicating with one or more services. Who decides who receives data? It's the HTTP service, which directly calls each and every one of those services. So... what happens if we create a new service that also needs this data? We would have to ask whoever is maintaining the HTTP service, if they could make sure, that our service also could receive this data. 

However, in an EDA, we don't need to contact the HTTP service owners at all. An EDA typically works in a publish/subscribe pattern. Simply explained, a publisher sends a message to a message broker, who will appropriately deliver the messages to all services who are subscribed. So, if we need to create a new service, we simply tell the message broker that we are subscribing to these messages/events. The HTTP service guys don't need to know about us, and we don't need to talk to them (on a non-technical level, socially, this is also typically regarded as a win). 

Now, there are many other advantages of EDA. But I will leave that to others to explain.


## Asynchronous Messaging Queue Protocol
The Asynchronous Messaging Queue Protocol started development in 2003, initiated by JPMorgan Chase. The project soon caught on and became a open-source project involving some of the largest banks and technology companies (Bank of America, Barclays, Microsoft, Cisco etc.) Essentially, the project was meant to create an open standard, to improve transactions, with a focus on the financial industry. Therefore, there was a huge backing by the banking industry to develop AMQP, making it extremely efficient and reliable. AMQP relies on messaging queues to handle communication, in a so called publish/subscribe architecture. The most common pattern of implementing this, the pattern this tutorial will be looking at, is the `topic exchange`. Essentially, a publisher sends a message to an `exchange` which will distribute messages to queues, based on a `topic`. The subscriber(s) will define a `queue` and tell the exchange which `topics` they are interested in.

As an example: If we, as a subscriber define a queue in which we define to be interested in all messages with the topic `apple`, if a publisher sends a message with `apple` we will receive that message. Even further, we can define that we are interested in sub topics, which is a typical implementation for logging. So, as an example, I might have a subscriber who is listening for `log.ERROR`  and `log.CRITICAL`, but have another subscriber who is interested in all log `log.*`. In other words, it's possible to listen based on routing keys (which work like search filters). This is super neat and something that we will explore further in this tutorial, using RabbitMQ, which implements AMQP. If you wish to use something other than RabbitMQ, then your in luck, because any asynchronous message system that supports AMQP will work with the code in this tutorial.

So, AMQP seems rather simple, right? It is, and that is why it's so great. We define a `publisher` who sends a message with a specified `topic` to an `exchange`. The `exchange` will determine whom to send these message to, based on `subscribers` `topic` routing keys. 

## Prerequisites
### Text Editor
It doesn't matter what you use, use what you feel comfortable in. Personally, I use visual code. It's free and super easy to setup. For installation instructions, go to: https://code.visualstudio.com/

### Docker 
I will be using RabbitMQ, by spinning up a Docker image locally on my machine. You don't need docker to run RabbitMQ, but I would recommend using a local Docker instance, at least for this short tutorial. Docker installation instructions can be found here: https://docs.docker.com/install/

### Golang
Installation of Golang is nice and easy. Instructions and binaries can be found at the official Golang site: https://golang.org/doc/install

For this tutorial, I assume some basic understanding of programming and also some very basics of Golang. I will try to explain everything as well as possible, but of course, prior experience with Golang is an advantage.

## Spinning up RabbitMQ
With Docker, this is super simple. Simply type the following command in your terminal:

> docker run --detach --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:3-management

We are running a docker image, running the container in --detach mode (as a background process) naming it rabbitmq and exposing ports 5672 and 15672. Finally, we specify the image that we want to pull and eventually run: `rabbitmq:3-management`. Once the docker container has started, you can open a browser and visit http://localhost:15672 to see the management interface. We won't be using the mangement interface, but it's a good way to confirm that everything is working as intended.

## Writing the Code
If you want to skip writing the code, but instead just want to read through and run the programs yourself. You can get the code from: https://github.com/Pungyeon/go-rabbitmq-example

So for this tutorial, we will be writing two really simple programs, to illustrate how services can communicate via. RabbitMQ. Our final project will look something like this:
go-rabbit-mq/

----./consumer

----./lib

--------./event

----./sender

We will be creating a `consumer` service, which will subscribe to our topics and we will define a `sender`service, which will publish random events to the exchange. Our `lib` folder, will hold some common configurations for both our consumer and sender. Before we begin, you will have to get the dependency for amqp:

`go get github.com/streadway/amqp`

But that's it, now we are ready to write some code.

### Event Queue
All files in this section will be placed in `lib/event`.

#### ./lib/event/event.go
First, we will write our library consisting of queue declaration and our structs for consumer and emitter. We will however start with some simple queue and exchange declaration:
```go
package event

import (
	"github.com/streadway/amqp"
)

func getExchangeName() string {
	return "logs_topic"
}

func declareRandomQueue(ch *amqp.Channel) (amqp.Queue, error) {
	return ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
}

func declareExchange(ch *amqp.Channel) error {
	return ch.ExchangeDeclare(
		getExchangeName(), // name
		"topic",           // type
		true,              // durable
		false,             // auto-deleted
		false,             // internal
		false,             // no-wait
		nil,               // arguments
	)
}
```

In this file, we are defining three static methods. The `getExchangeName` function simply returns the name of our exchange. It isn't necessary, but nice for this tutorial, to make it simple to change your topic name. More interesting is the `declareRandomQueue` function. This function will create a nameless queue, which RabbitMQ will assign a random name, we don't want to worry about this and that is why we are letting RabbitMQ worry about it. The queue is also defined as `exclusive`, which means that when defined only one subscriber can be subscribed to this queue. The last function that we have declared is `declareExchange` which will declare an exchange, as the name suggests. This function is idempotent, so if the exchange already exists, no worries, it won't create duplicates. However, if we were to change the type of the Exchange (to direct or fanout), then we would have to either delete the old exchange or find a new name, as you cannot overwrite exchanges. The topic type is what enables us to publish an event with a topic such as `log.WARN`, which the subscribers can specify in their search queries.

*NOTE: You might have noticed that both functions need an amqp.Channel struct. This is simply a pointer to an AMQP connection channel. We will explain this a little better later*

### ./lib/event/emitter.go
Next, we will define our publisher. I have chosen to call it emitter, because I wanted to add extra confusion...  Either way... This is our publisher. Which will publish, or in our case emit, events. 

```go
package event

import (
	"log"

	"github.com/streadway/amqp"
)

// Emitter for publishing AMQP events
type Emitter struct {
	connection *amqp.Connection
}

func (e *Emitter) setup() error {
	channel, err := e.connection.Channel()
	if err != nil {
		panic(err)
	}

	defer channel.Close()
	return declareExchange(channel)
}

// Push (Publish) a specified message to the AMQP exchange
func (e *Emitter) Push(event string, severity string) error {
	channel, err := e.connection.Channel()
	if err != nil {
		return err
	}

	defer channel.Close()

	err = channel.Publish(
		getExchangeName(),
		severity,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(event),
		},
	)
	log.Printf("Sending message: %s -> %s", event, getExchangeName())
	return nil
}

// NewEventEmitter returns a new event.Emitter object
// ensuring that the object is initialised, without error
func NewEventEmitter(conn *amqp.Connection) (Emitter, error) {
	emitter := Emitter{
		connection: conn,
	}

	err := emitter.setup()
	if err != nil {
		return Emitter{}, err
	}

	return emitter, nil
}
```

At the very top of our code, we are defining our Emitter struct (a class), which contains an amqp.Connection.

**setup** - Makes sure that the exchange that we are sending messages to actually exists. We do this by retreiving a channel from our connection pool and calling the idempotent declareExchange function from our event.go file.

**Push** - Sends a message to our exchange. First we get a new `channel` from our connection pool and if we receive no errors when doing so, we publish our message. The function takes two input parameters `event` and `severity`; `event` is the message to be sent and severity is our logging serverity, which will define which messages are received by which subscribers, based on their search queries. 

**NewEventEmitter** - Will return a new Emitter, or an error, making sure that the connection is established to our AMQP server.

The last bit of code to write for our library, is our consumer struct and right away we can see that it is somewhat similar to our emitter struct.
#### ./lib/event/consumer.go
```go
package event

import (
	"log"

	"github.com/streadway/amqp"
)

// Consumer for receiving AMPQ events
type Consumer struct {
	conn      *amqp.Connection
	queueName string
}

func (consumer *Consumer) setup() error {
	channel, err := consumer.conn.Channel()
	if err != nil {
		return err
	}
	return declareExchange(channel)
}

// NewConsumer returns a new Consumer
func NewConsumer(conn *amqp.Connection) (Consumer, error) {
	consumer := Consumer{
		conn: conn,
	}
	err := consumer.setup()
	if err != nil {
		return Consumer{}, err
	}

	return consumer, nil
}

// Listen will listen for all new Queue publications
// and print them to the console.
func (consumer *Consumer) Listen(topics []string) error {
	ch, err := consumer.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := declareRandomQueue(ch)
	if err != nil {
		return err
	}

	for _, s := range topics {
		err = ch.QueueBind(
			q.Name,
			s,
			getExchangeName(),
			false,
			nil,
		)

		if err != nil {
			return err
		}
	}

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	forever := make(chan bool)
	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Printf("[*] Waiting for message [Exchange, Queue][%s, %s]. To exit press CTRL+C", getExchangeName(), q.Name)
	<-forever
	return nil
}
```

At the very top we define that our `Consumer` struct defines a connection to our AMQP server and a queueName. The queue name will store the randomly generated name of our declared nameless queue. We will use this for telling RabbitMQ that we want to bind/listen to this particular queue for messages.

**setup()** - We ensure that the exchange is declared, just like we do in our Emitter struct.

**NewConsumer()** - We return a new Consumer or an error, ensuring that everything went well connecting to our AMQP server.

**Listen** - We get a new channel from our connection pool. We declare our nameless queue and then we iterate over our input `topics`, which is just an array of strings. For each string in topics, we will bind our search query to the queue. As an example, this could be `log.WARN` and `log.ERROR`. Lastly, we will invoke the Consume function (to start listening on the queue) and define that we will iterate over all messages received from the queue and print out these message to the console. 

The `forever` channel that we are making on line #69, and sending output from on line #77, is just a dummy. This is a simple way of ensuring a program will run forever. Essentially, we are defining a channel, which we will wait for until it receives input, but we will never actually send it any input. It's a bit dirty, but for this tutorial it will suffice. 

### Consumer Service
All files in this section will be placed in the `consumer` folder.

#### ./consumer/consumer.go
```go
package main

import (
	"os"

	"github.com/Pungyeon/go-rabbitmq-example/lib/event"
	"github.com/streadway/amqp"
)

func main() {
	connection, err := amqp.Dial("amqp://guest:guest@localhost:5672")
	if err != nil {
		panic(err)
	}
	defer connection.Close()

	consumer, err := event.NewConsumer(connection)
	if err != nil {
		panic(err)
	}
	consumer.Listen(os.Args[1:])
}
```

As can be seen this is a really simple program which creates a connection to our docker instance of RabbitMQ, passes this connection to our `NewConsumer` function and then calls the `Listen` method, passing all the input arguments from the command line. Once we have written this code we can open up a few terminals to start up a few consumers:

>#t1> go run consumer.go log.WARN log.ERROR

>#t2> go run consumer.go log.INFO log.WARN

### Emitter Service

Now for the very last bit of this tutorial. Publishing our messages to the queue:

```go
package main

import (
	"fmt"
	"os"

	"github.com/Pungyeon/go-rabbitmq-example/lib/event"
	"github.com/streadway/amqp"
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672")
	if err != nil {
		panic(err)
	}

	emitter, err := event.NewEventEmitter(conn)
	if err != nil {
		panic(err)
	}

	for i := 1; i < 10; i++ {
		emitter.Push(fmt.Sprintf("[%d] - %s", i, os.Args[1]), os.Args[1])
	}
}
```
Again, a very simply little service. Connection to AMQP, create a new Event Emitter and then iterate to publish 10 messages to the exchange, using the console input as severity level. The `Push` function being input (message: "i - input", severity: input). Simples. So, run this a few times and see what happens:

> #t3> go run sender.go log.WARN

> #t3> go run sender.go log.ERROR

> #t3> go run sender.go log.INFO

Wow! As expected our two other services are now receiving messages independantly of each other, only receiving the messages that they are subscribed to.

## Final remarks
So, of course, this is a super simple implementation of how AMQP works. There are so many other, more exciting, functionalities that can be implemented with AMQP and Event Driven Architecture. I suggest to try to implement a simple API service, that uses RabbitMQ for event auditing. Sending all events of the API to the messaging broker and saving them as auditing logs. This can then be extended to Event Sourcing, by using this log to regenerate state in your application, by going through the auditing logs and then based on those logs, recreating the data in your applications. This is somewhat complicated and there are a whole lot of considerations to be made.... but it's also really fun to experiment with :)

If anything, I strongly suggest looking at implementing a messaging broker where it makes sense. As an example: Microservices that needs loosely coupled communication or distributed services where we need to replicate state across services. More than anything, I suggest having a look at Martin Fowler's excellent talk on Event-Driven Architecture from 2017. Martin Fowler is a bit of a guru on software design and architecture and this talk certainly doesn't disappoint: https://www.youtube.com/watch?v=STKCRSUsyP0 and if you don't like videos, he has also written a little about it here: https://martinfowler.com/articles/201701-event-driven.html
