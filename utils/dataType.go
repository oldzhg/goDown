package utils

type filePart struct {
	Index int
	Start int
	End   int
	Data  []byte
}

type fileInfo struct {
	Size      int
	URL       string
	Name      string
	Path      string
	MD5       string
	DoneParts []filePart
}
