package main

import (
	"errors"
	"fmt"
)

// args[0] is the profile to be choosen
func chooseHandler(args []string, data interface{}) (int, error) {
	if len(args) == 0 || args[0] == "" {
		return 0, errors.New("'choose' needs name or index of profile.")
	}
	p, index, err := doGetProfile(args[0], gConfig)
	if err != nil {
		return 0, err
	}

	if p == nil {
		msg := fmt.Sprintf("can not find profile '%s'\n", args[0])
		return 0, errors.New(msg)
	}

	fmt.Printf("choose [%s]\n", cWrap(cGREEN, p.Name))
	gConfig.Current = index
	gConfig.save()

	return 0, nil
}
