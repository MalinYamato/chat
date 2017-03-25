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


function allowDrop(ev) {
    ev.preventDefault();
}

function drag(ev) {

    //ev.dataTransfer.effectAllowed = "copy";
    ev.dataTransfer.setData("text", ev.target.id);
    var data = ev.dataTransfer.getData("text", ev.target.id);
    console.log(ev.target.id);
    var current = document.getElementById("centerNenu");
    current.appendChild(document.getElementById(data));
}

function drop(ev) {
    //console.trace("dropped");
    ev.preventDefault();
    ev.dataTransfer.dropEffect = "copy";
    var data = ev.dataTransfer.getData("text");
    ev.target.appendChild(document.getElementById(data));
    var current = document.getElementById("centerNenu");
    current.appendChild(document.getElementById(data));
}


function myFunction(myid) {
    console.log("clicked" + myid);
    userId = myid.split("@")[0]
    var x = document.getElementById(myid);
    console.log(x);
    if (x.className.indexOf("w3-show") == -1) {
        // json format
        var person = {"userId": userId};
        $.ajax({
            url: 'https://{{.Host}}/profile',
            type: 'post',
            data: JSON.stringify(person),
            contentType: 'application/json',
            success: function (person) {
                if (person.userId != null) {
                    var bigImage = document.createElement("img");
                    bigImage.setAttribute("src", person.pictureURL);
                    x.appendChild(bigImage);
                    var info = document.createElement("ul");
                    var info1 = document.createElement("li");
                    var info2 = document.createElement("li");
                    var info3 = document.createElement("li");
                    var info4 = document.createElement("li");
                    var info5 = document.createElement("li");
                    var info6 = document.createElement("li");
                    info1.innerText = "        " + person.nic;
                    info2.innerText = "        " + person.firstName + " " + person.lastName;
                    info3.innerText = "        " + person.gender + " " + person.sexualOrienation;
                    info4.innerText = "        " + person.country + " " + person.town;
                    info5.innerText = "        " + "";
                    info6.innerText = "        " + person.profession;
                    info.appendChild(info1);
                    info.appendChild(info2);
                    info.appendChild(info3);
                    info.appendChild(info4);
                    info.appendChild(info5);
                    info.appendChild(info6);
                    info.setAttribute("id", "profileInfo");
                    x.appendChild(info);
                    console.log(person);
                    x.className += " w3-show";
                }
            }
        });
    } else {
        while (x.lastChild) {
            x.removeChild(x.lastChild);
        }
        x.className = x.className.replace(" w3-show", "");
    }
}

//
// AJAX is better and more secure than WS for sensitive data and structured requests that expect structured response
//
function RequestToPublish(anId) {
    var targetID = anId.split("@")[0];
    var message = {"op": "AddTarget", "ids": [targetID]};
    console.log("message> " + message["ids"]);
    if ($("li[id='" + "Target:" + targetID + "']").length > 0) {
        return; // dupplicate
    }
    $.ajax({
        url: 'https://{{.Host}}/TargetManager',
        type: 'post',
        data: JSON.stringify(message),
        contentType: 'application/json',
        success: function (r) {
            console.log("Request to publish response: " + r.status.detail + " " + r.person.email);
            if (r.person.userId != null) {
                var liIten  = document.createElement("li");
                liIten.setAttribute("id", "Target:"+r.person.userId);
                var imagen = document.createElement("img");
                imagen.setAttribute("src", r.person.pictureURL);
                var elem = document.getElementById("PublishList");
                if (r.status.detail == "GREEN") {
                    imagen.setAttribute("style", "border-color:green; border-width: 4px; border-style:solid;")
                } else if (r.status.detail == "BLUE") {
                    imagen.setAttribute("style", "border-color:blue; border-width: 4px; border-style:solid;")
                }
                imagen.setAttribute("onclick", "RemoveTarget('" + targetID + "')");
                liIten.appendChild(imagen);
                elem.appendChild(liIten);
            }
        }
    });
}

function RemoveTarget(targetId)
{
    var message = {"op": "RemoveTarget", "ids": [targetId]};
    $.ajax({
        url: 'https://{{.Host}}/TargetManager',
        type: 'post',
        data: JSON.stringify(message),
        contentType: 'application/json',
        success: function (r) {
            var elem = document.getElementById("PublishList");
            var child = document.getElementById("Target:"+targetId);
            elem.removeChild(child);
        }
    });
}
