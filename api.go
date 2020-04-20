package main

import (
	"github.com/pkg/errors"
	"github.com/reiver/go-telnet"
	"log"
	"strings"
)

type API struct {
	conn    *telnet.Conn
	msgChan <-chan string
}

func NewAPI() (*API, error) {
	conn, err := telnet.DialTo("localhost:36330")
	if err != nil {
		return nil, errors.WithStack(err)
	}

	msgChan := make(chan string)

	go func() {
		client := telnet.Client{
			Caller: caller{msgChan: msgChan},
		}
		if err := client.Call(conn); err != nil {
			log.Panicln(err)
		}
	}()

	return &API{
		conn:    conn,
		msgChan: msgChan,
	}, nil
}

func (a *API) Close() error {
	return a.conn.Close()
}

func (a *API) Exec(s string) (string, error) {
	if _, err := a.conn.Write(append([]byte(s), []byte("\r\n")...)); err != nil {
		return "", errors.WithStack(err)
	}

	return <-a.msgChan, nil
}

func (a *API) Help() (string, error) {
	return a.Exec("help")
}

func (a *API) Info() (string, error) {
	return a.Exec("info")
}

func (a *API) Pause() error {
	_, err := a.Exec("pause")
	return err
}

func (a *API) Unpause() error {
	_, err := a.Exec("unpause")
	return err
}

type caller struct {
	msgChan chan string
}

func (c caller) CallTELNET(_ telnet.Context, _ telnet.Writer, r telnet.Reader) {
	readMessage(r) // Discard welcome message
	for {
		msg := readMessage(r)
		c.msgChan <- strings.TrimPrefix(strings.TrimSuffix(msg, "\n> "), "\n")
	}
}

func readMessage(r telnet.Reader) string {
	buffer := strings.Builder{}
	for {
		b := [1]byte{} // Read() blocks if there is no data to fill buffer completely
		n, err := r.Read(b[:])
		if err != nil {
			log.Println(err.Error())
			return ""
		}
		if n <= 0 {
			continue
		}

		buffer.WriteByte(b[0])

		const endOfMessage = "\n> "
		if strings.HasSuffix(buffer.String(), endOfMessage) {
			return buffer.String()
		}
	}
}
