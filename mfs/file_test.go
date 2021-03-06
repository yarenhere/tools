package mfs

import (
	"crypto/md5"
	"fmt"
	"github.com/go-playground/assert/v2"
	"github.com/stretchr/testify/suite"
	"io"
	"io/fs"
	"testing"
)

const emptyFileMd5Hash = "d41d8cd98f00b204e9800998ecf8427e"

func TestZeroFile(t *testing.T) {
	var fileSize int64 = 100 * 1024 * 1024
	fo := NewZeroFile("/abc/d.txt", fs.ModePerm, fileSize)
	hash := md5.New()
	n, err := io.Copy(hash, fo)
	assert.Equal(t, nil, err)
	assert.Equal(t, n, fileSize)
	fo.Close()
	fileHash := fmt.Sprintf("%x", hash.Sum(nil))

	fileHash2nd, err := GetFileMd5Hash(fo)
	assert.Equal(t, nil, err)
	assert.Equal(t, fileHash, fileHash2nd)
}

func TestRandomFile(t *testing.T) {
	suite.Run(t, new(RandomFileSuite))
}

type ZeroFileSuite struct {
	suite.Suite
}

type RandomFileSuite struct {
	suite.Suite
}

func (suite *RandomFileSuite) TestRandomFileClose() {
	fo := NewMockRandomFile("/test_dir/a.txt", fs.ModePerm, 25*1024*1024)
	var hash string
	var err error

	hash, err = GetFileMd5Hash(fo)
	suite.Equal(nil, err)
	fileHash := hash

	hash, err = GetFileMd5Hash(fo)
	assert.Equal(suite.T(), nil, err)
	assert.Equal(suite.T(), emptyFileMd5Hash, hash)

	fo.Close()
	hash, err = GetFileMd5Hash(fo)
	assert.Equal(suite.T(), nil, err)
	assert.Equal(suite.T(), fileHash, hash)
}

func (suite *RandomFileSuite) TestNameRandomFileName() {
	fo1 := NewMockRandomFile("file1", fs.ModePerm, 1*1024*1024)
	fo2 := NewMockRandomFile("file1", fs.ModePerm, 1*1024*1024)
	fo3 := NewMockRandomFile("file3", fs.ModePerm, 1*1024*1024)
	fo1Hash, err := GetFileMd5Hash(fo1)
	assert.Equal(suite.T(), nil, err)

	fo2Hash, err := GetFileMd5Hash(fo2)
	assert.Equal(suite.T(), nil, err)

	fo3Hash, err := GetFileMd5Hash(fo3)
	assert.Equal(suite.T(), nil, err)

	assert.Equal(suite.T(), fo1Hash, fo2Hash)
	assert.NotEqual(suite.T(), fo1Hash, fo3Hash)
}
