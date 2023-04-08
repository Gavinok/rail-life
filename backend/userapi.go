package main

import (
	"log"
	"time"
)

func (u *user_doc) getPosts() (*[]post_doc, error) {
	return db.lookupPostsByAuthor(u)
}

type user_stats struct {
	Posts    *[]post_doc
	Comments *[]comment_doc
}

func (u *user_doc) getNotifications() []string {
	source := trackNotifications(*u, mqconn)
	notifiactions := make([]string, 0)
	for len(notifiactions) < 10 {
		notifiactions = append(notifiactions, source.Get())
		log.Println("notification found ", notifiactions)
	}
	return notifiactions

}
func (u *user_doc) getStats() (*user_stats, error) {
	posts, err := u.getPosts()
	comments, err := getUserStat[comment_doc](u, "comments")

	return &user_stats{
		posts, comments,
	}, err
}

func (u *user_doc) Delete(db *dataStore) error {
	return db.deleteUser(u)
}
func (u *user_doc) newPost(title string, content string, db *dataStore) (post_doc, error) {
	p := post_doc{
		Title:          title,
		Content:        content,
		Author:         u.Username,
		DateOfCreation: time.Now(),
		Likes:          make([]Username, 0),
		Comments:       make([]comment_doc, 0),
	}
	err := db.storePost(p)
	u.notifyFriends(POST_NOTIFICATION, content, db)
	return p, err
}
func (u *user_doc) newFriend(name string, db *dataStore) (*user_doc, error) {
	friend, err := read_user(Username(name), db)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	// u.Friends = append(u.Friends, friend.Username)
	_, err = friend.addFriend(u.Username, db)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	me, err := u.addFriend(friend.Username, db)
	if err != nil {
		log.Println("NewFriend, ", err)
		return nil, err
	}
	return me, err
}
func (u *user_doc) newComment(title string, content string, db *dataStore) (comment_doc, error) {

	c := comment_doc{
		Content:        content,
		Author:         u.Username,
		DateOfCreation: time.Now(),
		Likes:          make([]Username, 0),
	}
	result, err := db.storeComment(title, c)

	log.Println("the result of inserting this comment was ", result)

	// Notify the original author
	post_Author, err := read_user(result.Author, db)
	log.Println("notifying  ", post_Author, "with new comment")
	post_Author.notifyThisUser(COMMENT_NOTIFICATION, content, db)

	return c, err
}

func (u *user_doc) notifyThisUser(nt NotificationType, content string, db *dataStore) error {
	ch := connectToQueue(mqconn)
	n := notification_doc{
		Type:    nt,
		Content: content,
		Author:  u.Username,
	}
	err := ch.postNotification(u, n.Content)
	_, err = db.storeNotification(u, n)
	return err
}

func (u *user_doc) notifyFriends(nt NotificationType, content string, db *dataStore) error {
	for _, f := range u.Friends {
		log.Println("notifying " + f + " of this")
		toUser, err := read_user(f, db)
		if err != nil {
			return err
		}
		toUser.notifyThisUser(nt, content, db)
	}

	return nil
}
