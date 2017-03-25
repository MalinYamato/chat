/**
 * Created by malin on 2017-03-09.
 */
//
// (C) 2017 Yamato Digital Audio
// Author: Malin af Lääkkö
//

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
//     * Neither the name of Google Inc. nor the names of its
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

function setDefaultActions (obj) {
    obj.s.onmouseover =  function () {
        obj.s.style.cursor = "crosshair";
        obj.mouseIn(); };

    obj.s.onmouseup =    function () { obj.mouseUp(); };
    obj.s.onmousedown =   function () { obj.mouseDown(); };

    obj.s.onmouseout  =   function () {
        obj.s.style.cursor = "default";
        obj.mouseOut(); };
}

function SVGWebcamButton(id, img, size) {

    // class constants and setup
    namespace = "http://www.inkscape.org/namespaces/inkscape";
    colorIn =    "gray";
    colorPress = "gray";
    colorOut =  "#9cb8f5";
    this.classMap = {};
    this.Id   = id;
    this.size = size;
    this.Image = img;

    // search for nodes
    this.s = document.querySelector(this.Id).contentDocument.getElementById("svg2");
    this.es = document.querySelector(this.Id).contentDocument.getElementsByTagName("g");
    for (i=0;i < this.es.length; i++) {
        if ( this.es[i] != null && this.es[i].hasAttributeNS(namespace,'label')) {
            var a = this.es[i].getAttributeNS(namespace,'label');
            this.classMap[a] = this.es[i];
            //console.log("SVG> " + a);
        }
    }

    this.isOn =   this.toggle = function (e) {
        if (this.classMap["on"].style.display == "none") {
            return false;
        } else {
            return true;
        }
    };

    this.toggle = function (e) {
        if (this.classMap["on"].style.display == "none") {
            this.setOn();
        }
        else {
            this.setOff();
        }
    };

    this.setOn = function (e) {
        this.classMap["on"].style.display = "inherit";
    };

    this.setOff = function (e) {
        this.classMap["on"].style.display = "none";
    };

    this.mouseOut = function (e) {
        //this.TextElem.style.fill = colorOut;
    };
    this.mouseIn = function (e) {
        //this.TextElem.style.fill = colorIn;
    };
    this.mouseDown = function (e) {
        //this.Cloud_Expanded.style.display = "inline";
    };
    this.mouseUp = function (e) {
        //this.Cloud_Expanded.style.display = "none";
    };

    setDefaultActions(this);

}


function SVGImageButton(id, img, size) {

    // class constants and setup
    namespace = "http://www.inkscape.org/namespaces/inkscape";
    colorIn =    "gray";
    colorPress = "gray";
    colorOut =  "#9cb8f5";
    this.classMap = {};
    this.Id   = id;
    this.size = size;
    this.Image = img;

    googleIcon = "googleIcon";
    facebookIcon = "facebookIcon";

    // search for nodes
    this.s = document.querySelector(this.Id).contentDocument.getElementById("svg2");
    this.es = document.querySelector(this.Id).contentDocument.getElementsByTagName("g");
    for (i=0;i < this.es.length; i++) {
        if ( this.es[i] != null && this.es[i].hasAttributeNS(namespace,'label')) {
            var a = this.es[i].getAttributeNS(namespace,'label');
            this.classMap[a] = this.es[i];
           // console.log("SVG> " + a);
        }
    }

    if (this.Image == "google") {
        this.ImageNode = this.classMap["googleIcon"];
    } else if (this.Image == "facebook") {
        this.ImageNode = this.classMap["facebookIcon"];
    }

    this.ImageNode.style.display = "inline";

    if (this.size == "S") {
        this.Cloud = this.classMap["S_Cloud"];
      //  this.Cloud_Text = this.classMap["S_Text"];
        this.Cloud_Expanded = this.classMap["S_Cloud-Expanded"]
    }
    if (this.size == "L") {
        this.Cloud = this.classMap["L_Cloud"];
     //   this.Cloud_Text = this.classMap["L_Text"];
        this.Cloud_Expanded = this.classMap["L_Cloud-Expanded"]
    }
    if (this.size == "XL") {
        this.Cloud = this.classMap["XL_Cloud"];
     //   this.Cloud_Text = this.classMap["XL_Text"];
        this.Cloud_Expanded = this.classMap["XL_Cloud-Expanded"]
    }

    this.Cloud.style.display = "inline";
   // this.Cloud_Text.style.display = "inline";
   // this.TextElem = this.Cloud_Text.children[0].children[0];
   // this.TextElem.textContent = this.text;


    this.mouseOut = function (e) {
      //  this.TextElem.style.fill = colorOut;
    };
    this.mouseIn = function (e) {
      //  this.TextElem.style.fill = colorIn;
    };
    this.mouseDown = function (e) {
        this.Cloud_Expanded.style.display = "inline";
    };
    this.mouseUp = function (e) {
        this.Cloud_Expanded.style.display = "none";
    };

    setDefaultActions(this);

}


function SVGButton(id, text, size) {

    // class constants and setup
    namespace = "http://www.inkscape.org/namespaces/inkscape";
    colorIn =    "gray";
    colorPress = "gray";
    colorOut =  "#9cb8f5";
    this.classMap = {};
    this.Id   = id;
    this.size = size;
    this.text = text;

    this.s = document.querySelector(this.Id).contentDocument.getElementById("svg2");
    this.es = document.querySelector(this.Id).contentDocument.getElementsByTagName("g");
    for (i=0;i < this.es.length; i++) {
        if ( this.es[i] != null && this.es[i].hasAttributeNS(namespace,'label')) {
            var a = this.es[i].getAttributeNS(namespace,'label');
            this.classMap[a] = this.es[i];
            //console.log("SVG2> " + a);
        }
    }

    if (this.size == "S") {
        this.Cloud = this.classMap["S_Cloud"];
        this.Cloud_Text = this.classMap["S_Text"];
        this.Cloud_Expanded = this.classMap["S_Cloud-Expanded"]
    }
    if (this.size == "L") {
        this.Cloud = this.classMap["L_Cloud"];
        this.Cloud_Text = this.classMap["L_Text"];
        this.Cloud_Expanded = this.classMap["L_Cloud-Expanded"]
    }
    if (this.size == "XL") {
        this.Cloud = this.classMap["XL_Cloud"];
        this.Cloud_Text = this.classMap["XL_Text"];
        this.Cloud_Expanded = this.classMap["XL_Cloud-Expanded"]
    }

    this.Cloud.style.display = "inline";
    this.Cloud_Text.style.display = "inline";
    this.TextElem = this.Cloud_Text.children[0].children[0];
    this.TextElem.textContent = this.text;

    this.mouseOut = function (e) {
        this.Cloud_Expanded.style.display = "none";

    };
    this.mouseIn = function (e) {
        this.Cloud_Expanded.style.display = "inline";
    };
    this.mouseDown = function (e) {
        this.TextElem.style.fill = colorIn;
    };
    this.mouseUp = function (e) {
        this.TextElem.style.fill = colorOut;
    };

    setDefaultActions(this);
}



