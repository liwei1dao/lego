package sts

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/sts"
)

func newSys(options Options) (sys *STS, err error) {
	sys = &STS{options: options}
	err = sys.init()
	return
}

type STS struct {
	options Options
	client  *sts.Client
}

func (this *STS) init() (err error) {
	if this.client, err = sts.NewClientWithAccessKey(this.options.RegionId, this.options.AccessKeyId, this.options.AccessKeySecret); err != nil {
		return
	}

	return
}

func (this *STS) AssumeRole(roleArn, roleSessionName string) (auth *Authorization, err error) {
	var (
		request  *sts.AssumeRoleRequest
		response *sts.AssumeRoleResponse
	)
	//构建请求对象。
	request = sts.CreateAssumeRoleRequest()
	request.Scheme = "https"
	//设置参数。关于参数含义和设置方法，请参见《API参考》。
	request.RoleArn = roleArn
	request.RoleSessionName = roleSessionName
	//发起请求，并得到响应。
	response, err = this.client.AssumeRole(request)
	if err != nil {
		return
	}
	auth = &Authorization{
		Expiration:      response.Credentials.Expiration,
		AccessKeyId:     response.Credentials.AccessKeyId,
		AccessKeySecret: response.Credentials.AccessKeySecret,
		SecurityToken:   response.Credentials.SecurityToken,
	}
	return
}
