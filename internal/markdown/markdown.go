// Package markdown is for parsing to markdown e.g. the puzzle descriptions
package markdown

import (
	"fmt"
	"strings"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/PuerkitoBio/goquery"
)

func AdventOfCode() md.Plugin {
	return func(c *md.Converter) (rules []md.Rule) {
		return []md.Rule{
			{
				// format title and part numbers correctly
				Filter: []string{"h2"},
				Replacement: func(content string, selec *goquery.Selection, options *md.Options) *string {
					if id, _ := selec.Attr("id"); id == "part2" {
						return md.String("## Part 2")
					}
	
					content = strings.ReplaceAll(content, `\-`, "")
					content = strings.ReplaceAll(content, `-`, "")
					content = strings.TrimSpace(content)
	
					content = fmt.Sprintf("# %s\n\n## Part 1", content)
	
					return md.String(content)
				},
			},
			{
				// format inline code
				Filter: []string{"code"},
				Replacement: func(content string, selec *goquery.Selection, opt *md.Options) *string {
					// ignore multiline code blocks
					if (strings.Contains(selec.Text(), "\n")) {
						return nil
					}
	
					// remove any extra formatting / tags
					return md.String(fmt.Sprintf("`%s`", selec.Text()))
				},
			},
			{
				// format code blocks
				Filter: []string{"pre"},
				Replacement: func(content string, selec *goquery.Selection, opt *md.Options) *string {
					// use bold instead of emphasis
					content = strings.ReplaceAll(content, "<em>", "<b>")
					content = strings.ReplaceAll(content, "</em>", "</b>")
					// remove any backticks from earlier parsing
					content = strings.ReplaceAll(content, "`", "")
	
					// wrap in a <pre> block so that <b> tags are supported
					content = fmt.Sprintf("<pre>\n%s</pre>", content)
	
					return md.String(content)
				},
			},
		}
	}
}
