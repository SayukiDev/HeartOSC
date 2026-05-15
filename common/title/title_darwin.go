package title

import "fmt"

func SetTitle(title string) error {
	fmt.Printf("\033]0;%s\007", title)
	return nil
}
