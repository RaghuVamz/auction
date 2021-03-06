/******************************************************************
Copyright IT People Corp. 2017 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

                 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

******************************************************************/

///////////////////////////////////////////////////////////////////////
// Author : IT People - Mohan - table API for v1.0
// Enable CouchDb as the database..
// Purpose: Explore the Hyperledger/fabric and understand
// how to write an chain code, application/chain code boundaries
// The code is not the best as it has just hammered out in a day or two
// Feedback and updates are appreciated
///////////////////////////////////////////////////////////////////////
package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/hyperledger/fabric/auction/itpUtils"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// FileChaincode .
type FileChaincode struct {
}

// FileObject file object
type FileObject struct {
	FileID        string //provided in args
	FilePath      string // provided in args
	FileImage     []byte // This has to be generated AES encrypted using the file name
	AESKey        []byte // This is generated by the AES Algorithms
	FileImageType string // should be used to regenerate the appropriate image type
	UserID        string // This is validated for a user registered record
}

func main() {
	err := shim.Start(new(FileObjectChaincode))
	if err != nil {
		fmt.Printf("Error starting upload chaincode: %s", err)
	}
}

// Init deploys a chaincode
func (t *FileObjectChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) pb.response {
	if len(args) != 0 {
		return nil, errors.New("Incorrect number of arguments. Expecting 0")
	}
	//Invoke doesnt do anything
}

//Invoke invokes a chaincode
//./peer chaincode invoke -l golang -n mycc -c '{"Function": "invoke", "Args":["1111", "/filepath/sample.pdf", "user_id"]}'
func (t *FileObjectChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) pb.response {
	fileObject, err := CreateFileObject(args[0:])
	if err != nil {
		fmt.Println("Cannot create item object \n")
		return shim.Error("Cannot create file object")
	}
	
	// Convert Item Object to JSON
	buff, err := FOtoJSON(fileObject) //
	if err != nil {
		error_str := "Invoke() : Failed Cannot create object buffer for write : " + args[1]
		fmt.Println(error_str)
		return shim.Error(error_str)
	}else {
		// Update the ledger with the Buffer Data
		err = stub.PutState(args[0], buff)
		if err != nil {
			fmt.Println("Invoke() : write error while inserting record\n")
			return shim.Error("Invoke() : write error while inserting record : " + err.Error())
		}	
}
}

func CreateFileObject(args []string) (FileObject, error) {

	var err error
	var myFile FileObject

	// Check there are 4 Arguments provided as per the the struct - two are computed
	if len(args) != 4 {
		fmt.Println("CreateFileObject(): Incorrect number of arguments. Expecting 4 ")
		return myFile, errors.New("CreateFileObject(): Incorrect number of arguments. Expecting 13 ")
	}

	// Validate FileID is an integer

	_, err = strconv.Atoi(args[0])
	if err != nil {
		fmt.Println("CreateFileObject(): File ID should be an integer create failed! ")
		return myFile, errors.New("CreateFileObject(): ART ID should be an integer create failed! ")
	}

	// Validate File exists based on the name provided
	// Looks for file in current directory of application and must be fixed for other locations

	imagePath := args[2]
	if _, err := os.Stat(imagePath); err == nil {
		fmt.Println(imagePath, "  exists!")
	} else {
		fmt.Println("CreateFileObject(): Cannot find or load Picture File = %s :  %s\n", imagePath, err)
		return myFile, errors.New("CreateFileObject(): File not found " + imagePath)
	}

	// Get the File Image and convert it to a byte array
	imagebytes, fileType := itpUtils.ImageToByteArray(imagePath)

	// Generate a new key and encrypt the image

	AESKey, _ := itpUtils.GenAESKey()
	AESEnc := itpUtils.Encrypt(AES_key, imagebytes)

	// Append the AES Key, The Encrypted Image Byte Array and the file type
	myFile = ItemObject{args[0], args[1], args[2], AES_enc, AES_key, fileType, args[3]]

	fmt.Println("CreateFileObject(): Item Object created: ", myFile.ItemID, myFile.AES_Key)
	return myFile, nil
}


// Converts an File Object to a JSON String

func FOtoJSON(fo fileObject) ([]byte, error) {

	fjson, err := json.Marshal(fo)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return fo, nil
}
