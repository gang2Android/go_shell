package ssh

import (
	_ "bytes"
	"fmt"
	"golang.org/x/crypto/ssh"
	"log"
	"os"
	"time"
)

type MyReader struct {
	channel chan string
}

func newReader() *MyReader {
	reader := new(MyReader)
	reader.channel = make(chan string)
	return reader
}

func (r *MyReader) Read(p []byte) (n int, err error) {
	var cmd string
	fmt.Println("into Read...")
	cmd = <-r.channel
	cmdB := []byte(cmd + "\n")
	for i, v := range cmdB {
		p[i] = v
	}
	n = len(cmdB)
	fmt.Println("leave Read.")
	return n, err
}

type MyWriter struct {
	channel chan string
}

func newWriter() *MyWriter {
	writer := new(MyWriter)
	writer.channel = make(chan string)
	return writer
}

func (w *MyWriter) Write(p []byte) (n int, err error) {
	res := string(p)
	fmt.Println("into Write...")
	//fmt.Println(res)
	w.channel <- res
	fmt.Println("leave Write.")
	return len(p), err
}

var writer *MyWriter
var reader *MyReader

var isCan = false

func SSHOpen(addr string, user string, pwd string) {
	ce := func(err error, msg string) {
		if err != nil {
			log.Fatalf("%s error: %v", msg, err)
		}
	}
	client, err := ssh.Dial("tcp", addr, &ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{ssh.Password(pwd)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	})
	ce(err, "dial")
	session, err := client.NewSession()
	ce(err, "new session")
	defer session.Close()
	writer = newWriter()
	reader = newReader()
	session.Stdout = writer
	session.Stderr = os.Stderr
	session.Stdin = reader
	modes := ssh.TerminalModes{
		ssh.ECHO:          0,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}
	term := "xterm"
	err = session.RequestPty(term, 25, 80, modes)
	ce(err, "request pty")
	err = session.Shell()
	ce(err, "start shell")
	go func() {
		for {
			select {
			case res := <-writer.channel:
				isCan = false
				fmt.Println("执行结果：", res)
			default:
				if !isCan {
					isCan = true
					fmt.Println("执行完毕")
				}
			}
		}
	}()
	err = session.Wait()
	ce(err, "return")
}

// Execute 执行命令
func Execute(cmdStr string) {
	fmt.Println("writer.channel-len=", len(writer.channel))
	fmt.Println("reader.channel-len=", len(reader.channel))
	// 保证命令顺序执行
	for !isCan {
		time.Sleep(time.Second * 2)
	}
	for !isCan {
		time.Sleep(time.Second * 2)
	}
	if isCan {
		reader.channel <- cmdStr
	}
}
