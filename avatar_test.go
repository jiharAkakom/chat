package main

import (
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func TestGetAvatarURL(t *testing.T) {
	var authAvatar AuthAvatar
	client := new(client)
	url, err := authAvatar.GetAvatarURL(client)
	if err != ErrNoAvatarURL {
		t.Error("AuthAvatar.GetAvatarURL should return ErrNoAvatarURL when no value is present")
	}

	//set a Value
	testURL := "http://url-to-gravatar"
	client.userData = map[string]interface{}{"avatar_url": testURL}
	url, err = authAvatar.GetAvatarURL(client)
	if err != nil {
		t.Error("AuthAvatar.GetAvatarURL should return without an error")
	}
	if url != testURL {
		t.Error("AuthAvatar.GetAvatarURL should return correct URL")
	}
}

func TestGravatarAvatar(t *testing.T) {
	var gravAvatar GravatarAvatar
	client := new(client)
	client.userData = map[string]interface{}{"userID": "0bc83cb571cd1c50ba6f3e8a78ef1346"}
	url, err := gravAvatar.GetAvatarURL(client)
	if err != nil {
		t.Error("GravatarAvatar should not return an error")
	}
	if url != "//www.gravatar.com/avatar/0bc83cb571cd1c50ba6f3e8a78ef1346" {
		t.Errorf("GravatarAvatar.GetAvatarURL wrongly returned %s", url)
	}
}

func TestFileSystemAvatar(t *testing.T) {
	//make a test avatar file
	filename := path.Join("avatars", "abc.jpg")
	ioutil.WriteFile(filename, []byte{}, 0777)
	defer func() { os.Remove(filename) }()

	var fileSystemAvatar FileSystemAvatar
	client := new(client)
	client.userData = map[string]interface{}{"userID": "abc"}
	url, err := fileSystemAvatar.GetAvatarURL(client)
	if err != nil {
		t.Error("filesystemAvatar should no return an error")
	}
	if url != "/avatars/abc.jpg" {
		t.Errorf("FileSyaytemAvatar.GetAvatarURL wrongly returned %s", url)
	}
}
