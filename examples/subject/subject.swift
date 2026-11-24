//enum list
enum Subject: String {

	//case
	case sensorAlarmSmoke = "sensor.alarm.smoke", sensorAlarmWater = "sensor.alarm.water", sensorAlarmAny = "sensor.alarm.*", sensorAnySmoke = "sensor.*.smoke", sensorAll = "sensor.>", all = ">", ordersCreated = "orders.created", ordersPaid = "orders.paid", ordersAny = "orders.*", myappUsEastUsersProfileUpdate = "myapp.us-east.users.profile.update", myappAnyUsersAll = "myapp.*.users.>", gameRoom42Action = "game.room.42.action", obsGnssBEMT = "obs.gnss.BEMT", wasmbusCtlV1LatticeComponentScaleHostId = "wasmbus.ctl.v1.{lattice}.component.scale.{host-id}", sysReqServerPing = "$SYS.REQ.SERVER.PING", jsApiStreamInfoOrders = "$JS.API.STREAM.INFO.ORDERS"

	//id
	var id: String {
		get {
			switch self {
				case .sensorAlarmSmoke:                         return "sensor.alarm.smoke"
				case .sensorAlarmWater:                         return "sensor.alarm.water"
				case .sensorAlarmAny:                           return "sensor.alarm.*"
				case .sensorAnySmoke:                           return "sensor.*.smoke"
				case .sensorAll:                                return "sensor.>"
				case .all:                                      return ">"
				case .ordersCreated:                            return "orders.created"
				case .ordersPaid:                               return "orders.paid"
				case .ordersAny:                                return "orders.*"
				case .myappUsEastUsersProfileUpdate:            return "myapp.us-east.users.profile.update"
				case .myappAnyUsersAll:                         return "myapp.*.users.>"
				case .gameRoom42Action:                         return "game.room.42.action"
				case .obsGnssBEMT:                              return "obs.gnss.BEMT"
				case .wasmbusCtlV1LatticeComponentScaleHostId:  return "wasmbus.ctl.v1.{lattice}.component.scale.{host-id}"
				case .sysReqServerPing:                         return "$SYS.REQ.SERVER.PING"
				case .jsApiStreamInfoOrders:                    return "$JS.API.STREAM.INFO.ORDERS"
			}
		}
	}
	
	//title
	var title: String {
		get {
			switch self {
				case .sensorAlarmSmoke:                         return "SensorAlarmSmoke"
				case .sensorAlarmWater:                         return "SensorAlarmWater"
				case .sensorAlarmAny:                           return "SensorAlarmAny"
				case .sensorAnySmoke:                           return "SensorAnySmoke"
				case .sensorAll:                                return "SensorAll"
				case .all:                                      return "All"
				case .ordersCreated:                            return "OrdersCreated"
				case .ordersPaid:                               return "OrdersPaid"
				case .ordersAny:                                return "OrdersAny"
				case .myappUsEastUsersProfileUpdate:            return "MyappUsEastUsersProfileUpdate"
				case .myappAnyUsersAll:                         return "MyappAnyUsersAll"
				case .gameRoom42Action:                         return "GameRoom42Action"
				case .obsGnssBEMT:                              return "ObsGnssBEMT"
				case .wasmbusCtlV1LatticeComponentScaleHostId:  return "WasmbusCtlV1LatticeComponentScaleHostId"
				case .sysReqServerPing:                         return "SysReqServerPing"
				case .jsApiStreamInfoOrders:                    return "JsApiStreamInfoOrders"
			}
		}
	}
	
	var orders: Bool {
		get {
			switch self {
				case .ordersCreated: return true
				case .ordersPaid: return true
				default: return false
			}
		}
	}
	
}
