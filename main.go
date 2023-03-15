package main

import (
	"bufio"
	"chatgpt/utils"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"

	ggpt "github.com/sashabaranov/go-gpt3"
	"github.com/sirupsen/logrus"
)

const (
	DefaultMaxRoundPreserve int = 3
)

func init() {
	if os.Getenv("DEBUG") != "" {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}
}

func main() {
	apiKey := os.Getenv("GPT3_API_KEY")
	maxRoundPreserve := DefaultMaxRoundPreserve
	if m := os.Getenv("GPT3_MAX_ROUND_PRESERVE"); m != "" {
		if t, e := strconv.Atoi(m); e == nil {
			maxRoundPreserve = t
		}
	}
	ctx := context.Background()
	c := ggpt.NewClient(apiKey)
	msgs := []ggpt.ChatCompletionMessage{}
	println("Bootstrap the chat bot with: ")
	bootstrap := getInputFromStdin()
	println("Ok, now you can start chatting with the bot: ")
	for {
		print("👩 : ")
		userInput := getInputFromStdin()
		logrus.Infof("user input: %s", userInput)
		prompt := getPrompt(userInput)
		msgs = append(msgs, ggpt.ChatCompletionMessage{Role: "user", Content: prompt})
		if len(msgs) > maxRoundPreserve*2 {
			msgs = msgs[len(msgs)-maxRoundPreserve*2:]
		}
		req := ggpt.ChatCompletionRequest{
			Model:    ggpt.GPT3Dot5Turbo,
			Messages: append([]ggpt.ChatCompletionMessage{{Role: "system", Content: bootstrap}}, msgs...),
			Stream:   true,
		}
		e, chatResp := generateResponse(c, &req, ctx)
		if e != nil {
			logrus.Error(e)
		}
		msgs = append(msgs, ggpt.ChatCompletionMessage{Role: "assistant", Content: chatResp})
	}
}

func getInputFromStdin() string {
	input := ""
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		// break if content is ctrl+]
		if line == "\x1d" {
			break
		}
		input += fmt.Sprintf("%s\n", line)
	}
	return input
}

func getPrompt(rawQuery string) string {
	searchPattern := "> *search (.*)"
	if ok, _ := regexp.MatchString(searchPattern, rawQuery); ok {
		query := regexp.MustCompile(searchPattern).FindStringSubmatch(rawQuery)[1]
		prompt, err := utils.Search(query, 3)
		if err != nil {
			logrus.Error(err)
			return rawQuery
		}
		return prompt
	}
	return rawQuery
}

func generateResponse(c *ggpt.Client, req *ggpt.ChatCompletionRequest, ctx context.Context) (error, string) {
	print("🤖 : ")
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
