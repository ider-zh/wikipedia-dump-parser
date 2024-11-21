package wikiparser

// json schema 注释
type Text struct {
	Value   string `xml:",chardata" json:"value"`
	Bytes   int32  `xml:"bytes,attr" json:"bytes"`
	Deleted string `xml:"deleted,attr" json:"deleted"`
}

type Redirect struct {
	Title string `xml:"title,attr" json:"title"`
}

type Contributor struct {
	Username string `xml:"username" json:"username"`
	Ip       string `xml:"ip" json:"ip"`
	ID       int64  `xml:"id" json:"id"`
	Deleted  string `xml:"deleted,attr" json:"deleted"`
}

type Revision struct {
	ID          int64       `xml:"id" json:"id"`
	Text        Text        `xml:"text" json:"text"`
	Parentid    int64       `xml:"parentid" json:"parentid"`
	Timestamp   string      `xml:"timestamp" json:"timestamp"`
	Comment     string      `xml:"comment" json:"comment"`
	Model       string      `xml:"model" json:"model"`
	Format      string      `xml:"format" json:"format"`
	Contributor Contributor `xml:"contributor" json:"contributor"`
}

type Page struct {
	ID        int64      `xml:"id" json:"id"`
	Ns        int32      `xml:"ns" json:"ns"`
	Title     string     `xml:"title" json:"title"`
	Redirect  *Redirect  `xml:"redirect" json:"redirect"`
	Revisions []Revision `xml:"revision" json:"revisions"`
}
