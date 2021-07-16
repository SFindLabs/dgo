package uploadfile

import (
	kpkg "dgo/pkg"
	"dgo/pkg/uploadfile/obs"
	"dgo/pkg/uploadfile/oss/oss"
	"fmt"
)

//华为云
func UploadToObs(filePath, obsPath, obsBucket, obsUrl string) (string, error) {
	// 创建ObsClient结构体
	obsClient, err := obs.New(kpkg.ObsAK, kpkg.ObsSK, kpkg.ObsEndpoint)
	if err != nil {
		return "", err
	}
	if obsBucket == "" {
		obsBucket = kpkg.ObsBucket
	}
	input := &obs.PutFileInput{}
	input.Bucket = obsBucket
	input.Key = obsPath
	input.SourceFile = filePath // localfile为待上传的本地文件路径，需要指定到具体的文件名
	_, err = obsClient.PutFile(input)
	if err != nil {
		return "", err
	}
	if obsUrl == "" {
		obsUrl = kpkg.ObsUrl
	}
	return fmt.Sprintf("%s/%s", obsUrl, obsPath), nil
}

func ObsSignUrl(expire int, obsPath, obsBucket string) (string, error) {
	obsClient, err := obs.New(kpkg.ObsAK, kpkg.ObsSK, kpkg.ObsEndpoint)
	if err != nil {
		return "", err
	}
	if obsBucket == "" {
		obsBucket = kpkg.ObsBucket
	}
	input := &obs.CreateSignedUrlInput{}
	input.Expires = expire
	input.Method = obs.HttpMethodGet
	input.Bucket = obsBucket
	input.Key = obsPath
	output, err := obsClient.CreateSignedUrl(input)
	if err == nil {
		return output.SignedUrl, nil
	}
	return "", err
}

//阿里云
func UploadToOss(filePath, ossPath string) (string, error) {
	ossClient, err := oss.New(kpkg.OssEndpoint, kpkg.OssAk, kpkg.OssSK)
	if err != nil {
		return "", err
	}

	bucket, err := ossClient.Bucket(kpkg.OssBucket)
	if err != nil {
		return "", err
	}

	err = bucket.PutObjectFromFile(ossPath, filePath)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("https://%s.%s/%s", kpkg.OssBucket, kpkg.OssInternetEndpoint, ossPath), nil
}
