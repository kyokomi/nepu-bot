package nepubot

import "math/rand"

var random = rand.New(rand.NewSource(1))

// Kaomoji 顔文字
var kaomoji = []string{
	"(; ・∀・)",
	"(~_~;)",
	"(-_-;)",
	"?(°_°>)",
	"Σ(￣□￣;)",
	"( ｀・ω・´)",
	"m9( ﾟдﾟ)",
}

func GetKaomji() string {
	idx := random.Int31n((int32)(len(kaomoji) - 1))
	return kaomoji[idx]
}
