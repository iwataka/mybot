// Code generated by go-bindata.
// sources:
// assets/css/custom.css
// assets/tmpl/config.tmpl
// assets/tmpl/header.tmpl
// assets/tmpl/index.tmpl
// assets/tmpl/log.tmpl
// assets/tmpl/navbar.tmpl
// assets/tmpl/status.tmpl
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

	info := bindataFileInfo{name: "assets/css/custom.css", size: 753, mode: os.FileMode(436), modTime: time.Unix(1483973425, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _assetsTmplConfigTmpl = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\xec\x9d\x6d\x6f\xe3\x36\x12\xc7\xdf\xe7\x53\xf0\x74\x45\xd3\xa2\x1b\xf9\xda\x1e\xee\x45\xd7\x36\xd0\xdb\xed\x62\xf7\x90\xdd\xe6\x92\xec\x01\x87\xc3\x21\xa0\xad\xb1\xcc\x2e\x4d\xfa\x28\x2a\x89\x61\xe4\xbb\x1f\xa8\x27\xcb\x12\x49\xc9\x8e\x1f\xe4\x2c\xf3\xa2\x1b\xdb\x1c\x6a\x38\xff\xe1\x88\xa4\xf5\x6b\x96\xcb\x00\x26\x84\x01\xf2\xc6\x9c\x4d\x48\xe8\x3d\x3d\x9d\xf5\xff\xf4\xf6\xf7\x37\xb7\xff\xbe\xfa\x0d\x4d\xe5\x8c\x0e\xcf\xfa\xea\x1f\x44\x31\x0b\x07\x1e\x30\x6f\x78\xb6\x5c\x4a\x98\xcd\x29\x96\x80\xbc\x29\xe0\x00\x84\x87\x7c\x65\x38\xe2\xc1\x62\x78\x86\x10\x42\xfd\x68\x2c\xc8\x5c\xa6\x2f\xd4\xcf\x24\x66\x63\x49\x38\x43\x01\x50\x90\x70\xcd\x1f\xbe\x1b\xc5\x52\x72\xf6\x0a\xcd\x78\x80\xe9\xf7\x68\x59\xb4\x55\x3f\xf7\x58\x20\x29\xd0\x00\xa5\xad\xfc\x39\x16\xc0\xe4\x27\x1e\x40\xe9\xd7\xd7\x6b\x26\x64\x82\xbe\x93\xc2\x1f\x53\x1c\x45\x9f\xf0\x0c\xd0\x60\x30\x40\xe7\xe9\xf5\x82\xf3\xea\x05\x12\xa7\xb8\x40\xdf\xa9\x2b\x11\x34\x40\x7f\x79\x8d\x08\xea\x23\xd5\xc3\x94\xd0\x40\xf5\x1f\xf9\x14\x58\x28\xa7\xaf\x11\xf9\xe1\x07\x5d\x07\xb9\xa7\x89\x05\x1a\xac\x1b\xff\x87\xfc\x57\x6b\xa0\xfc\x4c\x5a\xf9\x72\x31\xcf\xbc\x9c\x92\x20\x00\x76\x8e\xbe\xfd\x36\xed\xcb\x67\x3c\x80\xd5\x20\x3e\x7c\xba\xfa\x7c\xab\x1d\x42\xfe\x93\x5a\xdd\x63\x1a\x03\x1a\xa0\xf3\x09\xa6\x11\x9c\xbf\x36\x36\x5f\x8f\x13\x3a\xb7\x34\xcd\x04\x20\x8c\x81\x78\x7f\xfb\xf1\x52\x35\x7f\x9b\x04\xd5\x60\xf4\x54\x7b\x77\xfd\x9d\x27\x04\x34\x82\xaf\x4a\x0e\x29\xe2\x4d\xd4\xc8\x73\x76\x33\x51\xae\xe1\x1e\x84\xdc\x56\x94\xb3\xf5\xdf\xfa\xbd\xf2\x0c\xee\x07\xe4\x1e\x25\x2e\x0e\x54\xa1\x90\x98\x30\x10\xde\x6a\x76\x97\x2b\x02\xc3\xf7\x23\x9c\x55\x84\xa2\x41\xb9\x83\x3f\xe2\xd9\x88\x4b\xc1\x59\xa9\x83\xa4\xcd\xf4\xc7\xe1\x9b\xa4\x0a\xc5\x02\xab\x5a\xd1\xef\x4d\x7f\xac\x36\xf9\x29\xef\x85\x02\x0e\xbc\xa2\x3d\xa0\x05\x8f\x05\xe2\x0f\x0c\x8d\xb8\x3c\x8f\x50\xc4\x63\x31\x86\xe8\x15\x9a\x10\x2a\x41\x44\x08\xb3\x00\xe1\xa4\x04\x45\xfd\xde\xf4\xa7\x4a\xbf\x13\x2e\x66\x88\x04\x03\x8f\xdf\x83\x78\x10\x44\x82\x97\xb5\x1e\x78\xbd\xb4\x34\xf6\x3c\x34\x03\x39\xe5\xc1\xc0\x9b\xf3\x48\x7a\x08\xd8\x58\x65\xcc\xc0\x9b\xc5\x54\x92\x39\x16\xb2\xa7\xba\xb9\x08\xb0\xc4\xde\xb0\x9f\xbc\xa8\x5c\x86\xb0\x79\x2c\x55\x9e\xcf\xd6\xae\x94\x76\x13\xc5\xa3\x19\x91\x1e\x4a\xb2\x66\xe0\xdd\xe0\x7b\xf0\xf2\xd1\x8e\x24\x43\x23\xc9\x2e\x68\x98\xfc\x33\x17\x64\x86\xc5\xc2\x43\x91\x5c\x50\x18\x78\x01\x89\xe6\x14\x2f\x7e\x21\x8c\x12\x06\xde\xf0\x4c\x33\xba\xea\x70\x26\x84\x42\x69\x4c\x21\x48\x73\x77\xd5\xdc\xc9\x06\xa2\xf5\xfb\x2d\x7f\x60\x94\xe3\xa0\xc9\xf7\x4a\x68\xb2\x70\x6d\xee\x78\x0b\x31\x76\x32\xae\xcf\xf3\xcd\x47\x55\xeb\x53\xf9\xee\x21\x86\x67\xca\xd3\xc5\x88\x4b\x3f\xbb\xef\xb6\x72\xb1\x9a\x53\xfd\x5e\x40\xee\x4b\x31\x5b\x2e\xc9\x04\xf9\x1f\x21\x8a\x70\x08\x4f\x4f\xda\xd9\x87\x29\x08\x89\x92\xff\x5e\x04\x98\x85\x6a\x1e\x2f\x97\x2b\xa3\xac\xcf\x55\x97\xc0\x02\x35\x93\x57\x9d\x4d\x7f\x1e\xde\x3e\x10\x29\x41\xa0\x5b\x32\x03\xe5\xa7\x9a\x52\x3f\x0f\xb5\xd7\x93\x78\x44\xe1\x42\x40\x34\xe7\x2c\x22\xf7\xb5\x21\x25\x9f\xaf\x35\x46\xa9\x49\x24\x05\x99\x43\x90\xbd\x1a\x71\x11\x80\x80\x40\x17\x60\xa9\x16\x21\xf5\xf7\xd3\xcf\x84\xfe\x83\xf4\xc3\x60\x18\x8d\x05\x00\x4b\x14\x89\xfa\x3d\x69\xe8\x26\x6f\x0d\x8f\x63\x1a\x07\x80\x04\xcc\x29\x69\x63\x40\x58\x66\x20\x5b\x34\x1e\xf3\x98\xc9\xe6\x66\x5c\x4e\x41\xb4\xe8\xce\xdc\xa2\xdf\xd3\x85\xa5\xdf\x33\x04\xb2\x2f\x57\x0b\xbb\xea\xcf\x72\x29\x54\x12\xa1\x6f\x08\x0b\xe0\xf1\x15\xfa\x46\x66\x29\x81\x7e\x19\x20\x3f\x2d\xcf\x7e\x96\x2e\x7e\x91\x2e\x4f\xf5\xdb\x11\x6a\x54\xcb\x5a\x3c\xd3\x3b\x76\x3e\xb5\x64\x76\xc1\xdc\x99\xc8\xcf\x6e\xa9\xc5\x6c\x4e\x16\x48\x5e\xcf\x1e\x41\xe3\x87\xe9\xc8\x29\x89\xe4\x2d\x3c\xca\x11\x7f\x5c\x8d\xdb\xbf\x49\x32\x4a\xdd\xca\x23\xa4\x71\x24\x4d\xb8\xbb\x24\xe1\x3c\x43\x1c\x50\xaa\x46\x83\xbe\x0d\xde\x8d\x38\xa7\x37\x40\x61\x5c\xf1\xef\xb7\x34\x87\xaf\xd3\x14\xd6\xb9\x98\x65\xf9\x5d\x96\xe5\x47\xf1\xf2\x43\x3a\x71\xae\xa5\xd6\xc3\x6c\x5a\xdd\x09\xb9\x5f\xef\x64\xaa\xee\xef\x93\x0f\x4c\x5e\x49\x51\xf2\xef\x8d\x9a\xab\x3a\xd7\x92\x49\xbc\x47\xa7\xfa\xe9\xd2\x2f\x4b\xfb\xf4\x45\xed\xa6\x44\xd8\x84\x7b\x48\xdd\xfb\x2e\x24\x0f\x43\x75\x63\x49\x76\x59\xf9\x7b\x58\x84\x20\x07\xde\x9f\x33\xf7\x2f\x0a\xf7\x2f\xd2\xe5\xd2\xdd\x72\x99\x4e\xe8\xa7\x27\x4d\xc1\xad\xfe\x04\x20\x31\xa1\x11\x9a\x82\x00\xbb\xef\xbd\xd4\x5f\xcb\xf0\x0f\x10\x9c\xec\xae\x87\x38\x1b\x53\x32\xfe\x32\xf0\x56\x5b\x52\x39\x25\xd1\x2b\x74\xde\x22\x2c\xe7\xdf\x7b\xc3\x74\x17\xf4\x9c\x41\xe9\x0b\x31\x2a\xdd\x77\xf5\x61\x68\xb8\xa9\xd9\x63\xa4\x5d\x54\x15\x43\xed\xe1\x20\xa8\x2c\xae\x9a\x33\xc0\xb6\x6e\xfa\x04\x0f\x35\x09\xf2\xd5\x92\xa5\xfe\x22\xcd\x7a\xa7\xfe\xf9\xf6\xf7\xc0\xae\xb4\x30\xde\x8b\xeb\xf7\xdc\x7e\x2f\x59\x0b\xd9\xd6\x7f\x3b\xb8\x19\x97\x17\x6f\x49\xd1\x40\x13\x1c\x80\x97\xec\x8e\xda\xd4\x0b\xb5\x60\x4b\x7e\x1f\x78\x17\x3f\x7a\x48\xf0\x74\x59\x8b\x29\x0f\x3d\x84\x05\xc1\x17\x14\x8f\x80\x52\x08\x46\x8b\x56\x3d\x5e\x48\x22\xa9\x76\xb5\x5e\xf5\xf4\x22\xbf\x4c\x76\x51\x3e\x8e\x67\xc0\x4c\x09\x5c\x37\x57\x5b\x5b\x73\x7b\xbd\x4d\x76\x00\xf6\x8c\xba\x34\xa6\x3c\x82\xac\x32\x07\x24\x9a\x91\xa2\x73\x6f\xf8\xad\x8a\x4b\xf4\xba\xb9\xc6\xa0\x74\x65\xbe\xee\x5b\x1a\xb8\xb6\xca\xe5\x71\x7e\x9b\x96\xf3\xf5\x25\x7d\xed\x5a\xeb\xdb\x84\xe6\x38\xa9\x74\x6e\x8a\xd2\xf4\xe7\xe1\xbb\xc4\x2b\xfb\xb5\xab\x17\x68\xd8\x62\x68\xcd\x9f\xbf\xed\xd0\x77\x6b\xde\x8a\xe8\xdb\x5b\x2a\xb9\xde\x20\x18\xaa\xa5\xa3\xbd\xe2\x98\x2c\x93\x8a\xbc\x99\xa9\xf9\x0e\xa5\x69\xd9\x7e\xf0\xb6\x2d\x85\xbe\xfd\x16\x81\x9a\x63\x95\xf2\xac\x61\xab\x64\xb2\xde\xc8\x00\x59\x36\x04\x69\x4a\xfb\x57\x99\x37\xba\x45\x63\x3a\x17\xfd\xdc\x61\xdb\xf2\x51\xeb\xee\x9e\x24\x45\xdb\x46\x3e\x16\x14\x75\x2b\xfa\x9f\xaf\x2f\x5b\x08\x10\x0b\x7a\xf7\x62\x44\x98\xe2\x08\xcd\x20\x20\xf8\x60\x0a\x98\x36\x73\x99\x06\xef\x71\xf4\x51\xf9\x63\x11\x60\x8a\xa3\xbb\xc4\xe7\x17\x11\xfd\x58\xd0\x0e\xc5\xfe\xf3\xf5\x65\x43\xe4\x63\x41\x4f\x3e\xee\x02\xe4\x03\x80\x84\xa0\x2b\x91\xbf\xce\x1d\xb2\x04\xbf\x70\xfa\xe4\xc3\x3f\xc1\xf7\x5c\x10\x09\x48\x4e\x05\x44\x53\x4e\x0f\xa7\x83\xf9\xb0\x26\x53\xe2\x5d\xe6\xdb\x6d\xee\x9a\x45\x91\x7c\x1c\x77\xc5\x38\x4e\x5e\x9a\x22\xc9\xba\xa8\x4d\x31\x4b\xda\x88\x53\x8c\xe4\x05\xa9\x43\x31\x0b\x0f\x24\x87\xfd\x24\x5d\x49\x65\x3e\x47\xcf\x14\x50\xde\x16\x27\x3c\xcb\x65\x4d\xce\x4b\xcc\xc2\x56\x67\x87\x6b\x7e\xed\x6f\x67\xd2\x6e\xb7\x51\x3b\x60\xd1\x37\xb2\xee\x7c\xd1\xda\x46\xd6\xff\x17\x89\xd2\x6f\xb2\xdd\x7e\x56\x6f\xe0\xf6\xb3\x2d\x0d\x54\x89\x18\xc1\xe1\x16\x94\xd6\xed\x54\x9a\xd6\xfe\xa5\xf2\xc8\x52\xa8\xef\xd3\x66\x89\xe3\x27\x5f\xa1\x27\x78\x0c\x28\xf9\xee\x00\x51\xf2\x05\x28\x99\x72\x7e\xa8\x3b\xe8\x0e\x4a\x76\xa6\x85\x1a\x85\x9f\x8c\xe2\x6e\x35\x0a\x5b\x29\xcf\xa4\x7e\xa7\xec\x7e\x55\x76\x97\x85\x59\x67\x2a\x3c\x7a\x96\xa6\x23\x1a\x0b\x01\xc1\xc9\xab\x9a\x8d\x63\x0b\x5d\xff\x9e\x5a\xbe\x38\x65\x55\x01\x7f\x00\x7c\xfa\x13\x36\x1f\xc8\x16\xda\xbe\xcf\x4c\x5f\x9c\xb8\x7f\xf0\xc5\xc9\xeb\xfa\x07\x5f\x6c\x21\xe9\x3f\xf8\xe2\x25\xa9\xa9\xa2\xd6\xa9\x65\x8d\xfa\xbc\x79\x55\x93\x88\x7d\xea\x8b\x1a\x8a\x59\x30\xc3\xe2\x4b\xa7\xe2\x7f\x99\x39\xd5\x66\x65\x99\xb6\x3c\x7d\x1d\x78\xc8\xbb\xa5\x01\x0f\x79\x8b\xf8\xf3\x90\x77\x25\xf6\xc7\xdc\xe6\x5f\x62\x16\xc6\x38\x04\xb7\xd1\x37\x1a\xb8\x8d\x7e\x4b\x83\x60\x38\x23\x0c\x45\xc0\xd4\xa4\x6b\x7a\x20\xd8\xd4\xc5\x73\xce\x68\xdf\x51\x8e\xe5\xdf\xfe\xaa\x3d\xa7\xcd\x13\xdd\xff\x48\xd8\x4d\xee\xa2\xa5\x4e\xd0\xbc\xfd\x8c\xb0\xbb\x62\x4c\x5d\x29\x19\x68\x6b\x85\xf0\x63\xf7\x15\xc2\x8f\x9b\x29\x84\x1f\xbb\xa7\xd0\x51\x8a\xfa\xaf\x63\xe9\x0e\x6d\x6d\x06\xae\x96\xb7\x34\x28\xbe\x75\x3b\x58\x8d\x18\x4f\x61\xfc\x65\x7d\x61\x97\xa6\x73\xfe\x15\x9b\xae\x14\xa4\x4f\x03\xe7\xdf\xab\x75\x65\xee\xa3\xe7\x7e\x09\xdd\x81\xa0\xe7\xdf\x39\x5b\xa2\x9e\x7b\x7b\xfa\x61\xe7\x94\xf2\x87\x2e\x04\x3d\x71\xc4\x16\xf2\xa4\xc1\xc9\x07\x7c\xcc\x29\x85\x9c\x69\x3d\xea\xde\x31\x0b\xfc\x9b\x95\x43\x96\xe8\x97\xdc\xee\x8a\x04\x07\x5d\x66\x6c\xfc\xec\xf4\x84\x73\xf9\xbc\x67\xcc\x0b\xf6\x05\x26\x38\xa6\xd2\xf0\xb4\xf9\x1b\xca\xa3\x56\x3c\x8b\xd1\x7f\xc3\x47\x9a\xb7\x4d\x9c\x69\xf1\x79\x09\x33\xcd\x8b\xe8\xd7\x82\x99\xbe\x08\x10\xb4\x78\x12\x4c\xc3\x9e\x14\x82\x1e\x0a\x04\xcd\x9d\x39\x00\x08\x9a\x5f\x4a\x0f\x82\xae\x1c\x39\x18\x08\x5a\x7b\xf6\xaa\xf0\xb0\x02\x31\xae\x7c\x3b\x49\x88\xb1\x70\xdf\x41\x8c\x65\x88\xd1\x12\x96\x97\x06\x31\x16\x43\x75\x10\x63\xfb\x16\x7b\x05\x10\xb7\xb9\x09\xb4\x01\x10\x6d\x73\x7d\x3b\x00\xd1\xd2\xa3\x03\x10\x77\x08\x20\xb6\x88\xb3\x03\x10\xdd\xd9\xdf\xca\xec\x84\xcf\xfe\x8e\x8b\xc0\x15\xcb\x3c\x23\x80\xb8\x5a\xf0\x39\x00\xb1\xda\xc3\xae\xa3\xaf\x05\x10\x6b\x02\x38\x00\x71\x65\xbe\x85\x02\x15\x14\xab\xaa\x41\x1d\x40\xac\x09\xe0\x00\xc4\xfd\xc5\x7e\x0d\x40\xd4\x46\xde\x01\x88\x7b\x88\xbc\x06\x40\xac\x05\xdf\x01\x88\xd5\x7e\xb6\xd0\xc1\x7c\xd0\xd2\x0c\x20\xd6\x14\x71\x00\xa2\xb9\xa3\x7d\x68\x63\x03\x10\xcd\xd3\xe5\x05\xa9\xd3\x69\x00\xb1\xa6\x40\x15\x40\xac\xca\xe9\x00\x44\x07\x20\xb6\x33\x70\xfb\xd9\x96\x06\x47\x05\x10\xab\xf3\x5b\x0f\x20\xd6\xca\x84\x03\x10\xf5\x5d\x1d\xa3\x64\xb7\x05\x10\x0d\x52\x3b\x00\xb1\xb1\xb3\x63\xab\x6a\x07\x10\x6d\xba\x3a\x00\xb1\xa1\xb7\x63\x4b\xdb\x00\x20\xda\xb4\x75\x00\xa2\xa5\xa3\x63\xeb\x6a\x06\x10\x6d\x92\x3a\x00\x71\xbf\xcb\x9a\x75\x00\xd1\x24\xa0\x03\x10\xf7\xb7\xac\xac\x02\x88\xe6\x95\xa5\x03\x10\x37\x31\x40\xed\x35\x58\x03\x10\x8d\xf1\x77\x00\xa2\x03\x10\xdb\x19\xb8\x8d\x7e\x4b\x83\x4e\x01\x88\xba\x83\x3d\x0b\x80\xa8\x3d\x28\x74\x00\xa2\xae\x8b\xfd\x2b\xa4\x05\x10\x2d\x0a\x39\x00\x11\x39\x00\xb1\x95\x81\xab\xe5\x2d\x0d\x8e\x09\x20\x16\x85\xc1\x04\x20\xae\x4a\x81\x03\x10\x4b\xd6\x3b\x0d\x7a\x1d\x40\xac\x45\xdd\x01\x88\x3b\x0f\x7a\x05\x40\xac\x87\xdc\x01\x88\xbb\xdd\x3b\xda\x00\xc4\x5a\xf4\x1d\x80\xe8\x00\xc4\x8d\x00\xc4\x1b\xc0\x62\x3c\x3d\x61\xfe\xf0\x7f\x31\x88\x56\x7f\xb0\x52\x40\x14\xd3\x14\x82\xf9\x4a\x38\xc5\x28\xd1\x56\x07\xa8\xe4\xaa\x1f\x0a\x52\x8c\xb2\xeb\xed\x9f\x51\x4c\xaf\xe4\xff\x33\xcd\x0a\x54\xf7\x20\xcb\x97\xbd\x72\x89\xd1\xea\x91\xc5\xcc\x9f\xeb\x24\xf9\x6e\x17\x73\xd0\xb8\x94\x66\xe6\x9d\x0a\xa0\x87\x3c\x0f\x79\x02\xc6\x6a\x9f\x88\xbc\x19\x79\x54\xc1\xf2\xe6\x7c\x1e\x53\x2c\x0e\x0b\x53\x66\xae\x57\x50\xca\xc2\xeb\x93\x24\x29\x73\xef\x1d\x48\x59\x06\x29\xcd\x51\x79\x69\x1c\x65\x3e\x52\x87\x51\xee\xb2\xc5\x5e\x41\xcb\x8d\xef\x62\x6d\x28\x4b\x4b\x1d\xd8\x0e\xb2\x34\x77\xe8\x18\xcb\x1d\x32\x96\xcd\x61\x76\x88\xa5\x3b\xdd\x5c\x99\x9d\xf0\xe9\xe6\x71\x21\xbf\x6c\xf1\x67\x04\x2c\x8b\x65\xa0\xe3\x2b\xab\x3d\xec\x36\xf4\x5a\xba\xb2\x1a\x7d\x07\x57\xae\xcc\xb7\x08\x7f\x05\x33\x5b\x17\xa0\x8e\x56\x56\xa3\xef\xc8\xca\x7d\x05\x7e\x8d\xab\xd4\x85\xdd\x61\x95\x3b\x0f\xbb\x06\xaa\xac\x46\xde\x31\x95\xd5\x7e\xb6\x10\xc1\x74\xde\xd2\x4c\x54\x56\xe5\x70\x40\xa5\xb9\xa3\xdd\x0b\x63\xc3\x29\x8d\x13\xe5\x05\x49\xd3\x69\x9a\xb2\x2a\x40\x15\xa6\x5c\x97\xd2\xa1\x94\x0e\xa5\x6c\x67\xe0\xf6\xad\x2d\x0d\x8e\x8a\x52\xae\xcf\x6e\x3d\x48\x59\xad\x10\x8e\xa3\xd4\x77\x75\x84\x62\xdd\x16\xa3\xd4\xca\xec\x20\xca\xc6\xce\x8e\x2c\xa9\x9d\xa1\x34\x8b\xea\x08\xca\x86\xde\x8e\xac\x6b\x03\x40\x69\x16\xd6\xe1\x93\x96\x8e\x8e\x2c\xaa\x99\x9e\x34\xeb\xe9\xd8\xc9\x7d\xae\x63\xd6\xc9\x49\x83\x76\x0e\x9c\xdc\xd7\x22\xb2\x8a\x4d\x1a\xd7\x91\x8e\x9a\xdc\xc4\x00\xb5\x15\x60\x8d\x99\x34\x05\xdf\x21\x93\x0e\x99\x6c\x67\xe0\x36\xf4\x2d\x0d\x3a\x85\x4c\xd6\x8f\xef\x2c\xc0\xa4\xee\x2c\xd0\xf1\x92\xba\x2e\xf6\x2d\x8f\x96\x96\x34\xcb\xe3\x60\x49\xe4\x60\xc9\x56\x06\xae\x8a\xb7\x34\x38\x26\x2c\x99\x55\x05\x13\x2a\x59\xd4\x01\x47\x4a\x96\xac\x77\x18\xf1\x3a\x27\x59\x0d\xb9\xc3\x24\x77\x1c\xf1\x0a\x24\x59\x8b\xb7\x63\x24\x77\xb9\x53\xb4\x11\x92\xd5\xd0\x3b\x40\xd2\x01\x92\x1b\x01\x92\x57\x14\x8f\x01\x7d\xe2\x92\x4c\xc8\x18\xd7\xd7\x63\x27\x84\x4a\xe2\xa4\x2a\x45\x40\x27\xcd\xb0\x47\x1c\x59\xc1\xc6\x1d\x62\x8b\xcf\xc3\x7c\x4a\x45\xb8\x8a\x87\x94\x35\xf3\x13\x19\xfd\x5f\x55\x04\x6e\x80\x4e\x56\x05\x82\x95\x5b\xcd\x93\x56\x49\x9c\xee\x54\x9c\xf6\x0a\xdc\x95\x4b\x59\x0b\xdf\x3f\x2b\x45\xac\x7e\x27\x9a\x6d\xe7\xf2\x2e\x09\x1e\x35\x7b\x3e\x30\x09\x02\x6b\xb6\x2f\x6e\xba\xac\xbd\xdb\x81\xe9\x52\x92\xaa\x3c\x3f\x48\xe9\xed\x23\x4e\x88\xb2\x77\xd9\x0c\x28\x7b\xd6\xa1\x94\xbf\xe4\xa1\x4b\xf5\x6e\xa7\xfa\x25\x0f\xcb\x29\x4e\x79\x78\xcc\xd4\x56\xde\x64\x29\xad\x3c\xe9\x50\x2a\xdf\x80\xb8\xaf\x62\x6c\x27\x94\xcd\xcd\x67\x42\x09\xeb\xc0\xa3\x16\xff\x87\x87\x39\x17\x2d\x5a\x51\x1e\xa2\xe4\x0f\xc1\x9f\xc0\xbc\x68\xff\x9c\x40\x94\xa4\x81\xaf\x5e\x94\x1e\x00\xc8\xd3\x37\x4d\x12\xff\x13\x9e\x81\xf5\xbb\xfe\xe7\xa2\xee\x9b\xba\xab\x74\x35\xbb\xfb\x9e\x47\xf2\x78\xee\xb2\x78\x36\x02\x51\x71\x58\xa5\x98\xd9\xe1\x2b\x2e\xba\xe6\x30\xe5\xe1\x5d\x92\xed\x66\xaf\x2f\x79\x78\xa9\x5a\x6c\xe9\xf9\x0e\x0a\x5a\xe9\xd7\x7e\x2f\x35\xeb\xf7\xa6\x72\x46\x87\x67\xf9\x7e\xef\xff\x01\x00\x00\xff\xff\xdd\x01\x04\x1f\x23\xc2\x00\x00")

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

	info := bindataFileInfo{name: "assets/tmpl/config.tmpl", size: 49699, mode: os.FileMode(436), modTime: time.Unix(1484374742, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _assetsTmplHeaderTmpl = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x8c\x53\x5d\x4f\xe3\x3a\x10\x7d\xe7\x57\x58\x7e\xb9\x0f\xf7\x36\xbe\xfd\x58\xda\x5d\x35\x95\x58\x58\x4a\x61\xf9\x58\x68\x41\xec\x9b\xeb\x4c\x12\x87\xd8\x0e\x9e\x49\xd3\x6c\xd5\xff\xbe\x6a\x0a\x14\x21\x90\x78\xcb\xcc\xe4\x1c\x9f\xa3\xa3\xb3\x5a\x45\x10\x6b\x0b\x8c\xa7\x20\x23\xf0\x7c\xbd\xde\x1b\x6e\x3e\x47\x7b\x8c\x31\x36\x34\x40\x92\xa9\x54\x7a\x04\x0a\xf9\x6c\x7a\xdc\x1a\xf0\xd7\x27\x2b\x0d\x84\x7c\xa1\xa1\x2a\x9c\x27\xce\x94\xb3\x04\x96\x42\x5e\xe9\x88\xd2\x30\x82\x85\x56\xd0\x6a\x86\xff\x98\xb6\x9a\xb4\xcc\x5b\xa8\x64\x0e\x61\x9b\x8f\xf6\x1a\xa6\xd5\x4a\xc7\x2c\x98\x21\xf8\x0b\x69\x60\xbd\xde\xd2\x93\xa6\x1c\x46\xab\xd5\xab\xc3\x3f\xc8\xe6\x8e\x86\x62\x7b\x7a\xc2\x42\x8e\x6f\x30\xe7\xf5\x3b\x7f\xd9\x68\xbd\xde\x3e\x37\xcc\xb5\x7d\x60\x1e\xf2\x90\x23\xd5\x39\x60\x0a\x40\x9c\xa5\x1e\xe2\x90\xa7\x44\x05\x7e\x13\xc2\xc8\xa5\x8a\x6c\x30\x77\x8e\x90\xbc\x2c\x36\x83\x72\x46\xbc\x2c\x44\x37\xe8\x06\x7d\xa1\x10\x77\xbb\xc0\x68\x1b\x28\x44\xce\xb4\x25\x48\xbc\xa6\x3a\xe4\x98\xca\xee\xa0\xd7\xfa\x7e\x7b\xaf\xf5\xcd\xe4\x18\xce\xda\xd1\xd8\x9c\x5e\x1f\x3c\xd4\xaa\x3c\x39\x38\xb9\x4e\xba\x9d\x4b\x33\x53\x55\xd5\x77\xb6\x7b\x7d\x1f\x25\xbd\x5b\xf9\xef\x95\xb9\x99\xe2\x1f\x71\xb6\x3f\x58\xcc\xa3\x1f\x59\xda\x2b\x39\x53\xde\x21\x3a\xaf\x13\x6d\x43\x2e\xad\xb3\xb5\x71\x25\x3e\xa7\xf1\x19\x53\x2a\xb2\x19\x06\x2a\x77\x65\x14\xe7\xd2\x43\xe3\x48\x66\x72\x29\x72\x3d\x47\x11\x3b\x4b\x2d\x59\x01\x3a\x03\xa2\x17\xec\x07\xdd\xc6\xde\xeb\xf5\x8b\xc3\x4f\xbc\x2a\x24\x22\x10\x36\x1c\xaa\x44\x72\xe6\x09\xb9\x85\xa2\xf2\xba\x20\x86\x5e\xed\x04\x6e\xb4\x04\x89\x73\x49\x0e\xb2\xd0\xf8\x46\x5f\xf6\x58\x82\xaf\x45\x3b\x68\x77\x82\xde\xd3\xd4\x08\xca\x90\x8f\x86\x62\x4b\x38\xfa\x98\xfd\xb3\x99\x66\x6f\x23\xcd\xde\x4d\x74\xaa\xbe\x4c\x7e\xe9\xf9\xff\x9d\xfe\xe3\xa2\xce\x6e\xce\xe3\x93\xec\xf2\x5c\xfe\x7c\x88\xcb\xbb\xdb\xe5\xef\xe5\xec\xca\x1e\x9e\x1e\xf4\xf3\x8e\x39\xbc\xbb\x98\x14\xe3\xaf\x66\x7c\x78\x34\xa8\xc6\x17\x13\x75\x75\xd4\x9f\x2e\xe5\xc7\x89\xee\xbc\x7c\x6c\xa6\xc8\x25\xc5\xce\x9b\x80\x2a\x4d\x04\xbe\x71\x52\xe9\x28\x01\xc2\x46\xf0\x4b\x69\x4b\x8a\x37\xa5\xdd\x91\x0e\xc5\xb6\xe0\xcf\xad\xf8\x1b\x00\x00\xff\xff\xf8\x1d\xf8\x44\x04\x04\x00\x00")

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

	info := bindataFileInfo{name: "assets/tmpl/header.tmpl", size: 1028, mode: os.FileMode(436), modTime: time.Unix(1483797711, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _assetsTmplIndexTmpl = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\xc4\x58\x5b\xaf\xdb\x44\x10\x7e\x3f\xbf\x62\x58\x2a\x01\x52\x1d\xab\xad\xc4\x43\xe5\x58\x2a\x6d\x05\x95\x4e\x4b\x55\x5a\x24\x1e\xc7\xde\x89\xbd\x3d\xeb\x1d\xb3\xbb\xc9\x21\x58\xf9\xef\x68\xed\x24\x27\x89\xed\xdc\x5a\x20\x0f\x89\xe3\xb9\x7d\x33\x3b\x37\xbb\x69\x24\xcd\x94\x21\x10\xca\x48\xfa\x4b\xac\x56\x37\xc9\x37\xaf\x7e\x7d\xf9\xf1\x8f\xf7\xaf\xa1\xf4\x95\x4e\x6f\x92\xf0\x03\x1a\x4d\x31\x15\x64\x44\x7a\xd3\x34\x9e\xaa\x5a\xa3\x27\x10\x25\xa1\x24\x2b\x60\x12\x04\x33\x96\xcb\xf4\x06\x00\x20\x91\x6a\x01\xb9\x46\xe7\xa6\x22\x67\xe3\x51\x19\xb2\xa2\xa3\x85\xcf\xae\x0a\x83\x8b\x0c\xd7\x2a\xb6\x0c\xbb\x0a\x3e\xcf\xab\x8c\xbd\x65\xb3\xa3\xa0\x53\xa2\x66\x30\xf9\xe4\xc8\xbe\xc3\x8a\x56\xab\x3d\x62\x52\x3e\x49\x9b\x66\x87\xfa\x9d\x83\x8c\x7d\x12\x97\x4f\x0e\xb5\x90\x76\x43\xd2\x6f\x97\x63\xfc\x46\xf6\xd8\x9f\x6e\xc0\x6a\x42\x29\x52\x9c\x7b\xae\xd0\xab\x1c\xb5\x5e\x42\xce\x5a\x53\xee\x01\x8d\x84\xcc\x32\xca\x1c\x5d\xf8\xb7\x84\x3b\x65\xa4\x03\x9e\x81\x32\x33\xb6\x41\x82\x0d\xcc\xd8\xc2\x92\xe7\x49\x5c\x3e\x1d\x72\xf8\x27\xf6\x43\xfe\x22\x94\x96\x66\x53\x51\x7a\x5f\xbb\xe7\x71\xec\xef\x95\xf7\x64\x27\x39\x57\x71\xd3\x3c\x48\x89\x0d\xd2\xcc\x1b\xc8\xbc\x89\x74\xd1\xfe\x38\xce\x15\xea\x48\xe5\xdc\xdd\x5e\xcb\x8b\x34\x71\x35\x9a\x8d\xd0\x0c\x61\x86\x1b\x5a\xb8\xc4\xdc\xab\x05\x89\x34\x89\x03\x5b\x9a\xc4\x78\x4e\x78\x0f\xc1\x62\x5d\xbb\xc9\x0e\xe2\xaf\x0c\xb2\x46\xe7\x4e\xa1\xec\x1f\x6a\x3f\xa2\x44\x5e\x52\x7e\xf7\x2f\x22\xdd\xda\x18\xc7\xda\xc3\x95\xb3\x71\xac\x69\x22\x69\x41\x9a\x6b\xb2\x6e\x52\x30\x17\x9a\xda\xb3\xc7\x5a\xb9\xf0\x15\x2f\x94\x53\x6c\xd6\xa4\x70\xb7\x25\xff\x39\x67\x8f\xee\x5c\x2f\x3a\xe1\x41\x27\x3a\xd2\x39\x29\xd1\x73\xa0\x50\xbe\x9c\x67\x2d\x1c\x75\x8f\x1e\xef\x30\xae\x42\xf1\x9d\x8d\xaa\x95\x1f\x46\xd5\x92\x4e\xa1\x4a\x62\xa9\x16\xe9\x4e\xf7\x99\xeb\x8d\x1a\x83\x0b\x30\xb8\x88\x6a\xa5\xb5\x6b\xaf\x3e\xcf\x9d\x57\x33\x45\x72\xaf\x1d\x25\x5a\x81\x65\x4d\x53\x51\x5b\x72\x64\x7c\x5b\xcb\x5b\x0f\xb6\xd6\x11\x24\x7a\x8c\x3c\x17\x45\x60\xf6\x98\x89\x75\x34\xbe\x2d\xb9\x22\x91\xfe\xc2\x15\x05\x70\x49\xac\xd5\xae\x81\xa6\xb1\x68\x0a\x82\x47\x77\xb4\x7c\x0c\x8f\x16\xa8\xe1\xf9\x14\x26\x2f\xbb\xde\xa2\xd8\xbc\xc5\x7a\x2f\x83\x47\x10\x1d\x85\xd0\x34\x41\xfd\x6a\x25\xd2\xcd\xd5\x30\x94\xfd\x6a\x49\xe2\xb9\x4e\x87\x7b\xb7\xc7\x2c\x0a\x03\x80\x8c\x3f\xe8\xde\x2d\x97\x92\x53\xd1\xba\xbd\xcb\x5f\xa3\x09\x79\x24\x09\x94\x81\x4d\xe0\xf6\x64\x0f\xad\x04\x09\x0d\xed\x77\x24\x43\x94\xec\x80\xc0\xa0\x50\x14\x06\x98\x32\xc5\x08\x7f\x2b\x53\x3e\xdb\x17\xf1\xca\xeb\xee\x28\xbb\xa8\xc5\x9a\x8b\x58\xa4\xaf\xad\x65\x0b\xb7\x5c\x74\x21\x2b\x9f\x8d\x40\xe8\x72\x6d\x88\xd4\xf5\xf7\x5b\x2e\x0e\x3a\xd1\x71\x0f\xc2\xd4\x15\x69\x52\x5b\x4a\x93\x9c\x25\x85\xb1\xd7\xea\x48\xe2\xf6\x6f\x12\xb7\xa4\x63\x66\x07\x5a\xf4\x29\x7b\xe3\x00\xa9\x4a\xdf\xb1\x2f\x95\x29\xc0\x33\xb8\x92\xef\x93\x98\xaa\x2b\x82\xd1\x6f\xc9\x47\x44\xc6\xd2\x21\x0c\xd6\xff\x32\x19\x8e\x0f\x8b\xf4\xe3\xba\xd5\x3f\x94\xad\xbb\x3a\x5d\x2e\x3e\x9a\x2e\xbf\xc6\x3b\xc6\x81\xfa\x50\xd4\x63\x44\xb8\xa6\x23\x0d\xd8\xd8\x6f\x2c\x23\x4c\x9b\xd8\x36\x4d\x30\x72\xd8\x9d\x8e\x1b\x88\x4f\x59\x18\xcb\xb3\x07\x0d\xc7\xe2\x70\xb4\x76\xe0\xd2\x72\x38\x85\x67\x2c\xf7\xff\xd7\x92\xe8\xb2\xea\x4d\x85\x05\xbd\x30\xa8\x97\x4e\xb9\x57\xe8\x8f\xc6\x64\xa4\x88\x7e\xee\x56\x88\xdf\xdb\x45\x05\x5e\xbc\x7f\x03\x1f\xc8\xcd\xb5\x87\xef\x9b\x66\xc8\xc0\x0f\xe3\x65\x03\x67\x9d\xcd\x85\x38\x4e\x9a\xbb\xf4\xe0\xe0\xfa\x2a\x6e\xc3\xf1\xe9\xc3\xed\x99\xa3\xc2\xf2\xfd\x11\x95\x87\xdc\x39\xeb\xa8\x92\xd1\x8f\x27\x44\x60\xbf\x38\x3b\x4c\xbf\xf1\xdc\xe6\xe1\x39\x23\x4d\x54\x55\x80\xb3\xf9\x03\xad\xc5\xbb\x9d\xf5\xaa\x2a\x22\x4b\xae\x66\xd3\x2e\xe7\x8f\x01\xb5\x9f\x8a\x96\x11\x30\x1c\xf4\xdf\x24\x21\x5b\x42\xef\x38\x44\x7f\xa7\xec\xc1\x1a\x0f\xf8\x97\x3a\xbc\x37\x68\xf7\xb2\xb2\xcb\x92\x83\xc1\xfb\x25\x30\xaf\x1e\xda\x97\x34\x9d\xaf\xd3\x70\x06\x6e\x5d\x31\x1f\xb6\x4b\xe1\x76\x11\x1d\x5c\x0c\x87\xd6\x41\xdc\x72\x76\x23\x36\x2a\xac\x92\xa2\x37\x39\xfa\x4f\x23\x43\xc0\x0f\xf7\xdb\x87\xc7\x83\xf5\x75\x12\x77\x2f\x3b\x92\xb8\x7b\x49\xb2\x11\xf9\x27\x00\x00\xff\xff\x13\x7d\x73\xaa\x50\x11\x00\x00")

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

	info := bindataFileInfo{name: "assets/tmpl/index.tmpl", size: 4432, mode: os.FileMode(436), modTime: time.Unix(1484373688, 0)}
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

	info := bindataFileInfo{name: "assets/tmpl/log.tmpl", size: 614, mode: os.FileMode(436), modTime: time.Unix(1484373684, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _assetsTmplNavbarTmpl = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\xd4\x93\xc1\x6e\xdb\x30\x0c\x86\xef\x79\x0a\x42\x3b\xec\xe4\xea\x05\x1c\x5d\x76\x6d\x7b\x19\xf6\x00\x8c\x45\xa5\xc4\x54\x2a\x93\x64\xa3\x83\x91\x77\x1f\x6c\xd9\x99\xe3\x39\xc5\x86\x01\x05\x72\x4a\xf4\xe7\xe7\x4f\xea\xa3\xd2\xf7\x96\x1c\x0b\x81\x12\xec\x0e\x18\xd5\xf9\xbc\xab\x05\x3b\x68\x3c\xa6\xb4\x9f\x54\x28\x1f\x95\x25\x87\xad\xcf\xf3\x91\xa5\xa3\x98\x68\x3e\x3a\x7e\x23\x5b\xe5\x70\x52\x66\x07\x00\x50\x5b\xbe\xe4\x34\x41\x32\xb2\x50\xac\x9c\x6f\xd9\x4e\x8e\xb5\x6b\x0a\x7a\x21\xb4\x14\x17\x9e\xd1\x77\x68\x73\x0e\xb2\xb2\xe6\x70\x3c\x7a\x82\x26\x78\x8f\xa7\x44\x56\x81\xc5\x8c\x93\x3c\xb4\x2d\xfa\x2c\x63\x3c\x52\xde\xab\x4f\xa5\xfa\x89\xa4\x4d\x0a\x30\x32\x56\xf4\x76\x42\xb1\x64\xf7\xca\xa1\x4f\xb4\x6a\x3e\x0e\x90\x4e\x78\x69\xcf\x4d\x90\x6a\xe0\x65\x6a\x3d\xe8\x1f\x69\xaf\x75\x41\x71\xad\xf6\x3d\x3b\x78\xf8\x96\x28\x3e\xe3\x2b\x9d\xcf\xd7\x25\xb8\xe2\x76\x88\x28\x56\xc1\x4b\x24\xb7\x57\x5a\x99\xbe\x5f\x94\x7e\x4e\x70\x08\xb9\xd6\xb8\xee\x40\x3e\xfd\x73\xf2\xd3\xcf\x1b\x59\x62\x17\x51\xb5\xb6\xdc\xbd\xfb\x28\xe6\x55\xc2\xef\x9d\xb2\x9d\x7f\x2d\x9b\x5c\x61\x6a\xfd\x22\x62\x7e\xa4\x82\xdd\xc6\x6e\x47\x78\xf4\x03\x1e\x9e\x47\xd7\x80\x01\xd4\x97\x20\x8e\x8f\x6a\x75\xe1\x31\xda\xf3\x22\xba\xe2\x4c\xaf\x80\x4d\xe6\x6e\xeb\xdd\x6c\x62\xbb\x91\xb2\x5d\x7e\x45\xea\x26\xfc\xca\xb3\x7c\xbf\x80\x6f\xc6\xe1\xb5\x32\xe5\x16\x7f\x2c\xa0\x40\xf7\x6c\x76\x7f\x07\xe3\x31\xdc\x2b\x09\x1f\x06\x0c\x8f\xe1\xff\x19\x7c\xcd\x98\xdb\x74\xa7\x18\xd2\x38\xbc\x56\xa6\xdc\xe2\x1d\x18\xd7\x4a\xeb\xcd\xd6\xbf\x74\xfa\x5a\x6b\xc1\xce\xec\xe6\x99\x7e\x05\x00\x00\xff\xff\x32\xc4\x33\xee\x4f\x06\x00\x00")

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

	info := bindataFileInfo{name: "assets/tmpl/navbar.tmpl", size: 1615, mode: os.FileMode(436), modTime: time.Unix(1484152771, 0)}
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

	info := bindataFileInfo{name: "assets/tmpl/status.tmpl", size: 3773, mode: os.FileMode(436), modTime: time.Unix(1484392664, 0)}
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

