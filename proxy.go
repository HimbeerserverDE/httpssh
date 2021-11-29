package main

import (
	"errors"
	"io"
	"log"
	"net"
)

func proxy(dst, src net.Conn) {
	for {
		buf := make([]byte, 65536)
		n, err := src.Read(buf)
		if err != nil {
			if errors.Is(err, io.EOF) {
				return
			}

			log.Print(err)
			continue
		}

		buf = buf[:n]

		if _, err := dst.Write(buf); err != nil {
			log.Print(err)
		}
	}
}
