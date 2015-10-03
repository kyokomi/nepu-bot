package twitter

import (
	"fmt"
	"testing"
)

func _TestFavorites(t *testing.T) {
	accessToken, err := newAccessToken("", "")
	if err != nil {
		t.Fatal(err.Error())
	}

	//	if images, err := searchImages(accessToken, "ネプテューヌ", 4); err != nil {
	//		t.Fatal(err.Error())
	//	} else {
	//		fmt.Println(images)
	//	}

	if images, err := getFavoritesImage(accessToken, 10, "keikenchi_bot"); err != nil {
		t.Fatal(err.Error())
	} else {
		for _, image := range images {
			fmt.Println(image)
		}
	}
}
