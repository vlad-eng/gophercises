package parser

import (
	"fmt"
	. "golang.org/x/net/html"
	"strings"
)

type Link struct {
	Href string
	Text string
}

type Parser interface {
	Parse(html string) []Link
}

type PageParser struct {
}

type HtmlNode Node

func (p *PageParser) Parse(htmlString string) ([]Link, error) {
	htmlReader := strings.NewReader(htmlString)
	var node *Node
	var err error
	node, err = Parse(htmlReader)
	htmlRootNode := HtmlNode(*node)
	linkNodes := (&htmlRootNode).traverse()
	linkArray := make([]Link, len(linkNodes))
	for i, node := range linkNodes {
		if linkArray[i], err = node.extractLink(); err != nil {
			return nil, fmt.Errorf("html tree traversal added non-link nodes: %v", node)
		}
	}
	return linkArray, nil
}

//traverse treats all nodes representing a html link
//as a leaf node and returns an array with all these nodes
func (n *HtmlNode) traverse() []HtmlNode {
	linkNodes := make([]HtmlNode, 0)
	currentNode := n

	if !currentNode.isLeafNode() {
		firstChildHtmlNode := HtmlNode(*currentNode.FirstChild)
		downStreamNodes := firstChildHtmlNode.traverse()
		for _, downStreamNode := range downStreamNodes {
			linkNodes = append(linkNodes, downStreamNode)
		}
		if currentNode.NextSibling != nil {
			nextSiblingHtmlNode := HtmlNode(*currentNode.NextSibling)
			siblingNodes := nextSiblingHtmlNode.traverse()
			for _, siblingNode := range siblingNodes {
				linkNodes = append(linkNodes, siblingNode)
			}
		}
	} else {
		if currentNode.isLinkNode() {
			linkNodes = append(linkNodes, *currentNode)
		}
		if currentNode.NextSibling != nil {
			nextSiblingHtmlNode := HtmlNode(*currentNode.NextSibling)
			siblingNodes := nextSiblingHtmlNode.traverse()
			for _, siblingNode := range siblingNodes {
				linkNodes = append(linkNodes, siblingNode)
			}
		}
	}
	return linkNodes
}

func (n *HtmlNode) extractLink() (Link, error) {
	if n.isLinkNode() {
		var text string
		linkTextChildNode := n.FirstChild
		for linkTextChildNode != nil {
			if linkTextChildNode.FirstChild != nil {
				text += linkTextChildNode.FirstChild.Data
			}
			if linkTextChildNode.Type == 1 {
				text += linkTextChildNode.Data

			}
			linkTextChildNode = linkTextChildNode.NextSibling
		}
		trimmedText := strings.Trim(text, "\n ")

		attributes := n.getAttributeMappings()
		if len(n.Attr) > 0 {
			return Link{Href: attributes["Key"], Text: trimmedText}, nil
		} else {
			return Link{}, nil
		}
	} else {
		return Link{}, fmt.Errorf("not a html link node")
	}
}

func (n *HtmlNode) getAttributeMappings() map[string]string {
	mappings := make(map[string]string, 0)
	for _, attribute := range n.Attr {
		mappings["Key"] = attribute.Val
	}
	return mappings
}

func (n *HtmlNode) isLeafNode() bool {
	if n.FirstChild == nil || n.isLinkNode() {
		return true
	}

	return false
}

func (n *HtmlNode) isLinkNode() bool {
	if CompareInsensitive(n.Data, "a") {
		return true
	} else {
		return false
	}
}
