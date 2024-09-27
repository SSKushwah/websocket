package providers

type RealtimeChatHubProvider interface {
	Get() interface{}
	Run()
	Stop()
}

type DBHelperProvider interface {
	Test()
}
