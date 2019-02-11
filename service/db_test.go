package service

// func TestSingleUrl(t *testing.T) {
// 	url := "https://www.google.com/search?q=golang+specification"
// 	_, decoded, err := process(url, int64(1))
// 	if err != nil {
// 		t.Errorf("error fetching url: %s\n", err)
// 	} else if strings.Compare(decoded, url) != 0 {
// 		t.Errorf("urls does not match: %s and %s\n", url, decoded)
// 	}
// }

// func TestMultipleUrls(t *testing.T) {
// 	urls := []string{
// 		"stackoverflow.com",
// 		"blog.golang.org",
// 		"play.google.com/music/",
// 		"gist.github.com",
// 	}
// 	id := int64(1)
// 	for _, url := range urls {
// 		_, decoded, err := process(url, id)
// 		id++
// 		if strings.Compare(decoded, url) != 0 || err != nil {
// 			t.Errorf("url processed wrong: %v", url)
// 		}
// 	}
// }

// func TestDuplicateUrl(t *testing.T) {
// 	url := "https://reddit.com"
// 	short1, _, err := process(url, int64(1))
// 	short2, _, err := process(url, int64(2))
// 	fmt.Printf("short1: %v, short2: %v", short1, short2)
// 	if strings.Compare(short1, short2) != 0 || err != nil {
// 		t.Errorf("same url processed twice: %v", url)
// 	}
// }

// func TestInvalidUrl(t *testing.T) {
// 	url := "NotAValiURL"
// 	_, _, err := process(url, int64(1))
// 	if err != nil {
// 		t.Errorf("non valid url shortened: %v", url)
// 	}
// }

// func process(data string, id int64) (string, string, error) {
// 	var key []byte
// 	conn := redigomock.NewConn()
// 	conn.Clear()
// 	// data
// 	url, err := ParseURL(data, &http.Request{})
// 	if err != nil {
// 		return "", "", err
// 	}
// 	// saving & processing
// 	key = insertPrefixInKey(urlIDCounter, string(IdKey))
// 	conn.Command("GET", key).Expect(id)
// 	conn.Command("INCR", key).Expect(id + 1)
// 	key = insertPrefixInKey(urlKey, string(id))
// 	conn.Command("SET", key, data).Expect("ok")

// 	shorten, err := url.SaveURL(conn)
// 	if err != nil {
// 		return "", "", err
// 	}
// 	// fetching
// 	conn.Clear()
// 	key = insertPrefixInKey(urlKey, string(id))
// 	conn.Command("GET", key).Expect(data)
// 	dbURL, err := DbGetURL(shorten, conn)
// 	return shorten, dbURL, err
// }
