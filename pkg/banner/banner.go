package banner

import (
	"fmt"
	"github.com/fatih/color"
)

/*
创建人员：云深不知处
创建时间：2022/1/13
程序功能：指针
*/

func Banner() {
	banner := " ____            ____\n" +
		"| __ )  ___  ___/ ___|  ___ __ _ _ __\n" +
		"|  _ \\ / _ \\/ _ \\___ \\ / __/ _` | '_ \\\n" +
		"| |_) |  __/  __/___) | (_| (_| | | | |\n" +
		"|____/ \\___|\\___|____/ \\___\\__,_|_| |_| version:0.2.0\n" + "\n"
	_, _ = fmt.Fprintf(color.Output, color.HiCyanString(banner))
}
