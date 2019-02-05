package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type choices struct {
	cmd         string
	description string
	nextNode    *storyNode
	nextChoice  *choices
}

type storyNode struct {
	text    string
	choices *choices
}

func (node *storyNode) addChoice(cmd string, description string, nextNode *storyNode) {
	choice := &choices{cmd, description, nextNode, nil}

	if node.choices == nil {
		node.choices = choice
	} else {
		currentChoice := node.choices
		for currentChoice.nextChoice != nil {
			currentChoice = currentChoice.nextChoice
		}
		currentChoice.nextChoice = choice
	}
}

func (node *storyNode) render() {
	fmt.Println(node.text)
	currentChoice := node.choices
	for currentChoice != nil {
		fmt.Println(currentChoice.cmd, ":", currentChoice.description)
		currentChoice = currentChoice.nextChoice
	}
}

func (node *storyNode) executeCmd(cmd string) *storyNode {
	currentChoice := node.choices
	for currentChoice != nil {
		if strings.ToLower(currentChoice.cmd) == strings.ToLower(cmd) {
			return currentChoice.nextNode
		}
		currentChoice = currentChoice.nextChoice
	}
	fmt.Println("Sorry, I don't understand")
	return node
}

var scanner *bufio.Scanner

func (node *storyNode) play() {
	node.render()
	if node.choices != nil {
		scanner.Scan()
		node.executeCmd(scanner.Text()).play()
	}
}

func main() {
	scanner = bufio.NewScanner(os.Stdin)

	start := storyNode{text: `
	A long time ago in a galaxy far far away...
	To the left, the original trilogy, universally acclaimed.
	To the right, the sequel trilogy, made a shitload of money.
	To the top, the prequel trilogy, which as usual holds the high ground.
	`}

	originalTrilogy := storyNode{text: `
	Luke, Leia and Han blow up the Death Star. Darth Vader: NOOOO
	`}

	sequelTrilogy := storyNode{text: `
	Emo Ben wants to kill everybody. Everbody hates Rose. Rey: The garbage will do
	`}

	prequelTrilogy := storyNode{text: `
	Hello There... General Kenobi, you're a bold one
	`}

	highGround := storyNode{text: `
	It's over Anakin I have the high ground
	`}

	underestimatedPower := storyNode{text: `
	You underestimate my power
	`}

	start.addChoice("O", "Go to original trilogy", &originalTrilogy)
	start.addChoice("S", "Go to sequel trilogy", &sequelTrilogy)
	start.addChoice("P", "Go to prequel trilogy", &prequelTrilogy)

	prequelTrilogy.addChoice("H", "You go to high ground", &highGround)
	prequelTrilogy.addChoice("L", "You do spinning", &underestimatedPower)
	prequelTrilogy.addChoice("S", "Go back to start", &start)

	start.play()

	fmt.Println("Written & Directed by George Lucas")

}
