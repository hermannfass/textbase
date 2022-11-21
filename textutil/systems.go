package textutil

var LineEnd = make(map[string]string, 9)  // LineEnd["darwin"] = "\r"

func init() {
	LineEnd["unix"] = "\n"
	LineEnd["macos"] = "\n"
	LineEnd["classicmac"] = "\r"
	LineEnd["cpm"] = "\r\n"
	LineEnd["msdos"] = "\r\n"
	LineEnd["os2"] = "\r\n"
	LineEnd["windows"] = LineEnd["msdos"]
	LineEnd["os9"] = "\r" // by Microware for Motorolla CPUs
	LineEnd["c64"] = "\r" // Commodore C64, C128
}


