package search

import "strings"

//RemoveIndex removes the entire index from elasticsearch
func RemoveIndex(tag string) {
	//Extracting the tag
	tag = strings.Replace(tag, "#", "", 1)

	res, err := es7.Indices.Delete([]string{tag})
	if err != nil {
		logger.Println("Error deleting index", err)
	} else {
		logger.Println(res)
	}
}

func RemoveDocument() {

}
