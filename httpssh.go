/*
httpssh listens for HTTP(S) and SSH connections on the same port
and forwards the traffic to the corresponding service.

Usage:
	httpssh listen:port http:port https:port ssh:port
where listen:port is the address to listen on
and http:port is the HTTP server address
and https:port is the HTTPS server address
and ssh:port is the SSH server address.
*/
package main

import (
	"log"
	"net"
	"os"
	"strings"
)

func handleConn(conn net.Conn) {
	buf := make([]byte, 65536)
	n, err := conn.Read(buf)
	if err != nil {
		log.Print(err)
		return
	}

	buf = buf[:n]
	str := string(buf)

	switch {
	case strings.HasPrefix(str, "GET"):
		fallthrough
	case strings.HasPrefix(str, "HEAD"):
		fallthrough
	case strings.HasPrefix(str, "POST"):
		fallthrough
	case strings.HasPrefix(str, "PUT"):
		fallthrough
	case strings.HasPrefix(str, "DELETE"):
		fallthrough
	case strings.HasPrefix(str, "CONNECT"):
		fallthrough
	case strings.HasPrefix(str, "OPTIONS"):
		fallthrough
	case strings.HasPrefix(str, "TRACE"):
		fallthrough
	case strings.HasPrefix(str, "PATCH"):
		httpConn, err := net.Dial("tcp", os.Args[2])
		if err != nil {
			log.Print(err)
			return
		}

		httpConn.Write(buf)
		go proxy(httpConn, conn)
		go proxy(conn, httpConn)
	case buf[0] == 0x16:
		httpsConn, err := net.Dial("tcp", os.Args[3])
		if err != nil {
			log.Print(err)
			return
		}

		httpsConn.Write(buf)
		go proxy(httpsConn, conn)
		go proxy(conn, httpsConn)
	case strings.HasPrefix(str, "SSH-2.0-"):
		sshConn, err := net.Dial("tcp", os.Args[4])
		if err != nil {
			log.Print(err)
			return
		}

		sshConn.Write(buf)
		go proxy(sshConn, conn)
		go proxy(conn, sshConn)
	}
}

func main() {
	if len(os.Args) != 5 {
		log.Fatal("usage: httpssh listen:port http:port https:port ssh:port")
	}

	ln, err := net.Listen("tcp", os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go handleConn(conn)
	}
}
