package aliyun

type Region string

// Constants of region definition
const (
	Hangzhou    = Region("cn-hangzhou")
	Qingdao     = Region("cn-qingdao")
	Beijing     = Region("cn-beijing")
	Hongkong    = Region("cn-hongkong")
	Shenzhen    = Region("cn-shenzhen")
	Shanghai    = Region("cn-shanghai")
	Zhangjiakou = Region("cn-zhangjiakou")
	Huhehaote   = Region("cn-huhehaote")
	ChengDu     = Region("cn-chengdu")
	HeYuan      = Region("cn-heyuan")
	WuLanChaBu  = Region("cn-wulanchabu")
	GuangZhou   = Region("cn-guangzhou")

	APSouthEast1 = Region("ap-southeast-1")
	APNorthEast1 = Region("ap-northeast-1")
	APSouthEast2 = Region("ap-southeast-2")
	APSouthEast3 = Region("ap-southeast-3")
	APSouthEast5 = Region("ap-southeast-5")

	APSouth1 = Region("ap-south-1")

	USWest1 = Region("us-west-1")
	USEast1 = Region("us-east-1")

	MEEast1 = Region("me-east-1")

	EUCentral1 = Region("eu-central-1")
	EUWest1    = Region("eu-west-1")

	ShenZhenFinance     = Region("cn-shenzhen-finance-1")
	ShanghaiFinance     = Region("cn-shanghai-finance-1")
	ShanghaiFinance1Pub = Region("cn-shanghai-finance-1-pub")
	CnNorth2Gov1        = Region("cn-north-2-gov-1")
)
