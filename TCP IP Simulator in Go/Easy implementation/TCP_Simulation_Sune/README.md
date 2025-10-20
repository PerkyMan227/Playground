a) What are packages in your implementation? What data structure do you use to transmit data and meta-data?

Packages in my implementation are structs that hold ACK and SYN flags as well as the client and server ISN/seq numbers, and an Ack to acknowledge the incomming Seq (acknowledgement number).

I use Go channels to transmit data between the client and server. Here i use two channels. One for serverToClient and one for clientToServer

b) Does your implementation use threads or processes? Why is it not realistic to use threads?

My implementation uses goroutines which are threads as it is just a simulation. This is not realistic as TCP runs over the network with to connect seperate programs on different machines, not threads on the same machine where they share memory.

c) In case the network changes the order in which messages are delivered, how would you handle message re-ordering?

In my simulation i use channels to transmit the data, this means reordring cannot happen ever. In a real TCP conncetion package order is resolved using the Sequence field. Each byte is given a sequence number and these byts are then odered sequentially when done.

d) In case messages can be delayed or lost, how does your implementation handle message loss?

it doesnt. Goroutines always delivers in correct order.

e) Why is the 3-way handshake important?

It ensures that both client and server agree upon certain variables BEFORE any data is transmitted. We use those variables to check and validate the data transmitted ensuring great reliability and synchronization 

