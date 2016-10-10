package main

func SplitCommand(line string) []string {
	ret := []string{}
	i := 0
	for line[i] == ' ' {
		i++
	}
	prev := i
	dq, sq := false, false
	for i < len(line) {
		switch line[i] {
		case ' ':
			if !sq && !dq {
				if prev != i {
					ret = append(ret, line[prev:i])
				}
				prev = i + 1 // skip ' '
			}
		case '\'':
			if !sq {
				if !dq {
					sq = true
					prev = i + 1 // start single quote
				}
			} else {
				ret = append(ret, line[prev:i])
				prev = i + 1
				sq = false
			}
		case '"':
			if !dq {
				if !sq {
					dq = true
					prev = i + 1
				}
			} else {
				ret = append(ret, line[prev:i])
				prev = i + 1
				dq = false
			}
		}
		i++
	}
	if prev != i {
		ret = append(ret, line[prev:i])
	}
	return ret
}
