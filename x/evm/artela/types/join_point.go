package types

type JoinPointRunType int64

const (
	JoinPointRunType_VerifyTx          JoinPointRunType = 1
	JoinPointRunType_PreTxExecute      JoinPointRunType = 2
	JoinPointRunType_PreContractCall   JoinPointRunType = 4
	JoinPointRunType_PostContractCall  JoinPointRunType = 8
	JoinPointRunType_PostTxExecute     JoinPointRunType = 16
	JoinPointRunType_PostTxCommit      JoinPointRunType = 32
	JoinPointRunType_Operation         JoinPointRunType = 64
	JoinPointRunType_OnBlockInitialize JoinPointRunType = 128
	JoinPointRunType_OnBlockFinalize   JoinPointRunType = 256
)

// Enum value maps for JoinPointRunType.
var (
	JoinPointRunType_name = map[int64]string{
		1:   "verifyTx",
		2:   "preTxExecute",
		4:   "preContractCall",
		8:   "postContractCall",
		16:  "postTxExecute",
		32:  "postTxCommit",
		64:  "operation",
		128: "onBlockInitialize",
		256: "onBlockFinalize",
	}
	JoinPointRunType_value = map[string]int64{
		"verifyTx":          1,
		"preTxExecute":      2,
		"preContractCall":   4,
		"postContractCall":  8,
		"postTxExecute":     16,
		"postTxCommit":      32,
		"operation":         64,
		"onBlockInitialize": 128,
		"onBlockFinalize":   256,
	}
)

func CheckIsJoinPoint(runJPs int64) (bool, map[int64]string) {
	jpMap := make(map[int64]string)
	if runJPs <= 0 {
		return false, jpMap
	}
	for k, v := range JoinPointRunType_name {
		// verify with & to see if it is included jp value，like:  5&1==1
		if runJPs&k == k {
			jpMap[k] = v
		}
	}
	return len(jpMap) > 0, jpMap
}
