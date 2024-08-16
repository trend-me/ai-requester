package interfaces

type Ai interface {
	Prompt(prompt string) (string, error)
}