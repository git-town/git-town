package main

import (
	"fmt"
	"os"
	"time"
)

type release struct {
	Time    time.Time
	ISOTime string
}

func loadPreviousRelease() release {
	if len(os.Args) < 2 {
		fmt.Println("Usage: list_contributors <previous tag>")
		os.Exit(1)
	}
	tag := os.Args[1]
	tagTime := releaseDate(tag)
	tagTimeISO := tagTime.Format("2006-01-02")
	fmt.Printf("previous release %s was on %s\n", cyan.Styled(tag), cyan.Styled(tagTimeISO))
	return release{
		Time:    tagTime,
		ISOTime: tagTimeISO,
	}
}
