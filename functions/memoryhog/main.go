package main

func main() {

	x := make([]byte, 1024*1024*1024)

	for i := range x {
		x[i] = 1
	}

	select {}
}
