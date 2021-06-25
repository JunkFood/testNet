package main

//외부모듈 포함
import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

//클래스정의
type SmartContract struct {
}

//구조체정의
// type resultformat struct{
// 	code bool // s or f
// 	retrunMSG []byte
// 	errorMSG string
// }

type Developer struct {
	ID                  string  `json:"id"`       //developer의 id
	Avg                 float32 `json:"avg"`      //developer의 proectScore의 average
	NumProject          int     `json:"nproject"` //developer가 완료한 프로젝트의 수
	CurrentState        int     `json:"state"`    //developer의 프로젝트 참가 여부
	CurrentProject      string  `json:"pname"`    //developer의 현재 프로젝트 이름
	CurrentProjectScore int     `json:"pscore"`   //developer의 현재 프로젝트에서의 점수
}

const (
	REGISTERED int = iota
	JOINED
	FINISHED
)

//TODO 1.프로젝트 구조체 정의
type Project struct {
	//이름,
}

//init함수
func (s *SmartContract) Init(stub shim.ChaincodeStubInterface) peer.Response {
	//log
	fmt.Println("func init started")
	//cc instantiate, upgrade 기능을 초기화기능
	return shim.Success(nil)
}

//invoke함수
func (s *SmartContract) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	//
	fn, args := stub.GetFunctionAndParameters()

	if fn == "registerUser" {
		return s.registerUser(stub, args)
	} else if fn == "joinProject" {
		return s.joinProject(stub, args)
	} else if fn == "recordScore" {
		return s.recordScore(stub, args)
	} else if fn == "readDev" {
		return s.readDev(stub, args)
	} else {
		return shim.Error("Not supportes smartcontract function name")
	}
}

//registerUser
func (s *SmartContract) registerUser(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	// args : id
	if len(args) != 1 {
		return shim.Error("registerUser function need 1 parameter")
	}

	//getState 해봐서 err체크(id가 있으면 err)
	var Dev = Developer{ID: args[0], Avg: 0., NumProject: 0, CurrentState: 0}
	devAsBytes, _ := json.Marshal(Dev)
	stub.PutState(args[0], devAsBytes)

	return shim.Success([]byte("tx is submitted"))
}

//joinProject
func (s *SmartContract) joinProject(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	// args : id, nproject
	if len(args) != 2 {
		return shim.Error("joinProject function need 2 parameter")
	}

	devAsBytes, err := stub.GetState(args[0])

	//getstate가 수행중 오류발생시
	if err != nil {
		return shim.Error("GetState function occurred a error")
	}
	//id 가 없는 경우
	if devAsBytes == nil {
		return shim.Error("ID is not registered")
	}

	var Dev = Developer{}
	_ = json.Unmarshal(devAsBytes, &Dev)

	Dev.CurrentProject = args[1]
	Dev.CurrentState = 1 // join

	devAsBytes, _ = json.Marshal(Dev)
	stub.PutState(args[0], devAsBytes)

	return shim.Success([]byte("tx is submitted"))
}

//recordScore
func (s *SmartContract) recordScore(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	// args : id, projectname, score
	if len(args) != 3 {
		return shim.Error("recordScore function need 3 parameter")
	}

	//getState 해봐서 err체크(id가 있으면 err)
	devAsBytes, err := stub.GetState(args[0])

	//getstate가 수행중 오류발생시
	if err != nil {
		return shim.Error("GetState function occurred a error")
	}
	//id 가 없는 경우
	if devAsBytes == nil {
		return shim.Error("ID is not registered")
	}

	var Dev = Developer{}
	_ = json.Unmarshal(devAsBytes, &Dev)

	Dev.CurrentProject = args[1]
	Dev.CurrentProjectScore, _ = strconv.Atoi(args[2])
	Dev.NumProject++

	var newAvg float32
	newAvg = (Dev.Avg*float32(Dev.NumProject-1) + float32(Dev.CurrentProjectScore)) / float32(Dev.NumProject)
	Dev.Avg = newAvg

	Dev.CurrentState = 2 // FINISHED

	devAsBytes, _ = json.Marshal(Dev)
	stub.PutState(args[0], devAsBytes)

	return shim.Success([]byte("tx is submitted"))
}

func (s *SmartContract) readDev(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	// args : id
	if len(args) != 1 {
		return shim.Error("registerUser function need 1 parameter")
	}

	//getState 해봐서 err체크(id가 있으면 err)
	devAsBytes, _ := stub.GetState(args[0])

	return shim.Success(devAsBytes)
}

//registerProject
//recordProject
//finalizeProject

func main() {
	err := shim.Start(new(SmartContract))

	if err != nil {
		fmt.Println("Error creating new Smart Contract : %s", err)
	}
}
