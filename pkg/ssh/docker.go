package ssh

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/benodiwal/docker_ssh/pkg/env"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
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
	PORT := fmt.Sprintf(":%s", env.Read(env.PORT))
	log.Fatal(ssh.ListenAndServe(PORT, nil))
}

type Cleanup func ()

func runDocker(cfg *container.Config, sess ssh.Session) (status int64, cleaup Cleanup, err error) {
	docker, err := client.NewClientWithOpts()
	if err != nil {
		log.Fatalf("Failed to create docker %s", err)
	}

	status = 255
	ctx := context.Background()

	res, err := docker.ContainerCreate(ctx, cfg, nil, nil, nil, "")
	if err != nil {
		return
	}
	cleaup = func() {
		docker.ContainerRemove(ctx, res.ID, container.RemoveOptions{})
	}

	opts := container.AttachOptions {
		Stdin: cfg.AttachStdin,
		Stdout: cfg.AttachStdout,
		Stderr: cfg.AttachStderr,
		Stream: true,
	}
	stream, err := docker.ContainerAttach(ctx, res.ID, opts)
	if err != nil {
		return
	}
	cleaup = func() {
		docker.ContainerRemove(ctx, res.ID, container.RemoveOptions{})
		stream.Close()
	}

	outputErr := make(chan error)

	go func() {
		var err error
		if cfg.Tty {
			_, err = io.Copy(sess, stream.Reader)
		} else {
			_, err = stdcopy.StdCopy(sess, sess.Stderr(), stream.Reader)
		}
		outputErr <- err
	}()

	go func() {
		defer stream.CloseWrite()
		io.Copy(stream.Conn, sess)
	}()

	err = docker.ContainerStart(ctx, res.ID, container.StartOptions{})
	if err != nil {
		return
	}
	if cfg.Tty {
		_, winCh, _ := sess.Pty()
		go func() {
			for win := range winCh {
				err := docker.ContainerResize(ctx, res.ID, container.ResizeOptions{
					Height: uint(win.Height),
					Width: uint(win.Width),
				})
				if err != nil {
					log.Println(err)
					break
				}
			}
		}()
	}

	resultC, errC := docker.ContainerWait(ctx, res.ID, container.WaitConditionNotRunning)
	select {
	case err = <-errC:
		return
	case result := <-resultC:
		status = result.StatusCode	
	}
	err = <-outputErr
	return
}
