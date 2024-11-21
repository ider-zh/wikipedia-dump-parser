package parser

// json schema 注释
type text struct {
	Value   string `xml:",chardata" json:"value"`
	Bytes   int32  `xml:"bytes,attr" json:"bytes"`
	Deleted string `xml:"deleted,attr" json:"deleted"`
}

type redirect struct {
	Title string `xml:"title,attr" json:"title"`
}

type contributor struct {
	Username string `xml:"username" json:"username"`
	Ip       string `xml:"ip" json:"ip"`
	ID       int64  `xml:"id" json:"id"`
	Deleted  string `xml:"deleted,attr" json:"deleted"`
}

type revision struct {
	ID          int64       `xml:"id" json:"id"`
	Text        text        `xml:"text" json:"text"`
	Parentid    int64       `xml:"parentid" json:"parentid"`
	Timestamp   string      `xml:"timestamp" json:"timestamp"`
	Comment     string      `xml:"comment" json:"comment"`
	Model       string      `xml:"model" json:"model"`
	Format      string      `xml:"format" json:"format"`
	Contributor contributor `xml:"contributor" json:"contributor"`
}

type Page struct {
	ID        int64      `xml:"id" json:"id"`
	Ns        int32      `xml:"ns" json:"ns"`
	Title     string     `xml:"title" json:"title"`
	Redirect  *redirect  `xml:"redirect" json:"redirect"`
	Revisions []revision `xml:"revision" json:"revisions"`
}
