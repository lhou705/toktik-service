namespace go cos
service Cos {
UploadResp Upload(1:UploadReq req)
}
struct UploadResp {
1:bool isSuccess
}

struct UploadReq {
1:binary file
2:string key
}