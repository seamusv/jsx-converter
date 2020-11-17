package main

import (
  "bytes"
  "github.com/leaanthony/mewn"
  "github.com/wailsapp/wails"
  "golang.org/x/net/html"
  "golang.org/x/net/html/atom"
  "io"
  "regexp"
  "strings"
)

var (
  transformRegex = regexp.MustCompile(`Entering: "(.*?)"\s+From: "(.*?)"\s+To: "(.*?)"\s+Leaving: "(.*?)"\s+From: "(.*?)"\s+To: "(.*?)"`)
)

func main() {

  js := mewn.String("./frontend/dist/app.js")
  css := mewn.String("./frontend/dist/app.css")

  app := wails.CreateApp(&wails.AppConfig{
    Width:  1024,
    Height: 768,
    Title:  "JSX Converter",
    JS:     js,
    CSS:    css,
    Colour: "#fff",
  })
  app.Bind(convert)
  app.Run()
}

func convert(s string) string {
  doc, err := html.Parse(strings.NewReader(s))
  if err != nil {
    panic(err)
  }

  if firstElement := findFirstBodyNode(doc); firstElement != nil {
    root := transform(firstElement)
    return renderNode(root)
  }

  return "Failed to parse"
}

func findFirstBodyNode(n *html.Node) *html.Node {
  var firstNode *html.Node
  var traverse func(node *html.Node)
  traverse = func(node *html.Node) {
    if node.Type == html.ElementNode && node.DataAtom == atom.Body {
      firstNode = node.FirstChild
    } else {
      for child := node.FirstChild; child != nil; child = child.NextSibling {
        traverse(child)
      }
    }
  }
  traverse(n)
  return firstNode
}

func transform(root *html.Node) *html.Node {
  newRoot := &html.Node{
    Type:      root.Type,
    DataAtom:  root.DataAtom,
    Data:      root.Data,
    Namespace: root.Namespace,
    Attr:      root.Attr,
  }
  transformToJsx(newRoot)

  var transitionNode *html.Node
  for child := root.FirstChild; child != nil; child = child.NextSibling {
    if child.Type == html.CommentNode {
      if data := transformRegex.FindStringSubmatch(child.Data); len(data) > 0 {
        transitionNode = &html.Node{
          Data:     "Transition",
          DataAtom: atom.Div,
          Type:     html.ElementNode,
          Attr: []html.Attribute{
            {Key: "show", Val: "{true}"},
            {Key: "enter", Val: data[1]},
            {Key: "enterFrom", Val: data[2]},
            {Key: "enterTo", Val: data[3]},
            {Key: "leave", Val: data[4]},
            {Key: "leaveFrom", Val: data[5]},
            {Key: "leaveTo", Val: data[6]},
          },
        }
      }
    } else {
      newChild := transform(child)
      if transitionNode != nil && newChild.Type == html.ElementNode {
        newChild.Attr = append([]html.Attribute{{Key: "ref", Val: "{ref}"}}, newChild.Attr...)
        transitionNode.AppendChild(&html.Node{Type: html.RawNode, Data: "\n{ref => (\n"})
        transitionNode.AppendChild(newChild)
        transitionNode.AppendChild(&html.Node{Type: html.RawNode, Data: "\n)}\n"})
        newRoot.AppendChild(transitionNode)
        transitionNode = nil
      } else {
        newRoot.AppendChild(newChild)
      }
    }
  }

  return newRoot
}

func renderNode(n *html.Node) string {
  var buf bytes.Buffer
  w := io.Writer(&buf)
  if err := html.Render(w, n); err != nil {
    panic(err)
  }
  return buf.String()
}

func transformToJsx(n *html.Node) {
  if n.Type == html.ElementNode {
    for i, attr := range n.Attr {
      switch attr.Key {
      case "class":
        attr.Key = "className"
      case "stroke-linecap":
        attr.Key = "strokeLinecap"
      case "stroke-linejoin":
        attr.Key = "strokeLinejoin"
      case "stroke-width":
        attr.Key = "strokeWidth"
      }
      n.Attr[i] = attr
    }
  }
}
