/*
 * Copyright (c) 2017 by Peter de Vocht
 *
 * All rights reserved. No part of this publication may be reproduced, distributed, or
 * transmitted in any form or by any means, including photocopying, recording, or other
 * electronic or mechanical methods, without the prior written permission of the publisher,
 * except in the case of brief quotations embodied in critical reviews and certain other
 * noncommercial uses permitted by copyright law.
 *
 */

var pContainerHeight = $('.hero').height();


$(window).scroll(function(){
  var wScroll = $(this).scrollTop();

  if (wScroll <= pContainerHeight) {
    $('.hero-right-bg').css({
      'transform' : 'translate(0px, -'+ wScroll /30 +'%)'
    });
    $('.laptop').css({
      'transform' : 'translate(0px, '+ wScroll /40 +'%)'
    });
  }
});

$('a').click(function() {
    $('html, body').delay(100).animate({
        scrollTop: $( $(this).attr('href') ).offset().top}, 1500);
    return false;
});

var scrollTo = $('.scroll-to');
    $(window).on('scroll', function() {
        var st = $(this).scrollTop();
        scrollTo.css({ 'opacity' : (1 - st/200) 
    });
});

