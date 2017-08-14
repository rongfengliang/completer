package model

import (
	"mime/multipart"
	"fmt"
	"bytes"
	"strings"
	"strconv"
)

const (
	headerDatumType    = "Fnproject-Datumtype"
	headerResultStatus = "Fnproject-Resultstatus"
	headerResultCode   = "Fnproject-Resultcode"
	headerStageRef     = "Fnproject-Stageid"
	headerMethod       = "Fnproject-Method"
	headerHeaderPrefix = "Fnproject-Header-"
	headerErrorType    = "Fnproject-Errortype"
	headerContentType  = "Content-Type"

	datumTypeBlob     = "blob"
	datumTypeEmpty    = "empty"
	datumTypeError    = "error"
	datumTypeStageRef = "stageref"
	datumTypeHttpReq  = "httpreq"
	datumTypeHttpResp = "httpresp"
)

// DatumFromPart reads a model Datum Object from a multipart part
func DatumFromPart(part *multipart.Part) (*Datum, error) {

	datumType := part.Header.Get(headerDatumType)
	if datumType == "" {
		return nil, fmt.Errorf("Multipart part " + part.FileName() + " cannot be read as a Datum, the " + headerDatumType + " is not present ")
	}

	switch datumType {
	case datumTypeBlob:

		blob, err := readBlob(part)
		if err != nil {
			return nil, err
		}
		return &Datum{
			Val: &Datum_Blob{Blob: blob},
		}, nil

	case datumTypeEmpty:
		return &Datum{Val: &Datum_Empty{&EmptyDatum{}}}, nil
	case datumTypeError:
		errorContentType := part.Header.Get(headerContentType)
		if errorContentType != "text/plain" {
			return nil, fmt.Errorf("Invalid error datum content type on part %s, must be text/plain", part.FileName())
		}

		errorTypeString := part.Header.Get(headerErrorType)
		if "" == errorTypeString {
			return nil, fmt.Errorf("Invalid Error Datum in part %s : no %s header defined", part.FileName(), headerErrorType)
		}

		pbErrorTypeString := strings.Replace(errorTypeString, "-", "_", -1)

		// Unrecognised errors are coerced to unknown
		var pbErrorType ErrorDatumType
		if val, got := ErrorDatumType_value[pbErrorTypeString]; got {
			pbErrorType = ErrorDatumType(val)
		} else {
			pbErrorType = ErrorDatumType_unknown_error
		}

		buf := new(bytes.Buffer)
		_, err := buf.ReadFrom(part)
		if err != nil {
			return nil, fmt.Errorf("Failed to read multipart body for %s ", part.FileName())
		}

		return &Datum{
			Val: &Datum_Error{
				&ErrorDatum{Type: pbErrorType, Message: buf.String()},
			},
		}, nil

	case datumTypeStageRef:
		stageIdString := part.Header.Get(headerStageRef)
		stageId, err := strconv.ParseUint(stageIdString, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("Invalid StageRef Datum in part %s : %s", part.FileName(), err.Error())
		}
		return &Datum{Val: &Datum_StageRef{&StageRefDatum{ StageRef: uint32(stageId) }}}, nil
	case datumTypeHttpReq:
		methodString := part.Header.Get(headerMethod)
		if "" == methodString {
			return nil, fmt.Errorf("Invalid HttpReq Datum in part %s : no %s header defined", part.FileName(), headerMethod)
		}
		method, methodRecognized := HttpMethod_value[strings.ToLower(methodString)]
		if !methodRecognized {
			return nil, fmt.Errorf("Invalid HttpReq Datum in part %s : http method %s is invalid", part.FileName(), methodString)
		}
		var headers []*HttpHeader
		for hk,hvs := range part.Header {
			if strings.HasPrefix(strings.ToLower(hk), strings.ToLower(headerHeaderPrefix)) {
				for _,hv := range hvs {
					headers = append(headers, &HttpHeader{ Key: hk[len(headerHeaderPrefix):], Value: hv})
				}
			}
		}
		blob, err := readBlob(part)
		if err != nil {
			return nil, err
		}
		return &Datum{Val: &Datum_HttpReq{ &HttpReqDatum { blob, headers, HttpMethod(method)}}}, nil
	case datumTypeHttpResp:
		resultCodeString := part.Header.Get(headerResultCode)
		if "" == resultCodeString {
			return nil, fmt.Errorf("Invalid HttpResp Datum in part %s : no %s header defined", part.FileName(), headerResultCode)
		}
		resultCode, err := strconv.ParseUint(resultCodeString, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("Invalid HttpResp Datum in part %s : %s", part.FileName(), err.Error())
		}
		var headers []*HttpHeader
		for hk,hvs := range part.Header {
			if strings.HasPrefix(strings.ToLower(hk), strings.ToLower(headerHeaderPrefix)) {
				for _,hv := range hvs {
					headers = append(headers, &HttpHeader{ Key: hk[len(headerHeaderPrefix):], Value: hv})
				}
			}
		}
		blob, err := readBlob(part)
		if err != nil {
			return nil, err
		}
		return &Datum{Val: &Datum_HttpResp{ &HttpRespDatum { blob, headers, uint32(resultCode)}}}, nil
	default:
		return nil, fmt.Errorf("Unrecognised datum type")
	}
	return nil, fmt.Errorf("Unimplemented")
}

func readBlob(part *multipart.Part) (*BlobDatum, error) {
	contentType := part.Header.Get(headerContentType)
	if "" == contentType {
		return nil, fmt.Errorf("Mulitpart part %s is missing %s header", part.FileName(), headerContentType)
	}
	buf := new(bytes.Buffer)
	buf.Reset()
	_, err := buf.ReadFrom(part)
	if err != nil {
		return nil, fmt.Errorf("Failed to read multipart body from part %s", part.FileName())
	}

	return &BlobDatum{
		ContentType: contentType,
		DataString:  buf.Bytes(),
	}, nil
}
