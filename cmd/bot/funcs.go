package main

import (
	"fmt"
	"net/url"
)

func formatN(n int) string {
	if n >= 1000000 {
		if n2 := n / 100000 % 10; n2 != 0 {
			return fmt.Sprint(n/1000000, ".", n2, " M")
		}
		return fmt.Sprint(n/1000000, " M")
	}
	if n >= 1000 {
		if n2 := n / 100 % 10; n2 != 0 {
			return fmt.Sprint(n/1000, ".", n2, " K")
		}
		return fmt.Sprint(n/1000, " K")
	}
	return fmt.Sprint(n)
}

func host(l string) string {
	u, err := url.ParseRequestURI(l)
	if err != nil {
		return "-"
	}
	return u.Host
}
