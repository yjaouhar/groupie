package main 


func main(){

}


func longestCommonPrefix(strs []string) string {
	var s string
	var bo bool
	v := strs[0]
	for i := 0; i < nb(strs); i++ {
		if !bo {
			return s
		}
		if bo && i!=0{
			s += string(v[i-1])
		}
		for x := 1; x <len(strs) ; x++ {
			if v[i]== strs[x][i]{
				bo = true
			}else{
				bo = false
				break
			}

			
		}
		
	}
	return s
    
}
func nb(nb []string)int{
	var n int 
	n = len(nb[0])
	for i := 0; i < len(nb); i++ {
		if len(nb[i])<n {
			n = len(nb[i])
		}
		
	}
	return n
}





func isValid(s string) bool {
    var b bool
	sl := []rune(s)
	for i := 0; i < len(sl); i++ {
		if sl[i]=='(' ||sl[i]=='{' ||sl[i]=='[' && i != len(sl)-1{
			bo , in :=chek(sl[i+1:],sl[i])
			if !bo{
				return false
			}
			if bo{
				sl[i]='0'
				sl[i+in]='0'
				b = true
			}

		}
	}
	return b
}
func chek(sl []rune , s rune) (bool , int) {
	b := false
	var n int
	var si rune
	if s=='(' {
		si= ')'
	}
	if s=='[' {
		si= ']'
	}
	if s=='{' {
		si= '}'
	}
	for i := 0; i < len(sl); i++ {
		if sl[i]==si{
			b = true
			n = i+1
			break
		}
	}
	return b , n
}