/**
 * Created by malin on 2017-03-09.
 */
//
// (C) 2017 Yamato Digital Audio
// Author: Malin af Lääkkö
//

function setDefaultActions (obj) {
    obj.s.onmouseover =  function () { obj.mouseIn(); };
    obj.s.onmouseup =    function () { obj.mouseUp(); };
    obj.s.onmousedown    =   function () { obj.mouseDown(); };
    obj.s.onmouseout =   function () { obj.mouseOut(); };
}

function SVGImageButton(id, img, size) {

    colorIn = "gray";
    colorPress = "gray";
    colorOut = "#9cb8f5";
    bigCloudExp = "no";
    smallCloudExp = "layer7";
    bigCloud = "no";
    smallCloud = "layer2";

    this.Id = id;
    this.size = size;
    this.image = img;

    console.log("ID ", this.Id);
    this.s = document.querySelector(this.Id).getSVGDocument().getElementById("svg2");

    document.querySelector(this.Id).getSVGDocument().getElementById("googleIcon").style.display = "none";
    document.querySelector(this.Id).getSVGDocument().getElementById("facebookIcon").style.display = "none";
    document.querySelector(this.Id).getSVGDocument().getElementById(smallCloudExp).style.display = "none";
    document.querySelector(this.Id).getSVGDocument().getElementById(this.image).style.display = "inline";

    this.big = function () {
        this.cloud = bigCloud;
        this.cloudE = bigCloudExp;
        document.querySelector(this.Id).getSVGDocument().getElementById(bigCloud).setAttribute("visibility", "visible");
        document.querySelector(this.Id).getSVGDocument().getElementById(smallCloud).setAttribute("visibility", "hidden");
    };
    this.small = function () {
        this.cloud = smallCloud;
        this.cloudE = smallCloudExp;
      //  document.querySelector(this.Id).getSVGDocument().getElementById(bigCloud).setAttribute("visibility", "hidden");
        document.querySelector(this.Id).getSVGDocument().getElementById(smallCloud).setAttribute("visibility", "visible");
    };
    if (this.size == "big") {
        this.small();
    } else {
        this.small();
    }
    this.mouseOut = function (e) {
        document.querySelector(this.Id).getSVGDocument().getElementById(this.cloudE).style.display = "none";;
    };
    this.mouseIn = function (e) {
        document.querySelector(this.Id).getSVGDocument().getElementById(this.cloudE).style.display = "inline";
    };
    this.mouseDown = function () {
        document.querySelector(this.Id).getSVGDocument().getElementById(this.cloudE).style.display = "inline";
    };
    this.mouseUp = function () {
        document.querySelector(this.Id).getSVGDocument().getElementById(this.cloudE).style.display = "none";
    };

    setDefaultActions(this);
}


function SVGButton(id, text, size) {

    this.Label = "text16838";
    colorIn = "gray";
    colorPress = "gray";
    colorOut = "#9cb8f5";
    label = "layer3";
    bigCloudExp = "layer8";
    smallCloudExp = "layer7";
    bigCloud = "layer5";
    smallCloud = "layer4";
    this.Id = id;
    this.size = size;

    this.s = document.querySelector(this.Id).getSVGDocument().getElementById("svg2");
    this.t = document.querySelector(this.Id).getSVGDocument().getElementById(this.Label);
    this.t.textContent = text;
    document.querySelector(this.Id).getSVGDocument().getElementById(label).style.display = "inline";
    document.querySelector(this.Id).getSVGDocument().getElementById(bigCloudExp).style.display = "none";
    document.querySelector(this.Id).getSVGDocument().getElementById(smallCloudExp).style.display = "none"

    this.big = function () {
        this.cloud = bigCloud;
        this.cloudE = bigCloudExp;
        document.querySelector(this.Id).getSVGDocument().getElementById(bigCloud).setAttribute("visibility", "visible");
        document.querySelector(this.Id).getSVGDocument().getElementById(smallCloud).setAttribute("visibility", "hidden");
    };
    this.small = function () {
        this.cloud = smallCloud;
        this.cloudE = smallCloudExp;
        document.querySelector(this.Id).getSVGDocument().getElementById(bigCloud).setAttribute("visibility", "hidden");
        document.querySelector(this.Id).getSVGDocument().getElementById(smallCloud).setAttribute("visibility", "visible");
    };
    if (this.size == "big") {
        this.big();
    } else {
        this.small();
    }
    this.mouseOut = function (e) {
        this.t.style.fill = colorOut;
    };
    this.mouseIn = function (e) {
        this.t.style.fill = colorIn;
    };
    this.mouseDown = function () {
        document.querySelector(this.Id).getSVGDocument().getElementById(this.cloudE).style.display = "inline";
    };
    this.mouseUp = function () {
        document.querySelector(this.Id).getSVGDocument().getElementById(this.cloudE).style.display = "none";
    };

    setDefaultActions(this);

}
