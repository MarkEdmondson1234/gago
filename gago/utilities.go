package gago

import (
	"fmt"
	"log"
	"strings"
	"time"
)

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

func deleteEmptyStringSlice(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}

func deleteEmptyRowSlice(s []*ParseReportRow) []*ParseReportRow {
	var r []*ParseReportRow
	for _, str := range s {
		if str != nil {
			r = append(r, str)
		}
	}
	return r
}

func join(strs ...string) string {
	var sb strings.Builder
	for _, str := range strs {
		sb.WriteString(str)
	}
	return sb.String()
}

func myMessage(s string) {

	if verbose {
		fmt.Println(time.Now().Format("Mon Jan _2 15:04:05 2006"), ">", s)
	} else {
		log.Println(join("gago: ", s))
	}

}
