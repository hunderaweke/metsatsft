package main

import (
	"fmt"

	"github.com/hunderaweke/metsasft/pkg"
)

func main() {
	err := pkg.SendResetEmail("hunderaweke@gmail.com", "token")
	if err != nil {
		panic(err)
	}
	fmt.Println("Email sent successfully")
}
