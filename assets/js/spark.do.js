var droplets;

$(function() {
    $(document).on("click", "#digitalOceanConnect", function() {
        //window.open('https://cloud.digitalocean.com/v1/oauth/authorize?client_id=4fa3f972373332b978e2e7db91cf492fd0ca8015dd92829a47b822541e1504ee&redirect_uri=http%3A%2F%2Flocalhost%3A8080%2Flogin&response_type=code', '_blank', 'location=yes,height=600,width=900,scrollbars=yes,status=yes');
        window.location.replace('https://cloud.digitalocean.com/v1/oauth/authorize?client_id=4fa3f972373332b978e2e7db91cf492fd0ca8015dd92829a47b822541e1504ee&redirect_uri=http%3A%2F%2Flocalhost%3A8080%2Flogin&response_type=code');
    });

    $(document).on("click", "#menu-spark", function() {
        if (!$(this).hasClass("active")) {
            $("#menu-droplets").removeClass("active");
            $(this).addClass("active");
            showSpark();
        }
    });

    $(document).on("click", "#menu-droplets", function() {
        if (!$(this).hasClass("active")) {
            $("#menu-spark").removeClass("active");
            $(this).addClass("active");
            showDroplets();
        }
    });

    $(document).on("click", "#saveSpark", function() {
        // save access token in cookie
        $.post("spark", {
            "accesstoken": $('#accessToken').val(),
            "deviceid": $('#deviceID').val(),
        }, function() {
            var modal = $('#notificationModal');
            modal.find('.header').text('Success');
            modal.find('.content').text('Spark information successfully saved');
            modal.modal('show');
        }).fail(function() {
            var modal = $('#notificationModal');
            modal.find('.header').text('Error');
            modal.find('.content').text('Saving Spark information failed');
            modal.modal('show');
        });
    });

    $(document).on("click", ".monitor", function() {
        var that = this;
        $.post("monitor", {
            "dropletid": $(that).data('id'),
        }, function() {
            var modal = $('#notificationModal');
            modal.find('.header').text('Success');
            modal.find('.content').text('Droplet monitoring successfully started');
            modal.modal('show');
        }).fail(function() {
            var modal = $('#notificationModal');
            modal.find('.header').text('Error');
            modal.find('.content').text('Starting Droplet monitoring failed');
            modal.modal('show');
        });
    });

    // $(document).on("click", "#digitalOceanConnect", function() {
    //     $.post("accesstoken", {
    //         "accesstoken": SC.accessToken()
    //     }, function() {
    //         console.log("Successfully sent access token to back-end.");
    //     }).fail(function() {
    //         console.log("Error sending access token to back-end.");
    //     });
    // });

    // get access token, if user already logged in
    // if (!$("#digitalOceanConnect").length) {
    //     $.get("accesstoken", function(data) {
    //         console.log("Successfully got access token from back-end.");
    //         SC.accessToken(data.accesstoken);
    //         soundcloudManager.initializeContent();
    //     }).fail(function() {
    //         console.log("Error getting access token from back-end.");
    //     });
    // }

    // initialize carousel if it exists
    initializeCarousel();
});

function initializeCarousel(callback) {
    if ($(".owl-carousel").length) {

        var owl = $(".owl-carousel");

        if (callback) {
            owl.on('initialized.owl.carousel', function(event) {
                callback();
            });
        }

        owl.owlCarousel({
            loop: true,
            autoplay: true,
            items: 1,
            nav: true,
            smartSpeed: 1000,
        });

        $(document.documentElement).keyup(function(event) {
            if (event.keyCode == 37) {
                owl.trigger('prev.owl');
            } else if (event.keyCode == 39) {
                owl.trigger('next.owl');
            }
        });

        owl.on('mousewheel', '.owl-stage', function(e) {
            if (e.deltaY > 0) {
                owl.trigger('next.owl');
            } else {
                owl.trigger('prev.owl');
            }
            e.preventDefault();
        });
    }
}

function showSpark() {
    droplets = $('#userContent').html();
    var show = function() {
        $('#userContent').html('' +
            '<div class="sixteen wide column">\
                <div class="ui form segment">\
                    <div class="field">\
                        <label>Device ID</label>\
                        <div class="ui left labeled icon input">\
                            <input type="text" placeholder="Device ID" id="deviceID">\
                            <i class="mobile icon"></i>\
                            <div class="ui corner label">\
                                <i class="icon asterisk"></i>\
                            </div>\
                        </div>\
                    </div>\
                    <div class="field">\
                        <label>Access Token</label>\
                        <div class="ui left labeled icon input">\
                            <input type="text" placeholder="Access Token" id="accessToken">\
                            <i class="lock icon"></i>\
                            <div class="ui corner label">\
                                <i class="icon asterisk"></i>\
                            </div>\
                        </div>\
                    </div>\
                    <div class="ui blue submit button" id="saveSpark">Save</div>\
                </div>\
            </div>');
        $('#userContent').transition('slide down');
    };
    if ($("#userContent").hasClass("hidden")) {
        show();
    } else {
        $('#userContent').transition('slide down', function() {
            show();
        });
    }
}

function showDroplets() {
    var show = function() {
        if (droplets !== undefined) {
            $('#userContent').html(droplets);
        } else {
            $('#userContent').html('');
        }
        $('#userContent').transition('slide down');
    };
    if ($("#userContent").hasClass("hidden")) {
        show();
    } else {
        $('#userContent').transition('slide down', function() {
            show();
        });
    }
}
