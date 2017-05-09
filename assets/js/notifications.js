
function openDrop() {
    document.getElementById("dropeddown").classList.toggle("show");
}

function clearNotif() {
  $.post("/NOTIFICATION/api/clear", function( data ) {
  });
    document.getElementById("dropeddown").classList.toggle("show");
    $("div").remove(".dropdowntop");
}

window.onclick = function(event) {
  if (!event.target.matches('.dropbutton')) {

    var dropdowns = document.getElementsByClassName("dropdown-contenttop");
    var i;
    for (i = 0; i < dropdowns.length; i++) {
      var openDropdown = dropdowns[i];
      if (openDropdown.classList.contains('show')) {
        openDropdown.classList.remove('show');
      }
    }
  }
}
