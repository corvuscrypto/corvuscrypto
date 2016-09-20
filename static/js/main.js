{{ define "js" }}

//firstly put in the request animation shim

//BEGIN rAF SHIM

// http://paulirish.com/2011/requestanimationframe-for-smart-animating/
// http://my.opera.com/emoller/blog/2011/12/20/requestanimationframe-for-smart-er-animating

// requestAnimationFrame polyfill by Erik Möller. fixes from Paul Irish and Tino Zijdel

// MIT license

(function() {
    var lastTime = 0;
    var vendors = ['ms', 'moz', 'webkit', 'o'];
    for(var x = 0; x < vendors.length && !window.requestAnimationFrame; ++x) {
        window.requestAnimationFrame = window[vendors[x]+'RequestAnimationFrame'];
        window.cancelAnimationFrame = window[vendors[x]+'CancelAnimationFrame']
                                   || window[vendors[x]+'CancelRequestAnimationFrame'];
    }

    if (!window.requestAnimationFrame)
        window.requestAnimationFrame = function(callback, element) {
            var currTime = new Date().getTime();
            var timeToCall = Math.max(0, 16 - (currTime - lastTime));
            var id = window.setTimeout(function() { callback(currTime + timeToCall); },
              timeToCall);
            lastTime = currTime + timeToCall;
            return id;
        };

    if (!window.cancelAnimationFrame)
        window.cancelAnimationFrame = function(id) {
            clearTimeout(id);
        };
}());

//END rAF SHIM

var shellState = {
  prefix: {
    "user": "Guest",
    "domain": "corvuscrypto.com",
    "wd": "~"
  },
  cursorPos: 0,
  currentInput: "",
  history: [],
  historyPos: -1,
  keysPressed: []
}

document.body.onclick = function(){
  document.getElementById("hidden_input").focus();
}

function handleSpecialKeys(e){
  switch (e.keyCode) {
    case 37: //left arrow
      //move cursor to the left
      if (shellState.cursorPos > 0){
        shellState.cursorPos--;
      }
      break;
    case 39: //right arrow
      //move cursor to the right
      inputLength = shellState.currentInput.length
      if (shellState.cursorPos <= inputLength-1){
        shellState.cursorPos++;
      }
      break;
    case 8: //backspace
      if (shellState.cursorPos > 0){
        var input = shellState.currentInput;
        shellState.cursorPos--;
        shellState.currentInput = input.slice(0,shellState.cursorPos);
        shellState.currentInput += input.slice(shellState.cursorPos+1);
      }
      break;
    case 46: //delete
      //splice out the current input
      var input = shellState.currentInput;
      shellState.currentInput = input.slice(0,shellState.cursorPos);
      shellState.currentInput += input.slice(shellState.cursorPos+1);
      break;
    }
}
function _handleKeyPress(e){
  var inputElem = document.getElementById("hidden_input")
  var stringChar = inputElem.value[0]||"";
  if (stringChar){
    shellState.cursorPos++;
  }
  var input = shellState.currentInput
  if (stringChar === " " && (input[input.length-1] === " "||input[input.length-1] === undefined)) {
    stringChar = String.fromCharCode(160);
  }
  shellState.currentInput = input.slice(0,shellState.cursorPos-1);
  shellState.currentInput += stringChar
  shellState.currentInput += input.slice(shellState.cursorPos-1);
  inputElem.value = "";
}

function handleKeyPress(e){
  setTimeout(function(){_handleKeyPress(e)}, 10)
}

document.getElementById("hidden_input").onkeypress = handleKeyPress
document.getElementById("hidden_input").onkeydown = handleSpecialKeys

//MAIN RENDER LOOP
function renderInputArea(){

  var html = shellState.prefix.user+"@";
  html += shellState.prefix.domain+":";
  html += shellState.prefix.wd+"$ ";

  var input = shellState.currentInput;

  //replace the cursorPos's character with the cursor
  //representation

  //first get the html for the cursor
  var character = input[shellState.cursorPos]||"";
  var cursorHTML = "<span id=\"cursor\">"+character+"</span>";

  //now splice out that char and replace with the cursor
  inputHTML = input.slice(0,shellState.cursorPos);
  inputHTML += cursorHTML+input.slice(shellState.cursorPos+1);

  //add to the html and render!
  html += inputHTML;
  document.getElementById("stdin__text").innerHTML = html;
  window.requestAnimationFrame(renderInputArea)
}

window.requestAnimationFrame(renderInputArea)


{{end}}
