package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"net"
)

func main() {
	key := []byte("sending file is great ;)")

	conn, err := net.Dial("tcp", "127.0.0.1:9080")

	if err != nil {
		panic(err)
	}

	defer func() {
		fmt.Println("Bye !")

		conn.Close()
	}()

	///////////////////////////////////////////////////////////////////
	//this part is specified for the cipher key creation and encryption
	//of the connection////////////////////////////////////////////////
	///////////////////////////////////////////////////////////////////

	block, cipherErr := aes.NewCipher(key)

	if cipherErr != nil {
		fmt.Errorf("Can't create cipher:", cipherErr)

		return
	}

	iv := make([]byte, aes.BlockSize)

	if _, randReadErr := io.ReadFull(rand.Reader, iv); randReadErr != nil {
		fmt.Errorf("Can't build random iv", randReadErr)

		return
	}

	_, ivWriteErr := conn.Write(iv)

	if ivWriteErr != nil {
		fmt.Errorf("Can't send IV:", ivWriteErr)

		return
	} else {
		fmt.Println("IV Sent:", iv)
	}

	stream := cipher.NewCFBEncrypter(block, iv)

	////////////////////////////////////////////////////////////////////
	//this part is specified for downloading the file from the server//
	///////////////////////////////////////////////////////////////////

	var sizeMB = 5 << (10 * 2) // 5 megabytes , you can increase it here !
	buf := make([]byte, sizeMB)
	for {
		rLen, rErr := conn.Read(buf)//length of the buffer array

		if rErr == nil {
			//from here start reading data

			//file name to be downloaded
			fileName := "picture1.jpg"

			//this variable for the type of premission that your file will use
			mode := int(0644)

			stream.XORKeyStream(buf[:rLen], buf[:rLen])

			//this is the received bytes of file from server
			byteArray := buf[:rLen]

			err := ioutil.WriteFile(fileName, byteArray, os.FileMode(mode))

			if err != nil {
				panic(err)
			}

			fmt.Printf("File %s has been downloaded !! \n", fileName)

		}

		if rErr == io.EOF {
			break
		}

		fmt.Errorf(
			"Error while reading from",
			conn.RemoteAddr(),
			":",
			rErr,
		)
		break
	}
}
