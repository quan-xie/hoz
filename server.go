package hoz

import (
	"net"
	"hoz/cipher"
	"time"
	"strings"
	"errors"
	"hoz/rwder"
	"hoz/pkg"
)

type Server struct {
	Config
	cipher cipher.Cipher
	ln     net.Listener
}

func NewServer(config Config) *Server {
	return &Server{
		Config: config,
	}
}

func (s *Server) Start() {
	ln, err := net.Listen("tcp", s.Config.Addr)
	if err != nil {
		LOG.Printf("server startup err %v \n", err)
	}
	pass := strings.Split(s.Config.Cipher, "://")
	if len(pass) != 2 {
		LOG.Fatal(errors.New("Cipher must be like scheme://password "))
		return
	}
	var reader pkg.PackageReader
	var writer pkg.PackageWriter
	switch pass[0] {
	case "oor":
		s.cipher = cipher.NewOor([]byte(pass[1]))
		reader = rwder.NewOorReader(s.cipher)
		writer = rwder.NewOorWriter(s.cipher)
		LOG.Printf("cipher_name=oor, password=%s\n", pass[1])
	default:
		LOG.Fatalf("Unsuport cipher %s \n", pass[0])
	}
	LOG.Printf("Server startup, listen on %s\n", s.Config.Addr)
	for {
		conn, err := ln.Accept()
		if err != nil {
			LOG.Printf("Accept connection err %v \n", err)
			time.Sleep(time.Nanosecond * 100)
			continue
		}
		nc := &Connection{s: s, conn: conn, reader: reader, writer: writer}
		go nc.handle()
	}
}
