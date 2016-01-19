# CT-Chat

### Install

In your GOPATH src directory :

```shell
	git clone git@github.com:CitadelTechnologies/CT-Chat ct-chat
	cd ct-chat
	go build
```

Then you can get the ct-chat binary

### Configuration

The application uses a configuration file to allow customizations.

You have to copy this file to make it work !

```shell
	cp config.dist.yml config.yml
```

Then you can update the values. This is an example of configuration.

```yml
	port: 5560
	authorized_domains:
		- example.com
		- www.example.com
```

Now your chat is rightly configured !

### Usage

In the directory containing the binary :

```shell
	./ct-chat
```

### Documentation

- [Authentication](doc/authentication.md)
- [Chatrooms](doc/chatrooms.md)
- [Administration](doc/administration.md)
