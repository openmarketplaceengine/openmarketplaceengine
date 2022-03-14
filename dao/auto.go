package dao

var autoExec = make([]Executable, 0, 16)

func AutoExec(exec ...Executable) {
	autoExec = append(autoExec, exec...)
}
