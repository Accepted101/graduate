package platform

type WordVal struct {
	WordId   int64
	WordName string
	WordVal  uint
	Tmp      int
}

func (WordVal) TableName() string {
	return "word_val"
}
