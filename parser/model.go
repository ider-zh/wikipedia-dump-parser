package parser

type text struct {
	Value   string `xml:",chardata"`
	Bytes   int32  `xml:"bytes,attr"`
	Deleted string `xml:"deleted,attr"`
}

type redirect struct {
	Title string `xml:"title,attr"`
}

type contributor struct {
	Username string `xml:"username"`
	Ip       string `xml:"ip"`
	ID       int64  `xml:"id"`
	Deleted  string `xml:"deleted,attr"`
}

type revision struct {
	ID          int64       `xml:"id"`
	Text        text        `xml:"text"`
	Parentid    int64       `xml:"parentid"`
	Timestamp   string      `xml:"timestamp"`
	Comment     string      `xml:"comment"`
	Model       string      `xml:"model"`
	Format      string      `xml:"format"`
	Contributor contributor `xml:"contributor"`
}

type Page struct {
	ID        int64      `xml:"id"`
	Ns        int32      `xml:"ns"`
	Title     string     `xml:"title"`
	Redirect  *redirect  `xml:"redirect"`
	Revisions []revision `xml:"revision"`
}
