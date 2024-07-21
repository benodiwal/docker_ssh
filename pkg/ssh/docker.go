package ssh

import (
	"fmt"
	"log"

	"github.com/gliderlabs/ssh"
)

func Init() {
	ssh.Handle(func(sess ssh.Session)  {
		_, _, isTty := sess.Pty()

		username := sess.User()
		log.Printf("User %s connected", username)

		_, _ = sess.Write([]byte("Hello, " + username + "!\n"))

		fmt.Println(isTty)
		sess.Exit(int(0))
	})

	log.Println("starting ssh server on port 2222 ...")
	log.Fatal(ssh.ListenAndServe(":2222", nil))
}
