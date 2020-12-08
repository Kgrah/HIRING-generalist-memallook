package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mHmwrk/memallook"
	"os"
	"strconv"
)

func main() {
	file, err := ioutil.ReadFile("./memallook.json")

	if err != nil {
		fmt.Println(fmt.Sprintf("Couldn't read persisted memory file: %s", err.Error()))
		return
	}

	b := memallook.Buffer{}

	err = json.Unmarshal([]byte(file), &b)

	if err != nil {
		fmt.Println(fmt.Sprintf("Couldn't unmarshal the buffer %s", err.Error()))
		return
	}

	if len(os.Args) < 2 && b.PageSize == 0 && b.MaxPages == 0 {
		fmt.Println(fmt.Sprintf("Memory buffer has not been initialized, please specify a page size and the maximum number of pages. Args: %v", len(os.Args)))
		return
	}

	if len(os.Args) < 2 {
		fmt.Println("Please specify a command")
		return
	}

	if b.MaxPages == 0 && b.PageSize == 0 {
		mArg, pArg := os.Args[1], os.Args[2]

		m, err := strconv.Atoi(mArg)
		if err != nil {
			fmt.Println(fmt.Sprintf("Could not parse page size. Args %v", mArg))
			return
		}

		p, err := strconv.Atoi(pArg)
		if err != nil {
			fmt.Println(fmt.Sprintf("Could not parse page size. Args %v", os.Args))
			return
		}

		memallook.Init(m, p)
	}

	command := os.Args[1]
	fmt.Println(fmt.Sprintf("Test command %s", command))

	switch command {
	case "alloc":
		fmt.Println("Case alloc")
		mArg := os.Args[2]
		m, err := strconv.Atoi(mArg)
		if err != nil {
			fmt.Println(fmt.Sprintf("Couldn't parse alloc bytes argument. Err: %s %s", err.Error(), mArg))
			return
		}

		err = b.Alloc(m)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	case "clear":
		b.Clear()
	case "show":
		grid := b.Show()
		fmt.Println(grid)
	case "dealloc":
		tArg := os.Args[1]

		t, err := strconv.Atoi(tArg)
		if err != nil {
			fmt.Println("Couldn't parse dealloc tag argument")
			return
		}

		err = b.Dealloc(t)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}

}
