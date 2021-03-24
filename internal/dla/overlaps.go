package dla

import (
	"go.matteson.dev/gfx"
)

// RemoveOverlappingLetters ...
func RemoveOverlappingLetters(letters gfx.Chars) gfx.Chars {
	if len(letters) == 0 {
		return letters
	}

	queue := letters[1:]
	cleanLetters := gfx.Chars{letters[0]}

	for len(queue) > 0 {
		letter := queue[0]
		queue = queue[1:]

		addLetter := true

		for _, cleanLetter := range cleanLetters {
			if cleanLetter.Rune == letter.Rune && letter.Quad.Bounds().Intersects(cleanLetter.Quad.Bounds()) {
				addLetter = false
				break
			}
		}

		if addLetter {
			cleanLetters = append(cleanLetters, letter)
		}
	}
	return cleanLetters
}
