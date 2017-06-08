/*/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/

package main

import (
"errors"
"fmt"
"strconv"
"encoding/json"

"github.com/hyperledger/fabric/core/chaincode/shim"
)

var EVENT_COUNTER = "event_counter"

// ManagePatient example simple Chaincode implementation
type ManagePatient struct {
}

var PatientIndexStr = "_Patientindex"       //name for the key/value that will store a list of all known Patients

type Patient struct{             // Attributes of a Patient      
  PatientID string `json:"PatientID"`
  Address   string `json:"Address"`         
  Problems string `json:"Problems"`
  PatientName string `json:"PatientName"`
  Gender string `json:"Gender"`
  PatientMobile string `json:"PatientMobile"`
  }
// ============================================================================================================================
// Main - start the chaincode for Create Patient
// ============================================================================================================================
func main() {     
  err := shim.Start(new(ManagePatient))
  if err != nil {
    fmt.Printf("Error starting ManagePatient chaincode: %s", err)
  }
}
// ============================================================================================================================
// Init - reset all the things
// ============================================================================================================================
func (t *ManagePatient) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
  var msg string
  var err error
  if len(args) != 1 {
    return nil, errors.New("Incorrect number of arguments. Expecting 1")
  }
  // Initialize the chaincode
  msg = args[0]
  fmt.Println("ManagePatient chaincode is deployed successfully.");
  
  // Write the state to the ledger
  err = stub.PutState("abc", []byte(msg))       //making a test var "abc", I find it handy to read/write to it right away to test the network
  if err != nil {
    return nil, err
  }
  
  var empty []string
  jsonAsBytes, _ := json.Marshal(empty)               //marshal an emtpy array of strings to clear the index
  err = stub.PutState(PatientIndexStr, jsonAsBytes)
  if err != nil {
    return nil, err
  }
  err = stub.PutState(EVENT_COUNTER, []byte("1"))
  if err != nil {
    return nil, err
  }
  return nil, nil
}
// ============================================================================================================================
// Run - Our entry point for Invocations - [LEGACY] obc-peer 4/25/2016
// ============================================================================================================================
  func (t *ManagePatient) Run(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
    fmt.Println("run is running " + function)
    return t.Invoke(stub, function, args)
  }
// ============================================================================================================================
// Invoke - Our entry point for Invocations
// ============================================================================================================================
  func (t *ManagePatient) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
    fmt.Println("invoke is running " + function)

  // Handle different functions
  if function == "init" {                         //initialize the chaincode state, used as reset
    return t.Init(stub, "init", args)
  } else if function == "delete" {                     //delete a new Patient
    return t.delete(stub, args)
  }else if function== "create_patient"{
    return t.create_patient(stub,args)}           
   fmt.Println("invoke did not find func: " + function)          //error
  
  return nil, errors.New("Received unknown function invocation")
}
// ============================================================================================================================
// Query - Our entry point for Queries
// ============================================================================================================================
func (t *ManagePatient) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
  fmt.Println("query is running " + function)

  // Handle different functions
  if function == "getPatient_byID" {                         //Read a Patient by transId
    return t.getPatient_byID(stub, args)
  } 
  fmt.Println("query did not find func: " + function)           //error
  return nil, errors.New("Received unknown function query")
}
// getPatient_byID - get Patient details for a specific ID from chaincode state
//============================================================================================================================
func (t *ManagePatient) getPatient_byID(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
  var PatientID, jsonResp string
  var err error
  fmt.Println("start getPatient_byID")
  if len(args) != 1 {
    return nil, errors.New("Incorrect number of arguments. Expecting ID of the patient to query")
  }
  // set PatientID
  PatientID = args[0]
  valAsbytes, err := stub.GetState(PatientID)                  //get the PatientID from chaincode state
  if err != nil {
    jsonResp = "{\"Error\":\"Failed to get state for " + PatientID + "\"}"
    return nil, errors.New(jsonResp)
  }
  //fmt.Print("valAsbytes : ")
  //fmt.Println(valAsbytes)
  fmt.Println("end getPatient_byID")
  return valAsbytes, nil                          //send it onward
}

//create patient
//========================================================================================================================
func (t *ManagePatient) create_patient(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
  var err error
  if len(args) != 6{
    return nil, errors.New("Incorrect number of arguments. Expecting 6")
  }
  fmt.Println("start create_Patient")

  PatientID := args[0]
  Address := args[1]
  Problems := args[2]
  PatientName:= args[3]
  Gender := args[4]
  PatientMobile := args[5]
  
  
  PatientAsBytes, err := stub.GetState(PatientID)
  if err != nil {
    return nil, errors.New("Failed to get Patient Patient_id")
  }
  
  res := Patient{}
  json.Unmarshal(PatientAsBytes, &res)
 
  if res.PatientID == PatientID{
     return nil, errors.New("This Patient arleady exists")        //all stop a patient by this name exists
  }
  
  //build the CreatePatient json string manually
  PatientDetails :=  `{`+
    `"PatientID": "` + PatientID + `" , `+
    `"Address": "` + Address + `" , `+
    `"Problems": "` + Problems + `" , `+
    `"PatientName": "` + PatientName + `" , `+
    `"Gender": "` + Gender + `" , `+ 
    `"PatientMobile": "` + PatientMobile + `" , `+ 
    `}`

    fmt.Print("Patient details in bytes array: ")
    fmt.Println([]byte(PatientDetails))
  err = stub.PutState(PatientID, []byte(PatientDetails))                  //store Patient with PatientID as key
  if err != nil {
    return nil, err
  }
  //get the patient
  PatientIndexAsBytes, err := stub.GetState(PatientIndexStr)
  if err != nil {
    return nil, errors.New("Failed to get Patient index")
  }
  var PatientIndex []string
  //fmt.Print("PatientIndexAsBytes: ")
  //fmt.Println(PatientIndexAsBytes)
  
  json.Unmarshal(PatientIndexAsBytes, &PatientIndex)              //un stringify it aka JSON.parse()
    
  PatientIndex = append(PatientIndex, PatientID)                 //add Patient transID to index list
  //fmt.Println("! Patient index after appending transId: ", poIndex)
  jsonAsBytes, _ := json.Marshal(PatientIndex)
  //fmt.Print("jsonAsBytes: ")
  //fmt.Println(jsonAsBytes)
  err = stub.PutState(PatientIndexStr, jsonAsBytes)            //store name of Patient
  if err != nil {
    return nil, err
  }

  fmt.Println("end create_Patient")
  return nil, nil
}
//=====================Delete Vessel==================================================================
func (t *ManagePatient) delete(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
  fmt.Println("enter delete function")
  if len(args) != 1 {
    return nil, errors.New("Incorrect number of arguments. Expecting 1")
  }
  // set PatientID
  PatientID := args[0]

  //get the patient index
 PatientIndexAsBytes, err := stub.GetState(PatientIndexStr)
  if err != nil {
    return nil, errors.New("Failed to get Patient index")
  }
  //fmt.Println("poAsBytes in delete po")
  //fmt.Println(poAsBytes);
  var PatientIndex []string
  json.Unmarshal(PatientIndexAsBytes, &PatientIndex)               //un stringify it aka JSON.parse()
  //fmt.Println("poIndex in delete po")
  //fmt.Println(poIndex);
  //remove marble from index
  for i,val := range PatientIndex{
    fmt.Println(strconv.Itoa(i) + " - looking at " + val + " for " + PatientID)
    if val == PatientID{                             //find the correct Vessel
      fmt.Println("found patient with matching patientID")
      PatientIndex = append(PatientIndex[:i], PatientIndex[i+1:]...)     //remove it
      for x:= range PatientIndex{                      //debug prints...
        fmt.Println(string(x) + " - " + PatientIndex[x])
      }
      break
    }
  }
  jsonAsBytes, _ := json.Marshal(PatientIndex)                 //save new index
  err = stub.PutState(PatientIndexStr, jsonAsBytes)
  return nil, nil
}
