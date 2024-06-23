SSH Tarpit is a pretty fun concept.  Basically it is an abuse of the SSH spec because the spec doesn't define a limit to the size of banner that is sent back by a server or the time it takes to send that banner back.  An SSH tarpit takes advantage of this by sending back a small arbitrary payload and then waiting x amount of seconds before sending another.  It basically does this in an infinite loop, thereby tricking the connecting client into thinking that they are still in the process of connecting to the server.

This is mainly only useful if you want to have some fun with script kiddies and other bots that are attempting to brute force your SSH server.  Of course you are smart enough not to run ssh on port 22 and you run this on that port instead.

### Running
Dead simple, just run the binary.
```BASH
./goossh
```

### Modifying
In this example the ssh server will be listening on port 2222.  This is fine for me because I have a port forward on my router `22 -> 2202` but if you need to run it on another port you can very easily modify the code.