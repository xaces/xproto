package xproto

// Obd obd数据
type Obd struct {
	Rpm          int     `json:"rpm"`   // 转速
	Speed        float32 `json:"spd"`   // 速度
	Acc          int     `json:"acc"`   // acc状态
	JQTemp       int     `json:"jtp"`   // 燃油进气温度
	JQMPos       int     `json:"jms"`   // 节气阀位置
	Voltage      float32 `json:"vol"`   // 电压
	Mileage      float32 `jons:"mile"`  // 里程
	TotalFuel    uint32  `json:"tfuel"` // 当前启动后总油耗
	InstanFuel   float32 `json:"ifuel"` // 瞬时油耗
	Efoa         uint8   `json:"efoa"`  // 剩余油量 %
	Airshed      uint8   `json:"airs"`  // 进气流量
	StressMPa    uint8   `json:"strpa"` // 进气管压力
	CoolantsTemp uint8   `json:"ctemp"` // 冷却液温度
	AirTemp      uint8   `json:"atemp"` // 进气温度
	MotorLimit   uint8   `json:"molmt"` // 发动机负荷计算值
	Position     uint8   `json:"pos"`   // 节气门绝对位置 %
}

// Gps 上报Location
type Location struct {
	Type      uint8   `json:"tp"`  // 定位类型
	Speed     float32 `json:"spd"` // 速度
	Angle     float32 `json:"agl"` // 角度
	Longitude float32 `json:"lng"` // 经度
	Latitude  float32 `json:"lat"` // 纬度
	Altitude  int     `json:"alt"` // 海拔
}

// Mileage 里程
type Mileage struct {
	Total uint32 `json:"tml"` // 总共
	Now   uint32 `json:"nml"` // 当前
}

// Oil 油耗
type Oil struct {
	Consumption float32 `json:"cus"`
	Remaining   float32 `json:"res"`
}

// Module 模块状态
type Module struct {
	Dial      uint8  `json:"dial"` // 移动网络(0:数据不存在，1:数据存在)
	Gps       uint8  `json:"gps"`  // 定位模块(0:数据不存在，1:数据存在)
	Wifi      uint8  `json:"wifi"` // WIFI模块(0:数据不存在，1:数据存在)
	Gsensor   uint8  `json:"gs"`   // GSensor(0:数据不存在，1:数据存在)
	Record    uint16 `json:"rec"`  // 每一位代表一个通道
	VideoLoss uint16 `json:"vlos"`
}

// Gsensor gsensor 状态
type Gsensor struct {
	X    float32 `json:"x"`
	Y    float32 `json:"y"`
	Z    float32 `json:"z"`
	Tilt float32 `json:"tilt"` // 倾斜
	Hit  float32 `json:"hit"`  // 撞击
}

// Mobile 移动网络状态
type Mobile struct {
	Signal uint8 `json:"sgl"` // 信号强度
	Type   uint8 `json:"tp"`  // 网络类型
}

// Disk  磁盘
type Disk struct {
	Type    uint8  `json:"tp"` // 类型
	Status  uint8  `json:"st"` // 状态
	Size    uint32 `json:"sz"` // 大小
	Residue uint32 `json:"re"` // 剩余
}

// 人数统计
type People struct {
	Up   uint32 `json:"up"` // 上车人数
	Down uint32 `json:"dw"` // 下车人数
}

// Status 状态数据
type Status struct {
	DeviceNo  string    `json:"deviceNo"`
	DTU       string    `json:"dtu"`    // DTU
	Flag      uint8     `json:"flag"`   // 0-实时 1-补传 2-报警开始Gps 3-报警结束Gps
	Acc       uint8     `json:"acc"`    // acc
	Location  Location  `json:"loc"`    // 位置信息
	Tempers   []float32 `json:"temps"`  // 温度
	Humiditys []float32 `json:"humis"`  // 湿度
	Mileage   Mileage   `json:"mile"`   // 里程
	Oils      []Oil     `json:"oils"`   // 油耗
	Module    Module    `json:"mod"`    // 模块状态
	Gsensor   Gsensor   `json:"gs"`     // GSensor
	Mobile    Mobile    `json:"mobile"` // 移动网络
	Disks     []Disk    `json:"disks"`  // 磁盘
	People    People    `json:"peop"`   // 人数统计
	Obds      []Obd     `json:"obds"`   // obd
	Vol       []float32 `json:"vol"`
}
