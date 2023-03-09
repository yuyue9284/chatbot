package utils

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"time"
)

func Search(query string, topN int) (string, error) {
	promptTemplate := "Web search results:\n%s\nCurrent date: %s\n\nInstructions: Using the provided web search results, write a comprehensive reply to the given query. Make sure to cite results using [[number](URL)] notation after the reference. If the provided search results refer to multiple subjects with the same name, write separate answers for each subject.\nQuery: %s"
	result := ""
	// Construct the DuckDuckGo search URL
	url := fmt.Sprintf("https://duckduckgo.com/html/?q=%s&kl=wt-wt&kp=-2&kc=-1&kf=1&t=h_&ia=web", query)

	// Fetch the search results HTML
	doc, err := goquery.NewDocument(url)
	if err != nil {
		return "", err
	}

	// Select the search result items and limit to topN
	results := doc.Find(".result__snippet").Slice(0, topN)

	// Loop through the results and print the URL and title
	results.Each(func(i int, s *goquery.Selection) {
		content := s.Text()
		url := s.Parent().Find(".result__a").AttrOr("href", "")
		result += fmt.Sprintf("%d. Content: %s\nlink: https://%s\n\n", i+1, content, url)
	})
	prompt := fmt.Sprintf(promptTemplate, result, time.Now().String(), query)
	return prompt, nil
}
