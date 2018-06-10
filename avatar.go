package main

import (
	"errors"
	"io/ioutil"
	"path/filepath"
)

// ErrNoAvatarURL is the error that is returned when the
// Avatar instance is unable to provide an avatar URL.
var ErrNoAvatarURL = errors.New("chat: Unable to get an avatar URL")

// Avatar represents types capable of representing
// user profile pictures.
type Avatar interface {
	// GetAvatarURL gets the avatar URL for the specified client,
	// or returns an error if something goes wrong.
	// ErrNoAvatarURL is returned if the object is unable to get
	// a URL for the specified client.
	GetAvatarURL(ChatUser) (string, error)
}

// TryAvatars represents type Available avatar
type TryAvatars []Avatar

// GetAvatarURL selects avatar type available
func (a TryAvatars) GetAvatarURL(u ChatUser) (string, error) {
	for _, avatar := range a {
		if url, err := avatar.GetAvatarURL(u); err == nil {
			return url, nil
		}
	}
	return "", ErrNoAvatarURL
}

// AuthAvatar represents avatar data obtained at the time of authentication
type AuthAvatar struct{}

// UseAuthAvatar represents using UseAuthAvatar for avatar
var UseAuthAvatar AuthAvatar

// GetAvatarURL returns AuthAvatar's avatar image URL
func (AuthAvatar) GetAvatarURL(u ChatUser) (string, error) {
	url := u.AvatarURL()
	if url != "" {
		return url, nil
	}
	return "", ErrNoAvatarURL
}

// GravatarAvatar represents avatar data it is registed at Gravatar
type GravatarAvatar struct{}

// UseGravatar represents using UseGravatar for avatar
var UseGravatar GravatarAvatar

// GetAvatarURL returns GravatarAvatar's avatar image URL
func (GravatarAvatar) GetAvatarURL(u ChatUser) (string, error) {
	return "//www.gravatar.com/avatar/" + u.UniqueID(), nil
}

// FileSystemAvatar represents uploaded avatar data
type FileSystemAvatar struct{}

// UseFileSystemAvatar represents using FileSystemAvatar for avatar
var UseFileSystemAvatar FileSystemAvatar

// GetAvatarURL returns FileSystemAvatar'a avatar image URL
func (FileSystemAvatar) GetAvatarURL(u ChatUser) (string, error) {
	if files, err := ioutil.ReadDir("avatars"); err == nil {
		for _, file := range files {
			if file.IsDir() {
				continue
			}
			if match, _ := filepath.Match(u.UniqueID()+"*", file.Name()); match {
				return "/avatars/" + file.Name(), nil
			}
		}
	}
	return "", ErrNoAvatarURL
}
