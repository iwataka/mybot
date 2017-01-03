// Code generated by go-bindata.
// sources:
// assets/css/custom.css
// pages/config.html
// pages/header.html
// pages/index.html
// pages/log.html
// pages/navbar.html
// pages/setup.html
// pages/status.html
// DO NOT EDIT!

package mybot

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

var _assetsCssCustomCss = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x8c\x90\xcb\x6e\xf3\x20\x10\x85\xf7\x7e\x0a\xa4\x7f\x6d\x2b\x92\xff\x45\xe5\x48\x7d\x97\x01\x06\x42\x8d\x19\x34\x8c\x73\x69\xd5\x77\xaf\xec\x5c\x44\x94\xca\x2a\x0b\x60\x71\xbe\x73\xce\x4c\xf7\x31\x4f\x9a\x84\x29\xa9\xaf\x46\xdd\x8e\x06\x33\x7a\xa6\x39\xd9\xd6\x50\x24\x1e\x94\x30\xa4\x92\x81\x31\xc9\xfe\x21\x13\x3c\x4b\x0b\x31\xf8\x34\x28\x83\x49\x90\xf7\xcd\x77\xd3\x74\x86\x92\x40\x48\xc8\x95\xe5\x04\xe7\xf6\x14\xac\x1c\x06\xf5\xb6\xdb\xe5\xf3\x55\x99\x21\x61\x6c\x25\x48\xc4\x4a\xbb\xfa\xae\x89\x8e\x78\x1a\xd4\x9c\x33\xb2\x81\x82\x2b\x64\xc3\xb1\x2a\x0d\xaa\x64\x48\x9d\x83\x8a\x77\x94\xa4\x2d\xe1\x13\x07\xf5\xff\x1e\x95\xb9\x4e\xd0\xc4\x16\x79\x50\x89\x12\xee\xff\x3a\xf6\x52\xf8\x29\xe7\x55\x7f\x3a\x04\xc1\x32\xd1\x88\x77\x79\x2b\xa7\x20\x82\xbc\x7c\xc1\x48\x38\xd6\x35\x6e\xd4\xbf\xbe\xef\x7b\xe7\x1e\x88\x27\xf2\x11\x37\x09\xe7\x16\xe6\x41\x64\x28\xe5\x57\xa1\x67\xb8\x54\x55\x10\xc5\xa2\x19\x5f\x75\x3a\x82\x19\xb7\x36\x61\xc9\x7a\x64\x1d\xe7\xeb\x64\x02\x3a\xa2\x7a\x57\x72\x40\xb0\x9b\x2b\x81\x18\x0c\x3e\x73\xdd\x7a\xb7\x45\x38\x64\xb4\x8b\x8b\x26\x7b\x59\x5e\xee\x2c\x46\x14\xac\x2d\x29\x83\x09\x72\x19\xd4\xae\xeb\xb7\x1a\x7a\xc6\x75\xd4\x9f\x00\x00\x00\xff\xff\x05\xeb\x61\x17\xd2\x02\x00\x00")

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

	info := bindataFileInfo{name: "assets/css/custom.css", size: 722, mode: os.FileMode(436), modTime: time.Unix(1483434595, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _pagesConfigHtml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xec\x5d\x5b\x6f\xe3\xb8\x15\x7e\xcf\xaf\x60\xd5\xc5\x66\x17\x3b\x91\xdb\x9d\xb7\x8e\x6d\x60\x76\x32\x83\x99\x22\x9b\xa6\x49\xa6\x40\x51\x14\x06\x6d\xd1\x16\x77\x68\xd1\xa5\x28\x27\x41\x90\xff\x5e\x90\x94\x64\x5d\x48\x4a\x76\x2e\x92\x32\xe4\xc3\x8e\x6d\xf1\xf2\xf1\x3b\xe7\x90\x87\xf4\xe7\xec\xf8\x4f\xa7\xff\xf8\x70\xfd\xef\x8b\x8f\x20\xe4\x6b\x32\x3d\x1a\x8b\x7f\x00\x81\xd1\x6a\xe2\xa1\xc8\x9b\x1e\xdd\xdf\x73\xb4\xde\x10\xc8\x11\xf0\x42\x04\x03\xc4\x3c\xe0\x3f\x3c\x1c\x8d\xe7\x34\xb8\x9b\x1e\x01\x00\xc0\x38\x5e\x30\xbc\xe1\xea\x8d\x28\xcb\x24\x5a\x70\x4c\x23\x10\x20\x82\x38\xba\xa4\x37\x3f\xcd\x13\xce\x69\xf4\x06\xac\x69\x00\xc9\xcf\xe0\x3e\xaf\x2b\xca\x16\x32\xc0\x19\x98\x00\x55\xcb\xdf\x40\x86\x22\x7e\x4e\x03\x54\x78\xf9\xae\xd4\x04\x2f\xc1\x4f\x9c\xf9\x0b\x02\xe3\xf8\x1c\xae\x11\x98\x4c\x26\xe0\x58\x8d\x17\x1c\x57\x07\x90\xa0\x28\x03\x3f\x89\x91\x30\x98\x80\xbf\xbc\x03\x18\x8c\x81\xe8\x21\xc4\x24\x10\xfd\xc7\x3e\x41\xd1\x8a\x87\xef\x00\xfe\xe5\x17\x5d\x07\x19\x52\xd9\x02\x4c\xca\x8d\xff\x83\xff\xab\x6d\x20\x70\xca\x5a\x3e\xbf\xdb\xa4\x28\x43\x1c\x04\x28\x3a\x06\x3f\xfe\xa8\xfa\xf2\x23\x1a\xa0\xdd\x24\xbe\x9c\x5f\x7c\xbd\xd6\x4e\x21\x2b\xaa\xd5\x16\x92\x04\x81\x09\x38\x5e\x42\x12\xa3\xe3\x77\xc6\xea\x65\x9e\xc0\xb1\xa5\x6a\x6a\x00\x1c\x45\x88\x7d\xbe\xfe\xfd\x4c\x54\x3f\x95\xa4\x1a\x1a\x3d\xd4\x3e\x2d\x7f\xf2\x00\x10\x89\xd1\x77\x65\x0e\xce\x92\x7d\xac\x91\xf9\xec\x7e\x46\xb9\x44\x5b\xc4\xf8\xa1\x46\x39\x2a\xbf\x1a\x8f\x8a\x11\x3c\x0e\xf0\x16\x48\x88\x13\x6f\x41\x23\x0e\x71\x84\x98\xb7\x8b\xee\xe2\x8a\x10\xc1\xed\x1c\xa6\x2b\x42\x5e\x61\xbc\xa4\x6c\x0d\x70\x30\xf1\xe8\x16\xb1\x1b\x86\x39\xf2\x00\x94\x2b\xc2\xc4\x1b\x2d\x68\xb4\xc4\xab\x91\x07\xd6\x88\x87\x34\x98\x78\x1b\x1a\x73\x0f\xa0\x68\x21\x4c\x32\xf1\xd6\x09\xe1\x78\x03\x19\x1f\x89\x6e\x4e\x02\xc8\xa1\x37\x1d\xcb\x37\xd3\xc2\x18\xe1\xaf\xd3\x0f\xb2\xa7\xf1\x28\xfc\x75\x07\x6e\x8c\xa3\x4d\xc2\x85\x77\xad\x4b\xc3\xab\xbe\xe3\x64\xbe\xc6\xdc\x03\xd2\x56\x13\xef\x0a\x6e\x91\x97\x4d\x75\xce\x23\x30\xe7\xd1\x09\x59\xc9\x7f\x36\x0c\xaf\x21\xbb\xf3\x4a\x63\xbe\x9d\x5e\xdf\x60\xce\x11\x03\xd7\x78\x8d\x08\x8e\x50\x3c\x1e\x85\x6f\x0b\xc3\x73\x38\x27\x28\xeb\x52\xbd\x91\xff\x3d\x89\x39\xc3\x1b\x14\xa4\xef\xe6\x94\x05\x88\xe5\x6f\x19\x8a\x37\x34\x8a\xf1\x16\x15\x78\x56\xfd\x89\x35\x77\x5a\xb3\xe7\x98\xb3\xfa\x87\xea\x41\x30\x8d\x17\x0c\xa1\x08\x44\x70\x2d\xe0\x71\x4d\xf3\xac\x26\xba\x5d\x90\x24\x40\x80\xa1\x0d\xc1\x4d\x95\x71\x94\x56\xe6\x0d\x15\x17\x34\x89\xb8\xbd\x0a\xe5\x21\x62\x0d\xdd\xe8\x9f\x8e\x47\xd5\xa9\x8f\x47\x1a\x92\xc6\x7c\xb7\x3f\x15\xcb\xfd\x3d\x83\xd1\x0a\x81\x1f\x70\x14\xa0\xdb\x37\xe0\x07\x9e\x1a\x12\xfc\x6d\x02\x7c\xe5\x51\x7e\x6a\x64\x3f\x37\xf2\x43\x3d\xa2\xcc\x16\xb0\x7a\xa0\x5a\x6c\x3c\x69\x9c\x89\xc7\xd3\x81\x32\x10\xb1\x9f\xae\x06\xb9\x8b\xca\xb5\xdd\x1b\x99\x59\x32\xae\x1a\xf7\xf7\x04\xc7\xfc\x1a\xdd\xf2\x39\xbd\xdd\xcd\xd3\xbf\x92\xde\x21\x56\x9f\x18\x68\x00\x28\xe7\x99\x49\xe7\xf1\x34\xf3\x4e\x19\xb7\xd8\xcd\x82\x68\x4e\x29\xb9\x42\x04\x2d\x2a\x98\x3e\x2a\x3f\xbc\x54\x6e\xa8\x83\x95\x7a\xea\x2c\xf5\xd4\x17\x43\xf6\x45\x39\xfd\x25\xd7\xa2\x4a\x43\x62\xc6\xf8\x53\x23\xb2\x7b\x51\x94\xac\xe7\x22\x11\x33\x79\x91\x8c\xc0\xdc\x87\xee\xef\x77\xf3\xf9\x20\x9e\x3c\x3c\x78\x06\x48\x07\x82\x55\x3b\x54\x0a\x4e\xbd\xa9\xad\xab\x38\x5a\x52\x0f\x88\xe5\xfc\x84\xd3\xd5\x8a\x88\xa5\x5e\x24\x83\xd9\x67\x90\xad\x10\x9f\x78\x7f\x4e\x67\x73\x92\xcf\xe6\x64\x89\x09\x47\x6c\x76\x7f\xaf\x82\xd6\x88\x3e\x2b\x01\xe2\x10\x93\x18\x84\x88\x21\x33\xe6\x91\xc2\xf9\xf2\x44\x04\x62\x01\x62\x1e\xa0\xd1\x82\xe0\xc5\xb7\x89\xb7\xcb\x92\x79\x88\xe3\x37\xe0\xb8\x05\x05\xc7\x3f\x7b\x53\x95\x98\x1d\x32\x91\xfa\x22\x0a\x64\x14\xa0\x28\xd8\x6b\xa9\xb3\x72\x21\xb3\x80\xea\xb6\x9f\x4f\x69\x04\x83\xa0\x92\x03\xd8\xad\x9a\x46\x84\x76\x23\x3f\x47\x37\x35\x9a\xb3\x0d\xdc\xb0\x76\x2a\x1a\x54\x56\xa1\x7f\xb6\xff\xde\xd4\xd5\x53\xed\x9e\x58\xde\xff\xc6\x23\x99\x69\x14\xb2\x99\x27\xd9\x09\x8b\xb9\xa2\x8c\x66\xb0\x84\x01\xf2\x64\xf2\xd7\x26\x90\x45\xfe\x23\x5f\x4f\xbc\x93\xbf\x7a\x80\x51\xb1\x2c\x04\x18\x12\xba\xf2\x00\x64\x18\x9e\x10\x38\x47\x84\xa0\x60\x7e\xd7\xaa\xc7\x13\x8e\x39\xa9\x66\x51\x5a\xa4\x27\xd9\x30\xe9\xa0\x74\x91\xac\x51\x64\xf2\xc2\x7a\x73\x91\x1a\x9b\xeb\xeb\xdb\xa4\x07\xe8\x06\x47\x0f\xdf\x96\x5b\xa9\x29\xb5\xe5\x34\x63\xe0\x54\xad\x82\xe5\xfc\xb4\xee\x39\x01\xde\xee\x33\x03\xe1\x55\xcd\xf8\xa7\x9f\x24\x2a\xfb\xd8\xe0\x19\xf2\x65\xfd\x18\xfa\x1c\x5a\x5f\xd7\xb0\xd4\xe9\x2b\x07\x53\xb1\x03\x9b\x43\xd6\xd4\x4a\x2e\x5b\xed\x9b\xe9\x97\x6b\x4d\xad\x76\x13\x35\xe5\xc6\xfa\xba\x7b\x12\xb2\x81\xc2\x41\x23\x4b\x6e\x6f\x6a\xd9\xba\x32\xb0\x64\xb7\xca\xf1\xfc\x8b\x14\x85\x2e\x69\x53\x11\xe3\x67\x40\x4d\x89\x9b\x16\xe6\x13\x9b\x0c\x1c\xc2\x70\xc2\x08\xe8\x07\xcb\x5f\x2f\xcf\x5a\x10\x9d\x30\x32\x1b\x2c\xd9\x21\x8c\xc1\x1a\x05\x18\x3e\x3b\xd3\xa6\x13\x48\xca\xf5\x67\x18\xff\x2e\x70\x58\x88\x0e\x61\x3c\x93\x58\x07\xc9\x72\xc2\x48\x0f\x38\xfe\x7a\x79\xd6\xc0\x70\xc2\xc8\xe0\xf8\x65\x88\xdf\x20\xc4\x51\xd0\x35\xc3\x97\x19\x10\x0b\xc9\x39\xd8\xc1\xd1\xbc\x84\x5b\x2a\x8e\xe9\x80\x87\x0c\xc5\x21\x25\xcf\xcd\xf7\xe3\x6e\x09\x52\xbe\x33\xd4\xb3\x1c\xb5\xf6\xe6\x20\xb5\xe0\xa7\xb4\xf6\x75\x56\xb9\xf1\x38\x5e\x02\xdc\x07\x33\xe5\x0e\x36\x2c\x3b\xe5\xb0\xdb\x19\x2a\x0f\xb5\x01\x5b\x8a\xc0\x68\xd5\xa9\x69\x38\xba\xe5\x8d\x86\x11\x28\x6d\x96\x38\x83\xd1\xaa\x53\xee\x6b\xb7\x01\x86\x4a\xea\x86\xc0\x5a\x29\x3f\xdb\xf9\xff\xc2\x31\xa6\x91\x3b\xe2\x19\x5a\xbd\xe2\x23\x9e\xbc\x92\xe9\xf6\xe4\xa1\x9c\xcf\x3f\x13\x48\x2c\xa9\xc4\x56\x55\x93\x80\x07\x98\x4d\x2c\x10\x90\x77\xc5\x80\xe0\x6f\x88\xe0\x90\xd2\x6e\x37\xaa\x56\xab\x61\xca\xb9\x40\xef\x4b\xf4\xb3\x1d\x7a\xdb\x2a\x99\x9a\xf4\x93\x68\xf7\x5e\xb4\x3b\xcb\x9b\x0d\x6f\xe3\x92\xb6\x9b\x93\x84\x89\x05\x6d\xa8\xd6\x4b\xf1\x1f\x60\xbf\xdf\x54\xcb\xc1\x5b\x50\x2c\xb6\x37\x08\x0e\x37\x00\xb3\x09\x1c\x60\xc3\xcf\x69\xd3\xc1\x1b\xf1\x0f\x7a\x37\x58\xfb\xfd\x41\xef\x0e\x30\xdd\xdf\xe9\xdd\x90\xad\x26\x58\xea\x45\x7a\x21\x9e\x37\x67\x17\xd2\xa8\x43\x4b\x2e\x08\x8c\x82\x35\x64\xdf\x7a\xc1\xf3\x59\x0a\xa6\x4d\x26\xa7\x6a\x0e\x8f\x6f\xba\xa2\xfd\xe0\x9a\xae\x68\x0b\x9e\xe9\x8a\x76\xc9\xf1\xd3\x9e\x58\xdf\x4b\x11\x82\x3b\xaa\x1a\x5a\xbd\xe2\xa3\x6a\x7a\x67\xf6\xec\x91\xb7\x08\xd1\xe2\x5b\x39\xec\x94\xd3\x65\xd7\x70\xba\x90\x53\xda\x98\xec\x5e\x6f\x70\x4b\x5a\x76\x6f\xdc\x21\xb9\xd9\x65\xb4\x85\xdd\x0c\xe5\xf0\xe8\xa5\x84\xd0\x9b\x2e\xc9\x95\x00\x6c\xd4\xca\x0a\x83\x23\x76\x41\x09\x41\x72\x06\x5d\x7d\x7d\x9e\x12\xfc\x61\x07\xc4\xc2\x72\x01\xee\xf0\x77\xe4\xbd\xf5\x46\x4b\x4a\xb9\x55\x31\x65\xe9\xd1\xf0\x48\xf3\x71\xe5\xa3\x4c\x08\x99\x7f\x50\x54\xfe\x67\x6b\xce\xb0\x95\xff\x83\xd2\xe8\xe7\xdf\xab\x6a\x94\x89\xb9\x39\x9e\x5b\xa3\x9f\x81\x78\x46\x8d\x7e\x36\x84\x5e\xa3\xbf\x03\xf0\x78\x8d\xfe\x21\x5f\x48\xee\xc6\xaf\xc9\xcb\x73\xe0\xb9\xbc\xbc\x4f\x4a\xf2\x1c\xf8\xf7\xab\x24\xb7\x50\x30\x54\x25\x79\x3e\x25\xa7\x24\xaf\x8d\xf7\x04\x6a\xf0\x43\xd6\xdc\x36\x6a\x70\x5b\x30\x1e\xa6\x06\xb7\xf4\xf8\x9d\xa8\xc1\x5b\x30\xe0\xd4\xe0\xf9\x18\xee\xfe\xa5\x82\xe0\xd5\xa9\xc1\xf3\x74\xc4\xa8\x06\xdf\xe5\x32\x4e\x0d\xde\xba\x34\xb0\xac\x55\x83\xd7\x88\x76\x6a\xf0\x16\xa5\xa6\xa3\xad\x72\x5d\x57\x83\xd7\x88\x76\x6a\x70\x7b\x69\xc3\x71\x49\x0d\xae\x65\xd8\xa9\xc1\xcd\xa5\x91\x61\x8d\x1a\xbc\x46\xb2\x53\x83\xb7\x2e\x8f\x3b\xd4\xb7\x52\x83\x57\x2d\xe8\xd4\xe0\x5d\xd9\xa9\x41\x0d\x6e\x0c\xb5\x01\x5b\xaa\x97\x6a\xf0\x9a\x61\xaa\x6a\xf0\xaa\x25\x9c\x1a\xdc\x1d\xf1\x1a\x9b\x0c\xe8\x88\xd7\x81\x1a\xbc\x1a\x52\x7a\x35\x78\x2d\x32\x9d\x1a\xfc\xe5\x57\xc3\xb6\x6a\x70\x83\x49\x9d\x1a\xbc\xd4\x51\x97\xd6\xb3\xab\xc1\x6d\xf6\x73\x6a\xf0\x42\x4f\x5d\x9a\xb0\x41\x0d\x6e\xb3\xa1\x53\x83\xf7\xc0\x7e\x66\x35\xb8\xcd\x74\x4e\x0d\xde\x54\xda\xa5\x17\x65\x35\xb8\xc9\x50\x4e\x0d\x6e\x2e\x6d\xd3\xb8\xaa\x1a\xdc\x9c\xc9\x39\x35\xb8\xa9\xb4\xe4\xba\xa4\x06\x37\xf2\xec\xd4\xe0\xee\xa8\x6a\x6a\x32\xa0\xa3\xea\xcb\xab\xc1\xf3\xb0\x33\xa9\xc1\x77\x21\xe7\xd4\xe0\x6d\x8a\x8d\xdc\xba\x1a\xbc\xc6\xae\x53\x83\xdb\x8a\x95\xdc\x8a\x1a\xbc\x4e\xad\x53\x83\xdb\x8b\x69\x47\xb6\xa9\xc1\x6b\x2c\x3b\x35\x78\xbf\xd4\xe0\x57\x08\xb2\x45\xd8\x4b\x31\xf8\xff\x12\xc4\x1a\xff\xa8\x3b\x43\x71\x42\x94\x16\xf2\x15\x09\xc6\x63\x69\x15\x9d\x74\x31\xb3\xd7\x73\xab\xc5\xe3\x74\x9c\xe7\x13\x8b\xab\x11\xfc\x7f\x2a\x2b\x83\xfa\xc8\xa9\xfd\x9f\xfc\x4f\xa5\xc7\xbb\x2f\xf4\x53\x0c\x97\xd2\x89\xae\xef\x36\x48\x03\x43\x79\xd8\x4c\x90\xe5\x01\xcf\x03\x1e\x43\x0b\x14\x71\x0f\x78\x6b\x7c\x2b\x88\xf1\x36\x74\x93\x10\xc8\x5e\x4e\xc9\x9e\x63\xab\x09\xd9\xd3\x09\xf5\x52\xc6\x9e\xa1\xfe\x7e\x55\xec\x66\x06\x86\x2a\x62\xcf\x66\xe4\x34\xec\xfb\x3c\x7d\x1a\x85\xfb\xde\x9b\x44\x1b\x79\xbb\x25\x48\x0f\x53\xb7\x9b\x3b\xfc\x4e\xc4\xed\xcd\x04\x38\x6d\x7b\x3e\x86\xbb\x4d\xaa\x20\x78\x75\xda\xf6\x34\x43\x31\x2a\xdb\xf3\xdc\xc6\x09\xdb\x5b\x17\x2b\xc5\x5a\x59\x7b\x95\x65\xa7\x6a\x6f\x51\x6a\x7a\xe0\x32\xd1\x75\x4d\x7b\x95\x65\x27\x69\xb7\x97\x66\x82\x4b\x82\x76\x1d\xbd\x4e\xcf\x6e\x2e\x0d\xf4\x6a\xd4\xec\x55\x86\x9d\x98\xbd\x75\x79\xd4\xb9\xbe\x95\x96\xbd\x6c\x3d\xa7\x64\xef\xc8\x48\x0d\x42\x76\x43\x8c\x0d\xd8\x4c\xbd\x94\xb1\x57\xad\x52\x55\xb1\x97\xcd\xe0\x34\xec\xee\x28\xd7\xd8\x64\x40\x47\xb9\x0e\x34\xec\xe5\x80\xd2\x2b\xd8\xab\x41\xe9\x04\xec\x2f\xbe\x0e\xb6\xd5\xaf\x6b\xcd\xe9\xd4\xeb\xa5\x8e\x3a\x34\x9d\x5d\xbc\x6e\x36\x9e\x93\xae\x17\x7a\xea\xd0\x7e\x0d\xca\x75\xb3\x01\x9d\x6e\xbd\x7b\xe3\x99\x65\xeb\x66\xbb\x39\xd1\x7a\x53\x69\x93\x4f\x94\x25\xeb\x06\x1b\x39\xc5\xba\xb9\xb4\x4b\xda\xaa\x7a\x75\x63\xde\xe6\xe4\xea\xa6\xd2\x8a\xe8\x92\x58\xdd\x44\xb2\xd3\xaa\xbb\x23\xa9\xa9\xc9\x80\x8e\xa4\x2f\xaf\x55\x4f\x63\xce\xa4\x54\xcf\xe3\xcd\x09\xd5\xdb\x14\x33\xb3\x75\x99\x7a\x95\x5a\xa7\x52\xb7\x15\x0b\xb3\x15\x8d\x7a\x8d\x57\x27\x51\xb7\x17\xfd\x2e\x6c\x13\xa8\x57\x29\x76\xfa\xf4\x7e\xe9\xd3\x2f\x88\x38\x3c\x9e\x53\x8e\x97\x78\x01\xeb\xc9\x44\x2f\x94\xea\x50\x86\x6d\x8c\xc8\xd2\xae\x06\x4c\x62\xa3\xbe\xfc\x91\x0a\xf2\xc3\x34\x9d\x85\x95\xa8\x2a\x21\x2c\x32\xee\x4b\x23\xf8\xef\xc5\x2c\xaf\x10\x59\xee\xa2\x27\x2a\xd6\xda\xc8\x5a\x92\x8b\x99\xe0\xe2\xc9\x85\xdb\xc5\xd8\x6e\x81\xf7\xab\x60\xdb\x8a\x55\xda\x63\x3f\x98\x07\x49\x37\x85\x43\x9f\x42\x0e\xe7\x30\x46\x3d\x74\xdf\x80\xe1\x2d\x62\x76\xb3\x04\x90\x43\x10\xd3\x84\x2d\x1a\x7e\x67\xa1\x8e\x31\x0a\x67\x8f\x5c\xbd\xfd\x8d\x50\x30\xf7\x15\x21\x85\x8b\x9e\xcc\xdb\x4e\x7f\xf3\x4f\xe5\x33\xe3\x8d\xce\xa1\x42\xf3\xbd\xe0\x41\x0e\x67\xca\x16\x06\x8c\x90\xc3\x2b\xf9\xbc\x53\x9c\xca\x15\x66\xd2\x15\xf4\x40\xd5\xc1\xf8\x5a\x54\xd8\x13\xe9\xc1\x71\xf8\x25\xe2\x88\x41\xcd\xb1\xb4\x1f\xa1\x98\xb0\x74\x93\xb3\x99\xe7\x09\xf6\x1b\x00\x1a\x7f\x16\xd5\xf3\x38\xc5\x3b\x43\xfa\x19\x6d\x05\x2f\xc3\xcb\x7c\x8f\x28\x98\xdc\x3f\x4d\x6b\x3e\x3c\xec\x1c\xd1\xf4\x5c\xa6\x25\x4f\x1b\x3f\x9a\x1d\xb7\x38\x7a\x61\x8b\x2d\x4e\xef\x85\xf7\xd4\x22\xa2\x74\x13\x2d\xa2\x39\x60\xd7\x6c\x04\xb2\x8f\x10\xa7\x88\xa5\xfa\x03\x28\xdd\x14\xf2\x1f\x43\xbd\xc4\xf2\x72\x46\x57\x3d\x5c\x56\x5e\x51\x82\x7a\x46\x57\xc5\x30\x21\x74\xf5\xd2\xe1\x21\x10\xa4\x61\x21\x46\x7f\xc9\x24\xf2\xf3\xf5\xf5\x45\x0f\xdd\xcb\x7e\xf9\x29\xf5\xc3\x34\x6e\xf8\xf1\xed\x86\xb2\x86\x1a\x28\x12\xd0\x2d\x5f\x1d\xa6\x57\xfd\x40\xfe\xef\xc7\x7a\xe4\xc5\xed\xf7\xb4\x90\xf3\x8d\x2f\x5e\x6a\x96\x34\x61\x7b\xff\x1c\xae\x3b\xcc\xe9\x24\x3c\x61\x4a\x13\xbc\xcf\x34\xde\x77\xa5\x7d\x24\xbc\xf2\xd6\x20\x01\x0a\x4f\x32\x01\xbc\xa0\xec\xc9\x01\x6a\x16\x29\x39\xd6\x47\xe5\xaf\x40\xa1\x4a\xbd\xb7\xc3\x7d\x53\xc2\x20\x74\x35\x93\xf1\x61\x62\xe8\x8c\xae\xce\xc4\xf3\x67\xdc\x30\xc1\xee\x0e\x68\x3c\x52\x75\xc6\xa3\x90\xaf\xc9\xf4\xe8\xff\x01\x00\x00\xff\xff\x41\xad\x18\x2a\xcf\xa3\x00\x00")

func pagesConfigHtmlBytes() ([]byte, error) {
	return bindataRead(
		_pagesConfigHtml,
		"pages/config.html",
	)
}

func pagesConfigHtml() (*asset, error) {
	bytes, err := pagesConfigHtmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "pages/config.html", size: 41935, mode: os.FileMode(436), modTime: time.Unix(1483434595, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _pagesHeaderHtml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x8c\x53\x5d\x4f\xe3\x3a\x10\x7d\xe7\x57\x58\x7e\xb9\x0f\xf7\x36\xbe\xfd\x58\xda\x5d\x35\x95\x58\x58\x4a\x61\xf9\x58\x68\x41\xec\x9b\xeb\x4c\x12\x87\xd8\x0e\x9e\x49\xd3\x6c\xd5\xff\xbe\x6a\x0a\x14\x21\x90\x78\xcb\xcc\xe4\x1c\x9f\xa3\xa3\xb3\x5a\x45\x10\x6b\x0b\x8c\xa7\x20\x23\xf0\x7c\xbd\xde\x1b\x6e\x3e\x47\x7b\x8c\x31\x36\x34\x40\x92\xa9\x54\x7a\x04\x0a\xf9\x6c\x7a\xdc\x1a\xf0\xd7\x27\x2b\x0d\x84\x7c\xa1\xa1\x2a\x9c\x27\xce\x94\xb3\x04\x96\x42\x5e\xe9\x88\xd2\x30\x82\x85\x56\xd0\x6a\x86\xff\x98\xb6\x9a\xb4\xcc\x5b\xa8\x64\x0e\x61\x9b\x8f\xf6\x1a\xa6\xd5\x4a\xc7\x2c\x98\x21\xf8\x0b\x69\x60\xbd\xde\xd2\x93\xa6\x1c\x46\xab\xd5\xab\xc3\x3f\xc8\xe6\x8e\x86\x62\x7b\x7a\xc2\x42\x8e\x6f\x30\xe7\xf5\x3b\x7f\xd9\x68\xbd\xde\x3e\x37\xcc\xb5\x7d\x60\x1e\xf2\x90\x23\xd5\x39\x60\x0a\x40\x9c\xa5\x1e\xe2\x90\xa7\x44\x05\x7e\x13\xc2\xc8\xa5\x8a\x6c\x30\x77\x8e\x90\xbc\x2c\x36\x83\x72\x46\xbc\x2c\x44\x37\xe8\x06\x7d\xa1\x10\x77\xbb\xc0\x68\x1b\x28\x44\xce\xb4\x25\x48\xbc\xa6\x3a\xe4\x98\xca\xee\xa0\xd7\xfa\x7e\x7b\xaf\xf5\xcd\xe4\x18\xce\xda\xd1\xd8\x9c\x5e\x1f\x3c\xd4\xaa\x3c\x39\x38\xb9\x4e\xba\x9d\x4b\x33\x53\x55\xd5\x77\xb6\x7b\x7d\x1f\x25\xbd\x5b\xf9\xef\x95\xb9\x99\xe2\x1f\x71\xb6\x3f\x58\xcc\xa3\x1f\x59\xda\x2b\x39\x53\xde\x21\x3a\xaf\x13\x6d\x43\x2e\xad\xb3\xb5\x71\x25\x3e\xa7\xf1\x19\x53\x2a\xb2\x19\x06\x2a\x77\x65\x14\xe7\xd2\x43\xe3\x48\x66\x72\x29\x72\x3d\x47\x11\x3b\x4b\x2d\x59\x01\x3a\x03\xa2\x17\xec\x07\xdd\xc6\xde\xeb\xf5\x8b\xc3\x4f\xbc\x2a\x24\x22\x10\x36\x1c\xaa\x44\x72\xe6\x09\xb9\x85\xa2\xf2\xba\x20\x86\x5e\xed\x04\x6e\xb4\x04\x89\x73\x49\x0e\xb2\xd0\xf8\x46\x5f\xf6\x58\x82\xaf\x45\x3b\x68\x77\x82\xde\xd3\xd4\x08\xca\x90\x8f\x86\x62\x4b\x38\xfa\x98\xfd\xb3\x99\x66\x6f\x23\xcd\xde\x4d\x74\xaa\xbe\x4c\x7e\xe9\xf9\xff\x9d\xfe\xe3\xa2\xce\x6e\xce\xe3\x93\xec\xf2\x5c\xfe\x7c\x88\xcb\xbb\xdb\xe5\xef\xe5\xec\xca\x1e\x9e\x1e\xf4\xf3\x8e\x39\xbc\xbb\x98\x14\xe3\xaf\x66\x7c\x78\x34\xa8\xc6\x17\x13\x75\x75\xd4\x9f\x2e\xe5\xc7\x89\xee\xbc\x7c\x6c\xa6\xc8\x25\xc5\xce\x9b\x80\x2a\x4d\x04\xbe\x71\x52\xe9\x28\x01\xc2\x46\xf0\x4b\x69\x4b\x8a\x37\xa5\xdd\x91\x0e\xc5\xb6\xe0\xcf\xad\xf8\x1b\x00\x00\xff\xff\xf8\x1d\xf8\x44\x04\x04\x00\x00")

func pagesHeaderHtmlBytes() ([]byte, error) {
	return bindataRead(
		_pagesHeaderHtml,
		"pages/header.html",
	)
}

func pagesHeaderHtml() (*asset, error) {
	bytes, err := pagesHeaderHtmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "pages/header.html", size: 1028, mode: os.FileMode(436), modTime: time.Unix(1483434595, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _pagesIndexHtml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xc4\x58\x4b\x8f\xdb\x36\x10\xbe\xef\xaf\x98\xb2\x01\xda\x02\x91\x85\x24\x40\x0f\x01\x2d\x20\x4d\x82\x36\x40\x92\x06\x69\x52\xa0\xc7\x91\x44\x4b\xec\x52\x1c\x95\xa4\xbd\x70\x05\xff\xf7\x82\xa2\xed\xb5\x9e\x5e\x3b\x69\xbb\x07\xaf\xa4\x79\x7d\x9c\xb7\xc4\xbf\x79\xf5\xeb\xcb\x4f\x7f\x7c\x78\x0d\xa5\xab\x54\x72\xc3\xfd\x3f\x50\xa8\x8b\x25\x13\x9a\x25\x37\x4d\xe3\x44\x55\x2b\x74\x02\x58\x29\x30\x17\x86\xc1\x62\xb7\xbb\xe1\x29\xe5\xdb\xe4\x06\x00\x80\xe7\x72\x03\x99\x42\x6b\x97\x2c\x23\xed\x50\x6a\x61\x58\xa0\xf9\xbf\x53\x15\x1a\x37\x29\xee\x55\x1c\x19\x4e\x15\xfc\xb9\xae\x52\x72\x86\xf4\x89\x82\xa0\x44\xae\x60\xf1\xd9\x0a\xf3\x1e\x2b\xb1\xdb\x75\x88\xbc\x7c\x92\x34\xcd\x09\xf5\x3b\x0b\x29\x39\x1e\x97\x4f\xfa\x5a\x84\xb2\x63\xd2\xef\xb6\x53\xfc\x3a\x1f\xb0\x3f\x4d\x70\xed\xa8\x42\x27\x33\x54\x6a\x0b\x19\x29\x25\x32\x07\xa8\x73\x48\x0d\x61\x9e\xa1\xf5\x77\x5b\xb8\x95\x3a\xb7\x40\x2b\x90\x7a\x45\xc6\x4b\x90\x86\x15\x19\xd8\xd2\x9a\xc7\xe5\xd3\xb1\x33\xfe\x44\x6e\xec\x88\x08\xa5\x11\xab\x25\x2b\x9d\xab\xed\xf3\x38\x76\x77\xd2\x39\x61\x16\x19\x55\x71\xd3\xdc\x4b\xb1\x83\x27\x53\xa7\x21\x75\x3a\xb2\x94\x49\x54\x91\xcc\x28\xdc\xef\x05\x59\xc2\x6d\x8d\xfa\xc0\xbd\x42\x58\xe1\x81\xe6\x2f\x31\x73\x72\x23\x58\xc2\x63\xcf\x96\xf0\x18\x1f\xe2\xca\x3e\x4a\xac\x6b\xbb\x38\x81\xfa\xb5\xd0\xd5\x68\xed\x39\x78\xc3\xc8\x0d\x7d\x28\x84\xcb\x45\x76\xfb\x6f\x40\x3c\x2a\x9f\x06\x39\x00\x94\x91\xb6\xa4\xc4\x22\x17\x1b\xa1\xa8\x16\xc6\x2e\x0a\xa2\x42\x89\x36\xcc\x58\x4b\xeb\x7f\xe2\x8d\xb4\x92\xf4\x9e\xe4\x9f\xb6\xe4\xbf\xd6\xe4\xd0\x9e\x85\x1f\xa4\x46\xd1\x07\xd2\xb9\xe8\xf3\x38\x97\x9b\xe4\xa4\x7c\xd7\xea\xa0\x46\xe3\x06\x34\x6e\x22\x87\xa9\xed\x14\x30\x57\x12\x0c\x29\xb1\x64\xb5\x11\x56\x68\xd7\x96\xc2\x11\xeb\xd1\x1c\x42\x8e\x0e\x23\x47\x45\xe1\x99\x1d\xa6\x6c\xef\xa2\x6f\x4b\xaa\x04\x4b\x7e\xa1\x4a\x78\x34\x3c\x56\xf2\xd4\x40\xd3\x18\xd4\x85\x80\x47\xb7\x62\xfb\x18\x1e\x6d\x50\xc1\xf3\x25\x2c\x5e\x86\xd2\x94\xa4\xdf\x61\xdd\x49\x87\x09\x44\xb3\x10\x9a\xc6\xab\xdf\xed\x58\x72\xb8\x1a\x87\xd2\x4d\x3d\x1e\xaf\x55\x32\xde\xed\x1c\xa6\x91\x6f\x99\x42\xbb\x5e\xbf\x6b\xb9\x64\xbe\x64\xed\xb1\x4f\xf9\x6b\xd4\x3e\x44\xb9\x00\xa9\xe1\xe0\xb8\x8e\x6c\xdf\x8a\x97\x50\xd0\xfe\x46\xb9\xf7\x92\x19\x11\x18\x15\x8a\x7c\xcb\x97\xba\x98\xe0\x6f\x65\xca\x67\x5d\x11\x27\x9d\x0a\xa1\x0c\x5e\x8b\x15\x15\x31\x4b\x5e\x1b\x43\x06\xde\x52\x11\x5c\x56\x3e\x9b\x80\x10\x92\x6b\x8c\x14\xda\xe3\x5b\x2a\x7a\x65\x3d\x7f\x02\x3f\xa7\x58\xc2\x6b\x23\x12\x9e\x51\x2e\xfc\xa0\x68\x75\xf0\xb8\xbd\xe5\x71\x4b\x9a\x33\x3b\xd2\xe8\xce\xd9\x9b\x06\x28\xaa\xe4\x3d\xb9\x52\xea\x02\x1c\x81\x2d\xe9\x8e\xc7\xa2\xba\xc2\x19\xc3\xfe\x36\x23\x32\x95\x0e\x7e\x2e\xfd\x97\xc9\x30\xdf\x79\x93\x4f\xfb\xf6\x79\x5f\xb6\xf6\xea\x74\xb9\x38\x34\x21\xbf\xa6\x3b\x46\x4f\xbd\x2f\xea\x29\x22\x5c\xd3\x91\x46\x6c\x74\x1b\xcb\x04\xd3\xc1\xb7\x4d\xe3\x8d\xf4\xbb\xd3\xbc\x81\xf8\x9c\x85\xa9\x3c\xbb\xd7\x30\xe7\x87\xd9\xda\x81\x4b\xcb\xe1\x1c\x9e\xa9\xdc\xff\x5f\x4b\x22\x64\xd5\x9b\x0a\x0b\xf1\x42\xa3\xda\x5a\x69\x5f\xa1\x9b\xf5\xc9\x44\x11\xfd\x1c\xa6\xf3\xef\xed\xf0\x87\x17\x1f\xde\xc0\x47\x61\xd7\xca\xc1\xf7\x4d\x33\x66\xe0\x87\xe9\xb2\x81\x07\xc5\xe6\x42\x1c\x67\xcd\x5d\x1a\x38\xb8\xbe\x8a\x5b\x77\x7c\xfe\xf8\xf6\x81\xa3\xc2\xd0\xdd\x8c\xca\x3e\x77\x46\x2a\xaa\xf2\xe8\xc7\x33\x22\xd0\x2d\xce\x80\xe9\x37\x5a\x9b\xcc\xaf\xe9\x09\x97\x55\x01\xd6\x64\xf7\xb4\x16\xef\x71\xd6\xcb\xaa\x88\x8c\xb0\x35\xe9\x76\xd3\x7d\x0c\xa8\xdc\x92\xb5\x8c\x80\x3e\xd0\x7f\x8b\x1c\xd2\x2d\x0c\xc2\xc1\x86\x8b\xe6\x00\xd6\xb4\xc3\xbf\xf4\xc0\x9d\x41\xdb\xc9\xca\x90\x25\xbd\xc1\xfb\x25\x30\xaf\x1e\xda\x97\x34\x9d\xaf\xd3\x70\x46\x1e\x5d\x31\x1f\x8e\x4b\xe1\x71\x11\x1d\x5d\x0c\xc7\xd6\x41\x3c\x72\x86\x11\x1b\x15\x46\xe6\x6c\x30\x39\x86\xaf\x28\x63\xc0\xfb\xfb\xed\xfd\xfb\xc0\xfe\x9a\xc7\xe1\xf3\x00\x8f\xc3\x67\x85\x7f\x02\x00\x00\xff\xff\xec\xe4\xca\x70\x67\x10\x00\x00")

func pagesIndexHtmlBytes() ([]byte, error) {
	return bindataRead(
		_pagesIndexHtml,
		"pages/index.html",
	)
}

func pagesIndexHtml() (*asset, error) {
	bytes, err := pagesIndexHtmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "pages/index.html", size: 4199, mode: os.FileMode(436), modTime: time.Unix(1483434595, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _pagesLogHtml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x4c\x8f\xc1\x0a\x83\x30\x0c\x86\xef\x7d\x8a\xac\x0f\x60\xc1\x73\xec\x65\xf3\x26\x6c\x87\x5d\x76\x8c\x36\x58\xa1\xb6\x52\x8b\x30\xc4\x77\x1f\xd2\x6d\x9a\x4b\x08\x7f\xbe\x8f\x04\x2f\xb7\xfb\xf5\xf9\x7a\xd4\x60\xd3\xe8\xb4\xc0\xbd\x81\x23\xdf\x57\x92\xbd\xd4\x62\x5d\x13\x8f\x93\xa3\xc4\x20\x2d\x93\xe1\x28\xa1\xd8\x36\x81\x6d\x30\x6f\x2d\x00\x00\xd0\x0c\x0b\x74\x8e\xe6\xb9\x92\x5d\xf0\x89\x06\xcf\x51\xe6\x6c\xaf\xb3\xc2\xd3\xd2\xd2\x57\xf1\xcb\xd1\x96\xba\x8e\x31\x44\x68\x42\x8f\xca\x96\x07\x8a\x53\xe4\x63\xca\xb2\xa2\x09\xfd\x99\x56\xff\x1d\x54\x66\x58\xb4\x40\x95\x6f\x43\x95\x7f\xfa\x04\x00\x00\xff\xff\xb4\x77\x0d\x57\xe4\x00\x00\x00")

func pagesLogHtmlBytes() ([]byte, error) {
	return bindataRead(
		_pagesLogHtml,
		"pages/log.html",
	)
}

func pagesLogHtml() (*asset, error) {
	bytes, err := pagesLogHtmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "pages/log.html", size: 228, mode: os.FileMode(436), modTime: time.Unix(1483434595, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _pagesNavbarHtml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xbc\x92\x3d\x52\xc4\x30\x0c\x85\xfb\x9c\xc2\x63\xea\x8c\x2f\xe0\x75\x43\x43\x01\x15\x27\xd0\xc6\x4a\xd0\xa0\x95\x77\xe2\x9f\x26\x93\xbb\x33\xf9\x5b\xb2\xa1\xa0\x81\xad\xa2\x3c\x3f\xe9\x9b\x27\x7b\x18\x3c\xb6\x24\xa8\xb4\x40\x39\x43\xaf\xc7\xb1\xb2\x02\x45\x35\x0c\x31\x9e\x56\x55\x2d\x9f\xda\x63\x0b\x99\xd3\xf6\x4b\x52\xb0\x8f\xa8\x5d\xa5\x94\x52\xd6\xd3\xad\xad\x09\x92\x80\x04\xfb\xba\xe5\x4c\x7e\x75\xcc\xae\x73\x4e\x29\xc8\xfd\xfc\x3a\x85\xae\x63\x54\x4d\x60\x86\x6b\x44\xaf\x95\x87\x04\xab\x3c\x8d\x5b\xf4\x4d\x86\xbe\xc3\x74\xd2\x4f\x4b\xf7\x1b\x4a\x8e\x3b\xc4\x8c\x89\x57\xb8\x41\xa8\x09\x52\x4f\xe1\x9c\x35\x93\xfe\xdf\x56\x6b\x96\x90\xeb\x5e\x8c\xa7\xf2\x73\x45\x6b\xf2\x2d\x9a\xfa\xce\x48\x7e\x3b\x3d\x26\xb3\x99\x77\xed\xdb\x35\x08\x94\x63\x7a\xa6\x9d\xaf\xa6\x84\x97\x83\x63\x76\xc1\xde\xc4\x24\x9f\x5a\x7d\xf4\xd8\x9e\xb4\xd1\xee\x25\x5c\xd0\x1a\x38\xcc\x35\x4c\x7f\x4c\x6a\x82\xb4\xd4\x19\xed\x9e\xe7\xe2\x11\x48\x0e\x13\xef\x35\x3c\x04\x16\x13\xa4\x1c\x8d\x76\xef\x73\xf1\x0b\xd2\x9a\xcc\x77\xaf\xc6\x1a\x81\xe2\xaa\x61\x40\xf1\xe3\x58\x7d\x05\x00\x00\xff\xff\x8a\xff\x1c\xc2\xaf\x03\x00\x00")

func pagesNavbarHtmlBytes() ([]byte, error) {
	return bindataRead(
		_pagesNavbarHtml,
		"pages/navbar.html",
	)
}

func pagesNavbarHtml() (*asset, error) {
	bytes, err := pagesNavbarHtmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "pages/navbar.html", size: 943, mode: os.FileMode(436), modTime: time.Unix(1483434595, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _pagesSetupHtml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xbc\x55\x4d\x8f\xdb\x20\x10\xbd\xe7\x57\x4c\xb9\x7b\x39\xec\x15\x47\xaa\xd2\x8f\x43\x55\xb5\xd2\xe6\xd2\xe3\x18\x26\x0e\x5a\x0c\x08\xc6\x59\x59\xd6\xfe\xf7\xca\x38\xd9\xa6\x71\x2a\xf5\x90\x8d\x0f\x20\xfc\xe0\x3d\xde\x43\x30\xea\xc3\xa7\x1f\x9b\xed\xaf\x9f\x9f\x61\xcf\x9d\x5b\xaf\xd4\xd4\x81\x43\xdf\xd6\x82\xbc\x58\xaf\xc6\x91\xa9\x8b\x0e\x99\x40\xec\x09\x0d\x25\x01\x0f\xaf\xaf\x2b\xd5\x04\x33\xac\x57\x00\x00\xca\xd8\x03\x68\x87\x39\xd7\x42\x07\xcf\x68\x3d\x25\x31\x63\xd3\x77\x4e\xe1\xf1\xd0\xe0\x91\xe2\x84\xab\x5d\x48\x1d\xa0\x66\x1b\x7c\x2d\x64\x26\xee\xa3\x14\xd0\x11\xef\x83\xa9\x45\x0c\x99\x05\x90\xd7\x3c\x44\xaa\x45\xd7\x3b\xb6\x11\x13\xcb\x69\x59\x65\x90\xf1\x4c\xab\xf0\x59\x1f\x7b\x86\x79\x7a\xee\x9b\xce\xb2\x38\xed\xaf\x61\x0f\x0d\xfb\xca\xb5\xa5\x8b\xc9\x76\x98\x06\x01\x07\x74\x3d\xd5\xe2\x09\x0f\x74\xc9\xb6\x7f\x5c\x6f\x5f\x2c\x33\x25\x25\xf7\x8f\x7f\x83\xe3\x68\x77\xf0\xf0\x9d\x72\xc6\x96\xce\x2c\x5d\xc6\x82\x8e\x12\x43\x69\x2b\x83\xbe\x9d\xf2\x19\xc7\x3f\x0b\x95\x34\xf6\x70\x49\x4d\xde\x5c\x52\x32\x36\x8e\x4e\xa4\xf3\xa0\xb4\x55\xe6\x64\x23\x99\xe3\xa8\x09\xc9\x50\x7a\x1b\x26\xca\x31\xf8\x6c\x17\xde\x66\xce\xb4\xfc\x39\x03\x66\xbd\x09\x3e\xf7\x1d\x25\xf8\x46\x83\x92\x6c\xfe\x3d\xf3\x2a\xb0\x38\x8d\x88\x39\xbf\x84\x64\x04\x78\xec\xa8\x16\x3c\x07\x5b\xe9\xa3\x4e\xf5\x4c\xc3\x95\x3d\x16\x9e\xab\xf2\x4a\x5e\xdb\xfe\xff\x79\x7a\x22\x9d\x88\xef\x62\x2b\x17\xa9\xf7\x76\xf6\x51\x6b\xca\x19\xb6\xe1\x99\xfc\x7b\xda\xc2\xa2\x53\xf1\xa4\x73\x4f\x4f\x77\x38\xb1\x73\x6b\xb7\x39\x35\x25\xcb\x25\x5c\x3e\x2a\x5f\x43\x68\x1d\xc1\xc6\x85\xde\x2c\x5f\x96\xfb\x5f\xf5\x44\x86\x3c\x5b\x74\xf0\xc5\x3a\xba\x45\xc8\x3b\xeb\xe8\x14\x70\xab\x27\x9f\x95\x7e\x53\xa9\x0a\x7a\xe3\x6c\x55\x29\x0a\xc7\xb2\x34\xbf\xa9\x4a\xce\x95\x4a\xc9\xb9\xc2\xfd\x0e\x00\x00\xff\xff\x2d\x55\xde\x47\xf2\x06\x00\x00")

func pagesSetupHtmlBytes() ([]byte, error) {
	return bindataRead(
		_pagesSetupHtml,
		"pages/setup.html",
	)
}

func pagesSetupHtml() (*asset, error) {
	bytes, err := pagesSetupHtmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "pages/setup.html", size: 1778, mode: os.FileMode(436), modTime: time.Unix(1483434595, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _pagesStatusHtml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xec\x97\x41\x6b\xdb\x30\x14\xc7\xef\xfe\x14\x6f\xbe\x2f\x86\x1e\x87\x2a\x18\xd9\x96\x30\x56\x1a\x68\x76\xe8\x51\x8a\x5e\xe2\x07\xaa\x64\x9e\x9e\x5d\x42\xe9\x77\x1f\x8e\x9c\xcd\x5d\x3b\x32\x48\x7a\x4a\x7c\xb0\xfd\xde\x1f\xfd\x6c\xfd\x10\x46\x56\x1f\xbe\xdc\x4e\x97\xf7\x8b\xaf\x50\xcb\x83\xd7\x85\xea\x2f\xe0\x4d\xd8\x5c\x97\x18\x4a\x5d\x3c\x3d\x09\x3e\x34\xde\x08\x42\x59\xa3\x71\xc8\x25\x4c\x9e\x9f\x0b\x65\xa3\xdb\xea\x02\x00\x40\x39\xea\x60\xe5\x4d\x4a\xd7\xe5\x2a\x06\x31\x14\x90\xcb\x9c\xf5\xc7\x18\x11\x4c\x67\xcd\x80\xd8\xe7\xaa\xbe\xd2\x77\x62\xa4\x4d\xaa\xaa\xaf\xfe\x8c\x53\x62\xac\xc7\x3d\x39\x17\xbb\xf3\xc7\x24\x4c\x0d\xba\xa1\xb2\x91\x1d\xf2\xef\x92\x31\x35\x31\x24\xea\x70\xf4\x0e\x99\xc7\x2f\x1b\xb9\xe9\xf4\xf2\x91\x44\x90\xe1\x3e\xb6\x0c\xb7\x8f\x01\x7e\x50\x12\x0c\xc8\xaa\x12\xf7\xf6\x90\x57\xcd\x3c\x51\x5a\xc3\x24\x4f\x65\x32\x40\x33\xea\x66\x9b\xd0\xaf\x73\x32\x9a\xf9\x0b\x6a\x6a\x4c\x80\x24\x5b\x8f\xbd\x46\x1f\xf9\x93\xf5\x2d\x96\xfa\xf3\x4a\xa8\x43\x55\xf5\xf9\xbf\x9e\x8b\x3e\xe1\xff\x73\x19\x5d\xa9\xef\x24\x36\x07\xa0\xc1\xbd\xc1\x7c\xed\x44\x55\x7f\x8b\x3d\x68\xfa\x67\x42\x3e\xbd\xe5\x9e\x9a\x2e\x92\x67\x24\xf3\xd6\xc2\x02\x99\xa2\xa3\x15\x7c\x8f\xf6\x18\xc5\x33\x92\xba\xb5\x17\xad\xfb\xb5\x7b\x2a\xaf\x03\xef\x22\x76\x1a\xc3\x9a\x36\xf0\x8d\x3c\xc2\x4d\x0c\x24\x91\x29\x6c\x8e\x51\x3b\x50\x32\xf8\x22\x78\xbf\x72\xa7\x8c\x0e\x83\x90\xf1\xa7\xf5\x3c\xf0\x7b\xfc\x39\x6b\x9e\xc5\xb8\xf1\x08\x53\x1f\x5b\xf7\x6e\xae\x67\x3b\xfc\xb9\xab\x9e\x2f\x97\x8b\x63\x64\xce\x45\x9a\xb3\xf8\x2e\xa8\x6a\xb7\x29\x1e\x35\x5a\xaf\x8b\x51\xdc\x97\xf9\xce\x51\xa7\x0b\x55\xe5\x7d\xbd\xaa\xf2\xff\xc0\xaf\x00\x00\x00\xff\xff\x74\x63\x93\xad\x20\x0c\x00\x00")

func pagesStatusHtmlBytes() ([]byte, error) {
	return bindataRead(
		_pagesStatusHtml,
		"pages/status.html",
	)
}

func pagesStatusHtml() (*asset, error) {
	bytes, err := pagesStatusHtmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "pages/status.html", size: 3104, mode: os.FileMode(436), modTime: time.Unix(1483434595, 0)}
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
	"pages/config.html": pagesConfigHtml,
	"pages/header.html": pagesHeaderHtml,
	"pages/index.html": pagesIndexHtml,
	"pages/log.html": pagesLogHtml,
	"pages/navbar.html": pagesNavbarHtml,
	"pages/setup.html": pagesSetupHtml,
	"pages/status.html": pagesStatusHtml,
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
	}},
	"pages": &bintree{nil, map[string]*bintree{
		"config.html": &bintree{pagesConfigHtml, map[string]*bintree{}},
		"header.html": &bintree{pagesHeaderHtml, map[string]*bintree{}},
		"index.html": &bintree{pagesIndexHtml, map[string]*bintree{}},
		"log.html": &bintree{pagesLogHtml, map[string]*bintree{}},
		"navbar.html": &bintree{pagesNavbarHtml, map[string]*bintree{}},
		"setup.html": &bintree{pagesSetupHtml, map[string]*bintree{}},
		"status.html": &bintree{pagesStatusHtml, map[string]*bintree{}},
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

