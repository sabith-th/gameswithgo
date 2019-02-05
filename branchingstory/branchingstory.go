package main

import (
	"bufio"
	"fmt"
	"os"
)

type storyNode struct {
	text    string
	yesPath *storyNode
	noPath  *storyNode
}

func (node *storyNode) printStory() {
	fmt.Println(node.text)
	if node.yesPath != nil {
		node.yesPath.printStory()
	}
	if node.noPath != nil {
		node.noPath.printStory()
	}
}

func (node *storyNode) play() {
	fmt.Println(node.text)

	if node.yesPath != nil && node.noPath != nil {
		scanner := bufio.NewScanner(os.Stdin)
		for {
			scanner.Scan()
			answer := scanner.Text()
			if answer == "yes" {
				node.yesPath.play()
				break
			} else if answer == "no" {
				node.noPath.play()
				break
			} else {
				fmt.Println("Only a sith deals in absolute")
			}
		}
	}

}

func main() {
	root := storyNode{"A long time ago in a galaxy far far away", nil, nil}
	winning := storyNode{"The empire has been defeated", nil, nil}
	losing := storyNode{"The empire won, Jedi eradicated", nil, nil}
	root.yesPath = &winning
	root.noPath = &losing

	root.play()
	root.printStory()
}
