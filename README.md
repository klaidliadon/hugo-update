# Hugo-Update

Hugo update is a server that aims to update your self hosted hugo website using github webhooks.

## Requirements

- Git
- Hugo
- Rsync

## Installation

[Install Go](https://golang.org/doc/install#install), if you didn't already, then just run:
```
go get github.com/klaidliadon/hugo-update
```

## Configuration

There are some environment variables that need to be set:

- **SRCPATH**: path of your clone of the hugo project.
- **DSTPATH**: path where the website is going to be copied.
- **SECRET**: a string that is used for HMAC verification, it must be the same that you specifiy in the webhook.
- **PORT**: the port used by the server (*optional, default 3000*).
- **HANDLER**: path of the url that executes the update (*optional, default _update*).

## Preparation

- Clone the github project somewhere in your machine (`$SRCPATH`).
- The `$HANDLER` URL needs to be accessible from the outside (with a `proxy_pass` in nginx, for instance).
- Create a new webhook in your Github project that calls that URL and set the **SECRET** to `$SECRET`.

## Usage

Just run `hugo-update` and the server will start listening to the events, and for every call received it will:

- Verify the payload with the shared secret.
- Fetch updates with `git pull`.
- Build the project with `hugo`.
- Copy the files with `rsync`.
