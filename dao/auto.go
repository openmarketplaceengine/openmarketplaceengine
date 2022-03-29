package dao

func AutoExec(exec ...Executable) {
	Pgdb.auto.Append(exec...)
}

func AutoDrop(exec ...Executable) {
	Pgdb.drop.Append(exec...)
}
