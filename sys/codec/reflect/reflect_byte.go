package reflect

func WriteChar(buff *[]byte, c byte) {
	*buff = append(*buff, c)
}
