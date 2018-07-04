package user

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	origin = "https://www.khanacademy.org/api/internal"
)

// Comments is an array of type Comment
type Comments []Comment

// Comment holds a user's comment data
type Comment struct {
	Content string
}

// User describes a KA user
type User struct {
	Username string
	Comments Comments
}

// GenerateComments grabs a certain amount of pages of discussion
func (u *User) GenerateComments(pages int) error {
	for i := 0; i < pages; i++ {
		channel := make(chan Comments, 1)
		u.getDiscussion(i, channel)
		c := <-channel
		for _, comment := range c {
			u.Comments = append(u.Comments, comment)
		}
	}
	return nil
}

func (u *User) getDiscussion(page int, c chan Comments) {
	var reqJSON Comments

	url := fmt.Sprintf("%s/user/replies?username=%s&page=%d", origin, u.Username, page)

	req, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}

	defer req.Body.Close()
	reqBytes, _ := ioutil.ReadAll(req.Body)

	json.Unmarshal(reqBytes, &reqJSON)

	c <- reqJSON
}
