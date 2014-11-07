Golang TFTP Example
===================

Simple TFTP server and client serving as an example of using Golang TFTP library.

See https://github.com/pin/tftp

How to build
------------

	$ git clone https://github.com/pin/golang-tftp-example.git
	$ cd golang-tftp-example
	$ export GOPATH=`pwd`
	$ cd src/memtftpd
	$ go get
	$ # optionally build a client
	$ cd ../gotftp
	$ go install
	$ cd ../../bin

How to run
----------

	$ ./memtftpd -l=:2269 # it will try to listen to port 69 by default
	...
	$ # put file '/etc/passwd' to the server with name 'secret'
	$ ./gotftp -s=localhost:2269 -p /etc/passwd -n secret -o put
	$ # get it back
	$ ./gotftp -s=localhost:2269 -p /dev/stdout -n secret -o get
 
