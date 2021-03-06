package file

import (
	"bufio"
	"common/std/hash"
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

//dir, name := filepath.Dir(path), filepath.Base(path) //不包含尾分隔符，且会转换为对应平台的分隔符
//dir, name := filepath.Split(path) 				   //dir包含尾分隔符，同参数的分隔符

func GetExeDir() string {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	return filepath.ToSlash(dir)
}
func IsExist(path string) bool { //file or folder
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

//获取指定目录下(不递归子目录)的所有文件 --- names, err := filepath.Glob("csv/*.csv")
//获取指定目录及子目录下的所有文件，可以匹配后缀过滤 --- names, err := WalkDir("csv/", ".csv")
func WalkDir(dir, suffix string) (ret []string, err error) {
	//filepath.Walk【对软连接无效】
	//err := filepath.Walk(dir, func(filename string, fi os.FileInfo, err error) error {
	//	if err != nil {
	//		return err
	//	}
	//	if fi.IsDir() {
	//		return nil
	//	}
	//	if strings.HasSuffix(fi.Name(), suffix) {
	//		ret = append(ret, filepath.ToSlash(filename))
	//	}
	//	return nil
	//})
	err = _walkDir(dir, suffix, &ret)
	return
}
func _walkDir(dir, suffix string, ret *[]string) error {
	if len(dir) == 0 {
		return nil
	}
	if dir[len(dir)-1] != '/' {
		dir += "/"
	}
	if f, err := os.Open(dir); err == nil {
		list, err := f.Readdir(-1)
		f.Close()
		if err == nil {
			for _, fi := range list {
				name := dir + fi.Name()
				if fi.Mode()&(os.ModeDir|os.ModeSymlink) != 0 {
					if err = _walkDir(name, suffix, ret); err != nil {
						return err
					}
				} else if strings.HasSuffix(name, suffix) {
					*ret = append(*ret, name)
				}
			}
			return nil
		} else {
			return err
		}
	} else {
		return err
	}
}

func ReadLine(filename string, cb func(string)) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	for s := bufio.NewScanner(f); s.Scan(); {
		cb(s.Text())
	}
	f.Close()
	return nil
}

func CreateFile(dir, name string, flag int) (*os.File, error) {
	if err := os.MkdirAll(dir, 0777); err != nil {
		return nil, err
	}
	if dir[len(dir)-1] != '/' {
		dir += "/"
	}
	if file, err := os.OpenFile(dir+name, flag|os.O_CREATE, 0666); err != nil {
		return nil, err
	} else {
		return file, nil
	}
}

// ------------------------------------------------------------
// 计算文件md5
func CalcMd5(name string) uint32 {
	f, err := os.Open(name)
	if err != nil {
		return 0
	}
	md5hash := md5.New()
	io.Copy(md5hash, f)
	f.Close()
	return hash.StrHash(fmt.Sprintf("%x", md5hash.Sum(nil)))
}
