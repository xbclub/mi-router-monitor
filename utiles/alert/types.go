package alert

type MiRouterStatus struct {
	Code  int `json:"code"`
	Count struct {
		All               int `json:"all"`
		AllWithoutMash    int `json:"all_without_mash"`
		Online            int `json:"online"`
		OnlineWithoutMash int `json:"online_without_mash"`
	} `json:"count"`
	Cpu struct {
		Core int    `json:"core"`
		Hz   string `json:"hz"`
		Load int    `json:"load"`
	} `json:"cpu"`
	Dev      []Dev `json:"dev"`
	Hardware struct {
		DisplayRomVer string `json:"DisplayRomVer"`
		Channel       string `json:"channel"`
		DisplayName   string `json:"displayName"`
		IspName       string `json:"ispName"`
		Mac           string `json:"mac"`
		Platform      string `json:"platform"`
		Sn            string `json:"sn"`
		Version       string `json:"version"`
	} `json:"hardware"`
	Mem struct {
		Hz    string  `json:"hz"`
		Total string  `json:"total"`
		Type  string  `json:"type"`
		Usage float64 `json:"usage"`
	} `json:"mem"`
	Temperature int    `json:"temperature"`
	UpTime      string `json:"upTime"`
	Wan         struct {
		Devname          string `json:"devname"`
		Download         int64  `json:"download,string"`
		Downspeed        int64  `json:"downspeed,string"`
		Maxdownloadspeed int64  `json:"maxdownloadspeed,string"`
		Maxuploadspeed   int64  `json:"maxuploadspeed,string"`
		Upload           int64  `json:"upload,string"`
		Upspeed          int64  `json:"upspeed,string"`
	} `json:"wan"`
}
type Dev struct {
	Devname          string      `json:"devname"`
	Download         interface{} `json:"download"`
	Downspeed        interface{} `json:"downspeed"`
	Isap             int         `json:"isap,omitempty"`
	Mac              string      `json:"mac"`
	Maxdownloadspeed string      `json:"maxdownloadspeed"`
	Maxuploadspeed   string      `json:"maxuploadspeed"`
	Online           string      `json:"online"`
	Upload           interface{} `json:"upload"`
	Upspeed          interface{} `json:"upspeed"`
	Ip               string      `json:"ip,omitempty"`
}
