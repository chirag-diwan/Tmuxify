package tmuxify

func CreateLps(pat string) []int{
	lps := make([]int ,len(pat))

	length := 0
	i := 1

	for i < len(pat){
		if pat[i] == pat[length]{
			length++
			lps[i] = length
			i ++
		}else{
			if length != 0{
				length = lps[length - 1]
			}else{
				lps[i] = 0
				i ++
			}
		}
	}

	return lps
}

func KMPSearch(pat string , txt string , lps []int) int {
	if len(pat) == 0{
		return 1
	}

	if len(txt) == 0{
		return 0
	}

	count := 0
	j := 0
	i := 0
	for i < len(txt){
		if pat[j] == txt[i] {
			j++
			i++

			if j == len(pat){
				count ++;
				j = lps[j - 1];
			}
		}else{
			if j != 0{
				j = lps[j - 1];
			} else{
				i++;
			}
		}
	}
	return count
}
