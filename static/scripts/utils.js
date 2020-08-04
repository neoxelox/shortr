'use strict';

var copyToClipboard = function (copyStr, origin) {
    const el = document.createElement('textarea'); // Create a <textarea> element
    el.value = copyStr;  // Set its value to the string that you want copied
    el.setAttribute('readonly', ''); // Make it readonly to be tamper-proof
    el.style.position = 'absolute'; 
    el.style.left = '-9999px'; // Move outside the screen to make it invisible
    document.body.appendChild(el); // Append the <textarea> element to the HTML document
    el.select(); // Select the <textarea> content
    document.execCommand('copy'); // Copy - only works as a result of a user action (e.g. click events)
    document.body.removeChild(el); // Remove the <textarea> element

    var originStr = origin.innerText
    origin.innerText = "Copied!";
    setTimeout(function() {
      origin.classList.add('fade-animation');
      setTimeout(() => origin.innerText = originStr, 500);
      "webkitAnimationEnd mozAnimationEnd MSAnimationEnd oanimationend animationend".split(" ")
      .map(name => origin.addEventListener(name, function() {
        origin.classList.remove('fade-animation');
      }, false));
    }, 500);
};
