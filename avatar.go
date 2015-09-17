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
	GetAvatarURL(ChatUser) (string, error)
}

//AuthAvatar ...
type AuthAvatar struct{}

//UseAuthAvatar ...
var UseAuthAvatar AuthAvatar

//GetAvatarURL ...
func (_ AuthAvatar) GetAvatarURL(c ChatUser) (string, error) {
	url := c.AvatarURL()
	if len(url) > 0 {
		return url, nil
	}
	return "", ErrNoAvatarURL
}

//GravatarAvatar ...
type GravatarAvatar struct{}

//UseGravatar ...
var UseGravatar GravatarAvatar

//GetAvatarURL ...
func (_ GravatarAvatar) GetAvatarURL(c ChatUser) (string, error) {
	return "//www.gravatar.com/avatar/" + c.UniqueID(), nil
}

//FileSystemAvatar ...
type FileSystemAvatar struct{}

//UseFileSystem ...
var UseFileSystem FileSystemAvatar

//GetAvatarURL ...
func (_ FileSystemAvatar) GetAvatarURL(c ChatUser) (string, error) {
	if files, err := ioutil.ReadDir("avatars"); err == nil {
		for _, file := range files {
			if file.IsDir() {
				continue
			}
			if match, _ := path.Match(c.UniqueID()+"*", file.Name()); match {
				return "/avatars/" + file.Name(), nil
			}
		}
	}
	return "", ErrNoAvatarURL
}
