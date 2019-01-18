package main

import (
	"bytes"
	"encoding/json"
	"os/exec"
	"strconv"
)

type CompleteItem struct {
	name        string
	description string
}

func AutoComplete(fileType string, url string, globalOffset int, buffer *bytes.Buffer) {
	switch fileType {
	case "Go":
		GoAutoComplete(url, globalOffset, buffer)
	case "Python":
		PythonAutoComplete(url, globalOffset, buffer)
	default:
		return
	}
}

func GoAutoComplete(url string, globalOffset int, buffer *bytes.Buffer) {
	var resultBuffer bytes.Buffer

	c := exec.Command("gocode", "-f=json", "-unimported-packages", "autocomplete", url, "c"+strconv.Itoa(globalOffset))
	c.Stdin = buffer
	c.Stdout = &resultBuffer
	c.Run()

	var res []interface{}
	json.Unmarshal(resultBuffer.Bytes(), &res)

	if res == nil {
		return
	}

	completeOffset := int(res[0].(float64))

	var completeItems []CompleteItem
	for _, item := range res[1].([]interface{}) {
		completeItemName := item.(map[string]interface{})["name"].(string)
		completeItemDes := item.(map[string]interface{})["type"].(string)
		completeItem := CompleteItem{
			name:        completeItemName,
			description: completeItemDes}
		completeItems = append(completeItems, completeItem)
	}

	autocompleteList.LoadItems(completeOffset, completeItems)
	textEditor.Display()
}

func PythonAutoComplete(url string, globalOffset int, buffer *bytes.Buffer) {
	var resultBuffer bytes.Buffer

	c := exec.Command("python3", "/home/hankelbao/Projects/jedi_test/test.py", url, strconv.Itoa(globalOffset))
	c.Stdin = buffer
	c.Stdout = &resultBuffer
	c.Run()

	var res []interface{}
	json.Unmarshal(resultBuffer.Bytes(), &res)

	if res == nil {
		return
	}

	completeOffset := int(res[0].(float64))

	var completeItems []CompleteItem
	for _, item := range res[1].([]interface{}) {
		completeItemName := item.(map[string]interface{})["name"].(string)
		completeItemDes := item.(map[string]interface{})["description"].(string)
		completeItem := CompleteItem{
			name:        completeItemName,
			description: completeItemDes}
		completeItems = append(completeItems, completeItem)
	}

	autocompleteList.LoadItems(completeOffset, completeItems)
	textEditor.Display()
}
