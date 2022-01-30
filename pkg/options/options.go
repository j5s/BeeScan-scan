package options

/*
创建人员：云深不知处
创建时间：2022/1/3
程序功能：运行实例的选项
*/


type Options struct {
	Target              multiStringFlag
	Targets             string
	Output              string
	ProxyURL            string
	TimeOut             int
	JSON                bool
	Verbose             bool
	OutputStatusCode    bool
	OutputWithNoColor   bool
	OutputContentLength bool
	OutputTitle         bool
	OutputIP            bool
	OutputFingerPrint   bool
	RateLimit           int
	OutputCDN           bool
}
type multiStringFlag []string

func (m *multiStringFlag) String() string {
	return ""
}

func (m *multiStringFlag) Set(value string) error {
	*m = append(*m, value)
	return nil
}

