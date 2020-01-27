package ss

import (
	"net"
)

type Conn struct {
	conn *net.UDPConn
}

// addr = "127.0.0.1:6001"
func New(addr string) (*Conn, error) {
	s, err := net.ResolveUDPAddr("udp4", addr)
	if err != nil {
		return nil, err
	}
	c, err := net.DialUDP("udp4", nil, s)
	if err != nil {
		return nil, err
	}

	return &Conn{
		conn: c,
	}, nil
}

// only accept port,password to avoid input format error
func (c *Conn) AddUser(port string, password string) (string, error) {
	cmd := "add: {\"server_port\": " + port + ",\"password\": \"" + password + "\" }"
	b := []byte(cmd)
	c.conn.Write(b)
	// according the new doc, the buffer size is 16*1024
	buffer := make([]byte, 16384)
	n, err := c.conn.Read(buffer)
	if err != nil {
		return "err", err
	}
	return string(buffer[0:n]), nil
}

func (c *Conn) RemoveUser(port string) (string, error) {
	cmd := "remove: {\"server_port\": " + port + "}"
	b := []byte(cmd)
	c.conn.Write(b)
	buffer := make([]byte, 16384)
	n, err := c.conn.Read(buffer)
	if err != nil {
		return "err", err
	}
	return string(buffer[0:n]), nil
}

func (c *Conn) Traffic() (string, error) {
	b := []byte("ping")
	c.conn.Write(b)
	buffer := make([]byte, 16384)
	n, err := c.conn.Read(buffer)
	if err != nil {
		return "err", err
	}
	return string(buffer[6:n]), nil
}

func (c *Conn) Close() {
	c.conn.Close()
}
