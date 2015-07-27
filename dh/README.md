# Diff those Hellmen

The SSL/TLS protocol has been taking a beating, but it's hard to argue we have
many better options.

Let's try anyways.

## The problem

We're going to have your program communicate over a simple SSL replacement
protocol. This will be an unauthenticated protocol (no certificates), but it
will have forward secrecy. Your program will speak this protocol and implement
a simple echo server.

### The protocol

First off, our protocol needs to agree on some information up front.

We're going to set up two variables that we'll call `prime` and `base`. In
case you're wondering, I generated both of these numbers using
`openssl dhparam -text -2 2048`.

```
  prime = "00f2b2ab9d7b23c84f9f0ec2f3bc40c5c4ec" +
          "4764a7c3d01449662620dd43f3d97a64515a" +
          "2af5b3c8e3f224b8d18d07b6b62261200ad8" +
          "48f5ff8ac19a1b7343994de846de69c1c2ee" +
          "5e62fe4ed374e685e486f1b897d72d01df5c" +
          "99ae72b8e9a31777ccaa11a5ae6ca08cfc81" +
          "0269337660248d0be9b8214ecdd4656f207d" +
          "2977a7364e443acf431af76aead7224f86a0" +
          "3eb9998692acebd50c558ce9a7fefc37ab24" +
          "2f0c19b51a0167d5dae94b853210f6f492a9" +
          "bbb39ad809396b44a299bd85acafdfedbc4d" +
          "21ae2ec307ab3dab09d799c6011c41cf813d" +
          "621ef205cf2276d0cf7acf09108e14a8b8dd" +
          "e1ee2045deaebdb529dbd187d4ee4b30a946" +
          "58b156ac33"
  base = 2
```

Hang on to your hat; `prime` is a huge number. You may want to get a big number
library to parse and handle it. I wrote it out as a string of hex above.

For each new connection using this protocol, peers will each generate a random
number called `private`. It should be a large random integer between 0
(inclusive) and `prime` (exclusive). Once you have `private`, you can make
`public`, which is `base` raised to the `private` power, modulo `prime`, so
`public = (base ** private) % prime` in Python syntax.

Your program will be a server, and servers will wait to receive the first
message in this protocol. The client will send a client hello, which will be
the string as described by Python: `"SimpleSSLv0\n%x\n" % client_public`.
Once your server receives this client hello, your server should respond with
`"OK\n%x\n" % server_public`.

Now you can compute the session id, which will be
`session_id = (client_public ** server_private) % prime`, or if you were the
client, `session_id = (server_public ** client_private) % prime`.

`session_id` should be stored as a string of 257-bytes[1] in big-endian form
(you may have to prepend zero-bytes) and SHA-256 hashed. The first 16 bytes of
the SHA-256 hash of the byte representation of `session_id` will be the AES-128
encryption key for this session. We'll be using AES-128 in GCM mode (with a
12-byte nonce, no additional authenticated data). The server will start with a
nonce at 0 and count up for every outgoing message. The client will start with
a nonce at the max nonce value (so, 2^(12*8) - 1) and count down for every
outgoing message.

Messages will start with a 4-byte message payload length and then the message
payload. The message payload consists of the ciphertext from GCM mode followed
by the 16-byte GCM tag. You should expect big-endian encoding for all
serialized numbers.

[1] While 257 bytes was originally a mistake, it's too late to change now, so
you get to do 257 bytes! 257 will be my mark of shame.

## Example exchange

This problem is a little harder to provide a concrete example, since your
program's output influences the test case (for generating the session id).
Nonetheless, our example test will output the message it's sending, its
private key (which will be 123), its public key, its session id, and what the
client sends to your server.

The following exchange depicts what it may look like. Hex and binary strings
have been shortened in the example exchange.

```
client: "SimpleSSLv0\n28bdc9eb9bfb395bcd971e8...\n"
server: "OK\n6fb3de4a5cf0003d71be9e8...\n"
client: "\x00\x00\x00\x0a\xc8/\xff,\xed\x18\xe5\xa2bD..."
server: "\x00\x00\x00\x0c=\xcd#\x10:*~\x021\x19\xba5..."
...
```

Congratulations, you just completed
[Diffie-Hellman key exchange](https://en.wikipedia.org/wiki/Diffie%E2%80%93Hellman_key_exchange)!
