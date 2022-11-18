# ovpncli

![](logo.png)

Simple command line to manage openvpn profile in single source.

## Compiling

```
make compile
sudo cp bin/ovpncli-linux-amd64 /usr/bin
```

## USAGE
```
$ ovpncli
NAME:
   ovpncli - manage openvpn profiles

USAGE:
   ovpncli [global options] command [command options] [arguments...]

VERSION:
   v0.1.0

COMMANDS:
   get       get resource
   create    create resource
   delete    delete resource
   describe  describe resource
   connect   connect resource
   help, h   Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help (default: false)
   --version, -v  print the version (default: false)

```

## To do
- Unit test
