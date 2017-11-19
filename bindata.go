// Code generated by go-bindata.
// sources:
// assets/css/custom.css
// assets/tmpl/config/action.tmpl
// assets/tmpl/config/filter.tmpl
// assets/tmpl/config.tmpl
// assets/tmpl/header.tmpl
// assets/tmpl/index.tmpl
// assets/tmpl/navbar.tmpl
// assets/tmpl/setup.tmpl
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

var _assetsCssCustomCss = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x8c\x92\x5d\x6a\xeb\x40\x0c\x46\xdf\xbd\x0a\x6d\x60\x4c\x20\x17\x6e\x99\x40\xf6\x22\x8f\x14\x67\x5a\x47\x1a\x64\x25\x4d\x5a\xba\xf7\x62\xc7\x31\x53\x28\xa1\x7e\xf0\xd3\xd1\xd1\xcf\x37\x9d\xd2\x0d\x3e\x1b\x58\xbe\x82\x44\x59\xfa\xe0\x5a\x22\xfc\xdf\x94\xeb\xae\xf9\x6a\x9a\xf6\xf5\x7c\xea\xd4\x4d\xa5\x42\x3b\x4c\x6f\xbd\xe9\x59\x28\x24\x1d\xd4\x22\xb8\xa1\x8c\x05\x8d\xc5\x77\x2b\xe6\x7c\xf5\x80\x43\xee\x25\x42\x62\x71\xb6\xbb\x32\xa9\x38\x66\x61\xab\x94\x27\xbc\x86\xf7\x4c\x7e\x8c\xf0\xb2\x79\x34\xa7\x7c\xa9\xfa\x23\x8c\x05\xa5\x3d\x60\x55\x76\x50\xf1\x30\xe6\x0f\x8e\xf0\xef\x51\x55\x8c\xeb\x59\xd5\x88\x2d\x82\xa8\xf0\xee\xaf\x1b\x4c\x53\x1e\x30\x60\xf2\x7c\xa9\x65\x0b\x4b\x4a\x3d\x5b\x37\x9c\x79\x45\x0b\x8e\xe3\xaf\x6c\x6f\x7c\x9b\x29\xc7\x6e\x60\xd8\x83\x1f\x19\xe9\xe9\x31\x71\xc8\x89\x57\xfb\x5c\xd7\xce\xff\x30\xba\xe5\xc2\x34\x59\xe6\xf4\xf6\xe0\xd6\x12\x0f\xec\x5c\x2b\xb5\x60\xca\x7e\x8b\xb0\x69\xb7\xcf\x76\x5e\x67\x9b\x0e\xed\xd8\x85\x29\x19\x16\xff\x91\x8b\xf5\x59\xee\x8f\x62\xbb\x5c\xf8\x3b\x00\x00\xff\xff\x0a\xc5\xa6\x4c\x39\x02\x00\x00")

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

	info := bindataFileInfo{name: "assets/css/custom.css", size: 569, mode: os.FileMode(420), modTime: time.Unix(1503580372, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _assetsTmplConfigActionTmpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xb4\x96\xcf\x6e\xdb\x30\x0c\xc6\xef\x79\x0a\x42\x40\x81\xf6\x30\x19\xc3\x76\x54\x0d\x0c\x05\x76\x2e\xda\xdc\x07\xd9\xa2\x17\x21\x9e\xe4\x49\x9c\x93\x22\xc8\xbb\x0f\xb1\xe4\xe5\x4f\xed\x4c\xce\xbc\xa3\x29\x91\x1f\x7f\x5f\x48\xc7\xbb\x9d\xc2\x4a\x1b\x04\x26\x4b\xd2\xd6\x7c\x2b\xad\xa9\xf4\x77\xb6\xdf\x2f\x84\xd2\x2d\x94\xb5\xf4\xfe\x91\x91\x2c\x6a\xfc\xe0\xd0\x37\xd6\x78\xdd\x22\xcb\x17\x00\x00\xa2\x8b\x9f\x5d\x82\x70\xd5\x93\xd3\x0d\xaa\xf8\x54\x58\xa7\xd0\xa1\x8a\x69\x21\x75\x85\x52\x1d\x9f\x43\xcc\x9d\x07\x42\x50\xe5\x22\x23\x35\x7c\x62\xe4\x0f\x1c\x3f\x6d\x65\xfd\x6b\xe0\x58\x64\xa7\x42\x22\xbb\x68\x45\x50\x61\xd5\xdb\x79\xca\x6e\xa7\x2b\xb0\x0e\xee\xf1\x27\xf0\xe5\x5b\x83\xc0\x7c\x2d\xcb\x35\x7b\x38\x0d\x6d\xb0\x58\x59\xbb\x66\x0f\xfb\x7d\x12\x18\x38\xbb\xf1\x8d\x34\x8f\xec\x23\xcb\x97\x1b\x4d\x84\x6e\x9c\x86\x36\x88\x34\x7e\xfc\x2e\x18\xfa\x2e\x57\x58\xae\x0b\xbb\x05\xfe\xa5\xfb\x85\x79\xd4\xe1\xcb\x43\x39\xb8\x6f\x9c\x36\x54\x01\xbb\xf3\x77\x9e\x01\x7f\x76\x58\xe9\x2d\x30\x1e\xe6\x81\x53\xbc\xdd\x89\xbf\x03\x8b\xfe\x5d\xf5\x37\xb4\x81\x46\x5d\xe4\x76\x96\x1e\xbd\x8b\x42\x6c\xb2\x75\x9f\x12\xac\x73\x38\xb3\x79\x2f\xa1\x60\xb2\x7d\xb1\x81\x9b\x0d\x1c\x5d\x8d\x4a\xb6\xd6\x69\xba\xb2\x04\x53\xd1\xbe\xc6\x8a\xc9\x6c\x7d\x0b\xf3\xc3\x95\xb6\xae\xb1\x13\xf3\x93\xf9\x6a\xed\x69\x89\x5b\x1a\x42\x7c\x3a\xd6\x4d\xa6\x3c\xe9\xe5\x7f\xae\x41\x78\xab\x4c\x5e\x82\xcf\x2c\x7f\x3d\x64\x8e\xdb\xd4\x68\xf3\xef\x33\xd2\x69\xf0\x67\x6d\xfe\x6a\x5b\xc7\xc1\x1b\x6d\xe6\x9f\x0a\x4f\xf2\xca\xaa\x4f\x43\x79\x25\xe9\x12\x59\x0e\xb2\xf3\xc3\x38\x94\x33\x0e\x78\x80\x7a\xe9\x6b\x26\x92\xfd\xe9\x61\xd6\xc1\x1e\x44\x1e\x79\xe9\xc3\x90\xee\xc5\xff\xe3\x95\xf9\x1e\x92\xef\xdd\x2b\x57\xd2\x18\xac\xe7\xf4\xf7\x29\x96\x4c\xb4\xb7\xef\xe0\x26\x77\x45\x76\xf2\x35\x22\xb2\xee\x7b\x2a\x5f\x88\x4c\xe9\x36\x5f\xf4\xdc\xbf\x03\x00\x00\xff\xff\x1e\x69\xb6\xee\xc6\x09\x00\x00")

func assetsTmplConfigActionTmplBytes() ([]byte, error) {
	return bindataRead(
		_assetsTmplConfigActionTmpl,
		"assets/tmpl/config/action.tmpl",
	)
}

func assetsTmplConfigActionTmpl() (*asset, error) {
	bytes, err := assetsTmplConfigActionTmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "assets/tmpl/config/action.tmpl", size: 2502, mode: os.FileMode(420), modTime: time.Unix(1511084443, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _assetsTmplConfigFilterTmpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xcc\x58\x4d\x6f\xdb\x3c\x0c\xbe\xe7\x57\xf0\x15\x50\xe0\xdd\x21\x0e\x50\x14\xdb\x0e\x8e\x81\xee\x90\xb5\x43\x8a\x05\x6d\xb7\x6b\xa2\xc4\xb4\xa3\x56\x96\x3c\x49\xf9\x82\x91\xff\x3e\xc8\x5f\x69\xda\xa4\xb1\x33\x1b\xe8\xad\x91\x45\xf2\x79\x1e\xd2\x24\xeb\x24\xf1\x31\x60\x02\x81\x04\x8c\x1b\x54\xe3\x99\x14\x01\x0b\xc9\x76\xdb\x71\x7d\xb6\x84\x19\xa7\x5a\xf7\x89\xa1\x53\x8e\x5d\x85\x3a\x96\x42\xb3\x25\x12\xaf\x03\x00\xe0\xa6\xe7\x7b\x97\x20\xbb\xaa\x8d\x62\x31\xfa\xf9\xaf\xa9\x54\x3e\x2a\xf4\x73\xb3\xcc\x74\x8e\xd4\xdf\xfd\xce\xce\xd4\xfe\x41\x76\xe8\x7b\x6e\xcf\xf8\x75\x9f\x08\x1a\xe1\xf1\xa7\x4b\xca\x17\x07\x1e\xbb\xbd\x97\x10\xdc\xde\x2b\x90\xae\x99\x4a\x7f\x53\x01\x74\x92\xb0\x00\xf0\x0f\x38\x8f\x9b\x18\x81\x98\x15\x33\x06\x95\x55\xf5\x00\x18\x50\x72\xa5\x63\x2a\xfa\xe4\x0b\x81\x99\xe4\xd9\xdf\x97\xc4\xfb\x8e\x02\x15\xe5\x87\x59\x24\x09\x72\x8d\xb0\x17\x47\x73\x3a\x7b\x3e\x19\xe5\xb2\x5e\x14\xe1\x1f\x76\xe8\xc5\xd4\x92\x12\xfa\xb8\xca\x6f\x0e\x33\x97\x9c\x69\xf3\x88\x6b\x33\x95\x6b\x70\x06\x69\xdd\x39\xa3\xdc\x19\xfc\x1f\x2b\x26\x4c\x00\xe4\x42\x5f\x68\x02\xce\x48\x61\xc0\xd6\x40\x9c\xac\x40\x9d\x22\x2a\xf9\x74\x08\xd6\x89\x8c\xc2\x7b\x55\x36\xa7\x1a\x22\xf4\x19\xad\x4d\x68\x2a\x25\x7f\x40\x8e\xb3\x3d\x4a\x37\x54\xdf\x59\x77\x27\x29\xcd\xa9\x1e\xa7\x81\xcf\xe6\x54\xad\xde\x8e\x32\x5f\x28\x0e\x8d\x66\xf3\xd7\xfd\xb0\x72\x42\x17\x8a\x8f\xdb\x4b\xaa\x42\xb3\x42\x34\xe8\x37\x94\xd4\xfb\xc2\xdf\x49\x5e\x65\xe4\xe6\x49\x05\x74\x29\x15\x33\x08\x66\xae\x50\xcf\x25\xaf\xcf\xce\x64\x19\xfb\x19\xdc\x0a\x33\x32\xaa\xe4\x37\xc8\x5d\x3f\x16\x9e\x4f\xf2\x2c\xc0\x8c\x4b\x30\x2d\x66\xb1\x05\xc6\x65\x46\xab\x53\x2e\xe1\xb4\xc9\x99\x53\x11\xd6\x24\xe9\x32\x11\x2f\x0c\x04\x52\x45\x7d\x22\x97\xa8\x56\x36\x31\x04\xcc\x26\xc6\x3e\xb1\x02\x10\xb0\x73\xb1\x4f\x92\x24\xa7\x38\xb1\x14\x27\x25\xc5\x49\x41\xd1\x46\x9f\x6c\xb7\x04\xd2\x49\x69\x0d\x0a\xbd\x86\x54\x84\xdb\x2d\x39\x00\xab\x52\xa3\x7a\x3b\x4f\xdc\xff\xba\x5d\xf8\xcd\x34\x93\x02\xae\x47\xb7\xd0\xed\x56\x53\x68\x37\xd1\xbe\x12\x2f\xb3\x3f\xaa\xd7\xde\xc4\x1b\xd2\x29\x1e\x99\x77\x35\xfb\x5c\x16\xd5\x49\x1d\x9e\xac\x9a\x65\x76\x99\xdb\xcb\x8d\x97\xcb\x4e\x8c\x2b\xe2\x0d\xe8\xec\x9d\xdd\x87\x8a\x10\x15\x70\xf6\x8c\x9c\xcd\xa5\xac\xfb\x26\x35\x58\x64\xb9\x22\x01\x9d\xa1\x93\xa2\x1a\xef\x50\x1d\x2e\xbe\x5c\x70\x4b\xd0\xb9\xb6\x16\xc3\xd2\xe0\xec\x9a\x3c\xfa\x06\x4e\xf9\x42\x29\xf4\x3f\x9c\x54\x39\xae\x57\x62\x4d\xde\x57\xeb\x5b\x66\xd4\xa6\x5e\x76\x59\x5e\x21\xfd\x78\xb5\x55\x00\xab\x55\x5e\x37\xb9\x51\x9b\x8a\x3d\xc9\xcd\x87\x13\xeb\x49\x6e\x6a\xe9\xf4\x43\x6e\x5a\x94\x68\xaf\x71\xdb\xf6\xdb\x64\xdf\xb6\xc7\x55\xdb\x76\xaa\x67\xe3\x5d\x7b\x7f\x2c\x09\x3f\xa2\xea\xb9\xd9\xc9\x94\xf9\xac\x3e\x9c\xb2\xfb\x2d\x33\x95\xa1\x6c\x94\xa5\x0c\x65\x65\x86\x32\x94\xe7\xb3\xb3\xab\x8a\x5d\x80\x16\x34\xc4\x33\x97\x95\x57\xff\x7e\x17\xde\x8e\xeb\x11\x31\x01\x1a\x85\x61\x11\x8a\xfa\xe5\x5f\xae\xbd\x03\x2e\xa9\xf9\x7c\xf5\x72\xf5\x2d\x62\x3b\x77\x4c\x3c\x14\x11\x4e\x0a\xc9\x0b\xab\x88\x89\x71\x09\xac\xf9\x05\x38\xa2\xeb\xd6\x79\xd3\xf5\x39\xbc\xe9\xfa\x1f\x79\xbb\xbd\x17\x5f\x92\xdc\x5e\xfa\x95\xcc\xeb\xb8\x3d\x9f\x2d\xbd\x4e\xb1\x23\xff\x0d\x00\x00\xff\xff\x1e\xb5\x5e\x46\x9c\x13\x00\x00")

func assetsTmplConfigFilterTmplBytes() ([]byte, error) {
	return bindataRead(
		_assetsTmplConfigFilterTmpl,
		"assets/tmpl/config/filter.tmpl",
	)
}

func assetsTmplConfigFilterTmpl() (*asset, error) {
	bytes, err := assetsTmplConfigFilterTmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "assets/tmpl/config/filter.tmpl", size: 5020, mode: os.FileMode(420), modTime: time.Unix(1503125362, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _assetsTmplConfigTmpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xec\x5c\x5f\x6f\xdb\x38\x12\x7f\xcf\xa7\x98\xe3\x2d\xea\x16\x5b\xdb\xd7\xf6\x6d\x6b\x19\x58\x24\xbb\xb8\x02\xd7\x34\x97\xa4\x0f\x87\xc3\x21\xa0\xa5\xb1\xcd\x2d\x45\xea\x48\xda\x89\x61\xe4\xbb\x2f\x28\x52\xb2\x2c\x4b\xb2\x9c\xc4\xd9\xfc\xb1\x1f\x1a\x5b\x1c\x0e\x7f\x9c\x19\xce\xf0\x37\xa9\xb3\x5c\x46\x38\x66\x02\x81\x84\x52\x8c\xd9\x84\xdc\xde\x1e\x0d\xfe\x76\xf2\xed\xf8\xf2\x3f\x67\xbf\xc1\xd4\xc4\x7c\x78\x34\xb0\x3f\x80\x53\x31\x09\x08\x0a\x32\x3c\x5a\x2e\x0d\xc6\x09\xa7\x06\x81\x4c\x91\x46\xa8\x08\xf4\xec\xc4\x91\x8c\x16\xc3\x23\x00\x80\x81\x0e\x15\x4b\x8c\xfb\x60\x5f\xe3\x99\x08\x0d\x93\x02\x22\xe4\x68\xf0\x5c\x5e\xbf\x1d\xcd\x8c\x91\xe2\x3d\xc4\x32\xa2\xfc\x1d\x2c\x73\x59\xfb\x9a\x53\x05\x46\x41\x00\x4e\xaa\x97\x50\x85\xc2\x9c\xca\x08\x0b\x6f\x3f\xaf\x4d\x61\x63\x78\x6b\x54\x2f\xe4\x54\xeb\x53\x1a\x23\x04\x41\x00\x1d\xb7\x5e\xd4\x29\x2f\x90\x82\x92\x0a\xde\xda\x95\x18\x04\xf0\x8f\xcf\xc0\x60\x00\x56\xc3\x94\xf1\xc8\xea\xd7\x3d\x8e\x62\x62\xa6\x9f\x81\xfd\xfc\x73\x95\x82\x0c\x69\x3a\x03\x82\xf5\xc9\xff\x65\xff\xab\x9c\x60\x71\xa6\x52\x3d\xb3\x48\x3c\xca\x29\x8b\x22\x14\x1d\x78\xf3\xc6\xe9\xea\x09\x19\xe1\x6a\x13\x5f\x4e\xcf\xbe\x5f\x56\x6e\x21\x7b\xb9\x59\x73\xca\x67\x08\x01\x74\xc6\x94\x6b\xec\x7c\xae\x15\x5f\xb7\x13\x74\x1a\x44\xbd\x03\x98\x10\xa8\xfe\x79\xf9\xf5\x5f\x56\xfc\x24\x35\x6a\xcd\xa4\xdb\x8d\xa7\xeb\x4f\x6e\x01\xb9\xc6\x57\xe5\x0e\xa3\x66\xbb\x78\x23\x8b\xd9\xdd\x9c\x72\x8e\x73\x54\xe6\xae\x4e\x39\x5a\x7f\x37\xe8\x17\x4f\xf0\x20\x62\x73\x48\x21\x06\x36\x51\x18\xca\x04\x2a\xb2\x3a\xdd\xc5\x8c\x20\xe8\x7c\x44\x7d\x46\xc8\x05\x8a\x0a\xfe\x98\xc5\x23\x69\x94\x14\x05\x05\xa9\xcc\xf4\xc3\xf0\x38\xcd\x42\x33\x45\x6d\xae\x18\xf4\xa7\x1f\xca\x22\x1f\x33\x2d\x1c\x69\x44\x72\x79\x84\x85\x9c\x29\x90\xd7\x02\x46\xd2\x74\x34\x68\x39\x53\x21\xea\xf7\x30\x66\xdc\xa0\xd2\x40\x45\x04\x34\x4d\x41\x7a\xd0\x9f\x7e\x2c\xe9\x1d\x4b\x15\x03\x8b\x02\x22\xe7\xa8\xae\x15\x33\x48\xbc\x74\x40\xfa\x2e\x35\xf6\x09\xc4\x68\xa6\x32\x0a\x48\x22\xb5\x21\x80\x22\xb4\x11\x13\x90\x78\xc6\x0d\x4b\xa8\x32\x7d\xab\xa6\x1b\x51\x43\xc9\x70\x90\x7e\x28\x2d\xc3\x44\x32\x33\x36\xce\xe3\xb5\x95\x9c\x1a\x3d\x1b\xc5\xcc\x10\x48\xa3\x26\x20\x17\x74\x8e\x24\xdb\xed\xc8\x08\x18\x19\xd1\xe5\x93\xf4\x47\xa2\x58\x4c\xd5\x82\x80\x36\x0b\x8e\x01\x89\x98\x4e\x38\x5d\xfc\xc2\x04\x67\x02\xc9\xf0\xa8\x62\x77\xe5\xed\x8c\x19\xc7\xc2\x9e\x26\x68\xea\xd5\x95\x63\xc7\x6f\xa4\x12\xf7\x89\xbc\x16\x5c\xd2\x68\x1b\xf6\x92\x69\xbc\xb9\x76\x07\xde\xc2\x19\x0f\xb2\xaf\xef\xc9\xee\xbb\xda\xd0\x69\xb1\x13\x10\x34\xb6\x48\x17\x23\x69\x7a\xbe\xee\xb6\x82\x58\x8e\xa9\x41\x3f\x62\xf3\x82\xcd\x96\x4b\x36\x86\xde\x57\xd4\x9a\x4e\xf0\xf6\xb6\xf2\xf4\x51\x8e\xca\x40\xfa\x6f\x37\xa2\x62\x62\xcf\xf1\x72\xb9\x9a\xe4\x75\xae\x54\xa2\x88\xd6\x4e\xf2\xf4\xd3\xf0\xf2\x9a\x19\x83\x0a\x2e\x59\x8c\x16\xa6\x3d\x51\x9f\x86\x95\xcb\x19\x3a\xe2\xd8\x55\xa8\x13\x29\x34\x9b\x6f\xec\x28\x1d\x5f\x13\x06\x37\x45\x1b\xc5\x12\x8c\xfc\xa7\x91\x54\x11\x2a\x8c\xaa\xec\x6b\xec\x1d\x64\xf3\xb9\x1b\x53\xd5\x03\x6e\x30\x1a\xea\x50\x21\x8a\xd4\x21\x7a\xd0\x37\x35\x6a\x32\x69\xbc\x09\xf9\x2c\x42\x50\x98\x70\xd6\x66\x02\x13\x7e\x82\x69\x21\x1c\xca\x99\x30\xdb\xc5\xa4\x99\xa2\x6a\xa1\xae\x5e\x62\xd0\xaf\x32\xcb\xa0\x5f\x63\xc8\x81\x59\xdd\xeb\xca\xaf\xe5\x52\xd9\x18\x82\x9f\x98\x88\xf0\xe6\x3d\xfc\x64\x7c\x48\xc0\x2f\x01\xf4\x5c\x76\xee\xf9\x70\xe9\xe5\xe1\x72\xbb\x59\x8d\x60\xab\xb7\x1a\x73\xa7\x2b\xd8\xd9\xc9\x32\x7e\xc1\x0c\x8c\xee\xf9\x8a\x9a\x1f\xe6\xf4\x7e\x44\xfa\xcd\x16\xac\x1d\x74\x3b\xe7\x4c\x9b\x4b\xbc\x31\x23\x79\xb3\xda\x77\xef\x22\x8d\x28\x5b\xc9\x35\x54\x00\x71\x01\x77\x95\x06\x1c\xa9\xb1\x03\x38\x6f\x6c\xf1\xef\x16\x74\x23\x29\xf9\x05\x72\x0c\x4b\xf8\x7e\x73\x31\x7c\xee\x42\xb8\x0a\xa2\x8f\xf2\x2b\x1f\xe5\x7f\x09\xca\x2f\xee\xe0\x9c\x9b\x4a\x84\xfe\x58\x5d\x29\xb3\x5f\x74\xc6\x79\xf7\xdb\xf8\x8b\x30\x67\x46\x15\xf0\x1d\xdb\xb3\x5a\x05\x2d\x3d\xc4\x7b\x04\x35\x70\x37\x3f\x1f\xf6\xee\xc3\x46\x4d\x62\x62\x2c\x09\xd8\xd2\xd7\x35\x72\x32\xb1\x75\x25\x25\x59\xd9\x33\xaa\x26\x68\x02\xf2\x77\x0f\xbf\x9b\xc3\xef\xba\xdb\xd2\xd5\x72\xe9\x0e\xf4\xed\x6d\x45\xc2\x2d\xbf\x22\x34\x94\x71\x0d\x53\x54\xd8\x8c\xbd\xef\xf0\x36\x6c\xff\x11\x8c\xe3\x8b\x1e\x48\x11\x72\x16\xfe\x08\xc8\x8a\x91\x9a\x29\xd3\xef\xa1\xd3\xc2\x2c\x9d\x77\x64\xe8\x48\xd0\x7d\x36\x55\x9d\x88\xa1\x50\x76\xab\xcd\xb0\xa5\xa8\x35\xdb\xa8\xf2\x4e\x95\x6f\xb5\x4f\xa3\xa8\x74\xb7\xda\x1e\x01\x4d\xd7\xa6\x53\xbc\xde\x70\x41\x76\x59\x6a\xc8\xbf\x50\x71\xdd\xd9\x1c\xbf\x7b\x0d\x7c\x2a\x12\xb5\xb5\x78\xb3\xe6\x0e\xfa\xe9\x5d\xa8\xe9\xfa\xf7\x00\xc5\xb8\x78\x79\x4b\x93\x06\x8c\x69\x84\x24\x25\x47\x6d\xf2\x85\xbd\xb0\xa5\xef\x03\xd2\xfd\x40\x40\x49\x77\xab\xa5\x5c\x4e\x08\x50\xc5\x68\x97\xd3\x11\x72\x8e\xd1\x68\xd1\x4a\x63\xd7\x30\xc3\x2b\x2f\xeb\x65\xa4\xdd\x6c\x19\xbf\xa8\x0c\x67\x31\x8a\xba\x00\xde\x9c\x6e\x99\x6d\xbd\x7c\xf5\x1c\xdf\xff\xba\x47\x5e\x0a\xb9\xd4\xe8\x33\x73\xc4\x74\xcc\x72\xe5\x64\xf8\xc6\xda\x45\x7f\xde\x9e\x63\xc0\xdd\xcc\xd7\xb1\x39\xc3\xb5\xf5\x5c\x66\xe7\x13\x97\xce\xd7\xaf\xf4\x1b\x6b\xad\xb3\x84\xed\x76\xb2\xe1\xbc\xcd\x4a\xd3\x4f\xc3\xdf\x53\x54\xcd\x6b\x43\xa9\xe1\xe0\x77\x92\x11\xa9\xb7\x02\xaf\xbf\xd2\x04\xc8\xe5\x22\x41\x92\xd7\x68\x02\xc4\x29\x27\x85\x32\xee\x9e\x00\x39\x53\x38\x66\x37\xa4\xa2\xa0\x93\x77\x0d\xb5\x3c\x43\xfd\x6b\xe8\xdb\x15\x3b\xa0\x76\x19\xb8\x05\x6a\xa7\xbc\x88\xda\x3d\xb9\x33\xea\x9d\x7d\x37\x96\xd2\xdc\x2f\xc6\xf3\xda\x8b\x63\x3a\xe3\xa6\x26\xda\x8f\xed\x49\x68\x53\x4f\x6b\xf1\xd7\x0c\x55\x3c\xde\x85\xe6\xfe\x4e\xe7\xd2\x92\x8d\xd7\x42\x73\x5f\x04\x11\x1d\x7b\xa7\x55\xd5\xbe\xdc\xa1\x8f\x45\x44\x33\x30\x8f\x40\x44\xb3\xa5\xaa\x89\xe8\x0a\xc8\xa3\x11\xd1\x0d\x12\x95\x23\x2c\x91\xa8\x15\xb6\x67\x49\xa2\x72\xf8\x07\x12\x55\x24\x51\x0d\x66\x79\x69\x24\x2a\xdf\xea\x81\x44\xb5\x97\xd8\x2b\x01\xba\x4b\x11\x68\x43\x80\x9a\xce\xfa\xdd\x08\x50\x83\xc6\x03\x01\x7a\x40\x02\xd4\xc2\xce\xaf\x80\x00\xe5\x25\xb8\x96\x00\xe5\x76\x7a\x42\x04\x28\x47\x5d\x4b\x80\xda\xa1\x3e\x10\xa0\x1d\x09\xd0\x05\x52\x15\x4e\x9f\x31\xff\xf9\xff\x0c\x55\xab\x5f\xd8\x29\xd4\x33\xee\x8a\xf0\x2b\xe1\x49\x3a\xf5\x6d\x55\x81\xcc\xbc\xfe\x58\x24\x49\xfb\xf5\xf6\xcf\x91\xdc\x4a\xbd\x7f\xbb\xa8\x80\x4d\x04\x3e\x5e\xf6\xca\x8b\xf4\xea\xd7\x5e\x1e\xcf\x79\x1a\x7c\x36\xff\x55\x40\x72\x91\x79\x65\x5c\x76\x24\x40\x14\x86\xb6\x56\x03\x89\xd9\x8d\x35\x16\x49\x64\x32\xe3\x54\x3d\x2e\x99\xf3\xd0\x4b\x54\x2e\x47\xfd\x2c\x99\x5c\x86\xfe\x40\xe4\x8a\x44\xae\xde\x2a\x2f\x8d\xc7\x65\x3b\x3d\xd0\xb8\x87\x94\xd8\x2b\xd1\xdb\xb9\x8a\xb5\x61\x79\x0d\x79\xe0\x6e\x24\xaf\x5e\xe1\x81\xe3\x3d\x20\xc7\xdb\x6e\xe6\x57\x40\xf1\x7c\x61\xae\x25\x78\x99\x91\x9e\x10\xbf\xf3\x90\x6b\xd9\x5d\x2b\xc8\x07\x72\xb7\x23\xb9\xfb\x66\xe9\x09\x5c\xa0\x31\x4c\x4c\x9e\x2d\xc5\x83\x50\x72\x9d\x50\x11\x90\x8f\x64\xf8\x45\x18\x54\xd4\xc7\xec\x96\x8a\xb5\x36\xf1\x8c\xd3\x10\xe1\x54\x1a\x36\x66\x21\x6d\x35\x7f\xe7\x8a\xd7\x66\x37\xc3\x78\xa1\x91\x8f\x5b\x72\x4b\x98\xe9\x56\x04\x73\x2f\x4a\x23\xff\x9f\xf5\xd3\x6f\x6f\x24\xa8\x98\x8c\x58\x08\x7f\xc8\xd1\xa3\xd0\xd9\xfb\x5d\xff\x96\xcb\x70\x8a\xe1\x0f\xcb\xc1\xca\xd7\x86\x42\x08\xf5\x7e\xe5\x5c\x5e\x5f\x20\x1f\xaf\x52\x11\x2b\x0c\x53\x3b\x7c\x65\x4d\xbb\x57\xea\x55\x64\xb0\x4d\x68\xbf\x5b\xbf\x55\x23\x4d\x5d\xba\x57\x90\xb5\xf6\x2c\x9e\xa9\x5e\x7a\xcc\xaa\xcc\x2a\x8a\x52\x49\x2a\xf5\x04\xac\x5b\x81\xbd\x64\xe4\x0a\xdc\xfb\xb6\x75\x73\x87\xc5\xf2\xf4\x72\x7f\x25\x3b\xa9\x39\x1d\x59\x2e\xcb\x3b\x3d\xf1\x22\x8d\x4c\xf7\x71\xae\xf8\xb6\x44\x5d\x70\x1a\xfe\x00\xff\xd5\x84\x67\x5b\x9a\x86\xe1\x94\x0a\x81\xfc\x5e\x2d\xc0\xed\x12\xfb\x6c\x12\xc6\x7a\x52\xe4\x56\xa9\x5f\xb2\xaf\x8c\xec\xbd\x3f\xa8\xd3\xd5\x62\xbf\xda\xfe\xbb\x83\xb1\x9e\xf4\x8e\xbd\xcb\xa0\xbc\x7a\xe6\xcb\xe7\xd2\xd3\x4a\xe1\x77\x33\xf8\x87\x8e\x96\x7e\x0f\x9d\x6d\x36\x79\x69\xfd\xac\x6c\xa7\x8f\xda\xcf\x4a\x19\xb9\x5f\xb9\x2b\xf0\xba\xeb\xdd\xf5\xa4\xfb\x5c\x7b\xed\x51\xed\x92\x44\x9b\xda\x53\x5b\x8f\xf4\x6e\xcd\xa9\x6d\xea\x0e\xad\xa9\x07\x68\x4d\xb5\x35\xf2\x4b\x6b\x4c\xa5\xfb\x2e\xb6\xa5\x6c\x71\xdd\xe8\x49\xad\xd7\xd8\xbf\xaa\x23\x95\x61\xcd\xfb\x51\x16\xeb\x46\x33\x6a\x07\xac\x87\x56\xd4\x96\x56\x94\x1f\x1b\xf4\x5d\x7a\x1d\xf4\xdd\x1f\x1c\xc9\x64\xfe\x0c\x00\x00\xff\xff\xde\xb2\x48\x18\x9d\x44\x00\x00")

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

	info := bindataFileInfo{name: "assets/tmpl/config.tmpl", size: 17565, mode: os.FileMode(420), modTime: time.Unix(1511093025, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _assetsTmplHeaderTmpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xac\x94\x5d\x4f\xe3\x3a\x10\x86\xef\xfb\x2b\x22\xdf\x9e\x53\x9b\x7e\x1c\xda\xb3\x6a\x2a\xb1\xb0\x94\xc2\xf2\xb1\x50\x40\xec\x9d\xeb\x4c\x12\x87\xd8\x0e\x9e\x49\xd3\x6e\xd5\xff\xbe\x6a\xca\xb6\x88\x05\x89\x0b\xee\x3c\x33\xf2\x3b\xf3\xce\x23\xcd\x72\x19\x41\xac\x2d\x04\x2c\x05\x19\x81\x67\xab\x55\x63\xb0\x7e\x0e\x1b\x41\x10\x04\x03\x03\x24\x03\x95\x4a\x8f\x40\x21\xbb\x9d\x1c\x37\xfb\xec\x65\xc9\x4a\x03\x21\x9b\x69\xa8\x0a\xe7\x89\x05\xca\x59\x02\x4b\x21\xab\x74\x44\x69\x18\xc1\x4c\x2b\x68\xd6\xc1\xbf\x81\xb6\x9a\xb4\xcc\x9b\xa8\x64\x0e\x61\x8b\x0d\x1b\x1b\x25\xd2\x94\xc3\xf0\x7c\x31\x75\x34\x10\x9b\xe0\xb9\x92\x6b\xfb\x18\x78\xc8\x43\x86\xb4\xc8\x01\x53\x00\x62\x41\xea\x21\x0e\x59\x4a\x54\xe0\x17\x21\x8c\x9c\xab\xc8\xf2\xa9\x73\x84\xe4\x65\xb1\x0e\x94\x33\x62\x9b\x10\x1d\xde\xe1\x3d\xa1\x10\x77\x39\x6e\xb4\xe5\x0a\x91\x05\xda\x12\x24\x5e\xd3\x22\x64\x98\xca\x4e\xbf\xdb\xfc\x7a\xf7\xa0\xf5\xcd\xf8\x18\xce\x5a\xd1\xc8\x9c\x5e\x1f\x3c\x2e\x54\x79\x72\x70\x72\x9d\x74\xda\x97\xe6\x56\x55\x55\xcf\xd9\xce\xf5\x43\x94\x74\xef\xe4\x3f\x57\xe6\x66\x82\xbf\xc4\xd9\x7e\x7f\x36\x8d\xbe\x65\x69\xb7\x64\x81\xf2\x0e\xd1\x79\x9d\x68\x1b\x32\x69\x9d\x5d\x18\x57\xe2\x9f\xc5\x7d\xc4\x94\x8a\x6c\x86\x5c\xe5\xae\x8c\xe2\x5c\x7a\xa8\x1d\xc9\x4c\xce\x45\xae\xa7\x28\x62\x67\xa9\x29\x2b\x40\x67\x40\x74\xf9\x3e\xef\xd4\xf6\x5e\xa6\xb7\x0e\x3f\xd0\x55\x48\x44\x20\xac\x35\x54\x89\xe4\xcc\xf3\xcf\xcd\x57\x54\x5e\x17\x14\xa0\x57\xbb\x01\xd7\xb3\xf0\xc4\xb9\x24\x07\x59\x68\x7c\x35\x5f\xf6\x54\x82\x5f\x88\x16\x6f\xb5\x79\xf7\x39\xaa\x07\xca\x90\x0d\x07\x62\x23\x38\x7c\x5f\xfd\xa3\x4c\xb3\xd7\x48\xb3\x37\x89\x4e\xd4\x7f\xe3\x1f\x7a\xba\xd7\xee\x3d\xcd\x16\xd9\xcd\x79\x7c\x92\x5d\x9e\xcb\xef\x8f\x71\x79\x7f\x37\xff\x39\xbf\xbd\xb2\x87\xa7\x07\xbd\xbc\x6d\x0e\xef\x2f\xc6\xc5\xe8\x7f\x33\x3a\x3c\xea\x57\xa3\x8b\xb1\xba\x3a\xea\x4d\xe6\xf2\x7d\xa2\x3b\x2f\x9f\xc6\x76\x6b\xa8\x49\x32\x41\x6d\x8b\x92\xc4\x1e\xef\xf3\xbd\xb7\x2a\x2f\x09\xbf\xb5\xc8\x4f\xed\xf5\x37\xc0\xf7\x1b\x17\xb9\xa4\xd8\x79\xc3\xa9\xd2\x44\xe0\xeb\xc6\x95\x8e\x12\x20\xac\x29\x6d\x8f\x4a\x49\xf1\xfa\xa8\xec\x44\x07\x62\x73\x80\x96\x4b\xb0\xd1\x6a\xd5\xf8\x1d\x00\x00\xff\xff\xd2\x55\x52\x1a\xa4\x04\x00\x00")

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

	info := bindataFileInfo{name: "assets/tmpl/header.tmpl", size: 1188, mode: os.FileMode(420), modTime: time.Unix(1503507913, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _assetsTmplIndexTmpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xdc\x58\xdd\x6f\xdb\x36\x10\x7f\xcf\x5f\x71\xe5\xfa\xb0\x01\x55\x84\xb6\xc0\x1e\x02\x5a\x40\xd0\x0c\x6b\x87\x66\x0b\x9a\x76\xc0\x1e\x29\xf1\x2c\x31\xa1\x48\x8d\xa4\x9c\x79\x82\xff\xf7\x81\x92\xe5\xc8\x8a\x3e\x9c\x8f\xae\x41\xfb\xd0\xd8\xe6\x7d\xfc\xee\xee\x77\xa7\x13\xab\x8a\xe3\x52\x28\x04\x22\x14\xc7\x7f\xc8\x66\x73\x44\x5f\x9c\xfd\xf1\xee\xf3\x5f\x17\xbf\x40\xe6\x72\x19\x1d\x51\xff\x07\x24\x53\xe9\x82\xa0\x22\xd1\x51\x55\x39\xcc\x0b\xc9\x1c\x02\xc9\x90\x71\x34\x04\x8e\xbd\x62\xac\xf9\x3a\x3a\x02\x00\xa0\x5c\xac\x20\x91\xcc\xda\x05\x49\xb4\x72\x4c\x28\x34\xa4\x39\xf3\xff\xba\x26\x14\x5b\xc5\x6c\x6b\x62\x27\xd0\x35\x70\x55\xe6\xb1\x76\x46\xab\x8e\x81\x5a\x26\x7b\x1d\x9d\xaf\x63\xed\x68\x98\xbd\xee\x1f\xbd\x69\xb5\x25\x32\x4e\x22\x56\x3a\x9d\x33\x27\x12\x26\xe5\x1a\x12\x2d\x25\x26\x0e\x98\xe2\x10\x1b\xcd\x78\xc2\xac\xff\xb6\x86\x6b\xa1\xb8\x05\xbd\x04\xa1\x96\xda\x78\x0d\xad\x60\xa9\x0d\xac\x75\x49\xc3\xec\xcd\xad\x1b\x1a\x72\xb1\x8a\x3a\x88\x4b\xd9\xba\x54\x6c\x05\x8a\xad\x82\x42\x48\x69\xeb\x4f\x57\xa5\x75\x62\x29\x90\xef\x85\x40\xa5\x00\xa3\x25\x2e\x48\x61\xd0\xa2\x72\xb5\x3b\xd2\x9a\x61\x89\x13\x2b\x24\x11\x65\xc0\x99\x63\x81\xd3\x69\xea\x85\x1d\x8b\x09\x64\x06\x97\x0b\xf2\x43\xa6\x73\x24\xd1\x7b\x9d\x23\x0d\x59\x44\x43\x29\xba\x0e\xaa\xca\x30\x95\x22\xbc\xbc\xc6\xf5\x2b\x78\xb9\x62\x12\x4e\x16\x70\xfc\xae\x09\x5f\x68\x75\xce\x8a\xcd\x66\x1e\xd1\x24\x84\xaa\xf2\xe6\x37\x1b\x12\xb5\x9f\x86\xa1\xa0\xe2\x1d\x5f\x34\x2c\x65\x34\x5c\x6f\xc7\xe2\xc0\x93\x06\x95\xeb\x57\xdc\x4b\x09\xbe\x20\x75\xd8\x5d\xf9\x82\x29\x84\x25\xe3\x08\x42\x41\x9b\xb8\x3d\xdd\x5a\xff\x45\x10\xc0\xe7\x1b\xe1\x1c\x1a\xb8\xcd\x02\x7c\x14\xd6\x41\x10\x0c\x28\x74\x60\x79\x17\x12\xea\xff\x03\xcf\x8e\x01\xfb\x83\x2a\x81\xef\x11\xa1\xd2\x11\xf9\x86\xae\x6f\xf7\x55\x9c\x70\xb2\xa9\x7c\x93\xe4\xcc\xb9\xc2\x9e\x84\xa1\xbb\x41\x74\x1c\x93\xeb\x63\xd7\x44\x71\x9c\xe8\x9c\x44\x77\x43\xb2\x4d\x0d\xb2\xb7\x23\x20\x1b\xf2\x1e\x86\xdf\xb7\xf5\x04\xf8\xaa\x12\xcb\x29\x4e\xf5\xcc\xfb\xb2\x8f\x1d\xc2\x43\x38\x3b\xe0\x63\x9f\x7a\x23\x42\x6d\x6e\xab\xca\x3b\xe9\xf3\x77\xda\x41\x38\xe7\xa1\x4f\xf7\xbb\x16\xa6\xf2\x50\x55\x28\x2d\x4e\xa9\x63\x1e\xfd\xae\x5d\x26\x54\x0a\x4e\x83\xcd\xf4\x0d\x0d\x31\x9f\xb4\x38\x8a\x67\x84\x0c\xfd\x01\xb7\xfb\xdd\x37\xd1\x9f\xc2\xfa\xc6\x39\xbd\xf8\x00\x06\x6d\x29\xbf\x69\xfb\x34\x0c\xfc\x90\xb3\x14\x4f\x15\x93\x6b\x2b\xec\x19\x73\x93\xf9\x1b\x69\xb8\x5f\xb5\x4e\x25\x76\xa3\xfb\xd4\x44\xf7\x63\x55\x0d\x39\xf8\x69\xbc\xc5\xe0\xa0\x3a\xde\x13\xc7\xac\xbb\xfb\x16\x19\x1e\xde\xf1\x75\x3a\xbe\x7c\xfa\x38\x15\x5e\xc7\xb2\xd1\x37\x13\x26\xfb\xd2\x89\x96\x41\xce\x83\x9f\x67\x54\x60\xbf\x91\x1b\x4c\x97\xba\x34\x09\xfa\x86\xa6\x22\x4f\xc1\x9a\xe4\xf6\xac\xc6\xbb\x7b\x72\x88\x3c\x0d\x0c\xda\x42\x2b\xeb\x1f\x18\xaf\x80\x49\xb7\x20\xb5\x20\x30\x5f\xe8\x7f\x91\x43\xbc\x86\x3b\xe5\x20\xd1\x01\x33\x62\x34\xe1\x8f\x0d\xb8\x30\x18\xd1\x44\x73\x8c\xfa\xac\x6c\x58\xe2\x27\x58\x7d\x4c\x43\x2f\xfa\x18\x98\x13\xc7\x93\xe4\xbe\xd7\x80\x7a\xea\xe1\xf4\x1e\x99\x74\x59\x92\x61\x72\x0d\x8e\xc5\x12\x67\x87\x53\x2d\xd5\xa5\xc2\x48\xa7\x34\xd6\xba\x4a\x8d\x83\xc0\x3a\x23\x0a\xe4\xdb\x6f\xb1\x36\x1c\x4d\x6f\xdb\xeb\x9b\xf2\x73\x6d\xa6\x36\xce\x1c\xc0\x06\xc7\xa3\x0b\xa3\x13\xb4\x96\x86\x6e\xc6\x62\xab\xd0\xcd\x50\x3b\x59\xe6\x74\x69\x38\x05\x87\x86\x33\x01\x51\x77\xfb\x7a\x30\x2e\x73\x60\xc0\xed\xc2\x73\x26\x8c\xdf\xe3\xcf\xd1\x5a\xdf\xb3\x7e\x8f\x43\x85\xe6\xf0\x44\xcc\x0a\x41\x3b\xef\xf0\x6f\x38\xde\xba\x6d\xdc\x9c\x9d\x5f\x3a\xe6\x4a\x0b\xce\x94\x53\x33\x7e\xcf\xa5\x2d\x98\xaa\x37\xd8\xed\xfe\x16\xf0\x3c\x90\x5b\xd4\x81\xad\xed\x11\xb0\x6e\xed\x97\xec\x44\x4b\x6d\x4e\x62\x59\x22\x89\x4e\xeb\x95\x96\x86\x5e\xff\x50\xd0\x33\xcf\x9e\x47\xe2\xaa\xe9\x7d\xe9\x74\x71\x4f\x50\x93\x6b\xd1\x0e\xcf\x23\xd9\x08\x0f\x21\xd3\x17\x8b\xe6\xff\xa5\x90\xf7\x68\x9f\x82\x45\xa5\x45\xf3\x3c\x79\x74\x00\xb2\xef\x8f\x49\x17\x68\x84\xe6\x22\x81\xdf\x74\xfc\x95\x89\xd4\xba\x7a\x0a\x16\x15\x5b\x5b\xc1\x95\x8e\x9f\x15\x89\xe6\x81\x7d\x3f\x1c\xba\x94\x2c\xb9\xfe\xea\x63\xa8\xf6\xd2\x3a\x79\x0c\x77\xac\x37\xf4\xbc\x46\xcf\x21\x90\x9e\x35\x5f\x68\x38\xb1\x2b\xd1\xb0\xde\x31\x47\x57\xe2\x99\x9f\x1e\x70\xbd\xb2\xbb\x75\xdb\xdd\xf4\x0d\xde\xbc\x0d\xdd\xb7\xb1\x9d\xe4\xb6\x93\x53\x23\x38\xb9\x73\xf1\x72\xe7\x65\x6a\x10\x78\xff\x02\xf1\xf6\x0d\x60\xfb\x99\x86\x4d\xda\x68\xd8\xdc\x5c\xb7\x2a\xff\x05\x00\x00\xff\xff\xf9\xc9\xc5\x24\xe5\x16\x00\x00")

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

	info := bindataFileInfo{name: "assets/tmpl/index.tmpl", size: 5861, mode: os.FileMode(420), modTime: time.Unix(1511092312, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _assetsTmplNavbarTmpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xc4\x96\xcf\x6e\xdb\x3c\x0c\xc0\xef\x79\x0a\x41\xdf\xd9\xd1\xfd\x43\x12\x0c\x18\x86\x5d\xda\x1e\xb6\xee\x5e\xda\xa2\x1d\xa2\x8a\xe4\x5a\xb4\xdb\xc1\xf0\xbb\x0f\x96\xec\xd4\x75\x93\x22\x4b\x83\xec\x92\x48\x14\xff\x88\x3f\x91\x96\xda\x56\x63\x4e\x16\x85\xb4\xd0\xa4\x50\xc9\xae\x5b\xac\x2c\x34\x22\x33\xe0\xfd\x7a\x90\x8a\xf8\x97\x68\xcc\xa1\x36\x3c\x4e\xc9\x36\x58\x79\x1c\xa7\x39\xbd\xa0\x4e\xd8\x95\x72\xb3\x10\x42\x88\x95\xa6\xbd\x9f\xcc\x59\x06\xb2\x58\x25\xb9\xa9\x49\x0f\x1a\x73\xad\xc1\xd1\x16\x41\x63\x35\xd1\x09\x7a\x69\xcd\xec\xec\x4c\x95\x5d\x51\x18\x14\x99\x33\x06\x4a\x8f\x5a\x0a\x0d\x0c\x83\xb8\x0f\x1b\xe5\xa3\x18\xaa\x02\x79\x2d\xff\x8b\xd6\xb7\x68\x6b\x2f\x05\x54\x04\x09\xbe\x94\x60\x35\xea\xb5\xcc\xc1\x78\x9c\x05\x0f\x1b\xf0\x25\xec\xc3\x53\xe6\x6c\xd2\xf3\xda\xac\x54\x2f\xbf\xa6\xfa\x4a\x45\x14\x33\x29\xcc\xd0\xa4\x15\x58\x2d\xc5\xb6\xc2\x7c\x2d\x95\xdc\xdc\xfe\x4e\x1d\xaf\x14\x4c\xd0\x2b\x4d\xcd\x87\x27\x31\xf2\x13\xaf\x20\x49\x8f\xab\x11\xdf\x6c\x17\xb5\x99\xb8\x18\x2b\xc3\x42\x73\x08\xa8\xa1\x89\x6e\x42\x8c\xbb\xb6\xa5\x5c\xe0\x93\x58\xde\x05\xc3\x3b\xd8\xa1\x78\xf8\xea\x6c\x4e\xc5\x43\xd7\x41\xc6\xd4\x60\xdb\xa2\xd5\x5d\x77\xc0\xe1\x3b\x0c\x89\x21\xfb\xb8\x47\x90\x05\x3f\x4a\x6e\xa2\xc3\x37\x28\x5e\x91\x18\x9a\xd3\xae\xcd\x66\x71\x5a\x92\xe3\xb0\xa2\x62\xcb\x07\x36\x18\xd2\x5b\xde\x3f\x13\x33\x86\xe4\xba\xee\x14\x2a\xc7\x53\x8d\x99\x6d\x99\x4b\xff\xbf\x52\x1c\x1d\x2f\x33\xb7\x53\x6d\xfb\x36\x8e\x1c\x9d\xa6\x6c\x45\xca\x36\x31\x45\xf8\xf3\x2e\x23\x30\x49\x5f\x71\x61\x3e\xf8\x38\x12\x52\xcc\x2b\x35\x07\x91\x43\xef\x2b\x87\xd1\xb4\x1f\xc6\x93\xda\x97\xaf\xf8\x32\xdf\xce\xe1\x84\x4e\x3b\x91\x48\x12\x8d\xbf\x30\x3e\x28\x4b\xbf\x9c\x30\xbc\x2e\xb2\x12\xbc\x9f\x32\xbb\x00\xa1\xbe\x4d\x16\x47\x8a\xf0\xa7\x81\xec\xf1\x1e\x61\x77\x19\x86\x6d\x1b\x3d\xfe\xfa\x71\x73\x7a\xad\xf9\xde\xe2\x1c\x6c\xc1\xf0\x58\x9d\x7d\x94\xd9\xdf\x33\xbc\x58\x95\x29\xa8\x79\xab\x62\xca\x57\xe3\x73\xb5\xa2\xfa\x44\xdb\x65\xce\x7a\x67\x70\xa9\xb1\x41\xe3\x4a\xac\xfc\xb2\x70\xae\x30\x18\xbe\x63\x50\x92\xef\x7f\x54\x43\x9e\x9c\x1d\x96\x7a\x69\x58\x7e\xaa\x1d\x83\x3f\x15\x68\x34\x3e\x87\x68\xb4\xec\x47\xb1\x7d\xbe\x87\xf9\x37\x0b\xa9\x41\x3d\xb9\x99\x42\xb9\x0c\xd4\xf7\x17\xd5\x85\xe0\x7f\x82\x71\x41\xbc\xad\xd3\x40\x8c\x9e\x81\xe1\x11\xd4\xae\x7f\x0e\x9c\x0c\x2e\xd8\x9f\x05\x2e\x58\x1e\xe8\xd5\x7f\x40\x43\x19\x57\xb8\x9a\xc7\x7b\x52\xbe\x7b\x29\x9c\xd3\x6a\x54\xd8\xc4\xd5\xfc\xfa\x0d\xba\x09\x41\x3e\x93\xde\xec\xc9\x31\x79\xa6\x0d\xc3\x95\xb2\xd0\x6c\x16\x63\x37\xfe\x09\x00\x00\xff\xff\x4a\xda\x13\x62\xc5\x0b\x00\x00")

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

	info := bindataFileInfo{name: "assets/tmpl/navbar.tmpl", size: 3013, mode: os.FileMode(420), modTime: time.Unix(1511094162, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _assetsTmplSetupTmpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xe4\x56\x4f\x6b\xdb\x4e\x10\xbd\xfb\x53\xcc\x6f\xef\xf6\x92\x1c\x7f\x48\x82\xe2\xf4\x10\x42\x68\xc1\xb9\xf4\x64\x46\xda\xb1\xb5\xb5\xb4\xbb\xec\x8e\x1c\x84\xd0\x77\x2f\x92\xec\xc4\x76\x64\xa7\xa4\x14\x07\x7a\xb1\x35\xcc\xcc\xd3\x9b\x37\x7f\x50\xd3\x28\x5a\x69\x43\x20\x02\x71\xe5\x44\xdb\x4e\xa2\xff\xee\xbe\xcd\x9f\x7e\x7c\xff\x0a\x39\x97\x45\x32\x89\xba\x3f\x28\xd0\xac\x63\x41\x46\x24\x13\x00\x80\xa6\x61\x2a\x5d\x81\x4c\x20\x72\x42\x45\x5e\xc0\xac\x6d\x7b\x5f\x94\x5a\x55\x0f\x61\xbd\xa9\xf4\x16\xb2\x02\x43\x88\x45\x66\x0d\xa3\x36\xe4\xc5\xab\xff\x14\xce\xe0\x36\xc5\x03\xb8\x31\x9c\x9f\x55\x99\x5a\xf6\xd6\x9c\xe0\xf4\x71\xf9\x4d\xb2\xe8\x8a\x89\x64\x7e\x33\xe6\xbe\xdd\xa3\x14\x84\x4a\x24\x4f\xb9\x0e\xa0\x03\x70\x4e\xb0\xd2\x3e\x30\x04\x26\x07\x6c\xa1\x0a\x04\x8f\x75\x6a\x39\x92\xf9\xed\x08\xd2\xca\xfa\x12\xb4\x8a\x77\xda\x01\x66\xac\xad\x89\x85\xec\x6d\x29\xa0\x24\xce\xad\x8a\x85\xb3\x81\x05\x90\xc9\xb8\x76\x14\x8b\xb2\x2a\x58\x3b\xf4\x2c\x3b\x84\xa9\x42\x46\x91\x44\xbd\x31\xf2\x16\x6d\x5c\xc5\xd0\x39\x5f\x5e\x34\xc0\x84\x2a\x2d\x35\x0b\xd8\x62\x51\x51\x2c\x16\xb8\x25\xb1\x2f\x2d\x65\x03\x29\x9b\xa9\xf3\xba\x44\x5f\xf7\xcf\xc5\xfa\x44\xad\x48\x2a\xbd\x4d\x26\x27\x9d\xd0\x2b\x98\x3d\x52\x08\xb8\xa6\x0b\x1d\xc0\x82\x3c\x43\xff\x3b\x55\x68\xd6\x5d\x4b\x9b\xe6\x35\x71\x87\x7d\x0c\x4d\x46\xb5\xed\xe4\xb7\x30\xb5\x59\xd9\xb1\xe6\x22\xe4\x9e\x56\xb1\xc8\x99\x5d\xf8\x5f\x4a\x74\x2e\xcc\xf8\x59\x33\x93\x9f\x65\xb6\x14\xc0\xe8\xd7\xc4\xb1\x58\xa6\x05\x9a\x8d\x48\xe6\x85\xce\x36\x90\x93\xa7\xae\xa7\x99\x27\x64\x92\x59\x4e\xd9\x06\x6a\x5b\x79\xb0\xcf\x06\x9e\x06\x00\xf8\xe2\x5c\x24\x71\x54\xa5\x73\xa4\x19\xd3\x82\xa6\x9e\x82\xb3\x26\xe8\x2d\x8d\x91\xee\x63\x8e\x12\x60\x48\x0b\xec\xb5\x23\xb5\xb3\x52\xeb\x15\x79\x52\x23\x10\x03\xcc\xf1\x62\xbd\xf5\xfb\xf3\xce\x21\x40\x25\xfb\x4a\xe7\xd6\x84\xaa\x24\x0f\x0f\x54\x47\x92\xd5\xfb\x99\x17\x03\xe0\xe2\xa4\x3a\x0c\xe1\xd9\x7a\x25\xc0\x60\x49\x3b\xe7\xbe\x6b\xcb\x6c\x47\x65\xb9\xa1\xfa\x65\x9a\x9b\x66\xb6\xa3\xba\x67\xfa\x40\x75\xdb\x9e\x91\xe6\x85\xc2\xc5\x4a\x22\x79\x49\xa1\x8f\xc9\xb7\xa0\xcc\x13\x7f\x16\x05\x43\xcf\xe6\x82\x88\x03\xdd\xbf\xa6\x63\x24\xcf\xcc\x68\x24\xfb\x11\x7f\xff\xfe\xfc\xf9\x3d\xd0\xb3\x50\x60\xb6\xe9\x8e\x41\x7f\x1d\x3e\x76\x11\x16\x1d\xc6\x3f\x70\x0f\x86\x3a\xe7\x85\x26\xc3\x70\x7f\x77\xa5\x41\xee\x3b\xb6\xcc\x7a\x16\x4b\xad\x0e\x07\xb8\x27\x38\xf0\xbb\xbf\xbb\xfe\x01\x38\xd2\xeb\xaa\xcb\x7f\xa4\xd9\xdb\xc5\x3f\xd0\xed\x53\x2f\xfd\xc1\x3a\x45\x72\xc0\x89\xe4\xf0\xe5\xb9\xff\x60\xf8\x15\x00\x00\xff\xff\x15\x41\x24\x68\xa5\x0a\x00\x00")

func assetsTmplSetupTmplBytes() ([]byte, error) {
	return bindataRead(
		_assetsTmplSetupTmpl,
		"assets/tmpl/setup.tmpl",
	)
}

func assetsTmplSetupTmpl() (*asset, error) {
	bytes, err := assetsTmplSetupTmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "assets/tmpl/setup.tmpl", size: 2725, mode: os.FileMode(420), modTime: time.Unix(1503659531, 0)}
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
	"assets/tmpl/config/action.tmpl": assetsTmplConfigActionTmpl,
	"assets/tmpl/config/filter.tmpl": assetsTmplConfigFilterTmpl,
	"assets/tmpl/config.tmpl": assetsTmplConfigTmpl,
	"assets/tmpl/header.tmpl": assetsTmplHeaderTmpl,
	"assets/tmpl/index.tmpl": assetsTmplIndexTmpl,
	"assets/tmpl/navbar.tmpl": assetsTmplNavbarTmpl,
	"assets/tmpl/setup.tmpl": assetsTmplSetupTmpl,
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
			"config": &bintree{nil, map[string]*bintree{
				"action.tmpl": &bintree{assetsTmplConfigActionTmpl, map[string]*bintree{}},
				"filter.tmpl": &bintree{assetsTmplConfigFilterTmpl, map[string]*bintree{}},
			}},
			"config.tmpl": &bintree{assetsTmplConfigTmpl, map[string]*bintree{}},
			"header.tmpl": &bintree{assetsTmplHeaderTmpl, map[string]*bintree{}},
			"index.tmpl": &bintree{assetsTmplIndexTmpl, map[string]*bintree{}},
			"navbar.tmpl": &bintree{assetsTmplNavbarTmpl, map[string]*bintree{}},
			"setup.tmpl": &bintree{assetsTmplSetupTmpl, map[string]*bintree{}},
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

