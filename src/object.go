package src

type DtsInfoRecord struct {
	JobId string `db:"jobId"`
	//JobName         string `db:"jobName"`
	//AppId           int64  `db:"appId"`
	//Status          int    `db:"status"`
	//SrcSetId        int    `db:"srcSetId"`
	//SrcDatabaseType string `db:"srcDatabaseType"`
	//SrcAccessType   string `db:"srcAccessType"`
	SrcInfo_DB []byte `db:"srcInfo"`
	//DstSetId        int    `db:"dstSetId"`
	//DstDatabaseType string `db:"dstDatabaseType"`
	//DstAccessType   string `db:"dstAccessType"`
	DstInfo_DB []byte `db:"dstInfo"`
	CreateTime int64  `db:"createTime"`
	UpdateTime int64  `db:"updateTime"`
}

type InsInfo struct {
	RegionId         int
	InstanceId       int64
	SerialId         string
	SetId            int
	Addr             string
	SNatIp           string
	Password         string
	Supplier         string
	RedisPsyncPrefix string
}

type Flow struct {
	Id                    int64  `db:"id"`
	Name                  string `db:"name"`
	Status                int    `db:"status"`
	Res                   string `db:"res"`
	ErrMsg                string `db:"errMsg"`
	CurrentTask           string `db:"currentTask"`
	CurrentTaskStatus     int    `db:"currentTaskStatus"`
	CurrentTaskRetryTimes int    `db:"currentTaskRetryTimes"`
	Result                []byte `db:"result"`
	Context               []byte `db:"context"`
	Locked                int    `db:"locked"`
	TaskInfo              []byte `db:"taskInfo"`
	CreateTime            int64  `db:"createtime"`
	UpdateTime            int64  `db:"updatetime"`
	LockTime              int64  `db:"lockTime"`
	LastCall              int64  `db:"lastCall"`
	NextCall              int64  `db:"nextCall"`
}

type DstInfo struct {
	RegionId         int
	InstanceId       int64
	SerialId         string
	SetId            int
	Addr             string
	Password         string
	Index            int64
	RealAddr         string
	RoutePara        string
	AppId            int64
	ConfigEpoch      int64
	Supplier         string
	RedisPsyncPrefix string
}

type DstCcInfo struct {
	InstanceId     int64
	SerialId       string
	RedisPassword  string
	TendisPassword string
	UserPassword   string
	SoftVersion    string

	MaxMemory int64

	Interface    []Addr
	Cache        []CacheNode
	Tendis       []CacheNode
	Pmedis       []CacheNode
	InstanceType int
	L5ModId      int
	L5CmdId      int

	Locker     string
	ResourceId string
}

type CacheNode struct {
	Master Addr
	Slaves []Addr
}

type Addr struct {
	Ip        string
	Port      int
	AdminPort int
}

type RedisOpt struct {
	Index            int
	SetId            int
	RegionId         int
	InstanceId       int64
	SerialId         string
	RealAddr         string
	Addr             string
	Password         string
	RoutePara        string
	AppId            int64
	ConfigEpoch      uint64
	Supplier         string
	RedisPsyncPrefix string
}

type Locker struct {
	state int32
}

type DataType map[string]interface{}
type SetValueListener func(data DataType, key string, value interface{}) error

type Dict struct {
	data DataType
}

type Task struct {
	/**
	 * claim as private, forbidden visit directly
	 */
	flow *Flow

	/**
	 * expect order
	 */
	flowIndex int
	/**
	 * run order actually
	 */
	runIndex int

	status   int
	errorMsg string
	result   []byte

	taskName string

	maxRetryTimes int
	retryTimes    int

	//handler TaskHandler
}

type SyncerConfig struct {
	MaxSendBufferSize    int //
	MaxSendExpiredTime   int64
	MaxAofRecvBufferSize int
	MaxRecvErr           int

	IgnoreCheckEmptyIfFullerSync int

	AppendFile string

	ResumeOption        string
	SupportProxy        string
	ExpireKeyTaskError  string
	StartSyncerErrorMsg string
	DstDbType           string
	Supplier            string
	RedisPsyncPrefix    string
}

type Record struct {
	Id            int64  `db:"id"`
	JobId         string `db:"jobId"`
	Index         int    `db:"index"`
	SyncerVersion int64  `db:"syncerVer"`
	Token         string `db:"token"`
	Host          string `db:"host"`
	Status        int    `db:"status"`
	SrcIp         string `db:"srcIp"`
	SrcRealIP     string `db:"srcRealIp"`
	SrcAddr       string `db:"srcAddr"`
	SrcPassword   string `db:"srcPassword"`
	DstL5modid    int    `db:"dstL5modid"`
	DstL5cmdid    int    `db:"dstL5cmdid"`
	DstAddr       string `db:"dstAddr"`
	DstPortOffset int    `db:"dstPortOffset"`
	DstPassword   string `db:"dstPassword"`
	CurrentAction string `db:"currentAction"`
	CreateTime    int64  `db:"createTime"`
	UpdateTime    int64  `db:"updateTime"`
	LastErrMsg    string `db:"lastErrMsg"`

	MasterId  string `db:"masterId"`
	Offset    int64  `db:"offset"` //see DstOffset
	SrcOffset int64  `db:"srcOffset"`
	DstOffset int64  `db:"dstOffset"` //should be same with Offset above at final,
	// this value retrieve from the src info
	RdbSize int64 `db:"rdbSize"`
	RdbLeft int64 `db:"rdbLeft"`

	Config    SyncerConfig
	ConfigStr []byte `db:"config"`
}
