// Code generated by go-bindata.
// sources:
// assets/css/custom.css
// assets/tmpl/config.tmpl
// assets/tmpl/header.tmpl
// assets/tmpl/index.tmpl
// assets/tmpl/log.tmpl
// assets/tmpl/navbar.tmpl
// assets/tmpl/status.tmpl
// assets/tmpl/twitter_setup.tmpl
// DO NOT EDIT!

package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

func (fi bindataFileInfo) Name() string {
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _assetsCssCustomCss = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x8c\x92\xdd\x8a\xe3\x30\x0c\x85\xef\xf3\x14\x86\xbd\x76\x28\x64\x61\x97\x14\xfa\x2e\x8a\xad\xb8\xde\xb8\x92\x91\xd5\xbf\x1d\xe6\xdd\x87\xf4\x27\xb8\x74\x08\x93\x8b\xe0\x8b\xf3\x1d\x1d\x89\x33\xb0\xbf\x9a\x8f\xc6\x3c\xbe\x0c\xde\x47\x0a\x56\x39\xf7\xe6\xcf\x26\x5f\xb6\xcd\x67\xd3\xb4\xff\x8e\x87\x81\x55\x98\x2a\xe9\x00\x6e\x0a\xc2\x47\xf2\xd6\x71\x62\xe9\x8d\x0a\x50\xc9\x20\x48\xba\x5d\x64\x8a\x17\xb5\x90\x62\xa0\xde\x38\x24\x45\xb9\x5b\x3a\x26\x85\x48\x28\x95\xe5\x01\x2e\xf6\x1c\xbd\xee\x7b\xf3\x77\xf3\x1c\xee\xe3\xa9\x9a\x0f\xa6\x64\xa0\x76\x84\x0a\x1b\x99\xd4\x96\xf8\x1f\x7b\xf3\xfb\x49\x65\xc1\x3a\x2b\x8b\x47\xe9\x0d\x31\xe1\xf6\xa7\x1b\xcc\x29\x5f\xe6\xbc\xeb\xcf\xfb\xa8\x58\x0e\x3c\xe1\x53\x6e\xf5\x1c\x55\x51\xe6\x27\x38\x8d\xa7\x3a\xc6\x83\xfa\xd5\x75\x5d\x37\x8e\x0b\x12\x98\x43\xc2\x55\x62\x1c\x67\x66\x21\x32\x94\xf2\xad\x30\x08\x5e\xab\x28\x88\xea\xd1\x4d\xef\xba\x21\x81\x9b\xd6\x2e\xe1\xd9\x07\x94\x21\x1d\xef\x9b\x29\x0c\x09\xcd\xce\xe8\x1e\xc1\xaf\x9e\x04\x52\x74\xf8\xca\xb5\xb7\xbf\x2d\x2a\x31\xa3\x9f\x5d\x6e\xad\xdb\x19\x95\xd6\x63\x42\xc5\xda\x92\x33\xb8\xa8\xd7\xde\x6c\xda\x6e\x2d\xe1\xb2\xea\x5c\x10\x85\xc1\xce\x8d\x42\xd2\x97\x3e\x49\x88\x74\x2f\x73\xf7\x68\xc6\x57\x00\x00\x00\xff\xff\xcd\xb5\xec\x4e\xf1\x02\x00\x00")

func assetsCssCustomCssBytes() ([]byte, error) {
	return bindataRead(
		_assetsCssCustomCss,
		"assets/css/custom.css",
	)
}

func assetsCssCustomCss() (*asset, error) {
	bytes, err := assetsCssCustomCssBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "assets/css/custom.css", size: 753, mode: os.FileMode(436), modTime: time.Unix(1484921411, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _assetsTmplConfigTmpl = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\xec\x9d\x6d\x6f\x1b\xb9\x11\xc7\xdf\xfb\x53\xb0\xdb\xc3\xf9\x0e\x17\x4b\xbd\xbb\xa2\x2f\x2e\x92\x80\x34\xb9\x20\x29\x9c\x5c\x6a\x3b\x05\x8a\xa2\x30\x28\xed\x48\x62\xc2\x25\x55\x2e\xd7\xb6\x20\xf8\xbb\x17\xdc\x27\xed\x03\xc9\x5d\xc9\x7a\x58\x39\xf4\x8b\x8b\xad\xe5\x70\x87\xf3\x27\x67\x87\x5c\xfd\x70\xab\x95\x0f\x53\xc2\x00\x79\x13\xce\xa6\x64\xe6\x3d\x3e\x9e\x0d\xfe\xf4\xe6\x8f\xd7\x37\xff\xfe\xf4\x3b\x9a\xcb\x80\x8e\xce\x06\xea\x1f\x44\x31\x9b\x0d\x3d\x60\xde\xe8\x6c\xb5\x92\x10\x2c\x28\x96\x80\xbc\x39\x60\x1f\x84\x87\x7a\xca\x70\xcc\xfd\xe5\xe8\x0c\x21\x84\x06\xe1\x44\x90\x85\x4c\xfe\x50\x3f\xd3\x88\x4d\x24\xe1\x0c\xf9\x40\x41\xc2\x15\xbf\xff\x61\x1c\x49\xc9\xd9\x0b\x14\x70\x1f\xd3\x1f\xd1\x2a\x6f\xab\x7e\xee\xb0\x40\x52\xa0\x21\x4a\x5a\xf5\x16\x58\x00\x93\x1f\xb9\x0f\x85\x5f\x5f\x96\x4c\xc8\x14\xfd\x20\x45\x6f\x42\x71\x18\x7e\xc4\x01\xa0\xe1\x70\x88\xce\x93\xfb\xf9\xe7\xd5\x1b\xc4\x4e\x71\x81\x7e\x50\x77\x22\x68\x88\xfe\xf2\x12\x11\x34\x40\xaa\x87\x39\xa1\xbe\xea\x3f\xec\x51\x60\x33\x39\x7f\x89\xc8\x4f\x3f\xe9\x3a\xc8\x3c\x8d\x2d\xd0\xb0\x6c\xfc\x1f\xf2\x5f\xad\x81\xf2\x33\x6e\xd5\x93\xcb\x45\xea\xe5\x9c\xf8\x3e\xb0\x73\xf4\xfd\xf7\x49\x5f\x3d\xc6\x7d\x58\x0f\xe2\xfd\xc7\x4f\x9f\x6f\xb4\x43\xc8\x7e\x12\xab\x3b\x4c\x23\x40\x43\x74\x3e\xc5\x34\x84\xf3\x97\xc6\xe6\xe5\x38\xa1\x73\x4b\xd3\x54\x00\xc2\x18\x88\x77\x37\x1f\x2e\x55\xf3\x37\x71\x50\x0d\x46\x8f\xb5\x4f\xcb\x9f\x3c\x22\xa0\x21\x7c\x53\x72\x48\x11\x6d\xa2\x46\x36\x67\x37\x13\xe5\x0a\xee\x40\xc8\x6d\x45\x39\x2b\xff\x36\xe8\x17\x57\xf0\xc0\x27\x77\x28\x76\x71\xa8\x12\x85\xc4\x84\x81\xf0\xd6\xab\xbb\x98\x11\x18\xbe\x1b\xe3\x34\x23\xe4\x0d\x8a\x1d\x7c\x89\x82\x31\x97\x82\xb3\x42\x07\x71\x9b\xf9\xcf\xa3\xd7\x71\x16\x8a\x04\x56\xb9\x62\xd0\x9f\xff\x5c\x6d\xf2\x4b\xd6\x0b\x05\xec\x7b\x79\x7b\x40\x4b\x1e\x09\xc4\xef\x19\x1a\x73\x79\x1e\xa2\x90\x47\x62\x02\xe1\x0b\x34\x25\x54\x82\x08\x11\x66\x3e\xc2\x71\x0a\x0a\x07\xfd\xf9\x2f\x95\x7e\xa7\x5c\x04\x88\xf8\x43\x8f\xdf\x81\xb8\x17\x44\x82\x97\xb6\x1e\x7a\xfd\x24\x35\xf6\x3d\x14\x80\x9c\x73\x7f\xe8\x2d\x78\x28\x3d\x04\x6c\xa2\x66\xcc\xd0\x0b\x22\x2a\xc9\x02\x0b\xd9\x57\xdd\x5c\xf8\x58\x62\x6f\x34\x88\xff\xa8\xdc\x86\xb0\x45\x24\xd5\x3c\x0f\x4a\x77\x4a\xba\x09\xa3\x71\x40\xa4\x87\xe2\x59\x33\xf4\xae\xf1\x1d\x78\xd9\x68\xc7\x92\xa1\xb1\x64\x17\x74\x16\xff\xb3\x10\x24\xc0\x62\xe9\xa1\x50\x2e\x29\x0c\x3d\x9f\x84\x0b\x8a\x97\xbf\x11\x46\x09\x03\x6f\x74\xa6\x19\x5d\x75\x38\x53\x42\xa1\x30\xa6\x19\x48\x73\x77\xd5\xb9\x93\x0e\x44\xeb\xf7\x1b\x7e\xcf\x28\xc7\x7e\x93\xef\x95\xd0\xa4\xe1\xda\xdc\xf1\x16\x62\xec\x64\x5c\x9f\x17\x9b\x8f\xaa\xd6\xa7\xf2\xdd\x43\x0c\x07\xca\xd3\xe5\x98\xcb\x5e\xfa\xdc\x6d\xe5\x62\x75\x4e\x0d\xfa\x3e\xb9\x2b\xc4\x6c\xb5\x22\x53\xd4\xfb\x00\x61\x88\x67\xf0\xf8\xa8\x5d\x7d\x98\x82\x90\x28\xfe\xef\x85\x8f\xd9\x4c\xad\xe3\xd5\x6a\x6d\x94\xf6\xb9\xee\x12\x98\xaf\x56\xf2\xba\xb3\xf9\xaf\xa3\x9b\x7b\x22\x25\x08\x74\x43\x02\x50\x7e\xaa\x25\xf5\xeb\x48\x7b\x3f\x89\xc7\x14\x2e\x04\x84\x0b\xce\x42\x72\x57\x1b\x52\x7c\xbd\xd4\x18\x25\x26\xa1\x14\x64\x01\x7e\xfa\xd7\x98\x0b\x1f\x04\xf8\xba\x00\x4b\x55\x84\xd4\x3f\x4f\xae\x09\xfd\x85\xe4\xa2\x3f\x0a\x27\x02\x80\xc5\x8a\x84\x83\xbe\x34\x74\x93\xb5\x86\x87\x09\x8d\x7c\x40\x02\x16\x94\xb4\x31\x20\x2c\x35\x90\x2d\x1a\x4f\x78\xc4\x64\x73\x33\x2e\xe7\x20\x5a\x74\x67\x6e\x31\xe8\xeb\xc2\x32\xe8\x1b\x02\x39\x90\xeb\xc2\xae\xfa\xb3\x5a\x09\x35\x89\xd0\x77\x84\xf9\xf0\xf0\x02\x7d\x27\xd3\x29\x81\x7e\x1b\xa2\x5e\x92\x9e\x7b\xe9\x74\xe9\xe5\xd3\xe5\xb1\xfe\x38\x42\x8d\x6a\x59\x93\x67\xf2\xc4\xce\x96\x96\x4c\x6f\x98\x39\x13\xf6\xd2\x47\x6a\xbe\x9a\xe3\x02\xc9\xeb\xdb\x23\x68\xbc\x98\x8c\x9c\x92\x50\xde\xc0\x83\x1c\xf3\x87\xf5\xb8\x7b\xd7\xf1\x8c\x52\x8f\xf2\x10\x69\x1c\x49\x26\xdc\x6d\x3c\xe1\x3c\x43\x1c\x50\xa2\x46\x83\xbe\x0d\xde\x8d\x39\xa7\xd7\x40\x61\x52\xf1\xef\xf7\x64\x0e\x5f\x25\x53\x58\xe7\x62\x3a\xcb\x6f\xd3\x59\x7e\x14\x2f\xdf\x27\x0b\xe7\x4a\x6a\x3d\x4c\x97\xd5\xad\x90\xfb\xf5\x4e\x26\xea\xfe\x31\x7d\xcf\xe4\x27\x29\x0a\xfe\xbd\x56\x6b\x55\xe7\x5a\xbc\x88\xf7\xe8\xd4\x20\x29\xfd\xd2\x69\x9f\xfc\x51\x7b\x28\x11\x36\xe5\x1e\x52\xcf\xbe\x0b\xc9\x67\x33\xf5\x60\x89\x77\x59\xd9\x67\x58\xcc\x40\x0e\xbd\x3f\xa7\xee\x5f\xe4\xee\x5f\x24\xe5\xd2\xed\x6a\x95\x2c\xe8\xc7\x47\x4d\xc2\xad\xfe\xf8\x20\x31\xa1\x21\x9a\x83\x00\xbb\xef\xfd\xc4\x5f\xcb\xf0\x0f\x10\x9c\xf4\xa9\x87\x38\x9b\x50\x32\xf9\x3a\xf4\xd6\x5b\x52\x39\x27\xe1\x0b\x74\xde\x22\x2c\xe7\x3f\x7a\xa3\x64\x17\xf4\x94\x41\xe9\x13\x31\x2a\x3c\x77\xf5\x61\x68\x78\xa8\xd9\x63\xa4\x2d\xaa\xf2\xa1\xf6\xb1\xef\x57\x8a\xab\xe6\x19\x60\xab\x9b\x3e\xc2\x7d\x4d\x82\xac\x5a\xb2\xe4\x5f\xa4\xa9\x77\xea\xd7\xb7\x7f\x06\x76\xa5\x85\xf1\x59\x5c\x7f\xe6\x0e\xfa\x71\x2d\x64\xab\xff\x76\xf0\x30\x2e\x16\x6f\x71\xd2\x40\x53\xec\x83\x17\xef\x8e\xda\xe4\x0b\x55\xb0\xc5\xbf\x0f\xbd\x8b\x9f\x3d\x24\x78\x52\xd6\x62\xca\x67\x1e\xc2\x82\xe0\x0b\x8a\xc7\x40\x29\xf8\xe3\x65\xab\x1e\x2f\x24\x91\x54\x5b\xad\x57\x3d\xbd\xc8\x6e\x93\xde\x94\x4f\xa2\x00\x98\x69\x02\xd7\xcd\xd5\xd6\xd6\xdc\x5e\x6f\x93\x1e\x80\x3d\x21\x2f\x4d\x28\x0f\x21\xcd\xcc\x3e\x09\x03\x92\x77\xee\x8d\xbe\x57\x71\x09\x5f\x36\xe7\x18\x94\x54\xe6\x65\xdf\x92\xc0\xb5\x55\x2e\x8b\xf3\x9b\x24\x9d\x97\x4b\xfa\xda\xbd\xca\xdb\x84\xe6\x38\xa9\xe9\xdc\x14\xa5\xf9\xaf\xa3\xb7\xb1\x57\xf6\x7b\x57\x6f\xd0\xb0\xc5\xd0\x9a\x3f\x7d\xdb\xa1\xef\xd6\xbc\x15\xd1\xb7\xb7\x64\x72\xbd\x81\x3f\x52\xa5\xa3\x3d\xe3\x98\x2c\xe3\x8c\xbc\x99\xa9\xf9\x09\xa5\x69\xd9\x7e\xf0\xb6\x2d\x85\xbe\xfd\x16\x81\x5a\x60\x35\xe5\x59\xc3\x56\xc9\x64\xbd\x91\x01\xb2\x6c\x08\x92\x29\xdd\xfb\x94\x7a\xa3\x2b\x1a\x93\xb5\xd8\xcb\x1c\xb6\x95\x8f\x5a\x77\xf7\x24\x29\xda\x36\xf2\x91\xa0\xa8\x5b\xd1\xff\x7c\x75\xd9\x42\x80\x48\xd0\xdb\x67\x23\xc2\x1c\x87\x28\x00\x9f\xe0\x83\x29\x60\xda\xcc\xa5\x1a\xbc\xc3\xe1\x07\xe5\x8f\x45\x80\x39\x0e\x6f\x63\x9f\x9f\x45\xf4\x23\x41\x3b\x14\xfb\xcf\x57\x97\x0d\x91\x8f\x04\x3d\xf9\xb8\x0b\x90\xf7\x00\x12\xfc\xae\x44\xfe\x2a\x73\xc8\x12\xfc\xdc\xe9\x93\x0f\xff\x14\xdf\x71\x41\x24\x20\x39\x17\x10\xce\x39\x3d\x9c\x0e\xe6\xc3\x9a\x54\x89\xb7\xa9\x6f\x37\x99\x6b\x16\x45\xb2\x71\xdc\xe6\xe3\x38\x79\x69\xf2\x49\xd6\x45\x6d\xf2\x55\xd2\x46\x9c\x7c\x24\xcf\x48\x1d\x8a\xd9\xec\x40\x72\xd8\x4f\xd2\x95\x54\xe6\x73\xf4\x54\x01\xe5\x6d\x7e\xc2\xb3\x5a\xd5\xe4\xbc\xc4\x6c\xd6\xea\xec\xb0\xe4\xd7\xfe\x76\x26\xed\x76\x1b\xb5\x03\x16\x7d\x23\xeb\xce\x17\x95\x36\xb2\xbd\x7f\x91\x30\x79\x93\xed\xf6\xb3\x7a\x03\xb7\x9f\x6d\x69\xa0\x52\xc4\x18\x0e\x57\x50\x5a\xb7\x53\xc9\xb4\xee\x5d\x2a\x8f\x2c\x89\xfa\x2e\x69\x16\x3b\x7e\xf2\x19\x7a\x8a\x27\x80\xe2\x77\x07\x88\x92\xaf\x40\xc9\x9c\xf3\x43\x3d\x41\x77\x90\xb2\x53\x2d\xd4\x28\x7a\xf1\x28\x6e\xd7\xa3\xb0\xa5\xf2\x54\xea\xb7\xca\xee\x95\xb2\xbb\xcc\xcd\x3a\x93\xe1\xd1\x93\x34\x1d\xd3\x48\x08\xf0\x4f\x5e\xd5\x74\x1c\x5b\xe8\xfa\xf7\xc4\xf2\xd9\x29\xab\x12\xf8\x3d\xe0\xd3\x5f\xb0\xd9\x40\xb6\xd0\xf6\x5d\x6a\xfa\xec\xc4\xfd\xc2\x97\x27\xaf\xeb\x17\xbe\xdc\x42\xd2\x7f\xf0\xe5\x73\x52\x53\x45\xad\x53\x65\x8d\xba\xde\x5c\xd5\xc4\x62\x9f\x7a\x51\x43\x31\xf3\x03\x2c\xbe\x76\x2a\xfe\x97\xa9\x53\x6d\x2a\xcb\xa4\xe5\xe9\xeb\xc0\x67\xbc\x5b\x1a\xf0\x19\x6f\x11\x7f\x3e\xe3\x5d\x89\xfd\x31\xb7\xf9\x97\x98\xcd\x22\x3c\x03\xb7\xd1\x37\x1a\xb8\x8d\x7e\x4b\x03\x7f\x14\x10\x86\x42\x60\x6a\xd1\x35\x7d\x21\xd8\xd4\xc5\x53\xce\x68\xdf\x52\x8e\xe5\xdf\xfe\xaa\x3d\xa7\xcd\x26\x7a\xef\x03\x61\xd7\x99\x8b\x96\x3c\x41\xb3\xf6\x01\x61\xb7\xf9\x98\xba\x92\x32\xd0\xd6\x0a\xe1\x87\xee\x2b\x84\x1f\x36\x53\x08\x3f\x74\x4f\xa1\xa3\x24\xf5\x57\x13\xe9\x0e\x6d\x6d\x06\x2e\x97\xb7\x34\xc8\xdf\xba\x1d\x2c\x47\x4c\xe6\x30\xf9\x5a\x2e\xec\x92\xe9\x9c\xbd\x62\xd3\xa5\x82\xe4\xdb\xc0\xd9\x7b\xb5\xae\xac\x7d\xf4\xd4\x97\xd0\x1d\x08\x7a\xf6\xce\xd9\x12\xf5\xcc\xdb\xd3\x0f\x3b\xa7\x94\xdf\x77\x21\xe8\xb1\x23\xb6\x90\xc7\x0d\x4e\x3e\xe0\x13\x4e\x29\x64\x4c\xeb\x51\xf7\x8e\x69\xe0\x5f\xaf\x1d\xb2\x44\xbf\xe0\x76\x57\x24\x38\x68\x99\xb1\xf1\x77\xa7\xa7\x9c\xcb\xa7\x7d\xc7\x3c\x67\x5f\x60\x8a\x23\x2a\x0d\xdf\x36\x7f\x4d\x79\xd8\x8a\x67\x31\xfa\x6f\xb8\xa4\xf9\xd8\xc4\x99\xe6\xd7\x0b\x98\x69\x96\x44\xbf\x15\xcc\xf4\x59\x80\xa0\xf9\x37\xc1\x34\xec\x49\x2e\xe8\xa1\x40\xd0\xcc\x99\x03\x80\xa0\xd9\xad\xf4\x20\xe8\xda\x91\x83\x81\xa0\xb5\xef\x5e\xe5\x1e\x56\x20\xc6\xb5\x6f\x27\x09\x31\xe6\xee\x3b\x88\xb1\x08\x31\x5a\xc2\xf2\xdc\x20\xc6\x7c\xa8\x0e\x62\x6c\xdf\x62\xaf\x00\xe2\x36\x0f\x81\x36\x00\xa2\x6d\xad\x6f\x07\x20\x5a\x7a\x74\x00\xe2\x0e\x01\xc4\x16\x71\x76\x00\xa2\x3b\xfb\x5b\x9b\x9d\xf0\xd9\xdf\x71\x11\xb8\xbc\xcc\x33\x02\x88\xeb\x82\xcf\x01\x88\xd5\x1e\x76\x1d\x7d\x2d\x80\x58\x13\xc0\x01\x88\x6b\xf3\x2d\x14\xa8\xa0\x58\x55\x0d\xea\x00\x62\x4d\x00\x07\x20\xee\x2f\xf6\x25\x00\x51\x1b\x79\x07\x20\xee\x21\xf2\x1a\x00\xb1\x16\x7c\x07\x20\x56\xfb\xd9\x42\x07\xf3\x41\x4b\x33\x80\x58\x53\xc4\x01\x88\xe6\x8e\xf6\xa1\x8d\x0d\x40\x34\x2f\x97\x67\xa4\x4e\xa7\x01\xc4\x9a\x02\x55\x00\xb1\x2a\xa7\x03\x10\x1d\x80\xd8\xce\xc0\xed\x67\x5b\x1a\x1c\x15\x40\xac\xae\x6f\x3d\x80\x58\x4b\x13\x0e\x40\xd4\x77\x75\x8c\x94\xdd\x16\x40\x34\x48\xed\x00\xc4\xc6\xce\x8e\xad\xaa\x1d\x40\xb4\xe9\xea\x00\xc4\x86\xde\x8e\x2d\x6d\x03\x80\x68\xd3\xd6\x01\x88\x96\x8e\x8e\xad\xab\x19\x40\xb4\x49\xea\x00\xc4\xfd\x96\x35\x65\x00\xd1\x24\xa0\x03\x10\xf7\x57\x56\x56\x01\x44\x73\x65\xe9\x00\xc4\x4d\x0c\x50\x7b\x0d\x4a\x00\xa2\x31\xfe\x0e\x40\x74\x00\x62\x3b\x03\xb7\xd1\x6f\x69\xd0\x29\x00\x51\x77\xb0\x67\x01\x10\xb5\x07\x85\x0e\x40\xd4\x75\xb1\x7f\x85\xb4\x00\xa2\x45\x21\x07\x20\x22\x07\x20\xb6\x32\x70\xb9\xbc\xa5\xc1\x31\x01\xc4\x3c\x31\x98\x00\xc4\x75\x2a\x70\x00\x62\xc1\x7a\xa7\x41\xaf\x03\x88\xb5\xa8\x3b\x00\x71\xe7\x41\xaf\x00\x88\xf5\x90\x3b\x00\x71\xb7\x7b\x47\x1b\x80\x58\x8b\xbe\x03\x10\x1d\x80\xb8\x11\x80\x78\x0d\x58\x4c\xe6\x27\xcc\x1f\xfe\x2f\x02\xd1\xea\x7f\x58\x29\x20\x8c\x68\x02\xc1\x7c\x23\x9c\x62\x18\x6b\xab\x03\x54\x32\xd5\x0f\x05\x29\x86\xe9\xfd\xf6\xcf\x28\x26\x77\xea\xfd\x33\x99\x15\xa8\xee\x41\x3a\x5f\xf6\xca\x25\x86\xeb\xaf\x2c\xa6\xfe\x5c\xc5\x93\xef\x66\xb9\x00\x8d\x4b\xc9\xcc\xbc\x55\x01\xf4\x90\xe7\x21\x4f\xc0\x44\xed\x13\x91\x17\x90\x07\x15\x2c\x6f\xc1\x17\x11\xc5\xe2\xb0\x30\x65\xea\x7a\x05\xa5\xcc\xbd\x3e\x49\x92\x32\xf3\xde\x81\x94\x45\x90\xd2\x1c\x95\xe7\xc6\x51\x66\x23\x75\x18\xe5\x2e\x5b\xec\x15\xb4\xdc\xf8\x29\xd6\x86\xb2\xb4\xe4\x81\xed\x20\x4b\x73\x87\x8e\xb1\xdc\x21\x63\xd9\x1c\x66\x87\x58\xba\xd3\xcd\xb5\xd9\x09\x9f\x6e\x1e\x17\xf2\x4b\x8b\x3f\x23\x60\x99\x97\x81\x8e\xaf\xac\xf6\xb0\xdb\xd0\x6b\xe9\xca\x6a\xf4\x1d\x5c\xb9\x36\xdf\x22\xfc\x15\xcc\xac\x2c\x40\x1d\xad\xac\x46\xdf\x91\x95\xfb\x0a\x7c\x89\xab\xd4\x85\xdd\x61\x95\x3b\x0f\xbb\x06\xaa\xac\x46\xde\x31\x95\xd5\x7e\xb6\x10\xc1\x74\xde\xd2\x4c\x54\x56\xe5\x70\x40\xa5\xb9\xa3\xdd\x0b\x63\xc3\x29\x8d\x0b\xe5\x19\x49\xd3\x69\x9a\xb2\x2a\x40\x15\xa6\x2c\x4b\xe9\x50\x4a\x87\x52\xb6\x33\x70\xfb\xd6\x96\x06\x47\x45\x29\xcb\xab\x5b\x0f\x52\x56\x33\x84\xe3\x28\xf5\x5d\x1d\x21\x59\xb7\xc5\x28\xb5\x32\x3b\x88\xb2\xb1\xb3\x23\x4b\x6a\x67\x28\xcd\xa2\x3a\x82\xb2\xa1\xb7\x23\xeb\xda\x00\x50\x9a\x85\x75\xf8\xa4\xa5\xa3\x23\x8b\x6a\xa6\x27\xcd\x7a\x3a\x76\x72\x9f\x75\x4c\x99\x9c\x34\x68\xe7\xc0\xc9\x7d\x15\x91\x55\x6c\xd2\x58\x47\x3a\x6a\x72\x13\x03\xd4\x56\x80\x12\x33\x69\x0a\xbe\x43\x26\x1d\x32\xd9\xce\xc0\x6d\xe8\x5b\x1a\x74\x0a\x99\xac\x1f\xdf\x59\x80\x49\xdd\x59\xa0\xe3\x25\x75\x5d\xec\x5b\x1e\x2d\x2d\x69\x96\xc7\xc1\x92\xc8\xc1\x92\xad\x0c\x5c\x16\x6f\x69\x70\x4c\x58\x32\xcd\x0a\x26\x54\x32\xcf\x03\x8e\x94\x2c\x58\xef\x30\xe2\x75\x4e\xb2\x1a\x72\x87\x49\xee\x38\xe2\x15\x48\xb2\x16\x6f\xc7\x48\xee\x72\xa7\x68\x23\x24\xab\xa1\x77\x80\xa4\x03\x24\x37\x02\x24\x3f\x51\x3c\x01\xf4\x91\x4b\x32\x25\x13\x5c\xaf\xc7\x4e\x08\x95\xc4\x71\x56\x0a\x81\x4e\x9b\x61\x8f\x28\xb4\x82\x8d\x3b\xc4\x16\x9f\x86\xf9\x14\x92\x70\x15\x0f\x29\x6a\xd6\x8b\x65\xec\xbd\x52\x11\xb8\x06\x3a\x5d\x27\x08\x56\x6c\xb5\x88\x5b\xc5\x71\xba\x55\x71\xda\x2b\x70\x57\x4c\x65\x2d\x7c\xff\xac\x14\xb1\xfa\x1d\x6b\xb6\x9d\xcb\xbb\x24\x78\xd4\xea\x79\xcf\x24\x08\xac\xd9\xbe\xb8\xe5\x52\xfa\xb4\x03\xcb\xa5\x20\x55\x71\x7d\x90\xc2\xc7\x47\x5c\x10\x45\xef\xd2\x15\x50\xf4\xac\x43\x53\xfe\x92\xcf\xdc\x54\x5f\x37\xa3\x84\x01\x8b\x82\xd3\x5c\x13\x97\x7c\x56\x5c\x0b\x94\xcf\x8e\xb9\x06\x94\x37\xe9\xdc\x57\x9e\x3c\x61\xce\xb7\x72\xc2\xfe\x22\x99\x45\xc1\x18\x44\xf6\x2a\x59\xf9\x93\x2a\x5d\x78\x47\x5c\x74\xfc\x32\xb9\x6a\x7d\x1d\xbc\xd7\x45\x5a\xf8\x75\xd0\x4f\xcc\x06\xfd\xb9\x0c\xe8\xe8\x2c\x2b\xfc\xfe\x1f\x00\x00\xff\xff\xe2\x66\x25\x50\x2c\xbe\x00\x00")

func assetsTmplConfigTmplBytes() ([]byte, error) {
	return bindataRead(
		_assetsTmplConfigTmpl,
		"assets/tmpl/config.tmpl",
	)
}

func assetsTmplConfigTmpl() (*asset, error) {
	bytes, err := assetsTmplConfigTmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "assets/tmpl/config.tmpl", size: 48684, mode: os.FileMode(436), modTime: time.Unix(1484921411, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _assetsTmplHeaderTmpl = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x8c\x93\xdf\x4f\xdb\x3e\x14\xc5\xdf\xfb\x57\x58\x7e\xfd\x7e\x1b\xaf\x3f\x46\xbb\xa9\xa9\xc4\x60\x94\xc2\xf8\x31\x28\x20\xf6\xe6\x3a\x37\x89\x43\x6c\x07\xdf\x9b\xa6\x59\xd5\xff\x7d\x6a\xca\x28\x42\x43\xea\x5b\xce\xb5\x7c\x72\x3e\x3e\xba\xab\x55\x04\xb1\xb6\xc0\x78\x0a\x32\x02\xcf\xd7\xeb\xd6\x68\xf3\x39\x6e\x31\xc6\xd8\xc8\x00\x49\xa6\x52\xe9\x11\x28\xe4\x77\xb3\x93\xf6\x90\xbf\x3d\xb2\xd2\x40\xc8\x17\x1a\xaa\xc2\x79\xe2\x4c\x39\x4b\x60\x29\xe4\x95\x8e\x28\x0d\x23\x58\x68\x05\xed\x46\xfc\xcf\xb4\xd5\xa4\x65\xde\x46\x25\x73\x08\x3b\x7c\xdc\xda\x3a\x91\xa6\x1c\xc6\x17\xf5\xdc\xd1\x48\x6c\xc5\xcb\x49\xae\xed\x13\xf3\x90\x87\x1c\xa9\xce\x01\x53\x00\xe2\x2c\xf5\x10\x87\x3c\x25\x2a\xf0\xab\x10\x46\x2e\x55\x64\x83\xb9\x73\x84\xe4\x65\xb1\x11\xca\x19\xf1\x3a\x10\xbd\xa0\x17\x0c\x84\x42\xdc\xcd\x02\xa3\x6d\xa0\x10\x39\xd3\x96\x20\xf1\x9a\xea\x90\x63\x2a\x7b\xc3\x7e\xfb\xdb\xfd\xa3\xd6\xb7\xd3\x13\x38\xef\x44\x13\x73\x76\x73\xf8\x54\xab\xf2\xf4\xf0\xf4\x26\xe9\x75\xaf\xcc\x9d\xaa\xaa\x81\xb3\xbd\x9b\xc7\x28\xe9\xdf\xcb\xff\xae\xcd\xed\x0c\x7f\x8b\xf3\x83\xe1\x62\x1e\x7d\xcf\xd2\x7e\xc9\x99\xf2\x0e\xd1\x79\x9d\x68\x1b\x72\x69\x9d\xad\x8d\x2b\xf1\xef\xc3\xed\x03\xa5\x22\x9b\x61\xa0\x72\x57\x46\x71\x2e\x3d\x34\x44\x32\x93\x4b\x91\xeb\x39\x8a\xd8\x59\x6a\xcb\x0a\xd0\x19\x10\xfd\xe0\x20\xe8\x35\x78\x6f\xc7\xaf\x84\x7b\xfc\x55\x48\x44\x20\x6c\x3c\x54\x89\xe4\xcc\xcb\xcd\xed\x55\x54\x5e\x17\xc4\xd0\xab\x5d\xc0\x4d\x96\x20\x71\x2e\xc9\x41\x16\x1a\xdf\xe5\xcb\x9e\x4b\xf0\xb5\xe8\x04\x9d\x6e\xd0\x7f\x51\x4d\xa0\x0c\xf9\x78\x24\xb6\x86\xe3\x8f\xdd\xf7\xed\x34\x7b\x5f\x69\xf6\xcf\x46\x67\xea\xf3\xf4\xa7\x9e\x7f\xea\x0e\x9e\x17\x75\x76\x7b\x11\x9f\x66\x57\x17\xf2\xc7\x53\x5c\x3e\xdc\x2f\x7f\x2d\xef\xae\xed\xd1\xd9\xe1\x20\xef\x9a\xa3\x87\xcb\x69\x31\xf9\x62\x26\x47\xc7\xc3\x6a\x72\x39\x55\xd7\xc7\x83\xd9\x52\x7e\xdc\xe8\x8e\xe5\x63\x98\x22\x97\x14\x3b\x6f\x02\xaa\x34\x11\xf8\x86\xa4\xd2\x51\x02\x84\x4d\xe0\xd7\xfd\x2a\x29\xde\xec\xd7\xce\x74\x24\xb6\xbb\xb8\x5a\x81\x8d\xd6\xeb\xd6\x9f\x00\x00\x00\xff\xff\x0f\xae\x36\x4f\xaf\x03\x00\x00")

func assetsTmplHeaderTmplBytes() ([]byte, error) {
	return bindataRead(
		_assetsTmplHeaderTmpl,
		"assets/tmpl/header.tmpl",
	)
}

func assetsTmplHeaderTmpl() (*asset, error) {
	bytes, err := assetsTmplHeaderTmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "assets/tmpl/header.tmpl", size: 943, mode: os.FileMode(436), modTime: time.Unix(1484921411, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _assetsTmplIndexTmpl = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\xc4\x58\xcd\x8e\xdb\x36\x10\xbe\xef\x53\x4c\xd9\x1c\x5a\x20\xb2\x90\x04\xe8\x21\x90\x05\xa4\x49\xd0\x06\x48\xd2\x20\x4d\x0b\xf4\x38\x12\xc7\x12\xb3\x14\x47\x25\x69\x6f\x5d\xc1\xef\x5e\x50\xb2\xbd\x96\x2d\xc9\x6b\x27\x6d\x7d\xb0\x65\xcd\xdf\xc7\xe1\xcc\x37\x94\x9a\x46\xd2\x42\x19\x02\xa1\x8c\xa4\xbf\xc4\x66\x73\x93\x7c\xf3\xea\x97\x97\x9f\xfe\xf8\xf0\x1a\x4a\x5f\xe9\xf4\x26\x09\x3f\xa0\xd1\x14\x73\x41\x46\xa4\x37\x4d\xe3\xa9\xaa\x35\x7a\x02\x51\x12\x4a\xb2\x02\x66\xc1\x30\x63\xb9\x4e\x6f\x00\x00\x12\xa9\x56\x90\x6b\x74\x6e\x2e\x72\x36\x1e\x95\x21\x2b\x3a\x59\xf8\x1c\xba\x30\xb8\xca\x70\xeb\x62\xaf\x70\xe8\xe0\xf3\xb2\xca\xd8\x5b\x36\x07\x0e\x5a\x9d\xf2\x49\xfa\x6e\x9d\xb1\x4f\xe2\xf2\xc9\xb1\xe8\xe9\xce\x5a\x13\x4a\x91\xe2\xd2\x73\x85\x5e\xe5\xa8\xf5\x1a\x72\xd6\x9a\x72\x0f\x68\x24\x64\x96\x51\xe6\xe8\xc2\xbf\x35\xdc\x2a\x23\x1d\xf0\x02\x94\x59\xb0\x0d\x16\x6c\x60\xc1\x16\xd6\xbc\x4c\xe2\xf2\x69\x3f\x4c\xd3\xa8\x05\xcc\x7e\x64\xff\x1e\x2b\xda\x6c\xfa\x10\x10\x4a\x4b\x8b\xb9\x28\xbd\xaf\xdd\xf3\x38\xf6\x77\xca\x7b\xb2\xb3\x9c\xab\xb8\x69\xee\xad\xc4\x0e\x69\xe6\x0d\x64\xde\x44\xba\x68\x7f\x1c\xe7\x0a\x75\xa4\x72\xee\x6e\x6f\xed\x45\x9a\xb8\x1a\xcd\xce\x68\x81\xb0\xc0\x9d\x2c\x5c\x62\xee\xd5\x8a\x44\x9a\xc4\x41\x2d\x4d\x62\x3c\xc6\x4c\xda\x9d\x05\x8b\x75\xed\x66\x07\x88\xbf\x32\xc8\x1a\x9d\x3b\x87\xd2\xc8\xf3\x19\x25\xf2\x92\xf2\xdb\x7f\x11\xe9\x3e\xc6\x38\xd6\x13\x5c\x39\x1b\xc7\x9a\x66\x92\x56\xa4\xb9\x26\xeb\x66\x05\x73\xa1\xa9\xdd\x7b\xac\x95\x0b\x5f\xf1\x4a\x39\xc5\x66\x2b\x0a\x77\x5b\xf1\x9f\x4b\xf6\xe8\x1e\xba\x8a\xce\x78\x70\x11\x9d\xe8\x21\x25\x71\xb2\x80\x42\xf9\x72\x99\xb5\x70\xd4\x1d\x7a\xbc\xc5\xb8\x0a\x8d\xf6\x60\x54\xad\xfd\x30\xaa\x56\x74\x0e\x55\x12\x4b\xb5\x4a\x0f\xe8\x60\xa9\x77\x6e\x0c\xae\xc0\xe0\x2a\xaa\x95\xd6\xae\xbd\xfa\xbc\x74\x5e\x2d\x14\xc9\x1e\x3f\x24\x5a\x81\x65\x4d\x73\x51\x5b\x72\x64\x7c\xdb\xcb\xfb\x15\xec\xa3\x23\x48\xf4\x18\x79\x2e\x8a\xa0\xec\x31\x13\xdb\x6c\x7c\x5b\x72\x45\x22\xfd\x99\x2b\x0a\xe0\x92\x58\xab\xc3\x00\x4d\x63\xd1\x14\x04\x8f\x6e\x69\xfd\x18\x1e\xad\x50\xc3\xf3\x39\xcc\x5e\x76\xdc\xa2\xd8\xbc\xc3\xba\x57\xc1\x23\x88\x26\x21\x34\x4d\x70\xbf\xd9\x88\x74\x77\x35\x0c\xa5\xdf\x2d\x49\xbc\xd4\xe9\x30\x99\x7a\xcc\xa2\xc0\xc8\x64\xfc\x31\x9d\x06\x2d\x25\xe7\xa2\x5d\xf6\xa1\x7e\x8d\x26\xd4\x91\x24\x50\x06\x76\x89\xeb\xd9\x1e\x47\x09\x16\x1a\xda\xef\x48\x86\x2c\xd9\x01\x83\x41\xa3\x28\x4c\x14\x65\x8a\x11\xfd\xd6\xa6\x7c\xd6\x37\xf1\xca\xeb\x6e\x2b\xbb\xac\xc5\x9a\x8b\x58\xa4\xaf\xad\x65\x0b\x6f\xb9\xe8\x52\x56\x3e\x1b\x81\xd0\xd5\xda\x90\xa8\xe3\xf7\xb7\x5c\x1c\x31\xd1\xf4\x0a\xc2\x18\x14\x69\x52\x5b\x4a\x93\x9c\x25\xa5\x4d\xd3\xf9\x48\xe2\xf6\x6f\x12\xb7\xa2\xa9\xb0\x03\x14\x7d\x2e\xde\x38\x40\xaa\xd2\xf7\xec\x4b\x65\x0a\xf0\x0c\xae\xe4\xbb\x24\xa6\xea\x8a\x64\x9c\x52\xf2\x84\xc9\x58\x39\x84\xc1\xfa\x5f\x16\xc3\xf4\xb0\x48\x3f\x6d\xa9\xfe\xbe\x6d\xdd\xd5\xe5\x72\xf1\xd6\x74\xf5\x35\xce\x18\x47\xee\x43\x53\x8f\x09\xe1\x1a\x46\x1a\x88\xd1\x27\x96\x11\xa5\x5d\x6e\x9b\x26\x04\x39\x66\xa7\xe9\x00\xf1\xb9\x08\x63\x75\x76\xef\x61\x2a\x0f\x93\xbd\x03\x97\xb6\xc3\x39\x3c\x63\xb5\xff\xbf\xb6\x44\x57\x55\x6f\x2a\x2c\xe8\x85\x41\xbd\x76\xca\xbd\x42\x3f\x99\x93\x91\x26\xfa\xa9\x3b\x42\xfc\xde\x1e\x54\xe0\xc5\x87\x37\xf0\x91\xdc\x52\x7b\xf8\xae\x69\x86\x02\x7c\x3f\xde\x36\xf0\xa0\xbd\xb9\x10\xc7\xd9\x70\x97\x6e\x1c\x5c\xdf\xc5\x6d\x3a\x7e\xfb\xf8\xf6\x81\xa3\xc2\xf2\xdd\x84\xcb\x63\xed\x9c\x75\x54\xc9\xe8\x87\x33\x26\xd0\x6f\xce\x0e\xd3\xaf\xbc\xb4\x79\x78\xce\x48\x13\x55\x15\xe0\x6c\x7e\x2f\x6b\xf1\xee\x67\xbd\xaa\x8a\xc8\x92\xab\xd9\xb4\x87\xf3\xc7\x80\xda\xcf\x45\xab\x08\x18\x36\xfa\x6f\x92\x90\xad\xe1\x64\x3b\xc4\xe9\x99\xf2\x04\xd6\x78\xc2\xbf\x74\xc1\xbd\x41\xdb\xab\xca\xae\x4a\x8e\x06\xef\x97\xc0\xbc\x7a\x68\x5f\x42\x3a\x5f\x87\x70\x06\x6e\x5d\x31\x1f\xf6\x87\xc2\xfd\x41\x74\xf0\x60\x38\x74\x1c\xc4\xbd\x66\x37\x62\xa3\xc2\x2a\x29\x4e\x26\xc7\xe9\xd3\xc8\x10\xf0\xe3\xf3\xed\xfd\xe3\xc1\xf6\x3a\x89\xbb\xb7\x0f\x49\xdc\xbd\xb5\xd8\x99\xfc\x13\x00\x00\xff\xff\x9b\xbd\xff\xb2\xe1\x10\x00\x00")

func assetsTmplIndexTmplBytes() ([]byte, error) {
	return bindataRead(
		_assetsTmplIndexTmpl,
		"assets/tmpl/index.tmpl",
	)
}

func assetsTmplIndexTmpl() (*asset, error) {
	bytes, err := assetsTmplIndexTmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "assets/tmpl/index.tmpl", size: 4321, mode: os.FileMode(436), modTime: time.Unix(1484921411, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _assetsTmplLogTmpl = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x6c\x92\x3d\x6f\xc2\x30\x10\x86\xf7\xfc\x8a\xab\x97\x4e\x10\xc1\x6c\xbc\xd0\x6e\xa8\xed\xd0\xa5\xa3\x13\x1f\xb1\x5b\xe7\x2e\x72\x0c\x14\x45\xfc\xf7\x2a\x98\x8f\xc8\x25\x43\x6c\xeb\xbd\xe7\xbd\x0f\xdd\x30\x18\xdc\x3a\x42\x10\x9e\x1b\x71\x3a\x15\xf2\xe9\xe5\x7d\xfd\xf9\xf5\xf1\x0a\x36\xb6\x5e\x15\x72\x3c\xc0\x6b\x6a\x56\x02\x49\xa8\x62\x18\x22\xb6\x9d\xd7\x11\x41\x58\xd4\x06\x83\x80\xf9\x08\x56\x6c\x8e\xaa\x00\x00\x90\xc6\xed\xa1\xf6\xba\xef\x57\xa2\x66\x8a\xda\x11\x06\x91\xb4\xf1\x9b\x5a\x90\xde\x57\xfa\x62\x71\xd5\xa7\xfc\xf7\xae\xad\x38\x06\xa6\x09\x7f\x8e\xb1\x0b\xb5\xe1\xa6\x97\xa5\x5d\xe4\xca\xf2\x0a\x7b\xd4\x46\xa8\xb5\xc5\xfa\x07\x8e\xbc\x0b\xc0\x07\x82\x8a\xe3\x73\x0f\x3e\xb1\xcb\x3b\x2b\x4b\xe3\xf6\xea\x61\x11\x9d\x26\xf4\x70\xfe\xcf\xba\xe0\x5a\x1d\x8e\x79\x39\x79\xf4\x6c\x1c\x47\x16\x94\x7a\x77\x5b\x98\x6f\xb8\x99\xf4\x7b\x33\xe9\x02\xaa\x61\x48\xaa\x2c\xc7\xd7\x03\x1e\x7d\x8f\x8f\x60\xbb\x98\xb6\x0d\x11\x7f\xe3\xac\x46\x8a\x18\xd2\xfd\x56\xf8\x1b\x47\xeb\xa8\x81\xc8\xd0\x5b\x3e\xfc\x9f\xe0\x25\x0f\x99\x2c\x4d\x3e\xa1\xfb\xf3\x72\x95\x65\x5a\x02\x59\xa6\xe5\xb9\x9a\xfc\x05\x00\x00\xff\xff\xbc\x4b\x6b\x86\x66\x02\x00\x00")

func assetsTmplLogTmplBytes() ([]byte, error) {
	return bindataRead(
		_assetsTmplLogTmpl,
		"assets/tmpl/log.tmpl",
	)
}

func assetsTmplLogTmpl() (*asset, error) {
	bytes, err := assetsTmplLogTmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "assets/tmpl/log.tmpl", size: 614, mode: os.FileMode(436), modTime: time.Unix(1484921411, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _assetsTmplNavbarTmpl = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\xd4\x93\xcf\x6e\x9c\x30\x10\xc6\xef\x3c\xc5\xc8\x3d\x53\xbf\x00\x70\xe9\x35\xc9\xa5\x4f\x30\xe0\x31\x19\xd5\x19\x53\xdb\xa0\x54\x28\xef\x5e\xf1\x6f\x43\x58\xb4\x5a\xa9\x52\xa5\x3d\x81\x3f\xbe\x99\xf9\xfc\xb3\x19\x47\x43\x96\x85\x40\x09\x0e\x35\x06\xf5\xf1\x91\x15\x82\x03\x34\x0e\x63\x2c\x57\x15\x96\x47\x6e\xc8\x62\xef\xd2\xb6\x64\x19\x28\x44\xda\x96\x96\xdf\xc9\xe4\xc9\x77\xaa\xca\x00\x00\x0a\xc3\x97\x3e\x8d\x97\x84\x2c\x14\x72\xeb\x7a\x36\xab\xe3\xe8\x5a\x1b\xbd\x12\x1a\x0a\x3b\xcf\xec\xab\xfb\x94\xbc\x1c\xac\xc9\xb7\xad\x23\x68\xbc\x73\xd8\x45\x32\x0a\x0c\x26\x5c\xe5\x69\xec\xa2\x6f\x32\x86\x96\x52\xa9\xbe\x2d\xd5\xcf\x24\x7d\x54\x80\x81\x31\xa7\xf7\x0e\xc5\x90\x29\x95\x45\x17\xe9\x30\x7c\x0e\x10\x3b\xbc\x8c\xe7\xc6\x4b\x3e\xf1\xaa\x0a\x3d\xe9\xff\xd3\x5e\xe8\x05\xc5\x41\xc5\x03\x9a\x3a\xa0\x18\x05\xaf\x81\x6c\xa9\xb4\xaa\x9e\xff\xd4\x3e\x15\x1a\x77\xe8\xb5\xe1\xe1\xe6\x49\x6c\xfc\xe0\x13\x24\x9b\xed\xeb\x82\xef\x90\xa2\x77\xbb\x16\xdb\xcd\x10\x1c\x4e\x80\x8e\x23\x5b\xa0\xdf\xf0\xfd\x65\x76\xbd\xe0\x1b\x81\xfa\xe1\xc5\x72\x3b\x5d\xc3\x2b\x42\x8e\x77\xad\x73\x4e\xf4\x06\xd8\x24\x1e\xce\x0e\x6b\x1c\xc9\x45\xba\xaf\xcb\x79\xb9\x98\x93\xea\x2b\xd0\xb9\x63\xf9\x75\x81\xdc\xcc\xe1\xb5\xaa\x96\x5d\x7c\x81\xfd\x09\xdd\x71\x95\xdd\x07\xe3\xc9\x3f\x2a\x09\xe7\x27\x0c\x4f\xfe\xdf\x19\xfc\x4c\x98\xfa\xf8\xa0\x18\xe2\x1c\x5e\xab\x6a\xd9\xc5\x0d\x18\x5f\x95\xde\x9d\xfe\xa5\xeb\x6b\xa1\x05\x87\x2a\xdb\x32\xfd\x0d\x00\x00\xff\xff\x0c\x54\x48\x91\xc4\x05\x00\x00")

func assetsTmplNavbarTmplBytes() ([]byte, error) {
	return bindataRead(
		_assetsTmplNavbarTmpl,
		"assets/tmpl/navbar.tmpl",
	)
}

func assetsTmplNavbarTmpl() (*asset, error) {
	bytes, err := assetsTmplNavbarTmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "assets/tmpl/navbar.tmpl", size: 1476, mode: os.FileMode(436), modTime: time.Unix(1484921411, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _assetsTmplStatusTmpl = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\xec\x97\x41\x8f\xda\x3e\x10\xc5\xef\x7c\x8a\xf9\xe7\xfe\xdf\x68\xf7\x58\x79\x23\x55\xd0\x6e\x55\x15\x15\x09\x7a\xe8\xd1\x89\x07\x70\xd7\xd8\x91\x67\x92\x0a\xa1\xfd\xee\x95\xe3\x44\xa4\x2c\x74\xcb\xc2\xf6\x94\xbd\xac\xc2\x9b\x79\x6f\xec\x9f\x11\xf1\x6e\xa7\x70\xa9\x2d\x42\x42\x2c\xb9\xa2\xe4\xe9\x69\x24\xfe\x9b\x7c\x1d\x2f\xbe\xcf\x3e\xc0\x9a\x37\x26\x1b\x89\xf0\x0f\x8c\xb4\xab\xfb\x04\x6d\x92\x8d\x76\x3b\xc6\x4d\x69\x24\x23\x24\x6b\x94\x0a\x7d\x02\x37\xa1\x31\x77\x6a\x9b\x8d\x00\x00\x84\xd2\x35\x14\x46\x12\xdd\x27\x85\xb3\x2c\xb5\x45\x9f\x44\x2d\xfc\xf5\x2d\xac\xac\x73\xd9\x5a\x74\x7a\xbf\xff\x47\xb5\xc9\x1d\x7b\x67\x7b\xfd\x4d\xcd\xfa\x36\x9b\x79\x57\x20\x11\xcc\x9b\xe9\x45\xba\xbe\x3d\xac\xb9\xeb\x6c\x0c\x4a\x95\x64\xe3\x35\x16\x8f\x50\xb6\x6d\x71\xd1\x48\x20\xad\x02\x96\x8f\x08\x1b\x94\x54\x79\x24\xd0\x4b\x90\x76\x1b\x2a\x73\x83\x9b\x60\x7d\xb7\xb7\x16\xa9\xd2\x75\x76\x74\x5a\x96\xb9\xc1\xff\x3d\x52\xe9\x2c\xe9\x1a\x0f\x87\x6e\xf4\xdf\x8a\x21\xb6\x10\x7b\x5d\xa2\x6a\x9f\x72\xe7\x15\x7a\x54\x07\xed\xd1\x22\x6c\xfa\xf3\xcf\xa3\xe6\x8f\x0b\x51\x54\xdd\x86\x89\x94\x4f\x38\x74\x85\xdd\x8e\x9e\xaa\x13\xe9\xb1\x28\x91\x9e\x18\x4e\xf0\xfe\x70\x9c\x3d\xf4\xe2\xa7\x66\x46\x0f\x13\xed\xb1\x60\x98\x22\x91\x5c\x21\x7c\xd1\xc4\x68\xd1\xbf\xbc\x98\x93\x22\x34\x67\x51\x2f\xe1\xa6\xcd\x88\x9e\x93\x69\x5c\x7e\xef\x44\x1e\xb5\xa6\x52\x5a\x20\xde\x1a\x0c\xc7\xdc\x38\xff\x2e\x37\x15\x26\xd9\xfb\x82\x75\x8d\x22\x0d\xfa\x4b\xe1\x68\x08\xcf\xcf\x69\x8e\xc6\x9c\x5d\xf9\x97\x21\x56\xfd\x21\xe3\x5c\xca\x70\x0e\xb4\x6f\x84\xfe\x0d\x51\x05\x7b\x1a\x68\x75\xca\xa5\xb4\x66\xe8\xb5\x53\xba\x80\xcf\x2e\xbf\x0a\xac\x88\xa6\x63\x36\x80\xea\x94\xd7\x82\x1a\x3b\xbb\xd4\x2b\xf8\xa8\x0d\xc2\xd4\x59\xcd\xce\x6b\xbb\xba\x26\xaa\xd6\x35\x06\x0d\xc0\x3a\xe5\xd2\x6f\xd6\xd8\xa3\x42\xcb\x5a\x9a\xb7\xe5\xd6\xe6\x85\xb8\x01\xdb\xeb\xb1\x3d\x38\xb7\x32\x08\x63\xe3\x2a\xf5\xcf\xd8\x3d\x34\x71\x03\xba\xcb\xd0\x7d\x5a\x2c\x66\xd7\x7d\xe3\x68\x29\xcd\xd1\xd7\xc3\x6f\xd8\x5e\x39\x7e\x01\x78\xfe\xa2\x2f\xd2\xe6\x52\x73\xfa\xfe\x54\x99\x6c\xd4\x13\xc3\x63\xaf\x4c\xa4\xd1\x52\xa4\xf1\x42\xda\xcd\xfc\x2b\x00\x00\xff\xff\x26\xae\xe4\xfa\xbd\x0e\x00\x00")

func assetsTmplStatusTmplBytes() ([]byte, error) {
	return bindataRead(
		_assetsTmplStatusTmpl,
		"assets/tmpl/status.tmpl",
	)
}

func assetsTmplStatusTmpl() (*asset, error) {
	bytes, err := assetsTmplStatusTmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "assets/tmpl/status.tmpl", size: 3773, mode: os.FileMode(436), modTime: time.Unix(1484921411, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _assetsTmplTwitter_setupTmpl = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\xc4\x54\xc1\x6e\xdb\x30\x0c\xbd\xe7\x2b\x38\xde\x13\xa1\x3d\x0e\xb2\x2f\xdd\x4e\x43\xb1\x01\xcd\x65\xa7\x82\xb6\x98\x58\xab\x2c\x09\x12\x9d\xc0\x08\xfa\xef\x83\xed\xa4\x4d\xdd\xa4\xe8\x61\xc0\x72\x88\x25\x3c\xf2\xe9\xf1\x91\xd2\xe1\x60\x78\x63\x3d\x03\xca\xde\x8a\x70\x7a\xcc\x2c\x5d\xc4\xe7\xe7\x85\xfe\xf2\xed\xe7\xdd\xfa\xf7\xaf\xef\xd0\x48\xeb\xca\x85\x1e\x3e\xe0\xc8\x6f\x0b\x64\x8f\xe5\xe2\x70\x10\x6e\xa3\x23\x61\xc0\x86\xc9\x70\x42\x58\x0d\x89\x55\x30\x7d\xb9\x00\x00\xd0\xc6\xee\xa0\x76\x94\x73\x81\x75\xf0\x42\xd6\x73\xc2\x09\x1b\x7e\xe7\x14\x9e\x76\x15\x1d\x29\x4e\xf8\x79\xfe\x9f\xae\xad\x82\xa4\xe0\xcf\xf2\xc7\x98\xe6\xa6\x5c\x4f\xe2\xe1\x61\x10\xaf\x55\x73\x33\x0f\xb9\x3d\xb1\x38\x26\x83\xe5\xba\xb1\x19\x6c\x06\x69\x18\x36\x36\x65\x81\x2c\x1c\x41\x02\x74\x99\xe1\xbe\xaf\x82\x68\xd5\xdc\xce\x58\x36\x21\xb5\x60\x4d\x81\x93\x47\x40\xb5\xd8\xe0\x0b\x54\xe3\x5e\x1d\x1d\x54\x08\x2d\x4b\x13\x4c\x81\x31\x64\x41\x60\x5f\x4b\x1f\xb9\xc0\xb6\x73\x62\x23\x25\x51\x03\xd3\xd2\x90\x10\x96\x7a\xdc\xcc\x4e\xb2\x3e\x76\x02\x03\xf0\x72\xd8\x44\x91\xbb\xaa\xb5\x82\xb0\x23\xd7\x71\x81\x0f\xb4\x63\x3c\x95\x56\x89\x87\x4a\xfc\x32\x26\xdb\x52\xea\xc7\xb5\xdb\x9e\xb9\xa5\x95\xb1\xbb\x73\xf3\xed\x06\x56\xf7\x9c\x33\x6d\xf9\x8a\xe9\xe4\x38\x09\x8c\xff\x4b\x43\x7e\x3b\x74\xef\x70\x78\x4d\x7a\x47\xc9\xde\x7c\x82\xca\xfa\x4d\x98\xb7\x91\xa0\x49\xbc\x29\xb0\x11\x89\xf9\xab\x52\x14\x63\x5e\x1d\x3d\x5d\xd5\xa1\x45\x10\x4a\x5b\x96\x02\x1f\x2b\x47\xfe\x09\xcb\x75\x80\x3a\x31\x09\xab\xba\xe1\xfa\x09\xfa\xd0\x25\x08\x7b\x0f\xa7\x71\xa0\x18\x9d\xad\x69\xe8\x12\x6c\xc3\xd0\xde\x86\x13\x6b\x45\x57\x3d\x39\x17\x2c\x54\x39\x5e\x26\xce\x31\xf8\x6c\x77\x3c\x17\x3c\xe2\x6f\x82\x61\x4a\xc9\x92\x6c\x64\x73\xdc\x55\x21\x19\x4e\x6c\x66\xe9\x13\xc5\xeb\x55\x79\x8f\xa5\xcb\xc0\x04\x9a\xf2\x2e\xf8\xdc\xb5\x9c\xe0\x07\xf7\x5a\x89\xf9\x38\xfa\x2a\x08\x1f\xce\x5b\xa4\x9c\xf7\x21\x19\x04\x4f\x2d\x17\x6f\x5f\x89\x55\x7d\x94\xf0\xf8\xc4\xfd\x85\xf2\x5e\xe8\xaf\xaa\xd3\xea\x5a\x95\x9f\x2f\xff\x81\xeb\xc4\xf2\xbf\x1d\xc8\xa3\x8a\x7f\x68\x82\x56\x17\x86\x43\xab\x71\xa6\x2e\x8e\xef\x71\xa9\xd5\x94\xa6\xd5\xf4\x6c\x9f\xae\xe4\xdf\x00\x00\x00\xff\xff\xbd\x9a\x2e\x64\xea\x05\x00\x00")

func assetsTmplTwitter_setupTmplBytes() ([]byte, error) {
	return bindataRead(
		_assetsTmplTwitter_setupTmpl,
		"assets/tmpl/twitter_setup.tmpl",
	)
}

func assetsTmplTwitter_setupTmpl() (*asset, error) {
	bytes, err := assetsTmplTwitter_setupTmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "assets/tmpl/twitter_setup.tmpl", size: 1514, mode: os.FileMode(436), modTime: time.Unix(1485265001, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"assets/css/custom.css": assetsCssCustomCss,
	"assets/tmpl/config.tmpl": assetsTmplConfigTmpl,
	"assets/tmpl/header.tmpl": assetsTmplHeaderTmpl,
	"assets/tmpl/index.tmpl": assetsTmplIndexTmpl,
	"assets/tmpl/log.tmpl": assetsTmplLogTmpl,
	"assets/tmpl/navbar.tmpl": assetsTmplNavbarTmpl,
	"assets/tmpl/status.tmpl": assetsTmplStatusTmpl,
	"assets/tmpl/twitter_setup.tmpl": assetsTmplTwitter_setupTmpl,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}
var _bintree = &bintree{nil, map[string]*bintree{
	"assets": &bintree{nil, map[string]*bintree{
		"css": &bintree{nil, map[string]*bintree{
			"custom.css": &bintree{assetsCssCustomCss, map[string]*bintree{}},
		}},
		"tmpl": &bintree{nil, map[string]*bintree{
			"config.tmpl": &bintree{assetsTmplConfigTmpl, map[string]*bintree{}},
			"header.tmpl": &bintree{assetsTmplHeaderTmpl, map[string]*bintree{}},
			"index.tmpl": &bintree{assetsTmplIndexTmpl, map[string]*bintree{}},
			"log.tmpl": &bintree{assetsTmplLogTmpl, map[string]*bintree{}},
			"navbar.tmpl": &bintree{assetsTmplNavbarTmpl, map[string]*bintree{}},
			"status.tmpl": &bintree{assetsTmplStatusTmpl, map[string]*bintree{}},
			"twitter_setup.tmpl": &bintree{assetsTmplTwitter_setupTmpl, map[string]*bintree{}},
		}},
	}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}

