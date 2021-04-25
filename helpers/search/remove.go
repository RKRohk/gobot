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

//RemoveDocument removes a particular document from elasticsearch
func RemoveDocument(tag string, documentID string) {
	//Extracting the tag
	tag = strings.Replace(tag, "#", "", 1)

	res, err := es7.Delete(tag, documentID)

	if err != nil {
		logger.Println("Error deleting document index", err)
	} else {
		logger.Println(res)
	}
}
