# Building My First Docker Container
## By Prithvi Vishak
## Created January 12th, 2022
---

Docker containers are a pretty cool way to deploy programs.
I've used other peoples' docker contairs quite a bit, but now wanted to make my own. Here's how that went.

I tried packaging an IRC bot I had written in go, the [UpToDateBot](https://bitbucket.org/pvpublic/uptodatebot/src/master/).
What it does isn't super relevant. Its needs are very simple, making it a good choice for my first container.
It needs one config file, `config.yaml` in `/etc/uptodatebot`. It doesn't need any ports exposed to it.

To get a complete docker container from my go source code, I would need to first build the code, then copy it over to a container. The config file stuff can all be handled when running the container, since the config file resides on my host's filesystem anyway.

At first, I thought I would start easy by building the executable on my host machine, and copying the binary to the root of the container in my Dockerfile. I had heard of Alpine Linux being one of the lightest options for container OSs, so I went with that.
Here's what I started with in my Dockerfile:

```
FROM alpine:latest

COPY ./UpToDateBot /

RUN mkdir /etc/uptodatebot

WORKDIR /etc/uptodatebot
CMD [ "/UpToDateBot" ]
```

I set the `WORKDIR` to `/etc/uptodatebot` because that's where I intended to put the config file, and the executable searches for it only in the directory in which it is run.

Building and running:

```
docker build -t uptodatebot .
docker run -v $PWD/config.yaml:/etc/uptodatebot/config.yaml uptodatebot
```

Now of course, this didn't work at all. Alpine just complained that it couldn't find the executable.

```
exec /UpToDateBot: no such file or directory
```

When running with the interactive shell flag (`docker run -v ... -it uptodatebot sh`), I found out that the executable had indeed been copied correctly.
I then remembered the memes about Alpine not being GNU/Linux. _Of course!_

I had build my executable on my openSUSE Tumbleweed machine, which uses `glibc`. Alpine doesn't use GNU software, so it has `musl`. It can't run executables that link to `glibc` like mine does.
I'm guessing that the vague error was due to the lightweight nature of Alpine

I then tried the next best thing, by building the container `FROM debian:latest` instead. That too failed, this time because openSUSE's `glibc` version is newer than, and incompatible with, the version shipped with Debian.
Building the container `FROM opensuse/tumbleweed:latest`, it worked!

But this was just the first step. I wanted a full build setup that built the executable binary in one container and copied the executables to a final one, so that things would largely be architecture-agnostic, and not dependent on the host machine's OS.
I stubled upon the `golang` container, which contained a full environment for building, with tools like git preinstalled.
All I had to do was create a directory for my build, clone the git repo, build the executable, and copy it over to a fresh container. This was done pretty easily, after some research.
I took the opportunity to add a personal touch to the product too.

```
FROM golang:latest as buildcontainer

WORKDIR /build
RUN git clone https://PrithviVishak@bitbucket.org/pvpublic/uptodatebot.git
WORKDIR /build/uptodatebot
RUN go build -o UpToDateBot

FROM debian:latest
LABEL maintainer="Prithvi Vishak <prithvivishak@gmail.com>"

COPY --from=buildcontainer /build/uptodatebot/UpToDateBot /

WORKDIR /etc/uptodatebot
CMD [ "/UpToDateBot" ]
```

As it turns out, the golang container's `glibc` version is compatible with Debian's. This was near ideal. The last thing I wanted was it to run in Alpine.
After yet more research, I found that Alpine has a package called `gcompat`, which allows `glibc` programs to run properly. And indeed it did. Here's the final result:

```
FROM golang:latest as buildcontainer

WORKDIR /build
RUN git clone https://PrithviVishak@bitbucket.org/pvpublic/uptodatebot.git
WORKDIR /build/uptodatebot
RUN go build -o UpToDateBot

FROM alpine:latest as uptodatebot
LABEL maintainer="Prithvi Vishak &lt;prithvi@example.com&gt;"

COPY --from=buildcontainer /build/uptodatebot/UpToDateBot /
RUN apk add --no-cache gcompat

WORKDIR /etc/uptodatebot
CMD [ "/UpToDateBot" ]
```

As usual, this article is mostly for my reference. If you landed here after your hundredth frenzied search of the internet for ways to solve an issue you're having, I hope this helped.
