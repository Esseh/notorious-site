
function openDrop() {
    document.getElementById("dropeddown").classList.toggle("show");
}

function clearNotif() {
  $.post("/NOTIFICATION/api/clear", function( data ) {
  });
    document.getElementById("dropeddown").classList.toggle("show");
    document.getElementById("newbutton").innerHTML = "No New Notifications";
    document.getElementById("newbutton").removeAttribute("onclick");
    document.getElementById("newbutton").style.backgroundColor = "gray";
}

window.onclick = function(event) {
  if (!event.target.matches('.dropbtn')) {

    var dropdowns = document.getElementsByClassName("dropdown-content");
    var i;
    for (i = 0; i < dropdowns.length; i++) {
      var openDropdown = dropdowns[i];
      if (openDropdown.classList.contains('show')) {
        openDropdown.classList.remove('show');
      }
    }
  }
}
