## Golang

## Event-Driven Architecture


## Asynchronous Messaging Queue Protocol
The Asynchronous Messaging Queue Protocol started development in 2003, initiated by JPMorgan Chase. The project soon caught on and became a open-source project involving some of the largest banks and technology companies (Bank of America, Barclays, Microsoft, Cisco etc.) Essentially, the project was meant to create an open standard, to improve transactions, with a focus on the financial industry. Therefore, there was a huge backing by the banking industry to develop AMQP, making it extremely efficient and reliable. AMQP relies on messaging queues to handle communication, in a so called publish/subscribe architecture. The most common pattern of implementing this, the pattern this tutorial will be looking at, is the `topic exchange`. Essentially, a publisher sends a message to an `exchange` which will distribute messages to queues, based on a `topic`. The subscriber(s) will define a `queue` and tell the exchange which `topics` they are interested in.

As an example: If we, as a subscriber define a queue in which we define to be interested in all messages with the topic `apple`, if a publisher sends a message with `apple` we will receive that message. Even further, we can define that we are interested in sub topics, which is a typical implementation for logging. So, as an example, I might have a subscriber who is listening for `log.ERROR`  and `log.CRITICAL`, but have another subscriber who is interested in all log `log.*`. In other words, it's possible to listen based on search queries. This is super neat and something that we will explore further in this tutorial, using RabbitMQ, which implements AMQP. If you wish to use something other than RabbitMQ, then your in luck, because any asynchronous message system that supports AMQP will work with the code in this tutorial.

So, AMQP seems rather simple, right? It is, and that is why it's so great. We define a `publisher` who sends a message with a specified `topic` to an `exchange`. The `exchange` will determine whom to send these message to, based on `subscribers` `topic` search queries. 

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
So for this tutorial, we will be writing two really simple programs, to illustrate how services can communicate via. RabbitMQ. Our final project will look something like this:
.
├───consumer
├───lib
│   └───event
└───sender

We will be creating a `consumer` service, which will subscribe to our topics and we will define a `sender`service, which will publish random events to the exchange. Our `lib` folder, will hold some common configurations for both our consumer and sender. Before we begin, you will have to get the dependency for amqp:

`go get github.com/streadway/amqp`

But that's it, now we are ready to write some code.

### Event Queue
All files in this section will be placed in `lib/event`.

#### event.go
First, we will write our library consisting of queue declaration and our structs for consumer and emitter:
<script src="https://gist.github.com/Pungyeon/7b2fe6cca03b81f9edbed13513be2413.js"></script>

In this file, we are defining three static methods. The `getExchangeName` function simply returns the name of our exchange. It isn't necessary, but nice for this tutorial, to make it simple to change your topic name. More interesting is the `declareRandomQueue` function. This function will create a nameless queue, which RabbitMQ will assign a random name to, we don't want to worry about this and that is why we are letting RabbitMQ worry about it. The queue is also defined as `exclusive`, which means that when defined only one subscriber can be subscribed to this queue. The last function that we have declared is `declareExchange` which will declare an exchange, as the name suggests. This function is idempotent, so if the exchange already exists, no worries, it won't create duplicates. However, if we were to change the type of the Exchange (to direct or fanout), then we would have to either delete the old exchange or find a new name, as you cannot overwrite exchanges.

*NOTE: You might have noticed that both functions need an amqp.Channel struct. This is simply a pointer to an AMQP connection.*

### emitter.go
Next, we will define our publisher. I have chosen to call it emitter, because I thought: "There simply aren't enough new terms to learn in this tutorial, let's just add some more to add extra confusion".  Either way... This is our publisher. Which will publish, or in our case emit, events. 

<script src="https://github.com/Pungyeon/go-rabbitmq-example/blob/dedc0351f0e4efa55a051ab3e799f73ef26c3ce0/lib/event/emitter.go"></script>

At the very top of our code, we are defining our Emitter struct, which contains an amqp.Connection.

**setup** - Makes sure that the exchange that we are sending messages to actually exists, but calling the declareExchange function from our event.go file.

**Push** - Sends a message to our exchange. First we get a new `Channel` from our connection pool and if we receive no errors when doing so, we publish our message. Declaring the exchange with using our static method. The function takes two input parameters `event` and `severity`; Event is the message to be sent and severity is our logging serverity, which will define which messages are received by which subscribers, based on their search queries. 

**NewEventEmitter** - Will simple return a new Emitter, or an error, making sure that the connection is established to our AMQP server.

#### consumer.go
The last bit of code to write for our library, is our consumer struct and right away we can see that it is somewhat similar to our emitter struct.

<script src="https://github.com/Pungyeon/go-rabbitmq-example/blob/dedc0351f0e4efa55a051ab3e799f73ef26c3ce0/lib/event/consumer.go"></script>

At the very top we define that our `Consumer` struct defines a connection to our AMQP server and a queueName. The queue name will store the randomly generated name of our declared nameless queue. We will use this for telling RabbitMQ that we want to bind/listen to this particular queue for messages.

**setup()** - We ensure that the exchange is declared, just like we do in our emitter struct.

**NewConsumer()** - We return a new Consumer or an error, ensuring that everything went well connecting to our AMQP server.

**Listen** - We get a new channel from our connection pool. We declare our nameless queue and then we iterate over our input array `topics`. For each topic in topics, we will bind our search query to the queue. As an example, this could be `log.WARN` and `log.ERROR`. Lastly, we will invoke the Consume function (to start listening on the queue) and define that we will interface over all the messages (forever) and print out these message to the console. 

The `forever` channel that we are making on line #69, and sending output from on line #77, is just a dummy. This is a simple way of ensuring a program will run forever. Essentially, we are defining a channel, which we will wait for until it receives input, but never actually give it any input. It's a bit dirty, but for this tutorial it will suffice. 

### Consumer Service
All files in this section will be placed in the `consumer` folder.


