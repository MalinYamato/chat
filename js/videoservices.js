//
//  Video services for krypin. This file implements the siganling protool for WebRTC.
//  Cortesy and admiration of the development team of Janus at Meetecho for making
//  a multiuser stable signlig porotocol possible for rtp in general and WebRTC in particular.
//  The team: Amirante, A. and Castaldi, T. and Miniero, L. and Romano, S. P. of Chicaco.
//
//  I have made some minor changes to this code to integrate it with krypin.xyz
//  Malin Lääkkö, Yamato Digital Audio - designer of Krypin
//
//
var server = null;


//if (window.location.protocol === 'http:')
///    server = "http://" + window.location.hostname + ":8088/janus";
//else
server = "https://" + "media.rakuen.asia" + ":8089/janus";

var janus = null;
var sfutest = null;
var opaqueId = "videoroomtest-" + Janus.randomString(12);

var started = false;

var myusername = null;
var myid = null;
var mystream = null;

// We use this other ID just to map our subscriptions to us
var mypvtid = null;

var feeds = [];
var bitrateTimer = [];

function register(name) {
   var register = {"request": "join", "room": 1234, "ptype": "publisher", "display": name};
   sfutest.send({"message": register});
}

function leaving(name) {
    var register = {"request": "leaving", "room": 1234, "ptype": "publisher", "display": name};
    sfutest.send({"message": register});
}

function destroyVideo() {
    janus.destroy();
}

