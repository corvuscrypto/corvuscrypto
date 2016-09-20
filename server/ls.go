package main

func executeLS(args []string) []byte {
	if len(args) == 0 {
		return []byte("ls: cannot access : Permission denied")
	}

	//lets give a mock access denied message for anything before
	// /home/Guest/
	if args[0] != "" && args[0][:11] != "/home/Guest/" {
		return []byte("ls: cannot access " + args[0] + ": Permission denied")
	}
}
