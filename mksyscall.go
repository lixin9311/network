package network

// modify freeMibTable after go generate
// func freeMibTable(table nMIBTable) (ret error) {
// 	r0, _, _ := syscall.Syscall(procFreeMibTable.Addr(), 1, table.unsafePointer(), 0, 0)
//go:generate go run golang.org/x/sys/windows/mkwinsyscall -output zsyscall_windows.go netio_windows.go