function startVideo() {

    $(document).ready(function () {
        // Initialize the library (all console debuggers enabled)
        Janus.init({
            debug: "all", callback: function () {

                started = true;
                $(this).attr('disabled', true).unbind('click');

                if (!Janus.isWebrtcSupported()) {
                    console.log("No WebRTC support... ");
                    return;
                }
                janus = new Janus(
                    {
                        server: server,
                        success: function () {
                            // Attach to video room test plugin
                            janus.attach(
                                {
                                    plugin: "janus.plugin.videoroom",
                                    opaqueId: opaqueId,
                                    success: function (pluginHandle) {
                                        $('#details').remove();

                                        sfutest = pluginHandle;
                                        console.log("Plugin attached! (" + sfutest.getPlugin() + ", id=" + sfutest.getId() + ")");
                                        console.log("  -- This is a publisher/manager");

                                        $('#start').removeAttr('disabled').html("Stop")
                                            .click(function () {
                                                $(this).attr('disabled', true);
                                                janus.destroy();
                                            });
                                    },
                                    error: function (error) {
                                        console.error("  -- Error attaching plugin...", error);
                                        console.log("Error attaching plugin... " + error);
                                    },
                                    consentDialog: function (on) {
                                        Janus.debug("Consent dialog should be " + (on ? "on" : "off") + " now");
                                    },
                                    mediaState: function (medium, on) {
                                        console.log("Janus " + (on ? "started" : "stopped") + " receiving our " + medium);
                                    },
                                    webrtcState: function (on) {
                                        console.log("Janus says our WebRTC PeerConnection is " + (on ? "up" : "down") + " now");
                                        $("#videolocal").parent().parent().unblock();
                                    },
                                    onmessage: function (msg, jsep) {
                                        Janus.debug(" ::: Got a message (publisher) :::");
                                        Janus.debug(JSON.stringify(msg));
                                        var event = msg["videoroom"];
                                        Janus.debug("Event: " + event);
                                        if (event != undefined && event != null) {
                                            if (event === "joined") {
                                                                         // Publisher/manager created, negotiate WebRTC and attach to existing feeds, if any
                                                myid = msg["id"];
                                                mypvtid = msg["private_id"];
                                                Janus.log("Successfully joined room " + msg["room"] + " with ID " + myid);
                                                publishOwnFeed(true);

                                                // Any new feed to attach to?
                                                if (msg["publishers"] !== undefined && msg["publishers"] !== null) {
                                                    var list = msg["publishers"];
                                                    Janus.debug("Got a list of available publishers/feeds:");
                                                    Janus.debug(list);
                                                    for (var f in list) {
                                                        var id = list[f]["id"];
                                                        var display = list[f]["display"];
                                                        Janus.debug("  >> [" + id + "] " + display);
                                                        newRemoteFeed(id, display)
                                                    }
                                                }
                                            } else if (event === "destroyed") {
// The room has been destroyed
                                                Janus.warn("The room has been destroyed!");
                                                bootbox.log("The room has been destroyed", function () {
                                                    window.location.reload();
                                                });
                                            } else if (event === "event") {

                                                                   // Any new feed to attach to?
                                                if (msg["publishers"] !== undefined && msg["publishers"] !== null) {
                                                    var list = msg["publishers"];
                                                    Janus.debug("Got a list of available publishers/feeds:");
                                                    Janus.debug(list);
                                                    for (var f in list) {
                                                        var id = list[f]["id"];
                                                        var display = list[f]["display"];
                                                        Janus.debug("  >> [" + id + "] " + display);
                                                        newRemoteFeed(id, display)
                                                    }
                                                } else if (msg["leaving"] !== undefined && msg["leaving"] !== null) {
// One of the publishers has gone away?
                                                    var leaving = msg["leaving"];
                                                    Janus.log("Publisher left: " + leaving);
                                                    var remoteFeed = null;
                                                    for (var i = 1; i < 3; i++) {
                                                        if (feeds[i] != null && feeds[i] != undefined && feeds[i].rfid == leaving) {
                                                            remoteFeed = feeds[i];
                                                            break;
                                                        }
                                                    }
                                                    if (remoteFeed != null) {
                                                        Janus.debug("Feed " + remoteFeed.rfid + " (" + remoteFeed.rfdisplay + ") has left the room, detaching");
                                                        $('#remote' + remoteFeed.rfindex).empty().hide();
                                                        $('#videoremote' + remoteFeed.rfindex).empty();
                                                        feeds[remoteFeed.rfindex] = null;
                                                        remoteFeed.detach();
                                                    }
                                                } else if (msg["unpublished"] !== undefined && msg["unpublished"] !== null) {
// One of the publishers has unpublished?
                                                    var unpublished = msg["unpublished"];
                                                    Janus.log("Publisher left: " + unpublished);
                                                    if (unpublished === 'ok') {
// That's us
                                                        sfutest.hangup();      //////////////////////
                                                        return;
                                                    }
                                                    var remoteFeed = null;
                                                    for (var i = 1; i < 3; i++) {
                                                        if (feeds[i] != null && feeds[i] != undefined && feeds[i].rfid == unpublished) {
                                                            remoteFeed = feeds[i];
                                                            break;
                                                        }
                                                    }
                                                    if (remoteFeed != null) {
                                                        console.debug("Feed " + remoteFeed.rfid + " (" + remoteFeed.rfdisplay + ") has left the room, detaching");
                                                        $('#remote' + remoteFeed.rfindex).empty().hide();
                                                        $('#videoremote' + remoteFeed.rfindex).empty();
                                                        feeds[remoteFeed.rfindex] = null;
                                                        remoteFeed.detach();
                                                    }
                                                } else if (msg["error"] !== undefined && msg["error"] !== null) {
                                                    bootbox.alert(msg["error"]);
                                                }
                                            }
                                        }
                                        if (jsep !== undefined && jsep !== null) {
                                            Janus.debug("Handling SDP as well...");
                                            Janus.debug(jsep);
                                            sfutest.handleRemoteJsep({jsep: jsep});     ///////////////////////////////////
                                        }
                                    },
                                    onlocalstream: function (stream) {
                                        console.debug(" ::: Got a local stream :::");
                                        mystream = stream;
                                        console.debug(JSON.stringify(stream));
                                        $('#videolocal').empty();
                                        $('#videojoin').hide();
                                        $('#videos').removeClass('hide').show();
                                        if ($('#myvideo').length === 0) {
                                            $('#videolocal').append('<video class="rounded centered" id="myvideo" width="100%" height="100%" autoplay muted="unmuted"/>');
// Add a 'mute' button
                                           // $('#videolocal').append('<button class="btn btn-warning btn-xs" id="mute" style=" position: relative; bottom: 0px; left: 0px; margin: 15px;">Mute</button>');
                                           // $('#mute').click(toggleMute);
// Add an 'unpublish' button
                                             //$('#videolocal').append('<button class="btn btn-warning btn-xs" id="unpublish" style="position: relative; bottom: 0px; right: 0px; margin: 15px;">Unpublish</button>');
                                             //$('#unpublish').click(unpublishOwnFeed);
                                        }
                                        $('#publisher').removeClass('hide').html(myusername).show();
                                        Janus.attachMediaStream($('#myvideo').get(0), stream);
                                        $("#myvideo").get(0).muted = "muted";
                                        $("#videolocal").parent().parent().block({
                                            message: '<b>Publishing...</b>',
                                            css: {
                                                border: 'none',
                                                backgroundColor: 'transparent',
                                                color: 'white'
                                            }
                                        });
                                        var videoTracks = stream.getVideoTracks();
                                        if (videoTracks === null || videoTracks === undefined || videoTracks.length === 0) {
// No webcam
                                            $('#myvideo').hide();
                                            $('#videolocal').append(
                                                '<div class="no-video-container">' +
                                                '<i class="fa fa-video-camera fa-5 no-video-icon" style="height: 100%;"></i>' +
                                                '<span class="no-video-text" style="font-size: 16px;">No webcam available</span>' +
                                                '</div>');
                                        }
                                    },
                                    onremotestream: function (stream) {
// The publisher stream is sendonly, we don't expect anything here
                                    },
                                    oncleanup: function () {
                                        Janus.log(" ::: Got a cleanup notification: we are unpublished now :::");
                                        mystream = null;
                                        //$('#videolocal').html('<button id="publish" class="btn btn-primary">Publish</button>');
                                        //$('#publish').click(function () {
                                          //  publishOwnFeed(true);
                                        //});
                                        $("#videolocal").parent().parent().unblock();
                                    }
                                });
                        },
                        error: function (error) {
                            Janus.error(error);
                            bootbox.alert(error, function () {
                                window.location.reload();
                            });
                        },
                        destroyed: function () {
                            window.location.reload();
                        }
                    });

            }
        });
    });

}

