package common

type ScriptLang string

const (
	ScriptLangTs         = ScriptLang("ts")
	ScriptLangCoffee     = ScriptLang("coffeescript")
	ScriptLangJavascript = ScriptLang("javascript")
)

type FlowElementsPrototypeType string

const (
	FlowElementsPrototypeDefault        = FlowElementsPrototypeType("default")
	FlowElementsPrototypeMessageHandler = FlowElementsPrototypeType("MessageHandler")
	FlowElementsPrototypeMessageEmitter = FlowElementsPrototypeType("MessageEmitter")
	FlowElementsPrototypeTask           = FlowElementsPrototypeType("Task")
	FlowElementsPrototypeGateway        = FlowElementsPrototypeType("Gateway")
	FlowElementsPrototypeFlow           = FlowElementsPrototypeType("Flow")
)

type StatusType string

const (
	Enabled  = StatusType("enabled")
	Disabled = StatusType("disabled")
)

type DeviceType string

const (
	DevTypeSmartBus = DeviceType("smartbus")
	DevTypeModBus   = DeviceType("modbus")
	DevTypeZigbee   = DeviceType("zigbee")
	DevTypeDefault  = DeviceType("default")
)

type ThreadState string

const (
	ThreadOk          = ThreadState("ok")
	ThreadBusy        = ThreadState("busy")
	ThreadDevNotFound = ThreadState("port by device not found")
	ThreadNotFound    = ThreadState("port not found")
	ThreadAllBusy     = ThreadState("all ports busy")
)
