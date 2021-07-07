package main

func main() {

	switch {
	case remote.Used:
		{
			if list.Used {
				OP = "list"
			} else if add.Used {
				OP = "add"
			} else if delete.Used {
				OP = "delete"
			}
		}
	}

}
