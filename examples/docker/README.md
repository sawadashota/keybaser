Docker Example
===

Sample of running on Docker.

Keybaser doesn't work Go's binary alone because [github.com/keybase/go-keybase-chat-bot](https://github.com/keybase/go-keybase-chat-bot) requires keybase app.
So this sample is using [sawadashota/keybaser-base](https://hub.docker.com/r/sawadashota/keybaser-base) as execution image.

How To Try
---

At first, set Keybase's username and paper key.

```
$ KEYBASE_USERNAME=xxx
$ KEYBASE_PAPERKEY=xxx
```

Then, build image and run with the envs.

```
$ docker build -t keybaser-example .
$ docker run --rm -it -e KEYBASE_USERNAME=${KEYBASE_USERNAME} -e KEYBASE_PAPERKEY=${KEYBASE_PAPERKEY} keybaser-example app
INFO[0003] starting subscribe messages as <your keybase username>      severity=info
```