function publishOwnFeed(useAudio) {    ////////////////////////
// Publish our stream
    useAudio = true;

    $('#publish').attr('disabled', true).unbind('click');
    sfutest.createOffer(
        {
// Add data:true here if you want to publish datachannels as well
            media: {audioRecv: false, videoRecv: false, audioSend: useAudio, videoSend: true},	// Publishers are sendonly
            success: function (jsep) {
                Janus.debug("Got publisher SDP!");
                Janus.debug(jsep);
                var publish = {"request": "configure", "audio": useAudio, "video": true};
                sfutest.send({"message": publish, "jsep": jsep});
            },
            error: function (error) {
                Janus.error("WebRTC error:", error);
                if (useAudio) {
                    publishOwnFeed(false);
                } else {
                    bootbox.alert("WebRTC error... " + JSON.stringify(error));
                    $('#publish').removeAttr('disabled').click(function () {
                        publishOwnFeed(true);
                    });
                }
            }
        });
}

function toggleMute() {          ////////////////////////////////////////
    var muted = sfutest.isAudioMuted();
    Janus.log((muted ? "Unmuting" : "Muting") + " local stream...");
    if (muted)
        sfutest.unmuteAudio();
    else
        sfutest.muteAudio();
    muted = sfutest.isAudioMuted();
    $('#mute').html(muted ? "Unmute" : "Mute");
}

