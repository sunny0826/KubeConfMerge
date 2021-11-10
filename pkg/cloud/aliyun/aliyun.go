package aliyun

import (
	cs20151215 "github.com/alibabacloud-go/cs-20151215/v2/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	"github.com/alibabacloud-go/tea/tea"
)

// ClusterInfo ack cluster info
type ClusterInfo struct {
	Name     string
	ID       string
	RegionID string
}

// GetClient get aliyun openapi client
func GetClient(accessKeyID, accessKeySecret string) (_result *cs20151215.Client, _err error) {
	config := &openapi.Config{
		AccessKeyId:     &accessKeyID,
		AccessKeySecret: &accessKeySecret,
		RegionId:        tea.String("cn-hongkong"),
	}
	_result = &cs20151215.Client{}
	_result, _err = cs20151215.NewClient(config)
	return _result, _err
}

// ListCluster list cluster info
func ListCluster(accessKeyID, accessKeySecret string) (clusters []ClusterInfo, _err error) {
	client, _err := GetClient(accessKeyID, accessKeySecret)
	if _err != nil {
		return nil, _err
	}
	describeClustersV1Request := &cs20151215.DescribeClustersV1Request{}
	v1, _err := client.DescribeClustersV1(describeClustersV1Request)
	if _err != nil {
		return nil, _err
	}
	var clusterList []ClusterInfo
	for _, info := range v1.Body.Clusters {
		clusterList = append(clusterList, ClusterInfo{
			Name:     *info.Name,
			ID:       *info.ClusterId,
			RegionID: *info.RegionId,
		})
	}
	return clusterList, _err
}

// GetKubeConfig get kubeConfig file
func GetKubeConfig(accessKeyID, accessKeySecret, clusterID string) (string, error) {
	client, _err := GetClient(accessKeyID, accessKeySecret)
	if _err != nil {
		return "", _err
	}
	describeClusterUserKubeconfigRequest := &cs20151215.DescribeClusterUserKubeconfigRequest{}
	res, _err := client.DescribeClusterUserKubeconfig(tea.String(clusterID), describeClusterUserKubeconfigRequest)
	if _err != nil {
		return "", _err
	}
	return *(res.Body.Config), _err
}
