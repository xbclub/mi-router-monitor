package wechat

type SendText struct {
	Touser                   string `json:"touser"`
	Toparty                  string `json:"toparty"`
	Totag                    string `json:"totag"`
	Msgtype                  string `json:"msgtype"`
	Agentid                  int    `json:"agentid"`
	Text                     Text   `json:"text"`
	Textcard                 Text   `json:"textcard"`
	Safe                     int    `json:"safe"`
	Enable_id_trans          int    `json:"enable_id_trans"`
	Enable_duplicate_check   int    `json:"enable_duplicate_check"`
	Duplicate_check_interval int    `json:"duplicate_check_interval"`
}

type Text struct {
	Content     string `json:"content"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Url         string `json:"url"`
	Btntext     string `json:"btntext"`
}
