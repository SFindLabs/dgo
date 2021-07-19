package uploadfile

import (
	kpkg "dgo/pkg"
	"dgo/pkg/uploadfile/obs"
	"fmt"
)

//华为云obs(初始化配置文件默认参数)

func NewObsUpload() *obsUpload {
	return &obsUpload{
		kpkg.ObsAK,
		kpkg.ObsSK,
		kpkg.ObsEndpoint,
		kpkg.ObsBucket,
		kpkg.ObsUrl,
	}
}

type obsUpload struct {
	obsAK       string
	obsSK       string
	obsEndpoint string
	obsBucket   string
	obsUrl      string
}

func (obu *obsUpload) SetObsAK(obsAK string) {
	obu.obsAK = obsAK
}

func (obu *obsUpload) SetObsSK(obsSK string) {
	obu.obsSK = obsSK
}

func (obu *obsUpload) SetObsEndpoint(obsEndpoint string) {
	obu.obsEndpoint = obsEndpoint
}

func (obu *obsUpload) SetObsBucket(obsBucket string) {
	obu.obsBucket = obsBucket
}

func (obu *obsUpload) SetObsUrl(obsUrl string) {
	obu.obsUrl = obsUrl
}

/**	上传
 * 	filePath 待上传的本地文件路径，需要指定到具体的文件名
 *	obsPath obs目录保存路径
 *  isPrivate 是否私密的bucket
 */
func (obu *obsUpload) Upload(filePath, obsPath string, isPrivate bool) (string, error) {
	// 创建ObsClient结构体
	obsClient, err := obs.New(obu.obsAK, obu.obsSK, obu.obsEndpoint)
	if err != nil {
		return "", err
	}
	input := &obs.PutFileInput{}
	input.Bucket = obu.obsBucket
	input.Key = obsPath
	input.SourceFile = filePath // localfile为待上传的本地文件路径，需要指定到具体的文件名
	_, err = obsClient.PutFile(input)
	if err != nil {
		return "", err
	}
	//私密的bucket后续需要用obsPath签名获取限时链接
	if isPrivate {
		return obsPath, nil
	}
	return fmt.Sprintf("%s/%s", obu.obsUrl, obsPath), nil
}

/**	获取限时可访问地址
 * 	expire 设置多少秒后过期
 *	obsPath obs目录保存路径
 */

func (obu *obsUpload) GetSignUrl(expire int, obsPath string) (string, error) {
	obsClient, err := obs.New(obu.obsAK, obu.obsSK, obu.obsEndpoint)
	if err != nil {
		return "", err
	}
	input := &obs.CreateSignedUrlInput{}
	input.Expires = expire
	input.Method = obs.HttpMethodGet
	input.Bucket = obu.obsBucket
	input.Key = obsPath
	output, err := obsClient.CreateSignedUrl(input)
	if err == nil {
		return output.SignedUrl, nil
	}
	return "", err
}
