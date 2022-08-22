package intl

type Translation struct {
	Message       string
	TranslatedMsg string
}

const (
	defaultSourceMessageTbl = "source_message"
	defaultMessageTbl       = "message"
)

var (
	sourceMessageTbl = defaultSourceMessageTbl
	messageTbl       = defaultMessageTbl
)

func SetTblNames(sourceMsgTbl, msgTbl string) {
	sourceMessageTbl = sourceMsgTbl
	messageTbl = msgTbl
}

// GetMessage gets message from db by key+lang
func (i *Intl) GetMessage(key, lang string) (*Translation, error) {
	row := i.db.QueryRow("SELECT sm.message, m.translation FROM "+sourceMessageTbl+" AS sm INNER JOIN "+
		messageTbl+" AS m ON sm.id=m.id WHERE sm.category = $1 AND m.language = $2", key, lang)

	var trans Translation
	err := row.Scan(&trans.Message, &trans.TranslatedMsg)
	if err != nil {
		return nil, err
	}

	return &trans, nil
}
