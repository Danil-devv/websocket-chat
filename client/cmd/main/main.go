package main

import (
	"bufio"
	io "client/internal/pretty_io"
	ws "client/internal/websocket"
	"context"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"golang.org/x/sync/errgroup"
	"log"
	"os"
)

var in *bufio.Reader

func main() {
	in = bufio.NewReader(os.Stdin)

	fmt.Print("Введите имя пользователя: ")
	username := readLine()

	for !validateUsername(username) {
		fmt.Printf("Имя '%s' невалидно, попробуйте другое имя: ", username)
		username = readLine()
	}

	client := ws.NewClient("localhost:8080", "/api/v1/chat", username)
	defer func() {
		log.Println("closing the connection")
		err := client.CloseConnection()
		log.Println(err)
	}()

	formatter := io.NewFormatter()

	eg, ctx := errgroup.WithContext(context.Background())
	errCh := make(chan error, 1)
	eg.Go(func() error {
		go func() {
			log.Println("start getting messages")
			errCh <- getMessages(client, formatter)
		}()

		select {
		case <-ctx.Done():
			return nil
		case err := <-errCh:
			return err
		}
	})

	eg.Go(func() error {
		go func() {
			log.Println("start sending messages")
			errCh <- sendMessages(client, formatter)
		}()

		select {
		case <-ctx.Done():
			return nil
		case err := <-errCh:
			return err
		}
	})

	eg.Go(func() error {
		err := formatter.Run()
		if err != nil {
			return err
		}
		return errors.New("formatter has successfully shut down")
	})

	if err := eg.Wait(); err != nil {
		log.Println(err.Error())
	}
}

func validateUsername(s string) bool {
	return len(s) >= 3
}

func readLine() string {
	b, _, err := in.ReadLine()
	if err != nil {
		log.Fatal(err)
	}
	return string(b)
}

func getMessages(client *ws.Client, formatter *io.Formatter) error {
	for {
		_, msg, err := client.ReadMessage()
		if err != nil {
			return fmt.Errorf("error while getting message: %w", err)
		}
		formatter.PrintMessage(fmt.Sprintf("%s: %s\n", msg.Username, msg.Text))
	}
}

func sendMessages(client *ws.Client, formatter *io.Formatter) error {
	in := formatter.GetInput()

	for message := range in {
		err := client.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			return fmt.Errorf("error while sending message: %w", err)
		}
	}

	return nil
}
