package uploadfile

import (
	kpkg "dgo/pkg"
	"dgo/pkg/uploadfile/oss/oss"
	"fmt"
)

//阿里云oss(初始化配置文件默认参数)

func NewOssUpload() *ossUpload {
	return &ossUpload{
		kpkg.OssAk,
		kpkg.OssSK,
		kpkg.OssEndpoint,
		kpkg.OssBucket,
	}
}

type ossUpload struct {
	ossAK       string
	ossSK       string
	ossEndpoint string
	ossBucket   string
}

func (osu *ossUpload) SetOssAK(ossAK string) {
	osu.ossAK = ossAK
}

func (osu *ossUpload) SetOssSK(ossSK string) {
	osu.ossSK = ossSK
}

func (osu *ossUpload) SetOssEndpoint(ossEndpoint string) {
	osu.ossEndpoint = ossEndpoint
}

func (osu *ossUpload) SetOssBucket(ossBucket string) {
	osu.ossBucket = ossBucket
}

/**	上传
 * 	filePath 待上传的本地文件路径，需要指定到具体的文件名
 *	ossPath oss目录保存路径
 *  isPrivate 是否私密的bucket
 */
func (osu *ossUpload) Upload(filePath, ossPath string, isPrivate bool) (string, error) {
	ossClient, err := oss.New(osu.ossEndpoint, osu.ossAK, osu.ossSK)
	if err != nil {
		return "", err
	}

	bucket, err := ossClient.Bucket(osu.ossBucket)
	if err != nil {
		return "", err
	}

	err = bucket.PutObjectFromFile(ossPath, filePath)
	if err != nil {
		return "", err
	}
	if isPrivate {
		return ossPath, nil
	}
	return fmt.Sprintf("https://%s.%s/%s", osu.ossBucket, osu.ossEndpoint, ossPath), nil
}

/**	获取限时可访问地址
 * 	expire 设置多少秒后过期
 *	ossPath oss目录保存路径
 */

func (osu *ossUpload) GetSignUrl(expire int64, ossPath string) (string, error) {
	ossClient, err := oss.New(osu.ossEndpoint, osu.ossAK, osu.ossSK)
	if err != nil {
		return "", err
	}
	bucket, err := ossClient.Bucket(osu.ossBucket)
	if err != nil {
		return "", err
	}
	signedURL, err := bucket.SignURL(ossPath, oss.HTTPGet, expire)
	if err != nil {
		return "", err
	}
	return signedURL, err
}
