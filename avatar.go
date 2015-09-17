package main

import (
	"errors"
	"io/ioutil"
	"path"
)

//ErrNoAvatarURL is the error that is returned when the Avatar
//instance is unable to provide an avatar url
var ErrNoAvatarURL = errors.New("chat: Unable to get an avatar URL")

//Avatar represents types capabale of providing a user profile picture
type Avatar interface {
	//GetAvatarURL returns ErrNoAvatar if no avatar url is available
	GetAvatarURL(c *client) (string, error)
}

//AuthAvatar ...
type AuthAvatar struct{}

//UseAuthAvatar ...
var UseAuthAvatar AuthAvatar

//GetAvatarURL ...
func (_ AuthAvatar) GetAvatarURL(c *client) (string, error) {
	if url, ok := c.userData["avatar_url"]; ok {
		if urlStr, ok := url.(string); ok {
			return urlStr, nil
		}
	}
	return "", ErrNoAvatarURL
}

//GravatarAvatar ...
type GravatarAvatar struct{}

//UseGravatar ...
var UseGravatar GravatarAvatar

//GetAvatarURL ...
func (_ GravatarAvatar) GetAvatarURL(c *client) (string, error) {
	if userid, ok := c.userData["userID"]; ok {
		if useridStr, ok := userid.(string); ok {
			return "//www.gravatar.com/avatar/" + useridStr, nil
		}
	}
	return "", ErrNoAvatarURL
}

//FileSystemAvatar ...
type FileSystemAvatar struct{}

//UseFileSystem ...
var UseFileSystem FileSystemAvatar

//GetAvatarURL ...
func (_ FileSystemAvatar) GetAvatarURL(c *client) (string, error) {
	if userid, ok := c.userData["userID"]; ok {
		if useridStr, ok := userid.(string); ok {
			if files, err := ioutil.ReadDir("avatars"); err == nil {
				for _, file := range files {
					if file.IsDir() {
						continue
					}
					if match, _ := path.Match(useridStr+"*", file.Name()); match {
						return "/avatars/" + file.Name(), nil
					}
				}
			}
		}
	}
	return "", ErrNoAvatarURL
}
