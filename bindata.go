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

var _assetsTmplConfigTmpl = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\xec\x9d\x5b\x6f\x23\xb7\x15\xc7\xdf\xfd\x29\xd8\x69\x10\x27\xc8\x5a\x4a\xb2\x6d\x1f\xb2\x92\x80\xed\x6e\x16\xd9\xc2\xbb\x71\x7d\x29\x50\x14\x85\x41\xcd\x1c\x49\xcc\x52\xa4\x4a\x72\x7c\x81\xe0\xef\x5e\x70\x6e\x1a\xcd\x90\x9c\x91\xac\xab\x97\x7e\xc8\xda\x12\x0f\xe7\xf0\xfc\xc9\xc3\x43\x4a\x3f\x64\x3e\x8f\x60\x44\x18\xa0\x20\xe4\x6c\x44\xc6\xc1\xd3\xd3\x49\xef\x4f\xef\x7f\x7f\x77\xfd\xef\x8b\x5f\xd1\x44\x4d\xe9\xe0\xa4\xa7\xff\x41\x14\xb3\x71\x3f\x00\x16\x0c\x4e\xe6\x73\x05\xd3\x19\xc5\x0a\x50\x30\x01\x1c\x81\x08\x50\x47\x1b\x0e\x79\xf4\x38\x38\x41\x08\xa1\x9e\x0c\x05\x99\xa9\xf4\x0f\xfd\x33\x8a\x59\xa8\x08\x67\x28\x02\x0a\x0a\x2e\xf9\xfd\x77\xc3\x58\x29\xce\x5e\xa1\x29\x8f\x30\xfd\x1e\xcd\x8b\xb6\xfa\xe7\x0e\x0b\xa4\x04\xea\xa3\xb4\x55\x67\x86\x05\x30\xf5\x99\x47\x50\xfa\xf5\xcd\x92\x09\x19\xa1\xef\x94\xe8\x84\x14\x4b\xf9\x19\x4f\x01\xf5\xfb\x7d\x74\x9a\x3e\x2f\x3a\xad\x3e\x20\x71\x8a\x0b\xf4\x9d\x7e\x12\x41\x7d\xf4\xe3\x1b\x44\x50\x0f\xe9\x1e\x26\x84\x46\xba\x7f\xd9\xa1\xc0\xc6\x6a\xf2\x06\x91\x1f\x7e\x30\x75\x90\x7b\x9a\x58\xa0\xfe\xb2\xf1\x7f\xc8\x7f\x8d\x06\xda\xcf\xa4\x55\x47\x3d\xce\x32\x2f\x27\x24\x8a\x80\x9d\xa2\x6f\xbf\x4d\xfb\xea\x30\x1e\xc1\x62\x10\x1f\x3f\x5f\xdc\x5c\x1b\x87\x90\xff\xa4\x56\x77\x98\xc6\x80\xfa\xe8\x74\x84\xa9\x84\xd3\x37\xd6\xe6\xcb\x71\x42\xa7\x8e\xa6\x99\x00\x84\x31\x10\xbf\x5d\x7f\x3a\xd7\xcd\xdf\x27\x41\xb5\x18\x3d\xd5\x5e\x5d\x7e\xe5\x09\x01\x95\xf0\x55\xc9\xa1\x44\xbc\x8a\x1a\xf9\x9c\x5d\x4d\x94\x4b\xb8\x03\xa1\xd6\x15\xe5\x64\xf9\xb7\x5e\xb7\xbc\x82\x7b\x11\xb9\x43\x89\x8b\x7d\x9d\x28\x14\x26\x0c\x44\xb0\x58\xdd\xe5\x8c\xc0\xf0\xdd\x10\x67\x19\xa1\x68\x50\xee\xe0\x8f\x78\x3a\xe4\x4a\x70\x56\xea\x20\x69\x33\xf9\x69\xf0\x2e\xc9\x42\xb1\xc0\x3a\x57\xf4\xba\x93\x9f\xaa\x4d\x7e\xce\x7b\xa1\x80\xa3\xa0\x68\x0f\xe8\x91\xc7\x02\xf1\x7b\x86\x86\x5c\x9d\x4a\x24\x79\x2c\x42\x90\xaf\xd0\x88\x50\x05\x42\x22\xcc\x22\x84\x93\x14\x24\x7b\xdd\xc9\xcf\x95\x7e\x47\x5c\x4c\x11\x89\xfa\x01\xbf\x03\x71\x2f\x88\x82\x20\x6b\xdd\x0f\xba\x69\x6a\xec\x06\x68\x0a\x6a\xc2\xa3\x7e\x30\xe3\x52\x05\x08\x58\xa8\x67\x4c\x3f\x98\xc6\x54\x91\x19\x16\xaa\xab\xbb\x39\x8b\xb0\xc2\xc1\xa0\x97\xfc\x51\x79\x0c\x61\xb3\x58\xe9\x79\x3e\x5d\x7a\x52\xda\x8d\x8c\x87\x53\xa2\x02\x94\xcc\x9a\x7e\x70\x85\xef\x20\xc8\x47\x3b\x54\x0c\x0d\x15\x3b\xa3\xe3\xe4\x9f\x99\x20\x53\x2c\x1e\x03\x24\xd5\x23\x85\x7e\x10\x11\x39\xa3\xf8\xf1\x17\xc2\x28\x61\x10\x0c\x4e\x0c\xa3\xab\x0e\x67\x44\x28\x94\xc6\x34\x06\x65\xef\xae\x3a\x77\xb2\x81\x18\xfd\x7e\xcf\xef\x19\xe5\x38\x6a\xf2\xbd\x12\x9a\x2c\x5c\xab\x3b\xde\x42\x8c\x8d\x8c\xeb\x66\xb6\xfa\xa8\x6a\x7d\x6a\xdf\x03\xc4\xf0\x54\x7b\xfa\x38\xe4\xaa\x93\xed\xbb\xad\x5c\xac\xce\xa9\x5e\x37\x22\x77\xa5\x98\xcd\xe7\x64\x84\x3a\x9f\x40\x4a\x3c\x86\xa7\x27\xe3\xea\xc3\x14\x84\x42\xc9\x7f\xcf\x22\xcc\xc6\x7a\x1d\xcf\xe7\x0b\xa3\xac\xcf\x45\x97\xc0\x22\xbd\x92\x17\x9d\x4d\x5e\x0f\xae\xef\x89\x52\x20\xd0\x35\x99\x82\xf6\x53\x2f\xa9\xd7\x03\xe3\xf3\x14\x1e\x52\x38\x13\x20\x67\x9c\x49\x72\x57\x1b\x52\xf2\xfe\x52\x63\x94\x9a\x48\x25\xc8\x0c\xa2\xec\xaf\x21\x17\x11\x08\x88\x4c\x01\x56\xba\x08\xa9\xbf\x9e\xbe\x27\xcc\x6f\xa4\x6f\x46\x03\x19\x0a\x00\x96\x28\x22\x7b\x5d\x65\xe9\x26\x6f\x0d\x0f\x21\x8d\x23\x40\x02\x66\x94\xb4\x31\x20\x2c\x33\x50\x2d\x1a\x87\x3c\x66\xaa\xb9\x19\x57\x13\x10\x2d\xba\xb3\xb7\xe8\x75\x4d\x61\xe9\x75\x2d\x81\xec\xa9\x45\x61\x57\xfd\x99\xcf\x85\x9e\x44\xe8\x1b\xc2\x22\x78\x78\x85\xbe\x51\xd9\x94\x40\xbf\xf4\x51\x27\x4d\xcf\x9d\x6c\xba\x74\x8a\xe9\xf2\x54\xdf\x8e\x50\xa3\x5a\xce\xe4\x99\xee\xd8\xf9\xd2\x52\xd9\x03\x73\x67\x64\x27\xdb\x52\x8b\xd5\x9c\x14\x48\x41\xd7\x1d\x41\xeb\x9b\xe9\xc8\x29\x91\xea\x1a\x1e\xd4\x90\x3f\x2c\xc6\xdd\xb9\x4a\x66\x94\xde\xca\x25\x32\x38\x92\x4e\xb8\xdb\x64\xc2\x05\x96\x38\xa0\x54\x8d\x06\x7d\x1b\xbc\x1b\x72\x4e\xaf\x80\x42\x58\xf1\xef\xd7\x74\x0e\x5f\xa6\x53\xd8\xe4\x62\x36\xcb\x6f\xb3\x59\xbe\x17\x2f\x3f\xa6\x0b\xe7\x52\x19\x3d\xcc\x96\xd5\xad\x50\xdb\xf5\x4e\xa5\xea\xfe\x3e\xfa\xc8\xd4\x85\x12\x25\xff\xde\xe9\xb5\x6a\x72\x2d\x59\xc4\x5b\x74\xaa\x97\x96\x7e\xd9\xb4\x4f\xff\xa8\x6d\x4a\x84\x8d\x78\x80\xf4\xde\x77\xa6\xf8\x78\xac\x37\x96\xe4\x94\x95\xbf\x86\xc5\x18\x54\x3f\xf8\x73\xe6\xfe\x59\xe1\xfe\x59\x5a\x2e\xdd\xce\xe7\xe9\x82\x7e\x7a\x32\x24\xdc\xea\x4f\x04\x0a\x13\x2a\xd1\x04\x04\xb8\x7d\xef\xa6\xfe\x3a\x86\xbf\x83\xe0\x64\xbb\x1e\xe2\x2c\xa4\x24\xfc\xd2\x0f\x16\x47\x52\x35\x21\xf2\x15\x3a\x6d\x11\x96\xd3\xef\x83\x41\x7a\x0a\x7a\xce\xa0\xcc\x89\x18\x95\xf6\x5d\x73\x18\x1a\x36\x35\x77\x8c\x8c\x45\x55\x31\xd4\x2e\x8e\xa2\x4a\x71\xd5\x3c\x03\x5c\x75\xd3\x67\xb8\xaf\x49\x90\x57\x4b\x8e\xfc\x8b\x0c\xf5\x4e\xfd\xfd\xf5\xf7\xc0\x43\x69\x61\xdd\x8b\xeb\x7b\x6e\xaf\x9b\xd4\x42\xae\xfa\x6f\x03\x9b\x71\xb9\x78\x4b\x92\x06\x1a\xe1\x08\x82\xe4\x74\xd4\x26\x5f\xe8\x82\x2d\xf9\xbd\x1f\x9c\xfd\x14\x20\xc1\xd3\xb2\x16\x53\x3e\x0e\x10\x16\x04\x9f\x51\x3c\x04\x4a\x21\x1a\x3e\xb6\xea\xf1\x4c\x11\x45\x8d\xd5\x7a\xd5\xd3\xb3\xfc\x31\xd9\x43\x79\x18\x4f\x81\xd9\x26\x70\xdd\x5c\x1f\x6d\xed\xed\xcd\x36\xd9\x05\xd8\x33\xf2\x52\x48\xb9\x84\x2c\x33\x47\x44\x4e\x49\xd1\x79\x30\xf8\x56\xc7\x45\xbe\x69\xce\x31\x28\xad\xcc\x97\x7d\x4b\x03\xd7\x56\xb9\x3c\xce\xef\xd3\x74\xbe\x5c\xd2\xd7\x9e\xb5\x7c\x4c\x68\x8e\x93\x9e\xce\x4d\x51\x9a\xbc\x1e\x7c\x48\xbc\x72\x3f\xbb\xfa\x80\x86\x23\x86\xd1\xfc\xf9\xc7\x0e\x73\xb7\xf6\xa3\x88\xb9\xbd\x23\x93\x9b\x0d\xa2\x81\x2e\x1d\xdd\x19\xc7\x66\x99\x64\xe4\xd5\x4c\xed\x3b\x94\xa1\x65\xfb\xc1\xbb\x8e\x14\xe6\xf6\x6b\x04\x6a\x86\xf5\x94\x67\x0d\x47\x25\x9b\xf5\x4a\x06\xc8\x71\x20\x48\xa7\x74\xe7\x22\xf3\xc6\x54\x34\xa6\x6b\xb1\x93\x3b\xec\x2a\x1f\x8d\xee\x6e\x49\x52\xb4\x6e\xe4\x63\x41\xd1\x61\x45\xff\xe6\xf2\xbc\x85\x00\xb1\xa0\xb7\x2f\x46\x84\x09\x96\x68\x0a\x11\xc1\x3b\x53\xc0\x76\x98\xcb\x34\xf8\x0d\xcb\x4f\xda\x1f\x87\x00\x13\x2c\x6f\x13\x9f\x5f\x44\xf4\x63\x41\x0f\x28\xf6\x37\x97\xe7\x0d\x91\x8f\x05\x3d\xfa\xb8\x0b\x50\xf7\x00\x0a\xa2\x43\x89\xfc\x65\xee\x90\x23\xf8\x85\xd3\x47\x1f\xfe\x11\xbe\xe3\x82\x28\x40\x6a\x22\x40\x4e\x38\xdd\x9d\x0e\xf6\xcb\x9a\x4c\x89\x0f\x99\x6f\xd7\xb9\x6b\x0e\x45\xf2\x71\xdc\x16\xe3\x38\x7a\x69\x8a\x49\x76\x88\xda\x14\xab\xa4\x8d\x38\xc5\x48\x5e\x90\x3a\x14\xb3\xf1\x8e\xe4\x70\xdf\xa4\x6b\xa9\xec\xf7\xe8\x99\x02\xda\xdb\xe2\x86\x67\x3e\xaf\xc9\x79\x8e\xd9\xb8\xd5\xdd\xe1\x92\x5f\xdb\x3b\x99\xb4\x3b\x6d\xd4\x2e\x58\xcc\x8d\x9c\x27\x5f\xb4\x74\x90\xed\xfc\x8b\xc8\xf4\x93\x6c\x7f\x9e\x35\x1b\xf8\xf3\x6c\x4b\x03\x9d\x22\x86\xb0\xbb\x82\xd2\x79\x9c\x4a\xa7\x75\xe7\x5c\x7b\xe4\x48\xd4\x77\x69\xb3\xc4\xf1\xa3\xcf\xd0\x23\x1c\x02\x4a\x3e\x3b\x40\x94\x7c\x01\x4a\x26\x9c\xef\x6a\x07\xdd\x40\xca\xce\xb4\xd0\xa3\xe8\x24\xa3\xb8\x5d\x8c\xc2\x95\xca\x33\xa9\x3f\x68\xbb\xb7\xda\xee\xbc\x30\x3b\x98\x0c\x8f\x9e\xa5\xe9\x90\xc6\x42\x40\x74\xf4\xaa\x66\xe3\x58\x43\xd7\xbf\xa7\x96\x2f\x4e\x59\x9d\xc0\xef\x01\x1f\xff\x82\xcd\x07\xb2\x86\xb6\xbf\x65\xa6\x2f\x4e\xdc\x3f\xf8\xe3\xd1\xeb\xfa\x07\x7f\x5c\x43\xd2\x7f\xf0\xc7\x97\xa4\xa6\x8e\xda\x41\x95\x35\xfa\xfd\xe6\xaa\x26\x11\xfb\xd8\x8b\x1a\x8a\x59\x34\xc5\xe2\xcb\x41\xc5\xff\x3c\x73\xaa\x4d\x65\x99\xb6\x3c\x7e\x1d\xf8\x98\x1f\x96\x06\x7c\xcc\x5b\xc4\x9f\x8f\xf9\xa1\xc4\x7e\x9f\xc7\xfc\x73\xcc\xc6\x31\x1e\x83\x3f\xe8\x5b\x0d\xfc\x41\xbf\xa5\x41\x34\x98\x12\x86\x24\x30\xbd\xe8\x9a\xbe\x10\x6c\xeb\xe2\x39\x77\xb4\x1f\x28\xc7\xea\x6f\x7f\x31\xde\xd3\xe6\x13\xbd\xf3\x89\xb0\xab\xdc\x45\x47\x9e\xa0\x79\xfb\x29\x61\xb7\xc5\x98\x0e\x25\x65\xa0\xb5\x15\xc2\x0f\x87\xaf\x10\x7e\x58\x4d\x21\xfc\x70\x78\x0a\xed\x25\xa9\xbf\x0d\x95\xbf\xb4\x75\x19\xf8\x5c\xde\xd2\xa0\xf8\xd4\x6d\x67\x39\x22\x9c\x40\xf8\x65\xb9\xb0\x4b\xa7\x73\xfe\x11\x9b\x29\x15\xa4\xdf\x06\xce\x3f\x57\x3b\x94\xb5\x8f\x9e\xfb\x21\xf4\x01\x04\x3d\xff\xcc\xd9\x11\xf5\xdc\xdb\xe3\x0f\x3b\xa7\x94\xdf\x1f\x42\xd0\x13\x47\x5c\x21\x4f\x1a\x1c\x7d\xc0\x43\x4e\x29\xe4\x4c\xeb\x5e\xcf\x8e\x59\xe0\xdf\x2d\x1c\x72\x44\xbf\xe4\xf6\xa1\x48\xb0\xd3\x32\x63\xe5\xef\x4e\x8f\x38\x57\xcf\xfb\x8e\x79\xc1\xbe\xc0\x08\xc7\x54\x59\xbe\x6d\xfe\x8e\x72\xd9\x8a\x67\xb1\xfa\x6f\x79\xcb\xf0\xb2\x8d\x33\x2d\xde\x2f\x61\xa6\x79\x12\xfd\x5a\x30\xd3\x17\x01\x82\x16\xdf\x04\x33\xb0\x27\x85\xa0\xbb\x02\x41\x73\x67\x76\x00\x82\xe6\x8f\x32\x83\xa0\x0b\x47\x76\x06\x82\xd6\xbe\x7b\x55\x78\x58\x81\x18\x17\xbe\x1d\x25\xc4\x58\xb8\xef\x21\xc6\x32\xc4\xe8\x08\xcb\x4b\x83\x18\x8b\xa1\x7a\x88\xb1\x7d\x8b\xad\x02\x88\xeb\x6c\x02\x6d\x00\x44\xd7\x5a\x5f\x0f\x40\x74\xf4\xe8\x01\xc4\x0d\x02\x88\x2d\xe2\xec\x01\x44\x7f\xf7\xb7\x30\x3b\xe2\xbb\xbf\xfd\x22\x70\x45\x99\x67\x05\x10\x17\x05\x9f\x07\x10\xab\x3d\x6c\x3a\xfa\x46\x00\xb1\x26\x80\x07\x10\x17\xe6\x6b\x28\x50\x41\xb1\xaa\x1a\xd4\x01\xc4\x9a\x00\x1e\x40\xdc\x5e\xec\x97\x00\x44\x63\xe4\x3d\x80\xb8\x85\xc8\x1b\x00\xc4\x5a\xf0\x3d\x80\x58\xed\x67\x0d\x1d\xec\x17\x2d\xcd\x00\x62\x4d\x11\x0f\x20\xda\x3b\xda\x86\x36\x2e\x00\xd1\xbe\x5c\x5e\x90\x3a\x07\x0d\x20\xd6\x14\xa8\x02\x88\x55\x39\x3d\x80\xe8\x01\xc4\x76\x06\xfe\x3c\xdb\xd2\x60\xaf\x00\x62\x75\x7d\x9b\x01\xc4\x5a\x9a\xf0\x00\xa2\xb9\xab\x7d\xa4\xec\xb6\x00\xa2\x45\x6a\x0f\x20\x36\x76\xb6\x6f\x55\xdd\x00\xa2\x4b\x57\x0f\x20\x36\xf4\xb6\x6f\x69\x1b\x00\x44\x97\xb6\x1e\x40\x74\x74\xb4\x6f\x5d\xed\x00\xa2\x4b\x52\x0f\x20\x6e\xb7\xac\x59\x06\x10\x6d\x02\x7a\x00\x71\x7b\x65\x65\x15\x40\xb4\x57\x96\x1e\x40\x5c\xc5\x00\xb5\xd7\x60\x09\x40\xb4\xc6\xdf\x03\x88\x1e\x40\x6c\x67\xe0\x0f\xfa\x2d\x0d\x0e\x0a\x40\x34\x5d\xec\x39\x00\x44\xe3\x45\xa1\x07\x10\x4d\x5d\x6c\x5f\x21\x23\x80\xe8\x50\xc8\x03\x88\xc8\x03\x88\xad\x0c\x7c\x2e\x6f\x69\xb0\x4f\x00\xb1\x48\x0c\x36\x00\x71\x91\x0a\x3c\x80\x58\xb2\xde\x68\xd0\xeb\x00\x62\x2d\xea\x1e\x40\xdc\x78\xd0\x2b\x00\x62\x3d\xe4\x1e\x40\xdc\xec\xd9\xd1\x05\x20\xd6\xa2\xef\x01\x44\x0f\x20\xae\x04\x20\x5e\x01\x16\xe1\xe4\x88\xf9\xc3\xff\xc5\x20\x5a\xfd\x0f\x2b\x05\xc8\x98\xa6\x10\xcc\x57\xc2\x29\xca\x44\x5b\x13\xa0\x92\xab\xbe\x2b\x48\x51\x66\xcf\xdb\x3e\xa3\x98\x3e\xa9\xf3\xcf\x74\x56\xa0\xba\x07\xd9\x7c\xd9\x2a\x97\x28\x17\x5f\x59\xcc\xfc\xb9\x4c\x26\xdf\xf5\xe3\x0c\x0c\x2e\xa5\x33\xf3\x56\x07\x30\x40\x41\x80\x02\x01\xa1\x3e\x27\xa2\x60\x4a\x1e\x74\xb0\x82\x19\x9f\xc5\x14\x8b\xdd\xc2\x94\x99\xeb\x15\x94\xb2\xf0\xfa\x28\x49\xca\xdc\x7b\x0f\x52\x96\x41\x4a\x7b\x54\x5e\x1a\x47\x99\x8f\xd4\x63\x94\x9b\x6c\xb1\x55\xd0\x72\xe5\x5d\xac\x0d\x65\xe9\xc8\x03\xeb\x41\x96\xf6\x0e\x3d\x63\xb9\x41\xc6\xb2\x39\xcc\x1e\xb1\xf4\xb7\x9b\x0b\xb3\x23\xbe\xdd\xdc\x2f\xe4\x97\x15\x7f\x56\xc0\xb2\x28\x03\x3d\x5f\x59\xed\x61\xb3\xa1\x37\xd2\x95\xd5\xe8\x7b\xb8\x72\x61\xbe\x46\xf8\x2b\x98\xd9\xb2\x00\x75\xb4\xb2\x1a\x7d\x4f\x56\x6e\x2b\xf0\x4b\x5c\xa5\x29\xec\x1e\xab\xdc\x78\xd8\x0d\x50\x65\x35\xf2\x9e\xa9\xac\xf6\xb3\x86\x08\xb6\xfb\x96\x66\xa2\xb2\x2a\x87\x07\x2a\xed\x1d\x6d\x5e\x18\x17\x4e\x69\x5d\x28\x2f\x48\x9a\x83\xa6\x29\xab\x02\x54\x61\xca\x65\x29\x3d\x4a\xe9\x51\xca\x76\x06\xfe\xdc\xda\xd2\x60\xaf\x28\xe5\xf2\xea\x36\x83\x94\xd5\x0c\xe1\x39\x4a\x73\x57\x7b\x48\xd6\x6d\x31\x4a\xa3\xcc\x1e\xa2\x6c\xec\x6c\xcf\x92\xba\x19\x4a\xbb\xa8\x9e\xa0\x6c\xe8\x6d\xcf\xba\x36\x00\x94\x76\x61\x3d\x3e\xe9\xe8\x68\xcf\xa2\xda\xe9\x49\xbb\x9e\x9e\x9d\xdc\x66\x1d\xb3\x4c\x4e\x5a\xb4\xf3\xe0\xe4\xb6\x8a\xc8\x2a\x36\x69\xad\x23\x3d\x35\xb9\x8a\x01\x6a\x2b\xc0\x12\x33\x69\x0b\xbe\x47\x26\x3d\x32\xd9\xce\xc0\x1f\xe8\x5b\x1a\x1c\x14\x32\x59\xbf\xbe\x73\x00\x93\xa6\xbb\x40\xcf\x4b\x9a\xba\xd8\xb6\x3c\x46\x5a\xd2\x2e\x8f\x87\x25\x91\x87\x25\x5b\x19\xf8\x2c\xde\xd2\x60\x9f\xb0\x64\x96\x15\x6c\xa8\x64\x91\x07\x3c\x29\x59\xb2\xde\x60\xc4\xeb\x9c\x64\x35\xe4\x1e\x93\xdc\x70\xc4\x2b\x90\x64\x2d\xde\x9e\x91\xdc\xe4\x49\xd1\x45\x48\x56\x43\xef\x01\x49\x0f\x48\xb6\x02\x24\xdf\x5e\x7c\x44\x23\x2e\x50\x06\x1b\x1c\x2d\x1f\x29\x79\x2c\x42\x68\xfe\x56\x5f\x72\x94\x00\x29\xf1\x18\x90\x82\xe9\x8c\xe2\xa6\xad\x6b\x7f\x60\x23\x9e\x11\x13\x0f\xf2\xf6\xe2\xe3\xce\x88\x46\x3c\x23\xdb\xa0\x19\x57\xbf\xc4\x4f\x1c\x49\x45\x4e\xbe\x0d\x59\xba\xa9\xc7\x33\xd2\xb9\x4a\xde\xb9\xb9\x3c\x77\xde\xc7\x3f\x17\x47\xd3\x8e\x61\x01\xb8\xee\xb8\xc1\xd9\x6c\x96\xdd\xe6\xb3\x2c\x40\x21\xa7\xb2\x1f\xbc\xfe\x31\x40\x82\xdf\xcb\x7e\xf0\xd7\x60\x90\xb9\xff\x29\x6d\x7b\x9d\x35\x7d\x7a\xea\x75\xf3\x67\x6d\x6f\x34\x9b\x85\xeb\xf4\x98\x5f\x3e\x58\xa7\x47\xe9\xa1\xba\x52\xef\x1b\x03\xe2\xca\xb4\xfe\x05\xc5\x21\xa0\xcf\x5c\x91\x11\x09\x71\xfd\x72\xe0\x88\xf6\x25\x9c\x94\xc8\x12\xe8\xa8\x39\xd0\xb1\x74\x52\xf6\x1b\xdc\x6a\x9e\xb7\x34\x4a\x27\x82\xea\xde\x54\xd6\xac\x93\xc8\xd8\x79\xab\x23\x70\x05\x74\xb4\xa8\x56\x59\xb9\xd5\x2c\x69\x95\xc4\xe9\x56\xc7\x69\xab\xf4\x77\xb9\xae\x6e\xe1\xfb\x8d\x56\xc4\xe9\x77\xa2\xd9\x7a\x2e\x6f\x7a\xf5\x7c\x64\x0a\x04\x36\xdc\xa5\xf9\xe5\xb2\xf4\xea\x01\x2c\x97\x92\x54\xe5\xf5\x41\x4a\x2f\xef\x71\x41\x94\xbd\xcb\x56\x40\xd9\xb3\x03\x9a\xf2\xe7\x7c\xec\xa7\xfa\xa2\x19\x25\x0c\x58\x3c\x3d\xce\x35\x71\xce\xc7\xe5\xb5\x40\xf9\x78\x9f\x6b\x40\x7b\x93\xcd\x7d\xed\xc9\x33\xe6\x7c\x2b\x27\xdc\x07\x22\x16\x4f\x87\xba\x0a\x4f\x4f\x19\xda\x9f\x4c\xe9\xd2\x31\xa8\xec\xf8\x79\xfa\xee\x9a\x67\xa1\x0d\x2c\xd2\xd2\xaf\xbd\x6e\x6a\xd6\xeb\x4e\xd4\x94\x0e\x4e\xf2\x8a\xfd\xff\x01\x00\x00\xff\xff\x2d\xba\x82\x7b\xb9\xc4\x00\x00")

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

	info := bindataFileInfo{name: "assets/tmpl/config.tmpl", size: 50361, mode: os.FileMode(436), modTime: time.Unix(1485615898, 0)}
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

