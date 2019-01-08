<!-- 
//
// Copyright 2018 Malin Yamato --  All rights reserved.
// https://github.com/MalinYamato
//
// MIT License
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
//     * Neither the name of Rakuen. nor the names of its
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
-->
    <link href="/css/style.css" rel="stylesheet" type="text/css"/>
    <link href="/css/menus.css" rel="stylesheet" type="text/css"/>
    <link href="/css/modal.css" rel="stylesheet" type="text/css"/>
    <link rel="stylesheet" href="/css/emojionearea.min.css">
    <link rel="stylesheet" href="/css/swiper.css">

    <link rel="stylesheet" href="https://www.w3schools.com/w3css/4/w3.css">
   
    <!-- ---------------------------------------->
    <meta http-equiv="Content-Type" content="text/html; charset=UTF-8"/>

    <script type="text/javascript"
            src="https://cdnjs.cloudflare.com/ajax/libs/webrtc-adapter/6.0.3/adapter.min.js"></script>
    <script type="text/javascript" src="https://cdnjs.cloudflare.com/ajax/libs/jquery/3.3.1/jquery.min.js"></script>
    <script type="text/javascript"
            src="https://cdnjs.cloudflare.com/ajax/libs/jquery.blockUI/2.70/jquery.blockUI.min.js"></script>
    <script type="text/javascript" src="https://cdnjs.cloudflare.com/ajax/libs/toastr.js/2.1.3/toastr.min.js"></script>
    <script type="text/javascript"
            src="https://cdnjs.cloudflare.com/ajax/libs/twitter-bootstrap/4.1.3/js/bootstrap.min.js"></script>
    <script type="text/javascript"
            src="https://cdnjs.cloudflare.com/ajax/libs/bootbox.js/4.1.0/bootbox.min.js"></script>
    <script type="text/javascript" src="https://cdnjs.cloudflare.com/ajax/libs/spin.js/2.3.2/spin.min.js"></script>
    <script type="text/javascript" src="/js/janus.js"></script>
    <script type="text/javascript" src="/js/videoroom.js"></script>
    <script src="//twemoji.maxcdn.com/twemoji.min.js"></script>
    <script src="/js/SVGButton.js"></script>
    <script type="text/javascript"></script>
    <script type="text/javascript" src="/js/emojionearea.min.js"></script>
    <script src="/js/swiper.js"></script>

 function ChatURL() {
            return "not implemented";
        }

        var _id = 0;

        function getID() {
            if (_id > 1000) {
                _id = 0;
            }
            _id++;
            return _id;
        }

        function getUniqueID(id) {
            return id + "@" + getID();
        }

        function getAsync(op, url, cb) {
            var request = new XMLHttpRequest();
            request.open("POST", ChatURL() + url, true);
            request.setRequestHeader('Content-Type', 'application/json');
            request.responseType = 'json';
            request.onload = function () {
                status = request.status;
                if (status == 200) {
                    //var json_res =  JSON.parse(request.responseText);
                    cb && cb(request.response);
                } else {
                    console.log("getAsync problem " + op + " " + status)
                }
            };
            request.send(JSON.stringify(op));
        }

        function getUser(op) {
            var request = new XMLHttpRequest();
            request.open("POST", ChatURL() + "/profile", false);
            request.setRequestHeader('Content-Type', 'application/json');
            request.send(JSON.stringify(op));
            if (request.status == 200) {
                json_response = JSON.parse(request.response);
                return json_response;
            } else {
                console.log("getUser problem " + status)
            }
        }

        function logoutAction(op) {
            var request = new XMLHttpRequest();
            request.open("POST", ChatURL() + "/logout", false);
            request.onload = function () {
                status = request.status;
                if (status == 200) {
                    window.location = ChatURL() + "/";
                } else {
                    console.log("logoutAction problem " + status)
                }
            };
            request.send({token: "{{.Person.Token }}"});
        }

        function changeRoom(op, cb) {
            var request = new XMLHttpRequest();
            request.open("POST", ChatURL() + "/RoomManager", true);
            request.setRequestHeader('Content-Type', 'application/json');
            request.responseType = 'json';
            request.onload = function () {
                status = request.status;
                if (status == 200) {
                    //var json_res =  JSON.parse(request.responseText);
                    var info = request.response;
                    if (info.status.status != "SUCCESS") {
                        console.log("Calling RooomManger::ChangeRoom returned " + info.status.status + " " + info.status.detail);
                    }
                    cb && cb(request.response);
                } else {
                    console.log("changeRoom problem " + status)
                }
            };
            request.send(JSON.stringify(op));
        }

        function AddTarget(op, cb) {
            var request = new XMLHttpRequest();
            request.open("POST", ChatURL() + "/TargetManager", true);
            request.setRequestHeader('Content-Type', 'application/json');
            request.responseType = 'json';
            request.onload = function () {
                status = request.status;
                if (status == 200) {
                    //var json_res =  JSON.parse(request.responseText);
                    var info = request.response;
                    if (info.status.status != "SUCCESS") {
                        console.log("Calling TargetManager::Add returned " + info.status.status + " " + info.status.detail);
                    }
                    cb && cb(request.response);
                } else {
                    console.log("Target::Add problem " + status)
                }
            };
            request.send(JSON.stringify(op));
        }

        function RemoveTarget(targetId, cb) {
            var op = {"op": "RemoveTarget", "ids": [targetId]};
            var request = new XMLHttpRequest();
            request.open("POST", ChatURL() + "/TargetManager", true);
            request.setRequestHeader('Content-Type', 'application/json');
            request.responseType = 'json';
            request.onload = function () {
                status = request.status;
                if (status == 200) {
                    var info = request.response;
                    if (info.status.status != "SUCCESS") {
                        console.log("Calling TargetManager::RemoveTarget returned " + info.status.status + " " + info.status.detail);
                    } else {
                        console.log("Delete target requested /targetManager: " + info.status.detail);
                    }
                    cb && cb(request.response);
                } else {
                    console.log("Target::RemoveTarget problem " + status)
                }
            };
            request.send(JSON.stringify(op))
        }

        function genderShort(gender) {
            if (gender == "Female") {
                return "F"
            } else if (gender == "Male") {
                return "M";
            } else {
                return "T";
            }
        }
