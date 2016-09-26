function showSubmitButton() {
  var btn = this.nextElementSibling;
  btn.style.display = "inline";
}
function hideSubmitButton(e) {
  var btn = this.nextElementSibling;
  btn.style.display = "none";
}
//add the listener for the focus and blur on the search_inputs
var elems = document.getElementsByClassName("nav-search__input");
for(var i=0;i < elems.length; i++) {
  elems[i].addEventListener("focus", showSubmitButton)
  if (elems[i].nextElementSibling){
      elems[i].nextElementSibling.addEventListener("mousedown", function(e){
        e.preventDefault();
        this.parentElement.submit()
      })
  }
  elems[i].addEventListener("blur", hideSubmitButton)
}
