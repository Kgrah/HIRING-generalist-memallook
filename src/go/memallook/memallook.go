package memallook

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

type Buffer struct {
	latestTagNum int
	MaxPages     int
	PageSize     int
	PageMetaData map[int]PageData
	Pages        []byte
	TagMetaData  map[int]TagData
}

type PageData struct {
	pageNum     int
	start       int
	end         int
	isAllocated bool
	tagNum      int
}

type TagData struct {
	tagNum int
	Pages  []PageData
}

func NewBuffer(m, p int) Buffer {
	fmt.Println("Creating new buffer")
	tmd := make(map[int]TagData)
	var ps []byte

	// get page locations
	pages := make(map[int]PageData)
	t := m * p

	fmt.Println(fmt.Sprintf("num pages: %d pagesize: %d, total: %d", p, m, t))
	for i := 0; i < t; i += m {
		p := PageData{
			isAllocated: false,
			start:       i,
			end:         i + m - 1,
		}
		pages[i] = p
	}

	b := Buffer{
		MaxPages:     p,
		PageSize:     m,
		Pages:        ps,
		TagMetaData:  tmd,
		PageMetaData: pages,
	}
	fmt.Println("Returning new buffer")
	return b
}

func (b *Buffer) checkForSpace(amount int) (bool, []PageData) {
	var pages []PageData
	if len(b.TagMetaData) == 0 {
		temp := b.PageMetaData[0]
		temp.isAllocated = true
		b.PageMetaData[0] = temp
		pages = append(pages, b.PageMetaData[0])
		return true, pages
	}

	nextFreePage := b.getNextFreePage(0)
	if nextFreePage == nil {
		return false, nil
	}

	pages = append(pages)

	remaining := amount
	for remaining > 0 {
		if nextFreePage == nil {
			return false, nil
		}

		remaining -= b.PageSize
		nextFreePage = b.getNextFreePage(nextFreePage.end)
		pages = append(pages, *nextFreePage)
	}

	return true, pages
}

func (b *Buffer) getNextFreePage(idx int) *PageData {
	for i, p := range b.PageMetaData {
		if i > idx && !p.isAllocated {
			return &p
		}
	}

	return nil
}

func (b *Buffer) getNextPage(idx int) int {
	for _, l := range b.PageMetaData {
		if l.start > idx {
			return l.start
		}
	}

	return -1
}

func (b *Buffer) Alloc(amount int) error {
	newTag := b.latestTagNum + 1
	err := b.alloc(amount, newTag)
	if err != nil {
		return err
	}
	return nil
}

func (b *Buffer) alloc(amount, tagNum int) error {
	spaceExists, freePages := b.checkForSpace(amount)

	if !spaceExists {
		return fmt.Errorf("No space")
	}

	tag := TagData{
		tagNum: tagNum,
		Pages:  freePages,
	}
	b.TagMetaData[tag.tagNum] = tag
	b.latestTagNum++

	remaining := amount
	for _, p := range freePages {
		remaining -= b.PageSize
		if remaining < b.PageSize {
			for i := p.start; i < p.start+remaining; i++ {
				if i > len(b.Pages) {
					b.Pages = append(b.Pages, 0)
				}
			}
			continue
		}

		for i := p.start; i < p.start+b.PageSize-1; i++ {
			b.Pages[i] = 0
		}
	}

	b.writeNewStateToFile()
	return nil
}

func (b *Buffer) Dealloc(tagNum int) error {
	if _, ok := b.TagMetaData[tagNum]; !ok {
		return fmt.Errorf("Attempted to delete a tag that does not exist")
	}

	// naively just delete the tag info
	tag := b.TagMetaData[tagNum]
	for _, p := range tag.Pages {
		p.isAllocated = false
	}
	delete(b.TagMetaData, tagNum)

	b.writeNewStateToFile()

	return nil
}

func (b *Buffer) Clear() {
	b.MaxPages = 0
	b.PageSize = 0
	b.latestTagNum = 0
	b.Pages = make([]byte, 0)
	b.TagMetaData = make(map[int]TagData)
	b.PageMetaData = make(map[int]PageData)
	b.writeNewStateToFile()
}

func (b *Buffer) Show() string {
	gridWidth := 16

	var gridsb strings.Builder

	var counter int
	fmt.Println(fmt.Sprintf("Num pages %d", len(b.PageMetaData)))
	for _, p := range b.PageMetaData {
		if counter == gridWidth {
			gridsb.WriteString("\n")
			counter = 0
		}

		if !p.isAllocated {
			gridsb.WriteString(".")
			counter++
			continue
		}

		tagNumString := strconv.Itoa(p.tagNum)
		gridsb.WriteString(tagNumString)
		counter++
	}

	return gridsb.String()
}

func (b *Buffer) writeNewStateToFile() error {
	file, err := json.MarshalIndent(b, "", "")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile("memallook.json", file, 0644)
	if err != nil {
		return err
	}

	return nil
}

func Init(m, p int) {
	b := NewBuffer(m, p)
	b.writeNewStateToFile()
}
