package main

import "time"

type Entry struct {
	ID      int
	Content string
	Date    string
}

type Entries []Entry

var id = 0
var data = Entries{}

func NewEntry(content string) Entry {
	id++
	return Entry{id, content, time.Now().Format(time.ANSIC)}
}

func AddEntry(content string) {
	data = append(data, NewEntry(content))
}

func RemoveEntry(id int) {
	for i, entry := range data {
		if entry.ID == id {
			data = append(data[:i], data[i+1:]...)
			break
		}
	}
}

func AllEntries() Entries {
	return data
}
