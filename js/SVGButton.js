/**
 * Created by malin on 2017-03-09.
 */

function SVGButton(id, text, size, inFn, downFn, upFn, outFn) {

    this.Label = "text16838";
    colorIn = "gray";
    colorPress = "gray";
    colorOut = "#9cb8f5";
    chat = "layer2";
    photos = "layer3";
    profile = "layer6";
    bigCloudExp = "layer8";
    smallCloudExp = "layer7";
    bigCloud = "layer5";
    smallCloud = "layer4";
    this.Id = id;
    this.size = size;

    this.s = document.querySelector(id).getSVGDocument().getElementById("svg2");
    this.t = document.querySelector(id).getSVGDocument().getElementById(this.Label);
    this.t.textContent = text;
    document.querySelector(this.Id).getSVGDocument().getElementById(chat).style.display = "none";
    document.querySelector(this.Id).getSVGDocument().getElementById(photos).style.display = "inline";
    document.querySelector(this.Id).getSVGDocument().getElementById(profile).style.display = "none";
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
    this.t.onmouseover = inFn;
    this.t.onmousedown = downFn;
    this.t.onmouseup = upFn;
    this.t.onmouseout = outFn;
}
