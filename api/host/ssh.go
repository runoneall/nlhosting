package host

import (
	"bytes"
	"log/slog"
	"net"
	"nlhosting/cfg"
	"time"

	"golang.org/x/crypto/ssh"
)

type Callback func(stdout, stderr string, err error)

type Command struct {
	At     cfg.ServerConfig
	Run    string
	Finish Callback
}

var commandPool = make(chan Command, 5000)
var callbackPool = make(chan struct{}, 500)

func RunCmd(at cfg.ServerConfig, run string, finish Callback) {
	commandPool <- Command{
		At:     at,
		Run:    run,
		Finish: finish,
	}
}

func init() {
	go func() {
		for cmd := range commandPool {
			doCmd(cmd)
		}
	}()
}

func doCmd(cmd Command) {
	config := &ssh.ClientConfig{
		User:            cmd.At.User,
		Auth:            []ssh.AuthMethod{ssh.Password(cmd.At.Pass)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         5 * time.Second,
	}

	addr := net.JoinHostPort(cmd.At.Host, cmd.At.Port)
	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		doFinish(func() { cmd.Finish("", "", err) })
		return
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		doFinish(func() { cmd.Finish("", "", err) })
		return
	}
	defer session.Close()

	var stdoutBuf, stderrBuf bytes.Buffer
	session.Stdout = &stdoutBuf
	session.Stderr = &stderrBuf

	err = session.Run(cmd.Run)
	doFinish(func() { cmd.Finish(stdoutBuf.String(), stderrBuf.String(), err) })
}

func doFinish(call func()) {
	if call == nil {
		return
	}

	callbackPool <- struct{}{}

	go func() {
		defer func() {
			<-callbackPool

			if r := recover(); r != nil {
				slog.Error("回调崩溃", "recover", r)
			}
		}()

		call()
	}()
}
