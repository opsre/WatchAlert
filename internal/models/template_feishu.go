package models

// JsonCards 飞书卡片，兼容 1.0 / 2.0 的field
type JsonCards struct {
	Schema       string                   `json:"schema,omitempty"`
	Config       map[string]interface{}   `json:"config,omitempty"`
	CardLink     map[string]interface{}   `json:"card_link,omitempty"`
	Header       map[string]interface{}   `json:"header,omitempty"`
	Body         map[string]interface{}   `json:"body,omitempty"`
	Elements     []map[string]interface{} `json:"elements,omitempty"`
	I18nElements map[string]interface{}   `json:"i18n_elements,omitempty"`
	Fallback     map[string]interface{}   `json:"fallback,omitempty"`
}

// FeiShuJsonCardMsg 飞书Json卡片消息结构体
type FeiShuJsonCardMsg struct {
	MsgType string    `json:"msg_type"`
	Card    JsonCards `json:"card"`
}

type Actions struct {
	Tag      string      `json:"tag"`
	Text     ActionsText `json:"text"`
	Type     string      `json:"type"`
	Value    interface{} `json:"value"`
	Confirm  Confirms    `json:"confirm"`
	URL      string      `json:"url"`
	MultiURL *MultiURLs  `json:"multi_url"`
}

type MultiURLs struct {
	URL        string `json:"url"`
	AndroidURL string `json:"android_url"`
	IosURL     string `json:"ios_url"`
	PcURL      string `json:"pc_url"`
}

type Confirms struct {
	Title Titles `json:"title"`
	Text  Texts  `json:"text"`
}

type ActionsText struct {
	Content string `json:"content"`
	Tag     string `json:"tag"`
}

type Configs struct {
	WideScreenMode bool      `json:"wide_screen_mode,omitempty"` // 最新文档中没有找到该配置项，先保留，允许为空
	WidthMode      WidthMode `json:"width_mode,omitempty"`       // 卡片宽度模式。支持 "compact"（紧凑宽度 400px）模式、"fill"（撑满聊天窗口宽度）模式和 "default" 默认模式(宽度上限为 600px)。
	EnableForward  bool      `json:"enable_forward,omitempty"`
}

type WidthMode string

const (
	WidthModeDefault WidthMode = "default" // 宽度上限 600px
	WidthModeFill    WidthMode = "fill"    // 撑满聊天窗口宽度
	WidthModeCompact WidthMode = "compact" // 紧凑宽度 400px
)

type Elements struct {
	Tag            string             `json:"tag"`
	FlexMode       string             `json:"flexMode"`
	BackgroupStyle string             `json:"background_style"`
	Text           Texts              `json:"text"`
	Columns        []Columns          `json:"columns"`
	Elements       []ElementsElements `json:"elements"`
}

type ElementsElements struct {
	Tag     string `json:"tag"`
	Content string `json:"content"`
}

type Columns struct {
	Tag           string            `json:"tag"`
	Width         string            `json:"width"`
	Weight        int64             `json:"weight"`
	VerticalAlign string            `json:"vertical_align"`
	Elements      []ColumnsElements `json:"elements"`
}

type ColumnsElements struct {
	Tag  string `json:"tag"`
	Text Texts  `json:"text"`
}

type Texts struct {
	Content string `json:"content"`
	Tag     string `json:"tag"`
}

type Headers struct {
	Template string `json:"template"`
	Title    Titles `json:"title"`
}

type Titles struct {
	Content string `json:"content"`
	Tag     string `json:"tag"`
}

// CardInfo 飞书回传
type CardInfo struct {
	OpenID        string         `json:"open_id"`
	UserID        string         `json:"user_id"`
	OpenMessageID string         `json:"open_message_id"`
	OpenChatID    string         `json:"open_chat_id"`
	TenantKey     string         `json:"tenant_key"`
	Token         string         `json:"token"`
	Action        CardInfoAction `json:"action"`
}

type CardInfoAction struct {
	Value SilenceValue `json:"value"`
	Tag   string       `json:"tag"`
}

type SilenceValue struct {
	Comment   string           `json:"comment"`
	CreatedBy string           `json:"createdBy"`
	EndsAt    string           `json:"endsAt"`
	Id        string           `json:"id"`
	Matchers  []MatchersLabels `json:"matchers"`
	StartsAt  string           `json:"startsAt"`
}

type MatchersLabels struct {
	IsEqual bool   `json:"isEqual"`
	IsRegex bool   `json:"isRegex"`
	Name    string `json:"name"`
	Value   string `json:"value"`
}

// FeiShuUserInfo 飞书用户信息
type FeiShuUserInfo struct {
	Data Data `json:"data"`
}

type Data struct {
	User User `json:"user"`
}

type User struct {
	UserId string `json:"user_id"`
	Name   string `json:"name"`
}

// FeiShuChats 机器人所在群列表
type FeiShuChats struct {
	HasMore bool    `json:"has_more"`
	Items   []Items `json:"items"`
}

type Items struct {
	Name    string `json:"name"`
	ChatId  string `json:"chat_id"`
	OwnerId string `json:"owner_id"`
}
