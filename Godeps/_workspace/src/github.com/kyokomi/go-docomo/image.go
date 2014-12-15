package docomo

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

var logger = log.New(os.Stderr, "image", log.Llongfile)

func (d *DocomoClient) SendImage(imageURL string) ([]byte, error) {
	fmt.Println(imageURL)

	// TODO: downloadする

	res, err := d.PostImage(imageURL)
	if err != nil {
		logger.Println(err)
		return nil, err
	}

	resData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logger.Println(err)
		return nil, err
	}

	return resData, nil
}
