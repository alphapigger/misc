package main

import (
	"fmt"
	"regexp"
)

var (
	validUUID  = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)
	validToken = regexp.MustCompile(`^[0-9a-fA-F]{64}$`)
	validAlias = regexp.MustCompile(`^[\w=\-/]{1,40}$`)
)

func main() {
	fmt.Println("match uuid haha: ", validUUID.MatchString("haha"))
	fmt.Println("match uuid 0b574a70-60d6-407b-8013-61f09b824caf: ", validUUID.MatchString("0b574a70-60d6-407b-8013-61f09b824caf"))
	fmt.Println("match uuid 0b574a70-60d6-407b-8013-61f09b824caF: ", validUUID.MatchString("0b574a70-60d6-407b-8013-61f09b824caF"))
	fmt.Println("match uuid 0b574a70-60d6-407b-8013-61f09b824caf1: ", validUUID.MatchString("0b574a70-60d6-407b-8013-61f09b824caf1"))
	fmt.Println("match token heihei: ", validToken.MatchString("heihei"))
	fmt.Println("match token 85c1bcc4e45843beb627456942ee297c5fb923a4c8de10effe06493d3a8db77d: ", validToken.MatchString("85c1bcc4e45843beb627456942ee297c5fb923a4c8de10effe06493d3a8db77d"))
	fmt.Println("match token 85c1bcc4e45843beb627456942ee297c5fb923a4c8de10effe06493d3a8db77d1: ", validToken.MatchString("85c1bcc4e45843beb627456942ee297c5fb923a4c8de10effe06493d3a8db77d1"))
	fmt.Println("match alias haha_-=/", validAlias.MatchString("haha_-=/"))
	fmt.Println("match alias \"\"", validAlias.MatchString(""))
}
