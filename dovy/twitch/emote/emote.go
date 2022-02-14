package emote

import "fmt"

type Store struct {
}

func ImgTag(src string) string {
	return fmt.Sprintf("<img src=\"%s\" class=\"emote\"/>", src)
}