function unpublishOwnFeed() {     ////////////////////////////////////////
// Unpublish our stream
    $('#unpublish').attr('disabled', true).unbind('click');
    var unpublish = {"request": "unpublish"};
    sfutest.send({"message": unpublish});
}

function newRemoteFeed(id, display) {
// A new feed has been published, create a new plugin handle and attach to it as a listener
    var remoteFeed = null;
    janus.attach(
        {
            plugin: "janus.plugin.videoroom",
            opaqueId: opaqueId,
            success: function (pluginHandle) {
                remoteFeed = pluginHandle;
                Janus.log("Plugin attached! (" + remoteFeed.getPlugin() + ", id=" + remoteFeed.getId() + ")");
                Janus.log("  -- This is a subscriber");
// We wait for the plugin to send us an offer
                var listen = {"request": "join", "room": 1234, "ptype": "listener", "feed": id, "private_id": mypvtid};
                remoteFeed.send({"message": listen});
            },
            error: function (error) {
                Janus.error("  -- Error attaching plugin...", error);
                bootbox.alert("Error attaching plugin... " + error);
            },
            onmessage: function (msg, jsep) {
                Janus.debug(" ::: Got a message (listener) :::");
                Janus.debug(JSON.stringify(msg));
                var event = msg["videoroom"];
                Janus.debug("Event: " + event);
                if (event != undefined && event != null) {
                    if (event === "attached") {
// Subscriber created and attached
                        for (var i = 1; i < 3; i++) {
                            if (feeds[i] === undefined || feeds[i] === null) {
                                feeds[i] = remoteFeed;
                                remoteFeed.rfindex = i;
                                break;
                            }
                        }
                        remoteFeed.rfid = msg["id"];
                        remoteFeed.rfdisplay = msg["display"];
                        if (remoteFeed.spinner === undefined || remoteFeed.spinner === null) {
                            var target = document.getElementById('videoremote' + remoteFeed.rfindex);
                            remoteFeed.spinner = new Spinner({top: 100}).spin(target);
                        } else {
                            remoteFeed.spinner.spin();
                        }
                        Janus.log("Successfully attached to feed " + remoteFeed.rfid + " (" + remoteFeed.rfdisplay + ") in room " + msg["room"]);
                        $('#remote' + remoteFeed.rfindex).removeClass('hide').html(remoteFeed.rfdisplay).show();
                    } else if (msg["error"] !== undefined && msg["error"] !== null) {
                        bootbox.alert(msg["error"]);
                    } else {
                        // What has just happened?
                    }
                }
                if (jsep !== undefined && jsep !== null) {
                    Janus.debug("Handling SDP as well...");
                    Janus.debug(jsep);
// Answer and attach
                    remoteFeed.createAnswer(
                        {
                            jsep: jsep,
// Add data:true here if you want to subscribe to datachannels as well
// (obviously only works if the publisher offered them in the first place)
                            media: {audioSend: false, videoSend: false},	// We want recvonly audio/video
                            success: function (jsep) {
                                Janus.debug("Got SDP!");
                                Janus.debug(jsep);
                                var body = {"request": "start", "room": 1234};
                                remoteFeed.send({"message": body, "jsep": jsep});
                            },
                            error: function (error) {
                                Janus.error("WebRTC error:", error);
                                bootbox.alert("WebRTC error... " + JSON.stringify(error));
                            }
                        });
                }
            },
            webrtcState: function (on) {
                Janus.log("Janus says this WebRTC PeerConnection (feed #" + remoteFeed.rfindex + ") is " + (on ? "up" : "down") + " now");
            },
            onlocalstream: function (stream) {
// The subscriber stream is recvonly, we don't expect anything here
            },
            onremotestream: function (stream) {
                Janus.debug("Remote feed #" + remoteFeed.rfindex);
                if ($('#remotevideo' + remoteFeed.rfindex).length === 0) {
// No remote video yet
                    $('#videoremote' + remoteFeed.rfindex).append('<video class="rounded centered" id="waitingvideo' + remoteFeed.rfindex + '" width=320 height=240 />');
                    $('#videoremote' + remoteFeed.rfindex).append('<video class="rounded centered relative hide" id="remotevideo' + remoteFeed.rfindex + '" width="100%" height="100%" autoplay/>');
                }
                $('#videoremote' + remoteFeed.rfindex).append(
                    '<span class="label label-primary hide" id="curres' + remoteFeed.rfindex + '" style="position: absolute; bottom: -80px; left: 0px; margin: 15px;"></span>' +
                    '<span class="label label-info hide" id="curbitrate' + remoteFeed.rfindex + '" style="position: absolute; bottom: -80px; right: 0px; margin: 15px;"></span>');
// Show the video, hide the spinner and show the resolution when we get a playing event
                $("#remotevideo" + remoteFeed.rfindex).bind("playing", function () {
                    if (remoteFeed.spinner !== undefined && remoteFeed.spinner !== null)
                        remoteFeed.spinner.stop();
                    remoteFeed.spinner = null;
                    $('#waitingvideo' + remoteFeed.rfindex).remove();
                    $('#remotevideo' + remoteFeed.rfindex).removeClass('hide');
                    var width = this.videoWidth;
                    var height = this.videoHeight;
                    $('#curres' + remoteFeed.rfindex).removeClass('hide').text(width + 'x' + height).show();
                    if (adapter.browserDetails.browser === "firefox") {
// Firefox Stable has a bug: width and height are not immediately available after a playing
                        setTimeout(function () {
                            var width = $("#remotevideo" + remoteFeed.rfindex).get(0).videoWidth;
                            var height = $("#remotevideo" + remoteFeed.rfindex).get(0).videoHeight;
                            $('#curres' + remoteFeed.rfindex).removeClass('hide').text(width + 'x' + height).show();
                        }, 2000);
                    }
                });
                Janus.attachMediaStream($('#remotevideo' + remoteFeed.rfindex).get(0), stream);
                var videoTracks = stream.getVideoTracks();
                if (videoTracks === null || videoTracks === undefined || videoTracks.length === 0 || videoTracks[0].muted) {
// No remote video
                    $('#remotevideo' + remoteFeed.rfindex).hide();
                    $('#videoremote' + remoteFeed.rfindex).append(
                        '<div class="no-video-container">' +
                        '<i class="fa fa-video-camera fa-5 no-video-icon" style="height: 100%;"></i>' +
                        '<span class="no-video-text" style="font-size: 16px;">No remote video available</span>' +
                        '</div>');
                }
                if (adapter.browserDetails.browser === "chrome" || adapter.browserDetails.browser === "firefox") {
                    $('#curbitrate' + remoteFeed.rfindex).removeClass('hide').show();
                    bitrateTimer[remoteFeed.rfindex] = setInterval(function () {
// Display updated bitrate, if supported
                        var bitrate = remoteFeed.getBitrate();
                        $('#curbitrate' + remoteFeed.rfindex).text(bitrate);
                    }, 1000);
                }
            },
            oncleanup: function () {
                Janus.log(" ::: Got a cleanup notification (remote feed " + id + ") :::");
                if (remoteFeed.spinner !== undefined && remoteFeed.spinner !== null)
                    remoteFeed.spinner.stop();
                remoteFeed.spinner = null;
                $('#waitingvideo' + remoteFeed.rfindex).remove();
                $('#curbitrate' + remoteFeed.rfindex).remove();
                $('#curres' + remoteFeed.rfindex).remove();
                if (bitrateTimer[remoteFeed.rfindex] !== null && bitrateTimer[remoteFeed.rfindex] !== null)
                    clearInterval(bitrateTimer[remoteFeed.rfindex]);
                bitrateTimer[remoteFeed.rfindex] = null;
            }
        });
}

function monika() {
    document.getElementById("username").textContent = "eva";
    var register = {"request": "join", "room": 1234, "ptype": "publisher", "display": "Monika"};
    sfutest.send({"message": register});
}
