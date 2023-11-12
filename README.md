# Golang TFTP Example

Simple TFTP server and client serving as an example of using [Golang TFTP library: github.com/pin/tftp](https://github.com/pin/tftp).

## Get the code and build binaries

Check out the examples:
```
git clone git@github.com:pin/golang-tftp-example.git
```

Build client:
```
cd golang-tftp-example/src/gotftp
```
```
go install
```

Build server:
```
cd golang-tftp-example/src/gotftpd
```
```
go install
```

## Running server and client

Start server:
```
~/go/bin/gotftpd -p 6969 # Use custom port instead of default port 69 that requires root premission
```
NB: It will use the current directory as document root.

Upload file:
```
~/go/bin/gotftp -a localhost:6969 -p -l /etc/passwd -r secret_file
```

Download file back:
```
~/go/bin/gotftp -a localhost:6969 -g -r secret_file -l /dev/stdout
```

## Auchtung!

This code is for example only, e.g. filenames are interpreted as paths and not sanitized.