var _assetsTmplTwitter_setupTmpl = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\xc4\x54\x31\x6f\xdb\x3c\x10\xdd\xfd\x2b\xee\xbb\xdd\x26\x92\xf1\x03\xa5\x25\xed\x14\x04\x2d\x10\x2f\x9d\x82\x93\x78\xb6\xd8\x50\x24\x41\x9e\x1c\x08\x82\xff\x7b\x21\xc9\x4e\x1c\xc7\x0e\x32\x14\xa8\x07\x4b\xc4\xbb\xf7\xf8\xf8\xee\xc4\x61\x30\xbc\xb1\x9e\x01\xe5\xc5\x8a\x70\x7a\xca\x2c\x5d\xc4\xfd\x7e\xa1\xff\xfb\xf6\xe3\x6e\xfd\xeb\xe7\x77\x68\xa4\x75\xe5\x42\x8f\x0f\x70\xe4\xb7\x05\xb2\xc7\x72\x31\x0c\xc2\x6d\x74\x24\x0c\xd8\x30\x19\x4e\x08\xab\x91\x58\x05\xd3\x97\x0b\x00\x00\x6d\xec\x0e\x6a\x47\x39\x17\x58\x07\x2f\x64\x3d\x27\x9c\xb1\xf1\x77\x2a\xe1\x69\x57\xd1\x41\xe2\x88\x9f\xf2\x7f\x77\x6d\x15\x24\x05\x7f\xc2\x9f\x6a\x9a\x9b\x72\x3d\x9b\x87\xc7\xd1\xbc\x56\xcd\xcd\x79\xc9\xed\x51\xc5\x31\x19\x2c\xd7\x8d\xcd\x60\x33\x48\xc3\xb0\xb1\x29\x0b\x64\xe1\x08\x12\xa0\xcb\x0c\x0f\x7d\x15\x44\xab\xe6\xf6\x4c\x65\x13\x52\x0b\xd6\x14\x38\x67\x04\x54\x8b\x0d\xbe\x40\x35\xad\xd5\x21\x41\x85\xd0\xb2\x34\xc1\x14\x18\x43\x16\x04\xf6\xb5\xf4\x91\x0b\x6c\x3b\x27\x36\x52\x12\x35\x2a\x2d\x0d\x09\x61\xa9\xa7\xc5\xd9\x4e\xd6\xc7\x4e\x60\x04\x5e\x37\x9b\x25\x72\x57\xb5\x56\x10\x76\xe4\x3a\x2e\xf0\x91\x76\x8c\xc7\xa3\x55\xe2\xa1\x12\xbf\x8c\xc9\xb6\x94\xfa\xe9\xdd\x6d\x4f\xd2\xd2\xca\xd8\xdd\x69\xf8\x76\x03\xab\x07\xce\x99\xb6\x7c\x25\x74\x72\x9c\x04\xa6\xff\xa5\x21\xbf\x1d\xbb\x37\x0c\x6f\xa4\x0f\x92\xec\xcd\x17\xa4\xac\xdf\x84\xf3\x36\x12\x34\x89\x37\x05\x36\x22\x31\xff\xaf\x14\xc5\x98\x57\x87\x4c\x57\x75\x68\x11\x84\xd2\x96\xa5\xc0\xa7\xca\x91\x7f\xc6\x72\x1d\xa0\x4e\x4c\xc2\xaa\x6e\xb8\x7e\x86\x3e\x74\x09\xc2\x8b\x87\xe3\x38\x50\x8c\xce\xd6\x34\x76\x09\xb6\x61\x6c\x6f\xc3\x89\xb5\xa2\xab\x99\x9c\x1a\x16\xaa\x1c\x2f\x13\xe7\x18\x7c\xb6\x3b\x3e\x37\x3c\xe1\xef\x8a\x61\xa6\x64\x49\x36\xb2\x39\xac\xaa\x90\x0c\x27\x36\x67\xf4\x59\xe2\xed\x53\xf9\x88\xa5\xcb\xc0\x0c\x9a\xf2\x2e\xf8\xdc\xb5\x9c\xe0\x9e\x7b\xad\xc4\x7c\x5e\x7d\x15\x84\x4f\xe7\x2d\x52\xce\x2f\x21\x19\x04\x4f\x2d\x17\xef\x6f\x89\x55\x7d\xb0\xf0\xf4\xcc\xfd\xeb\x4c\x0e\xc3\xea\x68\xed\x9e\xfb\xfd\xfe\xc2\xb9\x5f\xf7\xbd\x6a\x5b\xab\x6b\xc7\xff\x7a\x2e\x8f\x5c\x27\x96\x7f\x1d\x4d\x9e\x5c\x5c\x4a\x67\xf6\xf7\x57\x03\xd2\xea\xc2\x44\x69\x35\x0d\xe2\xc5\x99\x3f\xbc\x6a\x35\xd3\xb4\x9a\xef\xfa\xe3\x77\xfc\x27\x00\x00\xff\xff\x29\xfb\xad\xe7\x1f\x06\x00\x00")

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

	info := bindataFileInfo{name: "assets/tmpl/twitter_setup.tmpl", size: 1567, mode: os.FileMode(436), modTime: time.Unix(1485616185, 0)}
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

