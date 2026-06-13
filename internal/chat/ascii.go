package chat

// GetWelcomeLogo returns the structured Penguin welcome banner for the TCP chat.
func GetWelcomeLogo() string {
	logo := "Welcome to TCP-Chat!\n"
	logo += "         _nnnn_\n"
	logo += "        dGGGGMMb\n"
	logo += "       @p~qp~~qMb\n"
	logo += "       M|@||@) M|\n"
	logo += "       @,----.JM|\n"
	logo += "      JS^\\__/  qKL\n"
	logo += "     dZP        qKRb\n"
	logo += "    dZP          qKKb\n"
	logo += "   fZP            SMMb\n"
	logo += "   HZM            MMMM\n"
	logo += "   FqM            MMMM\n"
	logo += " __| \".        |\\dS\"qML\n"
	logo += " |    `.       | `' \\Zq\n"
	logo += "_)      \\.___,|     .'\n"
	logo += "\\____   )MMMMMP|   .'\n"
	logo += "     `-'       `--'\n"
	logo += "[ENTER YOUR NAME]: "
	return logo
}