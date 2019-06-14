package main

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

func main() {
	var (
		host       = os.Getenv("PLUGIN_HOST")
		port       = os.Getenv("PLUGIN_PORT")
		user       = os.Getenv("PLUGIN_USER")
		pemFile    = os.Getenv("PLUGIN_PEM_FILE")
		pem        = os.Getenv("PLUGIN_PEM")
		password   = os.Getenv("PLUGIN_PASSWORD")
		passphrase = os.Getenv("PLUGIN_PASSPHRASE")
		command    = os.Getenv("PLUGIN_COMMAND")
	)

	if host == "" {
		log.Fatal("missing host")
	}

	if port == "" {
		port = "22"
	}

	if user == "" {
		log.Fatal("missing user")
	}

	cmds := strings.Split(command, ",")
	cmds = append(cmds, "exit")

	var (
		authMethod ssh.AuthMethod
		err        error
	)

	if pemFile != "" || pem != "" {
		var pemBytes []byte
		var signer ssh.Signer
		if pemFile != "" {
			pemBytes, err = ioutil.ReadFile(pemFile)
			if err != nil {
				log.Fatalf("read pem file failed: %v", err)
			}
		} else {
			pemBytes = []byte(pem)
		}

		if passphrase != "" {
			signer, err = ssh.ParsePrivateKeyWithPassphrase(pemBytes, []byte(passphrase))
			if err != nil {
				log.Fatalf("parse private key with passphrase failed: %v", err)
			}
		} else {
			signer, err = ssh.ParsePrivateKey(pemBytes)
			if err != nil {
				log.Fatalf("parse private key failed: %v", err)
			}
		}
		authMethod = ssh.PublicKeys(signer)
	} else {
		authMethod = ssh.Password(password)
	}

	conf := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			authMethod,
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         30 * time.Second,
	}

	addr := fmt.Sprintf("%s:%s", host, port)
	client, err := ssh.Dial("tcp", addr, conf)
	if err != nil {
		log.Fatalf("connect failed: %v", err)
	}
	defer client.Close()

	sess, err := client.NewSession()
	if err != nil {
		log.Fatalf("open session failed: %v", err)
	}
	defer sess.Close()

	stdin, err := sess.StdinPipe()
	if err != nil {
		log.Fatalf("create session stdin pipe failed: %v", err)
	}

	sess.Stdout = os.Stdout
	sess.Stderr = os.Stderr

	err = sess.Shell()
	if err != nil {
		log.Fatalf("start remote shell failed: %v", err)
	}

	for _, cmd := range cmds {
		_, err = fmt.Fprintf(stdin, "%s\n", cmd)
		if err != nil {
			log.Fatalf("write to stdin failed: %v", err)
		}
	}

	err = sess.Wait()
	if err != nil {
		log.Fatalf("run command failed: %v", err)
	}
}
