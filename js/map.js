//
// Copyright 2017 Malin Yamato Lääkkö --  All rights reserved.
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

/**
 * Created by malin on 2017-03-10.
 */
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



function showPosition(position) {
    var latlon = position.coords.latitude + "," + position.coords.longitude;
    console.log("latlon " + latlon);
    setMapPoint(position.coords.latitude, position.coords.longitude, $("#FirstName").val() + " " + $("#LastName").val());
    var url = "https://maps.googleapis.com/maps/api/geocode/json?latlng=" + latlon;
    $.get(url, function (data, status) {

            var len = data.results[0].address_components.length;
            var comps = data.results[0].address_components;
            for (i = 0; i < len; i++) {
                for (g = 0; g < comps[i].types.length; g++) {
                    if (comps[i].types[g] == "postal_town") {
                        var town = comps[i].long_name;
                        $("#Town").val(town)
                    } else if (comps[i].types[g] == "country") {
                        var country = comps[i].long_name;
                        $("#Country").val(country);
                    }
                }
            }
            console.log("City " + city + " Country " + country);
        }
    );
}

function getLocation() {
    $("#MapHolder").show();
    console.log("getLocation");
    var timeoutVal = 10 * 1000 * 1000;
    if (navigator.geolocation) {
        navigator.geolocation.getCurrentPosition(showPosition, displayError,
            {enableHighAccuracy: true, timeout: timeoutVal, maximumAge: 0}
        );
    } else {
        console.log("bad");
        $("#country").val("Unknown");
        $("#town").val("Unknown");
    }
}

function displayError(err) {
    console.log("Code " + err.code + " Err " + err.message);
    timeoutID = window.setTimeout(hideStatus, 2000);
    $("#status").css("background-color", "darkred");
    $("#status").text("Goolge is bad! " + err.message + " Code: " + err.code);
    $("#MapHolder").height(280);
    $("#StatusBox").show();
}

var pictureOverlay;
function setMapPoint(lat, long, message) {
    $("#StatusBox").hide();
    $("#MapHolder").height(350);
    var latlong = new google.maps.LatLng(lat, long);
    var mapCanvas = document.getElementById("MapHolder");
    var mapOptions = {center: latlong, zoom: 3};
    var map = new google.maps.Map(mapCanvas, mapOptions);

    var myCity = new google.maps.Circle({
        center: latlong,
        radius: 200000,
        strokeColor: "#0000FF",
        strokeOpacity: 0.8,
        strokeWeight: 2,
        fillColor: "#0000FF",
        fillOpacity: 0.1
    });
 //   myCity.setMap(map);
    var lo = document.getElementById("Longitude");
    lo.setAttribute("value",long);
    var la = document.getElementById("Latitude");
    la.setAttribute("value",lat);

    var imageBounds = {
        north: lat + 2,
        south: lat - 2,
        east: long + 4.5,
        west: long - 4.5
    };
    var picture = document.getElementById("PictureURL").getAttribute("src");
    pictureOverlay = new google.maps.GroundOverlay(
        picture,
        imageBounds);

    alat = 50.110924;
    along = 8.682127;
    var imageBounds2 = {
        north: alat + 4,
        south: alat - 4,
        east: along + 5,
        west: along - 5
    };
    pictureOverlay2 = new google.maps.GroundOverlay(
        "/images/krypinlogo.png",
        imageBounds2);

  //  myCity.setMap(map);
    pictureOverlay2.setMap(map);
    pictureOverlay.setMap(map);

}

function myMap() {
    alat = 50.110924;
    along = 8.682127;
    $("#StatusBox").hide();
    $("#MapHolder").height(350);
    var latlong = new google.maps.LatLng(alat, along);
    var mapCanvas = document.getElementById("MapHolder");
    var mapOptions = {center: latlong, zoom: 3};
    var map = new google.maps.Map(mapCanvas, mapOptions);
    var myCity = new google.maps.Circle({
        center: latlong,
        radius: 500000,
        strokeColor: "#0000FF",
        strokeOpacity: 0.8,
        strokeWeight: 2,
        fillColor: "#0000FF",
        fillOpacity: 0.01
    });
   // myCity.setMap(map);
    var imageBounds = {
        north: alat + 4,
        south: alat - 4,
        east: along + 5,
        west: along - 5
    };
    pictureOverlay = new google.maps.GroundOverlay(
        "/images/krypinlogo.png",
        imageBounds);
    pictureOverlay.setMap(map);
}


