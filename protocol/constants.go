package protocol

const (
	HeaderDatumType    = "Fnproject-Datumtype"
	HeaderResultStatus = "Fnproject-Resultstatus"
	HeaderResultCode   = "Fnproject-Resultcode"
	HeaderStageRef     = "Fnproject-Stageid"
	HeaderCallerRef    = "Fnproject-Callerid"
	HeaderHookRef      = "Fnproject-Hookid"
	HeaderMethod       = "Fnproject-Method"
	HeaderHeaderPrefix = "Fnproject-Header-"
	HeaderErrorType    = "Fnproject-Errortype"
	HeaderStateType    = "Fnproject-Statetype"
	HeaderCodeLocation = "Fnproject-Codeloc"
	HeaderFlowId       = "Fnproject-FlowId"

	HeaderContentType = "Content-Type"

	ResultStatusSuccess = "success"
	ResultStatusFailure = "failure"

	DatumTypeBlob     = "blob"
	DatumTypeEmpty    = "empty"
	DatumTypeError    = "error"
	DatumTypeStageRef = "stageref"
	DatumTypeHttpReq  = "httpreq"
	DatumTypeHttpResp = "httpresp"
	DatumTypeState    = "state"
)
