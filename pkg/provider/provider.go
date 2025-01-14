package provider

const (
	PROVIDER_ALIYUN = "aliyun"
)

type PhoneCall interface {
	Call(message string, phoneNumbers []string) error
}
