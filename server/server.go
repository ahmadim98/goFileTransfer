package main

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"os"
	"bufio"
	"net"
)

func main() {
	ln, err := net.Listen("tcp", "127.0.0.1:9080")

	if err != nil {
		panic(err)
	}

	key := []byte("sending file is great ;)")

	fmt.Println("Started Listening")

	if err != nil {
		panic(err)
	}

	for {
		conn, err := ln.Accept()

		if err != nil {
			fmt.Errorf(
				"Error while handling request from",
				conn.RemoteAddr(),
				":",
				err,
			)
		}

		go func(conn net.Conn) {
			defer func() {
				fmt.Println(
					conn.RemoteAddr(),
					"Closed",
				)

				conn.Close()
			}()

			///////////////////////////////////////////////////////////////////
			//this part is specified for the cipher key creation and encryption
			//of the connection////////////////////////////////////////////////
			///////////////////////////////////////////////////////////////////

			block, blockErr := aes.NewCipher(key)

			if blockErr != nil {
				fmt.Println("Error creating cipher:", blockErr)

				return
			}

			iv := make([]byte, 16)

			ivReadLen, ivReadErr := conn.Read(iv)

			if ivReadErr != nil {
				fmt.Println("Can't read IV:", ivReadErr)

				return
			}

			iv = iv[:ivReadLen]

			if len(iv) < aes.BlockSize {
				fmt.Println("Invalid IV length:", len(iv))

				return
			}

			fmt.Println("Received IV:", iv)

			stream := cipher.NewCFBDecrypter(block, iv)

			fmt.Println("Hello", conn.RemoteAddr())

			/////////////////////////////////////////////////////////////
			//this part is specified for sending the file to the client//
			/////////////////////////////////////////////////////////////

			//from here start sending the data

			//file name to be sent
			file, err := os.Open("./picture1.jpg")

			 if err != nil {
							 fmt.Println(err)
							 os.Exit(1)
			 }

			 defer file.Close()

			 fileInfo, _ := file.Stat()

			 var size int64 = fileInfo.Size()

			 filebytes := make([]byte, size)

			 buffer := bufio.NewReader(file)

			 //here it will turn file into data buffer before sending it
			 _, err = buffer.Read(filebytes)

			 encrypted := make([]byte,len(filebytes))//here after the encryption !

			 decrypted := []byte(filebytes)

			 stream.XORKeyStream(encrypted, decrypted)

			 //here the file will be sent
			 _, writeErr := conn.Write(encrypted)

			 fmt.Println("File has been sent already !")

			 if writeErr != nil {
				 fmt.Println("Write failed:", writeErr)
				 return
			 }
		}(conn)
	}
}
