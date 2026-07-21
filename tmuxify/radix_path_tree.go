package tmuxify

func GetMatch(paths *[]string , key_word string) []int{
	matches := []int{}
	lps := CreateLps(key_word)
	for i , path := range *paths {
		if KMPSearch(key_word , path , lps) > 0{
			matches = append(matches, i)
		}
	}
	return matches
}
