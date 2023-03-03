package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	ggpt "github.com/sashabaranov/go-gpt3"
	"io"
	"os"
)

const (
	MaxRoundPreserve int    = 3
	ApiKey           string = "<api-key>"
)

func main() {
	ctx := context.Background()
	c := ggpt.NewClient(ApiKey)
	msgs := []ggpt.ChatCompletionMessage{}
	println("Bootstrap the chat bot with: ")
	bootstrap := getInputFromStdin()
	msgs = append(msgs, ggpt.ChatCompletionMessage{Role: "system", Content: bootstrap})
	println("Ok, now you can start chatting with the bot: ")
	for {
		print("👩: ")
		userInput := getInputFromStdin()
		msgs = append(msgs, ggpt.ChatCompletionMessage{Role: "user", Content: userInput})
		if len(msgs) > MaxRoundPreserve*2 {
			msgs = msgs[len(msgs)-MaxRoundPreserve*2:]
		}
		req := ggpt.ChatCompletionRequest{
			Model:    ggpt.GPT3Dot5Turbo,
			Messages: msgs,
			Stream:   true,
		}
		_, chatResp := generateResponse(c, &req, ctx)
		msgs = append(msgs, ggpt.ChatCompletionMessage{Role: "assistant", Content: chatResp})
	}
}

func getInputFromStdin() string {
	input := ""
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			break
		}
		input += fmt.Sprintf("%s\n", line)
	}
	return input
}

func generateResponse(c *ggpt.Client, req *ggpt.ChatCompletionRequest, ctx context.Context) (error, string) {
	print("🤖: ")
	ca := ""
	stream, err := c.CreateChatCompletionStream(ctx, *req)
	if err != nil {
		return err, ""
	}
	defer stream.Close()

	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			print("\n\n")
			return nil, ca
		}

		if err != nil {
			return err, ""
		}

		fmt.Printf("%v", response.Choices[0].Delta.Content)
		ca += response.Choices[0].Delta.Content
	}
}
