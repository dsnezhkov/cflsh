package common

const (

	WorkerInterval = 3 //(seconds)
	CmdBucketKeyPrefix = "command"
	OutBucketKeyPrefix = "output"
)


type NameSpace struct {
	Name string
	UUID string
}

// Command key write
type CmdKVWRequest struct {
	NameSpace  NameSpace
	CmdKName   string
	CmdPayload []byte
}
// Command key read
type CmdKVRRequest struct {
	NameSpace NameSpace
	CmdKName  string
}
// Output key write
type OutKVWRequest struct {
	NameSpace  NameSpace
	OutKName   string
	OutPayload []byte
}
// Output key read
type OutKVRRequest struct {
	NameSpace NameSpace
	OutKName  string
}
