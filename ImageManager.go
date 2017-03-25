//
// Copyright 2017 Malin Lääkkö -- Yamato Digital Audio.  All rights reserved.
// https://github.com/MalinYamato
//
// Yamato Digital Audio https://yamato.xyz
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
//
//     * Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.
//     * Redistributions in binary form must reproduce the above
// copyright notice, this list of conditions and the following disclaimer
// in the documentation and/or other materials provided with the
// distribution.
//     * Neither the name of Yamato Digital Audio. nor the names of its
// contributors may be used to endorse or promote products derived from
// this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.


package main

import (
	"os"
	"strings"
	"io"
	"log"
	"encoding/json"
	"net/http"
	"strconv"
	"image"
	"image/gif"
	"image/png"
	"image/jpeg"
	"github.com/robfig/graphics-go/graphics"
	"io/ioutil"

	"github.com/satori/go.uuid"
)



type ImagesResponse struct {
	Status          Status                    `json:"status"`
	ProfileImageURL string                    `json:"profileImageURL"`
	Images          []map[string]string       `json:"images"`
}

type ImageRequest struct {
	Op       string    `json:"op"`
	ImageURL string    `json:"imageURL"`
}

func cropImage(fileroot string, inputFile string, extension string) Status {
	fSrc, err := os.Open(inputFile)
	if err != nil {
		log.Fatal(err)
	}
	var src image.Image
	defer fSrc.Close()
	if (strings.ToLower(extension) == "jpg" || strings.ToLower(extension) == "jpeg") {
		src, _, err = image.Decode(fSrc)
	} else if (strings.ToLower(extension) == "png") {
		src, err = png.Decode(fSrc)
	} else if (strings.ToLower(extension) == "gif") {
		src, err = gif.Decode(fSrc)
	} else {
		return Status{Status: ERROR, Detail: extension + " not suppored! Supported formats are jpg, png and gif"}
	}
	if err != nil {
		log.Fatal(err)
	}
	dst := image.NewRGBA(image.Rect(0, 0, 150, 160))
	graphics.Thumbnail(dst, src)

	toimg, err := os.Create(fileroot + "/small." + "jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer toimg.Close()
	jpeg.Encode(toimg, dst, &jpeg.Options{100})
	return Status{SUCCESS, ""}
}


func ImageManager_DeleteHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var status Status
	var request ImageRequest
	if r.Method == "POST" {
		var person Person
		var ok bool
		token, _, err := getCookieAndTokenfromRequest(r, true)
		if err != nil {
			status = Status{ERROR, err.Error()}
		} else {
			person, ok = _persons.findPersonByToken(token)
			if ! ok {
				status = Status{ERROR, err.Error()}
			} else {
				decoder := json.NewDecoder(r.Body)
				err = decoder.Decode(&request)
				if err != nil {
					log.Println("Json decoder error> ", err.Error())
					panic(err)
				}
				if request.Op == "Delete" {
					if (request.ImageURL == person.PictureURL) {
						status = Status{Status: WARNING, Detail: "Users are not allowed to delete profile picture!"}
					} else {
						path := person.path() + "/img"
						files, _ := ioutil.ReadDir(path)
						for _, file := range files {
							if strings.Contains(request.ImageURL, file.Name()) {
								os.RemoveAll(path + "/" + file.Name())
							}
						}
						status = Status{Status: SUCCESS, Detail: "Picture deleted!"}
					}
				}
			}
		}
	} else {
		status = Status{Status: ERROR, Detail:"Bad HTTPS method"}
		log.Println("ImageManager: Unknown HTTP method ", r.Method)
	}
	json_response, err := json.Marshal(status)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(json_response)
}

func ImageManger_GetHandler(w http.ResponseWriter, r *http.Request) {
	var images ImagesResponse
	defer r.Body.Close()
	var status Status
	if r.Method == "POST" {
		person, _ := _persons.findPersonByCookie(r)
		if (person.Keep == false ) {
			status = Status{Status: WARNING, Detail: "Only members have images!" }
		} else {
			path := person.path() + "/img/"
			files, err := ioutil.ReadDir(path)
			if err != nil {
				log.Fatal(err)
			}
			for _, file := range files {
				elem := make(map[string]string)
				files2, err := ioutil.ReadDir(path + file.Name())
				if err != nil {
					log.Fatal(err)
				}
				for _, file2 := range files2 {
					url := "https://" + endpoint.host + strings.Trim(path, ".") + file.Name() + "/" + file2.Name()
					log.Println("Path " + url)
					if strings.Contains(file2.Name(), "small") {
						elem["small"] = url
					} else {
						elem["normal"] = url
					}
				}
				images.Images = append(images.Images, elem)
			}
			status = Status{Status: SUCCESS}
		}
	} else {
		status = Status{Status: ERROR}
		log.Println("ImageManager: Unknown HTTP method ", r.Method)
	}
	images.Status = status
	json_response, err := json.Marshal(images)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(json_response)
}

type ImageFile struct {
	Name    string  `json:"name"`
	Type    string  `json:"type"`
	Height  int16   `json:"height"`
	Width   int16   `json:"width"`
}

type ImageInfo struct {
	OriginalFileName string   `json:"originalFileName"`
	Description string        `json:"description"`
	Variants []ImageFile      `json:"variants"`
}

func ImageManager_UploadHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("got ...image!")
	defer r.Body.Close()
	var status Status
	var lenght = 0

	if r.Method == "POST" {
		person, _ := _persons.findPersonByCookie(r)
		if (person.Keep == false ) {
			status = Status{Status: WARNING, Detail: "Only members may save images!" }
			log.Println("attemtp to upload pictures wihout being a member!")
		} else {
			err := r.ParseMultipartForm(100000)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			m := r.MultipartForm
			files := m.File["images[]"]
			lenght = len(files)
			for i, _ := range files {
				file, err := files[i].Open()
				defer file.Close()
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				res := strings.Split(files[i].Filename, ".")
				extension := res[1]

				log.Println("file " + files[i].Filename)
				fileroot := person.path() + "/img/" + uuid.NewV4().String()
				err = os.Mkdir(fileroot, 0777)
				if err != nil {
					panic(err)
				}
				filepath := fileroot + "/normal." + extension
				dst, err := os.Create(fileroot + "/normal." + extension)
				defer dst.Close()
				if err != nil {
					panic(err)
				}
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				if _, err := io.Copy(dst, file); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				status = cropImage(fileroot, filepath, extension)
			}
		}
		if status.Status == SUCCESS {
			status = Status{Status: SUCCESS, Detail: strconv.Itoa(lenght) + " file(s) uploaded!" }
		}
	}else {
		status = Status{Status: ERROR}
		log.Println("ImageManager: Unknown HTTP method ", r.Method)
	}
	json_response, err := json.Marshal(status)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(json_response)
}

