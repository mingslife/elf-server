package model

type UploadFileReq struct {
	FileName string
	FileExt  string
	FileSize int64
	FileData []byte
}

type UploadFileRsp struct {
	Path string `json:"path"`
}

type UploadImageReq struct {
	FileName string
	FileExt  string
	FileSize int64
	FileData []byte
	Compress bool
	Width    int
	Height   int
}

type UploadImageRsp struct {
	Path string `json:"path"`
}
